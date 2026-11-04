package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/store"
)

// codeAlphabet has no 0/O or 1/I/L, so a code read off the screen pastes cleanly.
const codeAlphabet = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"

const codeLen = 10 // ~49 bits; single use, 10 minutes, rate-limited per IP

// CLI issues login codes (shown in the browser) and bearer tokens (kept by
// the `codequest` CLI). Only hashes are stored.
type CLI struct {
	Pool     *pgxpool.Pool
	CodeTTL  time.Duration // default 10 minutes
	TokenTTL time.Duration // default 90 days, extended on use
}

func (c *CLI) codeTTL() time.Duration {
	if c.CodeTTL <= 0 {
		return 10 * time.Minute
	}
	return c.CodeTTL
}

func (c *CLI) tokenTTL() time.Duration {
	if c.TokenTTL <= 0 {
		return 90 * 24 * time.Hour
	}
	return c.TokenTTL
}

// NormalizeCode upper-cases a pasted code and drops everything outside the
// alphabet (spaces, the dash, a trailing newline).
func NormalizeCode(s string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(s) {
		if strings.ContainsRune(codeAlphabet, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func newCode() string {
	buf := make([]byte, codeLen)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("crypto/rand: %v", err))
	}
	out := make([]byte, codeLen)
	for i, b := range buf {
		// 256 % 31 != 0, so reject-sample to keep every letter equally likely.
		for int(b) >= 256-256%len(codeAlphabet) {
			var one [1]byte
			if _, err := rand.Read(one[:]); err != nil {
				panic(fmt.Sprintf("crypto/rand: %v", err))
			}
			b = one[0]
		}
		out[i] = codeAlphabet[int(b)%len(codeAlphabet)]
	}
	return string(out)
}

// FormatCode shows a code as XXXXX-XXXXX.
func FormatCode(code string) string {
	if len(code) != codeLen {
		return code
	}
	return code[:5] + "-" + code[5:]
}

// CreateLoginCode issues a code for a signed-in user, formatted for display.
func (c *CLI) CreateLoginCode(ctx context.Context, user uuid.UUID) (string, time.Time, error) {
	code := newCode()
	exp := time.Now().Add(c.codeTTL())
	if err := store.New(c.Pool).CreateCLILoginCode(ctx, store.CreateCLILoginCodeParams{CodeHash: Hash(code), UserID: user, ExpiresAt: exp}); err != nil {
		return "", time.Time{}, fmt.Errorf("create cli code: %w", err)
	}
	return FormatCode(code), exp, nil
}

// ErrInvalidCode is an unknown, used or expired login code.
var ErrInvalidCode = errors.New("invalid or expired code")

// Redeem trades a pasted code for a new token (returned raw, once).
func (c *CLI) Redeem(ctx context.Context, pasted, label string) (string, uuid.UUID, error) {
	code := NormalizeCode(pasted)
	if len(code) != codeLen {
		return "", uuid.Nil, ErrInvalidCode
	}
	var raw string
	var user uuid.UUID
	err := pgx.BeginFunc(ctx, c.Pool, func(tx pgx.Tx) error {
		q := store.New(tx)
		u, err := q.UseCLILoginCode(ctx, Hash(code))
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidCode
		}
		if err != nil {
			return fmt.Errorf("use cli code: %w", err)
		}
		var hash string
		raw, hash = NewToken()
		user = u
		if len(label) > 64 {
			label = label[:64]
		}
		return q.CreateCLIToken(ctx, store.CreateCLITokenParams{ID: hash, UserID: u, Label: label, ExpiresAt: time.Now().Add(c.tokenTTL())})
	})
	if err != nil {
		if errors.Is(err, ErrInvalidCode) {
			return "", uuid.Nil, err
		}
		return "", uuid.Nil, fmt.Errorf("redeem cli code: %w", err)
	}
	return raw, user, nil
}

// Lookup returns the token's user (extending its expiry), or an invalid
// NullUUID for an unknown, expired or revoked token.
func (c *CLI) Lookup(ctx context.Context, raw string) (uuid.NullUUID, error) {
	if raw == "" || len(raw) > 128 {
		return uuid.NullUUID{}, nil
	}
	u, err := store.New(c.Pool).TouchCLIToken(ctx, store.TouchCLITokenParams{ID: Hash(raw), ExpiresAt: time.Now().Add(c.tokenTTL())})
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.NullUUID{}, nil
	}
	if err != nil {
		return uuid.NullUUID{}, fmt.Errorf("lookup cli token: %w", err)
	}
	return uuid.NullUUID{UUID: u, Valid: true}, nil
}

// Revoke deletes a token.
func (c *CLI) Revoke(ctx context.Context, raw string) error {
	if err := store.New(c.Pool).DeleteCLIToken(ctx, Hash(raw)); err != nil {
		return fmt.Errorf("revoke cli token: %w", err)
	}
	return nil
}
