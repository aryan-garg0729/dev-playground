package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aryangarg/codequest/api/internal/auth"
)

const cliLesson = "cl.hello-script"

// cliDo sends a request the way the CLI does: no cookie, no Origin, an
// optional bearer token.
func (e *env) cliDo(method, path, body, token string) (int, []byte) {
	e.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.Bytes()
}

// cliLogin gets a code with the browser session and trades it for a token.
func (e *env) cliLogin() string {
	e.t.Helper()
	code, body := e.do(http.MethodPost, "/api/v1/cli/codes", `{}`, nil)
	if code != http.StatusOK {
		e.t.Fatalf("create code = %d: %s", code, body)
	}
	var c struct{ Code string }
	if err := json.Unmarshal(body, &c); err != nil || len(c.Code) != 11 || c.Code[5] != '-' {
		e.t.Fatalf("code = %q (%v)", c.Code, err)
	}
	code, body = e.cliDo(http.MethodPost, "/api/v1/cli/token", `{"code":"`+strings.ToLower(c.Code)+`\n","label":"laptop"}`, "")
	if code != http.StatusOK {
		e.t.Fatalf("token = %d: %s", code, body)
	}
	var tok struct{ Token, Email string }
	if err := json.Unmarshal(body, &tok); err != nil || tok.Token == "" || tok.Email != "t@example.com" {
		e.t.Fatalf("token response = %s", body)
	}
	return tok.Token
}

func cliSubmitBody(stdout string) string {
	b, _ := json.Marshal(map[string]any{"exercise_id": cliLesson, "version": 1,
		"answer": map[string]any{"steps": []any{map[string]any{"cli": map[string]any{"exit_code": 0, "stdout": stdout}}}}})
	return string(b)
}

func TestCLILoginCodes(t *testing.T) {
	e := setup(t, true)
	// Signed-out browsers and CLI tokens cannot mint codes.
	anon := setup(t, false)
	if code, _ := anon.do(http.MethodPost, "/api/v1/cli/codes", `{}`, nil); code != http.StatusUnauthorized {
		t.Fatalf("anonymous code = %d", code)
	}
	tok := e.cliLogin()
	if code, _ := e.cliDo(http.MethodPost, "/api/v1/cli/codes", `{}`, tok); code != http.StatusUnauthorized {
		t.Fatalf("code via bearer = %d", code)
	}
	// A code works once.
	_, body := e.do(http.MethodPost, "/api/v1/cli/codes", `{}`, nil)
	var c struct{ Code string }
	_ = json.Unmarshal(body, &c)
	if code, _ := e.cliDo(http.MethodPost, "/api/v1/cli/token", `{"code":"`+c.Code+`"}`, ""); code != http.StatusOK {
		t.Fatalf("first use = %d", code)
	}
	code, body := e.cliDo(http.MethodPost, "/api/v1/cli/token", `{"code":"`+c.Code+`"}`, "")
	if code != http.StatusBadRequest || !strings.Contains(string(body), "invalid_code") {
		t.Fatalf("second use = %d %s", code, body)
	}
	// An expired code is refused.
	ctx := context.Background()
	exp, _, err := (&auth.CLI{Pool: e.pool, CodeTTL: time.Millisecond}).CreateLoginCode(ctx, e.user.UUID)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)
	if code, _ := e.cliDo(http.MethodPost, "/api/v1/cli/token", `{"code":"`+exp+`"}`, ""); code != http.StatusBadRequest {
		t.Fatalf("expired code = %d", code)
	}
	// Garbage is refused without a lookup.
	if code, _ := e.cliDo(http.MethodPost, "/api/v1/cli/token", `{"code":"nope"}`, ""); code != http.StatusBadRequest {
		t.Fatalf("garbage code = %d", code)
	}
	// /cli/token still needs a JSON content type.
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cli/token", strings.NewReader("code=x"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("form post = %d", rec.Code)
	}
}

func TestCLIBearerTokens(t *testing.T) {
	e := setup(t, true)
	tok := e.cliLogin()

	if code, body := e.cliDo(http.MethodGet, "/api/v1/me", "", tok); code != http.StatusOK || !strings.Contains(string(body), "t@example.com") {
		t.Fatalf("bearer /me = %d %s", code, body)
	}
	// Bearer + cookie is refused, so the two credentials never mix.
	if code, _ := e.do(http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer " + tok}); code != http.StatusBadRequest {
		t.Fatalf("bearer + cookie = %d", code)
	}
	for name, h := range map[string]string{"unknown": "Bearer nope", "basic": "Basic dTpw"} {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
		req.Header.Set("Authorization", h)
		rec := httptest.NewRecorder()
		e.h.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s token = %d", name, rec.Code)
		}
	}
	// A cookie request without Origin is still refused.
	if code, _ := e.do(http.MethodPost, "/api/v1/exercises/"+cliLesson+"/submit", cliSubmitBody("x"), map[string]string{"Origin": ""}); code != http.StatusForbidden {
		t.Fatalf("cookie without origin = %d", code)
	}
	// Revoke: the token stops working.
	if code, _ := e.cliDo(http.MethodPost, "/api/v1/cli/logout", `{}`, tok); code != http.StatusNoContent {
		t.Fatalf("logout = %d", code)
	}
	if code, _ := e.cliDo(http.MethodGet, "/api/v1/me", "", tok); code != http.StatusUnauthorized {
		t.Fatalf("revoked token = %d", code)
	}
	// An expired token is refused.
	tok2 := e.cliLogin()
	if _, err := e.pool.Exec(context.Background(), `UPDATE cli_tokens SET expires_at = now() - interval '1 second' WHERE id = $1`, auth.Hash(tok2)); err != nil {
		t.Fatal(err)
	}
	if code, _ := e.cliDo(http.MethodGet, "/api/v1/me", "", tok2); code != http.StatusUnauthorized {
		t.Fatalf("expired token = %d", code)
	}
}
