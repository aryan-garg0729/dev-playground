package economy_test

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/content"
	"github.com/aryangarg/codequest/api/internal/economy"
	"github.com/aryangarg/codequest/api/internal/exercise"
	_ "github.com/aryangarg/codequest/api/internal/exercise/all"
	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/testdb"
)

type clock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *clock) Now() time.Time  { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *clock) Set(t time.Time) { c.mu.Lock(); defer c.mu.Unlock(); c.t = t }

type fixture struct {
	t     *testing.T
	pool  *pgxpool.Pool
	clk   *clock
	econ  *economy.Service
	eng   *game.Engine
	user  uuid.UUID
	attem map[string]int
}

func setup(t *testing.T) *fixture {
	t.Helper()
	pool, _ := testdb.New(t)
	ctx := context.Background()
	b, err := content.Load("../content/testdata/valid")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := content.Build(ctx, pool, b); err != nil {
		t.Fatal(err)
	}
	clk := &clock{t: time.Date(2026, 11, 2, 9, 0, 0, 0, time.UTC)} // Monday
	econ := &economy.Service{Pool: pool, Now: clk.Now, Rand: func(int) int { return 0 }}
	if err := econ.Sync(ctx); err != nil {
		t.Fatal(err)
	}
	u, err := auth.UserForEmail(ctx, pool, "e@example.com", "E", "")
	if err != nil {
		t.Fatal(err)
	}
	return &fixture{t: t, pool: pool, clk: clk, econ: econ, eng: &game.Engine{Now: clk.Now, Rewards: econ}, user: u.ID, attem: map[string]int{}}
}

// pass grades a pass of an exercise for the user (as grading would).
func (f *fixture) grade(id string, status string) game.Outcome {
	f.t.Helper()
	f.attem[id]++
	var out game.Outcome
	err := pgx.BeginFunc(context.Background(), f.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(context.Background(), `INSERT INTO progress (user_id, exercise_id, status, attempts) VALUES ($1, $2, 'started', 1)
			ON CONFLICT (user_id, exercise_id) DO UPDATE SET attempts = progress.attempts + 1,
			status = CASE WHEN progress.status = 'passed' OR $3 THEN 'passed' ELSE 'started' END`, f.user, id, status == "pass"); err != nil {
			return err
		}
		var err error
		ex := exercise.Exercise{ID: id, XP: 5, Boss: id == "demo.basics.sum"}
		res := exercise.Result{Status: status, Checks: []exercise.CheckResult{{Name: "c", Passed: status == "pass", Visible: true}}}
		out, err = f.eng.OnGraded(context.Background(), tx, f.user, ex, res, game.Progress{Attempts: f.attem[id]})
		return err
	})
	if err != nil {
		f.t.Fatal(err)
	}
	return out
}

func (f *fixture) stats() (streak, freezes, bytes int) {
	f.t.Helper()
	if err := f.pool.QueryRow(context.Background(), `SELECT streak_days, streak_freezes, bytes FROM user_stats WHERE user_id = $1`, f.user).
		Scan(&streak, &freezes, &bytes); err != nil {
		f.t.Fatal(err)
	}
	return
}

func (f *fixture) exec(sql string, args ...any) {
	f.t.Helper()
	if _, err := f.pool.Exec(context.Background(), sql, args...); err != nil {
		f.t.Fatal(err)
	}
}

func (f *fixture) day(n int) { f.clk.Set(time.Date(2026, 11, 2+n, 9, 0, 0, 0, time.UTC)) }

// Criterion: a Streak Freeze is spent automatically on a missed day (by the
// daily rollover, with no pass needed), the streak survives, and a third
// freeze cannot be held.
func TestStreakFreezeIsAutomatic(t *testing.T) {
	f := setup(t)
	f.grade("demo.basics.pick", "pass") // Monday: streak 1
	f.exec(`UPDATE user_stats SET streak_freezes = 1 WHERE user_id = $1`, f.user)

	f.day(2) // Wednesday: Tuesday was missed
	if n, err := f.econ.RunRollovers(context.Background()); err != nil || n != 1 {
		t.Fatalf("rollover: %d %v", n, err)
	}
	if s, fr, _ := f.stats(); s != 1 || fr != 0 {
		t.Fatalf("after a missed day: streak %d, freezes %d (want 1, 0)", s, fr)
	}
	if out := f.grade("demo.basics.query", "pass"); out.StreakDays != 2 {
		t.Fatalf("streak after the freeze = %d, want 2 (it survived)", out.StreakDays)
	}
	// The rollover runs once per local day.
	if n, _ := f.econ.RunRollovers(context.Background()); n != 0 {
		t.Fatalf("second rollover the same day processed %d users", n)
	}
	// The database refuses a third held freeze, whatever the code path.
	if _, err := f.pool.Exec(context.Background(), `UPDATE user_stats SET streak_freezes = 3 WHERE user_id = $1`, f.user); err == nil {
		t.Fatal("a third freeze was stored")
	}
}

// Without freezes the streak is lost at the rollover, and Streak Repair brings
// it back within 48 hours (once a month).
func TestStreakRepair(t *testing.T) {
	f := setup(t)
	for d := range 4 { // a 4-day streak
		f.day(d)
		f.grade("demo.basics.pick", "pass")
	}
	f.day(5) // day 4 missed, no freeze
	if _, err := f.econ.RunRollovers(context.Background()); err != nil {
		t.Fatal(err)
	}
	if s, _, _ := f.stats(); s != 0 {
		t.Fatalf("streak after a missed day without freezes = %d", s)
	}
	f.exec(`UPDATE user_stats SET bytes = 3000 WHERE user_id = $1`, f.user)
	f.exec(`INSERT INTO byte_events (user_id, amount, reason, ref_id) VALUES ($1, 3000, 'test', 'seed')`, f.user)
	if _, err := f.econ.Use(context.Background(), f.user, economy.StreakRepair, economy.UseRequest{RequestID: "r1"}); err != nil {
		t.Fatal(err)
	}
	if s, _, b := f.stats(); s != 4 || b != 1000 {
		t.Fatalf("after repair: streak %d, bytes %d", s, b)
	}
	if out := f.grade("demo.basics.query", "pass"); out.StreakDays != 5 {
		t.Fatalf("pass after repair: streak %d, want 5", out.StreakDays)
	}
	_, err := f.econ.Use(context.Background(), f.user, economy.StreakRepair, economy.UseRequest{RequestID: "r2"})
	var limit economy.LimitError
	if !errors.As(err, &limit) {
		t.Fatalf("second repair: %v", err)
	}
}

// Criterion: medals are evaluated after a submission, the league close and
// the daily rollover, and each tier is awarded (and paid) once.
func TestMedalsEachTierOnce(t *testing.T) {
	f := setup(t)
	medals := func() []string {
		rows, _ := f.pool.Query(context.Background(), `SELECT medal_id || ':' || tier FROM user_medals WHERE user_id = $1 ORDER BY 1`, f.user)
		defer rows.Close()
		var out []string
		for rows.Next() {
			var s string
			_ = rows.Scan(&s)
			out = append(out, s)
		}
		return out
	}

	// After a submission: the 10th lesson earns Scholar bronze.
	f.exec(`INSERT INTO counters (user_id, key, value) VALUES ($1, 'lessons', 9)`, f.user)
	out := f.grade("demo.basics.pick", "pass")
	if out.Rewards == nil || !slices.ContainsFunc(out.Rewards.Medals, func(m game.MedalAward) bool { return m.ID == "scholar" && m.Tier == "bronze" }) {
		t.Fatalf("submission medals = %+v", out.Rewards)
	}
	// After the daily rollover: a longest streak set elsewhere is picked up.
	f.exec(`UPDATE user_stats SET longest_streak = 7 WHERE user_id = $1`, f.user)
	f.day(1)
	if _, err := f.econ.RunRollovers(context.Background()); err != nil {
		t.Fatal(err)
	}
	// After the league close: a top-3 week earns Podium and Contender.
	leagues := &game.Leagues{Pool: f.pool, Now: f.clk.Now, Rewards: f.econ}
	f.day(7) // next Monday
	if _, err := leagues.CloseDue(context.Background()); err != nil {
		t.Fatal(err)
	}
	got := medals()
	for _, want := range []string{"scholar:bronze", "streak:bronze", "podium:bronze", "contender:bronze"} {
		if !slices.Contains(got, want) {
			t.Errorf("missing %s in %v", want, got)
		}
	}
	// Re-evaluating everything awards nothing twice and pays nothing twice.
	_, _, before := f.stats()
	f.day(8)
	_, _ = f.econ.RunRollovers(context.Background())
	f.grade("demo.basics.query", "pass")
	if again := medals(); len(again) != len(got) {
		t.Fatalf("medals awarded twice: %v -> %v", got, again)
	}
	var paid int
	_ = f.pool.QueryRow(context.Background(), `SELECT count(*) FROM byte_events WHERE user_id = $1 AND reason = 'medal'`, f.user).Scan(&paid)
	if paid != len(got) {
		t.Fatalf("%d medal payouts for %d tiers", paid, len(got))
	}
	if _, _, after := f.stats(); after < before {
		t.Fatalf("balance went down: %d -> %d", before, after)
	}
}

// Double or Nothing: 7 days on goal pays 2x; a missed finished day loses.
func TestDoubleOrNothing(t *testing.T) {
	f := setup(t)
	f.exec(`UPDATE user_stats SET bytes = 500, daily_goal = 1 WHERE user_id = $1`, f.user)
	f.exec(`INSERT INTO byte_events (user_id, amount, reason, ref_id) VALUES ($1, 500, 'test', 'seed')`, f.user)
	if _, err := f.econ.Use(context.Background(), f.user, economy.DoubleOrNothing, economy.UseRequest{RequestID: "b1", Stake: 200}); err != nil {
		t.Fatal(err)
	}
	if _, err := f.econ.Use(context.Background(), f.user, economy.DoubleOrNothing, economy.UseRequest{RequestID: "b2", Stake: 100}); err == nil {
		t.Fatal("a second open bet was accepted")
	}
	for d := range 7 { // a lesson a day (goal 1); re-solves still count for the day
		f.day(d)
		f.grade("demo.basics.pick", "pass")
	}
	f.day(7)
	if _, err := f.econ.RunRollovers(context.Background()); err != nil {
		t.Fatal(err)
	}
	var won, sum int
	_ = f.pool.QueryRow(context.Background(), `SELECT COALESCE(sum(amount) FILTER (WHERE reason = 'bet_win'), 0), sum(amount)
		FROM byte_events WHERE user_id = $1`, f.user).Scan(&won, &sum)
	if _, _, b := f.stats(); won != 400 || b != sum {
		t.Fatalf("won bet paid %d (want 2x200), balance %d, ledger %d", won, b, sum)
	}

	// A new bet, then a missed day: lost at the next rollover.
	if _, err := f.econ.Use(context.Background(), f.user, economy.DoubleOrNothing, economy.UseRequest{RequestID: "b3", Stake: 100}); err != nil {
		t.Fatal(err)
	}
	f.day(9) // the bet started on day 7, which ended with no lesson
	if _, err := f.econ.RunRollovers(context.Background()); err != nil {
		t.Fatal(err)
	}
	var status string
	_ = f.pool.QueryRow(context.Background(), `SELECT status FROM bets WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`, f.user).Scan(&status)
	if status != "lost" {
		t.Fatalf("bet status = %s", status)
	}
}

// Quests count distinct lessons per day; the first pass of a lesson counts once.
func TestQuestProgress(t *testing.T) {
	f := setup(t)
	f.grade("demo.basics.pick", "pass")
	f.grade("demo.basics.pick", "pass") // re-solve, same day: no extra progress
	v, err := f.econ.Quests(context.Background(), f.user)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, q := range v.Daily {
		ids[q.ID] = true
		want := map[string]int{"lessons": 1, "first_try": 1, "xp": 9, "quiz": 1, "code": 0}[q.Event]
		if q.Progress != min(want, q.Target) {
			t.Errorf("%s progress = %d, want %d", q.ID, q.Progress, min(want, q.Target))
		}
	}
	if len(ids) != 3 {
		t.Fatalf("daily quests not distinct: %v", ids)
	}
	if v.Weekly.Event == "lessons" && v.Weekly.Progress != 1 {
		t.Errorf("weekly progress = %d", v.Weekly.Progress)
	}
}

// Half of all chests hold no power-up: the "" slot of ChestOdds pays only Bytes.
func TestChestOddsHalfEmpty(t *testing.T) {
	total, empty := 0, 0
	for _, o := range economy.ChestOdds {
		total += o.Weight
		if o.PowerUp == "" {
			empty += o.Weight
		}
		if o.PowerUp == economy.SkeletonKey {
			t.Fatal("a chest can drop the Skeleton Key")
		}
	}
	if total != 100 || empty != 50 {
		t.Fatalf("odds total %d, empty %d; want 100 and 50", total, empty)
	}
	f := setup(t)
	f.econ.Rand = func(n int) int { return n - 1 } // the last slot: no power-up
	out := f.grade("demo.basics.sum", "pass")      // the boss: opens the boss chest
	c := out.Rewards.Boss.Chest
	if c == nil || c.PowerUp != "" || c.Bytes != economy.BossBytes {
		t.Fatalf("boss chest = %+v", c)
	}
	var held int
	if err := f.pool.QueryRow(context.Background(), `SELECT COALESCE(sum(count), 0) FROM inventory WHERE user_id = $1`, f.user).Scan(&held); err != nil {
		t.Fatal(err)
	}
	if held != 0 {
		t.Fatalf("an empty chest added %d power-ups", held)
	}
}

// The "hands-on lesson" daily quest counts project_cli lessons, so CLI-only
// courses can finish it.
func TestHandsOnQuestCountsCLI(t *testing.T) {
	f := setup(t)
	f.exec(`UPDATE exercises SET type = 'project_cli' WHERE id = 'demo.basics.pick'`)
	for d := range 20 { // find a day the quest is drawn for this user
		f.day(d)
		v, err := f.econ.Quests(context.Background(), f.user)
		if err != nil {
			t.Fatal(err)
		}
		for _, q := range v.Daily {
			if q.ID != "d.code1" {
				continue
			}
			f.grade("demo.basics.pick", "pass")
			v, err = f.econ.Quests(context.Background(), f.user)
			if err != nil {
				t.Fatal(err)
			}
			for _, q := range v.Daily {
				if q.ID == "d.code1" && q.Progress != 1 {
					t.Fatalf("hands-on quest after a CLI pass = %d", q.Progress)
				}
			}
			return
		}
	}
	t.Fatal("d.code1 was never drawn in 20 days")
}
