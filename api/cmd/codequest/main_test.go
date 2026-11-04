package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aryangarg/codequest/api/internal/exercise/projectcli"
)

func intp(n int) *int       { return &n }
func strp(s string) *string { return &s }

func TestRunCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses sh")
	}
	sr := newStepRunner()
	ctx := context.Background()
	r := sr.runCommand(ctx, `echo out; echo err >&2; exit 3`)
	if r.ExitCode != 3 || !strings.Contains(r.Stdout, "out") || !strings.Contains(r.Stdout, "err") {
		t.Fatalf("result = %+v", r)
	}
	if r := sr.runCommand(ctx, `definitely-not-a-command-xyz`); r.ExitCode != 127 {
		t.Fatalf("missing command exit = %d", r.ExitCode)
	}
	big := sr.runCommand(ctx, `head -c 100000 /dev/zero | tr '\0' x`)
	if len(big.Stdout) > projectcli.MaxOutputBytes || !strings.Contains(big.Stdout, "[output truncated]") {
		t.Fatalf("big output: %d bytes", len(big.Stdout))
	}
	bad := sr.runCommand(ctx, `printf '\377\376'`)
	if !strings.Contains(bad.Stdout, "�") {
		t.Fatalf("invalid UTF-8 not replaced: %q", bad.Stdout)
	}
	sr.cmdTimeout = 200 * time.Millisecond
	start := time.Now()
	r = sr.runCommand(ctx, `sleep 10`)
	if r.ExitCode != -1 || !strings.Contains(r.Stdout, "killed") || time.Since(start) > 5*time.Second {
		t.Fatalf("timeout result = %+v after %v", r, time.Since(start))
	}
}

func TestRunStepsCapturesAndHTTP(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses sh")
	}
	var gotAuth, gotWho, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a := r.Header.Get("Authorization"); a != "" {
			gotAuth = a
		}
		if strings.HasPrefix(r.URL.Path, "/users/") && r.Method == http.MethodPost {
			gotWho = r.Header.Get("X-Who")
			b, _ := io.ReadAll(r.Body)
			gotBody = string(b)
		}
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/users/alice", http.StatusFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"u1","path":"` + r.URL.Path + `"}`))
	}))
	defer srv.Close()
	s := projectcli.Spec{
		BaseURL: "http://localhost:1",
		Steps: []projectcli.Step{
			{Description: "make", CLICommand: &projectcli.CLICommand{Command: "echo created alice",
				Captures: []projectcli.Capture{{Name: "user", Regex: `created (\w+)`}}, Tests: []projectcli.CLITest{{ExitCode: intp(0)}}}},
			{Description: "fetch", HTTPRequest: &projectcli.HTTPRequest{
				Request:  projectcli.Request{Method: "POST", URL: "${baseURL}/users/${user}", Headers: map[string]string{"X-Who": "${user}"}, BodyJSON: map[string]any{"name": "${user}"}},
				Captures: []projectcli.Capture{{Name: "uid", JSONPath: ".id"}}, Tests: []projectcli.HTTPTest{{StatusCode: intp(200)}}}},
			{Description: "use", CLICommand: &projectcli.CLICommand{Command: "echo ${uid}", Tests: []projectcli.CLITest{{StdoutContainsAll: []string{"u1"}}}}},
			{Description: "redirect", HTTPRequest: &projectcli.HTTPRequest{
				Request: projectcli.Request{Method: "GET", URL: "${baseURL}/redirect"}, Tests: []projectcli.HTTPTest{{StatusCode: intp(302)}}}},
			{Description: "needs nothing", CLICommand: &projectcli.CLICommand{Command: "echo ${nope}", Tests: []projectcli.CLITest{{ExitCode: intp(0)}}}},
		},
	}
	a := newStepRunner().run(context.Background(), s, srv.URL+"/")
	if a.BaseURL != srv.URL {
		t.Fatalf("base url = %q", a.BaseURL)
	}
	if h := a.Steps[1].HTTP; h == nil || h.Status != 200 || !strings.Contains(h.Body, `"path":"/users/alice"`) || gotWho != "alice" || gotBody != `{"name":"alice"}` {
		t.Fatalf("http step = %+v who=%q body=%q", h, gotWho, gotBody)
	}
	if !strings.Contains(a.Steps[2].CLI.Stdout, "u1") {
		t.Fatalf("json capture not used: %+v", a.Steps[2].CLI)
	}
	if a.Steps[3].HTTP.Status != http.StatusFound {
		t.Fatalf("redirect followed without follow_redirects: %d", a.Steps[3].HTTP.Status)
	}
	if !a.Steps[4].Skipped {
		t.Fatal("a step with an undefined variable ran")
	}
	if gotAuth != "" {
		t.Fatalf("lesson request carried Authorization %q", gotAuth)
	}
	// Connection refused is reported as an error, not a status.
	s2 := projectcli.Spec{BaseURL: "http://127.0.0.1:1", Steps: s.Steps[3:4]}
	if h := newStepRunner().run(context.Background(), s2, "").Steps[0].HTTP; h.Status != 0 || h.Error == "" {
		t.Fatalf("refused = %+v", h)
	}
}

func TestConfigPermissions(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cq")
	t.Setenv("CODEQUEST_CONFIG_DIR", dir)
	c, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	c.Token = "secret"
	if err := c.save(); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(c.path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && fi.Mode().Perm() != 0o600 {
		t.Fatalf("config mode = %v, want 0600", fi.Mode().Perm())
	}
	di, _ := os.Stat(dir)
	if runtime.GOOS != "windows" && di.Mode().Perm() != 0o700 {
		t.Fatalf("dir mode = %v, want 0700", di.Mode().Perm())
	}
	c2, err := loadConfig()
	if err != nil || c2.Token != "secret" {
		t.Fatalf("reload = %+v, %v", c2, err)
	}
}

func TestCheckServiceURL(t *testing.T) {
	for _, ok := range []string{"https://codequest.dev", "http://localhost:8080", "http://127.0.0.1:1"} {
		if err := checkServiceURL(ok); err != nil {
			t.Errorf("%s: %v", ok, err)
		}
	}
	for _, bad := range []string{"http://codequest.dev", "ftp://x", "nope"} {
		if err := checkServiceURL(bad); err == nil {
			t.Errorf("%s accepted", bad)
		}
	}
}

func TestParseArgsInterspersed(t *testing.T) {
	fs := flag.NewFlagSet("t", flag.ContinueOnError)
	yes := fs.Bool("y", false, "")
	pos, err := parseArgs(fs, []string{"cl.x", "-y"})
	if err != nil || len(pos) != 1 || pos[0] != "cl.x" || !*yes {
		t.Fatalf("pos = %v, y = %v, %v", pos, *yes, err)
	}
}

// fakeAPI is a minimal CodeQuest API for the CLI.
type fakeAPI struct {
	mu        sync.Mutex
	spec      projectcli.Spec
	submitted map[string]any
	auths     []string
	revoked   bool
	srv       *httptest.Server
}

func newFakeAPI(t *testing.T, spec projectcli.Spec) *fakeAPI {
	f := &fakeAPI{spec: spec}
	mux := http.NewServeMux()
	write := func(w http.ResponseWriter, code int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(v)
	}
	authed := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f.mu.Lock()
			f.auths = append(f.auths, r.Header.Get("Authorization"))
			f.mu.Unlock()
			if r.Header.Get("Authorization") != "Bearer tok-1" {
				write(w, 401, map[string]any{"error": map[string]string{"code": "unauthorized", "message": "no"}})
				return
			}
			h(w, r)
		}
	}
	mux.HandleFunc("POST /api/v1/cli/token", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Code string }
		_ = json.NewDecoder(r.Body).Decode(&req)
		if strings.ToUpper(strings.TrimSpace(req.Code)) != "ABCDE-FGHJK" { // the real API normalizes too
			write(w, 400, map[string]any{"error": map[string]string{"code": "invalid_code", "message": "bad"}})
			return
		}
		write(w, 200, map[string]string{"token": "tok-1", "name": "Ada", "email": "ada@example.com"})
	})
	mux.HandleFunc("POST /api/v1/cli/logout", authed(func(w http.ResponseWriter, _ *http.Request) {
		f.revoked = true
		w.WriteHeader(204)
	}))
	mux.HandleFunc("GET /api/v1/cli/next", authed(func(w http.ResponseWriter, _ *http.Request) {
		write(w, 200, map[string]any{"exercise_id": "cl.t", "title": "T", "course_title": "C", "chapter_no": 1, "lesson_no": 1})
	}))
	mux.HandleFunc("GET /api/v1/exercises/cl.t", authed(func(w http.ResponseWriter, _ *http.Request) {
		write(w, 200, map[string]any{"id": "cl.t", "type": "project_cli", "version": 4, "title": "T", "public_view": f.spec})
	}))
	mux.HandleFunc("POST /api/v1/exercises/cl.t/submit", authed(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.mu.Lock()
		f.submitted = body
		f.mu.Unlock()
		// Pending first, then the poll returns the result.
		write(w, 202, map[string]string{"submission_id": "s1", "status": "pending"})
	}))
	mux.HandleFunc("GET /api/v1/submissions/s1", authed(func(w http.ResponseWriter, _ *http.Request) {
		write(w, 200, map[string]any{"submission_id": "s1", "status": "pass", "xp_awarded": 20,
			"result":  map[string]any{"status": "pass", "score": 1, "checks": []map[string]any{{"name": "Step 1: exit code is 0", "passed": true}}},
			"outcome": map[string]any{"xp_awarded": 20, "total_xp": 120, "streak_days": 2, "streak_extended": true}})
	}))
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func newTestApp(t *testing.T, apiURL, input string) (*app, *bytes.Buffer) {
	t.Helper()
	t.Setenv("CODEQUEST_CONFIG_DIR", t.TempDir())
	c, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	c.APIURL, c.AppURL = apiURL, "http://localhost:5173"
	var out bytes.Buffer
	return &app{cfg: c, out: &out, err: &out, in: bufio.NewReader(strings.NewReader(input))}, &out
}

func TestLoginSubmitLogout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses sh")
	}
	// A learner's local server: it must never see the CodeQuest token.
	var lessonAuth []string
	var mu sync.Mutex
	local := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		lessonAuth = append(lessonAuth, r.Header.Get("Authorization")+r.Header.Get("Cookie"))
		mu.Unlock()
		_, _ = w.Write([]byte("ok"))
	}))
	defer local.Close()
	f := newFakeAPI(t, projectcli.Spec{Prompt: "p", BaseURL: local.URL, Steps: []projectcli.Step{
		{Description: "say hi", CLICommand: &projectcli.CLICommand{Command: "echo hi", Tests: []projectcli.CLITest{{ExitCode: intp(0)}}}},
		{Description: "call my server", HTTPRequest: &projectcli.HTTPRequest{Request: projectcli.Request{Method: "GET", URL: "${baseURL}/"},
			Tests: []projectcli.HTTPTest{{BodyContains: strp("ok")}}}},
	}})
	ctx := context.Background()

	a, out := newTestApp(t, f.srv.URL, "wrong\n")
	if code := a.run(ctx, []string{"login"}); code != 1 || !strings.Contains(out.String(), "unknown, used or expired") {
		t.Fatalf("bad code: %d %s", code, out)
	}
	a, out = newTestApp(t, f.srv.URL, "abcde-fghjk\n")
	var opened string
	a.open = func(u string) error { opened = u; return nil }
	if code := a.run(ctx, []string{"login"}); code != 0 || a.cfg.Token != "tok-1" || opened != "http://localhost:5173/cli/login" {
		t.Fatalf("login: %d token=%q opened=%q %s", code, a.cfg.Token, opened, out)
	}

	// Declining the confirmation runs nothing.
	a.in = bufio.NewReader(strings.NewReader("n\n"))
	if code := a.run(ctx, []string{"submit"}); code != 1 || f.submitted != nil {
		t.Fatalf("declined submit ran: %d", code)
	}
	// run prints output and no verdict.
	out.Reset()
	if code := a.run(ctx, []string{"run", "cl.t", "-y", "--no-color"}); code != 0 || !strings.Contains(out.String(), "│ hi") || strings.Contains(out.String(), "PASS") || f.submitted != nil {
		t.Fatalf("run: %d %s", code, out)
	}
	// submit sends the raw output with the lesson version and prints the verdict.
	out.Reset()
	a.in = bufio.NewReader(strings.NewReader("\n"))
	if code := a.run(ctx, []string{"submit"}); code != 0 {
		t.Fatalf("submit: %d %s", code, out)
	}
	s := out.String()
	if !strings.Contains(s, "✓ Step 1: exit code is 0") || !strings.Contains(s, "PASS") || !strings.Contains(s, "+20 XP") || !strings.Contains(s, "Streak: 2 days") {
		t.Fatalf("submit output: %s", s)
	}
	ans, _ := json.Marshal(f.submitted)
	if !strings.Contains(string(ans), `"version":4`) || !strings.Contains(string(ans), `"stdout":"hi\n"`) {
		t.Fatalf("submitted %s", ans)
	}

	mu.Lock()
	if len(lessonAuth) != 2 || lessonAuth[0] != "" || lessonAuth[1] != "" {
		t.Fatalf("lesson server saw credentials: %q", lessonAuth)
	}
	mu.Unlock()

	// logout revokes and forgets the token.
	if code := a.run(ctx, []string{"logout"}); code != 0 || !f.revoked || a.cfg.Token != "" {
		t.Fatalf("logout: %d revoked=%v", code, f.revoked)
	}
	c, _ := loadConfig()
	if c.Token != "" {
		t.Fatal("token still on disk after logout")
	}
}

// Criterion: logout deletes the token locally even when the API is unreachable.
func TestLogoutOffline(t *testing.T) {
	a, out := newTestApp(t, "http://127.0.0.1:1", "")
	a.cfg.Token = "tok-1"
	if err := a.cfg.save(); err != nil {
		t.Fatal(err)
	}
	if code := a.run(context.Background(), []string{"logout"}); code != 0 || !strings.Contains(out.String(), "warning") {
		t.Fatalf("logout offline: %d %s", code, out)
	}
	c, _ := loadConfig()
	if c.Token != "" {
		t.Fatal("token kept")
	}
}
