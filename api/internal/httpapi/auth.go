package httpapi

import (
	"crypto/subtle"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/ratelimit"
	"github.com/aryangarg/codequest/api/internal/store"
)

const (
	sessionCookie = "cq_session"
	oauthCookie   = "cq_oauth"
)

func (a *api) setSession(w http.ResponseWriter, r *http.Request, user uuid.UUID) bool {
	raw, exp, err := a.Sessions.Create(r.Context(), user)
	if err != nil {
		a.internal(w, r, err)
		return false
	}
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: raw, Path: "/", Expires: exp, MaxAge: int(time.Until(exp).Seconds()),
		HttpOnly: true, Secure: a.SecureCookies, SameSite: http.SameSiteLaxMode,
	})
	return true
}

func (a *api) clearCookie(w http.ResponseWriter, name, path string) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: path, MaxAge: -1, HttpOnly: true, Secure: a.SecureCookies, SameSite: http.SameSiteLaxMode})
}

type providersResponse struct {
	GitHub bool `json:"github"`
	Google bool `json:"google"`
	Dev    bool `json:"dev"`
}

func (a *api) providers(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, providersResponse{GitHub: a.GitHub != nil, Google: a.Google != nil, Dev: a.DevAuth})
}

type signedInResponse struct {
	User userDTO `json:"user"`
	Next string  `json:"next"`
}

type devSignInRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

// devSignIn signs in without email. Registered only when DEV_AUTH=1.
func (a *api) devSignIn(w http.ResponseWriter, r *http.Request) {
	var req devSignInRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	starter := req.Email == ""
	if starter {
		req.Email, req.Name = "dev@codequest.local", "Dev User"
	}
	email, err := auth.NormalizeEmail(req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	// The default dev user starts with XP and Bytes (once) so power-ups can be
	// tried; dev sign-ins with an email (the e2e learners) start from zero.
	if starter {
		u, err := auth.UserForEmail(r.Context(), a.Pool, email, req.Name, "")
		if err == nil {
			err = a.Economy.GrantDevStarter(r.Context(), u.ID)
		}
		if err != nil {
			a.internal(w, r, err)
			return
		}
	}
	a.finishSignIn(w, r, email, req.Name, "", req.Timezone, "/")
}

// finishSignIn finds or creates the user, records their timezone on first
// sign-in, starts a session and answers with the user.
func (a *api) finishSignIn(w http.ResponseWriter, r *http.Request, email, name, avatar, tz, next string) {
	u, err := auth.UserForEmail(r.Context(), a.Pool, email, name, avatar)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	if validTimezone(tz) && u.Timezone == "UTC" {
		if err := store.New(a.Pool).SetUserTimezone(r.Context(), store.SetUserTimezoneParams{ID: u.ID, Timezone: tz}); err != nil {
			a.internal(w, r, err)
			return
		}
		u.Timezone = tz
	}
	if !a.setSession(w, r, u.ID) {
		return
	}
	writeJSON(w, http.StatusOK, signedInResponse{User: toUserDTO(u), Next: auth.SafeNext(next)})
}

func (a *api) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		if err := a.Sessions.Delete(r.Context(), c.Value); err != nil {
			a.internal(w, r, err)
			return
		}
	}
	a.clearCookie(w, sessionCookie, "/")
	w.WriteHeader(http.StatusNoContent)
}

// oauthProvider is the configured provider named in the URL, or nil.
func (a *api) oauthProvider(r *http.Request) auth.Provider {
	switch chi.URLParam(r, "provider") {
	case "github":
		if a.GitHub != nil {
			return a.GitHub
		}
	case "google":
		if a.Google != nil {
			return a.Google
		}
	}
	return nil
}

func oauthPath(p auth.Provider) string { return "/api/v1/auth/" + p.Name() }

// oauthStart redirects to the provider. The state is a random value kept in a
// short-lived HttpOnly cookie (with the post-sign-in path and the browser's
// timezone, saved on first sign-in) and checked on return.
func (a *api) oauthStart(w http.ResponseWriter, r *http.Request) {
	p := a.oauthProvider(r)
	if p == nil {
		writeError(w, http.StatusNotFound, "not_found", "this sign-in provider is not configured")
		return
	}
	state, _ := auth.NewToken()
	next := auth.SafeNext(r.URL.Query().Get("next"))
	tz := r.URL.Query().Get("tz")
	if !validTimezone(tz) {
		tz = ""
	}
	http.SetCookie(w, &http.Cookie{
		Name: oauthCookie, Value: state + "|" + url.QueryEscape(next) + "|" + url.QueryEscape(tz), Path: oauthPath(p), MaxAge: 600,
		HttpOnly: true, Secure: a.SecureCookies, SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, p.AuthCodeURL(state), http.StatusFound)
}

func (a *api) oauthCallback(w http.ResponseWriter, r *http.Request) {
	p := a.oauthProvider(r)
	if p == nil {
		writeError(w, http.StatusNotFound, "not_found", "this sign-in provider is not configured")
		return
	}
	fail := func(reason string, err error) {
		if err != nil {
			a.Log.Warn("oauth sign-in failed", "provider", p.Name(), "reason", reason, "err", err)
		}
		a.clearCookie(w, oauthCookie, oauthPath(p))
		http.Redirect(w, r, "/sign-in?error="+p.Name(), http.StatusFound)
	}
	var cookie string
	if c, err := r.Cookie(oauthCookie); err == nil {
		cookie = c.Value
	}
	state, rest, _ := strings.Cut(cookie, "|")
	nextEnc, tzEnc, _ := strings.Cut(rest, "|")
	q := r.URL.Query()
	if state == "" || subtle.ConstantTimeCompare([]byte(q.Get("state")), []byte(state)) != 1 {
		fail("state mismatch", nil)
		return
	}
	if q.Get("error") != "" || q.Get("code") == "" {
		fail("denied", nil)
		return
	}
	token, err := p.Exchange(r.Context(), q.Get("code"))
	if err != nil {
		fail("exchange", err)
		return
	}
	prof, err := p.Profile(r.Context(), token)
	if err != nil {
		fail("profile", err)
		return
	}
	user, err := a.oauthUser(r, p.Name(), prof)
	if err == nil {
		err = a.firstTimezone(r, user, tzEnc)
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	a.clearCookie(w, oauthCookie, oauthPath(p))
	if !a.setSession(w, r, user) {
		return
	}
	next, _ := url.QueryUnescape(nextEnc)
	http.Redirect(w, r, auth.SafeNext(next), http.StatusFound)
}

// firstTimezone records the browser's timezone for a user still on the default (UTC).
func (a *api) firstTimezone(r *http.Request, user uuid.UUID, tzEnc string) error {
	tz, _ := url.QueryUnescape(tzEnc)
	if !validTimezone(tz) {
		return nil
	}
	q := store.New(a.Pool)
	u, err := q.GetUser(r.Context(), user)
	if err != nil || u.Timezone != "UTC" {
		return err
	}
	return q.SetUserTimezone(r.Context(), store.SetUserTimezoneParams{ID: user, Timezone: tz})
}

// oauthUser returns the account linked to this provider identity, else the
// account with the same verified email (linking it), else a new account.
func (a *api) oauthUser(r *http.Request, provider string, p auth.Profile) (uuid.UUID, error) {
	var id uuid.UUID
	err := pgx.BeginFunc(r.Context(), a.Pool, func(tx pgx.Tx) error {
		q := store.New(tx)
		linked, err := q.GetOAuthUser(r.Context(), store.GetOAuthUserParams{Provider: provider, ProviderUserID: p.ID})
		if err == nil {
			id = linked
			return nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		email, err := auth.NormalizeEmail(p.Email)
		if err != nil {
			return err
		}
		u, err := auth.UserForEmail(r.Context(), tx, email, p.Name, p.Avatar)
		if err != nil {
			return err
		}
		id = u.ID
		return q.LinkOAuth(r.Context(), store.LinkOAuthParams{Provider: provider, ProviderUserID: p.ID, UserID: u.ID})
	})
	return id, err
}

// limitSignIns caps sign-in attempts per client IP (20 per hour).
func (a *api) limitSignIns(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		ok, retry, err := a.Limiter.Allow(r.Context(), ratelimit.Rule{Name: "signin-ip", Limit: 20, Window: time.Hour}, ip)
		if err != nil {
			a.Log.Error("rate limiter", "err", err)
			writeError(w, http.StatusServiceUnavailable, "unavailable", "try again shortly")
			return
		}
		if !ok {
			w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())+1))
			writeError(w, http.StatusTooManyRequests, "rate_limited", "too many sign-in attempts; try again later")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func validTimezone(tz string) bool {
	if tz == "" || len(tz) > 64 || tz == "Local" {
		return false
	}
	_, err := time.LoadLocation(tz)
	return err == nil
}
