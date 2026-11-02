package code_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"testing/fstest"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/exercise"
	_ "github.com/aryangarg/codequest/api/internal/exercise/code"
	"github.com/aryangarg/codequest/api/internal/runner"
)

const specYAML = `
language: python
runner: docker-python
prompt: prompt.md
entry: main.py
starter: starter.py
solution: solution.py
limits: { timeout_ms: 2000, memory_mb: 128 }
examples:
  - { input: "nums = [1, 2, 3]", output: "6", explanation: "1 + 2 + 3 = 6" }
constraints: ["0 ≤ len(nums)"]
hints: ["What should total start at?"]
checks:
  - kind: unit_tests
    file: tests_visible.py
    visible: true
  - kind: unit_tests
    file: tests_hidden.py
    weight: 2
  - kind: forbidden_calls
    names: [sum]
`

const (
	solution = "# SOLUTION_SECRET\ndef sum_list(nums):\n    total = 0\n    for n in nums:\n        total += n\n    return total\n"
	visible  = "from cq import expect\nfrom main import sum_list\n\ndef test_one():\n    \"sum_list([5])\"\n    expect(sum_list([5]), 5)\n\ndef test_three():\n    \"sum_list([1, 2, 3])\"\n    expect(sum_list([1, 2, 3]), 6)\n"
	hidden   = "# HIDDEN_SECRET\nfrom main import sum_list\n\ndef test_negative():\n    \"handles negative numbers\"\n    assert sum_list([-1000, 999]) == -1\n\ndef test_empty():\n    assert sum_list([]) == 0\n"
)

func files() fstest.MapFS {
	return fstest.MapFS{
		"prompt.md":             {Data: []byte("Return the total of every integer in `nums`.")},
		"starter.py":            {Data: []byte("def sum_list(nums):\n    total = 0\n    return total\n")},
		"solution.py":           {Data: []byte(solution)},
		"tests_visible.py":      {Data: []byte(visible)},
		"tests_hidden.py":       {Data: []byte(hidden)},
		"wrong/early.py":        {Data: []byte("# WRONG_SECRET\ndef sum_list(nums):\n    for n in nums:\n        return n\n")},
		"wrong/uses_builtin.py": {Data: []byte("def sum_list(nums):\n    return sum(nums)\n")},
	}
}

func parse(t *testing.T, src string, fsys fstest.MapFS) (exercise.Type, exercise.Spec, error) {
	t.Helper()
	typ, ok := exercise.Get("code")
	if !ok {
		t.Fatal("code type not registered")
	}
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(src), &node); err != nil {
		t.Fatal(err)
	}
	spec, err := typ.ParseSpec(*node.Content[0], fsys)
	return typ, spec, err
}

func mustExercise(t *testing.T) (exercise.Type, exercise.Exercise) {
	t.Helper()
	typ, spec, err := parse(t, specYAML, files())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return typ, exercise.Exercise{ID: "t.sum", Type: "code", Version: 1, Spec: spec}
}

func TestParseSpecRejects(t *testing.T) {
	tests := []struct {
		name, find, repl string
		drop             string
	}{
		{"unknown field", "entry: main.py", "entry: main.py\nbogus: 1", ""},
		{"bad language", "language: python", "language: ruby", ""},
		{"bad runner", "runner: docker-python", "runner: docker-go", ""},
		{"entry path", "entry: main.py", "entry: ../main.py", ""},
		{"timeout too big", "timeout_ms: 2000", "timeout_ms: 60000", ""},
		{"memory too big", "memory_mb: 128", "memory_mb: 4096", ""},
		{"bad check config", "names: [sum]", "names: []", ""},
		{"unknown check", "kind: forbidden_calls", "kind: nope", ""},
		{"mcq check", "kind: forbidden_calls\n    names: [sum]", "kind: exact_choice", ""},
		{"missing solution", "", "", "solution.py"},
		{"missing hidden tests", "", "", "tests_hidden.py"},
		{"missing prompt", "", "", "prompt.md"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := files()
			delete(fsys, tt.drop)
			_, _, err := parse(t, strings.Replace(specYAML, tt.find, tt.repl, 1), fsys)
			if err == nil {
				t.Fatal("want error")
			}
		})
	}
}

func TestParseSpecNeedsChecks(t *testing.T) {
	head, _, _ := strings.Cut(specYAML, "checks:")
	noChecks := head + "checks: []\n"
	if _, _, err := parse(t, noChecks, files()); err == nil || !strings.Contains(err.Error(), "at least one") {
		t.Fatalf("err = %v, want 'need at least one'", err)
	}
}

func TestPublicViewHasNoServerOnlyContent(t *testing.T) {
	typ, ex := mustExercise(t)
	pv, err := typ.PublicView(ex)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"SOLUTION_SECRET", "HIDDEN_SECRET", "WRONG_SECRET", "tests_hidden"} {
		if strings.Contains(string(pv), secret) {
			t.Errorf("public view contains %q: %s", secret, pv)
		}
	}
	var v struct {
		Starter      string `json:"starter"`
		VisibleTests []struct {
			File string `json:"file"`
		} `json:"visible_tests"`
	}
	if err := json.Unmarshal(pv, &v); err != nil {
		t.Fatal(err)
	}
	if v.Starter == "" || len(v.VisibleTests) != 1 || v.VisibleTests[0].File != "tests_visible.py" {
		t.Errorf("public view = %s", pv)
	}
}

func TestPrivateSpecRoundTripsWithoutSolution(t *testing.T) {
	typ, ex := mustExercise(t)
	priv, err := json.Marshal(ex.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(priv), "SOLUTION_SECRET") || strings.Contains(string(priv), "WRONG_SECRET") {
		t.Error("solution and wrong/ must not be stored in private_spec")
	}
	if _, err := typ.LoadSpec(priv); err != nil {
		t.Fatal(err)
	}
}

// fakeRunner returns a canned harness report and records the job.
type fakeRunner struct {
	out runner.Output
	err error
	job runner.Job
}

func (f *fakeRunner) Run(_ context.Context, job runner.Job) (runner.Output, error) {
	f.job = job
	return f.out, f.err
}

func TestGradeInvalidAnswers(t *testing.T) {
	typ, ex := mustExercise(t)
	ctx := runner.NewContext(context.Background(), &fakeRunner{})
	for _, ans := range []string{
		`{}`,
		`{"files":{}}`,
		`{"files":{"other.py":"x"}}`,
		`{"files":{"main.py":"x","tests_hidden.py":"y"}}`,
		`{"files":{"main.py":"x"},"passed":true}`,
		`{"files":{"main.py":"` + strings.Repeat("x", 60<<10) + `"}}`,
		`{"files":{"main.py":"a\u0000b"}}`,
	} {
		_, err := typ.Grade(ctx, ex, exercise.Submission{Answer: json.RawMessage(ans)})
		if !errors.Is(err, exercise.ErrInvalidAnswer) {
			t.Errorf("%.60s: err = %v, want ErrInvalidAnswer", ans, err)
		}
	}
}

func TestGradeBuildsJobAndMapsReport(t *testing.T) {
	typ, ex := mustExercise(t)
	report := `{"version":1,"syntax_error":null,"forbidden":[],"program":null,"timed_out":false,"runtime_ms":40,
	 "tests":{"stdout":"hi\n","stderr":"","exit_code":0,"timed_out":false,"runtime_ms":30,"completed":true,"files":[
	  {"file":"tests_visible.py","cases":[{"name":"test_one","label":"sum_list([5])","ran":true,"passed":true},
	                                      {"name":"test_three","label":"sum_list([1, 2, 3])","ran":true,"passed":false,"got":"1","expected":"6"}]},
	  {"file":"tests_hidden.py","cases":[{"name":"test_negative","label":"handles negative numbers","ran":true,"passed":true},
	                                     {"name":"test_empty","label":"","ran":true,"passed":true}]}]}}`
	f := &fakeRunner{out: runner.Output{Stdout: strings.ReplaceAll(report, "\n", "") + "\n"}}
	res, err := typ.Grade(runner.NewContext(context.Background(), f), ex,
		exercise.Submission{Answer: json.RawMessage(`{"files":{"main.py":"def sum_list(n): pass"}}`)})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != exercise.StatusFail || len(res.Checks) != 5 || res.Stdout != "hi\n" || res.RuntimeMs != 30 {
		t.Fatalf("result = %+v", res)
	}
	// visible 1 + 1 (weight 1), hidden 2 + 2 (weight 2), forbidden 1: 6 of 7 passing.
	if res.Score < 0.857 || res.Score > 0.858 {
		t.Errorf("score = %v, want 6/7", res.Score)
	}
	var job struct {
		Files     map[string]string `json:"files"`
		Tests     []struct{ File string }
		Forbidden []string `json:"forbidden"`
		RunEntry  bool     `json:"run_entry"`
	}
	if err := json.Unmarshal(f.job.Stdin, &job); err != nil {
		t.Fatal(err)
	}
	if f.job.Image != "codequest/python:3.12" || f.job.Limits.TimeoutMs != 2000 || len(job.Tests) != 2 ||
		job.Forbidden[0] != "sum" || job.RunEntry || len(job.Files) != 1 {
		t.Errorf("job = %+v / %+v", f.job, job)
	}
}

func TestGradeRunnerOutcomes(t *testing.T) {
	typ, ex := mustExercise(t)
	sub := exercise.Submission{Answer: json.RawMessage(`{"files":{"main.py":"x = 1"}}`)}
	tests := []struct {
		name   string
		f      *fakeRunner
		status string
		err    bool
	}{
		{"hard timeout", &fakeRunner{out: runner.Output{TimedOut: true}}, exercise.StatusTimeout, false},
		{"oom kill", &fakeRunner{out: runner.Output{ExitCode: 137}}, exercise.StatusError, false},
		{"garbage output", &fakeRunner{out: runner.Output{Stdout: "{\"version\":1}\n{\"version\":1}"}}, "", true},
		{"runner down", &fakeRunner{err: runner.ErrUnavailable}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := typ.Grade(runner.NewContext(context.Background(), tt.f), ex, sub)
			if (err != nil) != tt.err || (!tt.err && res.Status != tt.status) {
				t.Fatalf("res=%+v err=%v", res, err)
			}
		})
	}
	if _, err := typ.Grade(context.Background(), ex, sub); err == nil {
		t.Error("grading without a runner in the context must fail")
	}
}
