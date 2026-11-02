// Package testredis gives tests a Redis client and a key prefix of their own.
//
// Set TEST_REDIS_URL (e.g. redis://localhost:6379/15). Without it tests are
// skipped, unless CQ_REQUIRE_DB=1 (set by `make test`), in which case they fail.
package testredis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// New returns a client and a unique key prefix; keys under it are deleted when the test ends.
func New(t *testing.T) (*redis.Client, string) {
	t.Helper()
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		if os.Getenv("CQ_REQUIRE_DB") == "1" {
			t.Fatal("TEST_REDIS_URL is not set (run `make test`, which starts Redis)")
		}
		t.Skip("TEST_REDIS_URL not set; skipping Redis test")
	}
	opt, err := redis.ParseURL(url)
	if err != nil {
		t.Fatalf("parse TEST_REDIS_URL: %v", err)
	}
	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("ping redis: %v", err)
	}
	prefix := "cqtest:" + uuid.NewString() + ":"
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		iter := rdb.Scan(ctx, 0, prefix+"*", 100).Iterator()
		for iter.Next(ctx) {
			rdb.Del(ctx, iter.Val())
		}
		_ = rdb.Close()
	})
	return rdb, prefix
}
