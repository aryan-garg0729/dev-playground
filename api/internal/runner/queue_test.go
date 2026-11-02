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
