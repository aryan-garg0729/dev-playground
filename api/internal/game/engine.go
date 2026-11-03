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
	TotalXP        int        `json:"total_xp"`
	Level          int        `json:"level"`
	LevelUp        bool       `json:"level_up"`
	LevelStartXP   int        `json:"level_start_xp"`
	NextLevelXP    int        `json:"next_level_xp"`
	StreakDays     int        `json:"streak_days"`
	StreakExtended bool       `json:"streak_extended"`
	FreezesUsed    int        `json:"freezes_used"`
	FreezeEarned   bool       `json:"freeze_earned"`
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
// rows are written. It reads only the Result, the exercise's XP and the user's
// history; it never changes the Result.
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
	today := LocalDate(now, Location(st.Timezone))
	streak := Streak{Days: int(st.StreakDays), Freezes: int(st.StreakFreezes)}
	if st.LastActiveDate != nil {
		streak.Last = *st.LastActiveDate
	}
	out = Outcome{TotalXP: int(st.TotalXp), Level: int(st.Level), StreakDays: Current(streak, today)}
	defer func() { out.LevelStartXP, out.NextLevelXP = LevelStart(out.Level), LevelStart(out.Level+1) }()
	if r.Status != exercise.StatusPass {
		return out, nil
	}

	if err := q.RecordDayPass(ctx, store.RecordDayPassParams{UserID: user, Day: today, ExerciseID: ex.ID}); err != nil {
		return Outcome{}, fmt.Errorf("record day: %w", err)
	}
	change := ApplyPass(streak, today)
	out.StreakDays, out.StreakExtended = change.Days, change.Extended
	out.FreezesUsed, out.FreezeEarned = change.FreezesUsed, change.FreezeEarned

	out.FirstTry = p.Attempts == 1
	b := AwardFor(ex.XP, out.FirstTry, p.UsedHint)
	_, err = q.AwardXP(ctx, store.AwardXPParams{UserID: user, Amount: int32(b.Total()), Reason: ReasonPass, RefID: ex.ID})
	switch {
	case errors.Is(err, pgx.ErrNoRows): // already passed before: no XP
	case err != nil:
		return Outcome{}, fmt.Errorf("award xp: %w", err)
	default:
		out.XPAwarded, out.Breakdown = b.Total(), &b
		out.TotalXP += b.Total()
		if lvl := LevelFor(out.TotalXP); lvl > out.Level {
			out.Level, out.LevelUp = lvl, true
		}
		ref, err := joinLeague(ctx, q, user, st.LeagueTier, WeekStart(now))
		if err != nil {
			return Outcome{}, err
		}
		out.League = &ref
	}

	last := change.Last
	if err := q.UpdateUserStats(ctx, store.UpdateUserStatsParams{
		UserID: user, TotalXp: int32(out.TotalXP), Level: int32(out.Level),
		StreakDays: int32(change.Days), StreakFreezes: int32(change.Freezes), LastActiveDate: &last,
	}); err != nil {
		return Outcome{}, fmt.Errorf("update stats: %w", err)
	}
	return out, nil
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
