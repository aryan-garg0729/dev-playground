package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aryangarg/codequest/api/internal/auth"
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

func tokenFromMail(t *testing.T, body string) string {
	t.Helper()
	_, after, ok := strings.Cut(body, "/sign-in/verify#token=")
	if !ok {
		t.Fatalf("no link in email: %q", body)
	}
	return strings.Fields(after)[0]
}

func TestMagicLinkSignIn(t *testing.T) {
	e := setup(t, false)
	code, body := e.do(http.MethodPost, "/api/v1/auth/magic-link", `{"email":" New@Example.com ","next":"/learn/foundations/fd.py.sum-list"}`, nil)
	if code != http.StatusAccepted {
		t.Fatalf("send = %d: %s", code, body)
	}
	m := e.mail.last()
	if m.to != "new@example.com" || !strings.HasPrefix(m.body, "Sign in to CodeQuest:\n\n"+origin+"/sign-in/verify#token=") {
		t.Fatalf("email = %+v", m)
	}
	token := tokenFromMail(t, m.body)

	rec := e.doRaw(jsonReq(http.MethodPost, "/api/v1/auth/magic-link/verify", `{"token":"`+token+`","timezone":"Asia/Kolkata"}`))
	if rec.Code != http.StatusOK {
		t.Fatalf("verify = %d: %s", rec.Code, rec.Body)
	}
	c := sessionFrom(t, rec)
	if !c.HttpOnly || c.SameSite != http.SameSiteLaxMode || c.Path != "/" {
		t.Fatalf("cookie flags: %+v", c)
	}
	var signed struct {
		User struct{ Email, Timezone string }
		Next string
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &signed)
	if signed.User.Email != "new@example.com" || signed.User.Timezone != "Asia/Kolkata" || signed.Next != "/learn/foundations/fd.py.sum-list" {
		t.Fatalf("verify body = %s", rec.Body)
	}

	// The session works; the database holds only a hash of it.
	e.cookie = c.Value
	var me struct{ User struct{ Email string } }
	e.getJSON("/api/v1/me", &me)
	if me.User.Email != "new@example.com" {
		t.Fatalf("me = %+v", me)
	}
	var stored int
	_ = e.pool.QueryRow(t.Context(), `SELECT count(*) FROM sessions WHERE id = $1`, c.Value).Scan(&stored)
	if stored != 0 {
		t.Fatal("raw session id stored in the database")
	}

	// A link works once.
	code, body = e.do(http.MethodPost, "/api/v1/auth/magic-link/verify", `{"token":"`+token+`"}`, nil)
	if code != http.StatusBadRequest || !strings.Contains(string(body), "invalid_link") {
		t.Fatalf("reuse = %d: %s", code, body)
	}

	// Sign out ends the session server-side.
	if code, _ := e.do(http.MethodPost, "/api/v1/auth/logout", `{}`, nil); code != http.StatusNoContent {
		t.Fatalf("logout = %d", code)
	}
	if code, _ := e.do(http.MethodGet, "/api/v1/me", "", nil); code != http.StatusUnauthorized {
		t.Fatalf("me after logout = %d", code)
	}
}

func TestMagicLinkRules(t *testing.T) {
	e := setup(t, false)
	// Existing and unknown addresses get the same answer.
	e.signIn("known@example.com", "K")
	c1, b1 := e.do(http.MethodPost, "/api/v1/auth/magic-link", `{"email":"known@example.com"}`, nil)
	c2, b2 := e.do(http.MethodPost, "/api/v1/auth/magic-link", `{"email":"unknown@example.com"}`, nil)
	if c1 != c2 || string(b1) != string(b2) {
		t.Fatalf("responses differ: %d %s / %d %s", c1, b1, c2, b2)
	}
	for _, bad := range []string{`{"email":"nope"}`, `{"email":"a@b"}`, `{"email":"Name <a@example.com>"}`, `{"email":"a@example.com","extra":1}`} {
		if code, _ := e.do(http.MethodPost, "/api/v1/auth/magic-link", bad, nil); code != http.StatusBadRequest {
			t.Errorf("%s = %d, want 400", bad, code)
		}
	}
	// Off-site next paths are dropped.
	e.do(http.MethodPost, "/api/v1/auth/magic-link", `{"email":"x@example.com","next":"//evil.example/x"}`, nil)
	code, body := e.do(http.MethodPost, "/api/v1/auth/magic-link/verify", `{"token":"`+tokenFromMail(t, e.mail.last().body)+`"}`, nil)
	if code != http.StatusOK || !strings.Contains(string(body), `"next":"/"`) {
		t.Fatalf("verify with off-site next = %d: %s", code, body)
	}
	// Expired links fail.
	e.do(http.MethodPost, "/api/v1/auth/magic-link", `{"email":"y@example.com"}`, nil)
	tok := tokenFromMail(t, e.mail.last().body)
	if _, err := e.pool.Exec(t.Context(), `UPDATE magic_links SET expires_at = now() - interval '1 second'`); err != nil {
		t.Fatal(err)
	}
	if code, _ := e.do(http.MethodPost, "/api/v1/auth/magic-link/verify", `{"token":"`+tok+`"}`, nil); code != http.StatusBadRequest {
		t.Fatalf("expired link = %d", code)
	}
	// Five emails per address per hour.
	for i := range 6 {
		code, _ := e.do(http.MethodPost, "/api/v1/auth/magic-link", `{"email":"flood@example.com"}`, nil)
		if want := map[bool]int{true: 202, false: 429}[i < 5]; code != want {
			t.Fatalf("email %d = %d, want %d", i+1, code, want)
		}
	}
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
	// An existing magic-link account with the same verified email gets linked.
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

func TestGitHubNotConfigured(t *testing.T) {
	e := setup(t, false)
	if code, _ := e.do(http.MethodGet, "/api/v1/auth/github/start", "", nil); code != http.StatusNotFound {
		t.Fatalf("start = %d", code)
	}
}
