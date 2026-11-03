package httpapi

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/store"
)

type courseSummary struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ChapterCount  int    `json:"chapter_count"`
	ExerciseCount int    `json:"exercise_count"`
	PassedCount   int    `json:"passed_count"`
}

func (a *api) listCourses(w http.ResponseWriter, r *http.Request) {
	rows, err := store.New(a.Pool).ListCourses(r.Context(), currentUser(r))
	if err != nil {
		a.internal(w, r, err)
		return
	}
	out := make([]courseSummary, len(rows))
	for i, c := range rows {
		out[i] = courseSummary{
			ID: c.ID, Title: c.Title, Description: c.Description,
			ChapterCount: int(c.ChapterCount), ExerciseCount: int(c.ExerciseCount), PassedCount: int(c.PassedCount),
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"courses": out})
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
	out := courseDetail{
		ID: c.ID, Title: c.Title, Description: c.Description, Chapters: make([]chapterDTO, 0, len(chapters)),
		Skills: skills, XPEarned: int(xp.Earned), XPTotal: int(xp.Total),
	}
	for _, ch := range chapters {
		list := byChapter[ch.ID]
		if list == nil {
			list = []exerciseSummary{}
		}
		out.Chapters = append(out.Chapters, chapterDTO{ID: ch.ID, Title: ch.Title, Exercises: list})
	}
	writeJSON(w, http.StatusOK, out)
}

func progressStatus(s string) string {
	if s == "" {
		return "not_started"
	}
	return s
}
