// Command api is the CodeQuest HTTP server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata" // user timezones work without OS zoneinfo

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/config"
	"github.com/aryangarg/codequest/api/internal/economy"
	_ "github.com/aryangarg/codequest/api/internal/exercise/all" // register exercise types
	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/grading"
	"github.com/aryangarg/codequest/api/internal/httpapi"
	"github.com/aryangarg/codequest/api/internal/ratelimit"
	"github.com/aryangarg/codequest/api/internal/runner"
	"github.com/aryangarg/codequest/api/internal/store"
	"github.com/aryangarg/codequest/api/internal/tutor"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(log); err != nil {
		log.Error("api stopped", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	ropt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("REDIS_URL: %w", err)
	}
	rdb := redis.NewClient(ropt)
	defer func() { _ = rdb.Close() }()

	board := &game.Board{Redis: rdb}
	econ := &economy.Service{Pool: pool, Log: log}
	if err := econ.Sync(ctx); err != nil {
		return err
	}
	grader := grading.NewService(pool, &runner.Queue{Redis: rdb}, grading.Options{
		Log: log, Board: board, Game: &game.Engine{Rewards: econ},
	})
	defer grader.Close()
	// Anything still pending from before this start was lost with the old process.
	if n, err := grader.FailStale(ctx, 5*time.Minute); err != nil {
		log.Error("fail stale submissions", "err", err)
	} else if n > 0 {
		log.Warn("marked stale pending submissions as errors", "count", n)
	}

	var gh *auth.GitHub
	if cfg.GitHubClientID != "" {
		gh = &auth.GitHub{ClientID: cfg.GitHubClientID, ClientSecret: cfg.GitHubClientSecret, RedirectURL: cfg.AppURL + "/api/v1/auth/github/callback"}
	}
	var google *auth.Google
	if cfg.GoogleClientID != "" {
		google = &auth.Google{ClientID: cfg.GoogleClientID, ClientSecret: cfg.GoogleClientSecret, RedirectURL: cfg.AppURL + "/api/v1/auth/google/callback"}
	}
	if cfg.DevAuth {
		log.Warn("DEV_AUTH=1: POST /auth/dev signs anyone in without email; development only")
	}
	// The AI tutor is development-only for now (config refuses it in production).
	var tutorModel tutor.Provider
	var tutorQuota *tutor.Quota
	if cfg.Tutor == "ollama" {
		tutorModel = &tutor.Ollama{URL: cfg.OllamaURL, Name: cfg.TutorModel}
		tutorQuota = &tutor.Quota{Redis: rdb, Limit: cfg.TutorDailyLimit}
		log.Info("AI tutor on", "provider", cfg.Tutor, "model", cfg.TutorModel, "url", cfg.OllamaURL, "daily_limit", cfg.TutorDailyLimit)
	}

	// Close finished league weeks now and every minute (safe on many replicas).
	leagues := &game.Leagues{Pool: pool, Log: log, Rewards: econ}
	go leagues.Run(ctx, time.Minute)
	// Daily rollover: spends freezes on missed days, settles bets, evaluates medals.
	go econ.RunDaily(ctx, 10*time.Minute)
	go cleanAuth(ctx, pool, log)

	handler := httpapi.New(httpapi.Deps{
		Pool: pool, Grader: grader, Log: log,
		AllowedOrigins: cfg.AllowedOrigins,
		Sessions:       &auth.Sessions{Pool: pool},
		GitHub:         gh, Google: google, DevAuth: cfg.DevAuth, SecureCookies: cfg.SecureCookies(), Board: board, Economy: econ,
		Limiter: &ratelimit.Limiter{Redis: rdb},
		Tutor:   tutorModel, TutorQuota: tutorQuota,
	})
	if cfg.WebDir != "" {
		if handler, err = httpapi.WithWeb(handler, cfg.WebDir); err != nil {
			return fmt.Errorf("WEB_DIR: %w", err)
		}
	}
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	errc := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", cfg.HTTPAddr)
		errc <- srv.ListenAndServe()
	}()
	select {
	case err := <-errc:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve: %w", err)
		}
	case <-ctx.Done():
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}

// cleanAuth deletes expired sessions and CLI codes every hour, and
// analytics events and client errors older than 90 days.
func cleanAuth(ctx context.Context, pool *pgxpool.Pool, log *slog.Logger) {
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for {
		if err := store.New(pool).DeleteExpiredAuth(ctx); err != nil && ctx.Err() == nil {
			log.Error("clean expired sessions", "err", err)
		}
		if err := store.New(pool).DeleteExpiredCLIAuth(ctx); err != nil && ctx.Err() == nil {
			log.Error("clean expired cli tokens", "err", err)
		}
		if err := store.New(pool).DeleteOldTelemetry(ctx); err != nil && ctx.Err() == nil {
			log.Error("clean old telemetry", "err", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}
