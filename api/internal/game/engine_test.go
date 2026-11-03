package game_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/content"
	"github.com/aryangarg/codequest/api/internal/exercise"
	_ "github.com/aryangarg/codequest/api/internal/exercise/all"
	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/store"
	"github.com/aryangarg/codequest/api/internal/testdb"
)

// Exercises from content/testdata/valid.
var (
	pick = exercise.Exercise{ID: "demo.basics.pick", XP: 5}
	sum  = exercise.Exercise{ID: "demo.basics.sum", XP: 20}
)

var (
	pass = exercise.Result{Status: exercise.StatusPass, Score: 1}
	fail = exercise.Result{Status: exercise.StatusFail}
)

func setup(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool, _ := testdb.New(t)
	b, err := content.Load("../content/testdata/valid")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := content.Build(context.Background(), pool, b); err != nil {
		t.Fatal(err)
	}
	return pool
}

func newUser(t *testing.T, pool *pgxpool.Pool, email, tz string) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	u, err := auth.UserForEmail(ctx, pool, email, "User "+email, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.New(pool).SetUserTimezone(ctx, store.SetUserTimezoneParams{ID: u.ID, Timezone: tz}); err != nil {
		t.Fatal(err)
	}
	return u.ID
}

// clock is a settable fake clock.
type clock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *clock) Now() time.Time  { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *clock) Set(t time.Time) { c.mu.Lock(); defer c.mu.Unlock(); c.t = t }

func graded(t *testing.T, pool *pgxpool.Pool, e *game.Engine, user uuid.UUID, ex exercise.Exercise, r exercise.Result, attempts int) game.Outcome {
	t.Helper()
	var out game.Outcome
	err := pgx.BeginFunc(context.Background(), pool, func(tx pgx.Tx) error {
		var err error
		out, err = e.OnGraded(context.Background(), tx, user, ex, r, game.Progress{Attempts: attempts})
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestFirstPassAwardsXPOnce(t *testing.T) {
	pool := setup(t)
	clk := &clock{t: time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)}
	e := &game.Engine{Now: clk.Now}
	user := newUser(t, pool, "a@example.com", "UTC")

	if out := graded(t, pool, e, user, sum, fail, 1); out.XPAwarded != 0 {
		t.Fatalf("fail awarded %d XP", out.XPAwarded)
	}
	// Second attempt passes: no first-try bonus, +25% for no hint.
	out := graded(t, pool, e, user, sum, pass, 2)
	if out.XPAwarded != 25 || out.Breakdown == nil || *out.Breakdown != (game.Breakdown{Base: 20, NoHint: 5}) || out.FirstTry {
		t.Fatalf("first pass outcome = %+v (breakdown %+v)", out, out.Breakdown)
	}
	if out.TotalXP != 25 || out.Level != 1 || out.NextLevelXP != 100 || out.StreakDays != 1 || !out.StreakExtended {
		t.Fatalf("stats after first pass = %+v", out)
	}
	// Re-solving gives nothing, on the same day and later.
	if out := graded(t, pool, e, user, sum, pass, 3); out.XPAwarded != 0 || out.Breakdown != nil || out.TotalXP != 25 {
		t.Fatalf("re-solve outcome = %+v", out)
	}
	clk.Set(clk.Now().Add(48 * time.Hour))
	if out := graded(t, pool, e, user, sum, pass, 4); out.XPAwarded != 0 || out.TotalXP != 25 {
		t.Fatalf("re-solve days later = %+v", out)
	}

	// First-try pass of another lesson: +50% and +25%.
	out = graded(t, pool, e, user, pick, pass, 1)
	if out.XPAwarded != 5+3+1 || !out.FirstTry {
		t.Fatalf("first-try outcome = %+v (breakdown %+v)", out, out.Breakdown)
	}

	var events, total int
	if err := pool.QueryRow(context.Background(), `SELECT count(*), COALESCE(sum(amount),0) FROM xp_events WHERE user_id = $1`, user).Scan(&events, &total); err != nil {
		t.Fatal(err)
	}
	if events != 2 || total != 34 {
		t.Fatalf("ledger has %d events totalling %d, want 2 / 34", events, total)
	}
}

// Many passing submissions at once (double-clicks, retries) still award once:
// the ledger's unique key decides, and stats stay consistent with the ledger.
func TestConcurrentPassesAwardOnce(t *testing.T) {
	pool := setup(t)
	e := &game.Engine{}
	user := newUser(t, pool, "c@example.com", "UTC")
	var wg sync.WaitGroup
	errs := make(chan error, 10)
	for range 10 {
		wg.Go(func() {
			errs <- pgx.BeginFunc(context.Background(), pool, func(tx pgx.Tx) error {
				_, err := e.OnGraded(context.Background(), tx, user, sum, pass, game.Progress{Attempts: 1})
				return err
			})
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	var events, statsXP int
	ctx := context.Background()
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM xp_events WHERE user_id = $1`, user).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT total_xp FROM user_stats WHERE user_id = $1`, user).Scan(&statsXP); err != nil {
		t.Fatal(err)
	}
	if events != 1 || statsXP != 35 {
		t.Fatalf("events = %d, total_xp = %d; want 1 and 35", events, statsXP)
	}
}

func TestLevelUp(t *testing.T) {
	pool := setup(t)
	e := &game.Engine{}
	user := newUser(t, pool, "l@example.com", "UTC")
	if _, err := pool.Exec(context.Background(), `UPDATE user_stats SET total_xp = 90 WHERE user_id = $1`, user); err != nil {
		t.Fatal(err)
	}
	out := graded(t, pool, e, user, sum, pass, 1) // +35 → 125 XP
	if !out.LevelUp || out.Level != 2 || out.LevelStartXP != 100 || out.NextLevelXP != 283 {
		t.Fatalf("level outcome = %+v", out)
	}
}

// The streak follows the user's own calendar: two passes 1h apart around
// local midnight in Kolkata are two days; a missed day without a freeze resets.
func TestStreakInUserTimezone(t *testing.T) {
	pool := setup(t)
	clk := &clock{}
	e := &game.Engine{Now: clk.Now}
	user := newUser(t, pool, "k@example.com", "Asia/Kolkata")

	clk.Set(time.Date(2026, 11, 3, 18, 0, 0, 0, time.UTC)) // 23:30 IST, Nov 3
	if out := graded(t, pool, e, user, sum, pass, 1); out.StreakDays != 1 {
		t.Fatalf("day 1 streak = %d", out.StreakDays)
	}
	clk.Set(time.Date(2026, 11, 3, 19, 0, 0, 0, time.UTC)) // 00:30 IST, Nov 4
	if out := graded(t, pool, e, user, sum, pass, 2); out.StreakDays != 2 || !out.StreakExtended {
		t.Fatalf("day 2 streak = %+v", out)
	}
	clk.Set(time.Date(2026, 11, 6, 6, 0, 0, 0, time.UTC)) // Nov 6 IST: Nov 5 missed, no freeze
	if out := graded(t, pool, e, user, pick, pass, 1); out.StreakDays != 1 {
		t.Fatalf("after a missed day streak = %d, want 1", out.StreakDays)
	}
	var days int
	if err := pool.QueryRow(context.Background(), `SELECT count(DISTINCT day) FROM day_passes WHERE user_id = $1`, user).Scan(&days); err != nil {
		t.Fatal(err)
	}
	if days != 3 {
		t.Fatalf("active days = %d, want 3", days)
	}
}

func TestFailDoesNotTouchStreak(t *testing.T) {
	pool := setup(t)
	e := &game.Engine{}
	user := newUser(t, pool, "f@example.com", "UTC")
	out := graded(t, pool, e, user, sum, fail, 1)
	if out.StreakDays != 0 || out.StreakExtended || out.XPAwarded != 0 || out.League != nil {
		t.Fatalf("fail outcome = %+v", out)
	}
}

// XP joins the user to this week's league group; groups fill to 30.
func TestLeagueGroupsFillUp(t *testing.T) {
	pool := setup(t)
	clk := &clock{t: time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)} // a Wednesday
	e := &game.Engine{Now: clk.Now}
	groups := map[int]int{}
	for i := range game.LeagueGroupSize + 3 {
		user := newUser(t, pool, fmt.Sprintf("g%d@example.com", i), "UTC")
		out := graded(t, pool, e, user, pick, pass, 1)
		if out.League == nil || out.League.Tier != "bronze" || !out.League.Week.Equal(time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("league = %+v", out.League)
		}
		groups[out.League.Group]++
	}
	if groups[1] != game.LeagueGroupSize || groups[2] != 3 {
		t.Fatalf("group sizes = %v", groups)
	}
}
