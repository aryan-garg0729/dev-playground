package projectcli

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"testing/fstest"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

func node(t *testing.T, src string) yaml.Node {
	t.Helper()
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(src), &doc); err != nil {
		t.Fatal(err)
	}
	return *doc.Content[0]
}

const validSpec = `
prompt: prompt.md
base_url: http://localhost:8000
solution: solution.json
steps:
  - description: Create a user
    cli_command:
      command: python3 make_user.py
      captures: [{ name: user, regex: "created (\\w+)" }]
      tests:
        - exit_code: 0
        - stdout_contains_all: ["created"]
        - stdout_contains_none: ["Traceback"]
        - stdout_lines_gt: 0
  - description: Fetch the user
    http_request:
      request: { method: POST, url: "${baseURL}/users/${user}", headers: { X-Who: "${user}" }, body_json: { name: "${user}", n: 1 } }
      captures: [{ name: uid, json_path: ".id" }]
      tests:
        - status_code: 200
        - body_contains: "${user}"
        - json_body: { path: ".tags[1]", op: eq, value: b }
        - json_body: { path: ".count", op: gt, value: 2 }
        - json_body: { path: ".tags", op: contains, value: a }
        - json_body: { path: ".id", op: eq, value: "${uid}" }
        - header: { name: content-type, contains: json }
`

const solution = `{"steps":[
  {"cli":{"exit_code":0,"stdout":"created alice\n"}},
  {"http":{"status":200,"headers":{"Content-Type":"application/json"},"body":"{\"id\":\"u1\",\"name\":\"alice\",\"count\":3,\"tags\":[\"a\",\"b\"]}"}}
]}`

func fsWith(files map[string]string) fstest.MapFS {
	m := fstest.MapFS{"prompt.md": {Data: []byte("Make a user.")}, "solution.json": {Data: []byte(solution)}}
	for k, v := range files {
		m[k] = &fstest.MapFile{Data: []byte(v)}
	}
	return m
}

func parse(t *testing.T, src string, dir fstest.MapFS) (exercise.Exercise, error) {
	t.Helper()
	typ, _ := exercise.Get(TypeName)
	s, err := typ.ParseSpec(node(t, src), dir)
	return exercise.Exercise{ID: "t.cli", Type: TypeName, Spec: s}, err
}

func mustParse(t *testing.T) exercise.Exercise {
	t.Helper()
	ex, err := parse(t, validSpec, fsWith(nil))
	if err != nil {
		t.Fatal(err)
	}
	return ex
}

func gradeAnswer(t *testing.T, ex exercise.Exercise, answer string) exercise.Result {
	t.Helper()
	typ, _ := exercise.Get(TypeName)
	r, err := typ.Grade(context.Background(), ex, exercise.Submission{Answer: json.RawMessage(answer)})
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestParseSpecErrors(t *testing.T) {
	tests := []struct{ name, src, want string }{
		{"unknown field", validSpec + "colour: red\n", "colour"},
		{"no steps", "prompt: prompt.md\nsolution: solution.json\nsteps: []\n", "need 1 to 20 steps"},
		{"both kinds", strings.Replace(validSpec, "    http_request:", "    cli_command: { command: ls, tests: [{exit_code: 0}] }\n    http_request:", 1), "exactly one of cli_command and http_request"},
		{"no tests", "prompt: prompt.md\nsolution: solution.json\nsteps:\n  - { description: x, cli_command: { command: ls, tests: [] } }\n", "need at least one test"},
		{"two fields in one test", strings.Replace(validSpec, "- exit_code: 0", "- { exit_code: 0, stdout_lines_gt: 1 }", 1), "exactly one of exit_code"},
		{"undefined var", strings.Replace(validSpec, "/users/${user}", "/users/${nobody}", 1), "${nobody}"},
		{"capture used before it exists", strings.Replace(validSpec, "command: python3 make_user.py", "command: echo ${user}", 1), "${user}"},
		{"bad method", strings.Replace(validSpec, "method: POST", "method: TRACE", 1), "method"},
		{"bad url", strings.Replace(validSpec, `url: "${baseURL}/users/${user}"`, `url: "ftp://x"`, 1), "url"},
		{"bad regex", strings.Replace(validSpec, `regex: "created (\\w+)"`, `regex: "("`, 1), "regex"},
		{"json_path on cli", strings.Replace(validSpec, `regex: "created (\\w+)"`, `json_path: ".x"`, 1), "json_path is for http_request only"},
		{"reserved capture", strings.Replace(validSpec, "name: user,", "name: baseURL,", 1), "not a valid name"},
		{"bad op", strings.Replace(validSpec, "op: gt", "op: ge", 1), "op"},
		{"gt needs number", strings.Replace(validSpec, "op: gt, value: 2", "op: gt, value: x", 1), "needs a number"},
		{"bad path", strings.Replace(validSpec, `path: ".count"`, `path: "count"`, 1), "must start with"},
		{"bad os", validSpec + "allowed_os: [plan9]\n", "plan9"},
		{"missing solution", strings.Replace(validSpec, "solution: solution.json", "", 1), "solution: required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(t, tt.src, fsWith(nil))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want it to mention %q", err, tt.want)
			}
		})
	}
}

func TestParseSpecRecordedAnswers(t *testing.T) {
	if _, err := parse(t, validSpec, fsWith(map[string]string{"solution.json": `{"steps":[]}`})); err == nil || !strings.Contains(err.Error(), "reported 0 steps") {
		t.Fatalf("short solution: err = %v", err)
	}
	if _, err := parse(t, validSpec, fsWith(map[string]string{"solution.json": strings.Replace(solution, `"steps"`, `"stepz":[],"steps"`, 1)})); err == nil {
		t.Fatal("unknown field in solution.json accepted")
	}
	if _, err := parse(t, validSpec, fsWith(map[string]string{"wrong/x.txt": "hi"})); err == nil || !strings.Contains(err.Error(), ".json") {
		t.Fatalf("non-json wrong file: err = %v", err)
	}
	ex, err := parse(t, validSpec, fsWith(map[string]string{"wrong/exit.json": strings.Replace(solution, `"exit_code":0`, `"exit_code":1`, 1)}))
	if err != nil {
		t.Fatal(err)
	}
	typ, _ := exercise.Get(TypeName)
	wrong, _ := typ.(exercise.Testable).WrongAnswers(ex.Spec)
	if len(wrong) != 1 || wrong["wrong/exit.json"] == nil {
		t.Fatalf("wrong answers = %v", wrong)
	}
}

func TestSolutionPasses(t *testing.T) {
	ex := mustParse(t)
	r := gradeAnswer(t, ex, solution)
	if r.Status != exercise.StatusPass || r.Score != 1 {
		t.Fatalf("status %s score %v: %+v", r.Status, r.Score, r.Checks)
	}
	if len(r.Checks) != 11 {
		t.Fatalf("got %d rows, want one per test (11)", len(r.Checks))
	}
	if r.Checks[0].Name != "Step 1: exit code is 0" || !r.Checks[0].Visible {
		t.Fatalf("first row = %+v", r.Checks[0])
	}
	if !strings.Contains(r.Checks[5].Name, `body contains "alice"`) {
		t.Fatalf("captured value not interpolated into the name: %q", r.Checks[5].Name)
	}
	if !strings.Contains(r.Stdout, "$ python3 make_user.py\ncreated alice") || !strings.Contains(r.Stdout, "> POST http://localhost:8000/users/alice") {
		t.Fatalf("transcript = %q", r.Stdout)
	}
}

func TestFailures(t *testing.T) {
	ex := mustParse(t)
	tests := []struct {
		name, from, to, failing string
	}{
		{"exit code", `"exit_code":0`, `"exit_code":2`, "Step 1: exit code is 0"},
		{"missing text", `created alice\n`, `made alice\n`, "Step 1: output contains"},
		{"forbidden text", `created alice\n`, `created alice\nTraceback\n`, "does not contain"},
		{"status", `"status":200`, `"status":404`, "status code is 200"},
		{"json eq", `[\"a\",\"b\"]`, `[\"a\",\"c\"]`, `JSON .tags[1] is "b"`},
		{"json gt", `\"count\":3`, `\"count\":2`, "greater than"},
		{"header", `"Content-Type":"application/json"`, `"Content-Type":"text/html"`, "header content-type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := strings.Replace(solution, tt.from, tt.to, 1)
			if ans == solution {
				t.Fatal("replacement did not apply")
			}
			r := gradeAnswer(t, ex, ans)
			if r.Status != exercise.StatusFail {
				t.Fatalf("status = %s", r.Status)
			}
			for _, c := range r.Checks {
				if strings.Contains(c.Name, tt.failing) {
					if c.Passed {
						t.Fatalf("%q passed", c.Name)
					}
					return
				}
			}
			t.Fatalf("no row named %q", tt.failing)
		})
	}
}

func TestCapturesAreRecomputedOnTheServer(t *testing.T) {
	ex := mustParse(t)
	// The step 1 output captures "bob", but the client claims user=alice and
	// reports a body about alice. The body test must use the server's value.
	ans := strings.Replace(solution, `created alice`, `created bob`, 1)
	ans = strings.Replace(ans, `{"steps"`, `{"variables":{"user":"alice"},"steps"`, 1)
	r := gradeAnswer(t, ex, ans)
	if r.Status != exercise.StatusFail {
		t.Fatalf("forged variables passed: %+v", r.Checks)
	}
	for _, c := range r.Checks {
		if strings.HasPrefix(c.Name, "Step 2: body contains") && (c.Passed || c.Expected != "bob") {
			t.Fatalf("body row = %+v, want it to expect the server-captured bob", c)
		}
	}
}

func TestDependencyAndSkip(t *testing.T) {
	ex := mustParse(t)
	// No capture match in step 1: step 2 needs ${user} and fails entirely.
	r := gradeAnswer(t, ex, strings.Replace(solution, `created alice`, `nothing`, 1))
	for _, c := range r.Checks {
		if strings.HasPrefix(c.Name, "Step 2") && (c.Passed || !strings.Contains(c.Feedback, "${user}")) {
			t.Fatalf("step 2 row = %+v", c)
		}
	}
	// A skipped step fails all of its tests.
	skipped := `{"steps":[{"cli":{"exit_code":0,"stdout":"created alice"}},{"skipped":true}]}`
	r = gradeAnswer(t, ex, skipped)
	if r.Status != exercise.StatusFail || r.Checks[len(r.Checks)-1].Feedback != "This step was not run" {
		t.Fatalf("skipped: %+v", r.Checks)
	}
	// No response at all.
	noResp := `{"steps":[{"cli":{"exit_code":0,"stdout":"created alice"}},{"http":{"status":0,"body":"","error":"connection refused"}}]}`
	r = gradeAnswer(t, ex, noResp)
	if last := r.Checks[len(r.Checks)-1]; last.Passed || !strings.Contains(last.Feedback, "connection refused") {
		t.Fatalf("no response row = %+v", last)
	}
}

func TestInvalidAnswers(t *testing.T) {
	ex := mustParse(t)
	typ, _ := exercise.Get(TypeName)
	big := strings.Repeat("x", MaxOutputBytes+1)
	tests := map[string]string{
		"not json":        `nope`,
		"wrong count":     `{"steps":[{"skipped":true}]}`,
		"two kinds":       `{"steps":[{"skipped":true,"cli":{"exit_code":0,"stdout":""}},{"skipped":true}]}`,
		"none":            `{"steps":[{},{"skipped":true}]}`,
		"wrong kind":      `{"steps":[{"http":{"status":200,"body":""}},{"skipped":true}]}`,
		"stdout too big":  `{"steps":[{"cli":{"exit_code":0,"stdout":"` + big + `"}},{"skipped":true}]}`,
		"bad status":      `{"steps":[{"skipped":true},{"http":{"status":42,"body":""}}]}`,
		"headers too big": `{"steps":[{"skipped":true},{"http":{"status":200,"body":"","headers":{"a":"` + strings.Repeat("y", MaxHeaderBytes) + `"}}}]}`,
	}
	for name, ans := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := typ.Grade(context.Background(), ex, exercise.Submission{Answer: json.RawMessage(ans)})
			if !errors.Is(err, exercise.ErrInvalidAnswer) {
				t.Fatalf("err = %v, want ErrInvalidAnswer", err)
			}
		})
	}
}

func TestBaseURLOverrideAndPublicView(t *testing.T) {
	ex := mustParse(t)
	ans := strings.Replace(solution, `{"steps"`, `{"base_url":"http://127.0.0.1:9000/","steps"`, 1)
	r := gradeAnswer(t, ex, ans)
	if !strings.Contains(r.Stdout, "> POST http://127.0.0.1:9000/users/alice") {
		t.Fatalf("override not used: %q", r.Stdout)
	}
	typ, _ := exercise.Get(TypeName)
	pv, err := typ.PublicView(ex)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(pv), "created alice") || strings.Contains(string(pv), "solution") {
		t.Fatalf("public view leaks the recorded solution: %s", pv)
	}
	// private_spec round-trip
	b, _ := json.Marshal(ex.Spec)
	s, err := typ.LoadSpec(b)
	if err != nil {
		t.Fatal(err)
	}
	ex.Spec = s
	if r := gradeAnswer(t, ex, solution); r.Status != exercise.StatusPass {
		t.Fatalf("after LoadSpec: %+v", r.Checks)
	}
}

func TestPath(t *testing.T) {
	var v any
	_ = json.Unmarshal([]byte(`{"a":{"b":[1,{"c":"x"}]}}`), &v)
	if got, ok := lookup(v, ".a.b[1].c"); !ok || got != "x" {
		t.Fatalf("lookup = %v %v", got, ok)
	}
	if _, ok := lookup(v, ".a.b[5]"); ok {
		t.Fatal("out of range found")
	}
	if got, ok := lookup(v, "."); !ok || got == nil {
		t.Fatal("root not found")
	}
	for _, bad := range []string{"", "a", ".a[", ".a[-1]", ".."} {
		if _, err := parsePath(bad); err == nil {
			t.Fatalf("parsePath(%q) accepted", bad)
		}
	}
}
