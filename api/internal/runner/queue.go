package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Queue is a Runner that hands jobs to cmd/worker through Redis and waits for
// the output. The worker needs Docker but no database access.
//
//	<prefix>jobs           list; the API LPUSHes, workers BRPOP
//	<prefix>result:<id>    list; the worker LPUSHes one reply, the API BLPOPs it
type Queue struct {
	Redis  *redis.Client
	Prefix string // key prefix, default "cq:"
	// MaxQueueWait is how long a job may sit in the queue before workers skip
	// it. Default 60s.
	MaxQueueWait time.Duration
}

type queuedJob struct {
	ID       string    `json:"id"`
	Job      Job       `json:"job"`
	Deadline time.Time `json:"deadline"` // workers drop the job after this
}

type reply struct {
	Output Output `json:"output"`
	Error  string `json:"error,omitempty"`
}

const replyTTL = 2 * time.Minute

func prefix(p string) string {
	if p == "" {
		return "cq:"
	}
	return p
}

// Run enqueues job and blocks until a worker replies or ctx ends.
func (q *Queue) Run(ctx context.Context, job Job) (Output, error) {
	wait := q.MaxQueueWait
	if wait <= 0 {
		wait = time.Minute
	}
	p := prefix(q.Prefix)
	qj := queuedJob{ID: uuid.NewString(), Job: job, Deadline: time.Now().Add(wait)}
	b, err := json.Marshal(qj)
	if err != nil {
		return Output{}, fmt.Errorf("queue: encode job: %w", err)
	}
	if err := q.Redis.LPush(ctx, p+"jobs", b).Err(); err != nil {
		return Output{}, fmt.Errorf("%w: enqueue: %w", ErrUnavailable, err)
	}
	// The reply can take up to the queue wait plus the run's own hard limit.
	l := job.Limits.withDefaults()
	timeout := wait + 2*time.Duration(l.TimeoutMs)*time.Millisecond + hardGrace
	if dl, ok := ctx.Deadline(); ok {
		timeout = min(timeout, time.Until(dl))
	}
	res, err := q.Redis.BLPop(ctx, timeout, p+"result:"+qj.ID).Result()
	if errors.Is(err, redis.Nil) {
		return Output{}, fmt.Errorf("%w: no worker replied within %s", ErrUnavailable, timeout.Round(time.Second))
	}
	if err != nil {
		return Output{}, fmt.Errorf("%w: wait for result: %w", ErrUnavailable, err)
	}
	var rep reply
	if err := json.Unmarshal([]byte(res[1]), &rep); err != nil {
		return Output{}, fmt.Errorf("queue: decode reply: %w", err)
	}
	if rep.Error != "" {
		return rep.Output, fmt.Errorf("%w: worker: %s", ErrUnavailable, rep.Error)
	}
	return rep.Output, nil
}

// Worker pulls jobs from the queue and runs them with Runner.
type Worker struct {
	Redis       *redis.Client
	Prefix      string
	Runner      Runner
	Concurrency int // default 2
	Log         *slog.Logger
}

// Serve runs until ctx ends, then waits for in-flight jobs.
func (w *Worker) Serve(ctx context.Context) error {
	n := w.Concurrency
	if n <= 0 {
		n = 2
	}
	log := w.Log
	if log == nil {
		log = slog.Default()
	}
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() { w.loop(ctx, log) })
	}
	wg.Wait()
	return nil
}

func (w *Worker) loop(ctx context.Context, log *slog.Logger) {
	p := prefix(w.Prefix)
	for ctx.Err() == nil {
		res, err := w.Redis.BRPop(ctx, time.Second, p+"jobs").Result()
		if errors.Is(err, redis.Nil) || ctx.Err() != nil {
			continue
		}
		if err != nil {
			log.Error("worker: dequeue", "err", err)
			select {
			case <-ctx.Done():
			case <-time.After(time.Second):
			}
			continue
		}
		var qj queuedJob
		if err := json.Unmarshal([]byte(res[1]), &qj); err != nil || qj.ID == "" {
			log.Error("worker: bad job", "err", err)
			continue
		}
		if time.Now().After(qj.Deadline) {
			log.Warn("worker: skipping expired job", "id", qj.ID)
			continue
		}
		// Finish a started job even during shutdown; the API is waiting for it.
		out, err := w.Runner.Run(context.WithoutCancel(ctx), qj.Job)
		rep := reply{Output: out}
		if err != nil {
			rep.Error = err.Error()
		}
		b, _ := json.Marshal(rep) // plain structs; cannot fail
		key := p + "result:" + qj.ID
		pipe := w.Redis.TxPipeline()
		pipe.LPush(context.WithoutCancel(ctx), key, b)
		pipe.Expire(context.WithoutCancel(ctx), key, replyTTL)
		if _, err := pipe.Exec(context.WithoutCancel(ctx)); err != nil {
			log.Error("worker: reply", "id", qj.ID, "err", err)
		}
		log.Info("worker: ran job", "id", qj.ID, "ms", out.RuntimeMs, "exit", out.ExitCode, "timed_out", out.TimedOut, "err", rep.Error)
	}
}
