package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/aryangarg/codequest/api/internal/testredis"
)

func TestAllow(t *testing.T) {
	rdb, prefix := testredis.New(t)
	now := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	l := &Limiter{Redis: rdb, Prefix: prefix, Now: func() time.Time { return now }}
	rule := Rule{Name: "submit", Limit: 3, Window: time.Minute}
	ctx := context.Background()
	for i := range 3 {
		ok, _, err := l.Allow(ctx, rule, "u1")
		if err != nil || !ok {
			t.Fatalf("event %d: ok=%v err=%v", i+1, ok, err)
		}
	}
	ok, retry, err := l.Allow(ctx, rule, "u1")
	if err != nil || ok {
		t.Fatalf("4th event: ok=%v err=%v, want refused", ok, err)
	}
	if retry <= 0 || retry > time.Minute {
		t.Fatalf("retry after %s", retry)
	}
	if ok, _, _ := l.Allow(ctx, rule, "u2"); !ok {
		t.Fatal("other key must have its own budget")
	}
	now = now.Add(time.Minute)
	if ok, _, _ := l.Allow(ctx, rule, "u1"); !ok {
		t.Fatal("next window must reset the count")
	}
}
