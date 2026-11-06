package economy

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/store"
)

// Errors the API maps to responses. Nothing changes when one is returned.
var (
	ErrInsufficientBytes = errors.New("not enough Bytes")
	ErrLimit             = errors.New("limit reached")
	ErrUnavailable       = errors.New("not available")
	ErrNotOwned          = errors.New("you don't have this power-up")
	ErrNothingToPeek     = errors.New("nothing to peek at")
	ErrRequestReused     = errors.New("request id already used for another action")
	ErrBadRequest        = errors.New("bad request")
)

// LimitError says which limit was hit, for the user.
type LimitError struct{ Msg string }

func (e LimitError) Error() string { return e.Msg }
func (e LimitError) Unwrap() error { return ErrLimit }

// Service is the economy. It is also the game engine's Rewards and the league
// close's LeagueRewards.
type Service struct {
	Pool *pgxpool.Pool
	Now  func() time.Time
	// Rand returns a uniform int in [0, n) (crypto/rand when nil); tests fix it.
	Rand func(n int) int
	Log  *slog.Logger
}

func (s *Service) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s *Service) rand(n int) int {
	if s.Rand != nil {
		return s.Rand(n)
	}
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(fmt.Sprintf("crypto/rand: %v", err))
	}
	return int(v.Int64())
}

// Sync writes the medal and quest catalogues to their tables (at startup).
func (s *Service) Sync(ctx context.Context) error {
	return pgx.BeginFunc(ctx, s.Pool, func(tx pgx.Tx) error {
		q := store.New(tx)
		for i, m := range Medals {
			t, _ := json.Marshal(m.Tiers)
			if err := q.UpsertMedal(ctx, store.UpsertMedalParams{ID: m.ID, Family: m.Family, Name: m.Name, Secret: m.Secret, Tiers: t, Ord: int32(i)}); err != nil {
				return fmt.Errorf("sync medal %s: %w", m.ID, err)
			}
		}
		for _, qu := range append(append([]Quest{}, DailyPool...), WeeklyPool...) {
			rule, _ := json.Marshal(map[string]any{"event": qu.Event, "target": qu.Target, "name": qu.Name})
			reward, _ := json.Marshal(map[string]any{"bytes": qu.Bytes})
			if err := q.UpsertQuest(ctx, store.UpsertQuestParams{ID: qu.ID, Cadence: qu.Cadence, Rule: rule, Reward: reward}); err != nil {
				return fmt.Errorf("sync quest %s: %w", qu.ID, err)
			}
		}
		return nil
	})
}

// credit adds amount (negative to spend) to the ledger and the cached balance,
// once per (reason, ref). It reports false when that entry already exists, and
// ErrInsufficientBytes (leaving the ledger entry to the rolled-back
// transaction) when the balance would go negative.
func credit(ctx context.Context, q *store.Queries, user uuid.UUID, amount int, reason, ref string) (bool, error) {
	if err := q.EnsureUserStats(ctx, user); err != nil {
		return false, fmt.Errorf("ensure stats: %w", err)
	}
	_, err := q.InsertByteEvent(ctx, store.InsertByteEventParams{UserID: user, Amount: int32(amount), Reason: reason, RefID: ref})
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("byte event: %w", err)
	}
	if _, err := q.AddBytes(ctx, store.AddBytesParams{UserID: user, Amount: int32(amount)}); errors.Is(err, pgx.ErrNoRows) {
		return false, ErrInsufficientBytes
	} else if err != nil {
		return false, fmt.Errorf("balance: %w", err)
	}
	return true, nil
}

// once runs fn in tx at most once per (user, request id): a replay returns the
// first response without running fn again. Failed calls store nothing (their
// transaction rolls back), so they may be retried with the same id.
func once(ctx context.Context, tx pgx.Tx, user uuid.UUID, requestID, action string, fn func(q *store.Queries) (any, error)) (json.RawMessage, error) {
	if requestID == "" || len(requestID) > 64 {
		return nil, fmt.Errorf("%w: request_id must be 1-64 characters", ErrBadRequest)
	}
	q := store.New(tx)
	if err := q.LockIdempotencyKey(ctx, store.LockIdempotencyKeyParams{UserID: user.String(), RequestID: requestID}); err != nil {
		return nil, fmt.Errorf("lock request: %w", err)
	}
	prev, err := q.GetIdempotent(ctx, store.GetIdempotentParams{UserID: user, RequestID: requestID})
	if err == nil {
		if prev.Action != action {
			return nil, ErrRequestReused
		}
		return prev.Response, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("read request: %w", err)
	}
	res, err := fn(q)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("encode response: %w", err)
	}
	if err := q.PutIdempotent(ctx, store.PutIdempotentParams{UserID: user, RequestID: requestID, Action: action, Response: b}); err != nil {
		return nil, fmt.Errorf("store request: %w", err)
	}
	return b, nil
}

// Tx runs fn in a transaction.
func (s *Service) tx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, s.Pool, fn)
}

// userDay is the user's local date and the week start for quests.
type userDay struct {
	Loc   *time.Location
	Today time.Time // local date, midnight UTC
	Week  time.Time // Monday of the local week
}

func dayOf(now time.Time, tz string) userDay {
	loc := game.Location(tz)
	today := game.LocalDate(now, loc)
	return userDay{Loc: loc, Today: today, Week: today.AddDate(0, 0, -((int(today.Weekday()) + 6) % 7))}
}

func bump(ctx context.Context, q *store.Queries, user uuid.UUID, key string, n int) error {
	if n == 0 {
		return nil
	}
	if _, err := q.BumpCounter(ctx, store.BumpCounterParams{UserID: user, Key: key, N: int32(n)}); err != nil {
		return fmt.Errorf("counter %s: %w", key, err)
	}
	return nil
}

// drop picks a chest's random power-up by ChestOdds.
func (s *Service) drop() string {
	total := 0
	for _, o := range ChestOdds {
		total += o.Weight
	}
	r := s.rand(total)
	for _, o := range ChestOdds {
		if r < o.Weight {
			return o.PowerUp
		}
		r -= o.Weight
	}
	return ChestOdds[0].PowerUp
}

// openChest records a chest once and pays it out: Bytes plus one power-up.
func (s *Service) openChest(ctx context.Context, q *store.Queries, user uuid.UUID, kind, ref string, bytes int) (*game.ChestDrop, error) {
	p := s.drop()
	_, err := q.OpenChest(ctx, store.OpenChestParams{UserID: user, Kind: kind, RefID: ref, Bytes: int32(bytes), Powerup: p})
	if errors.Is(err, pgx.ErrNoRows) {
		c, err := q.GetChest(ctx, store.GetChestParams{UserID: user, Kind: kind, RefID: ref})
		if err != nil {
			return nil, fmt.Errorf("chest: %w", err)
		}
		return &game.ChestDrop{Bytes: int(c.Bytes), PowerUp: c.Powerup}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open chest: %w", err)
	}
	if _, err := credit(ctx, q, user, bytes, kind+"_chest", ref); err != nil {
		return nil, err
	}
	if err := q.AddInventory(ctx, store.AddInventoryParams{UserID: user, Powerup: p, N: 1}); err != nil {
		return nil, fmt.Errorf("chest power-up: %w", err)
	}
	return &game.ChestDrop{Bytes: bytes, PowerUp: p}, nil
}

// Dev starter grant: what the default dev user (DEV_AUTH only) starts with, so
// level-gated UI and every power-up can be tried without grinding.
const (
	DevStarterXP    = 1000
	DevStarterBytes = 5000 // one of every power-up plus a Streak Repair and a max bet
	devGrantReason  = "dev_grant"
	devGrantRef     = "starter"
)

// GrantDevStarter gives user DevStarterXP and DevStarterBytes once, through
// the XP and Byte ledgers (keyed on reason/ref, so repeat calls do nothing).
// Only the DEV_AUTH sign-in calls it. It never touches grading or progress.
func (s *Service) GrantDevStarter(ctx context.Context, user uuid.UUID) error {
	return s.tx(ctx, func(tx pgx.Tx) error {
		q := store.New(tx)
		if err := q.EnsureUserStats(ctx, user); err != nil {
			return fmt.Errorf("ensure stats: %w", err)
		}
		_, err := q.AwardXP(ctx, store.AwardXPParams{UserID: user, Amount: DevStarterXP, Reason: devGrantReason, RefID: devGrantRef})
		switch {
		case errors.Is(err, pgx.ErrNoRows): // granted before
		case err != nil:
			return fmt.Errorf("dev xp: %w", err)
		default:
			st, err := q.LockUserStats(ctx, user)
			if err != nil {
				return fmt.Errorf("lock stats: %w", err)
			}
			xp := int(st.TotalXp) + DevStarterXP
			if err := q.UpdateUserStats(ctx, store.UpdateUserStatsParams{
				UserID: user, TotalXp: int32(xp), Level: int32(max(int(st.Level), game.LevelFor(xp))),
				StreakDays: st.StreakDays, StreakFreezes: st.StreakFreezes, LastActiveDate: st.LastActiveDate,
			}); err != nil {
				return fmt.Errorf("dev stats: %w", err)
			}
		}
		_, err = credit(ctx, q, user, DevStarterBytes, devGrantReason, devGrantRef)
		return err
	})
}
