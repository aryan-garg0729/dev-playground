package httpapi_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func (e *env) ledger() (events, total int) {
	e.t.Helper()
	if err := e.pool.QueryRow(e.t.Context(), `SELECT count(*), COALESCE(sum(amount), 0) FROM xp_events WHERE user_id = $1`, e.user.UUID).Scan(&events, &total); err != nil {
		e.t.Fatal(err)
	}
	return events, total
}

type submitOut struct {
	XPAwarded int `json:"xp_awarded"`
	Outcome   struct {
		XPAwarded int  `json:"xp_awarded"`
		FirstTry  bool `json:"first_try"`
		TotalXP   int  `json:"total_xp"`
		Breakdown *struct {
			Base, FirstTry, NoHint int
		} `json:"breakdown"`
	} `json:"outcome"`
}

// Criterion: forged client requests cannot award XP. Only a server-graded
// first pass writes the ledger; nothing the client sends can.
func TestForgedRequestsCannotAwardXP(t *testing.T) {
	e := setup(t, true)
	submit := "/api/v1/exercises/" + exID + "/submit"

	forged := []struct{ name, body string }{
		{"claims a pass", `{"exercise_id":"` + exID + `","version":2,"answer":{"choices":["a"]},"status":"pass"}`},
		{"claims XP", `{"exercise_id":"` + exID + `","version":2,"answer":{"choices":["a"]},"xp_awarded":1000}`},
		{"sends a result", `{"exercise_id":"` + exID + `","version":2,"answer":{"choices":["a"]},"result":{"status":"pass","score":1}}`},
	}
	for _, f := range forged {
		if code, body := e.do(http.MethodPost, submit, f.body, nil); code != http.StatusBadRequest {
			t.Errorf("%s: %d %s, want 400", f.name, code, body)
		}
	}
	// No endpoint takes XP, level, streak or league fields.
	for _, body := range []string{`{"total_xp":9999}`, `{"level":50}`, `{"streak_days":365}`, `{"league_tier":"diamond"}`} {
		if code, _ := e.do(http.MethodPatch, "/api/v1/me", body, nil); code != http.StatusBadRequest {
			t.Errorf("PATCH /me %s = %d, want 400", body, code)
		}
	}
	// Run never awards XP (and mcq has no server Run).
	e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/run", submitBody("b"), nil)
	// A wrong answer is graded as a fail: no XP.
	var out submitOut
	code, body := e.do(http.MethodPost, submit, submitBody("a"), nil)
	if code != http.StatusOK || json.Unmarshal(body, &out) != nil || out.XPAwarded != 0 {
		t.Fatalf("wrong answer = %d %s", code, body)
	}
	// Someone else's session cannot be used without the cookie; no session, no XP.
	anon := *e
	anon.cookie = ""
	if code, _ := anon.do(http.MethodPost, submit, submitBody("b"), nil); code != http.StatusUnauthorized {
		t.Fatalf("anonymous submit = %d", code)
	}
	if n, _ := e.ledger(); n != 0 {
		t.Fatalf("ledger has %d events before any real pass", n)
	}

	// The real pass (second attempt): base 5 + no-hint 1, no first-try bonus.
	code, body = e.do(http.MethodPost, submit, submitBody("b"), nil)
	if code != http.StatusOK || json.Unmarshal(body, &out) != nil || out.XPAwarded != 6 || out.Outcome.FirstTry || out.Outcome.TotalXP != 6 {
		t.Fatalf("pass = %d %s", code, body)
	}
	// Replaying the passing request earns nothing more.
	for range 3 {
		code, body = e.do(http.MethodPost, submit, submitBody("b"), nil)
		if code != http.StatusOK || json.Unmarshal(body, &out) != nil || out.XPAwarded != 0 || out.Outcome.Breakdown != nil {
			t.Fatalf("replayed pass = %d %s", code, body)
		}
	}
	if n, total := e.ledger(); n != 1 || total != 6 {
		t.Fatalf("ledger = %d events, %d XP; want 1 and 6", n, total)
	}
}

// Parallel passing submits (double-click, retries) award once.
func TestParallelSubmitsAwardOnce(t *testing.T) {
	e := setup(t, true)
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() { e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("b"), nil) })
	}
	wg.Wait()
	if n, total := e.ledger(); n != 1 || total != 9 {
		t.Fatalf("ledger = %d events, %d XP; want 1 and 9 (first try)", n, total)
	}
}

func TestMeHomeAndLeague(t *testing.T) {
	now := time.Date(2026, 11, 4, 9, 0, 0, 0, time.UTC) // Wednesday
	e := setupWith(t, true, opts{now: func() time.Time { return now }})

	var league struct {
		Joined bool
		Tier   string
	}
	e.getJSON("/api/v1/leagues/current", &league)
	if league.Joined || league.Tier != "bronze" {
		t.Fatalf("league before XP = %+v", league)
	}

	e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("b"), nil)

	var me struct {
		Stats struct {
			TotalXP      int `json:"total_xp"`
			Level        int
			NextLevelXP  int `json:"next_level_xp"`
			StreakDays   int `json:"streak_days"`
			DailyGoal    int `json:"daily_goal"`
			TodayLessons int `json:"today_lessons"`
		}
	}
	e.getJSON("/api/v1/me", &me)
	if s := me.Stats; s.TotalXP != 9 || s.Level != 1 || s.NextLevelXP != 100 || s.StreakDays != 1 || s.DailyGoal != 2 || s.TodayLessons != 1 {
		t.Fatalf("me.stats = %+v", s)
	}

	var home struct {
		Continue *struct {
			CourseID   string `json:"course_id"`
			ExerciseID string `json:"exercise_id"`
			Lessons    []string
		}
		Today struct {
			Goal, Done int
			Week       []struct {
				Date          string
				Active, Today bool
			}
		}
		Paths  []struct{ ID, Status string }
		League struct {
			Joined   bool
			Rank, Of int
			WeeklyXP int `json:"weekly_xp"`
			Zone     string
		}
	}
	body := e.getJSON("/api/v1/home", &home)
	if home.Continue == nil || home.Continue.ExerciseID == exID {
		t.Fatalf("continue should move past the passed lesson: %s", body)
	}
	if len(home.Today.Week) != 7 || home.Today.Week[0].Date != "2026-11-02" || !home.Today.Week[2].Active || !home.Today.Week[2].Today || home.Today.Done != 1 {
		t.Fatalf("today = %+v", home.Today)
	}
	// Weekly XP is not checked: xp_events.created_at is the database's real clock, the week is the test's fake one.
	if !home.League.Joined || home.League.Rank != 1 || home.League.Of != 1 || home.League.Zone != "promote" {
		t.Fatalf("league = %+v", home.League)
	}
	if strings.Contains(string(body), "SERVER_ONLY_SENTINEL") {
		t.Fatal("home leaked server-only content")
	}

	// A user's settings: daily goal and timezone, validated.
	if code, _ := e.do(http.MethodPatch, "/api/v1/me", `{"daily_goal":3,"timezone":"Asia/Kolkata"}`, nil); code != http.StatusOK {
		t.Fatalf("patch me = %d", code)
	}
	for _, bad := range []string{`{"daily_goal":0}`, `{"daily_goal":11}`, `{"timezone":"Mars/Olympus"}`} {
		if code, _ := e.do(http.MethodPatch, "/api/v1/me", bad, nil); code != http.StatusBadRequest {
			t.Errorf("PATCH %s = %d", bad, code)
		}
	}
}

func TestCourseXP(t *testing.T) {
	e := setup(t, true)
	e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("b"), nil)
	var c struct {
		Skills   []string
		XPEarned int `json:"xp_earned"`
		XPTotal  int `json:"xp_total"`
	}
	e.getJSON("/api/v1/courses/demo", &c)
	if c.XPEarned != 9 || c.XPTotal != 40 {
		t.Fatalf("course xp = %+v", c)
	}
}
