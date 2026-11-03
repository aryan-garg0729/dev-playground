package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/store"
)

type leagueSummary struct {
	Joined      bool      `json:"joined"` // false until the first XP of the week
	Tier        string    `json:"tier"`
	Rank        int       `json:"rank"`
	Of          int       `json:"of"`
	WeeklyXP    int       `json:"weekly_xp"`
	EndsAt      time.Time `json:"ends_at"`
	Promote     int       `json:"promote"` // how many move up this week
	Demote      int       `json:"demote"`  // how many move down
	XPToPromote int       `json:"xp_to_promote"`
	Zone        string    `json:"zone"` // promote | demote | stay (where the user stands now)
}

type leagueEntry struct {
	Rank     int    `json:"rank"`
	Name     string `json:"name"`  // display name only; other users' emails are never sent
	Title    string `json:"title"` // their worn medal title
	WeeklyXP int    `json:"weekly_xp"`
	You      bool   `json:"you"`
	Zone     string `json:"zone"`
}

type leagueResponse struct {
	leagueSummary
	Standings []leagueEntry `json:"standings"`
}

// league computes the user's standing this week. Rank comes from the Redis
// board; names and XP from Postgres (the source of truth).
func (a *api) league(ctx context.Context, user uuid.UUID, tier string) (leagueResponse, error) {
	week := game.WeekStart(a.now())
	out := leagueResponse{leagueSummary: leagueSummary{Tier: tier, EndsAt: game.EndsAt(week), Zone: "stay"}, Standings: []leagueEntry{}}
	q := store.New(a.Pool)
	m, err := q.GetLeagueMember(ctx, store.GetLeagueMemberParams{Week: week, UserID: user})
	if errors.Is(err, pgx.ErrNoRows) {
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("league member: %w", err)
	}
	ref := game.LeagueRef{Week: week, Tier: m.Tier, Group: int(m.GroupNo)}
	rows, err := q.LeagueStandings(ctx, store.LeagueStandingsParams{Week: week, Tier: m.Tier, GroupNo: m.GroupNo})
	if err != nil {
		return out, fmt.Errorf("standings: %w", err)
	}
	moves := game.Moves(m.Tier, len(rows))
	out.Joined, out.Tier, out.Of = true, m.Tier, len(rows)
	for i, r := range rows {
		zone := map[string]string{game.Promoted: "promote", game.Demoted: "demote", game.Stayed: "stay"}[moves[i]]
		e := leagueEntry{Rank: i + 1, Name: r.Name, Title: r.Title, WeeklyXP: int(r.WeeklyXp), You: r.UserID == user, Zone: zone}
		if e.Name == "" {
			e.Name = "Learner"
		}
		if moves[i] == game.Promoted {
			out.Promote++
		}
		if moves[i] == game.Demoted {
			out.Demote++
		}
		if e.You {
			out.Rank, out.WeeklyXP, out.Zone = e.Rank, e.WeeklyXP, zone
		}
		out.Standings = append(out.Standings, e)
	}
	if rank, of, err := a.Board.Rank(ctx, a.Pool, ref, user); err == nil && of == len(rows) {
		out.Rank = rank
	} else if err != nil {
		a.Log.Warn("league board unavailable; using database order", "err", err)
	}
	if out.Promote > 0 && out.Zone != "promote" {
		out.XPToPromote = int(rows[out.Promote-1].WeeklyXp) - out.WeeklyXP + 1
	}
	return out, nil
}

func (a *api) leagueSummary(ctx context.Context, user uuid.UUID, tier string) (leagueSummary, error) {
	l, err := a.league(ctx, user, tier)
	return l.leagueSummary, err
}

func (a *api) currentLeague(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !u.Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in first")
		return
	}
	stats, _, err := a.userStats(r.Context(), u.UUID)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	out, err := a.league(r.Context(), u.UUID, stats.LeagueTier)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
