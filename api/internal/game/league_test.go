package game_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/testredis"
)

var week = time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC) // a Monday

// member puts user in a group with xp earned at the given time.
func member(t *testing.T, pool *pgxpool.Pool, user uuid.UUID, tier string, group, xp int, at time.Time) {
	t.Helper()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `INSERT INTO user_stats (user_id, league_tier) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET league_tier = $2`, user, tier); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO league_members (week, user_id, tier, group_no) VALUES ($1, $2, $3, $4)`, week, user, tier, group); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO xp_events (user_id, amount, reason, ref_id, created_at) VALUES ($1, $2, 'pass', $3, $4)`,
		user, xp, "x-"+user.String(), at); err != nil {
		t.Fatal(err)
	}
}

func tierOf(t *testing.T, pool *pgxpool.Pool, user uuid.UUID) string {
	t.Helper()
	var tier string
	if err := pool.QueryRow(context.Background(), `SELECT league_tier FROM user_stats WHERE user_id = $1`, user).Scan(&tier); err != nil {
		t.Fatal(err)
	}
	return tier
}

// The weekly job closes a week only after it ends (Monday 00:00 UTC), promotes
// the top 5 and demotes the bottom 5 of each group, and runs exactly once.
func TestWeeklyLeagueClose(t *testing.T) {
	pool := setup(t)
	clk := &clock{}
	l := &game.Leagues{Pool: pool, Now: clk.Now}
	ctx := context.Background()

	// Silver group of 12: user i earns 100 - i XP (user 0 first).
	silver := make([]uuid.UUID, 12)
	for i := range silver {
		silver[i] = newUser(t, pool, fmt.Sprintf("s%d@example.com", i), "UTC")
		member(t, pool, silver[i], "silver", 1, 100-i, week.Add(time.Duration(i)*time.Hour))
	}
	// Tie in a bronze group of 3: equal XP, whoever reached it first ranks higher.
	early := newUser(t, pool, "early@example.com", "UTC")
	late := newUser(t, pool, "late@example.com", "UTC")
	low := newUser(t, pool, "low@example.com", "UTC")
	member(t, pool, late, "bronze", 1, 50, week.Add(48*time.Hour))
	member(t, pool, early, "bronze", 1, 50, week.Add(24*time.Hour))
	member(t, pool, low, "bronze", 1, 10, week.Add(time.Hour))
	// Diamond group of 6: nobody moves up from diamond; the bottom 5 move down.
	diamond := make([]uuid.UUID, 6)
	for i := range diamond {
		diamond[i] = newUser(t, pool, fmt.Sprintf("d%d@example.com", i), "UTC")
		member(t, pool, diamond[i], "diamond", 1, 60-i, week.Add(time.Hour))
	}

	// Sunday 23:59: the week is still running.
	clk.Set(week.AddDate(0, 0, 7).Add(-time.Minute))
	if n, err := l.CloseDue(ctx); err != nil || n != 0 {
		t.Fatalf("before the week ends: closed %d, err %v", n, err)
	}
	if got := tierOf(t, pool, silver[0]); got != "silver" {
		t.Fatalf("tier changed early: %s", got)
	}

	// Monday 00:00 UTC: the week closes.
	clk.Set(week.AddDate(0, 0, 7))
	if n, err := l.CloseDue(ctx); err != nil || n != 1 {
		t.Fatalf("at the week's end: closed %d, err %v", n, err)
	}
	for i, u := range silver {
		want := "silver"
		switch {
		case i < 5:
			want = "gold"
		case i >= 7:
			want = "bronze"
		}
		if got := tierOf(t, pool, u); got != want {
			t.Errorf("silver rank %d: tier %s, want %s", i+1, got, want)
		}
	}
	// Bronze group of 3: all three are promoted (top 5); tie order is recorded.
	var earlyRank, lateRank int
	_ = pool.QueryRow(ctx, `SELECT final_rank FROM league_members WHERE user_id = $1`, early).Scan(&earlyRank)
	_ = pool.QueryRow(ctx, `SELECT final_rank FROM league_members WHERE user_id = $1`, late).Scan(&lateRank)
	if earlyRank != 1 || lateRank != 2 || tierOf(t, pool, low) != "silver" {
		t.Errorf("bronze: early rank %d, late rank %d, low tier %s", earlyRank, lateRank, tierOf(t, pool, low))
	}
	if tierOf(t, pool, diamond[0]) != "diamond" || tierOf(t, pool, diamond[1]) != "gold" || tierOf(t, pool, diamond[5]) != "gold" {
		t.Errorf("diamond: rank 1 %s, rank 6 %s; want diamond and gold (nobody above diamond; ranks 2-6 demoted)",
			tierOf(t, pool, diamond[0]), tierOf(t, pool, diamond[5]))
	}

	// Running again (another replica, the next tick) changes nothing.
	if _, err := pool.Exec(ctx, `UPDATE user_stats SET league_tier = 'silver' WHERE user_id = $1`, silver[0]); err != nil {
		t.Fatal(err)
	}
	if n, err := l.CloseDue(ctx); err != nil || n != 0 {
		t.Fatalf("second run: closed %d, err %v", n, err)
	}
	if got := tierOf(t, pool, silver[0]); got != "silver" {
		t.Fatalf("second run moved tiers again: %s", got)
	}
}

// The Redis board ranks exactly like Postgres, including ties, and heals
// itself when its key is missing.
func TestBoardMatchesDatabase(t *testing.T) {
	pool := setup(t)
	rdb, prefix := testredis.New(t)
	b := &game.Board{Redis: rdb, Prefix: prefix}
	ctx := context.Background()
	a := newUser(t, pool, "a@example.com", "UTC")
	bb := newUser(t, pool, "b@example.com", "UTC")
	c := newUser(t, pool, "c@example.com", "UTC")
	member(t, pool, a, "bronze", 1, 30, week.Add(3*time.Hour))
	member(t, pool, bb, "bronze", 1, 30, week.Add(1*time.Hour)) // same XP, earlier
	member(t, pool, c, "bronze", 1, 90, week.Add(5*time.Hour))
	ref := game.LeagueRef{Week: week, Tier: "bronze", Group: 1}

	// No key yet: Rank syncs from Postgres first.
	for user, want := range map[uuid.UUID]int{c: 1, bb: 2, a: 3} {
		rank, of, err := b.Rank(ctx, pool, ref, user)
		if err != nil || rank != want || of != 3 {
			t.Errorf("rank = %d of %d (err %v), want %d of 3", rank, of, err, want)
		}
	}
	// Syncing twice is harmless.
	if err := b.Sync(ctx, pool, ref); err != nil {
		t.Fatal(err)
	}
	if err := b.Sync(ctx, pool, ref); err != nil {
		t.Fatal(err)
	}
	if n, _ := rdb.ZCard(ctx, prefix+"cq:league:2026-11-02:bronze:1").Result(); n != 3 {
		t.Fatalf("board size = %d", n)
	}
}
