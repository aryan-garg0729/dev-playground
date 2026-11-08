package code_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/runner"
	"github.com/aryangarg/codequest/api/internal/runner/wasm"
)

const goSpecYAML = `
language: go
runner: wasm
prompt: prompt.md
entry: main.go
starter: starter.go
solution: solution.go
limits: { timeout_ms: 2000, memory_mb: 128 }
checks:
  - kind: io_cases
    name: Sample
    cases:
      - { stdin: "3\n", stdout: "1\n2\nFizz\n" }
      - { stdin: "5\n", stdout: "1\n2\nFizz\n4\nBuzz\n", name: "Up to Buzz" }
      - { stdin: "15\n", stdout: "1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n", visible: false }
  - kind: exit_code
    expected: 0
`

const goSolution = `package main

// SOLUTION_SECRET
import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}
}
`

func goFiles() fstest.MapFS {
	return fstest.MapFS{
		"prompt.md":            {Data: []byte("Print FizzBuzz up to n.")},
		"starter.go":           {Data: []byte("package main\n\nfunc main() {\n}\n")},
		"solution.go":          {Data: []byte(goSolution)},
		"wrong/no_fizzbuzz.go": {Data: []byte(strings.Replace(goSolution, `case i%15 == 0:`+"\n\t\t\tfmt.Println(\"FizzBuzz\")\n", "", 1))},
	}
}

func goExercise(t *testing.T) (exercise.Type, exercise.Exercise) {
	t.Helper()
	typ, spec, err := parse(t, goSpecYAML, goFiles())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return typ, exercise.Exercise{ID: "t.fizz", Type: "code", Version: 1, Spec: spec}
}

func TestWasmParseSpecRejects(t *testing.T) {
	tests := []struct{ name, find, repl string }{
		{"python runner for go", "runner: wasm", "runner: docker-python"},
		{"wasm runner for python", "language: go", "language: python"},
		{"browser run", "runner: wasm", "runner: wasm\nrun: browser"},
		{"python entry", "entry: main.go", "entry: main.py"},
		{"unit tests", "  - kind: exit_code\n    expected: 0", "  - kind: unit_tests\n    file: t.py"},
		{"forbidden calls", "  - kind: exit_code\n    expected: 0", "  - kind: forbidden_calls\n    names: [sum]"},
		{"two io_cases", "  - kind: exit_code\n    expected: 0", "  - kind: io_cases\n    cases: [{ stdin: \"\", stdout: \"x\" }]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := goFiles()
			fsys["t.py"] = &fstest.MapFile{Data: []byte("def test_x():\n    pass\n")}
			if _, _, err := parse(t, strings.Replace(goSpecYAML, tt.find, tt.repl, 1), fsys); err == nil {
				t.Fatal("want error")
			}
		})
	}
	fsys := goFiles()
	fsys["wrong/x.py"] = &fstest.MapFile{Data: []byte("x")}
	if _, _, err := parse(t, goSpecYAML, fsys); err == nil || !strings.Contains(err.Error(), ".go") {
		t.Fatalf("wrong/x.py in a Go lesson: %v", err)
	}
	if _, _, err := parse(t, strings.Replace(specYAML, "  - kind: forbidden_calls\n    names: [sum]", "  - kind: io_cases\n    cases: [{ stdin: \"\", stdout: \"x\" }]", 1), files()); err == nil {
		t.Fatal("io_cases in a Python lesson: want error")
	}
}

func TestWasmPublicViewAndCounts(t *testing.T) {
	typ, ex := goExercise(t)
	pv, err := typ.PublicView(ex)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"SOLUTION_SECRET", "15\\n", "14\\nFizzBuzz"} {
		if strings.Contains(string(pv), secret) {
			t.Errorf("public view contains %q: %s", secret, pv)
		}
	}
	if !strings.Contains(string(pv), `"run":"server"`) || !strings.Contains(string(pv), `"check_counts":{"visible":3,"hidden":1}`) {
		t.Errorf("public view: %s", pv)
	}
}

func TestWasmGradeJob(t *testing.T) {
	typ, ex := goExercise(t)
	f := &fakeRunner{out: runner.Output{Stdout: `{"version":1,"forbidden":[],"program":{"stdout":"","exit_code":0},"cases":[{"stdout":"1\n2\nFizz\n"},{"stdout":"x"},{"stdout":"y"}]}` + "\n"}}
	res, err := typ.Grade(runner.NewContext(context.Background(), f), ex,
		exercise.Submission{Answer: json.RawMessage(`{"files":{"main.go":"package main"}}`)})
	if err != nil {
		t.Fatal(err)
	}
	if f.job.Image != "wasm:go" || len(res.Checks) != 4 || !res.Checks[0].Passed || res.Checks[1].Name != "Up to Buzz" || res.Checks[2].Visible {
		t.Fatalf("job %s, result %+v", f.job.Image, res)
	}
	var job struct {
		RunEntry bool     `json:"run_entry"`
		Cases    []string `json:"cases"`
	}
	_ = json.Unmarshal(f.job.Stdin, &job)
	if !job.RunEntry || len(job.Cases) != 3 || job.Cases[2] != "15\n" {
		t.Fatalf("job = %+v", job)
	}
}

// Server Run sends only visible cases; the hidden stdin never reaches the runner.
func TestWasmRunUsesVisibleCasesOnly(t *testing.T) {
	typ, ex := goExercise(t)
	f := &fakeRunner{out: runner.Output{Stdout: `{"version":1,"forbidden":[],"program":{"stdout":"1\n2\nFizz\n"},"cases":[{"stdout":"1\n2\nFizz\n"},{"stdout":"1\n2\nFizz\n4\nBuzz\n"}]}` + "\n"}}
	out, err := typ.(exercise.Runnable).Run(runner.NewContext(context.Background(), f), ex, json.RawMessage(`{"files":{"main.go":"package main"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(f.job.Stdin), `15\n`) {
		t.Fatal("hidden case stdin sent on Run")
	}
	// Two visible cases plus the exit_code row; the output pane shows case 1.
	if out.Stdout != "1\n2\nFizz\n" || len(out.Checks) != 3 || out.Checks[0].Name != "Sample 1" || out.Checks[1].Name != "Up to Buzz" ||
		!out.Checks[1].Passed || out.Checks[2].Name != "Exit code" {
		t.Fatalf("out = %+v", out)
	}
}

var (
	wasmOnce sync.Once
	wasmRun  *wasm.Runner
	wasmErr  error
)

func wasmOrSkip(t *testing.T) runner.Runner {
	t.Helper()
	wasmOnce.Do(func() {
		dir := os.Getenv("CQ_WASM_DIR")
		if dir == "" {
			dir = filepath.Join("..", "..", "..", "..", "runners", "wasm", "dist")
		}
		wasmRun, wasmErr = wasm.New(context.Background(), dir, filepath.Join(os.TempDir(), "codequest-wasm-cache"))
	})
	if wasmErr != nil {
		if os.Getenv("CQ_REQUIRE_WASM") == "1" {
			t.Fatal(wasmErr)
		}
		t.Skip(wasmErr)
	}
	return wasmRun
}

// End to end on the real wasm runner: the solution passes, the wrong answer
// fails only the hidden case, a compile error marks its line.
func TestWasmGradeEndToEnd(t *testing.T) {
	r := runner.NewContext(context.Background(), wasmOrSkip(t))
	typ, ex := goExercise(t)
	grade := func(src string) exercise.Result {
		ans, _ := json.Marshal(map[string]any{"files": map[string]string{"main.go": src}})
		res, err := typ.Grade(r, ex, exercise.Submission{Answer: ans})
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	if res := grade(goSolution); res.Status != exercise.StatusPass {
		t.Fatalf("solution: %+v", res)
	}
	wrong := goFiles()["wrong/no_fizzbuzz.go"]
	res := grade(string(wrong.Data))
	if res.Status != exercise.StatusFail || !res.Checks[0].Passed || res.Checks[2].Passed || res.Checks[2].Got != "" {
		t.Fatalf("wrong: %+v", res)
	}
	res = grade("package main\n\nfunc main() {\n\tundefinedThing()\n}\n")
	if res.Status != exercise.StatusFail || res.Checks[0].Line != 4 {
		t.Fatalf("compile error: %+v", res)
	}
}
