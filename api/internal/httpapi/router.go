// Package httpapi is the JSON HTTP API under /api/v1, built on chi.
package httpapi

import (
	"context"
	"log/slog"
	"mime"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/grading"
	"github.com/aryangarg/codequest/api/internal/ratelimit"
)

// Deps are the handler dependencies.
type Deps struct {
	Pool           *pgxpool.Pool
	Grader         *grading.Service
	Log            *slog.Logger
	AllowedOrigins []string
	// DevUser, when valid, is the user every request acts as (Phase 1, no auth).
	DevUser uuid.NullUUID
	// Limiter rate-limits submissions per user and per client IP. Required.
	Limiter *ratelimit.Limiter
	// SubmitPerUser and SubmitPerIP are per-minute limits (defaults 30 and 120).
	SubmitPerUser, SubmitPerIP int
}

type api struct{ Deps }

// New returns the API handler.
func New(d Deps) http.Handler {
	if d.Log == nil {
		d.Log = slog.Default()
	}
	if d.SubmitPerUser <= 0 {
		d.SubmitPerUser = 30
	}
	if d.SubmitPerIP <= 0 {
		d.SubmitPerIP = 120
	}
	a := &api{d}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(a.logRequests)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(a.guardStateChanges)
	r.Use(a.devUser)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not_found", "no such endpoint")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", a.healthz)
		r.Get("/me", a.me)
		r.Get("/courses", a.listCourses)
		r.Get("/courses/{courseID}", a.getCourse)
		r.Get("/exercises/{exerciseID}", a.getExercise)
		r.With(a.limitSubmits).Post("/exercises/{exerciseID}/submit", a.submit)
		r.With(a.limitRuns).Post("/exercises/{exerciseID}/run", a.run)
		r.Get("/submissions/{submissionID}", a.getSubmission)
	})
	return r
}

type ctxKey int

const userKey ctxKey = iota

// currentUser returns the signed-in user, if any.
func currentUser(r *http.Request) uuid.NullUUID {
	u, _ := r.Context().Value(userKey).(uuid.NullUUID)
	return u
}

func (a *api) devUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.DevUser.Valid {
			r = r.WithContext(context.WithValue(r.Context(), userKey, a.DevUser))
		}
		next.ServeHTTP(w, r)
	})
}

// guardStateChanges is the CSRF defence for non-GET requests: a JSON content
// type (which a cross-site form cannot send without a CORS preflight) and, when
// the browser sends one, an allowed Origin.
func (a *api) guardStateChanges(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		if origin := r.Header.Get("Origin"); origin != "" && !slices.Contains(a.AllowedOrigins, origin) {
			writeError(w, http.StatusForbidden, "forbidden_origin", "origin not allowed")
			return
		}
		mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mt != "application/json" {
			writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "send Content-Type: application/json")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *api) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		next.ServeHTTP(ww, r)
		a.Log.Info("http", "method", r.Method, "path", r.URL.Path, "status", ww.Status(),
			"ms", time.Since(start).Milliseconds(), "req_id", middleware.GetReqID(r.Context()))
	})
}

func (a *api) internal(w http.ResponseWriter, r *http.Request, err error) {
	a.Log.Error("internal error", "err", err, "path", r.URL.Path, "req_id", middleware.GetReqID(r.Context()))
	writeError(w, http.StatusInternalServerError, "internal", "something went wrong")
}

func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	if err := a.Pool.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "unavailable", "database unreachable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
