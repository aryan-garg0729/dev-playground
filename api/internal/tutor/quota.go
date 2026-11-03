package tutor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Quota counts hints per user per local day in Redis.
type Quota struct {
	Redis  *redis.Client
	Prefix string // default "cq:"
	Limit  int    // hints per day
}

func (q *Quota) key(user uuid.UUID, day time.Time) string {
	p := q.Prefix
	if p == "" {
		p = "cq:"
	}
	return p + "tutor:" + user.String() + ":" + day.Format(time.DateOnly)
}

// Used returns how many hints the user has taken on day.
func (q *Quota) Used(ctx context.Context, user uuid.UUID, day time.Time) (int, error) {
	n, err := q.Redis.Get(ctx, q.key(user, day)).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("tutor quota: %w", err)
	}
	return n, nil
}

// Take records one hint. It reports false (and records nothing) when the
// day's limit is used up, and returns the hints left after this one.
func (q *Quota) Take(ctx context.Context, user uuid.UUID, day time.Time) (bool, int, error) {
	k := q.key(user, day)
	pipe := q.Redis.TxPipeline()
	incr := pipe.Incr(ctx, k)
	pipe.Expire(ctx, k, 48*time.Hour) // a local day lasts at most ~26 hours
	if _, err := pipe.Exec(ctx); err != nil {
		return false, 0, fmt.Errorf("tutor quota: %w", err)
	}
	n := int(incr.Val())
	if n > q.Limit {
		if err := q.Redis.Decr(ctx, k).Err(); err != nil {
			return false, 0, fmt.Errorf("tutor quota: %w", err)
		}
		return false, 0, nil
	}
	return true, q.Limit - n, nil
}

// Refund gives back a hint that produced no answer (the model failed).
func (q *Quota) Refund(ctx context.Context, user uuid.UUID, day time.Time) error {
	if err := q.Redis.Decr(ctx, q.key(user, day)).Err(); err != nil {
		return fmt.Errorf("tutor quota: %w", err)
	}
	return nil
}
