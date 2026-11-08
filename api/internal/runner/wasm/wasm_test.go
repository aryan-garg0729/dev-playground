package wasm

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aryangarg/codequest/api/internal/checks"
	"github.com/aryangarg/codequest/api/internal/runner"
)

var (
	shared     *Runner
	sharedErr  error
	sharedOnce sync.Once
)

// testRunner loads the modules from CQ_WASM_DIR (default runners/wasm/dist).
// Without them the tests skip, unless CQ_REQUIRE_WASM=1 (make test sets it).
func testRunner(t *testing.T) *Runner {
	t.Helper()
	sharedOnce.Do(func() {
		dir := os.Getenv("CQ_WASM_DIR")
		if dir == "" {
			dir = filepath.Join("..", "..", "..", "..", "runners", "wasm", "dist")
		}
		cache := filepath.Join(os.TempDir(), "codequest-wasm-cache")
		shared, sharedErr = New(context.Background(), dir, cache)
	})
	if sharedErr != nil {
		if os.Getenv("CQ_REQUIRE_WASM") == "1" {
			t.Fatal(sharedErr)
		}
		t.Skip(sharedErr)
	}
	return shared
}

func run(t *testing.T, lang, src string, j Job, limits runner.Limits) *checks.Report {
	t.Helper()
	if j.Entry == "" {
		j.Entry = "main." + lang
	}
	j.Files = map[string]string{j.Entry: src}
	b, _ := json.Marshal(j)
	out, err := testRunner(t).Run(context.Background(), runner.Job{Image: runner.WasmImage(lang), Stdin: b, Limits: limits})
	if err != nil {
		t.Fatal(err)
	}
	rep, err := checks.ParseReport(out.Stdout)
	if err != nil {
		t.Fatalf("%v: %q", err, out.Stdout)
	}
	return rep
}

const goEcho = `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Println(strings.ToUpper(s.Text()))
	}
}
`

func TestProgramAndCases(t *testing.T) {
	for _, tc := range []struct{ lang, src string }{
		{"go", goEcho},
		{"js", `const s = require('fs').readFileSync(0, 'utf8'); for (const l of s.split('\n')) if (l) console.log(l.toUpperCase());`},
	} {
		t.Run(tc.lang, func(t *testing.T) {
			rep := run(t, tc.lang, tc.src, Job{RunEntry: true, Stdin: "hi\n", Cases: []string{"a\nb\n", ""}}, runner.Limits{})
			if rep.SyntaxError != nil || rep.Program == nil || rep.Program.Stdout != "HI\n" || rep.Program.ExitCode != 0 {
				t.Fatalf("program: %+v %+v", rep.Program, rep.SyntaxError)
			}
			if len(rep.Cases) != 2 || rep.Cases[0].Stdout != "A\nB\n" || rep.Cases[1].Stdout != "" {
				t.Fatalf("cases: %+v", rep.Cases)
			}
		})
	}
}

func TestSyntaxError(t *testing.T) {
	rep := run(t, "go", "package main\n\nfunc main() {\n\tfmt.Println(1)\n}\n", Job{RunEntry: true, Cases: []string{"", ""}}, runner.Limits{})
	if rep.SyntaxError == nil || rep.SyntaxError.Line != 4 || rep.Program != nil || len(rep.Cases) != 0 {
		t.Fatalf("go: %+v program=%+v cases=%d", rep.SyntaxError, rep.Program, len(rep.Cases))
	}
	rep = run(t, "js", "let x = ;\n", Job{RunEntry: true}, runner.Limits{})
	if rep.SyntaxError == nil || rep.SyntaxError.Line != 1 || !strings.Contains(rep.SyntaxError.Message, "SyntaxError") {
		t.Fatalf("js: %+v", rep.SyntaxError)
	}
}

func TestExitCodesAndErrors(t *testing.T) {
	rep := run(t, "go", "package main\n\nimport \"os\"\n\nfunc main() { os.Exit(4) }\n", Job{RunEntry: true}, runner.Limits{})
	if rep.Program.ExitCode != 4 {
		t.Fatalf("os.Exit: %+v", rep.Program)
	}
	rep = run(t, "go", "package main\n\nfunc main() {\n\tvar a []int\n\t_ = a[3]\n}\n", Job{RunEntry: true}, runner.Limits{})
	if rep.Program.ExitCode != 2 || !strings.Contains(rep.Program.Stderr, "panic") {
		t.Fatalf("panic: %+v", rep.Program)
	}
	rep = run(t, "js", "console.log('a'); throw new TypeError('nope')", Job{RunEntry: true}, runner.Limits{})
	if rep.Program.ExitCode != 1 || rep.Program.Stdout != "a\n" || !strings.Contains(rep.Program.Stderr, "TypeError: nope") {
		t.Fatalf("throw: %+v", rep.Program)
	}
}

// Learner code sees only its own file: no host files, no environment, no network.
func TestIsolation(t *testing.T) {
	src := `package main

import (
	"fmt"
	"os"
)

func main() {
	for _, p := range []string{"/etc/passwd", "/proc/self/environ", "../../../etc/hostname", "/work/main.go"} {
		_, err := os.ReadFile(p)
		fmt.Println(p, err == nil)
	}
	fmt.Println(len(os.Environ()))
	fmt.Println(os.WriteFile("/work/x", []byte("x"), 0o644) == nil)
}
`
	rep := run(t, "go", src, Job{RunEntry: true}, runner.Limits{})
	want := "/etc/passwd false\n/proc/self/environ false\n../../../etc/hostname false\n/work/main.go true\n0\nfalse\n"
	if rep.Program.Stdout != want {
		t.Fatalf("got %q, stderr %q", rep.Program.Stdout, rep.Program.Stderr)
	}
	for _, imp := range []string{"net", "net/http", "os/exec", "syscall", "unsafe"} {
		rep = run(t, "go", "package main\n\nimport _ \""+imp+"\"\n\nfunc main() {}\n", Job{RunEntry: true}, runner.Limits{})
		if rep.SyntaxError == nil || !strings.Contains(rep.SyntaxError.Message, "not available") {
			t.Errorf("import %s: %+v %+v", imp, rep.SyntaxError, rep.Program)
		}
	}
	rep = run(t, "js", "console.log(typeof fetch, typeof XMLHttpRequest); require('child_process')", Job{RunEntry: true}, runner.Limits{})
	if !strings.HasPrefix(rep.Program.Stdout, "undefined undefined") || rep.Program.ExitCode == 0 {
		t.Fatalf("js globals: %+v", rep.Program)
	}
}

// Abuse ends within the limits and the runner keeps serving.
func TestLimits(t *testing.T) {
	limits := runner.Limits{TimeoutMs: 500, MemoryMB: 64}
	start := time.Now()
	rep := run(t, "js", "while (true) {}", Job{RunEntry: true, TimeoutMs: 500, Cases: []string{"", "", "", ""}}, limits)
	if !rep.TimedOut || !rep.Program.TimedOut || len(rep.Cases) != 4 || !rep.Cases[3].TimedOut {
		t.Fatalf("loop: %+v", rep)
	}
	if d := time.Since(start); d > 2*time.Second {
		t.Fatalf("five looping phases took %s; the job budget is 2 x 500 ms", d)
	}
	rep = run(t, "js", "let a = []; while (true) a.push(new Array(1e6).fill(1));", Job{RunEntry: true, TimeoutMs: 3000}, runner.Limits{TimeoutMs: 3000, MemoryMB: 64})
	if rep.Program.ExitCode == 0 {
		t.Fatalf("memory bomb exited 0: %+v", rep.Program)
	}
	rep = run(t, "go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfor {\n\t\tfmt.Println(\"spam spam spam\")\n\t}\n}\n", Job{RunEntry: true, TimeoutMs: 500, OutputLimit: 1000}, limits)
	if len(rep.Program.Stdout) != 1000 || !rep.Program.TimedOut {
		t.Fatalf("flood: %d bytes, timed out %v", len(rep.Program.Stdout), rep.Program.TimedOut)
	}
	rep = run(t, "go", goEcho, Job{RunEntry: true, Stdin: "still ok\n"}, runner.Limits{})
	if rep.Program.Stdout != "STILL OK\n" {
		t.Fatalf("after abuse: %+v", rep.Program)
	}
}

func TestRejectsBadJobs(t *testing.T) {
	r := testRunner(t)
	for name, j := range map[string]Job{
		"two files":  {Files: map[string]string{"main.go": "", "x.go": ""}, Entry: "main.go"},
		"path entry": {Files: map[string]string{"../main.go": ""}, Entry: "../main.go"},
		"tests":      {Files: map[string]string{"main.go": ""}, Entry: "main.go", Tests: []json.RawMessage{[]byte(`{}`)}},
	} {
		b, _ := json.Marshal(j)
		if _, err := r.Run(context.Background(), runner.Job{Image: "wasm:go", Stdin: b}); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if _, err := r.Run(context.Background(), runner.Job{Image: "wasm:cobol", Stdin: []byte(`{}`)}); err == nil {
		t.Error("unknown language: no error")
	}
}

func TestConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			in := strings.Repeat("x", i+1) + "\n"
			rep := run(t, "go", goEcho, Job{RunEntry: true, Stdin: in}, runner.Limits{})
			if rep.Program.Stdout != strings.ToUpper(in) {
				t.Errorf("run %d: %q", i, rep.Program.Stdout)
			}
		})
	}
	wg.Wait()
}

// One allocation past the memory cap fails instead of growing past it.
func TestMemoryCap(t *testing.T) {
	src := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tb := make([]byte, 100<<20)\n\tb[len(b)-1] = 1\n\tfmt.Println(\"allocated\", len(b))\n}\n"
	rep := run(t, "go", src, Job{RunEntry: true}, runner.Limits{MemoryMB: 64})
	if rep.Program.ExitCode == 0 || strings.Contains(rep.Program.Stdout, "allocated") {
		t.Fatalf("100 MB under a 64 MB cap: %+v", rep.Program)
	}
	if rep.Program.Stderr != "Memory limit: your program used more than 64 MB\n" {
		t.Fatalf("stderr %q", rep.Program.Stderr)
	}
	rep = run(t, "go", src, Job{RunEntry: true}, runner.Limits{MemoryMB: 256})
	if rep.Program.ExitCode != 0 || !strings.Contains(rep.Program.Stdout, "allocated") {
		t.Fatalf("100 MB under a 256 MB cap: %+v", rep.Program)
	}
}

func TestLanguagesMatchRunner(t *testing.T) {
	if len(Languages) != len(runner.WasmLanguages) {
		t.Fatalf("runner.WasmLanguages %v, modules %v", runner.WasmLanguages, Languages)
	}
	for _, l := range runner.WasmLanguages {
		if Languages[l] == "" {
			t.Errorf("no module for %s", l)
		}
	}
}
