package economy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/store"
)

func streakOf(st store.GetEconomyStatsRow) game.Streak {
	s := game.Streak{Days: int(st.StreakDays), Freezes: int(st.StreakFreezes)}
	if st.LastActiveDate != nil {
		s.Last = *st.LastActiveDate
	}
	return s
}

// canRepair: a streak lost within the repair window (recorded, or lost by
// missed days the rollover has not processed yet).
func canRepair(st store.GetEconomyStatsRow, today time.Time) bool {
	if st.LostStreakDays > 0 && st.LostOn != nil && !st.LostOn.Before(today.AddDate(0, 0, -RepairWindowDays)) {
		return true
	}
	s := streakOf(st)
	if _, _, lost := game.Rollover(s, today); lost > 0 {
		return !s.Last.Before(today.AddDate(0, 0, -(RepairWindowDays + 1)))
	}
	return false
}

// rollover settles a user's day once per local date: missed days spend Streak
// Freezes (or lose the streak, kept for Streak Repair), the open bet is
// settled, and medals are evaluated. Safe to call any time (it is idempotent).
func (s *Service) rollover(ctx context.Context, q *store.Queries, user uuid.UUID) error {
	st, err := q.GetEconomyStats(ctx, user)
	if err != nil {
		return fmt.Errorf("stats: %w", err)
	}
	d := dayOf(s.now(), st.Timezone)
	if st.RolledOverOn != nil && !st.RolledOverOn.Before(d.Today) {
		return nil
	}
	next, _, lost := game.Rollover(streakOf(st), d.Today)
	lostDays, lostOn := st.LostStreakDays, st.LostOn
	if lost > 0 {
		lostDays, lostOn = int32(lost), &d.Today
	}
	var last *time.Time
	if !next.Last.IsZero() {
		last = &next.Last
	}
	if err := q.SetStreakState(ctx, store.SetStreakStateParams{
		UserID: user, StreakDays: int32(next.Days), StreakFreezes: int32(next.Freezes), LastActiveDate: last,
		LostStreakDays: lostDays, LostOn: lostOn,
	}); err != nil {
		return fmt.Errorf("streak: %w", err)
	}
	if err := s.settleBet(ctx, q, user, d.Today, nil); err != nil {
		return err
	}
	if err := s.evaluateMedals(ctx, q, user, nil); err != nil {
		return err
	}
	return q.SetRolledOver(ctx, store.SetRolledOverParams{UserID: user, Day: &d.Today})
}

// RunRollovers processes every user whose local day changed (the daily job).
func (s *Service) RunRollovers(ctx context.Context) (int, error) {
	users, err := store.New(s.Pool).DueRollovers(ctx, s.now())
	if err != nil {
		return 0, fmt.Errorf("due rollovers: %w", err)
	}
	for _, u := range users {
		if err := s.tx(ctx, func(tx pgx.Tx) error { return s.rollover(ctx, store.New(tx), u) }); err != nil {
			return 0, fmt.Errorf("rollover %s: %w", u, err)
		}
	}
	return len(users), nil
}

// RunDaily runs the rollover job now and then every interval until ctx ends.
func (s *Service) RunDaily(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		if _, err := s.RunRollovers(ctx); err != nil && ctx.Err() == nil && s.Log != nil {
			s.Log.Error("daily rollover", "err", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// repair restores a recently lost streak (Streak Repair).
func (s *Service) repair(ctx context.Context, q *store.Queries, user uuid.UUID, requestID string, today time.Time, res *UseResult) error {
	if err := s.rollover(ctx, q, user); err != nil {
		return err
	}
	st, err := q.GetEconomyStats(ctx, user)
	if err != nil {
		return err
	}
	if !canRepair(st, today) {
		return LimitError{"Your streak is safe: nothing to repair"}
	}
	used, err := q.CountPowerupUses(ctx, store.CountPowerupUsesParams{UserID: user, Powerup: StreakRepair, Since: s.now().AddDate(0, 0, -RepairEveryDays)})
	if err != nil {
		return err
	}
	if used > 0 {
		return LimitError{"Streak Repair works once a month"}
	}
	p, _ := PowerUpByID(StreakRepair)
	if _, err := credit(ctx, q, user, -p.Cost, "buy:"+StreakRepair, requestID); err != nil {
		return err
	}
	cur := streakOf(st)
	days, last := int(st.LostStreakDays)+cur.Days, cur.Last
	if cur.Days == 0 {
		last = today.AddDate(0, 0, -1) // alive again: pass today to extend it
	}
	if err := q.SetStreakState(ctx, store.SetStreakStateParams{
		UserID: user, StreakDays: int32(days), StreakFreezes: int32(cur.Freezes), LastActiveDate: &last,
	}); err != nil {
		return err
	}
	if err := q.SetLongestStreak(ctx, store.SetLongestStreakParams{UserID: user, Days: int32(days)}); err != nil {
		return err
	}
	res.Streak = days
	return nil
}

// placeBet stakes Bytes on hitting the daily goal 7 days running, from today.
func (s *Service) placeBet(ctx context.Context, q *store.Queries, user uuid.UUID, requestID string, stake, goal int, today time.Time, res *UseResult) error {
	if stake < BetMin || stake > BetMax {
		return fmt.Errorf("%w: stake must be %d-%d Bytes", ErrBadRequest, BetMin, BetMax)
	}
	if _, err := q.OpenBet(ctx, user); err == nil {
		return LimitError{"One bet at a time"}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if _, err := credit(ctx, q, user, -stake, "bet", requestID); err != nil {
		return err
	}
	if _, err := q.CreateBet(ctx, store.CreateBetParams{UserID: user, Stake: int32(stake), Goal: int32(goal), StartsOn: today}); err != nil {
		return err
	}
	res.Bet = &Bet{Stake: stake, Goal: goal, StartsOn: today.Format(time.DateOnly)}
	return nil
}

// openBet settles and returns the open bet with its progress.
func (s *Service) openBet(ctx context.Context, q *store.Queries, user uuid.UUID, today time.Time) (*Bet, error) {
	if err := s.settleBet(ctx, q, user, today, nil); err != nil {
		return nil, err
	}
	b, err := q.OpenBet(ctx, user)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	days, err := q.DayLessons(ctx, store.DayLessonsParams{UserID: user, FromDay: b.StartsOn, ToDay: b.StartsOn.AddDate(0, 0, 6)})
	if err != nil {
		return nil, err
	}
	met := 0
	for _, d := range days {
		if d.Lessons >= b.Goal {
			met++
		}
	}
	return &Bet{Stake: int(b.Stake), Goal: int(b.Goal), StartsOn: b.StartsOn.Format(time.DateOnly), DaysMet: met}, nil
}

// settleBet decides an open bet: lost once a finished day missed the goal,
// won (2x the stake back) when all 7 days from its start met it.
func (s *Service) settleBet(ctx context.Context, q *store.Queries, user uuid.UUID, today time.Time, out *game.RewardsOutcome) error {
	b, err := q.OpenBet(ctx, user)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open bet: %w", err)
	}
	end := b.StartsOn.AddDate(0, 0, 6)
	rows, err := q.DayLessons(ctx, store.DayLessonsParams{UserID: user, FromDay: b.StartsOn, ToDay: end})
	if err != nil {
		return err
	}
	met := map[string]bool{}
	for _, r := range rows {
		met[r.Day.Format(time.DateOnly)] = r.Lessons >= b.Goal
	}
	won := true
	for d := b.StartsOn; !d.After(end); d = d.AddDate(0, 0, 1) {
		if met[d.Format(time.DateOnly)] {
			continue
		}
		won = false
		if d.Before(today) { // a finished day missed the goal
			return q.SettleBet(ctx, store.SettleBetParams{ID: b.ID, Status: "lost"})
		}
	}
	if !won {
		return nil
	}
	if err := q.SettleBet(ctx, store.SettleBetParams{ID: b.ID, Status: "won"}); err != nil {
		return err
	}
	if _, err := credit(ctx, q, user, 2*int(b.Stake), "bet_win", b.ID.String()); err != nil {
		return err
	}
	if out != nil {
		out.Bytes = append(out.Bytes, game.Gain{Reason: "bet_win", Label: "Double or Nothing won", Amount: 2 * int(b.Stake)})
	}
	return bump(ctx, q, user, "bets_won", 1)
}

// QuestView is one quest with progress.
type QuestView struct {
	Quest
	Progress int  `json:"progress"`
	Done     bool `json:"done"`
}

// QuestsView is GET /quests.
type QuestsView struct {
	Daily    []QuestView `json:"daily"`
	Weekly   QuestView   `json:"weekly"`
	ResetsAt time.Time   `json:"resets_at"` // next local midnight
	Chest    struct {
		Ready  bool            `json:"ready"`
		Opened *game.ChestDrop `json:"opened"`
		Odds   any             `json:"odds"`
		Bytes  [2]int          `json:"bytes"`
	} `json:"chest"`
}

// Quests returns today's quests, this week's quest and the daily chest.
func (s *Service) Quests(ctx context.Context, user uuid.UUID) (QuestsView, error) {
	var v QuestsView
	err := s.tx(ctx, func(tx pgx.Tx) error {
		q := store.New(tx)
		if err := q.EnsureUserStats(ctx, user); err != nil {
			return err
		}
		st, err := q.GetEconomyStats(ctx, user)
		if err != nil {
			return err
		}
		d := dayOf(s.now(), st.Timezone)
		rows, err := q.ListUserQuests(ctx, store.ListUserQuestsParams{UserID: user, Day: d.Today, Week: d.Week})
		if err != nil {
			return err
		}
		prog := map[string]int{}
		for _, r := range rows {
			prog[r.QuestID+r.Period.Format(time.DateOnly)] = int(r.Progress)
		}
		daily, weekly := todaysQuests(user.String(), d.Today, d.Week)
		all := true
		for _, qu := range daily {
			p := prog[qu.ID+d.Today.Format(time.DateOnly)]
			qv := QuestView{Quest: qu, Progress: min(p, qu.Target), Done: p >= qu.Target}
			all = all && qv.Done
			v.Daily = append(v.Daily, qv)
		}
		wp := prog[weekly.ID+d.Week.Format(time.DateOnly)]
		v.Weekly = QuestView{Quest: weekly, Progress: min(wp, weekly.Target), Done: wp >= weekly.Target}
		midnight := time.Date(d.Today.Year(), d.Today.Month(), d.Today.Day(), 0, 0, 0, 0, d.Loc).AddDate(0, 0, 1)
		v.ResetsAt = midnight
		v.Chest.Ready, v.Chest.Odds, v.Chest.Bytes = all, ChestOdds, [2]int{ChestMinBytes, ChestMaxBytes}
		if c, err := q.GetChest(ctx, store.GetChestParams{UserID: user, Kind: "daily", RefID: d.Today.Format(time.DateOnly)}); err == nil {
			v.Chest.Opened = &game.ChestDrop{Bytes: int(c.Bytes), PowerUp: c.Powerup}
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		return nil
	})
	return v, err
}

// OpenDailyChest opens today's chest once all daily quests are done.
func (s *Service) OpenDailyChest(ctx context.Context, user uuid.UUID, requestID string) (json.RawMessage, error) {
	view, err := s.Quests(ctx, user)
	if err != nil {
		return nil, err
	}
	if !view.Chest.Ready {
		return nil, LimitError{"Finish today's three quests to open the chest"}
	}
	var out json.RawMessage
	err = s.tx(ctx, func(tx pgx.Tx) error {
		var err error
		out, err = once(ctx, tx, user, requestID, "chest:daily", func(q *store.Queries) (any, error) {
			st, err := q.GetEconomyStats(ctx, user)
			if err != nil {
				return nil, err
			}
			day := dayOf(s.now(), st.Timezone).Today.Format(time.DateOnly)
			bytes := ChestMinBytes + s.rand(ChestMaxBytes-ChestMinBytes+1)
			return s.openChest(ctx, q, user, "daily", day, bytes)
		})
		return err
	})
	return out, err
}

// LockedChapters returns the chapters of a course behind an undefeated boss.
func (s *Service) LockedChapters(ctx context.Context, user uuid.NullUUID, courseID string) (map[string]bool, error) {
	ids, err := store.New(s.Pool).LockedChapters(ctx, store.LockedChaptersParams{UserID: user, CourseID: courseID})
	if err != nil {
		return nil, fmt.Errorf("locked chapters: %w", err)
	}
	out := map[string]bool{}
	for _, id := range ids {
		out[id] = true
	}
	return out, nil
}
