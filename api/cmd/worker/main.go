// Command worker runs grading jobs from a Redis queue.
//
// WORKER_RUNNER=docker (default) runs the jobs queue in Docker containers: it
// needs the docker CLI and a Docker daemon. Run it on a host that runs no
// other services. WORKER_RUNNER=wasm runs the jobs:wasm queue as WASI modules
// in-process: it needs the modules (WASM_DIR) and no Docker socket, so it runs
// in a hardened container of its own. Neither has database credentials.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/aryangarg/codequest/api/internal/config"
	"github.com/aryangarg/codequest/api/internal/runner"
	"github.com/aryangarg/codequest/api/internal/runner/wasm"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(log); err != nil {
		log.Error("worker stopped", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.LoadWorker()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("REDIS_URL: %w", err)
	}
	rdb := redis.NewClient(opt)
	defer func() { _ = rdb.Close() }()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}
	w := &runner.Worker{
		Redis: rdb, Runner: runner.Docker{Binary: cfg.DockerBin},
		Concurrency: cfg.Concurrency, Log: log,
	}
	if cfg.Runner == "wasm" {
		start := time.Now()
		wr, err := wasm.New(ctx, cfg.WasmDir, cfg.WasmCacheDir)
		if err != nil {
			return err
		}
		defer func() { _ = wr.Close(context.Background()) }()
		log.Info("wasm modules ready", "dir", cfg.WasmDir, "ms", time.Since(start).Milliseconds())
		w.Runner, w.Queue = wr, runner.JobsKey(runner.WasmImage(""))
	}
	log.Info("worker started", "runner", cfg.Runner, "concurrency", cfg.Concurrency)
	return w.Serve(ctx)
}
