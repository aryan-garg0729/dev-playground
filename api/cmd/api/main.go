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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/config"
	_ "github.com/aryangarg/codequest/api/internal/exercise/all" // register exercise types
	"github.com/aryangarg/codequest/api/internal/grading"
	"github.com/aryangarg/codequest/api/internal/httpapi"
	"github.com/aryangarg/codequest/api/internal/store"
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

	var devUser uuid.NullUUID
	if cfg.DevUserEmail != "" {
		u, err := store.New(pool).EnsureUserByEmail(ctx, store.EnsureUserByEmailParams{Email: cfg.DevUserEmail, Name: "Dev User"})
		if err != nil {
			return fmt.Errorf("ensure dev user: %w", err)
		}
		devUser = uuid.NullUUID{UUID: u.ID, Valid: true}
		log.Warn("DEV_USER_EMAIL is set: every request acts as this user; never use in production", "email", cfg.DevUserEmail)
	}

	srv := &http.Server{
		Addr: cfg.HTTPAddr,
		Handler: httpapi.New(httpapi.Deps{
			Pool: pool, Grader: grading.NewService(pool), Log: log,
			AllowedOrigins: cfg.AllowedOrigins, DevUser: devUser,
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}
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
