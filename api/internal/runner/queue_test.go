package runner

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aryangarg/codequest/api/internal/testredis"
)

type echoRunner struct{ calls atomic.Int32 }

func (e *echoRunner) Run(_ context.Context, job Job) (Output, error) {
	e.calls.Add(1)
	if string(job.Stdin) == "boom" {
		return Output{}, errors.New("docker exploded")
	}
	return Output{Stdout: "ran " + string(job.Stdin), RuntimeMs: 7}, nil
}

func TestQueueRoundTrip(t *testing.T) {
	rdb, prefix := testredis.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inner := &echoRunner{}
	w := &Worker{Redis: rdb, Prefix: prefix, Runner: inner, Concurrency: 2}
	done := make(chan struct{})
	go func() { _ = w.Serve(ctx); close(done) }()

	q := &Queue{Redis: rdb, Prefix: prefix}
	out, err := q.Run(ctx, Job{Image: "x", Stdin: []byte("job-1")})
	if err != nil || out.Stdout != "ran job-1" || out.RuntimeMs != 7 {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	if _, err := q.Run(ctx, Job{Image: "x", Stdin: []byte("boom")}); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("worker error should surface as ErrUnavailable, got %v", err)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("worker did not stop")
	}
}

func TestQueueNoWorker(t *testing.T) {
	rdb, prefix := testredis.New(t)
	q := &Queue{Redis: rdb, Prefix: prefix, MaxQueueWait: 100 * time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	_, err := q.Run(ctx, Job{Image: "x", Limits: Limits{TimeoutMs: 100}})
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("err = %v, want ErrUnavailable", err)
	}
}

func TestWorkerSkipsExpiredJobs(t *testing.T) {
	rdb, prefix := testredis.New(t)
	q := &Queue{Redis: rdb, Prefix: prefix, MaxQueueWait: time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	_, _ = q.Run(ctx, Job{Image: "x"}) // nobody is serving yet: enqueued, then given up on
	time.Sleep(10 * time.Millisecond)

	inner := &echoRunner{}
	wctx, wcancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer wcancel()
	_ = (&Worker{Redis: rdb, Prefix: prefix, Runner: inner, Concurrency: 1}).Serve(wctx)
	if n := inner.calls.Load(); n != 0 {
		t.Fatalf("worker ran %d expired job(s)", n)
	}
}

func TestDockerArgsHardening(t *testing.T) {
	args := Docker{}.Args("cq-run-x", Job{Image: "img", Limits: Limits{MemoryMB: 256}})
	joined := " " + strings.Join(args, " ") + " "
	for _, want := range []string{
		" --network none ", " --read-only ", " --cap-drop ALL ", " --security-opt no-new-privileges ",
		" --memory 256m ", " --memory-swap 256m ", " --pids-limit 64 ", " --user 10001:10001 ", " --cpus 1 ",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("docker args missing %q: %s", want, joined)
		}
	}
	if args[len(args)-1] != "img" {
		t.Errorf("image must be last, got %v", args)
	}
	for _, a := range args {
		if a == "-v" || a == "--volume" || a == "--privileged" || a == "--mount" {
			t.Errorf("forbidden docker arg %q", a)
		}
	}
}

// Wasm jobs go to their own queue: a Docker worker never sees them, and a wasm
// worker refuses a Docker job pushed onto its queue.
func TestQueueRoutesWasmJobs(t *testing.T) {
	rdb, prefix := testredis.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	docker, wasm := &echoRunner{}, &echoRunner{}
	for _, w := range []*Worker{
		{Redis: rdb, Prefix: prefix, Runner: docker, Concurrency: 1},
		{Redis: rdb, Prefix: prefix, Queue: "jobs:wasm", Runner: wasm, Concurrency: 1},
	} {
		go func() { _ = w.Serve(ctx) }()
	}
	q := &Queue{Redis: rdb, Prefix: prefix}
	for range 3 {
		if _, err := q.Run(ctx, Job{Image: WasmImage("go"), Stdin: []byte("w")}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := q.Run(ctx, Job{Image: "codequest/python:3.12", Stdin: []byte("d")}); err != nil {
		t.Fatal(err)
	}
	if docker.calls.Load() != 1 || wasm.calls.Load() != 3 {
		t.Fatalf("docker ran %d, wasm ran %d; want 1 and 3", docker.calls.Load(), wasm.calls.Load())
	}

	// A Docker job forced onto the wasm queue gets an error reply, not a run.
	b := []byte(`{"id":"forced","job":{"image":"codequest/python:3.12"},"deadline":"2999-01-01T00:00:00Z"}`)
	if err := rdb.LPush(ctx, prefix+"jobs:wasm", b).Err(); err != nil {
		t.Fatal(err)
	}
	res, err := rdb.BLPop(ctx, 5*time.Second, prefix+"result:forced").Result()
	if err != nil || !strings.Contains(res[1], "on queue jobs:wasm") {
		t.Fatalf("reply = %v, %v", res, err)
	}
	if wasm.calls.Load() != 3 {
		t.Fatal("wasm worker ran a Docker job")
	}
}

func TestRoute(t *testing.T) {
	docker, wasm := &echoRunner{}, &echoRunner{}
	r := Route{Wasm: wasm, Docker: docker}
	_, _ = r.Run(context.Background(), Job{Image: "wasm:js"})
	_, _ = r.Run(context.Background(), Job{Image: "codequest/python:3.12"})
	if docker.calls.Load() != 1 || wasm.calls.Load() != 1 {
		t.Fatalf("docker %d wasm %d", docker.calls.Load(), wasm.calls.Load())
	}
	if _, err := (Route{Docker: docker}).Run(context.Background(), Job{Image: "wasm:go"}); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("missing side: %v", err)
	}
}
