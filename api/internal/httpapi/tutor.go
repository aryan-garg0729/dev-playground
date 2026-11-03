package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/store"
	"github.com/aryangarg/codequest/api/internal/tutor"
)

const (
	maxTutorQuestion = 2000     // characters
	maxTutorCode     = 64 << 10 // bytes
	tutorTimeout     = 2 * time.Minute
)

type tutorMessageDTO struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type tutorThread struct {
	DailyLimit int               `json:"daily_limit"`
	Remaining  int               `json:"remaining"`
	Messages   []tutorMessageDTO `json:"messages"`
}

type hintRequest struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// tutorReady answers 401/404 when there is no user or no tutor; it reports
// whether the handler may go on.
func (a *api) tutorReady(w http.ResponseWriter, r *http.Request) bool {
	if !currentUser(r).Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in to ask the tutor")
		return false
	}
	if a.Tutor == nil || a.TutorQuota == nil {
		writeError(w, http.StatusNotFound, "tutor_unavailable", "the tutor is not available yet")
		return false
	}
	return true
}

// userToday returns the user's local date.
func (a *api) userToday(ctx context.Context, q *store.Queries, r *http.Request) (time.Time, error) {
	u, err := q.GetUser(ctx, currentUser(r).UUID)
	if err != nil {
		return time.Time{}, fmt.Errorf("user: %w", err)
	}
	return game.LocalDate(a.now(), game.Location(u.Timezone)), nil
}

// tutorThread returns the conversation for an exercise and today's hints left.
func (a *api) tutorThread(w http.ResponseWriter, r *http.Request) {
	if !a.tutorReady(w, r) {
		return
	}
	ctx := r.Context()
	q := store.New(a.Pool)
	user := currentUser(r).UUID
	id := chi.URLParam(r, "exerciseID")
	if _, err := q.GetExercise(ctx, id); errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "not_found", "exercise not found")
		return
	} else if err != nil {
		a.internal(w, r, err)
		return
	}
	today, err := a.userToday(ctx, q, r)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	used, err := a.TutorQuota.Used(ctx, user, today)
	if err != nil {
		a.Log.Error("tutor quota", "err", err)
		writeError(w, http.StatusServiceUnavailable, "unavailable", "try again shortly")
		return
	}
	rows, err := q.ListTutorMessages(ctx, store.ListTutorMessagesParams{UserID: user, ExerciseID: id, Lim: 50})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	out := tutorThread{DailyLimit: a.TutorQuota.Limit, Remaining: max(0, a.TutorQuota.Limit-used), Messages: []tutorMessageDTO{}}
	for _, m := range slices.Backward(rows) {
		out.Messages = append(out.Messages, tutorMessageDTO{Role: m.Role, Content: m.Content, CreatedAt: m.CreatedAt})
	}
	writeJSON(w, http.StatusOK, out)
}

// hint asks the tutor and streams the reply as Server-Sent Events:
// `delta` {text}, then `done` {remaining} or `error` {code, message}.
// Errors before the stream starts use the normal JSON envelope.
func (a *api) hint(w http.ResponseWriter, r *http.Request) {
	if !a.tutorReady(w, r) {
		return
	}
	var req hintRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	req.Message = strings.TrimSpace(req.Message)
	switch {
	case req.Message == "":
		writeError(w, http.StatusBadRequest, "bad_request", "message is required")
		return
	case utf8.RuneCountInString(req.Message) > maxTutorQuestion:
		writeError(w, http.StatusBadRequest, "bad_request", fmt.Sprintf("message is longer than %d characters", maxTutorQuestion))
		return
	case len(req.Code) > maxTutorCode:
		writeError(w, http.StatusBadRequest, "bad_request", "code is too long")
		return
	}

	askedAt := a.now()
	ctx, cancel := context.WithTimeout(r.Context(), tutorTimeout)
	defer cancel()
	q := store.New(a.Pool)
	user := currentUser(r).UUID
	ex, err := q.GetExercise(ctx, chi.URLParam(r, "exerciseID"))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "not_found", "exercise not found")
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	today, err := a.userToday(ctx, q, r)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	ok, remaining, err := a.TutorQuota.Take(ctx, user, today)
	if err != nil {
		a.Log.Error("tutor quota", "err", err)
		writeError(w, http.StatusServiceUnavailable, "unavailable", "try again shortly")
		return
	}
	if !ok {
		writeError(w, http.StatusTooManyRequests, "limit_reached",
			fmt.Sprintf("You've used today's %d tutor hints. They come back tomorrow.", a.TutorQuota.Limit))
		return
	}
	// From here on a failure before any text refunds the hint.
	refund := func() {
		if err := a.TutorQuota.Refund(context.WithoutCancel(ctx), user, today); err != nil {
			a.Log.Error("tutor refund", "err", err)
		}
	}

	hist, err := q.ListTutorMessages(ctx, store.ListTutorMessagesParams{UserID: user, ExerciseID: ex.ID, Lim: tutor.HistoryLen})
	if err != nil {
		refund()
		a.internal(w, r, err)
		return
	}
	last, err := q.GetLatestResult(ctx, store.GetLatestResultParams{UserID: user, ExerciseID: ex.ID})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		refund()
		a.internal(w, r, err)
		return
	}
	in := tutor.Exercise{Title: ex.Title, Type: ex.Type, Public: ex.PublicView, Code: req.Code, Result: last, Question: req.Message}
	for _, m := range slices.Backward(hist) {
		in.History = append(in.History, tutor.Message{Role: m.Role, Content: m.Content})
	}

	sse := newSSE(w)
	var reply strings.Builder
	usage, err := a.Tutor.Stream(ctx, tutor.BuildRequest(in), func(s string) error {
		reply.WriteString(s)
		return sse.send("delta", map[string]string{"text": s})
	})
	text := strings.TrimSpace(reply.String())
	reqID := middleware.GetReqID(r.Context())
	if err != nil && text == "" {
		refund()
		a.Log.Error("tutor failed", "err", err, "req_id", reqID)
		_ = sse.send("error", errorDetail{Code: "tutor_failed", Message: "The tutor could not answer. Your hint was not used; try again."})
		return
	}
	if err != nil { // a partial answer still counts; keep what was said
		a.Log.Warn("tutor stream cut short", "err", err, "req_id", reqID)
	}

	// Save even if the browser went away, so the hint is never lost. The reply
	// is always stamped after the question so the thread keeps its order.
	save := context.WithoutCancel(ctx)
	answeredAt := a.now()
	if !answeredAt.After(askedAt) {
		answeredAt = askedAt.Add(time.Microsecond)
	}
	err = pgx.BeginFunc(save, a.Pool, func(tx pgx.Tx) error {
		tq := store.New(tx)
		if err := tq.InsertTutorMessage(save, store.InsertTutorMessageParams{
			UserID: user, ExerciseID: ex.ID, Role: "user", Content: req.Message, Model: a.Tutor.Model(), CreatedAt: askedAt,
		}); err != nil {
			return err
		}
		if err := tq.InsertTutorMessage(save, store.InsertTutorMessageParams{
			UserID: user, ExerciseID: ex.ID, Role: "assistant", Content: text, Model: a.Tutor.Model(),
			InputTokens: int32(usage.InputTokens), OutputTokens: int32(usage.OutputTokens), CreatedAt: answeredAt,
		}); err != nil {
			return err
		}
		return tq.MarkHintUsed(save, store.MarkHintUsedParams{UserID: user, ExerciseID: ex.ID})
	})
	if err != nil {
		a.Log.Error("save tutor reply", "err", err, "req_id", reqID)
	}
	a.Log.Info("tutor", "user_id", user, "exercise_id", ex.ID, "model", a.Tutor.Model(),
		"input_tokens", usage.InputTokens, "output_tokens", usage.OutputTokens, "req_id", reqID)
	_ = sse.send("done", map[string]any{"remaining": remaining, "used_hint": true})
}

// sse writes Server-Sent Events, flushing after each one.
type sse struct {
	w  http.ResponseWriter
	rc *http.ResponseController
}

func newSSE(w http.ResponseWriter) *sse {
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-store")
	h.Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	return &sse{w: w, rc: http.NewResponseController(w)}
}

func (s *sse) send(event string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("sse: %w", err)
	}
	if _, err := fmt.Fprintf(s.w, "event: %s\ndata: %s\n\n", event, b); err != nil {
		return fmt.Errorf("sse: %w", err)
	}
	if err := s.rc.Flush(); err != nil {
		return fmt.Errorf("sse: %w", err)
	}
	return nil
}
