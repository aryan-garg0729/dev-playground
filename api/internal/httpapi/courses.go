package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/store"
)

type chapterPos struct {
	Index int    `json:"index"` // 1-based
	Title string `json:"title"`
}

type courseSummary struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Icon          string   `json:"icon"`
	Level         string   `json:"level"` // beginner | intermediate | advanced
	Tags          []string `json:"tags"`
	Topics        string   `json:"topics"`
	Runs          []string `json:"runs"` // where lessons run: "browser" and/or "cli"
	ChapterCount  int      `json:"chapter_count"`
	ExerciseCount int      `json:"exercise_count"`
	PassedCount   int      `json:"passed_count"`
	ChaptersDone  int      `json:"chapters_done"`
	FreeCount     int      `json:"free_count"`
	Status        string   `json:"status"` // not_started | in_progress | done
	// The chapter of the first unpassed lesson, while the course is in progress.
	CurrentChapter *chapterPos `json:"current_chapter"`
}

type pathStepDTO struct {
	CourseID string `json:"course_id"`
	Title    string `json:"title"`
	Note     string `json:"note"`
	Status   string `json:"status"` // the course's status
}

type careerPathDTO struct {
	ID            string        `json:"id"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	Yours         bool          `json:"yours"` // matches the user's onboarding goal
	Steps         []pathStepDTO `json:"steps"`
	ExerciseCount int           `json:"exercise_count"`
	PassedCount   int           `json:"passed_count"`
}

// courseSummaries lists published courses with the user's progress (u may be
// invalid: signed-out visitors see the catalog without progress).
func (a *api) courseSummaries(ctx context.Context, u uuid.NullUUID) ([]courseSummary, error) {
	q := store.New(a.Pool)
	rows, err := q.ListCourses(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("courses: %w", err)
	}
	current := map[string]chapterPos{}
	if u.Valid {
		cur, err := q.CourseCurrentChapters(ctx, u.UUID)
		if err != nil {
			return nil, fmt.Errorf("current chapters: %w", err)
		}
		for _, c := range cur {
			current[c.CourseID] = chapterPos{Index: int(c.ChapterIndex), Title: c.Title}
		}
	}
	out := make([]courseSummary, len(rows))
	for i, c := range rows {
		s := courseSummary{
			ID: c.ID, Title: c.Title, Description: c.Description, Icon: c.Icon, Level: c.Level,
			Tags: c.Tags, Topics: c.Topics, Runs: []string{},
			ChapterCount: int(c.ChapterCount), ExerciseCount: int(c.ExerciseCount), PassedCount: int(c.PassedCount),
			ChaptersDone: int(c.ChaptersDone), FreeCount: int(c.FreeCount), Status: "not_started",
		}
		if s.Tags == nil {
			s.Tags = []string{}
		}
		if c.HasBrowser {
			s.Runs = append(s.Runs, "browser")
		}
		if c.HasCli {
			s.Runs = append(s.Runs, "cli")
		}
		switch {
		case c.ExerciseCount > 0 && c.PassedCount >= c.ExerciseCount:
			s.Status = "done"
		case c.PassedCount > 0:
			s.Status = "in_progress"
			if ch, ok := current[c.ID]; ok {
				s.CurrentChapter = &ch
			}
		}
		out[i] = s
	}
	return out, nil
}

// careerPaths lists published paths over the given course summaries. Steps
// whose course is no longer published are left out.
func (a *api) careerPaths(ctx context.Context, goal string, courses []courseSummary) ([]careerPathDTO, error) {
	rows, err := store.New(a.Pool).ListCareerPaths(ctx)
	if err != nil {
		return nil, fmt.Errorf("career paths: %w", err)
	}
	byID := map[string]courseSummary{}
	for _, c := range courses {
		byID[c.ID] = c
	}
	out := []careerPathDTO{}
	for _, r := range rows {
		var steps []struct{ Course, Note string }
		if err := json.Unmarshal(r.Steps, &steps); err != nil {
			return nil, fmt.Errorf("path %s steps: %w", r.ID, err)
		}
		p := careerPathDTO{ID: r.ID, Title: r.Title, Description: r.Description, Yours: goal != "" && r.Goal == goal, Steps: []pathStepDTO{}}
		for _, st := range steps {
			c, ok := byID[st.Course]
			if !ok {
				continue
			}
			p.Steps = append(p.Steps, pathStepDTO{CourseID: c.ID, Title: c.Title, Note: st.Note, Status: c.Status})
			p.ExerciseCount += c.ExerciseCount
			p.PassedCount += c.PassedCount
		}
		if len(p.Steps) > 0 {
			out = append(out, p)
		}
	}
	return out, nil
}

// userGoal is the signed-in user's onboarding goal ("" when signed out).
func (a *api) userGoal(ctx context.Context, u uuid.NullUUID) (string, error) {
	if !u.Valid {
		return "", nil
	}
	row, err := store.New(a.Pool).GetUser(ctx, u.UUID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("user: %w", err)
	}
	return row.Goal, nil
}

func (a *api) listCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := currentUser(r)
	courses, err := a.courseSummaries(ctx, u)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	goal, err := a.userGoal(ctx, u)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	paths, err := a.careerPaths(ctx, goal, courses)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"courses": courses, "paths": paths})
}

type exerciseSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	XP         int    `json:"xp"`
	Difficulty int    `json:"difficulty"`
	Free       bool   `json:"free"`
	Boss       bool   `json:"boss"`
	Status     string `json:"status"` // not_started | started | passed
}

type chapterDTO struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Locked    bool              `json:"locked"` // behind an earlier chapter's undefeated boss
	Exercises []exerciseSummary `json:"exercises"`
}

type courseDetail struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Chapters    []chapterDTO `json:"chapters"`
	Skills      []string     `json:"skills"`    // distinct lesson tags, in course order
	XPEarned    int          `json:"xp_earned"` // from the user's XP ledger
	XPTotal     int          `json:"xp_total"`  // base XP of all lessons (bonuses not included)
}

func (a *api) getCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := store.New(a.Pool)
	id := chi.URLParam(r, "courseID")
	c, err := q.GetCourse(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "not_found", "course not found")
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	chapters, err := q.ListChapters(ctx, id)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	exs, err := q.ListCourseExercises(ctx, store.ListCourseExercisesParams{UserID: currentUser(r), CourseID: id})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	byChapter := map[string][]exerciseSummary{}
	for _, e := range exs {
		byChapter[e.ChapterID] = append(byChapter[e.ChapterID], exerciseSummary{
			ID: e.ID, Title: e.Title, Type: e.Type, XP: int(e.Xp), Difficulty: int(e.Difficulty),
			Free: e.Free, Boss: e.Boss, Status: progressStatus(e.ProgressStatus),
		})
	}
	xp, err := q.CourseXP(ctx, store.CourseXPParams{UserID: currentUser(r), CourseID: id})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	skills, err := q.CourseTags(ctx, id)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	if skills == nil {
		skills = []string{}
	}
	locked, err := a.Economy.LockedChapters(ctx, currentUser(r), id)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	out := courseDetail{
		ID: c.ID, Title: c.Title, Description: c.Description, Chapters: make([]chapterDTO, 0, len(chapters)),
		Skills: skills, XPEarned: int(xp.Earned), XPTotal: int(xp.Total),
	}
	for _, ch := range chapters {
		list := byChapter[ch.ID]
		if list == nil {
			list = []exerciseSummary{}
		}
		out.Chapters = append(out.Chapters, chapterDTO{ID: ch.ID, Title: ch.Title, Locked: locked[ch.ID], Exercises: list})
	}
	writeJSON(w, http.StatusOK, out)
}

func progressStatus(s string) string {
	if s == "" {
		return "not_started"
	}
	return s
}
