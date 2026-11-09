// Package httpapi is the JSON HTTP API under /api/v1, built on chi.
package httpapi

import (
	"context"
	"log/slog"
	"mime"
	"net/http"
	"slices"
	"strings"
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
	"github.com/aryangarg/codequest/api/internal/tutor"
)

// Deps are the handler dependencies.
type Deps struct {
	Pool           *pgxpool.Pool
	Grader         *grading.Service
	Log            *slog.Logger
	AllowedOrigins []string
	// Sessions resolve the session cookie. Required.
	Sessions *auth.Sessions
	// CLI issues login codes and resolves the CLI's bearer tokens
	// (&auth.CLI{Pool: Pool} when nil).
	CLI *auth.CLI
	// GitHub and Google are nil when that sign-in is not configured.
	GitHub *auth.GitHub
	Google *auth.Google
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
	// Tutor is the AI tutor's model; nil turns the tutor off (always in
	// production for now). TutorQuota is its daily hint limit; required with Tutor.
	Tutor      tutor.Provider
	TutorQuota *tutor.Quota
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
	if d.CLI == nil {
		d.CLI = &auth.CLI{Pool: d.Pool}
	}
	a := &api{d}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(a.logRequests)
	r.Use(a.recoverer)
	r.Use(a.guardStateChanges)
	r.Use(a.sessionUser)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not_found", "no such endpoint")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	})

	r.Route("/api/v1", func(r chi.Router) {
		// The tutor streams for longer than the usual timeout; it sets its own.
		r.With(a.requireUnlocked).Post("/exercises/{exerciseID}/hint", a.hint)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Timeout(30 * time.Second))
			r.Get("/healthz", a.healthz)
			r.Get("/me", a.me)
			r.Patch("/me", a.updateMe)
			r.Get("/home", a.home)
			r.Get("/profile", a.profile)
			r.Get("/leagues/current", a.currentLeague)

			r.Get("/auth/providers", a.providers)
			r.Get("/auth/{provider}/start", a.oauthStart)
			r.With(a.limitSignIns).Get("/auth/{provider}/callback", a.oauthCallback)
			r.Post("/auth/logout", a.logout)
			r.Post("/cli/codes", a.createCLICode)
			r.With(a.limitSignIns).Post("/cli/token", a.cliToken)
			r.Post("/cli/logout", a.cliLogout)
			r.Get("/cli/next", a.cliNext)
			if d.DevAuth {
				r.Post("/auth/dev", a.devSignIn)
			}
			r.Get("/courses", a.listCourses)
			r.Get("/courses/{courseID}", a.getCourse)
			r.Get("/assets/{file}", a.getAsset)
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
			r.With(a.requireUnlocked).Get("/exercises/{exerciseID}/tutor", a.tutorThread)

			r.With(a.limitTelemetry).Post("/events", a.events)
			r.With(a.limitTelemetry).Post("/client-errors", a.clientError)
		})
	})
	return r
}

type ctxKey int

const (
	userKey ctxKey = iota
	bearerKey
)

// currentUser returns the signed-in user, if any.
func currentUser(r *http.Request) uuid.NullUUID {
	u, _ := r.Context().Value(userKey).(uuid.NullUUID)
	return u
}

// bearerToken returns the request's `Authorization: Bearer` token, if any.
func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", false
	}
	scheme, tok, _ := strings.Cut(h, " ")
	return strings.TrimSpace(tok), strings.EqualFold(scheme, "Bearer")
}

// viaCLI reports whether the request was authenticated by a CLI bearer token.
func viaCLI(r *http.Request) bool {
	v, _ := r.Context().Value(bearerKey).(bool)
	return v
}

// sessionUser resolves the session cookie, or the CLI's bearer token, to the
// signed-in user. A request may carry one credential, not both.
func (a *api) sessionUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, cookieErr := r.Cookie(sessionCookie)
		if tok, ok := bearerToken(r); ok || r.Header.Get("Authorization") != "" {
			if !ok {
				writeError(w, http.StatusUnauthorized, "unauthorized", "use Authorization: Bearer <token>")
				return
			}
			if cookieErr == nil {
				writeError(w, http.StatusBadRequest, "bad_request", "send a session cookie or a bearer token, not both")
				return
			}
			u, err := a.CLI.Lookup(r.Context(), tok)
			if err != nil {
				a.internal(w, r, err)
				return
			}
			if !u.Valid {
				writeError(w, http.StatusUnauthorized, "unauthorized", "the CLI token is invalid or expired; run `codequest login`")
				return
			}
			ctx := context.WithValue(r.Context(), userKey, u)
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, bearerKey, true)))
			return
		}
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
//
// Two exceptions skip the Origin check (never the content type): requests
// with an Authorization header (the CLI; a cross-site page cannot attach one
// without a CORS preflight, which this API never answers, and sessionUser
// refuses a request that also carries the cookie), and POST /cli/token, which
// sets no cookie and needs a one-time code shown only to the signed-in user.
func (a *api) guardStateChanges(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		cli := r.Header.Get("Authorization") != "" || r.URL.Path == "/api/v1/cli/token"
		if origin := r.Header.Get("Origin"); !cli && !slices.Contains(a.AllowedOrigins, origin) {
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

// limitTelemetry bounds analytics and error reports per user and IP.
func (a *api) limitTelemetry(next http.Handler) http.Handler {
	return a.limit("telemetry", 60, 240, "too many reports; wait a moment", next)
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
