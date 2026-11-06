package httpapi_test

import (
	"net/http"
	"testing"
	"time"
)

type catalogCourse struct {
	ID, Icon, Level, Status string
	Tags, Runs              []string
	FreeCount               int `json:"free_count"`
	ChaptersDone            int `json:"chapters_done"`
	CurrentChapter          *struct {
		Index int
		Title string
	} `json:"current_chapter"`
}

type catalog struct {
	Courses []catalogCourse
	Paths   []struct {
		ID    string
		Yours bool
		Steps []struct {
			CourseID string `json:"course_id"`
			Note     string
			Status   string
		}
	}
}

func (c catalog) course(id string) catalogCourse {
	for _, x := range c.Courses {
		if x.ID == id {
			return x
		}
	}
	return catalogCourse{}
}

func TestCoursesCatalogAndPaths(t *testing.T) {
	anon := setup(t, false)
	var out catalog
	anon.getJSON("/api/v1/courses", &out)
	fd, cli := out.course("foundations"), out.course("command-line")
	if fd.Icon != "fd" || fd.Level != "beginner" || len(fd.Tags) != 2 || len(fd.Runs) != 1 || fd.Runs[0] != "browser" || fd.FreeCount == 0 {
		t.Fatalf("foundations = %+v", fd)
	}
	if len(cli.Runs) != 1 || cli.Runs[0] != "cli" {
		t.Fatalf("command-line runs = %v", cli.Runs)
	}
	if len(out.Paths) != 2 || out.Paths[0].ID != "ai-engineer" || out.Paths[0].Yours || out.Paths[0].Steps[0].Note == "" {
		t.Fatalf("paths = %+v", out.Paths)
	}

	e := setup(t, true)
	if code, body := e.do(http.MethodPatch, "/api/v1/me", `{"goal":"ai"}`, nil); code != http.StatusOK {
		t.Fatalf("patch goal = %d %s", code, body)
	}
	e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("b"), nil)
	out = catalog{}
	e.getJSON("/api/v1/courses", &out)
	demo := out.course("demo")
	if demo.Status != "in_progress" || demo.CurrentChapter == nil || demo.CurrentChapter.Index != 1 || demo.CurrentChapter.Title != "Basics" || demo.ChaptersDone != 0 {
		t.Fatalf("demo after one pass = %+v", demo)
	}
	if !out.Paths[0].Yours || out.Paths[1].Yours {
		t.Fatalf("your path should follow the goal: %+v", out.Paths)
	}
}

func TestHomeActivityAndPaths(t *testing.T) {
	now := time.Date(2026, 11, 4, 9, 0, 0, 0, time.UTC)
	e := setupWith(t, true, opts{now: func() time.Time { return now }})
	e.do(http.MethodPatch, "/api/v1/me", `{"goal":"ai"}`, nil)
	e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("b"), nil)

	var home struct {
		Paths []struct {
			ID       string
			Status   string
			NextBoss *string `json:"next_boss"`
			UpNext   bool    `json:"up_next"`
			After    string
		}
		Activity struct {
			From, To string
			Days     []struct {
				Date    string
				Lessons int
			}
			LongestStreak int `json:"longest_streak"`
		}
	}
	body := e.getJSON("/api/v1/home", &home)
	a := home.Activity
	if a.To != "2026-11-04" || a.From != "2025-10-30" || len(a.Days) != 1 || a.Days[0].Date != "2026-11-04" || a.Days[0].Lessons != 1 || a.LongestStreak != 1 {
		t.Fatalf("activity = %+v", a)
	}
	if len(home.Paths) != 2 {
		t.Fatalf("paths: want the started course and the next on the path: %s", body)
	}
	if p := home.Paths[0]; p.ID != "demo" || p.Status != "in_progress" || p.NextBoss == nil || *p.NextBoss != "Sum a list" || p.UpNext {
		t.Fatalf("started course card = %+v", p)
	}
	if p := home.Paths[1]; p.ID != "foundations" || !p.UpNext || p.After != "" {
		t.Fatalf("up next card = %+v", p)
	}
}

func TestProfile(t *testing.T) {
	e := setup(t, true)
	e.do(http.MethodPatch, "/api/v1/me", `{"goal":"backend"}`, nil)
	e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("b"), nil)

	var p struct {
		JoinedAt        time.Time `json:"joined_at"`
		Path            *struct{ ID, Title string }
		ProjectsShipped int `json:"projects_shipped"`
		Activity        struct{ Days []struct{ Lessons int } }
		Projects        []struct {
			ExerciseID string `json:"exercise_id"`
			CourseID   string `json:"course_id"`
			Status     string
		}
		Skills *struct {
			CourseID string `json:"course_id"`
			Chapters []struct {
				Title, Tier     string
				Passed, Lessons int
			}
		}
		Recent []struct {
			ExerciseID string `json:"exercise_id"`
		}
	}
	body := e.getJSON("/api/v1/profile", &p)
	if p.JoinedAt.IsZero() || p.Path == nil || p.Path.ID != "backend-developer" || p.ProjectsShipped != 0 || len(p.Activity.Days) != 1 {
		t.Fatalf("profile = %s", body)
	}
	status := map[string]string{}
	for _, pr := range p.Projects {
		status[pr.ExerciseID] = pr.Status
	}
	// demo is started, so its boss is next; foundations' boss waits.
	if status[codeID] != "next" || status["fd.py.fizzbuzz-boss"] != "upcoming" {
		t.Fatalf("projects = %+v", p.Projects)
	}
	if p.Skills == nil || p.Skills.CourseID != "demo" || p.Skills.Chapters[0].Tier != "bronze" || p.Skills.Chapters[1].Tier != "" {
		t.Fatalf("skills = %+v", p.Skills)
	}
	if len(p.Recent) != 1 || p.Recent[0].ExerciseID != exID {
		t.Fatalf("recent = %+v", p.Recent)
	}
	if code, _ := setup(t, false).do(http.MethodGet, "/api/v1/profile", "", nil); code != http.StatusUnauthorized {
		t.Fatalf("signed-out profile = %d", code)
	}
}
