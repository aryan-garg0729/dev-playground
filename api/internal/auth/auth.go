// Package auth handles sign-in: sessions, magic links and GitHub OAuth.
// Secrets (session ids, link tokens) are random 32-byte values; only their
// SHA-256 is stored, so a database leak does not leak live credentials.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/store"
)

// NewToken returns a random URL-safe token and its hash for storage.
func NewToken() (raw, hash string) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand: %v", err)) // never fails on supported platforms
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	return raw, Hash(raw)
}

// Hash is the stored form of a token.
func Hash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// Sessions are server-side sessions referenced by an opaque cookie.
type Sessions struct {
	Pool *pgxpool.Pool
	TTL  time.Duration // default 30 days
}

func (s *Sessions) ttl() time.Duration {
	if s.TTL <= 0 {
		return 30 * 24 * time.Hour
	}
	return s.TTL
}

// Create starts a session for user and returns the cookie value.
func (s *Sessions) Create(ctx context.Context, user uuid.UUID) (string, time.Time, error) {
	raw, hash := NewToken()
	exp := time.Now().Add(s.ttl())
	if err := store.New(s.Pool).CreateSession(ctx, store.CreateSessionParams{ID: hash, UserID: user, ExpiresAt: exp}); err != nil {
		return "", time.Time{}, fmt.Errorf("create session: %w", err)
	}
	return raw, exp, nil
}

// Lookup returns the session's user, or an invalid NullUUID for an unknown or
// expired session.
func (s *Sessions) Lookup(ctx context.Context, raw string) (uuid.NullUUID, error) {
	if raw == "" || len(raw) > 128 {
		return uuid.NullUUID{}, nil
	}
	id, err := store.New(s.Pool).GetSessionUser(ctx, Hash(raw))
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.NullUUID{}, nil
	}
	if err != nil {
		return uuid.NullUUID{}, fmt.Errorf("lookup session: %w", err)
	}
	return uuid.NullUUID{UUID: id, Valid: true}, nil
}

// Delete ends a session.
func (s *Sessions) Delete(ctx context.Context, raw string) error {
	if err := store.New(s.Pool).DeleteSession(ctx, Hash(raw)); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

// NormalizeEmail validates a bare address and lower-cases it.
func NormalizeEmail(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" || len(s) > 254 {
		return "", errors.New("enter a valid email address")
	}
	a, err := mail.ParseAddress(s)
	if err != nil || a.Address != s || a.Name != "" || !strings.Contains(s[strings.LastIndexByte(s, '@')+1:], ".") {
		return "", errors.New("enter a valid email address")
	}
	return strings.ToLower(s), nil
}

// SafeNext returns next if it is a same-site path ("/learn/x"), else "/".
// It rejects absolute and protocol-relative URLs and backslash tricks.
func SafeNext(next string) string {
	if next == "" || len(next) > 512 || next[0] != '/' || strings.HasPrefix(next, "//") ||
		strings.ContainsAny(next, "\\\r\n\t") || strings.HasPrefix(next, "/api/") {
		return "/"
	}
	return next
}

// UserForEmail returns the user with this (normalized) email, creating one.
func UserForEmail(ctx context.Context, db store.DBTX, email, name, avatar string) (store.GetUserRow, error) {
	q := store.New(db)
	u, err := q.GetUserByEmail(ctx, email)
	if err == nil {
		if err := q.UpdateUserProfile(ctx, store.UpdateUserProfileParams{ID: u.ID, Name: name, Avatar: avatar}); err != nil {
			return store.GetUserRow{}, fmt.Errorf("update profile: %w", err)
		}
		return store.GetUserRow(u), nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return store.GetUserRow{}, fmt.Errorf("find user: %w", err)
	}
	c, err := q.CreateUser(ctx, store.CreateUserParams{Email: email, Name: name, Avatar: avatar})
	if err != nil {
		return store.GetUserRow{}, fmt.Errorf("create user: %w", err)
	}
	if err := q.EnsureUserStats(ctx, c.ID); err != nil {
		return store.GetUserRow{}, fmt.Errorf("create stats: %w", err)
	}
	return store.GetUserRow(c), nil
}
