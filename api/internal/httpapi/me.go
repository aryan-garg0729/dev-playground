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

type userDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Timezone  string `json:"timezone"`
	Goal      string `json:"goal"`      // onboarding goal: "", zero, backend, interviews, ai
	Onboarded bool   `json:"onboarded"` // finished (or skipped) onboarding
}

func toUserDTO(u store.GetUserRow) userDTO {
	return userDTO{ID: u.ID.String(), Email: u.Email, Name: u.Name, Avatar: u.Avatar, Timezone: u.Timezone,
		Goal: u.Goal, Onboarded: u.OnboardedAt != nil}
}

// featuresDTO says which optional features this server has turned on.
type featuresDTO struct {
	Tutor bool `json:"tutor"`
}

// onboardingGoals are the goals a learner can pick (users.goal check).
var onboardingGoals = map[string]bool{"zero": true, "backend": true, "interviews": true, "ai": true}

type statsDTO struct {
	TotalXP       int    `json:"total_xp"`
	Level         int    `json:"level"`
	LevelStartXP  int    `json:"level_start_xp"`
	NextLevelXP   int    `json:"next_level_xp"`
	StreakDays    int    `json:"streak_days"`
	StreakFreezes int    `json:"streak_freezes"`
	LeagueTier    string `json:"league_tier"`
	DailyGoal     int    `json:"daily_goal"`
	TodayLessons  int    `json:"today_lessons"`
	Bytes         int    `json:"bytes"`
	Title         string `json:"title"`
	LongestStreak int    `json:"longest_streak"`
}

type meResponse struct {
	User     userDTO     `json:"user"`
	Stats    statsDTO    `json:"stats"`
	Features featuresDTO `json:"features"`
}

// userStats loads the user's stats and their local "today".
func (a *api) userStats(ctx context.Context, user uuid.UUID) (statsDTO, time.Time, error) {
	q := store.New(a.Pool)
	if err := q.EnsureUserStats(ctx, user); err != nil {
		return statsDTO{}, time.Time{}, fmt.Errorf("ensure stats: %w", err)
	}
	s, err := q.GetUserStats(ctx, user)
	if err != nil {
		return statsDTO{}, time.Time{}, fmt.Errorf("stats: %w", err)
	}
	today := game.LocalDate(a.now(), game.Location(s.Timezone))
	streak := game.Streak{Days: int(s.StreakDays), Freezes: int(s.StreakFreezes)}
	if s.LastActiveDate != nil {
		streak.Last = *s.LastActiveDate
	}
	days, err := q.ListDayPasses(ctx, store.ListDayPassesParams{UserID: user, FromDay: today, ToDay: today})
	if err != nil {
		return statsDTO{}, time.Time{}, fmt.Errorf("today: %w", err)
	}
	out := statsDTO{
		TotalXP: int(s.TotalXp), Level: int(s.Level),
		LevelStartXP: game.LevelStart(int(s.Level)), NextLevelXP: game.LevelStart(int(s.Level) + 1),
		StreakDays: game.Current(streak, today), StreakFreezes: int(s.StreakFreezes),
		LeagueTier: s.LeagueTier, DailyGoal: int(s.DailyGoal),
		Bytes: int(s.Bytes), Title: s.Title, LongestStreak: int(s.LongestStreak),
	}
	if len(days) > 0 {
		out.TodayLessons = int(days[0].Lessons)
	}
	return out, today, nil
}

func (a *api) me(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !u.Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in first")
		return
	}
	row, err := store.New(a.Pool).GetUser(r.Context(), u.UUID)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in first")
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	stats, _, err := a.userStats(r.Context(), u.UUID)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, meResponse{User: toUserDTO(row), Stats: stats, Features: featuresDTO{Tutor: a.Tutor != nil}})
}

// updateMeRequest is everything a user may change about themselves. XP,
// level, streak and league fields are not here: only grading writes them.
type updateMeRequest struct {
	Name      *string `json:"name"`
	Timezone  *string `json:"timezone"`
	DailyGoal *int    `json:"daily_goal"`
	Title     *string `json:"title"` // one of the user's earned medal titles, or "" for none
	Goal      *string `json:"goal"`  // onboarding goal
	// Onboarded marks onboarding finished; only true is accepted (it cannot be undone).
	Onboarded *bool `json:"onboarded"`
}

func (a *api) updateMe(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !u.Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in first")
		return
	}
	var req updateMeRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	switch {
	case req.Timezone != nil && !validTimezone(*req.Timezone):
		writeError(w, http.StatusBadRequest, "bad_request", "timezone must be an IANA zone such as Asia/Kolkata")
		return
	case req.DailyGoal != nil && (*req.DailyGoal < 1 || *req.DailyGoal > game.MaxDailyGoal):
		writeError(w, http.StatusBadRequest, "bad_request", fmt.Sprintf("daily_goal must be 1-%d", game.MaxDailyGoal))
		return
	case req.Name != nil && len(*req.Name) > 80:
		writeError(w, http.StatusBadRequest, "bad_request", "name is too long")
		return
	case req.Goal != nil && !onboardingGoals[*req.Goal]:
		writeError(w, http.StatusBadRequest, "bad_request", "goal must be zero, backend, interviews or ai")
		return
	case req.Onboarded != nil && !*req.Onboarded:
		writeError(w, http.StatusBadRequest, "bad_request", "onboarded can only be set to true")
		return
	}
	if req.Title != nil {
		if err := a.Economy.SetTitle(r.Context(), u.UUID, *req.Title); err != nil {
			a.econError(w, r, err)
			return
		}
	}
	err := pgx.BeginFunc(r.Context(), a.Pool, func(tx pgx.Tx) error {
		q := store.New(tx)
		if req.Timezone != nil {
			if err := q.SetUserTimezone(r.Context(), store.SetUserTimezoneParams{ID: u.UUID, Timezone: *req.Timezone}); err != nil {
				return err
			}
		}
		if req.Name != nil {
			if err := q.SetUserName(r.Context(), store.SetUserNameParams{ID: u.UUID, Name: *req.Name}); err != nil {
				return err
			}
		}
		if req.Goal != nil {
			if err := q.SetUserGoal(r.Context(), store.SetUserGoalParams{ID: u.UUID, Goal: *req.Goal}); err != nil {
				return err
			}
		}
		if req.Onboarded != nil {
			if err := q.SetUserOnboarded(r.Context(), u.UUID); err != nil {
				return err
			}
		}
		if req.DailyGoal != nil {
			if err := q.EnsureUserStats(r.Context(), u.UUID); err != nil {
				return err
			}
			return q.SetDailyGoal(r.Context(), store.SetDailyGoalParams{UserID: u.UUID, DailyGoal: int32(*req.DailyGoal)})
		}
		return nil
	})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	a.me(w, r)
}
