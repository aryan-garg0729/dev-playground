package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/economy"
)

// doRaw sends a request and returns the recorder, for cookie and header checks.
func (e *env) doRaw(req *http.Request) *httptest.ResponseRecorder {
	e.t.Helper()
	rec := httptest.NewRecorder()
	e.h.ServeHTTP(rec, req)
	return rec
}

func jsonReq(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", origin)
	return req
}

func sessionFrom(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range rec.Result().Cookies() {
		if c.Name == "cq_session" && c.Value != "" {
			return c
		}
	}
	t.Fatalf("no session cookie in %v", rec.Header().Values("Set-Cookie"))
	return nil
}

func TestDevAuth(t *testing.T) {
	off := setup(t, false)
	if code, _ := off.do(http.MethodPost, "/api/v1/auth/dev", `{}`, nil); code != http.StatusNotFound {
		t.Fatalf("dev sign-in without DEV_AUTH = %d, want 404", code)
	}
	var p struct{ Dev bool }
	off.getJSON("/api/v1/auth/providers", &p)
	if p.Dev {
		t.Fatal("providers says dev is on")
	}

	on := setupWith(t, false, opts{devAuth: true})
	rec := on.doRaw(jsonReq(http.MethodPost, "/api/v1/auth/dev", `{}`))
	if rec.Code != http.StatusOK {
		t.Fatalf("dev sign-in = %d: %s", rec.Code, rec.Body)
	}
	on.cookie = sessionFrom(t, rec).Value
	var me struct{ User struct{ Email string } }
	on.getJSON("/api/v1/me", &me)
	if me.User.Email != "dev@codequest.local" {
		t.Fatalf("me = %+v", me)
	}

	// The default dev user starts with XP and Bytes, granted once.
	on.doRaw(jsonReq(http.MethodPost, "/api/v1/auth/dev", `{}`))
	var stats struct {
		Stats struct {
			TotalXP int `json:"total_xp"`
			Bytes   int `json:"bytes"`
		}
	}
	on.getJSON("/api/v1/me", &stats)
	if stats.Stats.TotalXP != economy.DevStarterXP || stats.Stats.Bytes != economy.DevStarterBytes {
		t.Fatalf("dev stats = %+v, want %d XP and %d Bytes once", stats.Stats, economy.DevStarterXP, economy.DevStarterBytes)
	}

	// A dev sign-in with an email is an ordinary fresh user.
	other := on.doRaw(jsonReq(http.MethodPost, "/api/v1/auth/dev", `{"email":"fresh@example.com"}`))
	on.cookie = sessionFrom(t, other).Value
	on.getJSON("/api/v1/me", &stats)
	if stats.Stats.TotalXP != 0 || stats.Stats.Bytes != 0 {
		t.Fatalf("fresh dev user stats = %+v, want zero", stats.Stats)
	}
}

// State-changing requests need an allowed Origin, even with a valid session.
func TestCSRF(t *testing.T) {
	e := setup(t, true)
	for name, hdr := range map[string]map[string]string{
		"no origin":    {"Origin": ""},
		"other origin": {"Origin": "https://evil.example"},
	} {
		req := jsonReq(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("b"))
		req.AddCookie(&http.Cookie{Name: "cq_session", Value: e.cookie})
		for k, v := range hdr {
			if v == "" {
				req.Header.Del(k)
			} else {
				req.Header.Set(k, v)
			}
		}
		if rec := e.doRaw(req); rec.Code != http.StatusForbidden {
			t.Errorf("%s: %d, want 403", name, rec.Code)
		}
	}
}

// fakeGitHub serves the three GitHub endpoints sign-in uses.
func fakeGitHub(t *testing.T, email string) *auth.GitHub {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("code") != "good-code" || r.Form.Get("client_secret") != "secret" {
			_, _ = w.Write([]byte(`{"error":"bad_verification_code"}`))
			return
		}
		_, _ = w.Write([]byte(`{"access_token":"tok"}`))
	})
	mux.HandleFunc("GET /user", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte(`{"id":4242,"login":"octo","name":"Octo Cat","avatar_url":"https://avatars.example/4242"}`))
	})
	mux.HandleFunc("GET /user/emails", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"email":"other@example.com","primary":false,"verified":true},{"email":"` + email + `","primary":true,"verified":true}]`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return &auth.GitHub{
		ClientID: "id", ClientSecret: "secret", RedirectURL: origin + "/api/v1/auth/github/callback",
		AuthorizeURL: srv.URL + "/login/oauth/authorize", TokenURL: srv.URL + "/login/oauth/access_token", APIURL: srv.URL,
	}
}

func TestGitHubSignIn(t *testing.T) {
	e := setupWith(t, false, opts{github: fakeGitHub(t, "Octo@Example.com")})
	// An existing account with the same verified email gets linked.
	existing, _ := e.signIn("octo@example.com", "")

	rec := e.doRaw(httptest.NewRequest(http.MethodGet, "/api/v1/auth/github/start?next=/courses/foundations", nil))
	if rec.Code != http.StatusFound {
		t.Fatalf("start = %d", rec.Code)
	}
	loc, _ := url.Parse(rec.Header().Get("Location"))
	state := loc.Query().Get("state")
	if state == "" || loc.Query().Get("client_id") != "id" {
		t.Fatalf("authorize URL = %s", loc)
	}
	var stateCookie *http.Cookie
	for _, c := range rec.Result().Cookies() {
		if c.Name == "cq_oauth" {
			stateCookie = c
		}
	}
	if stateCookie == nil || !stateCookie.HttpOnly {
		t.Fatalf("state cookie = %+v", stateCookie)
	}

	callback := func(q string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/github/callback?"+q, nil)
		req.AddCookie(stateCookie)
		return e.doRaw(req)
	}
	// A forged state is refused before talking to GitHub.
	if rec := callback("code=good-code&state=forged"); rec.Header().Get("Location") != "/sign-in?error=github" {
		t.Fatalf("forged state: %d %s", rec.Code, rec.Header().Get("Location"))
	}
	if rec := callback("code=bad-code&state=" + state); rec.Header().Get("Location") != "/sign-in?error=github" {
		t.Fatalf("bad code: %d %s", rec.Code, rec.Header().Get("Location"))
	}
	rec = callback("code=good-code&state=" + state)
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/courses/foundations" {
		t.Fatalf("callback = %d %s", rec.Code, rec.Header().Get("Location"))
	}
	e.cookie = sessionFrom(t, rec).Value
	var me struct {
		User struct{ ID, Email, Name, Avatar string }
	}
	e.getJSON("/api/v1/me", &me)
	if me.User.ID != existing.UUID.String() || me.User.Name != "Octo Cat" || me.User.Avatar == "" {
		t.Fatalf("me = %+v, want the linked existing account", me)
	}
}

func TestOAuthNotConfigured(t *testing.T) {
	e := setup(t, false)
	for _, p := range []string{"github", "google", "nope"} {
		if code, _ := e.do(http.MethodGet, "/api/v1/auth/"+p+"/start", "", nil); code != http.StatusNotFound {
			t.Errorf("%s start = %d", p, code)
		}
	}
	var prov map[string]bool
	e.getJSON("/api/v1/auth/providers", &prov)
	if prov["github"] || prov["google"] || len(prov) != 3 {
		t.Fatalf("providers = %v", prov)
	}
}

// fakeGoogle serves Google's token and userinfo endpoints.
func fakeGoogle(t *testing.T, email string, verified bool) *auth.Google {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("code") != "good-code" || r.Form.Get("client_secret") != "gsecret" || r.Form.Get("grant_type") != "authorization_code" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
			return
		}
		_, _ = w.Write([]byte(`{"access_token":"gtok","token_type":"Bearer"}`))
	})
	mux.HandleFunc("GET /userinfo", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer gtok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		b, _ := json.Marshal(map[string]any{"sub": "1098", "email": email, "email_verified": verified, "name": "Goo Gle", "picture": "https://pics.example/1098"})
		_, _ = w.Write(b)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return &auth.Google{
		ClientID: "gid", ClientSecret: "gsecret", RedirectURL: origin + "/api/v1/auth/google/callback",
		AuthorizeURL: srv.URL + "/auth", TokenURL: srv.URL + "/token", UserinfoURL: srv.URL + "/userinfo",
	}
}

// googleSignIn runs start + callback and returns the callback response.
func googleSignIn(t *testing.T, e *env, code string) *httptest.ResponseRecorder {
	t.Helper()
	rec := e.doRaw(httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/start?next=/courses/git&tz=Asia%2FKolkata", nil))
	loc, _ := url.Parse(rec.Header().Get("Location"))
	if rec.Code != http.StatusFound || loc.Query().Get("client_id") != "gid" || loc.Query().Get("scope") != "openid email profile" {
		t.Fatalf("start = %d %s", rec.Code, loc)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?code="+code+"&state="+loc.Query().Get("state"), nil)
	for _, c := range rec.Result().Cookies() {
		if c.Name == "cq_oauth" {
			if c.Path != "/api/v1/auth/google" || !c.HttpOnly {
				t.Fatalf("state cookie = %+v", c)
			}
			req.AddCookie(c)
		}
	}
	return e.doRaw(req)
}

// Google sign-in links to the account with the same verified email, and
// refuses an unverified one.
func TestGoogleSignIn(t *testing.T) {
	e := setupWith(t, false, opts{google: fakeGoogle(t, "Goo@Example.com", true)})
	existing, _ := e.signIn("goo@example.com", "")
	var prov map[string]bool
	e.getJSON("/api/v1/auth/providers", &prov)
	if !prov["google"] || prov["github"] {
		t.Fatalf("providers = %v", prov)
	}
	if rec := googleSignIn(t, e, "bad-code"); rec.Header().Get("Location") != "/sign-in?error=google" {
		t.Fatalf("bad code: %d %s", rec.Code, rec.Header().Get("Location"))
	}
	rec := googleSignIn(t, e, "good-code")
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/courses/git" {
		t.Fatalf("callback = %d %s", rec.Code, rec.Header().Get("Location"))
	}
	e.cookie = sessionFrom(t, rec).Value
	var me struct {
		User struct{ ID, Name, Avatar, Timezone string }
	}
	e.getJSON("/api/v1/me", &me)
	if me.User.ID != existing.UUID.String() || me.User.Name != "Goo Gle" || me.User.Avatar == "" || me.User.Timezone != "Asia/Kolkata" {
		t.Fatalf("me = %+v, want the linked existing account", me)
	}

	unverified := setupWith(t, false, opts{google: fakeGoogle(t, "new@example.com", false)})
	if rec := googleSignIn(t, unverified, "good-code"); rec.Header().Get("Location") != "/sign-in?error=google" {
		t.Fatalf("unverified email: %d %s", rec.Code, rec.Header().Get("Location"))
	}
}
