package httpapi_test

import (
	"net/http"
	"testing"
	"time"
)

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
}
