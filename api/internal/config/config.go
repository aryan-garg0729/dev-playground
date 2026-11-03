// Package config loads settings from environment variables, once, in main.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds runtime settings.
type Config struct {
	DatabaseURL    string
	HTTPAddr       string
	AllowedOrigins []string
	RedisURL       string

	// Production (APP_ENV=production) refuses unsafe settings at startup.
	Production bool
	// AppURL is the web app's public URL (links in emails, OAuth redirect).
	AppURL string
	// DevAuth enables POST /auth/dev, which signs in without email. For local
	// development and browser tests only; refused in production.
	DevAuth bool

	SMTPAddr, SMTPUsername, SMTPPassword, MailFrom string
	GitHubClientID, GitHubClientSecret             string
}

// Worker holds cmd/worker settings.
type Worker struct {
	RedisURL    string
	Concurrency int    // WORKER_CONCURRENCY, default 2
	DockerBin   string // DOCKER_BIN, default "docker"
}

// LoadWorker reads the worker's environment.
func LoadWorker() (Worker, error) {
	w := Worker{RedisURL: getenv("REDIS_URL", "redis://localhost:6379/0"), DockerBin: getenv("DOCKER_BIN", "docker")}
	n, err := strconv.Atoi(getenv("WORKER_CONCURRENCY", "2"))
	if err != nil || n < 1 || n > 64 {
		return Worker{}, fmt.Errorf("WORKER_CONCURRENCY must be 1-64, got %q", os.Getenv("WORKER_CONCURRENCY"))
	}
	w.Concurrency = n
	return w, nil
}

// Load reads the environment.
func Load() (Config, error) {
	c := Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		HTTPAddr:           getenv("HTTP_ADDR", ":8080"),
		AllowedOrigins:     splitList(getenv("ALLOWED_ORIGINS", "http://localhost:5173")),
		RedisURL:           getenv("REDIS_URL", "redis://localhost:6379/0"),
		Production:         os.Getenv("APP_ENV") == "production",
		AppURL:             strings.TrimRight(getenv("APP_URL", "http://localhost:5173"), "/"),
		DevAuth:            os.Getenv("DEV_AUTH") == "1",
		SMTPAddr:           os.Getenv("SMTP_ADDR"),
		SMTPUsername:       os.Getenv("SMTP_USERNAME"),
		SMTPPassword:       os.Getenv("SMTP_PASSWORD"),
		MailFrom:           getenv("MAIL_FROM", "CodeQuest <no-reply@codequest.local>"),
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}
	if c.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if (c.GitHubClientID == "") != (c.GitHubClientSecret == "") {
		return Config{}, errors.New("set both GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET, or neither")
	}
	if c.Production {
		switch {
		case c.DevAuth:
			return Config{}, errors.New("DEV_AUTH must not be set in production")
		case c.SMTPAddr == "":
			return Config{}, errors.New("SMTP_ADDR is required in production (magic links)")
		case !strings.HasPrefix(c.AppURL, "https://"):
			return Config{}, errors.New("APP_URL must be https in production")
		}
	}
	return c, nil
}

// SecureCookies reports whether cookies need the Secure flag.
func (c Config) SecureCookies() bool { return strings.HasPrefix(c.AppURL, "https://") }

// DatabaseURL returns DATABASE_URL or an error, for tools that need only that.
func DatabaseURL() (string, error) {
	u := os.Getenv("DATABASE_URL")
	if u == "" {
		return "", errors.New("DATABASE_URL is required")
	}
	return u, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func splitList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
