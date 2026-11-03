package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/store"
)

type codeLine struct {
	N    int    `json:"n"`
	Text string `json:"text"`
}

// lastAttempt is the user's latest failing answer, for "where you left off".
// The answer is the user's own; nothing server-only appears here.
type lastAttempt struct {
	Status    string     `json:"status"`
	Passed    int        `json:"passed"`
	Total     int        `json:"total"`
	File      string     `json:"file"`
	Lines     []codeLine `json:"lines"`
	ErrorLine int        `json:"error_line"`
}

type continueCard struct {
	CourseID      string       `json:"course_id"`
	CourseTitle   string       `json:"course_title"`
	ChapterTitle  string       `json:"chapter_title"`
	ChapterIndex  int          `json:"chapter_index"`
	ExerciseID    string       `json:"exercise_id"`
	ExerciseTitle string       `json:"exercise_title"`
	ExerciseType  string       `json:"exercise_type"`
	Position      position     `json:"position"`
	Lessons       []string     `json:"lessons"` // statuses of the chapter's lessons, in order
	Last          *lastAttempt `json:"last"`
}

type dayCell struct {
	Date   string `json:"date"` // YYYY-MM-DD, local
	Active bool   `json:"active"`
	Today  bool   `json:"today"`
}

type todayCard struct {
	Goal       int       `json:"goal"`
	Done       int       `json:"done"`
	StreakDays int       `json:"streak_days"`
	Week       []dayCell `json:"week"` // Monday..Sunday of the local week
}

type pathCard struct {
	courseSummary
	Status string `json:"status"` // not_started | in_progress | done
}

type homeResponse struct {
	Continue *continueCard `json:"continue"`
	Today    todayCard     `json:"today"`
	Paths    []pathCard    `json:"paths"`
	League   leagueSummary `json:"league"`
}

func (a *api) home(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !u.Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in first")
		return
	}
	ctx := r.Context()
	stats, today, err := a.userStats(ctx, u.UUID)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	out := homeResponse{Today: todayCard{Goal: stats.DailyGoal, Done: stats.TodayLessons, StreakDays: stats.StreakDays}}
	if out.Today.Week, err = a.week(ctx, u.UUID, today); err != nil {
		a.internal(w, r, err)
		return
	}
	courses, err := store.New(a.Pool).ListCourses(ctx, u)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	for _, c := range courses {
		p := pathCard{courseSummary: courseSummary{
			ID: c.ID, Title: c.Title, Description: c.Description,
			ChapterCount: int(c.ChapterCount), ExerciseCount: int(c.ExerciseCount), PassedCount: int(c.PassedCount),
		}, Status: "not_started"}
		switch {
		case c.ExerciseCount > 0 && c.PassedCount >= c.ExerciseCount:
			p.Status = "done"
		case c.PassedCount > 0:
			p.Status = "in_progress"
		}
		out.Paths = append(out.Paths, p)
	}
	if out.Paths == nil {
		out.Paths = []pathCard{}
	}
	if out.Continue, err = a.continueCard(ctx, u, courses); err != nil {
		a.internal(w, r, err)
		return
	}
	if out.League, err = a.leagueSummary(ctx, u.UUID, stats.LeagueTier); err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (a *api) week(ctx context.Context, user uuid.UUID, today time.Time) ([]dayCell, error) {
	monday := today.AddDate(0, 0, -((int(today.Weekday()) + 6) % 7))
	rows, err := store.New(a.Pool).ListDayPasses(ctx, store.ListDayPassesParams{UserID: user, FromDay: monday, ToDay: monday.AddDate(0, 0, 6)})
	if err != nil {
		return nil, fmt.Errorf("week: %w", err)
	}
	active := map[string]bool{}
	for _, r := range rows {
		active[r.Day.Format(time.DateOnly)] = r.Lessons > 0
	}
	out := make([]dayCell, 7)
	for i := range out {
		d := monday.AddDate(0, 0, i).Format(time.DateOnly)
		out[i] = dayCell{Date: d, Active: active[d], Today: d == today.Format(time.DateOnly)}
	}
	return out, nil
}

// continueCard picks the lesson to resume: the latest submitted one if not
// passed yet, else the next unpassed lesson in that course; with no history,
// the first lesson of the first course.
func (a *api) continueCard(ctx context.Context, u uuid.NullUUID, courses []store.ListCoursesRow) (*continueCard, error) {
	q := store.New(a.Pool)
	last, err := q.LatestActivity(ctx, u.UUID)
	hasLast := err == nil
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("latest activity: %w", err)
	}
	var courseID string
	switch {
	case hasLast:
		courseID = last.CourseID
	case len(courses) > 0:
		courseID = courses[0].ID
	default:
		return nil, nil
	}
	exs, err := q.ListCourseExercises(ctx, store.ListCourseExercisesParams{UserID: u, CourseID: courseID})
	if err != nil || len(exs) == 0 {
		return nil, err
	}
	at := 0
	if hasLast {
		for i, e := range exs {
			if e.ID == last.ExerciseID {
				at = i
			}
		}
		if exs[at].ProgressStatus == "passed" {
			at = nextUnpassed(exs, at)
		}
	}
	target := exs[at]

	chapters, err := q.ListChapters(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("chapters: %w", err)
	}
	c := &continueCard{CourseID: courseID, ExerciseID: target.ID, ExerciseTitle: target.Title, ExerciseType: target.Type, Lessons: []string{}}
	for _, co := range courses {
		if co.ID == courseID {
			c.CourseTitle = co.Title
		}
	}
	for i, ch := range chapters {
		if ch.ID == target.ChapterID {
			c.ChapterTitle, c.ChapterIndex = ch.Title, i+1
		}
	}
	for _, e := range exs {
		if e.ChapterID == target.ChapterID {
			c.Lessons = append(c.Lessons, progressStatus(e.ProgressStatus))
			if e.ID == target.ID {
				c.Position.Index = len(c.Lessons)
			}
		}
	}
	c.Position.Total = len(c.Lessons)
	if hasLast && last.ExerciseID == target.ID && last.Status != exercise.StatusPass {
		c.Last = summarizeAttempt(last.Status, last.Answer, last.Result)
	}
	return c, nil
}

func nextUnpassed(exs []store.ListCourseExercisesRow, from int) int {
	for i := 1; i <= len(exs); i++ {
		if j := (from + i) % len(exs); exs[j].ProgressStatus != "passed" {
			return j
		}
	}
	return from
}

// summarizeAttempt shows up to three lines of the answer ending at the first
// failing line (or the last lines), plus the check counts.
func summarizeAttempt(status string, answer, result json.RawMessage) *lastAttempt {
	out := &lastAttempt{Status: status, Lines: []codeLine{}}
	var res exercise.Result
	if json.Unmarshal(result, &res) == nil {
		out.Total = len(res.Checks)
		for _, c := range res.Checks {
			if c.Passed {
				out.Passed++
			} else if out.ErrorLine == 0 && c.Line > 0 {
				out.ErrorLine = c.Line
			}
		}
	}
	var ans struct {
		Files map[string]string `json:"files"`
		Query string            `json:"query"`
	}
	if json.Unmarshal(answer, &ans) != nil {
		return out
	}
	src := ans.Query
	out.File = "query.sql"
	if len(ans.Files) > 0 {
		names := make([]string, 0, len(ans.Files))
		for n := range ans.Files {
			names = append(names, n)
		}
		sort.Strings(names)
		out.File, src = names[0], ans.Files[names[0]]
	}
	lines := strings.Split(strings.TrimRight(src, "\n"), "\n")
	end := len(lines)
	if out.ErrorLine > 0 && out.ErrorLine <= len(lines) {
		end = out.ErrorLine
	}
	for i := max(0, end-3); i < end; i++ {
		text := lines[i]
		if len(text) > 120 {
			text = text[:120]
		}
		out.Lines = append(out.Lines, codeLine{N: i + 1, Text: text})
	}
	return out
}
