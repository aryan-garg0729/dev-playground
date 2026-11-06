package httpapi

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/grading"
	"github.com/aryangarg/codequest/api/internal/ratelimit"
	"github.com/aryangarg/codequest/api/internal/runner"
	"github.com/aryangarg/codequest/api/internal/store"
)

type chapterRef struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type position struct {
	Index int `json:"index"` // 1-based, within the chapter
	Total int `json:"total"`
}

type lastSubmission struct {
	SubmissionID    string          `json:"submission_id"`
	ExerciseVersion int             `json:"exercise_version"`
	Answer          json.RawMessage `json:"answer"`
	Result          json.RawMessage `json:"result"`
	CreatedAt       time.Time       `json:"created_at"`
}

type exerciseDetail struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Version    int             `json:"version"`
	Title      string          `json:"title"`
	XP         int             `json:"xp"`
	Difficulty int             `json:"difficulty"`
	Tags       []string        `json:"tags"`
	Free       bool            `json:"free"`
	Boss       bool            `json:"boss"`
	CourseID   string          `json:"course_id"`
	Chapter    chapterRef      `json:"chapter"`
	Position   position        `json:"position"`
	PrevID     *string         `json:"prev_id"`
	NextID     *string         `json:"next_id"`
	PublicView json.RawMessage `json:"public_view"`
	// Lesson images: name used in the content ("assets/loop.svg") -> URL.
	Assets map[string]string `json:"assets"`
	Last   *lastSubmission   `json:"last_submission"`
}

func (a *api) getExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := store.New(a.Pool)
	user := currentUser(r)
	e, err := q.GetExercise(ctx, chi.URLParam(r, "exerciseID"))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "not_found", "exercise not found")
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	// Phase 5: billing.HasAccess check here, 402 { code: "locked" }.

	siblings, err := q.ListCourseExercises(ctx, store.ListCourseExercisesParams{UserID: user, CourseID: e.CourseID})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	out := exerciseDetail{
		ID: e.ID, Type: e.Type, Version: int(e.Version), Title: e.Title, XP: int(e.Xp),
		Difficulty: int(e.Difficulty), Tags: e.Tags, Free: e.Free, Boss: e.Boss, CourseID: e.CourseID,
		Chapter: chapterRef{ID: e.ChapterID, Title: e.ChapterTitle}, PublicView: e.PublicView,
	}
	if out.Tags == nil {
		out.Tags = []string{}
	}
	if out.Assets, err = assetURLs(e.Assets); err != nil {
		a.internal(w, r, err)
		return
	}
	for i, s := range siblings {
		if s.ChapterID == e.ChapterID {
			out.Position.Total++
			if s.ID == e.ID {
				out.Position.Index = out.Position.Total
			}
		}
		if s.ID == e.ID {
			if i > 0 {
				out.PrevID = &siblings[i-1].ID
			}
			if i+1 < len(siblings) {
				out.NextID = &siblings[i+1].ID
			}
		}
	}

	if user.Valid {
		last, err := q.GetLatestSubmission(ctx, store.GetLatestSubmissionParams{UserID: user.UUID, ExerciseID: e.ID})
		switch {
		case errors.Is(err, pgx.ErrNoRows):
		case err != nil:
			a.internal(w, r, err)
			return
		default:
			out.Last = &lastSubmission{
				SubmissionID: last.ID.String(), ExerciseVersion: int(last.ExerciseVersion),
				Answer: last.Answer, Result: last.Result, CreatedAt: last.CreatedAt,
			}
		}
	}
	writeJSON(w, http.StatusOK, out)
}

type submitResponse struct {
	SubmissionID string          `json:"submission_id"`
	Result       exercise.Result `json:"result"`
	XPAwarded    int             `json:"xp_awarded"`
	Outcome      game.Outcome    `json:"outcome"`
}

// pendingResponse is the 202 body: grading continues; poll GET /submissions/{id}.
type pendingResponse struct {
	SubmissionID string `json:"submission_id"`
	Status       string `json:"status"` // "pending"
}

// readSubmission decodes and checks a submit/run body; it writes the error itself.
func readSubmission(w http.ResponseWriter, r *http.Request) (exercise.Submission, bool) {
	var sub exercise.Submission
	if err := decodeJSON(w, r, &sub); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return sub, false
	}
	if sub.ExerciseID != chi.URLParam(r, "exerciseID") {
		writeError(w, http.StatusBadRequest, "bad_request", "exercise_id does not match the URL")
		return sub, false
	}
	if len(sub.Answer) == 0 || string(sub.Answer) == "null" {
		writeError(w, http.StatusBadRequest, "bad_request", "answer is required")
		return sub, false
	}
	return sub, true
}

// run is the ungraded server Run (`run: server`): nothing is stored.
func (a *api) run(w http.ResponseWriter, r *http.Request) {
	if !currentUser(r).Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in to run code")
		return
	}
	sub, ok := readSubmission(w, r)
	if !ok {
		return
	}
	out, err := a.Grader.Run(r.Context(), sub)
	switch {
	case errors.Is(err, grading.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "exercise not found")
	case errors.Is(err, grading.ErrStaleVersion):
		writeError(w, http.StatusConflict, "stale_version", "this exercise changed; reload to get the latest version")
	case errors.Is(err, exercise.ErrInvalidAnswer):
		writeError(w, http.StatusBadRequest, "invalid_answer", err.Error())
	case errors.Is(err, exercise.ErrRunUnavailable):
		writeError(w, http.StatusBadRequest, "run_unavailable", err.Error())
	case errors.Is(err, runner.ErrUnavailable):
		a.Log.Error("runner unavailable", "err", err, "req_id", middleware.GetReqID(r.Context()))
		writeError(w, http.StatusServiceUnavailable, "grader_unavailable", "the runner is busy or offline; try again shortly")
	case err != nil:
		a.internal(w, r, err)
	default:
		writeJSON(w, http.StatusOK, out)
	}
}

func (a *api) submit(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if !user.Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in to submit")
		return
	}
	sub, ok := readSubmission(w, r)
	if !ok {
		return
	}
	g, err := a.Grader.Submit(r.Context(), user.UUID, sub)
	switch {
	case errors.Is(err, grading.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "exercise not found")
	case errors.Is(err, grading.ErrStaleVersion):
		writeError(w, http.StatusConflict, "stale_version", "this exercise changed; reload to get the latest version")
	case errors.Is(err, exercise.ErrInvalidAnswer):
		writeError(w, http.StatusBadRequest, "invalid_answer", err.Error())
	case errors.Is(err, runner.ErrUnavailable):
		a.Log.Error("grader unavailable", "err", err, "req_id", middleware.GetReqID(r.Context()))
		writeError(w, http.StatusServiceUnavailable, "grader_unavailable", "the grader is busy or offline; try again shortly")
	case err != nil:
		a.internal(w, r, err)
	case g.Pending:
		writeJSON(w, http.StatusAccepted, pendingResponse{SubmissionID: g.SubmissionID.String(), Status: "pending"})
	default:
		writeJSON(w, http.StatusOK, submitResponse{
			SubmissionID: g.SubmissionID.String(), Result: g.Result, XPAwarded: g.Outcome.XPAwarded, Outcome: g.Outcome,
		})
	}
}

type submissionResponse struct {
	SubmissionID    string          `json:"submission_id"`
	ExerciseID      string          `json:"exercise_id"`
	ExerciseVersion int             `json:"exercise_version"`
	Status          string          `json:"status"` // pending | pass | fail | error | timeout
	Result          json.RawMessage `json:"result"` // null while pending
	XPAwarded       int             `json:"xp_awarded"`
	Outcome         *game.Outcome   `json:"outcome"` // null while pending or after a grading error
	CreatedAt       time.Time       `json:"created_at"`
}

func (a *api) getSubmission(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if !user.Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in to see submissions")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "submissionID"))
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "submission not found")
		return
	}
	s, err := a.Grader.Get(r.Context(), user.UUID, id)
	switch {
	case errors.Is(err, grading.ErrNoSubmission):
		writeError(w, http.StatusNotFound, "not_found", "submission not found")
	case err != nil:
		a.internal(w, r, err)
	default:
		res := s.Result
		if len(res) == 0 {
			res = json.RawMessage("null")
		}
		var out *game.Outcome
		if len(s.Outcome) > 0 {
			out = &game.Outcome{}
			if err := json.Unmarshal(s.Outcome, out); err != nil {
				a.internal(w, r, err)
				return
			}
		}
		resp := submissionResponse{
			SubmissionID: s.ID.String(), ExerciseID: s.Exercise, ExerciseVersion: s.Version,
			Status: s.Status, Result: res, Outcome: out, CreatedAt: s.CreatedAt,
		}
		if out != nil {
			resp.XPAwarded = out.XPAwarded
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// limitSubmits applies the per-user and per-IP submit limits. If Redis is
// down, submits are refused rather than let through unlimited.
func (a *api) limitSubmits(next http.Handler) http.Handler {
	return a.limit("submit", a.SubmitPerUser, a.SubmitPerIP, "too many submissions; wait a moment", next)
}

// limitRuns applies the server Run limits (Run is ungraded but still uses a worker).
func (a *api) limitRuns(next http.Handler) http.Handler {
	return a.limit("run", 2*a.SubmitPerUser, 2*a.SubmitPerIP, "too many runs; wait a moment", next)
}

func (a *api) limit(name string, perUser, perIP int, msg string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		type limit struct {
			rule ratelimit.Rule
			key  string
		}
		checks := []limit{{ratelimit.Rule{Name: name + "-ip", Limit: perIP, Window: time.Minute}, ip}}
		if u := currentUser(r); u.Valid {
			checks = append(checks, limit{ratelimit.Rule{Name: name + "-user", Limit: perUser, Window: time.Minute}, u.UUID.String()})
		}
		for _, c := range checks {
			ok, retry, err := a.Limiter.Allow(r.Context(), c.rule, c.key)
			if err != nil {
				a.Log.Error("rate limiter", "err", err)
				writeError(w, http.StatusServiceUnavailable, "unavailable", "try again shortly")
				return
			}
			if !ok {
				w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())+1))
				writeError(w, http.StatusTooManyRequests, "rate_limited", msg)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
