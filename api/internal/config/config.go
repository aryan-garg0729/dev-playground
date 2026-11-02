// Package config loads settings from environment variables, once, in main.
package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds runtime settings.
type Config struct {
	DatabaseURL string
	HTTPAddr    string
	// DevUserEmail, when set, makes every request act as this user. Phase 1
	// stand-in for auth; Phase 4 removes it. Never set it in production.
	DevUserEmail   string
	AllowedOrigins []string
}

// Load reads the environment.
func Load() (Config, error) {
	c := Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		HTTPAddr:       getenv("HTTP_ADDR", ":8080"),
		DevUserEmail:   os.Getenv("DEV_USER_EMAIL"),
		AllowedOrigins: splitList(getenv("ALLOWED_ORIGINS", "http://localhost:5173")),
	}
	if c.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	return c, nil
}

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
