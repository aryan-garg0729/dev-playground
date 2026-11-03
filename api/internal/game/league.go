package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/aryangarg/codequest/api/internal/store"
)

// Board keeps a Redis sorted set per league group, scored by weekly XP. It is
// rebuilt from Postgres (the source of truth) on every sync, so syncing is
// idempotent and a lost or flushed key heals itself.
type Board struct {
	Redis  *redis.Client
	Prefix string // key prefix (tests isolate themselves with one)
}

const boardTTL = 21 * 24 * time.Hour

func (b *Board) key(ref LeagueRef) string {
	return fmt.Sprintf("%scq:league:%s:%s:%d", b.Prefix, ref.Week.Format(time.DateOnly), ref.Tier, ref.Group)
}

// score orders by weekly XP, then by who reached it first, so Redis ranks
// match Postgres (LeagueStandings) exactly.
func score(xp int, last time.Time) float64 {
	const span = 1e10 // > any unix time in seconds we will see
	t := float64(last.Unix())
	if t >= span || t < 0 {
		t = span - 1
	}
	return float64(xp)*span + (span - 1 - t)
}

// Sync rewrites the group's set from Postgres. Call it after the grading
// transaction commits.
func (b *Board) Sync(ctx context.Context, db store.DBTX, ref LeagueRef) error {
	rows, err := store.New(db).LeagueStandings(ctx, store.LeagueStandingsParams{Week: ref.Week, Tier: ref.Tier, GroupNo: int32(ref.Group)})
	if err != nil {
		return fmt.Errorf("league standings: %w", err)
	}
	members := make([]redis.Z, len(rows))
	for i, r := range rows {
		members[i] = redis.Z{Score: score(int(r.WeeklyXp), r.LastXpAt), Member: r.UserID.String()}
	}
	key := b.key(ref)
	_, err = b.Redis.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Del(ctx, key)
		if len(members) > 0 {
			p.ZAdd(ctx, key, members...)
			p.Expire(ctx, key, boardTTL)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("sync league board: %w", err)
	}
	return nil
}

// Rank is the user's 1-based rank in the group and the group size.
func (b *Board) Rank(ctx context.Context, db store.DBTX, ref LeagueRef, user uuid.UUID) (rank, of int, err error) {
	key := b.key(ref)
	for attempt := 0; attempt < 2; attempt++ {
		r, err := b.Redis.ZRevRank(ctx, key, user.String()).Result()
		if err == nil {
			n, err := b.Redis.ZCard(ctx, key).Result()
			if err != nil {
				return 0, 0, fmt.Errorf("league size: %w", err)
			}
			return int(r) + 1, int(n), nil
		}
		if !errors.Is(err, redis.Nil) {
			return 0, 0, fmt.Errorf("league rank: %w", err)
		}
		if err := b.Sync(ctx, db, ref); err != nil {
			return 0, 0, err
		}
	}
	return 0, 0, errors.New("user is not in this league group")
}

// Leagues closes finished league weeks.
type Leagues struct {
	Pool *pgxpool.Pool
	Now  func() time.Time
	Log  *slog.Logger
	// Rewards, when set, runs in the close transaction (league Bytes, medals).
	Rewards LeagueRewards
}

// LeagueResult is one member's week.
type LeagueResult struct {
	User    uuid.UUID
	Tier    string // the week's tier
	NewTier string
	Rank    int
	Of      int
	Result  string // promoted | demoted | stayed
	Shield  bool   // a League Shield turned a demotion into a stay
}

// LeagueRewards reacts to a closed week inside the close transaction.
type LeagueRewards interface {
	OnLeagueClosed(ctx context.Context, tx pgx.Tx, week time.Time, results []LeagueResult) error
}

// CloseDue closes every week before the current one that has members and is
// not closed yet: ranks each group, records results, and moves tiers. It is
// safe to run concurrently and repeatedly (an advisory lock plus the
// league_weeks row make each week close exactly once).
func (l *Leagues) CloseDue(ctx context.Context) (int, error) {
	now := time.Now
	if l.Now != nil {
		now = l.Now
	}
	weeks, err := store.New(l.Pool).UnclosedLeagueWeeks(ctx, WeekStart(now()))
	if err != nil {
		return 0, fmt.Errorf("unclosed weeks: %w", err)
	}
	closed := 0
	for _, w := range weeks {
		did, err := l.closeWeek(ctx, w)
		if err != nil {
			return closed, err
		}
		if did {
			closed++
		}
	}
	return closed, nil
}

func (l *Leagues) closeWeek(ctx context.Context, week time.Time) (bool, error) {
	did := false
	err := pgx.BeginFunc(ctx, l.Pool, func(tx pgx.Tx) error {
		q := store.New(tx)
		if err := q.LockLeagueKey(ctx, "league-close:"+week.Format(time.DateOnly)); err != nil {
			return fmt.Errorf("lock: %w", err)
		}
		if done, err := q.IsLeagueWeekClosed(ctx, week); err != nil || done {
			return err
		}
		groups, err := q.ListLeagueGroups(ctx, week)
		if err != nil {
			return fmt.Errorf("groups: %w", err)
		}
		var results []LeagueResult
		for _, g := range groups {
			rows, err := q.LeagueStandings(ctx, store.LeagueStandingsParams{Week: week, Tier: g.Tier, GroupNo: g.GroupNo})
			if err != nil {
				return fmt.Errorf("standings: %w", err)
			}
			moves := Moves(g.Tier, len(rows))
			for i, r := range rows {
				res := LeagueResult{User: r.UserID, Tier: g.Tier, Rank: i + 1, Of: len(rows), Result: moves[i]}
				if res.Result == Demoted {
					shielded, err := q.IsLeagueShielded(ctx, store.IsLeagueShieldedParams{UserID: r.UserID, Week: week.Format(time.DateOnly)})
					if err != nil {
						return fmt.Errorf("shield: %w", err)
					}
					if shielded {
						res.Result, res.Shield = Stayed, true
					}
				}
				res.NewTier = NextTier(g.Tier, res.Result)
				if err := q.SetLeagueResult(ctx, store.SetLeagueResultParams{
					Week: week, UserID: r.UserID, FinalRank: int32(res.Rank), Result: res.Result,
				}); err != nil {
					return fmt.Errorf("result: %w", err)
				}
				if err := q.SetLeagueTier(ctx, store.SetLeagueTierParams{UserID: r.UserID, LeagueTier: res.NewTier}); err != nil {
					return fmt.Errorf("tier: %w", err)
				}
				results = append(results, res)
			}
		}
		if l.Rewards != nil {
			if err := l.Rewards.OnLeagueClosed(ctx, tx, week, results); err != nil {
				return fmt.Errorf("league rewards: %w", err)
			}
		}
		did = true
		return q.CloseLeagueWeek(ctx, week)
	})
	if err != nil {
		return false, fmt.Errorf("close league week %s: %w", week.Format(time.DateOnly), err)
	}
	if did && l.Log != nil {
		l.Log.Info("league week closed", "week", week.Format(time.DateOnly))
	}
	return did, nil
}

// Run closes due weeks now and then every interval until ctx ends.
func (l *Leagues) Run(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		if _, err := l.CloseDue(ctx); err != nil && l.Log != nil && ctx.Err() == nil {
			l.Log.Error("close league weeks", "err", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// EndsAt is when the league week starting at week closes.
func EndsAt(week time.Time) time.Time { return week.AddDate(0, 0, 7) }

// TierTitle is a display name ("Gold").
func TierTitle(tier string) string {
	if tier == "" {
		return ""
	}
	return string(tier[0]-'a'+'A') + tier[1:]
}
