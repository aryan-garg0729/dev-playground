package httpapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/aryangarg/codequest/api/internal/runner"
)

// grant gives the user Bytes the way the economy does: ledger row plus balance.
func (e *env) grant(n int) {
	e.t.Helper()
	ctx := e.t.Context()
	if _, err := e.pool.Exec(ctx, `INSERT INTO user_stats (user_id) VALUES ($1) ON CONFLICT DO NOTHING`, e.user.UUID); err != nil {
		e.t.Fatal(err)
	}
	if _, err := e.pool.Exec(ctx, `INSERT INTO byte_events (user_id, amount, reason, ref_id) VALUES ($1, $2, 'test_grant', gen_random_uuid()::text)`, e.user.UUID, n); err != nil {
		e.t.Fatal(err)
	}
	if _, err := e.pool.Exec(ctx, `UPDATE user_stats SET bytes = bytes + $2 WHERE user_id = $1`, e.user.UUID, n); err != nil {
		e.t.Fatal(err)
	}
}

// balance returns the cached balance and checks it equals the ledger.
func (e *env) balance() int {
	e.t.Helper()
	var bytes, sum int
	if err := e.pool.QueryRow(e.t.Context(), `SELECT s.bytes, COALESCE((SELECT sum(amount) FROM byte_events b WHERE b.user_id = s.user_id), 0)
		FROM user_stats s WHERE s.user_id = $1`, e.user.UUID).Scan(&bytes, &sum); err != nil {
		e.t.Fatal(err)
	}
	if bytes != sum {
		e.t.Fatalf("balance %d != ledger sum %d", bytes, sum)
	}
	return bytes
}

func (e *env) owned(powerup string) int {
	e.t.Helper()
	var n int
	_ = e.pool.QueryRow(e.t.Context(), `SELECT COALESCE((SELECT count FROM inventory WHERE user_id = $1 AND powerup = $2), 0)`, e.user.UUID, powerup).Scan(&n)
	return n
}

func (e *env) post(path, body string, want int) []byte {
	e.t.Helper()
	code, out := e.do(http.MethodPost, path, body, nil)
	if code != want {
		e.t.Fatalf("POST %s %s = %d (want %d): %s", path, body, code, want, out)
	}
	return out
}

func req(id string, extra ...string) string {
	return `{"request_id":"` + id + `"` + strings.Join(extra, "") + `}`
}

// Criteria: the balance always equals the ledger; too few Bytes is refused
// and changes nothing; a replayed buy or use changes state once.
func TestBytesLedgerAndIdempotency(t *testing.T) {
	e := setup(t, true)
	e.grant(100)

	// Too poor for Focus Boost (200): refused, nothing changes.
	body := e.post("/api/v1/shop/focus_boost/buy", req("r1"), http.StatusConflict)
	if !strings.Contains(string(body), "insufficient_bytes") || e.balance() != 100 || e.owned("focus_boost") != 0 {
		t.Fatalf("refused buy changed state: %s balance %d", body, e.balance())
	}

	e.grant(1000)
	first := e.post("/api/v1/shop/focus_boost/buy", req("r2"), http.StatusOK)
	again := e.post("/api/v1/shop/focus_boost/buy", req("r2"), http.StatusOK)
	if string(first) != string(again) || e.balance() != 900 || e.owned("focus_boost") != 1 {
		t.Fatalf("replayed buy: %s / %s, balance %d, owned %d", first, again, e.balance(), e.owned("focus_boost"))
	}
	// The same id cannot be reused for another action.
	e.post("/api/v1/shop/peek/buy", req("r2"), http.StatusConflict)

	// Using: once per request id.
	e.post("/api/v1/powerups/focus_boost/use", req("u1"), http.StatusOK)
	e.post("/api/v1/powerups/focus_boost/use", req("u1"), http.StatusOK)
	if e.owned("focus_boost") != 0 {
		t.Fatalf("replayed use spent %d", 1-e.owned("focus_boost"))
	}
	// One Focus Boost a day, even with another in hand.
	e.post("/api/v1/shop/focus_boost/buy", req("r3"), http.StatusOK)
	e.post("/api/v1/powerups/focus_boost/use", req("u2"), http.StatusConflict)

	// Streak Freeze: hold at most 2; the third is refused and costs nothing.
	e.post("/api/v1/shop/streak_freeze/buy", req("f1"), http.StatusOK)
	e.post("/api/v1/shop/streak_freeze/buy", req("f2"), http.StatusOK)
	before := e.balance()
	e.post("/api/v1/shop/streak_freeze/buy", req("f3"), http.StatusConflict)
	if e.balance() != before {
		t.Fatalf("refused third freeze cost Bytes: %d -> %d", before, e.balance())
	}
	// Not for sale: Pro and coming-soon items, and Streak Repair (used, not bought).
	for _, p := range []string{"holiday_mode", "extra_hint", "streak_repair", "nope"} {
		e.post("/api/v1/shop/"+p+"/buy", req("x-"+p), http.StatusBadRequest)
	}
	var shop struct {
		Balance int
		Items   []struct{ ID, State string }
	}
	e.getJSON("/api/v1/shop", &shop)
	states := map[string]string{}
	for _, it := range shop.Items {
		states[it.ID] = it.State
	}
	if shop.Balance != e.balance() || states["streak_freeze"] != "max_held" || states["holiday_mode"] != "pro" || states["extra_hint"] != "coming_soon" || states["streak_repair"] != "streak_safe" {
		t.Fatalf("shop = %+v", shop)
	}
}

// hiddenReport is a harness report for demo.basics.sum: visible tests pass,
// both hidden tests fail; the first declares an input, the second does not.
func hiddenReport() string {
	c := func(name, input string, pass bool) map[string]any {
		return map[string]any{"name": name, "label": "", "input": input, "ran": true, "passed": pass, "message": "",
			"error": "", "got": "0", "expected": "-1", "line": 3}
	}
	r := map[string]any{
		"version": 1, "syntax_error": nil, "forbidden": []any{}, "program": nil, "timed_out": false, "runtime_ms": 5,
		"tests": map[string]any{
			"stdout": "", "stderr": "", "exit_code": 0, "timed_out": false, "runtime_ms": 5, "completed": true,
			"files": []any{
				map[string]any{"file": "tests_visible.py", "cases": []any{c("test_one", "", true), c("test_three", "", true)}},
				map[string]any{"file": "tests_hidden.py", "cases": []any{
					c("test_negative", "sum_list([-1000, 999]) == -1", false), c("test_empty", "", false),
				}},
			},
		},
	}
	b, _ := json.Marshal(r)
	return string(b) + "\n"
}

// Criterion: Peek returns only the input of a failed hidden check, never the
// expected output or test source; the input never appears anywhere else.
func TestPeekRevealsOnlyTheInput(t *testing.T) {
	e := setupWith(t, true, opts{runner: &gatedRunner{out: runner.Output{Stdout: hiddenReport()}}})
	e.grant(1000)
	sub := e.post("/api/v1/exercises/"+codeID+"/submit", codeBody("def sum_list(n):\n    return 0\n"), http.StatusOK)
	for _, leak := range []string{"sum_list([-1000", "input", "hidden_inputs"} {
		if strings.Contains(string(sub), leak) {
			t.Fatalf("submit response leaks %q: %s", leak, sub)
		}
	}
	if _, body := e.do(http.MethodGet, "/api/v1/exercises/"+codeID, "", nil); strings.Contains(string(body), "sum_list([-1000") {
		t.Fatal("exercise view leaks the hidden input")
	}

	// Nothing held: refused.
	e.post("/api/v1/powerups/peek/use", req("p0", `,"exercise_id":"`+codeID+`"`), http.StatusConflict)
	e.post("/api/v1/shop/peek/buy", req("b1"), http.StatusOK)
	e.post("/api/v1/shop/peek/buy", req("b2"), http.StatusOK)
	body := e.post("/api/v1/powerups/peek/use", req("p1", `,"exercise_id":"`+codeID+`"`), http.StatusOK)
	var out struct {
		Peek map[string]any
	}
	_ = json.Unmarshal(body, &out)
	if out.Peek["input"] != "sum_list([-1000, 999])" || out.Peek["name"] != "Hidden test 1" || len(out.Peek) != 3 {
		t.Fatalf("peek = %s", body)
	}
	for _, leak := range []string{"expected", "got", "SERVER_ONLY_SENTINEL", "assert", "== -1"} {
		if strings.Contains(string(body), leak) {
			t.Fatalf("peek leaks %q: %s", leak, body)
		}
	}
	// The other failed hidden check declares no input: nothing to show, nothing spent.
	e.post("/api/v1/powerups/peek/use", req("p2", `,"exercise_id":"`+codeID+`"`), http.StatusConflict)
	if e.owned("peek") != 1 {
		t.Fatalf("a refused peek was spent: %d left", e.owned("peek"))
	}
	var peeked bool
	_ = e.pool.QueryRow(t.Context(), `SELECT peeked FROM progress WHERE user_id = $1 AND exercise_id = $2`, e.user.UUID, codeID).Scan(&peeked)
	if !peeked {
		t.Fatal("peek did not remove the first-try bonus")
	}
}

// Criterion: no power-up can change a Result. A user with every power-up
// active gets exactly the Result of a user with none.
func TestNoPowerUpChangesResult(t *testing.T) {
	fail := &gatedRunner{out: runner.Output{Stdout: report(false)}}
	e := setupWith(t, true, opts{runner: fail})
	plain := *e
	plain.user, plain.cookie = e.signIn("plain@example.com", "Plain")

	e.grant(5000)
	// Arm everything that can be armed.
	e.post("/api/v1/exercises/"+exID+"/submit", submitBody("a"), http.StatusOK) // one failed submit
	e.post("/api/v1/exercises/"+codeID+"/submit", codeBody("x = 1\n"), http.StatusOK)
	for i, p := range []string{"focus_boost", "second_wind", "second_wind", "peek", "league_shield", "streak_freeze"} {
		e.post("/api/v1/shop/"+p+"/buy", req(fmt.Sprintf("buy%d", i)), http.StatusOK)
	}
	e.post("/api/v1/powerups/focus_boost/use", req("u1"), http.StatusOK)
	e.post("/api/v1/powerups/league_shield/use", req("u2"), http.StatusOK)
	e.post("/api/v1/powerups/second_wind/use", req("u3", `,"exercise_id":"`+exID+`"`), http.StatusOK)
	e.post("/api/v1/powerups/double_or_nothing/use", req("u4", `,"stake":200`), http.StatusOK)

	type graded struct {
		Result json.RawMessage `json:"result"`
	}
	for _, c := range []struct{ path, body string }{
		{"/api/v1/exercises/" + exID + "/submit", submitBody("a")},
		{"/api/v1/exercises/" + codeID + "/submit", codeBody("def sum_list(n):\n    return 0\n")},
	} {
		var withAll, without graded
		_ = json.Unmarshal(e.post(c.path, c.body, http.StatusOK), &withAll)
		_ = json.Unmarshal(plain.post(c.path, c.body, http.StatusOK), &without)
		var r struct{ Status string }
		_ = json.Unmarshal(withAll.Result, &r)
		if r.Status != "fail" || string(withAll.Result) != string(without.Result) {
			t.Fatalf("%s: with power-ups %s, without %s", c.path, withAll.Result, without.Result)
		}
	}
}

// Criterion: a boss shows HP = its checks, each passing check is a hit, and
// defeating it unlocks the next chapter (with its reward). The fixture's boss
// is demo.basics.sum (5 checks); chapter demo.more sits behind it.
func TestBossLocksTheNextChapter(t *testing.T) {
	g := &gatedRunner{out: runner.Output{Stdout: report(false)}}
	e := setupWith(t, true, opts{runner: g})
	after := "/api/v1/exercises/demo.more.after"
	if code, body := e.do(http.MethodGet, after, "", nil); code != http.StatusForbidden || !strings.Contains(string(body), "chapter_locked") {
		t.Fatalf("locked chapter = %d %s", code, body)
	}
	e.post(after+"/submit", `{"exercise_id":"demo.more.after","version":1,"answer":{"choices":["b"]}}`, http.StatusForbidden)
	var course struct {
		Chapters []struct {
			ID     string
			Locked bool
		}
	}
	e.getJSON("/api/v1/courses/demo", &course)
	if len(course.Chapters) != 2 || course.Chapters[0].Locked || !course.Chapters[1].Locked {
		t.Fatalf("course chapters = %+v", course.Chapters)
	}
	var view struct {
		PublicView struct {
			CheckCounts struct{ Visible, Hidden int } `json:"check_counts"`
		} `json:"public_view"`
	}
	e.getJSON("/api/v1/exercises/"+codeID, &view)
	// 2 visible tests + forbidden_calls, and 2 hidden tests: 5 HP.
	if view.PublicView.CheckCounts.Visible != 3 || view.PublicView.CheckCounts.Hidden != 2 {
		t.Fatalf("boss HP (check counts) = %+v", view.PublicView.CheckCounts)
	}

	attack := func(stdout string) map[string]any {
		g.out = runner.Output{Stdout: stdout}
		var out struct {
			Outcome struct {
				Rewards struct {
					Boss map[string]any
				}
			}
		}
		_ = json.Unmarshal(e.post("/api/v1/exercises/"+codeID+"/submit", codeBody("def sum_list(n):\n    return 0\n"), http.StatusOK), &out)
		return out.Outcome.Rewards.Boss
	}
	for _, step := range []struct {
		stdout string
		hp     float64
	}{{report(false), 4}, {hiddenReport(), 2}} { // forbidden_calls always hits; then the visible tests too
		if b := attack(step.stdout); b["max_hp"] != 5.0 || b["hp"] != step.hp || b["defeated"] != false {
			t.Fatalf("attack: %v, want hp %v", b, step.hp)
		}
	}
	if code, _ := e.do(http.MethodGet, after, "", nil); code != http.StatusForbidden {
		t.Fatal("chapter unlocked before the boss fell")
	}
	b := attack(report(true))
	chest, _ := b["chest"].(map[string]any)
	if b["hp"] != 0.0 || b["defeated"] != true || b["unlocked_chapter"] != "demo.more" || chest["bytes"] != 300.0 {
		t.Fatalf("defeat: %v", b)
	}
	if code, _ := e.do(http.MethodGet, after, "", nil); code != http.StatusOK {
		t.Fatalf("after defeat = %d", code)
	}
	if bal := e.balance(); bal < 300 {
		t.Fatalf("boss chest not paid: %d", bal)
	}
	before := e.balance() // defeating it again pays nothing more
	attack(report(true))
	if e.balance() != before {
		t.Fatalf("second defeat paid %d", e.balance()-before)
	}
}

// Quests fill on passes; the chest opens once, when all three daily quests are done.
func TestQuestsAndDailyChest(t *testing.T) {
	e := setupWith(t, true, opts{rand: func(int) int { return 0 }})
	var q struct {
		Daily []struct {
			ID       string
			Event    string
			Target   int
			Progress int
			Done     bool
		}
		Chest struct{ Ready bool }
	}
	e.getJSON("/api/v1/quests", &q)
	if len(q.Daily) != 3 || q.Chest.Ready {
		t.Fatalf("quests = %+v", q)
	}
	e.post("/api/v1/quests/chest/open", req("c0"), http.StatusConflict) // not ready
	// Force today's quests done (the per-event counting is unit-tested), then open.
	for _, d := range q.Daily {
		if _, err := e.pool.Exec(t.Context(), `INSERT INTO user_quests (user_id, quest_id, period, progress, done_at)
			VALUES ($1, $2, (now() AT TIME ZONE 'UTC')::date, $3, now())`, e.user.UUID, d.ID, d.Target); err != nil {
			t.Fatal(err)
		}
	}
	first := e.post("/api/v1/quests/chest/open", req("c1"), http.StatusOK)
	second := e.post("/api/v1/quests/chest/open", req("c2"), http.StatusOK) // another id: same chest, no new payout
	if e.balance() != 40 || !strings.Contains(string(first), `"bytes":40`) || string(first) != string(second) || e.owned("focus_boost") != 1 {
		t.Fatalf("chest: %s / %s, balance %d", first, second, e.balance())
	}
}

// Skeleton Key: sold only after 5 bosses, held one at a time; it opens the
// first locked chapter (demo.more, behind demo.basics.sum) without passing
// that boss, and only once per course.
func TestSkeletonKey(t *testing.T) {
	e := setup(t, true)
	ctx := context.Background()
	e.grant(13000)
	if body := e.post("/api/v1/shop/skeleton_key/buy", req("k1"), http.StatusConflict); !strings.Contains(string(body), "5 bosses") || e.balance() != 13000 {
		t.Fatalf("key before 5 bosses: %s, balance %d", body, e.balance())
	}
	if _, err := e.pool.Exec(ctx, `INSERT INTO counters (user_id, key, value) VALUES ($1, 'bosses', 5)`, e.user.UUID); err != nil {
		t.Fatal(err)
	}
	e.post("/api/v1/shop/skeleton_key/buy", req("k2"), http.StatusOK)
	if e.balance() != 7000 || e.owned("skeleton_key") != 1 {
		t.Fatalf("after buying: balance %d, owned %d", e.balance(), e.owned("skeleton_key"))
	}
	if body := e.post("/api/v1/shop/skeleton_key/buy", req("k3"), http.StatusConflict); !strings.Contains(string(body), "already hold") || e.balance() != 7000 {
		t.Fatalf("second key held: %s, balance %d", body, e.balance())
	}

	// An open chapter and an unknown one are refused and keep the key.
	e.post("/api/v1/powerups/skeleton_key/use", req("u1", `,"chapter_id":"demo.basics"`), http.StatusConflict)
	e.post("/api/v1/powerups/skeleton_key/use", req("u2", `,"chapter_id":"nope"`), http.StatusBadRequest)
	if e.owned("skeleton_key") != 1 {
		t.Fatal("a refused use spent the key")
	}

	after := "/api/v1/exercises/demo.more.after"
	if code, _ := e.do(http.MethodGet, after, "", nil); code != http.StatusForbidden {
		t.Fatalf("demo.more before the key = %d", code)
	}
	var used struct {
		Remaining       int    `json:"remaining"`
		UnlockedChapter string `json:"unlocked_chapter"`
	}
	_ = json.Unmarshal(e.post("/api/v1/powerups/skeleton_key/use", req("u3", `,"chapter_id":"demo.more"`), http.StatusOK), &used)
	if used.UnlockedChapter != "demo.more" || used.Remaining != 0 {
		t.Fatalf("use = %+v", used)
	}
	if code, body := e.do(http.MethodGet, after, "", nil); code != http.StatusOK {
		t.Fatalf("demo.more after the key = %d %s", code, body)
	}
	var course struct{ Chapters []struct{ Locked bool } }
	e.getJSON("/api/v1/courses/demo", &course)
	if course.Chapters[1].Locked {
		t.Fatal("course map still shows demo.more locked")
	}
	// The boss stays undefeated: not passed, no boss chest.
	var passed, chests int
	if err := e.pool.QueryRow(ctx, `SELECT
		(SELECT count(*) FROM progress WHERE user_id = $1 AND exercise_id = 'demo.basics.sum' AND status = 'passed'),
		(SELECT count(*) FROM chests WHERE user_id = $1)`, e.user.UUID).Scan(&passed, &chests); err != nil {
		t.Fatal(err)
	}
	if passed != 0 || chests != 0 {
		t.Fatalf("key passed the boss (%d) or opened a chest (%d)", passed, chests)
	}

	// Once per course, even with another key in hand.
	e.post("/api/v1/shop/skeleton_key/buy", req("k4"), http.StatusOK)
	if body := e.post("/api/v1/powerups/skeleton_key/use", req("u4", `,"chapter_id":"demo.more"`), http.StatusConflict); !strings.Contains(string(body), "this course") {
		t.Fatalf("second key in a course: %s", body)
	}
}
