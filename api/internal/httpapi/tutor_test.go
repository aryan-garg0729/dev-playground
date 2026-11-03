package httpapi_test

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/aryangarg/codequest/api/internal/runner"
	"github.com/aryangarg/codequest/api/internal/tutor"
)

type sseEvent struct {
	name string
	data map[string]any
}

func parseSSE(t *testing.T, body []byte) []sseEvent {
	t.Helper()
	var out []sseEvent
	var cur sseEvent
	sc := bufio.NewScanner(strings.NewReader(string(body)))
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "event: "):
			cur.name = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &cur.data); err != nil {
				t.Fatalf("bad data line %q: %v", line, err)
			}
		case line == "":
			if cur.name != "" {
				out = append(out, cur)
			}
			cur = sseEvent{}
		}
	}
	return out
}

func hintBody(msg string) string {
	b, _ := json.Marshal(map[string]string{"message": msg, "code": "def sum_list(xs):\n    return 0\n"})
	return string(b)
}

func (e *env) usedHint(exercise string) bool {
	e.t.Helper()
	var used bool
	err := e.pool.QueryRow(context.Background(),
		`SELECT used_hint FROM progress WHERE user_id = $1 AND exercise_id = $2`, e.user.UUID, exercise).Scan(&used)
	if err != nil {
		return false
	}
	return used
}

func TestTutorStreamsAndRecords(t *testing.T) {
	fake := &tutor.Fake{Reply: "What does line 2 return for [1, 2]?"}
	e := setupWith(t, true, opts{tutor: fake, tutorLimit: 2, runner: &gatedRunner{out: runner.Output{Stdout: report(false)}}})
	path := "/api/v1/exercises/" + codeID + "/hint"

	// A failed submission first, so the tutor sees the latest Result.
	if code, body := e.do(http.MethodPost, "/api/v1/exercises/"+codeID+"/submit", codeBody("def sum_list(xs):\n    return 0\n"), nil); code != http.StatusOK {
		t.Fatalf("submit = %d: %s", code, body)
	}

	code, body := e.do(http.MethodPost, path, hintBody("why does it fail?"), nil)
	if code != http.StatusOK {
		t.Fatalf("hint = %d: %s", code, body)
	}
	evs := parseSSE(t, body)
	var text strings.Builder
	for _, ev := range evs[:len(evs)-1] {
		if ev.name != "delta" {
			t.Fatalf("event %q before done", ev.name)
		}
		text.WriteString(ev.data["text"].(string))
	}
	last := evs[len(evs)-1]
	if text.String() != fake.Reply || last.name != "done" || last.data["remaining"] != float64(1) {
		t.Fatalf("stream = %q, last = %+v", text.String(), last)
	}
	if !e.usedHint(codeID) {
		t.Fatal("a tutor hint must set progress.used_hint")
	}

	// The model got the prompt and the learner's code, and never server-only files.
	reqs := fake.Requests()
	if len(reqs) != 1 {
		t.Fatalf("model calls = %d", len(reqs))
	}
	all := reqs[0].System
	for _, m := range reqs[0].Messages {
		all += m.Content
	}
	if strings.Contains(all, sentinel) {
		t.Fatalf("tutor request contains server-only content:\n%s", all)
	}
	if !strings.Contains(all, "return 0") || !strings.Contains(all, "why does it fail?") || !strings.Contains(all, "Latest result: fail") {
		t.Fatalf("tutor request lacks the code or question:\n%s", all)
	}

	// The thread is saved and sent as history next time.
	var th struct {
		DailyLimit int `json:"daily_limit"`
		Remaining  int `json:"remaining"`
		Messages   []struct{ Role, Content string }
	}
	e.getJSON("/api/v1/exercises/"+codeID+"/tutor", &th)
	if th.DailyLimit != 2 || th.Remaining != 1 || len(th.Messages) != 2 || th.Messages[0].Role != "user" || th.Messages[1].Content != fake.Reply {
		t.Fatalf("thread = %+v", th)
	}
	if code, body := e.do(http.MethodPost, path, hintBody("and now?"), nil); code != http.StatusOK {
		t.Fatalf("second hint = %d: %s", code, body)
	}
	if h := fake.Requests()[1].Messages; len(h) != 3 || h[0].Content != "why does it fail?" {
		t.Fatalf("history sent = %+v", h)
	}
}

func TestTutorDailyLimit(t *testing.T) {
	e := setupWith(t, true, opts{tutor: &tutor.Fake{Reply: "Look at the loop."}, tutorLimit: 2})
	path := "/api/v1/exercises/" + codeID + "/hint"
	for i := range 2 {
		if code, body := e.do(http.MethodPost, path, hintBody("hint please"), nil); code != http.StatusOK {
			t.Fatalf("hint %d = %d: %s", i+1, code, body)
		}
	}
	code, body := e.do(http.MethodPost, path, hintBody("one more"), nil)
	if code != http.StatusTooManyRequests || !strings.Contains(string(body), "limit_reached") {
		t.Fatalf("third hint = %d: %s", code, body)
	}
	// Another user has their own limit.
	e.user, e.cookie = e.signIn("other@example.com", "O")
	if code, body := e.do(http.MethodPost, path, hintBody("hint"), nil); code != http.StatusOK {
		t.Fatalf("other user = %d: %s", code, body)
	}
}

func TestTutorFailureRefundsHint(t *testing.T) {
	fake := &tutor.Fake{Err: errors.New("model down")}
	e := setupWith(t, true, opts{tutor: fake, tutorLimit: 1})
	path := "/api/v1/exercises/" + codeID + "/hint"
	code, body := e.do(http.MethodPost, path, hintBody("help"), nil)
	evs := parseSSE(t, body)
	if code != http.StatusOK || len(evs) != 1 || evs[0].name != "error" || evs[0].data["code"] != "tutor_failed" {
		t.Fatalf("failed hint = %d: %s", code, body)
	}
	if e.usedHint(codeID) {
		t.Fatal("a failed hint must not remove the no-hint bonus")
	}
	fake.Err = nil
	fake.Reply = "ok"
	if code, body := e.do(http.MethodPost, path, hintBody("help"), nil); code != http.StatusOK || !strings.Contains(string(body), "event: done") {
		t.Fatalf("hint after refund = %d: %s", code, body)
	}
}

func TestTutorOffAndRejects(t *testing.T) {
	off := setup(t, true)
	var me struct {
		Features struct{ Tutor bool } `json:"features"`
	}
	off.getJSON("/api/v1/me", &me)
	if me.Features.Tutor {
		t.Fatal("features.tutor must be false without a provider")
	}
	for _, c := range []struct{ method, path, body string }{
		{http.MethodPost, "/api/v1/exercises/" + codeID + "/hint", hintBody("hi")},
		{http.MethodGet, "/api/v1/exercises/" + codeID + "/tutor", ""},
	} {
		if code, body := off.do(c.method, c.path, c.body, nil); code != http.StatusNotFound || !strings.Contains(string(body), "tutor_unavailable") {
			t.Fatalf("%s %s with tutor off = %d: %s", c.method, c.path, code, body)
		}
	}

	e := setupWith(t, false, opts{tutor: &tutor.Fake{Reply: "x"}})
	path := "/api/v1/exercises/" + codeID + "/hint"
	if code, _ := e.do(http.MethodPost, path, hintBody("hi"), nil); code != http.StatusUnauthorized {
		t.Fatalf("signed out = %d", code)
	}
	e.user, e.cookie = e.signIn("t@example.com", "T")
	e.getJSON("/api/v1/me", &me)
	if !me.Features.Tutor {
		t.Fatal("features.tutor must be true with a provider")
	}
	for name, c := range map[string]struct {
		path, body string
		want       int
	}{
		"empty message":    {path, `{"message":"  ","code":""}`, http.StatusBadRequest},
		"long message":     {path, `{"message":"` + strings.Repeat("a", 2001) + `"}`, http.StatusBadRequest},
		"unknown field":    {path, `{"message":"hi","solution":"x"}`, http.StatusBadRequest},
		"unknown exercise": {"/api/v1/exercises/nope/hint", hintBody("hi"), http.StatusNotFound},
	} {
		if code, body := e.do(http.MethodPost, c.path, c.body, nil); code != c.want {
			t.Errorf("%s = %d, want %d: %s", name, code, c.want, body)
		}
	}
}
