package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/store"
)

// analyticsEvents are the event names the browser may send. Anything else is
// rejected, so the table only ever holds events someone decided to collect.
var analyticsEvents = map[string]bool{
	"page_view":            true,
	"landing_cta":          true,
	"sign_in":              true,
	"onboarding_step":      true,
	"onboarding_completed": true,
	"lesson_opened":        true,
	"run_clicked":          true,
	"submit_clicked":       true,
	"tutor_opened":         true,
	"tutor_asked":          true,
	"shop_bought":          true,
}

const (
	maxEventsPerBatch = 20
	maxPropsBytes     = 1024
)

var anonIDPattern = regexp.MustCompile(`^[A-Za-z0-9-]{8,64}$`)

type eventIn struct {
	Name  string          `json:"name"`
	Props json.RawMessage `json:"props"`
	Path  string          `json:"path"`
	At    *time.Time      `json:"at"`
}

type eventsRequest struct {
	AnonID string    `json:"anon_id"`
	Events []eventIn `json:"events"`
}

// events stores a batch of analytics events (self-hosted; no third party).
func (a *api) events(w http.ResponseWriter, r *http.Request) {
	var req eventsRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	if req.AnonID != "" && !anonIDPattern.MatchString(req.AnonID) {
		writeError(w, http.StatusBadRequest, "bad_request", "anon_id must be 8-64 letters, digits or dashes")
		return
	}
	if len(req.Events) == 0 || len(req.Events) > maxEventsPerBatch {
		writeError(w, http.StatusBadRequest, "bad_request", fmt.Sprintf("send 1-%d events", maxEventsPerBatch))
		return
	}
	now := a.now()
	rows := make([]store.InsertAnalyticsEventParams, 0, len(req.Events))
	for _, e := range req.Events {
		if !analyticsEvents[e.Name] {
			writeError(w, http.StatusBadRequest, "bad_request", fmt.Sprintf("unknown event %q", e.Name))
			return
		}
		props, err := flatProps(e.Props)
		if err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", fmt.Sprintf("event %s: %v", e.Name, err))
			return
		}
		at := now
		// Trust the browser's time only within the last day (batches are sent late, not forged).
		if e.At != nil && e.At.Before(now) && now.Sub(*e.At) < 24*time.Hour {
			at = *e.At
		}
		rows = append(rows, store.InsertAnalyticsEventParams{
			UserID: currentUser(r), AnonID: req.AnonID, Name: e.Name, Props: props,
			Path: truncate(e.Path, 300), CreatedAt: at,
		})
	}
	err := pgx.BeginFunc(r.Context(), a.Pool, func(tx pgx.Tx) error {
		q := store.New(tx)
		for _, row := range rows {
			if err := q.InsertAnalyticsEvent(r.Context(), row); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// flatProps accepts a small JSON object of string, number or boolean values.
func flatProps(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return json.RawMessage(`{}`), nil
	}
	if len(raw) > maxPropsBytes {
		return nil, fmt.Errorf("props larger than %d bytes", maxPropsBytes)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("props must be an object")
	}
	for k, v := range m {
		switch v.(type) {
		case string, float64, bool:
		default:
			return nil, fmt.Errorf("prop %q must be a string, number or boolean", k)
		}
	}
	return raw, nil
}

type clientErrorRequest struct {
	Message   string `json:"message"`
	Stack     string `json:"stack"`
	URL       string `json:"url"`
	UserAgent string `json:"user_agent"`
	Release   string `json:"release"`
}

// clientError records a browser error, grouped by fingerprint per day, and
// logs it so it shows next to server errors.
func (a *api) clientError(w http.ResponseWriter, r *http.Request) {
	var req clientErrorRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "message is required")
		return
	}
	msg, stack := truncate(req.Message, 1000), truncate(req.Stack, 8000)
	fp := errorFingerprint(msg, stack)
	err := store.New(a.Pool).UpsertClientError(r.Context(), store.UpsertClientErrorParams{
		Fingerprint: fp, Day: a.now().UTC().Truncate(24 * time.Hour), Message: msg, Stack: stack,
		Url: truncate(req.URL, 500), UserAgent: truncate(req.UserAgent, 300), Release: truncate(req.Release, 64),
		UserID: currentUser(r),
	})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	a.Log.Error("client error", "fingerprint", fp, "message", msg, "url", truncate(req.URL, 500), "release", req.Release)
	w.WriteHeader(http.StatusNoContent)
}

// errorFingerprint groups errors by message and top stack frame, with line and
// column numbers and hashed bundle names removed so a rebuild keeps the group.
var volatile = regexp.MustCompile(`:\d+(:\d+)?|-[A-Za-z0-9_]{8}\.js|\?[^\s)]*`)

func errorFingerprint(msg, stack string) string {
	top := ""
	for _, l := range strings.Split(stack, "\n") {
		if l = strings.TrimSpace(l); strings.HasPrefix(l, "at ") || strings.Contains(l, "@") {
			top = l
			break
		}
	}
	sum := sha256.Sum256([]byte(volatile.ReplaceAllString(msg, "") + "\n" + volatile.ReplaceAllString(top, "")))
	return hex.EncodeToString(sum[:8])
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	// Cut on a rune boundary.
	for n > 0 && n < len(s) && (s[n]&0xC0) == 0x80 {
		n--
	}
	return s[:n]
}

// recoverer turns a panic into a 500 and logs it with its stack and request id
// (error monitoring for the API).
func (a *api) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if rec == http.ErrAbortHandler { //nolint:errorlint // sentinel compared as net/http does
				panic(rec)
			}
			a.Log.Error("panic", "panic", fmt.Sprint(rec), "path", r.URL.Path,
				"req_id", middleware.GetReqID(r.Context()), "stack", string(debug.Stack()))
			writeError(w, http.StatusInternalServerError, "internal", "something went wrong")
		}()
		next.ServeHTTP(w, r)
	})
}
