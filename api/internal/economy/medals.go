package economy

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/store"
)

// medalStats is every value a medal rule can read: counters plus user stats.
func medalStats(ctx context.Context, q *store.Queries, user uuid.UUID) (map[string]int, error) {
	rows, err := q.ListCounters(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("counters: %w", err)
	}
	v := map[string]int{}
	for _, r := range rows {
		v[r.Key] = int(r.Value)
	}
	st, err := q.GetEconomyStats(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("stats: %w", err)
	}
	v["longest_streak"] = int(st.LongestStreak)
	v["total_xp"] = int(st.TotalXp)
	return v, nil
}

// evaluateMedals awards every tier the user has reached and not yet received,
// once each (user_medals' key), paying its Bytes. Run after a submission, the
// league close and the daily rollover.
func (s *Service) evaluateMedals(ctx context.Context, q *store.Queries, user uuid.UUID, out *game.RewardsOutcome) error {
	stats, err := medalStats(ctx, q, user)
	if err != nil {
		return err
	}
	for _, m := range Medals {
		if m.ComingSoon {
			continue
		}
		for _, t := range m.Tiers {
			if stats[m.Stat] < t.Gte {
				break
			}
			_, err := q.AwardMedal(ctx, store.AwardMedalParams{UserID: user, MedalID: m.ID, Tier: t.Tier})
			if errors.Is(err, pgx.ErrNoRows) {
				continue // earned before
			}
			if err != nil {
				return fmt.Errorf("award medal: %w", err)
			}
			if _, err := credit(ctx, q, user, t.Bytes, "medal", m.ID+":"+t.Tier); err != nil {
				return err
			}
			if out != nil {
				out.Medals = append(out.Medals, game.MedalAward{ID: m.ID, Name: m.Name, Tier: t.Tier, Title: t.Title})
				out.Bytes = append(out.Bytes, game.Gain{Reason: "medal", Label: m.Name + " medal", Amount: t.Bytes})
			}
		}
	}
	return nil
}

// OnLeagueClosed pays the league's top 5, records podiums and the highest tier
// reached, counts shield saves, and evaluates medals, in the close transaction.
func (s *Service) OnLeagueClosed(ctx context.Context, tx pgx.Tx, week time.Time, results []game.LeagueResult) error {
	q := store.New(tx)
	for _, r := range results {
		if r.Rank <= 5 {
			if _, err := credit(ctx, q, r.User, LeagueTop5Bytes, "league_top5", week.Format(time.DateOnly)); err != nil {
				return err
			}
			if err := bump(ctx, q, r.User, "top5", 1); err != nil {
				return err
			}
		}
		if err := bump(ctx, q, r.User, "podium", b2i(r.Rank <= 3)); err != nil {
			return err
		}
		if err := bump(ctx, q, r.User, "shield_saves", b2i(r.Shield)); err != nil {
			return err
		}
		if err := q.MaxCounter(ctx, store.MaxCounterParams{UserID: r.User, Key: "tier_max", N: int32(slices.Index(game.Tiers, r.NewTier) + 1)}); err != nil {
			return fmt.Errorf("tier max: %w", err)
		}
		if err := s.evaluateMedals(ctx, q, r.User, nil); err != nil {
			return err
		}
	}
	return nil
}

// MedalTier is a tier as the profile shows it.
type MedalTier struct {
	Tier     string     `json:"tier"`
	Gte      int        `json:"gte"`
	Title    string     `json:"title"`
	Bytes    int        `json:"bytes"`
	EarnedAt *time.Time `json:"earned_at"`
}

// MedalView is one medal with the user's progress. Unearned secret medals
// are masked ("???") and carry no rule.
type MedalView struct {
	ID          string      `json:"id"`
	Family      string      `json:"family"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Unit        string      `json:"unit"`
	Secret      bool        `json:"secret"`
	ComingSoon  bool        `json:"coming_soon"`
	Value       int         `json:"value"`
	Tier        string      `json:"tier"`      // highest earned ("" if none)
	NextTier    string      `json:"next_tier"` // "" when all earned
	NextGte     int         `json:"next_gte"`
	Tiers       []MedalTier `json:"tiers"`
}

// MedalsView is the profile's medals tab.
type MedalsView struct {
	Earned   int         `json:"earned"` // medals with at least one tier
	Total    int         `json:"total"`
	Medals   []MedalView `json:"medals"`
	Families []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"families"`
	Titles []string `json:"titles"` // titles the user may wear
	Title  string   `json:"title"`
}

// Medals returns the user's medals and progress.
func (s *Service) Medals(ctx context.Context, user uuid.UUID) (MedalsView, error) {
	var out MedalsView
	err := s.tx(ctx, func(tx pgx.Tx) error {
		q := store.New(tx)
		if err := q.EnsureUserStats(ctx, user); err != nil {
			return err
		}
		stats, err := medalStats(ctx, q, user)
		if err != nil {
			return err
		}
		earned, err := q.ListUserMedals(ctx, user)
		if err != nil {
			return fmt.Errorf("user medals: %w", err)
		}
		at := map[string]time.Time{}
		for _, e := range earned {
			at[e.MedalID+":"+e.Tier] = e.EarnedAt
		}
		st, err := q.GetEconomyStats(ctx, user)
		if err != nil {
			return err
		}
		out.Title, out.Titles = st.Title, []string{}
		for _, f := range Families {
			out.Families = append(out.Families, struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{f.ID, f.Name})
		}
		for _, m := range Medals {
			v := MedalView{ID: m.ID, Family: m.Family, Name: m.Name, Description: m.Description, Unit: m.Unit,
				Secret: m.Secret, ComingSoon: m.ComingSoon, Value: stats[m.Stat], Tiers: []MedalTier{}}
			for _, t := range m.Tiers {
				mt := MedalTier{Tier: t.Tier, Gte: t.Gte, Title: t.Title, Bytes: t.Bytes}
				if when, ok := at[m.ID+":"+t.Tier]; ok {
					mt.EarnedAt = &when
					v.Tier = t.Tier
					out.Titles = append(out.Titles, t.Title)
				} else if v.NextTier == "" {
					v.NextTier, v.NextGte = t.Tier, t.Gte
				}
				v.Tiers = append(v.Tiers, mt)
			}
			if !m.ComingSoon {
				out.Total++
				if v.Tier != "" {
					out.Earned++
				}
			}
			if m.Secret && v.Tier == "" { // hidden until earned
				v = MedalView{ID: m.ID, Family: m.Family, Name: "???", Description: "Hidden until earned", Secret: true, Tiers: []MedalTier{}}
			}
			out.Medals = append(out.Medals, v)
		}
		return nil
	})
	return out, err
}

// SetTitle wears an earned title ("" removes it).
func (s *Service) SetTitle(ctx context.Context, user uuid.UUID, title string) error {
	if title != "" {
		v, err := s.Medals(ctx, user)
		if err != nil {
			return err
		}
		if !slices.Contains(v.Titles, title) {
			return fmt.Errorf("%w: earn that title first", ErrUnavailable)
		}
	}
	if err := store.New(s.Pool).SetTitle(ctx, store.SetTitleParams{UserID: user, Title: title}); err != nil {
		return fmt.Errorf("set title: %w", err)
	}
	return nil
}
