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

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/economy"
	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/grading"
	"github.com/aryangarg/codequest/api/internal/ratelimit"
)

// Deps are the handler dependencies.
type Deps struct {
	Pool           *pgxpool.Pool
	Grader         *grading.Service
	Log            *slog.Logger
	AllowedOrigins []string
	// Sessions resolve the session cookie. Required.
	Sessions *auth.Sessions
	// Links sends and redeems magic links. Required.
	Links *auth.MagicLinks
	// GitHub is nil when GitHub sign-in is not configured.
	GitHub *auth.GitHub
	// DevAuth enables POST /auth/dev (sign in without email). Never in production.
	DevAuth bool
	// SecureCookies sets the Secure flag (when the app is served over https).
	SecureCookies bool
	// Board is the Redis league board. Required.
	Board *game.Board
	// Economy is the shop, quests, medals and chapter locks. Required.
	Economy *economy.Service
	// Now is the clock (time.Now when nil).
	Now func() time.Time
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
	r.Use(a.sessionUser)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not_found", "no such endpoint")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", a.healthz)
		r.Get("/me", a.me)
		r.Patch("/me", a.updateMe)
		r.Get("/home", a.home)
		r.Get("/leagues/current", a.currentLeague)

		r.Get("/auth/providers", a.providers)
		r.With(a.limitSignIns).Post("/auth/magic-link", a.sendMagicLink)
		r.With(a.limitSignIns).Post("/auth/magic-link/verify", a.verifyMagicLink)
		r.Get("/auth/github/start", a.githubStart)
		r.Get("/auth/github/callback", a.githubCallback)
		r.Post("/auth/logout", a.logout)
		if d.DevAuth {
			r.Post("/auth/dev", a.devSignIn)
		}
		r.Get("/courses", a.listCourses)
		r.Get("/courses/{courseID}", a.getCourse)
		r.With(a.requireUnlocked).Get("/exercises/{exerciseID}", a.getExercise)
		r.With(a.requireUnlocked, a.limitSubmits).Post("/exercises/{exerciseID}/submit", a.submit)
		r.With(a.requireUnlocked, a.limitRuns).Post("/exercises/{exerciseID}/run", a.run)

		r.Get("/shop", a.shop)
		r.Get("/inventory", a.shop)
		r.Post("/shop/{powerup}/buy", a.buy)
		r.Post("/powerups/{powerup}/use", a.usePowerup)
		r.Get("/quests", a.quests)
		r.Post("/quests/chest/open", a.openChest)
		r.Get("/medals", a.medals)
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

// sessionUser resolves the session cookie to the signed-in user.
func (a *api) sessionUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(sessionCookie); err == nil {
			u, err := a.Sessions.Lookup(r.Context(), c.Value)
			if err != nil {
				a.internal(w, r, err)
				return
			}
			if u.Valid {
				r = r.WithContext(context.WithValue(r.Context(), userKey, u))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (a *api) now() time.Time {
	if a.Now != nil {
		return a.Now()
	}
	return time.Now()
}

// guardStateChanges is the CSRF defence for non-GET requests: an allowed
// Origin header (browsers always send one on these methods) and a JSON content
// type (which a cross-site form cannot send without a CORS preflight). The
// session cookie is also SameSite=Lax.
func (a *api) guardStateChanges(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		if origin := r.Header.Get("Origin"); !slices.Contains(a.AllowedOrigins, origin) {
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
