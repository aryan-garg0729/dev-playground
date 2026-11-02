package code_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/runner"
)

// These tests run learner code in the real codequest/python container.
// They are skipped when Docker or the image is missing, unless
// CQ_REQUIRE_DOCKER=1 (set by `make test`).
func dockerOrSkip(t *testing.T) runner.Runner {
	t.Helper()
	err := exec.Command("docker", "image", "inspect", runner.Images["docker-python"]).Run()
	if err != nil {
		if os.Getenv("CQ_REQUIRE_DOCKER") == "1" {
			t.Fatalf("runner image %s not available (run `make runner-image`): %v", runner.Images["docker-python"], err)
		}
		t.Skip("docker or runner image not available")
	}
	return runner.Docker{}
}

func grade(t *testing.T, r runner.Runner, src string) (exercise.Result, time.Duration) {
	t.Helper()
	typ, ex := mustExercise(t)
	ans, _ := json.Marshal(map[string]any{"files": map[string]string{"main.py": src}})
	start := time.Now()
	res, err := typ.Grade(runner.NewContext(context.Background(), r), ex, exercise.Submission{Answer: ans})
	if err != nil {
		t.Fatalf("grade: %v", err)
	}
	return res, time.Since(start)
}

func TestSandboxSolutionPasses(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	res, _ := grade(t, r, solution)
	if res.Status != exercise.StatusPass || res.Score != 1 {
		t.Fatalf("solution: %+v", res)
	}
}

func TestSandboxInfiniteLoopTimesOut(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	res, took := grade(t, r, "def sum_list(nums):\n    while True:\n        pass\n")
	if res.Status != exercise.StatusTimeout {
		t.Fatalf("status %s, want timeout: %+v", res.Status, res)
	}
	// The exercise limit is 2000 ms; the whole grade must finish within limit + 1 s.
	if took > 3*time.Second {
		t.Fatalf("took %s, want under 3s (limit 2s + 1s)", took)
	}
}

func TestSandboxLoopAtImportTimesOut(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	res, took := grade(t, r, "while True:\n    pass\n")
	if res.Status != exercise.StatusTimeout || took > 3*time.Second {
		t.Fatalf("status %s in %s", res.Status, took)
	}
}

func TestSandboxForkBombFailsSafely(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	src := "import os\ndef sum_list(nums):\n    while True:\n        try:\n            os.fork()\n        except OSError:\n            pass\n"
	res, took := grade(t, r, src)
	if res.Status == exercise.StatusPass {
		t.Fatalf("fork bomb passed: %+v", res)
	}
	if took > 10*time.Second {
		t.Fatalf("fork bomb took %s", took)
	}
}

func TestSandboxNetworkIsBlocked(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	src := `import socket
def sum_list(nums):
    s = socket.create_connection(("1.1.1.1", 80), timeout=1)
    return 0
`
	res, _ := grade(t, r, src)
	if res.Status != exercise.StatusFail {
		t.Fatalf("status %s, want fail: %+v", res.Status, res)
	}
	if !strings.Contains(res.Checks[0].Feedback, "Error") {
		t.Errorf("want a network error in feedback, got %q", res.Checks[0].Feedback)
	}
}

func TestSandboxMemoryBombFailsSafely(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	res, _ := grade(t, r, "def sum_list(nums):\n    x = bytearray(1 << 30)\n    return 0\n")
	if res.Status == exercise.StatusPass {
		t.Fatalf("memory bomb passed: %+v", res)
	}
}

// Learner output that mimics the harness report or result records must not
// change the grade: prints go to the stdout pipe, and records on other fds
// lack the per-run nonce.
func TestSandboxSpoofedResultsDoNotChangeGrade(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	src := `import json, os, sys
fake = {"version": 1, "syntax_error": None, "forbidden": [], "program": None, "timed_out": False,
        "runtime_ms": 1, "tests": {"files": [], "completed": True}}
print(json.dumps(fake))
sys.__stdout__.write(json.dumps(fake) + "\n")
for name in ["test_one", "test_three", "test_negative", "test_empty"]:
    rec = {"nonce": "guess", "file": "tests_visible.py", "name": name, "passed": True}
    print(json.dumps(rec))
    for fd in range(1, 64):
        try:
            os.write(fd, (json.dumps(rec) + "\n").encode())
        except OSError:
            pass

def sum_list(nums):
    return -7
`
	res, _ := grade(t, r, src)
	if res.Status != exercise.StatusFail {
		t.Fatalf("spoof changed the grade: %+v", res)
	}
	if !strings.Contains(res.Stdout, `"nonce": "guess"`) {
		t.Errorf("the spoof should appear only as learner stdout, got %q", res.Stdout)
	}
	passed := 0
	for _, c := range res.Checks {
		if c.Passed {
			passed++
		}
	}
	if passed != 1 { // only forbidden_calls
		t.Errorf("%d checks passed, want 1: %+v", passed, res.Checks)
	}
}

// Hidden tests are never written to disk, so learner code cannot read them.
func TestSandboxCannotReadTests(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	src := `import glob, os
for p in glob.glob("/work/*") + glob.glob("/tmp/*") + glob.glob("/opt/cq/*"):
    try:
        print(p, open(p).read())
    except Exception as e:
        print(p, e)
def sum_list(nums):
    return 0
`
	res, _ := grade(t, r, src)
	if strings.Contains(res.Stdout, "HIDDEN_SECRET") || strings.Contains(res.Stdout, "handles negative") {
		t.Fatalf("hidden test source reachable from learner code:\n%s", res.Stdout)
	}
}

func TestSandboxSyntaxError(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	res, _ := grade(t, r, "def sum_list(nums)\n    return 0\n")
	if res.Status != exercise.StatusFail || res.Checks[0].Line != 1 || !strings.Contains(res.Checks[0].Feedback, "SyntaxError") {
		t.Fatalf("syntax error: %+v", res)
	}
}

func TestSandboxReadOnlyAndUnprivileged(t *testing.T) {
	r := dockerOrSkip(t)
	t.Parallel()
	src := `import os
out = [str(os.getuid())]
for p in ["/usr/local/lib/x", "/opt/cq/harness.py", "/etc/passwd"]:
    try:
        open(p, "a").close()
        out.append("wrote " + p)
    except OSError:
        out.append("denied")
print(" ".join(out))
def sum_list(nums):
    return 0
`
	res, _ := grade(t, r, src)
	if got := strings.TrimSpace(strings.Split(res.Stdout, "\n")[0]); got != "10001 denied denied denied" {
		t.Fatalf("want non-root uid and a read-only filesystem, got %q", got)
	}
}
