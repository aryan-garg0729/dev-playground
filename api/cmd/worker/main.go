// Command worker runs grading jobs from the Redis queue in Docker containers.
//
// It needs the docker CLI and a Docker daemon, and nothing else: no database
// credentials. Run it on a host that runs no other services.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"

	"github.com/aryangarg/codequest/api/internal/config"
	"github.com/aryangarg/codequest/api/internal/runner"
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
	log.Info("worker started", "concurrency", cfg.Concurrency)
	w := &runner.Worker{
		Redis: rdb, Runner: runner.Docker{Binary: cfg.DockerBin},
		Concurrency: cfg.Concurrency, Log: log,
	}
	return w.Serve(ctx)
}
