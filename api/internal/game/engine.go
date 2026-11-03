package game

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/store"
)

// Engine applies the game rules to graded submissions.
type Engine struct {
	// Now is the clock (time.Now when nil); tests pass a fake.
	Now func() time.Time
	// Rewards, when set, runs after the XP and streak updates in the same
	// transaction (the economy: Bytes, quests, medals, boss chests).
	Rewards Rewards
}

// Rewards reacts to a graded submission inside the grading transaction. It
// sees the Result but cannot change it.
type Rewards interface {
	OnGraded(ctx context.Context, tx pgx.Tx, ev Graded) (*RewardsOutcome, error)
}

// Graded is what the rewards hook gets about a submission.
type Graded struct {
	User      uuid.UUID
	Exercise  exercise.Exercise
	Result    exercise.Result
	Outcome   Outcome   // XP and streak changes so far
	FirstPass bool      // this pass earned the exercise's XP
	NewToday  bool      // first pass of this exercise on this local day
	LocalNow  time.Time // now, in the user's timezone
	Today     time.Time // local date (midnight UTC)
}

// Gain is one Byte reward.
type Gain struct {
	Reason string `json:"reason"`
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

// MedalAward is a newly earned medal tier.
type MedalAward struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Tier  string `json:"tier"`
	Title string `json:"title"`
}

// ChestDrop is what a chest held.
type ChestDrop struct {
	Bytes   int    `json:"bytes"`
	PowerUp string `json:"powerup"`
}

// BossResult is a boss fight's state after a submission.
type BossResult struct {
	MaxHP           int        `json:"max_hp"`
	HP              int        `json:"hp"`
	Defeated        bool       `json:"defeated"`
	Chest           *ChestDrop `json:"chest"`
	UnlockedChapter string     `json:"unlocked_chapter"`
}

// RewardsOutcome is what the economy gave for a submission.
type RewardsOutcome struct {
	Bytes           []Gain       `json:"bytes"`
	Medals          []MedalAward `json:"medals"`
	QuestsDone      []string     `json:"quests_done"`
	ChestReady      bool         `json:"chest_ready"`
	ChapterComplete bool         `json:"chapter_complete"`
	Boss            *BossResult  `json:"boss"`
}

func (e *Engine) now() time.Time {
	if e == nil || e.Now == nil {
		return time.Now()
	}
	return e.Now()
}

// Progress is the user's progress on the exercise after this attempt.
type Progress struct {
	Attempts int // including this one
	UsedHint bool
	Peeked   bool // Peek used on this exercise: never a first-try bonus
}

// LeagueRef names a league group.
type LeagueRef struct {
	Week  time.Time `json:"week"`
	Tier  string    `json:"tier"`
	Group int       `json:"group"`
}

// Outcome is what one graded submission did to the user's game state.
type Outcome struct {
	XPAwarded      int        `json:"xp_awarded"`
	Breakdown      *Breakdown `json:"breakdown"` // set only when XP was awarded
	FirstTry       bool       `json:"first_try"`
	SecondWindKept bool       `json:"second_wind_kept"` // first-try bonus kept on the second attempt
	FocusBoost     bool       `json:"focus_boost"`
	TotalXP        int        `json:"total_xp"`
	Level          int        `json:"level"`
	LevelUp        bool       `json:"level_up"`
	LevelStartXP   int        `json:"level_start_xp"`
	NextLevelXP    int        `json:"next_level_xp"`
	StreakDays     int        `json:"streak_days"`
	StreakExtended bool       `json:"streak_extended"`
	FreezesUsed    int        `json:"freezes_used"`
	FreezeEarned   bool       `json:"freeze_earned"`
	// Rewards is the economy's part (nil when there is none).
	Rewards *RewardsOutcome `json:"rewards"`
	// League is the group the user's weekly XP counts in (set when XP was
	// awarded), for the leaderboard update after commit.
	League *LeagueRef `json:"-"`
}

// Location returns the IANA zone name's location, or UTC.
func Location(name string) *time.Location {
	if loc, err := time.LoadLocation(name); err == nil && name != "" {
		return loc
	}
	return time.UTC
}

// OnGraded runs in the grading transaction, after the submission and progress
// rows are written. It reads only the Result, the exercise's XP, the user's
// history and their active power-ups; it never changes the Result.
//
// XP: only the first pass of an exercise earns XP (the ledger's unique key
// enforces it, also under concurrency). Streak: any pass counts for its local day.
func (e *Engine) OnGraded(ctx context.Context, tx pgx.Tx, user uuid.UUID, ex exercise.Exercise, r exercise.Result, p Progress) (out Outcome, err error) {
	q := store.New(tx)
	if err := q.EnsureUserStats(ctx, user); err != nil {
		return Outcome{}, fmt.Errorf("ensure stats: %w", err)
	}
	st, err := q.LockUserStats(ctx, user)
	if err != nil {
		return Outcome{}, fmt.Errorf("lock stats: %w", err)
	}
	now := e.now()
	loc := Location(st.Timezone)
	today := LocalDate(now, loc)
	streak := Streak{Days: int(st.StreakDays), Freezes: int(st.StreakFreezes)}
	if st.LastActiveDate != nil {
		streak.Last = *st.LastActiveDate
	}
	out = Outcome{TotalXP: int(st.TotalXp), Level: int(st.Level), StreakDays: Current(streak, today)}
	defer func() { out.LevelStartXP, out.NextLevelXP = LevelStart(out.Level), LevelStart(out.Level+1) }()

	// Second Wind covers exactly the attempt after one failed submit.
	secondWind, err := q.EffectActive(ctx, store.EffectActiveParams{UserID: user, Effect: "second_wind", RefID: ex.ID, Now: now})
	if err != nil {
		return Outcome{}, fmt.Errorf("second wind: %w", err)
	}
	if secondWind && p.Attempts >= 2 {
		if err := q.ClearEffect(ctx, store.ClearEffectParams{UserID: user, Effect: "second_wind", RefID: ex.ID}); err != nil {
			return Outcome{}, fmt.Errorf("clear second wind: %w", err)
		}
	}

	ev := Graded{User: user, Exercise: ex, Result: r, LocalNow: now.In(loc), Today: today}
	if r.Status == exercise.StatusPass {
		if err := e.pass(ctx, q, &out, &ev, st.LeagueTier, streak, p, secondWind, now); err != nil {
			return Outcome{}, err
		}
	}
	if e.Rewards != nil {
		ev.Outcome = out
		rw, err := e.Rewards.OnGraded(ctx, tx, ev)
		if err != nil {
			return Outcome{}, fmt.Errorf("rewards: %w", err)
		}
		out.Rewards = rw
	}
	return out, nil
}

func (e *Engine) pass(ctx context.Context, q *store.Queries, out *Outcome, ev *Graded, tier string, streak Streak, p Progress, secondWind bool, now time.Time) error {
	user, ex, today := ev.User, ev.Exercise, ev.Today
	_, err := q.RecordDayPass(ctx, store.RecordDayPassParams{UserID: user, Day: today, ExerciseID: ex.ID})
	switch {
	case err == nil:
		ev.NewToday = true
	case !errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("record day: %w", err)
	}
	change := ApplyPass(streak, today)
	out.StreakDays, out.StreakExtended = change.Days, change.Extended
	out.FreezesUsed, out.FreezeEarned = change.FreezesUsed, change.FreezeEarned
	if change.Reset {
		if err := q.SetLostStreak(ctx, store.SetLostStreakParams{UserID: user, LostStreakDays: int32(streak.Days), LostOn: &today}); err != nil {
			return fmt.Errorf("lost streak: %w", err)
		}
	}

	out.FirstTry = !p.Peeked && (p.Attempts == 1 || (p.Attempts == 2 && secondWind))
	out.SecondWindKept = out.FirstTry && p.Attempts == 2
	b := AwardFor(ex.XP, out.FirstTry, p.UsedHint)
	if ex.Boss {
		b = b.WithBoss()
	}
	focus, err := q.EffectActive(ctx, store.EffectActiveParams{UserID: user, Effect: "focus_boost", RefID: "", Now: now})
	if err != nil {
		return fmt.Errorf("focus boost: %w", err)
	}
	if focus {
		b = b.WithBoost()
		out.FocusBoost = true
	}
	_, err = q.AwardXP(ctx, store.AwardXPParams{UserID: user, Amount: int32(b.Total()), Reason: ReasonPass, RefID: ex.ID})
	switch {
	case errors.Is(err, pgx.ErrNoRows): // already passed before: no XP
		out.FirstTry, out.SecondWindKept, out.FocusBoost = false, false, false
	case err != nil:
		return fmt.Errorf("award xp: %w", err)
	default:
		ev.FirstPass = true
		out.XPAwarded, out.Breakdown = b.Total(), &b
		out.TotalXP += b.Total()
		if lvl := LevelFor(out.TotalXP); lvl > out.Level {
			out.Level, out.LevelUp = lvl, true
		}
		if out.FirstTry {
			if err := q.SetFirstTry(ctx, store.SetFirstTryParams{UserID: user, ExerciseID: ex.ID}); err != nil {
				return fmt.Errorf("first try: %w", err)
			}
		}
		ref, err := joinLeague(ctx, q, user, tier, WeekStart(now))
		if err != nil {
			return err
		}
		out.League = &ref
	}

	last := change.Last
	if err := q.UpdateUserStats(ctx, store.UpdateUserStatsParams{
		UserID: user, TotalXp: int32(out.TotalXP), Level: int32(out.Level),
		StreakDays: int32(change.Days), StreakFreezes: int32(change.Freezes), LastActiveDate: &last,
	}); err != nil {
		return fmt.Errorf("update stats: %w", err)
	}
	if err := q.SetLongestStreak(ctx, store.SetLongestStreakParams{UserID: user, Days: int32(change.Days)}); err != nil {
		return fmt.Errorf("longest streak: %w", err)
	}
	return nil
}

// joinLeague returns the user's group for week, assigning one on the first XP
// of the week: the lowest-numbered group of their tier with room, else a new one.
func joinLeague(ctx context.Context, q *store.Queries, user uuid.UUID, tier string, week time.Time) (LeagueRef, error) {
	m, err := q.GetLeagueMember(ctx, store.GetLeagueMemberParams{Week: week, UserID: user})
	if err == nil {
		return LeagueRef{Week: week, Tier: m.Tier, Group: int(m.GroupNo)}, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return LeagueRef{}, fmt.Errorf("league member: %w", err)
	}
	if err := q.LockLeagueKey(ctx, "league-join:"+week.Format(time.DateOnly)+":"+tier); err != nil {
		return LeagueRef{}, fmt.Errorf("lock league: %w", err)
	}
	group, err := q.OpenLeagueGroup(ctx, store.OpenLeagueGroupParams{Week: week, Tier: tier, Size: LeagueGroupSize})
	if errors.Is(err, pgx.ErrNoRows) {
		maxGroup, err := q.MaxLeagueGroup(ctx, store.MaxLeagueGroupParams{Week: week, Tier: tier})
		if err != nil {
			return LeagueRef{}, fmt.Errorf("league groups: %w", err)
		}
		group = maxGroup + 1
	} else if err != nil {
		return LeagueRef{}, fmt.Errorf("open league group: %w", err)
	}
	if err := q.InsertLeagueMember(ctx, store.InsertLeagueMemberParams{Week: week, UserID: user, Tier: tier, GroupNo: group}); err != nil {
		return LeagueRef{}, fmt.Errorf("join league: %w", err)
	}
	return LeagueRef{Week: week, Tier: tier, Group: int(group)}, nil
}
