// Package ratelimit is a fixed-window counter in Redis.
package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Rule is a limit of Limit events per Window.
type Rule struct {
	Name   string
	Limit  int
	Window time.Duration
}

// Limiter counts events per rule and key.
type Limiter struct {
	Redis  *redis.Client
	Prefix string // default "cq:"
	Now    func() time.Time
}

// Allow records one event for key under rule. It reports whether the event is
// within the limit and, if not, how long until the window resets.
func (l *Limiter) Allow(ctx context.Context, rule Rule, key string) (bool, time.Duration, error) {
	now := time.Now()
	if l.Now != nil {
		now = l.Now()
	}
	p := l.Prefix
	if p == "" {
		p = "cq:"
	}
	win := now.UnixNano() / int64(rule.Window)
	k := p + "rl:" + rule.Name + ":" + key + ":" + strconv.FormatInt(win, 10)
	pipe := l.Redis.TxPipeline()
	incr := pipe.Incr(ctx, k)
	pipe.Expire(ctx, k, rule.Window+time.Second)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, 0, fmt.Errorf("ratelimit: %w", err)
	}
	if incr.Val() <= int64(rule.Limit) {
		return true, 0, nil
	}
	reset := time.Unix(0, (win+1)*int64(rule.Window))
	return false, reset.Sub(now), nil
}
