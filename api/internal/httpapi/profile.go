package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/store"
)

// activityDays is how far back the activity heatmap goes (53 weeks, so the
// first column is always a full week).
const activityDays = 53 * 7

type activityDay struct {
	Date    string `json:"date"` // YYYY-MM-DD, the user's local date
	Lessons int    `json:"lessons"`
}

// activityDTO is lessons passed per local day over the last year; only days
// with a pass are listed.
type activityDTO struct {
	From          string        `json:"from"`
	To            string        `json:"to"` // today
	Days          []activityDay `json:"days"`
	LongestStreak int           `json:"longest_streak"`
}

func (a *api) activity(ctx context.Context, user uuid.UUID, today time.Time, longest int) (activityDTO, error) {
	from := today.AddDate(0, 0, -(activityDays - 1))
	rows, err := store.New(a.Pool).ListDayPasses(ctx, store.ListDayPassesParams{UserID: user, FromDay: from, ToDay: today})
	if err != nil {
		return activityDTO{}, fmt.Errorf("activity: %w", err)
	}
	out := activityDTO{From: from.Format(time.DateOnly), To: today.Format(time.DateOnly), Days: []activityDay{}, LongestStreak: longest}
	for _, r := range rows {
		out.Days = append(out.Days, activityDay{Date: r.Day.Format(time.DateOnly), Lessons: int(r.Lessons)})
	}
	return out, nil
}

// projectDTO is a chapter boss: the project a learner ships to finish a chapter.
type projectDTO struct {
	ExerciseID   string     `json:"exercise_id"`
	CourseID     string     `json:"course_id"`
	CourseTitle  string     `json:"course_title"`
	Title        string     `json:"title"`
	ChapterIndex int        `json:"chapter_index"`
	ChapterTitle string     `json:"chapter_title"`
	Status       string     `json:"status"` // defeated | next | upcoming
	FirstTry     bool       `json:"first_try"`
	DefeatedAt   *time.Time `json:"defeated_at"`
}

type skillDTO struct {
	Title   string `json:"title"`
	Tier    string `json:"tier"` // "" (not started) | bronze | silver | gold
	Passed  int    `json:"passed"`
	Lessons int    `json:"lessons"`
}

type skillsDTO struct {
	CourseID    string     `json:"course_id"`
	CourseTitle string     `json:"course_title"`
	Chapters    []skillDTO `json:"chapters"`
}

type recentPassDTO struct {
	Date        string `json:"date"`
	ExerciseID  string `json:"exercise_id"`
	Title       string `json:"title"`
	Boss        bool   `json:"boss"`
	CourseID    string `json:"course_id"`
	CourseTitle string `json:"course_title"`
}

type pathRef struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// profileResponse is the signed-in user's own profile (never public). Name,
// stats and medals come from /me and /medals.
type profileResponse struct {
	JoinedAt        time.Time       `json:"joined_at"`
	Path            *pathRef        `json:"path"` // the path matching the onboarding goal
	ProjectsShipped int             `json:"projects_shipped"`
	Activity        activityDTO     `json:"activity"`
	Projects        []projectDTO    `json:"projects"`
	Skills          *skillsDTO      `json:"skills"`
	Recent          []recentPassDTO `json:"recent"`
}

func (a *api) profile(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !u.Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in first")
		return
	}
	ctx := r.Context()
	out, err := a.buildProfile(ctx, u)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (a *api) buildProfile(ctx context.Context, u uuid.NullUUID) (profileResponse, error) {
	q := store.New(a.Pool)
	var out profileResponse
	joined, err := q.GetUserJoined(ctx, u.UUID)
	if err != nil {
		return out, fmt.Errorf("joined: %w", err)
	}
	out.JoinedAt = joined
	stats, today, err := a.userStats(ctx, u.UUID)
	if err != nil {
		return out, err
	}
	if out.Activity, err = a.activity(ctx, u.UUID, today, stats.LongestStreak); err != nil {
		return out, err
	}
	courses, err := a.courseSummaries(ctx, u)
	if err != nil {
		return out, err
	}
	goal, err := a.userGoal(ctx, u)
	if err != nil {
		return out, err
	}
	paths, err := a.careerPaths(ctx, goal, courses)
	if err != nil {
		return out, err
	}
	for _, p := range paths {
		if p.Yours {
			out.Path = &pathRef{ID: p.ID, Title: p.Title}
		}
	}
	if out.Projects, err = a.projects(ctx, u.UUID, courses); err != nil {
		return out, err
	}
	for _, p := range out.Projects {
		if p.Status == "defeated" {
			out.ProjectsShipped++
		}
	}
	if out.Skills, err = a.skills(ctx, u.UUID, courses, pickPath(paths)); err != nil {
		return out, err
	}
	recent, err := q.RecentPasses(ctx, u.UUID)
	if err != nil {
		return out, fmt.Errorf("recent passes: %w", err)
	}
	out.Recent = make([]recentPassDTO, len(recent))
	for i, p := range recent {
		out.Recent[i] = recentPassDTO{Date: p.Day.Format(time.DateOnly), ExerciseID: p.ID, Title: p.Title, Boss: p.Boss, CourseID: p.CourseID, CourseTitle: p.CourseTitle}
	}
	return out, nil
}

// projects lists every boss: defeated ones, the first undefeated boss of each
// course the user has started ("next"), and the rest ("upcoming").
func (a *api) projects(ctx context.Context, user uuid.UUID, courses []courseSummary) ([]projectDTO, error) {
	rows, err := store.New(a.Pool).ListBosses(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("bosses: %w", err)
	}
	started := map[string]bool{}
	for _, c := range courses {
		started[c.ID] = c.Status != "not_started"
	}
	hasNext := map[string]bool{}
	out := make([]projectDTO, len(rows))
	for i, b := range rows {
		p := projectDTO{
			ExerciseID: b.ID, CourseID: b.CourseID, CourseTitle: b.CourseTitle, Title: b.Title,
			ChapterIndex: int(b.ChapterIndex), ChapterTitle: b.ChapterTitle, Status: "upcoming",
		}
		switch {
		case b.Status == "passed":
			p.Status, p.FirstTry, p.DefeatedAt = "defeated", b.FirstTry, b.FirstPassedAt
		case started[b.CourseID] && !hasNext[b.CourseID]:
			p.Status = "next"
			hasNext[b.CourseID] = true
		}
		out[i] = p
	}
	return out, nil
}

// skills rates each chapter of the course the user is working on: bronze once
// a lesson is passed, silver when every lesson is, gold when every lesson
// passed on the first try. Reviews will refine this when they exist.
func (a *api) skills(ctx context.Context, user uuid.UUID, courses []courseSummary, path *careerPathDTO) (*skillsDTO, error) {
	q := store.New(a.Pool)
	courseID := ""
	last, err := q.LatestActivity(ctx, user)
	switch {
	case err == nil:
		courseID = last.CourseID
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, fmt.Errorf("latest activity: %w", err)
	case path != nil:
		courseID = path.Steps[0].CourseID
	case len(courses) > 0:
		courseID = courses[0].ID
	}
	var title string
	for _, c := range courses {
		if c.ID == courseID {
			title = c.Title
		}
	}
	if title == "" {
		return nil, nil
	}
	rows, err := q.ChapterSkills(ctx, store.ChapterSkillsParams{UserID: user, CourseID: courseID})
	if err != nil {
		return nil, fmt.Errorf("skills: %w", err)
	}
	out := &skillsDTO{CourseID: courseID, CourseTitle: title, Chapters: make([]skillDTO, len(rows))}
	for i, c := range rows {
		s := skillDTO{Title: c.Title, Passed: int(c.Passed), Lessons: int(c.Lessons)}
		switch {
		case c.Lessons > 0 && c.FirstTry == c.Lessons:
			s.Tier = "gold"
		case c.Lessons > 0 && c.Passed == c.Lessons:
			s.Tier = "silver"
		case c.Passed > 0:
			s.Tier = "bronze"
		}
		out.Chapters[i] = s
	}
	return out, nil
}
