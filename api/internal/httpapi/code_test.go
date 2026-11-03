package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aryangarg/codequest/api/internal/runner"
)

const codeID = "demo.basics.sum"

func codeBody(src string) string {
	b, _ := json.Marshal(map[string]any{
		"exercise_id": codeID, "version": 1,
		"answer": map[string]any{"files": map[string]string{"main.py": src}},
	})
	return string(b)
}

// report is a harness report in which every test of the fixture passes or fails.
func report(pass bool) string {
	c := func(name string) map[string]any {
		return map[string]any{"name": name, "label": "", "ran": true, "passed": pass, "message": "", "error": "", "got": "", "expected": "", "line": 0}
	}
	r := map[string]any{
		"version": 1, "syntax_error": nil, "forbidden": []any{}, "program": nil, "timed_out": false, "runtime_ms": 5,
		"tests": map[string]any{
			"stdout": "", "stderr": "", "exit_code": 0, "timed_out": false, "runtime_ms": 5, "completed": true,
			"files": []any{
				map[string]any{"file": "tests_visible.py", "cases": []any{c("test_one"), c("test_three")}},
				map[string]any{"file": "tests_hidden.py", "cases": []any{c("test_negative"), c("test_empty")}},
			},
		},
	}
	b, _ := json.Marshal(r)
	return string(b) + "\n"
}

// gatedRunner blocks every run until release is closed, then returns out.
type gatedRunner struct {
	release chan struct{}
	out     runner.Output
	err     error
	mu      sync.Mutex
	jobs    int
}

func (g *gatedRunner) Run(ctx context.Context, _ runner.Job) (runner.Output, error) {
	g.mu.Lock()
	g.jobs++
	g.mu.Unlock()
	if g.release != nil {
		select {
		case <-g.release:
		case <-ctx.Done():
			return runner.Output{}, ctx.Err()
		}
	}
	return g.out, g.err
}

type submitResp struct {
	SubmissionID string `json:"submission_id"`
	Status       string `json:"status"`
	Result       *struct {
		Status string `json:"status"`
	} `json:"result"`
}

func decode(t *testing.T, body []byte, v any) {
	t.Helper()
	if err := json.Unmarshal(body, v); err != nil {
		t.Fatalf("decode %s: %v", body, err)
	}
}

func TestCodeSubmitGradesWithRunner(t *testing.T) {
	r := &gatedRunner{out: runner.Output{Stdout: report(true)}}
	e := setupWith(t, true, opts{runner: r})
	code, body := e.do(http.MethodPost, "/api/v1/exercises/"+codeID+"/submit", codeBody("def sum_list(n): ..."), nil)
	var res submitResp
	decode(t, body, &res)
	if code != http.StatusOK || res.Result == nil || res.Result.Status != "pass" {
		t.Fatalf("submit = %d %s", code, body)
	}
}

func TestCodeSubmitInvalidAnswerNotStored(t *testing.T) {
	r := &gatedRunner{out: runner.Output{Stdout: report(true)}}
	e := setupWith(t, true, opts{runner: r})
	body := strings.Replace(codeBody("x"), `"main.py"`, `"tests_hidden.py"`, 1)
	code, resp := e.do(http.MethodPost, "/api/v1/exercises/"+codeID+"/submit", body, nil)
	if code != 400 || !strings.Contains(string(resp), "invalid_answer") {
		t.Fatalf("submit = %d %s", code, resp)
	}
	if r.jobs != 0 {
		t.Fatal("an invalid answer must not reach the runner")
	}
	assertSubmissions(t, e, 0)
}

func TestGraderUnavailableIs503AndNotStored(t *testing.T) {
	e := setupWith(t, true, opts{runner: &gatedRunner{err: runner.ErrUnavailable}})
	code, body := e.do(http.MethodPost, "/api/v1/exercises/"+codeID+"/submit", codeBody("x = 1"), nil)
	if code != http.StatusServiceUnavailable || !strings.Contains(string(body), "grader_unavailable") {
		t.Fatalf("submit = %d %s", code, body)
	}
	assertSubmissions(t, e, 0)
}

func TestSlowGradingReturns202ThenPolls(t *testing.T) {
	r := &gatedRunner{release: make(chan struct{}), out: runner.Output{Stdout: report(false)}}
	e := setupWith(t, true, opts{runner: r, wait: 50 * time.Millisecond})

	code, body := e.do(http.MethodPost, "/api/v1/exercises/"+codeID+"/submit", codeBody("x = 1"), nil)
	var res submitResp
	decode(t, body, &res)
	if code != http.StatusAccepted || res.Status != "pending" || res.SubmissionID == "" {
		t.Fatalf("submit = %d %s, want 202 pending", code, body)
	}
	poll := "/api/v1/submissions/" + res.SubmissionID
	var sub submitResp
	decode(t, e.getJSON(poll, nil), &sub)
	if sub.Status != "pending" || sub.Result != nil {
		t.Fatalf("before release: %+v", sub)
	}

	close(r.release)
	deadline := time.Now().Add(5 * time.Second)
	for sub.Status == "pending" && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
		decode(t, e.getJSON(poll, nil), &sub)
	}
	if sub.Status != "fail" || sub.Result == nil || sub.Result.Status != "fail" {
		t.Fatalf("after release: %+v", sub)
	}
	var attempts int
	if err := e.pool.QueryRow(context.Background(),
		`SELECT attempts FROM progress WHERE exercise_id = $1`, codeID).Scan(&attempts); err != nil || attempts != 1 {
		t.Fatalf("attempts = %d (%v), want 1", attempts, err)
	}
}

func TestSubmissionsArePrivate(t *testing.T) {
	r := &gatedRunner{out: runner.Output{Stdout: report(true)}}
	e := setupWith(t, true, opts{runner: r})
	ctx := context.Background()
	var other, id string
	if err := e.pool.QueryRow(ctx, `INSERT INTO users (email) VALUES ('other@example.com') RETURNING id`).Scan(&other); err != nil {
		t.Fatal(err)
	}
	if err := e.pool.QueryRow(ctx, `INSERT INTO submissions (user_id, exercise_id, exercise_version, answer, status)
		VALUES ($1, $2, 1, '{}', 'pass') RETURNING id`, other, codeID).Scan(&id); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"/api/v1/submissions/" + id, "/api/v1/submissions/not-a-uuid"} {
		if code, body := e.do(http.MethodGet, p, "", nil); code != 404 {
			t.Errorf("GET %s = %d %s, want 404", p, code, body)
		}
	}
}

func TestSubmitRateLimited(t *testing.T) {
	e := setupWith(t, true, opts{submitPerUser: 2})
	path := "/api/v1/exercises/" + exID + "/submit"
	for i := range 2 {
		if code, body := e.do(http.MethodPost, path, submitBody("a"), nil); code != 200 {
			t.Fatalf("submit %d = %d %s", i+1, code, body)
		}
	}
	code, body := e.do(http.MethodPost, path, submitBody("b"), nil)
	if code != http.StatusTooManyRequests || !strings.Contains(string(body), "rate_limited") {
		t.Fatalf("3rd submit = %d %s, want 429", code, body)
	}
	assertSubmissions(t, e, 2)
}

// TestCodeSubmitEndToEnd grades the fixture through the API with the real
// Docker runner: the reference solution passes, the known-bad one fails.
func TestCodeSubmitEndToEnd(t *testing.T) {
	if err := exec.Command("docker", "image", "inspect", runner.Images["docker-python"]).Run(); err != nil {
		if os.Getenv("CQ_REQUIRE_DOCKER") == "1" {
			t.Fatalf("runner image missing: %v", err)
		}
		t.Skip("docker runner image not available")
	}
	e := setupWith(t, true, opts{runner: runner.Docker{}})
	for _, tt := range []struct{ file, want string }{
		{"solution.py", "pass"},
		{"wrong/early_return.py", "fail"},
		{"wrong/uses_sum.py", "fail"},
	} {
		src, err := os.ReadFile("../content/testdata/valid/demo/01-basics/02-sum/" + tt.file)
		if err != nil {
			t.Fatal(err)
		}
		code, body := e.do(http.MethodPost, "/api/v1/exercises/"+codeID+"/submit", codeBody(string(src)), nil)
		var res submitResp
		decode(t, body, &res)
		if code != 200 || res.Result == nil || res.Result.Status != tt.want {
			t.Errorf("%s: %d %s, want %s", tt.file, code, body, tt.want)
		}
	}
}

func assertSubmissions(t *testing.T, e *env, want int) {
	t.Helper()
	var n int
	if err := e.pool.QueryRow(context.Background(), `SELECT count(*) FROM submissions`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != want {
		t.Fatalf("%d submissions stored, want %d", n, want)
	}
}

func TestServerRun(t *testing.T) {
	r := &gatedRunner{out: runner.Output{Stdout: `{"version":1,"timed_out":false,"runtime_ms":3,"error":"","sql":{"results":[{"learner":{"columns":["name"],"rows":[["a"],["b"]],"truncated":false},"solution":null}]}}` + "\n"}}
	e := setupWith(t, true, opts{runner: r})
	body := `{"exercise_id":"demo.basics.query","version":1,"answer":{"query":"SELECT name FROM t"}}`
	code, resp := e.do(http.MethodPost, "/api/v1/exercises/demo.basics.query/run", body, nil)
	if code != 200 || !strings.Contains(string(resp), `"table":{"columns":["name"]`) || !strings.Contains(string(resp), `"passed":true`) {
		t.Fatalf("run = %d %s", code, resp)
	}
	if strings.Contains(string(resp), sentinel) {
		t.Fatal("run output leaks server-only content")
	}
	// The browser-run fixture refuses server Run.
	code, resp = e.do(http.MethodPost, "/api/v1/exercises/"+codeID+"/run", codeBody("x = 1"), nil)
	if code != 400 || !strings.Contains(string(resp), "run_unavailable") {
		t.Fatalf("browser exercise run = %d %s", code, resp)
	}
	assertSubmissions(t, e, 0)
}
