package economy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/store"
)

// ShopItem is a power-up with the user's state.
type ShopItem struct {
	PowerUp
	Owned int `json:"owned"`
	// State: buy | max_held | streak_safe | repair | pro | coming_soon | bet_open | place_bet
	State string `json:"state"`
	// Use: how a held one is used: auto | now | lessons | "" (not held)
	Use string `json:"use"`
	// UsableNow: a "now" power-up can be used at this moment.
	UsableNow bool `json:"usable_now"`
}

// Bet is an open Double or Nothing bet.
type Bet struct {
	Stake    int    `json:"stake"`
	Goal     int    `json:"goal"`
	StartsOn string `json:"starts_on"`
	DaysMet  int    `json:"days_met"`
}

// Effect is a power-up in force.
type Effect struct {
	Effect    string    `json:"effect"`
	RefID     string    `json:"ref_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ShopView is GET /shop (and /inventory).
type ShopView struct {
	Balance     int        `json:"balance"`
	EarnedToday int        `json:"earned_today"`
	Items       []ShopItem `json:"items"`
	Effects     []Effect   `json:"effects"`
	Bet         *Bet       `json:"bet"`
	Earn        []struct {
		Label  string `json:"label"`
		Amount string `json:"amount"`
	} `json:"earn"`
}

var usage = map[string]string{
	StreakFreeze: "auto", FocusBoost: "now", LeagueShield: "now", SecondWind: "lessons", Peek: "lessons",
}

// Shop returns the shop with the user's balance, inventory and effects.
func (s *Service) Shop(ctx context.Context, user uuid.UUID) (ShopView, error) {
	var v ShopView
	err := s.tx(ctx, func(tx pgx.Tx) error {
		q := store.New(tx)
		if err := q.EnsureUserStats(ctx, user); err != nil {
			return err
		}
		if err := s.rollover(ctx, q, user); err != nil { // settles bets and spent freezes first
			return err
		}
		st, err := q.GetEconomyStats(ctx, user)
		if err != nil {
			return err
		}
		now := s.now()
		d := dayOf(now, st.Timezone)
		v.Balance = int(st.Bytes)
		start := time.Date(d.Today.Year(), d.Today.Month(), d.Today.Day(), 0, 0, 0, 0, d.Loc)
		earned, err := q.BytesEarnedSince(ctx, store.BytesEarnedSinceParams{UserID: user, Since: start})
		if err != nil {
			return err
		}
		v.EarnedToday = int(earned)
		inv, err := q.ListInventory(ctx, user)
		if err != nil {
			return err
		}
		held := map[string]int{StreakFreeze: int(st.StreakFreezes)}
		for _, r := range inv {
			held[r.Powerup] = int(r.Count)
		}
		bet, err := s.openBet(ctx, q, user, d.Today)
		if err != nil {
			return err
		}
		v.Bet = bet
		repairable := canRepair(st, d.Today)
		for _, p := range PowerUps {
			it := ShopItem{PowerUp: p, Owned: held[p.ID], State: "buy", Use: usage[p.ID]}
			switch {
			case p.Unavailable != "":
				it.State = p.Unavailable
			case p.ID == StreakFreeze && it.Owned >= MaxFreezes:
				it.State = "max_held"
			case p.ID == StreakRepair && !repairable:
				it.State = "streak_safe"
			case p.ID == StreakRepair:
				it.State = "repair"
			case p.ID == DoubleOrNothing && bet != nil:
				it.State = "bet_open"
			case p.ID == DoubleOrNothing:
				it.State = "place_bet"
			}
			if it.Use == "now" && it.Owned > 0 {
				it.UsableNow = s.usableNow(ctx, q, user, p.ID, now, d) == nil
			}
			v.Items = append(v.Items, it)
		}
		effects, err := q.ListEffects(ctx, store.ListEffectsParams{UserID: user, Now: now})
		if err != nil {
			return err
		}
		v.Effects = []Effect{}
		for _, e := range effects {
			v.Effects = append(v.Effects, Effect{Effect: e.Effect, RefID: e.RefID, ExpiresAt: e.ExpiresAt})
		}
		for _, e := range [][2]string{
			{"Daily chest (3 quests)", fmt.Sprintf("%d–%d", ChestMinBytes, ChestMaxBytes)},
			{"Chapter complete", strconv.Itoa(ChapterBytes)}, {"Boss defeated", strconv.Itoa(BossBytes)},
			{"Medal earned", "50–500"}, {"League top 5", strconv.Itoa(LeagueTop5Bytes)},
		} {
			v.Earn = append(v.Earn, struct {
				Label  string `json:"label"`
				Amount string `json:"amount"`
			}{e[0], e[1]})
		}
		return nil
	})
	return v, err
}

// BuyResult is the response to a buy.
type BuyResult struct {
	PowerUp string `json:"powerup"`
	Owned   int    `json:"owned"`
	Balance int    `json:"balance"`
}

// Buy spends Bytes on a power-up, once per request id. With too few Bytes
// (or a limit hit) nothing changes.
func (s *Service) Buy(ctx context.Context, user uuid.UUID, powerup, requestID string) (json.RawMessage, error) {
	p, ok := PowerUpByID(powerup)
	if !ok || p.Cost == 0 || p.Unavailable != "" || p.ID == StreakRepair {
		return nil, fmt.Errorf("%w: %s cannot be bought", ErrUnavailable, powerup)
	}
	var out json.RawMessage
	err := s.tx(ctx, func(tx pgx.Tx) error {
		var err error
		out, err = once(ctx, tx, user, requestID, "buy:"+powerup, func(q *store.Queries) (any, error) {
			if _, err := credit(ctx, q, user, -p.Cost, "buy:"+powerup, requestID); err != nil {
				return nil, err
			}
			res := BuyResult{PowerUp: powerup}
			if powerup == StreakFreeze {
				n, err := q.AddFreezes(ctx, store.AddFreezesParams{UserID: user, N: 1, Cap: MaxFreezes})
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, LimitError{"You already hold 2 Streak Freezes"}
				} else if err != nil {
					return nil, err
				}
				res.Owned = int(n)
			} else {
				if err := q.AddInventory(ctx, store.AddInventoryParams{UserID: user, Powerup: powerup, N: 1}); err != nil {
					return nil, err
				}
				n, err := q.InventoryCount(ctx, store.InventoryCountParams{UserID: user, Powerup: powerup})
				if err != nil {
					return nil, err
				}
				res.Owned = int(n)
			}
			st, err := q.GetEconomyStats(ctx, user)
			if err != nil {
				return nil, err
			}
			res.Balance = int(st.Bytes)
			return res, nil
		})
		return err
	})
	return out, err
}

// UseRequest is POST /powerups/{powerup}/use.
type UseRequest struct {
	RequestID  string `json:"request_id"`
	ExerciseID string `json:"exercise_id"` // second_wind, peek
	Stake      int    `json:"stake"`       // double_or_nothing
}

// UseResult is the response to a use.
type UseResult struct {
	PowerUp   string     `json:"powerup"`
	Remaining int        `json:"remaining"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Peek      *PeekShown `json:"peek,omitempty"`
	Bet       *Bet       `json:"bet,omitempty"`
	Streak    int        `json:"streak,omitempty"`
	Balance   int        `json:"balance"`
}

// PeekShown is what Peek reveals: a failed hidden check's input, nothing else.
type PeekShown struct {
	Check int    `json:"check"` // row index in the latest result
	Name  string `json:"name"`  // the row's (masked) name, e.g. "Hidden test 2"
	Input string `json:"input"`
}

// usableNow says whether a "now" power-up can be used at the moment.
func (s *Service) usableNow(ctx context.Context, q *store.Queries, user uuid.UUID, powerup string, now time.Time, d userDay) error {
	switch powerup {
	case FocusBoost:
		start := time.Date(d.Today.Year(), d.Today.Month(), d.Today.Day(), 0, 0, 0, 0, d.Loc)
		n, err := q.CountPowerupUses(ctx, store.CountPowerupUsesParams{UserID: user, Powerup: FocusBoost, Since: start})
		if err != nil {
			return err
		}
		if n > 0 {
			return LimitError{"One Focus Boost a day"}
		}
	case LeagueShield:
		week := game.WeekStart(now)
		for _, w := range []time.Time{week, week.AddDate(0, 0, -7)} {
			on, err := q.IsLeagueShielded(ctx, store.IsLeagueShieldedParams{UserID: user, Week: w.Format(time.DateOnly)})
			if err != nil {
				return err
			}
			if on {
				return LimitError{"League Shield can't be used two weeks in a row"}
			}
		}
	}
	return nil
}

// Use spends a held power-up (or Bytes, for Streak Repair and bets), once per request id.
func (s *Service) Use(ctx context.Context, user uuid.UUID, powerup string, req UseRequest) (json.RawMessage, error) {
	if _, ok := PowerUpByID(powerup); !ok {
		return nil, fmt.Errorf("%w: unknown power-up %s", ErrUnavailable, powerup)
	}
	var out json.RawMessage
	err := s.tx(ctx, func(tx pgx.Tx) error {
		var err error
		out, err = once(ctx, tx, user, req.RequestID, "use:"+powerup, func(q *store.Queries) (any, error) {
			if err := q.EnsureUserStats(ctx, user); err != nil {
				return nil, err
			}
			st, err := q.GetEconomyStats(ctx, user) // locks the user's stats row
			if err != nil {
				return nil, err
			}
			now := s.now()
			d := dayOf(now, st.Timezone)
			res := UseResult{PowerUp: powerup}
			switch powerup {
			case FocusBoost, LeagueShield:
				err = s.useNow(ctx, q, user, powerup, now, d, &res)
			case SecondWind:
				err = s.useSecondWind(ctx, q, user, req.ExerciseID, &res)
			case Peek:
				err = s.usePeek(ctx, q, user, req.ExerciseID, &res)
			case StreakRepair:
				err = s.repair(ctx, q, user, req.RequestID, d.Today, &res)
			case DoubleOrNothing:
				err = s.placeBet(ctx, q, user, req.RequestID, req.Stake, int(st.DailyGoal), d.Today, &res)
			default:
				err = fmt.Errorf("%w: %s is not used by hand", ErrUnavailable, powerup)
			}
			if err != nil {
				return nil, err
			}
			if err := q.RecordPowerupUse(ctx, store.RecordPowerupUseParams{UserID: user, Powerup: powerup, RefID: useRef(powerup, req, now)}); err != nil {
				return nil, err
			}
			after, err := q.GetEconomyStats(ctx, user)
			if err != nil {
				return nil, err
			}
			res.Balance = int(after.Bytes)
			return res, nil
		})
		return err
	})
	return out, err
}

func useRef(powerup string, req UseRequest, now time.Time) string {
	switch powerup {
	case LeagueShield:
		return game.WeekStart(now).Format(time.DateOnly)
	case SecondWind, Peek:
		return req.ExerciseID
	}
	return ""
}

func take(ctx context.Context, q *store.Queries, user uuid.UUID, powerup string) (int, error) {
	n, err := q.TakeInventory(ctx, store.TakeInventoryParams{UserID: user, Powerup: powerup})
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotOwned
	}
	if err != nil {
		return 0, fmt.Errorf("inventory: %w", err)
	}
	return int(n), nil
}

func (s *Service) useNow(ctx context.Context, q *store.Queries, user uuid.UUID, powerup string, now time.Time, d userDay, res *UseResult) error {
	if err := s.usableNow(ctx, q, user, powerup, now, d); err != nil {
		return err
	}
	n, err := take(ctx, q, user, powerup)
	if err != nil {
		return err
	}
	res.Remaining = n
	exp, ref := now.Add(FocusBoostMinutes*time.Minute), ""
	if powerup == LeagueShield {
		week := game.WeekStart(now)
		exp, ref = game.EndsAt(week).Add(24*time.Hour), week.Format(time.DateOnly)
	}
	res.ExpiresAt = &exp
	return q.SetEffect(ctx, store.SetEffectParams{UserID: user, Effect: powerup, RefID: ref, ExpiresAt: exp})
}

// useSecondWind arms Second Wind on an exercise right after exactly one failed submit.
func (s *Service) useSecondWind(ctx context.Context, q *store.Queries, user uuid.UUID, exerciseID string, res *UseResult) error {
	p, err := q.GetProgress(ctx, store.GetProgressParams{UserID: user, ExerciseID: exerciseID})
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && (p.Status != "started" || p.Attempts != 1)) {
		return LimitError{"Second Wind works right after your first failed submit"}
	}
	if err != nil {
		return err
	}
	meta, err := q.ExerciseChapter(ctx, exerciseID)
	if err != nil {
		return fmt.Errorf("%w: unknown exercise", ErrBadRequest)
	}
	used, err := q.CountPowerupUsesInChapter(ctx, store.CountPowerupUsesInChapterParams{UserID: user, Powerup: SecondWind, ChapterID: meta.ChapterID})
	if err != nil {
		return err
	}
	if used > 0 {
		return LimitError{"One Second Wind per chapter"}
	}
	if res.Remaining, err = take(ctx, q, user, SecondWind); err != nil {
		return err
	}
	return q.SetEffect(ctx, store.SetEffectParams{UserID: user, Effect: SecondWind, RefID: exerciseID, ExpiresAt: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)})
}

// usePeek reveals the declared input of one failed hidden check of the user's
// latest submission (the first not shown yet) and removes the first-try bonus.
// Never the expected output, never test source; nothing to show costs nothing.
func (s *Service) usePeek(ctx context.Context, q *store.Queries, user uuid.UUID, exerciseID string, res *UseResult) error {
	sub, err := q.LatestPeekable(ctx, store.LatestPeekableParams{UserID: user, ExerciseID: exerciseID})
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && len(sub.HiddenInputs) == 0) {
		return ErrNothingToPeek
	}
	if err != nil {
		return err
	}
	var inputs map[string]string
	var result struct {
		Checks []struct{ Name string } `json:"checks"`
	}
	if json.Unmarshal(sub.HiddenInputs, &inputs) != nil || json.Unmarshal(sub.Result, &result) != nil {
		return ErrNothingToPeek
	}
	seen, err := q.PeekedChecks(ctx, store.PeekedChecksParams{UserID: user, SubmissionID: sub.ID.String()})
	if err != nil {
		return err
	}
	pick := -1
	for i := range result.Checks {
		if _, ok := inputs[strconv.Itoa(i)]; ok && !contains(seen, strconv.Itoa(i)) {
			pick = i
			break
		}
	}
	if pick < 0 {
		return ErrNothingToPeek
	}
	if res.Remaining, err = take(ctx, q, user, Peek); err != nil {
		return err
	}
	if err := q.SetPeeked(ctx, store.SetPeekedParams{UserID: user, ExerciseID: exerciseID}); err != nil {
		return err
	}
	if err := q.RecordPeek(ctx, store.RecordPeekParams{UserID: user, SubmissionID: sub.ID.String(), CheckIndex: strconv.Itoa(pick)}); err != nil {
		return err
	}
	res.Peek = &PeekShown{Check: pick, Name: result.Checks[pick].Name, Input: inputs[strconv.Itoa(pick)]}
	return nil
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
