package httpapi_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func (e *env) count(sql string, args ...any) int {
	e.t.Helper()
	var n int
	if err := e.pool.QueryRow(context.Background(), sql, args...).Scan(&n); err != nil {
		e.t.Fatal(err)
	}
	return n
}

func TestAnalyticsEvents(t *testing.T) {
	e := setup(t, false)
	ok := `{"anon_id":"a1b2c3d4-e5","events":[{"name":"page_view","path":"/","props":{"ref":"x","n":1,"b":true}},{"name":"landing_cta","props":null}]}`
	if code, body := e.do(http.MethodPost, "/api/v1/events", ok, nil); code != http.StatusNoContent {
		t.Fatalf("events = %d: %s", code, body)
	}
	if n := e.count(`SELECT count(*) FROM analytics_events WHERE anon_id = 'a1b2c3d4-e5' AND user_id IS NULL`); n != 2 {
		t.Fatalf("stored %d events", n)
	}

	e.user, e.cookie = e.signIn("t@example.com", "T")
	if code, body := e.do(http.MethodPost, "/api/v1/events", `{"events":[{"name":"lesson_opened","props":{"id":"x"}}]}`, nil); code != http.StatusNoContent {
		t.Fatalf("signed-in events = %d: %s", code, body)
	}
	if n := e.count(`SELECT count(*) FROM analytics_events WHERE user_id = $1`, e.user.UUID); n != 1 {
		t.Fatalf("user events = %d", n)
	}

	many := `{"events":[` + strings.TrimSuffix(strings.Repeat(`{"name":"page_view"},`, 21), ",") + `]}`
	for name, body := range map[string]string{
		"unknown event":  `{"events":[{"name":"steal_cookies"}]}`,
		"nested props":   `{"events":[{"name":"page_view","props":{"a":{"b":1}}}]}`,
		"array props":    `{"events":[{"name":"page_view","props":[1]}]}`,
		"huge props":     `{"events":[{"name":"page_view","props":{"a":"` + strings.Repeat("x", 1100) + `"}}]}`,
		"too many":       many,
		"none":           `{"events":[]}`,
		"bad anon id":    `{"anon_id":"<script>","events":[{"name":"page_view"}]}`,
		"unknown fields": `{"events":[{"name":"page_view"}],"user_id":"x"}`,
	} {
		if code, b := e.do(http.MethodPost, "/api/v1/events", body, nil); code != http.StatusBadRequest {
			t.Errorf("%s = %d: %s", name, code, b)
		}
	}
	// Events from another site are refused like every state change.
	if code, _ := e.do(http.MethodPost, "/api/v1/events", ok, map[string]string{"Origin": "https://evil.example"}); code != http.StatusForbidden {
		t.Fatalf("cross-site events = %d", code)
	}
}

func TestClientErrorsGroupByFingerprint(t *testing.T) {
	e := setup(t, true)
	report := func(stack string) {
		body := `{"message":"TypeError: x is undefined","stack":"` + stack + `","url":"http://localhost:5173/learn/a/b","user_agent":"test","release":"dev"}`
		if code, b := e.do(http.MethodPost, "/api/v1/client-errors", body, nil); code != http.StatusNoContent {
			t.Fatalf("client error = %d: %s", code, b)
		}
	}
	// The same error from two builds (different bundle hashes and columns) is one group.
	report(`TypeError: x is undefined\n    at render (http://localhost:5173/assets/index-AbC12345.js:10:20)`)
	report(`TypeError: x is undefined\n    at render (http://localhost:5173/assets/index-ZzZ98765.js:11:5)`)
	report(`TypeError: x is undefined\n    at other (http://localhost:5173/assets/index-AbC12345.js:1:1)`)
	if n := e.count(`SELECT count(*) FROM client_errors`); n != 2 {
		t.Fatalf("groups = %d, want 2", n)
	}
	if n := e.count(`SELECT max(count) FROM client_errors`); n != 2 {
		t.Fatalf("max count = %d, want 2", n)
	}
	if n := e.count(`SELECT max(cardinality(users)) FROM client_errors`); n != 1 {
		t.Fatalf("users per group = %d, want 1", n)
	}
	if code, _ := e.do(http.MethodPost, "/api/v1/client-errors", `{"message":" "}`, nil); code != http.StatusBadRequest {
		t.Fatalf("empty message = %d", code)
	}
}
