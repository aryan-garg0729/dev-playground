package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/store"
)

// ErrInvalidLink means the link is unknown, used or expired.
var ErrInvalidLink = errors.New("this sign-in link is invalid or has expired")

// Mailer sends one plain-text email.
type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

// MagicLinks issues and redeems single-use sign-in links.
//
// The link points at the web app (`/sign-in/verify#token=...`), not the API:
// the token travels in the URL fragment, so it never reaches a server log, and
// the page redeems it with a POST after a click, so mail scanners that fetch
// links cannot use it up.
type MagicLinks struct {
	Pool   *pgxpool.Pool
	Mailer Mailer
	AppURL string        // e.g. http://localhost:5173
	TTL    time.Duration // default 15 minutes
}

// Send emails a sign-in link to email (already normalized).
func (m *MagicLinks) Send(ctx context.Context, email, next string) error {
	ttl := m.TTL
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	raw, hash := NewToken()
	if err := store.New(m.Pool).CreateMagicLink(ctx, store.CreateMagicLinkParams{
		TokenHash: hash, Email: email, Next: SafeNext(next), ExpiresAt: time.Now().Add(ttl),
	}); err != nil {
		return fmt.Errorf("create magic link: %w", err)
	}
	link := strings.TrimRight(m.AppURL, "/") + "/sign-in/verify#token=" + raw
	body := "Sign in to CodeQuest:\n\n" + link + "\n\nThe link works once and expires in " +
		fmt.Sprintf("%d minutes", int(ttl.Minutes())) + ". If you did not ask for it, ignore this email.\n"
	if err := m.Mailer.Send(ctx, email, "Your CodeQuest sign-in link", body); err != nil {
		return fmt.Errorf("send magic link: %w", err)
	}
	return nil
}

// Redeem uses up a link and returns its email and next path.
func (m *MagicLinks) Redeem(ctx context.Context, raw string) (email, next string, err error) {
	if raw == "" || len(raw) > 128 {
		return "", "", ErrInvalidLink
	}
	row, err := store.New(m.Pool).UseMagicLink(ctx, Hash(raw))
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", ErrInvalidLink
	}
	if err != nil {
		return "", "", fmt.Errorf("redeem magic link: %w", err)
	}
	return row.Email, row.Next, nil
}

// LogMailer writes emails to the log instead of sending them (development only).
type LogMailer struct{ Log *slog.Logger }

// Send logs the message.
func (l LogMailer) Send(_ context.Context, to, subject, body string) error {
	l.Log.Warn("email not sent (no SMTP configured); development only", "to", to, "subject", subject, "body", body)
	return nil
}

// SMTPMailer sends through an SMTP relay (STARTTLS when the server offers it).
type SMTPMailer struct {
	Addr     string // host:port
	From     string
	Username string
	Password string
}

// Send delivers the message.
func (s SMTPMailer) Send(_ context.Context, to, subject, body string) error {
	if strings.ContainsAny(to+subject, "\r\n") {
		return errors.New("invalid header value")
	}
	msg := "From: " + s.From + "\r\nTo: " + to + "\r\nSubject: " + subject +
		"\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n" +
		strings.ReplaceAll(body, "\n", "\r\n")
	var a smtp.Auth
	if s.Username != "" {
		host := s.Addr
		if i := strings.LastIndexByte(host, ':'); i >= 0 {
			host = host[:i]
		}
		a = smtp.PlainAuth("", s.Username, s.Password, host)
	}
	if err := smtp.SendMail(s.Addr, a, s.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp: %w", err)
	}
	return nil
}
