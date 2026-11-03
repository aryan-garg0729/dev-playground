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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/ratelimit"
	"github.com/aryangarg/codequest/api/internal/store"
)

const (
	sessionCookie = "cq_session"
	oauthCookie   = "cq_oauth"
	githubPath    = "/api/v1/auth/github"
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
	MagicLink bool `json:"magic_link"`
	GitHub    bool `json:"github"`
	Dev       bool `json:"dev"`
}

func (a *api) providers(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, providersResponse{MagicLink: true, GitHub: a.GitHub != nil, Dev: a.DevAuth})
}

type magicLinkRequest struct {
	Email string `json:"email"`
	Next  string `json:"next"`
}

// sendMagicLink always answers the same way for a valid address, whether or
// not an account exists (no account enumeration).
func (a *api) sendMagicLink(w http.ResponseWriter, r *http.Request) {
	var req magicLinkRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	email, err := auth.NormalizeEmail(req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	ok, retry, err := a.Limiter.Allow(r.Context(), ratelimit.Rule{Name: "magic-email", Limit: 5, Window: time.Hour}, email)
	if err != nil {
		a.Log.Error("rate limiter", "err", err)
		writeError(w, http.StatusServiceUnavailable, "unavailable", "try again shortly")
		return
	}
	if !ok {
		w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())+1))
		writeError(w, http.StatusTooManyRequests, "rate_limited", "too many sign-in emails; try again later")
		return
	}
	if err := a.Links.Send(r.Context(), email, req.Next); err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "sent"})
}

type verifyRequest struct {
	Token    string `json:"token"`
	Timezone string `json:"timezone"`
}

type signedInResponse struct {
	User userDTO `json:"user"`
	Next string  `json:"next"`
}

func (a *api) verifyMagicLink(w http.ResponseWriter, r *http.Request) {
	var req verifyRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	email, next, err := a.Links.Redeem(r.Context(), req.Token)
	if errors.Is(err, auth.ErrInvalidLink) {
		writeError(w, http.StatusBadRequest, "invalid_link", err.Error())
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	a.finishSignIn(w, r, email, "", "", req.Timezone, next)
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
	if req.Email == "" {
		req.Email, req.Name = "dev@codequest.local", "Dev User"
	}
	email, err := auth.NormalizeEmail(req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
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

// githubStart redirects to GitHub. The state is a random value kept in a
// short-lived HttpOnly cookie (with the post-sign-in path) and checked on return.
func (a *api) githubStart(w http.ResponseWriter, r *http.Request) {
	if a.GitHub == nil {
		writeError(w, http.StatusNotFound, "not_found", "GitHub sign-in is not configured")
		return
	}
	state, _ := auth.NewToken()
	next := auth.SafeNext(r.URL.Query().Get("next"))
	http.SetCookie(w, &http.Cookie{
		Name: oauthCookie, Value: state + "|" + url.QueryEscape(next), Path: githubPath, MaxAge: 600,
		HttpOnly: true, Secure: a.SecureCookies, SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, a.GitHub.AuthCodeURL(state), http.StatusFound)
}

func (a *api) githubCallback(w http.ResponseWriter, r *http.Request) {
	fail := func(reason string, err error) {
		if err != nil {
			a.Log.Warn("github sign-in failed", "reason", reason, "err", err)
		}
		a.clearCookie(w, oauthCookie, githubPath)
		http.Redirect(w, r, "/sign-in?error=github", http.StatusFound)
	}
	if a.GitHub == nil {
		writeError(w, http.StatusNotFound, "not_found", "GitHub sign-in is not configured")
		return
	}
	var cookie string
	if c, err := r.Cookie(oauthCookie); err == nil {
		cookie = c.Value
	}
	state, nextEnc, _ := strings.Cut(cookie, "|")
	q := r.URL.Query()
	if state == "" || subtle.ConstantTimeCompare([]byte(q.Get("state")), []byte(state)) != 1 {
		fail("state mismatch", nil)
		return
	}
	if q.Get("error") != "" || q.Get("code") == "" {
		fail("denied", nil)
		return
	}
	token, err := a.GitHub.Exchange(r.Context(), q.Get("code"))
	if err != nil {
		fail("exchange", err)
		return
	}
	p, err := a.GitHub.Profile(r.Context(), token)
	if err != nil {
		fail("profile", err)
		return
	}
	user, err := a.githubUser(r, p)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	a.clearCookie(w, oauthCookie, githubPath)
	if !a.setSession(w, r, user) {
		return
	}
	next, _ := url.QueryUnescape(nextEnc)
	http.Redirect(w, r, auth.SafeNext(next), http.StatusFound)
}

// githubUser returns the account linked to this GitHub id, else the account
// with the same verified email (linking it), else a new account.
func (a *api) githubUser(r *http.Request, p auth.Profile) (uuid.UUID, error) {
	var id uuid.UUID
	err := pgx.BeginFunc(r.Context(), a.Pool, func(tx pgx.Tx) error {
		q := store.New(tx)
		linked, err := q.GetOAuthUser(r.Context(), store.GetOAuthUserParams{Provider: "github", ProviderUserID: p.ID})
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
		return q.LinkOAuth(r.Context(), store.LinkOAuthParams{Provider: "github", ProviderUserID: p.ID, UserID: u.ID})
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
