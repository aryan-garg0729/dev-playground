package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/content"
	"github.com/aryangarg/codequest/api/internal/economy"
	_ "github.com/aryangarg/codequest/api/internal/exercise/all"
	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/grading"
	"github.com/aryangarg/codequest/api/internal/httpapi"
	"github.com/aryangarg/codequest/api/internal/ratelimit"
	"github.com/aryangarg/codequest/api/internal/runner"
	"github.com/aryangarg/codequest/api/internal/testdb"
	"github.com/aryangarg/codequest/api/internal/testredis"
	"github.com/aryangarg/codequest/api/internal/tutor"
)

const (
	sentinel = "SERVER_ONLY_SENTINEL" // appears only in server-only files and the mcq explanation
	origin   = "http://localhost:5173"
	exID     = "demo.basics.pick"
)

type env struct {
	t      *testing.T
	pool   *pgxpool.Pool
	h      http.Handler
	user   uuid.NullUUID
	cookie string // session cookie value when signed in
	rdb    *redis.Client
	prefix string
}

// opts tweak setup for one test.
type opts struct {
	runner        runner.Runner
	wait          time.Duration // grading wait before 202
	submitPerUser int
	devAuth       bool
	github        *auth.GitHub
	google        *auth.Google
	now           func() time.Time
	rand          func(n int) int
	tutor         tutor.Provider // nil: tutor off
	tutorLimit    int            // hints per day (default 3)
}

func setup(t *testing.T, withUser bool) *env {
	t.Helper()
	return setupWith(t, withUser, opts{})
}

func setupWith(t *testing.T, withUser bool, o opts) *env {
	t.Helper()
	pool, _ := testdb.New(t)
	ctx := context.Background()
	// Build publishes exactly the bundle it is given, so merge the sample
	// content and the fixture (whose server-only files carry the sentinel).
	bundle := &content.Bundle{}
	for _, dir := range []string{"../../../content", "../content/testdata/valid"} {
		b, err := content.Load(dir)
		if err != nil {
			t.Fatalf("load %s: %v", dir, err)
		}
		bundle.Courses = append(bundle.Courses, b.Courses...)
		bundle.Paths = append(bundle.Paths, b.Paths...)
	}
	if _, err := content.Build(ctx, pool, bundle); err != nil {
		t.Fatalf("build: %v", err)
	}
	rdb, prefix := testredis.New(t)
	board := &game.Board{Redis: rdb, Prefix: prefix}
	econ := &economy.Service{Pool: pool, Now: o.now, Rand: o.rand, DevUnlockAll: o.devAuth}
	if err := econ.Sync(ctx); err != nil {
		t.Fatal(err)
	}
	grader := grading.NewService(pool, o.runner, grading.Options{
		Wait: o.wait, Log: slog.New(slog.DiscardHandler), Board: board, Game: &game.Engine{Now: o.now, Rewards: econ},
	})
	t.Cleanup(grader.Close)
	sessions := &auth.Sessions{Pool: pool}
	var quota *tutor.Quota
	if o.tutor != nil {
		if o.tutorLimit == 0 {
			o.tutorLimit = 3
		}
		quota = &tutor.Quota{Redis: rdb, Prefix: prefix, Limit: o.tutorLimit}
	}
	h := httpapi.New(httpapi.Deps{
		Tutor: o.tutor, TutorQuota: quota,
		Log: slog.New(slog.DiscardHandler), Pool: pool, Grader: grader, AllowedOrigins: []string{origin},
		Sessions: sessions,
		GitHub:   o.github, Google: o.google, DevAuth: o.devAuth, Board: board, Economy: econ, Now: o.now,
		Limiter: &ratelimit.Limiter{Redis: rdb, Prefix: prefix}, SubmitPerUser: o.submitPerUser,
	})
	e := &env{t: t, pool: pool, h: h, rdb: rdb, prefix: prefix}
	if withUser {
		e.user, e.cookie = e.signIn("t@example.com", "T")
	}
	return e
}

// signIn creates (or finds) a user and a session for them.
func (e *env) signIn(email, name string) (uuid.NullUUID, string) {
	e.t.Helper()
	ctx := context.Background()
	u, err := auth.UserForEmail(ctx, e.pool, email, name, "")
	if err != nil {
		e.t.Fatal(err)
	}
	raw, _, err := (&auth.Sessions{Pool: e.pool}).Create(ctx, u.ID)
	if err != nil {
		e.t.Fatal(err)
	}
	return uuid.NullUUID{UUID: u.ID, Valid: true}, raw
}

func (e *env) do(method, path, body string, hdr map[string]string) (int, []byte) {
	e.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", origin)
	}
	if e.cookie != "" {
		req.AddCookie(&http.Cookie{Name: "cq_session", Value: e.cookie})
	}
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	e.h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.Bytes()
}

func (e *env) getJSON(path string, v any) []byte {
	e.t.Helper()
	code, body := e.do(http.MethodGet, path, "", nil)
	if code != http.StatusOK {
		e.t.Fatalf("GET %s = %d: %s", path, code, body)
	}
	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			e.t.Fatalf("decode %s: %v", path, err)
		}
	}
	return body
}

func submitBody(choices ...string) string {
	b, _ := json.Marshal(map[string]any{"exercise_id": exID, "version": 2, "answer": map[string]any{"choices": choices}})
	return string(b)
}

type courseDetail struct {
	Chapters []struct {
		ID        string `json:"id"`
		Locked    bool   `json:"locked"`
		Exercises []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"exercises"`
	} `json:"chapters"`
}

// forbiddenKeys must never appear in any browser-facing response.
var forbiddenKeys = []string{"private_spec", "explain", "solution", "tests_hidden", "tests", "wrong", "checks"}

func findKeys(v any, path string, out *[]string) {
	switch t := v.(type) {
	case map[string]any:
		for k, child := range t {
			for _, f := range forbiddenKeys {
				// The learner's own last result legitimately has "checks".
				if k == f && (k != "checks" || !strings.Contains(path, "last_submission")) {
					*out = append(*out, path+"."+k)
				}
			}
			// "answer" is allowed only as the learner's own last answer.
			if k == "answer" && !strings.HasSuffix(path, "last_submission") {
				*out = append(*out, path+"."+k)
			}
			findKeys(child, path+"."+k, out)
		}
	case []any:
		for _, child := range t {
			findKeys(child, path+"[]", out)
		}
	}
}

func TestNoServerOnlyContentInAnyResponse(t *testing.T) {
	for _, withUser := range []bool{false, true} {
		e := setup(t, withUser)
		var list struct {
			Courses []struct{ ID string } `json:"courses"`
		}
		bodies := map[string][]byte{"/api/v1/courses": e.getJSON("/api/v1/courses", &list)}
		if len(list.Courses) < 2 {
			t.Fatalf("want sample course and fixture, got %+v", list.Courses)
		}
		for _, c := range list.Courses {
			var cd courseDetail
			p := "/api/v1/courses/" + c.ID
			bodies[p] = e.getJSON(p, &cd)
			for _, ch := range cd.Chapters {
				for _, ex := range ch.Exercises {
					p := "/api/v1/exercises/" + ex.ID
					code, body := e.do(http.MethodGet, p, "", nil)
					switch {
					case code == http.StatusOK:
					case code == http.StatusForbidden && ch.Locked && bytes.Contains(body, []byte("chapter_locked")):
					default:
						t.Fatalf("GET %s = %d: %s", p, code, body)
					}
					bodies[p] = body
				}
			}
		}
		if withUser {
			for _, p := range []string{"/api/v1/me", "/api/v1/home", "/api/v1/profile", "/api/v1/shop", "/api/v1/quests", "/api/v1/medals", "/api/v1/leagues/current"} {
				bodies[p] = e.getJSON(p, nil)
			}
		}
		if _, ok := bodies["/api/v1/exercises/"+exID]; !ok {
			t.Fatal("fixture exercise was not crawled")
		}
		for path, body := range bodies {
			if bytes.Contains(body, []byte(sentinel)) {
				t.Errorf("%s leaks server-only content: %s", path, body)
			}
			var v any
			if err := json.Unmarshal(body, &v); err != nil {
				t.Fatal(err)
			}
			var bad []string
			findKeys(v, "", &bad)
			if len(bad) > 0 {
				t.Errorf("%s exposes forbidden keys %v", path, bad)
			}
		}
	}
}

func TestCourseListsManifestOrder(t *testing.T) {
	e := setup(t, false)
	var cd courseDetail
	e.getJSON("/api/v1/courses/langchain-basics", &cd)
	var got []string
	for _, ch := range cd.Chapters {
		got = append(got, ch.ID)
		for _, ex := range ch.Exercises {
			got = append(got, ex.ID)
		}
	}
	want := []string{
		"lc.chat-models", "lc.chat.message-roles", "lc.chat.build-messages", "lc.chat.runnable-methods", "lc.chat.stream-tokens",
		"lc.chat.stateless-models", "lc.chat.continue-conversation", "lc.chat.token-usage",
		"lc.prompt-templates", "lc.prompts.input-variables", "lc.prompts.translation-prompt", "lc.prompts.messages-placeholder",
		"lc.prompts.escape-braces", "lc.prompts.partial-functions", "lc.prompts.dated-prompt",
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("order:\n got %v\nwant %v", got, want)
	}
}

func TestSubmitFlow(t *testing.T) {
	e := setup(t, true)

	code, body := e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("a"), nil)
	if code != http.StatusOK {
		t.Fatalf("wrong answer submit = %d: %s", code, body)
	}
	var res struct {
		SubmissionID string `json:"submission_id"`
		Result       struct {
			Status string `json:"status"`
			Checks []struct {
				Passed   bool   `json:"passed"`
				Feedback string `json:"feedback"`
			} `json:"checks"`
		} `json:"result"`
		XPAwarded int `json:"xp_awarded"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}
	if res.Result.Status != "fail" || !strings.Contains(res.Result.Checks[0].Feedback, "revealed only after submit") {
		t.Fatalf("want fail with explanation, got %s", body)
	}

	code, body = e.do(http.MethodPost, "/api/v1/exercises/"+exID+"/submit", submitBody("b"), nil)
	if code != http.StatusOK || !strings.Contains(string(body), `"status":"pass"`) {
		t.Fatalf("right answer = %d: %s", code, body)
	}

	var cd courseDetail
	e.getJSON("/api/v1/courses/demo", &cd)
	if st := cd.Chapters[0].Exercises[0].Status; st != "passed" {
		t.Fatalf("progress status = %q, want passed", st)
	}
	var detail struct {
		Last *struct {
			Result struct {
				Status string `json:"status"`
			} `json:"result"`
		} `json:"last_submission"`
	}
	e.getJSON("/api/v1/exercises/"+exID, &detail)
	if detail.Last == nil || detail.Last.Result.Status != "pass" {
		t.Fatalf("last_submission not returned: %+v", detail.Last)
	}

	var attempts int
	if err := e.pool.QueryRow(context.Background(),
		`SELECT attempts FROM progress WHERE exercise_id = $1`, exID).Scan(&attempts); err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestSubmitRejects(t *testing.T) {
	e := setup(t, true)
	anon := setup(t, false)
	path := "/api/v1/exercises/" + exID + "/submit"
	tests := []struct {
		name     string
		env      *env
		path     string
		body     string
		hdr      map[string]string
		wantCode int
		wantErr  string
	}{
		{"no user", anon, path, submitBody("b"), nil, 401, "unauthorized"},
		{"stale version", e, path, strings.Replace(submitBody("b"), `"version":2`, `"version":1`, 1), nil, 409, "stale_version"},
		{"invalid option", e, path, submitBody("z"), nil, 400, "invalid_answer"},
		{"two choices on single", e, path, submitBody("a", "b"), nil, 400, "invalid_answer"},
		{"id mismatch", e, "/api/v1/exercises/lc.chat.message-roles/submit", submitBody("b"), nil, 400, "bad_request"},
		{"unknown field", e, path, `{"exercise_id":"` + exID + `","version":2,"answer":{"choices":["b"]},"passed":true}`, nil, 400, "bad_request"},
		{"missing answer", e, path, `{"exercise_id":"` + exID + `","version":2}`, nil, 400, "bad_request"},
		{"unknown exercise", e, "/api/v1/exercises/nope/submit", strings.Replace(submitBody("b"), exID, "nope", 1), nil, 404, "not_found"},
		{"foreign origin", e, path, submitBody("b"), map[string]string{"Origin": "https://evil.example"}, 403, "forbidden_origin"},
		{"form content type", e, path, submitBody("b"), map[string]string{"Content-Type": "text/plain"}, 415, "unsupported_media_type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body := tt.env.do(http.MethodPost, tt.path, tt.body, tt.hdr)
			if code != tt.wantCode {
				t.Fatalf("status %d, want %d: %s", code, tt.wantCode, body)
			}
			var eb struct {
				Error struct{ Code, Message string } `json:"error"`
			}
			if err := json.Unmarshal(body, &eb); err != nil || eb.Error.Code != tt.wantErr {
				t.Fatalf("error code %q, want %q: %s", eb.Error.Code, tt.wantErr, body)
			}
		})
	}
	var n int
	if err := e.pool.QueryRow(context.Background(), `SELECT count(*) FROM submissions`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("rejected submissions were stored: %d", n)
	}
}

func TestNotFoundUsesErrorEnvelope(t *testing.T) {
	e := setup(t, false)
	for _, p := range []string{"/api/v1/courses/nope", "/api/v1/exercises/nope", "/api/v1/nope"} {
		code, body := e.do(http.MethodGet, p, "", nil)
		if code != 404 || !strings.Contains(string(body), `"code":"not_found"`) {
			t.Errorf("%s = %d %s", p, code, body)
		}
	}
}

func TestLessonImages(t *testing.T) {
	e := setup(t, false)
	var detail struct {
		Assets map[string]string `json:"assets"`
	}
	e.getJSON("/api/v1/exercises/"+exID, &detail)
	url := detail.Assets["assets/choice.svg"]
	if !strings.HasPrefix(url, "/api/v1/assets/") || !strings.HasSuffix(url, ".svg") {
		t.Fatalf("assets = %v", detail.Assets)
	}

	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	e.h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !bytes.HasPrefix(rec.Body.Bytes(), []byte("<svg")) {
		t.Fatalf("GET %s = %d: %s", url, rec.Code, rec.Body.String())
	}
	want := map[string]string{
		"Content-Type":            "image/svg+xml",
		"X-Content-Type-Options":  "nosniff",
		"Content-Security-Policy": "default-src 'none'; style-src 'unsafe-inline'; sandbox",
		"Cache-Control":           "public, max-age=31536000, immutable",
	}
	for k, v := range want {
		if got := rec.Header().Get(k); got != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}

	req = httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("If-None-Match", rec.Header().Get("ETag"))
	rec = httptest.NewRecorder()
	e.h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotModified {
		t.Fatalf("conditional GET = %d", rec.Code)
	}

	sha := strings.TrimSuffix(strings.TrimPrefix(url, "/api/v1/assets/"), ".svg")
	for _, p := range []string{
		"/api/v1/assets/" + strings.Repeat("0", 64) + ".svg", // unknown hash
		"/api/v1/assets/" + sha + ".png",                     // wrong type for these bytes
		"/api/v1/assets/" + sha + ".html",                    // not an image type
		"/api/v1/assets/../me",                               // not a hash
	} {
		if code, _ := e.do(http.MethodGet, p, "", nil); code != http.StatusNotFound {
			t.Errorf("GET %s = %d, want 404", p, code)
		}
	}
}
