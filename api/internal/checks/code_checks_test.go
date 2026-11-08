package checks

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

func run(t *testing.T, kind, cfg string, r *Report) []exercise.CheckResult {
	t.Helper()
	if err := Validate(kind, json.RawMessage(cfg)); err != nil {
		t.Fatalf("validate %s: %v", cfg, err)
	}
	rows, err := Run(context.Background(), kind, CheckInput{Config: json.RawMessage(cfg), Report: r})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	return rows
}

func testsReport(file string, cases ...TestCase) *Report {
	return &Report{Version: 1, Tests: &TestPhase{Completed: true, Files: []TestFile{{File: file, Cases: cases}}}}
}

func TestUnitTestsVisible(t *testing.T) {
	r := testsReport("tests_visible.py",
		TestCase{Name: "test_one", Label: "sum_list([5])", Ran: true, Passed: true},
		TestCase{Name: "test_three", Label: "sum_list([1, 2, 3])", Ran: true, Got: "1", Expected: "6"},
		TestCase{Name: "test_boom", Ran: true, Error: "TypeError: secret detail", Line: 4},
	)
	rows := run(t, "unit_tests", `{"file":"tests_visible.py","visible":true}`, r)
	if len(rows) != 3 {
		t.Fatalf("got %d rows", len(rows))
	}
	if !rows[0].Passed || rows[0].Name != "sum_list([5])" || !rows[0].Visible {
		t.Errorf("row 0 = %+v", rows[0])
	}
	if rows[1].Passed || rows[1].Got != "1" || rows[1].Expected != "6" {
		t.Errorf("row 1 = %+v, want got 1 expected 6", rows[1])
	}
	if rows[2].Name != "boom" || rows[2].Feedback != "TypeError: secret detail" || rows[2].Line != 4 {
		t.Errorf("row 2 = %+v", rows[2])
	}
}

func TestUnitTestsHiddenRevealsNothing(t *testing.T) {
	r := testsReport("tests_hidden.py",
		TestCase{Name: "test_a", Label: "sum_list([7, 8])", Ran: true, Got: "7", Expected: "15"},
		TestCase{Name: "test_b", Label: "handles negatives", Ran: true, Error: "KeyError: 'secret-input'"},
		TestCase{Name: "test_c", Ran: true, Message: "authored hint"},
		TestCase{Name: "test_d", Ran: true},
		TestCase{Name: "test_e", Ran: true, Passed: true},
	)
	rows := run(t, "unit_tests", `{"file":"tests_hidden.py","weight":2}`, r)
	want := []struct{ name, feedback string }{
		{"Hidden test 1", "Wrong result"}, // label looks like code: replaced
		{"handles negatives", "Your code raised KeyError"},
		{"Hidden test 3", "authored hint"},
		{"Hidden test 4", "Assertion failed"},
		{"Hidden test 5", ""},
	}
	for i, w := range want {
		row := rows[i]
		if row.Name != w.name || row.Feedback != w.feedback || row.Visible {
			t.Errorf("row %d = %+v, want name %q feedback %q", i, row, w.name, w.feedback)
		}
		if row.Got != "" || row.Expected != "" {
			t.Errorf("row %d leaks values: %+v", i, row)
		}
		if row.Weight != 2 {
			t.Errorf("row %d weight %v", i, row.Weight)
		}
	}
	for _, row := range rows {
		if strings.Contains(row.Name+row.Feedback, "secret") || strings.Contains(row.Name, "[7, 8]") {
			t.Errorf("hidden row leaks: %+v", row)
		}
	}
}

func TestUnitTestsNotRun(t *testing.T) {
	r := testsReport("t.py", TestCase{Name: "test_x"})
	r.TimedOut = true
	if rows := run(t, "unit_tests", `{"file":"t.py"}`, r); rows[0].Passed || !strings.Contains(rows[0].Feedback, "time limit") {
		t.Errorf("timed out: %+v", rows[0])
	}
	r.SyntaxError = &SyntaxError{File: "main.py", Line: 3, Message: "SyntaxError: invalid syntax"}
	if rows := run(t, "unit_tests", `{"file":"t.py"}`, r); rows[0].Line != 3 || !strings.Contains(rows[0].Feedback, "SyntaxError") {
		t.Errorf("syntax error: %+v", rows[0])
	}
}

func TestStdoutEquals(t *testing.T) {
	cfg := `{"expected":"a\nb\n"}`
	tests := []struct {
		name   string
		prog   *Phase
		passed bool
		fb     string
	}{
		{"exact", &Phase{Stdout: "a\nb\n"}, true, ""},
		{"trailing whitespace ignored", &Phase{Stdout: "a  \nb\n\n\n"}, true, ""},
		{"different", &Phase{Stdout: "a\nc\n"}, false, "Output differs"},
		{"crashed", &Phase{Stdout: "a\n", ExitCode: 1, Stderr: "Traceback...\nValueError: x"}, false, "Your program stopped with an error: ValueError: x"},
		{"timed out", &Phase{TimedOut: true}, false, "Did not finish: time limit exceeded"},
		{"not run", nil, false, "Did not run"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := run(t, "stdout_equals", cfg, &Report{Version: 1, Program: tt.prog})
			if rows[0].Passed != tt.passed || rows[0].Feedback != tt.fb {
				t.Fatalf("got %+v", rows[0])
			}
		})
	}
	rows := run(t, "stdout_equals", cfg, &Report{Version: 1, Program: &Phase{Stdout: "a\nc"}})
	if rows[0].Got != "a\nc" || rows[0].Expected != "a\nb" {
		t.Errorf("visible stdout check should show got/expected: %+v", rows[0])
	}
	rows = run(t, "stdout_equals", `{"expected":"x","visible":false}`, &Report{Version: 1, Program: &Phase{Stdout: "y"}})
	if rows[0].Got != "" || rows[0].Expected != "" {
		t.Errorf("hidden stdout check leaks: %+v", rows[0])
	}
}

func TestStdoutRegex(t *testing.T) {
	cfg := `{"pattern":"(?m)^total: \\d+$","message":"print the total"}`
	if rows := run(t, "stdout_regex", cfg, &Report{Version: 1, Program: &Phase{Stdout: "hi\ntotal: 42\n"}}); !rows[0].Passed {
		t.Errorf("match: %+v", rows[0])
	}
	if rows := run(t, "stdout_regex", cfg, &Report{Version: 1, Program: &Phase{Stdout: "total: many"}}); rows[0].Passed || rows[0].Feedback != "print the total" {
		t.Errorf("no match: %+v", rows[0])
	}
}

func TestForbiddenCalls(t *testing.T) {
	cfg := `{"names":["sum"],"message":"Use a loop"}`
	if rows := run(t, "forbidden_calls", cfg, &Report{Version: 1}); !rows[0].Passed || rows[0].Name != "Does not call sum()" {
		t.Errorf("clean: %+v", rows[0])
	}
	r := &Report{Version: 1, Forbidden: []Forbidden{{Name: "len", Line: 1}, {Name: "sum", File: "main.py", Line: 2}}}
	rows := run(t, "forbidden_calls", cfg, r)
	if rows[0].Passed || rows[0].Feedback != "Use a loop (line 2)" || rows[0].Line != 2 {
		t.Errorf("hit: %+v", rows[0])
	}
}

func TestValidateRejects(t *testing.T) {
	tests := []struct{ kind, cfg string }{
		{"unit_tests", `{}`},
		{"unit_tests", `{"file":"../x.py"}`},
		{"unit_tests", `{"file":"t.py","extra":1}`},
		{"unit_tests", `{"file":"t.py","weight":-1}`},
		{"stdout_equals", `{}`},
		{"stdout_regex", `{"pattern":"("}`},
		{"forbidden_calls", `{"names":[]}`},
		{"forbidden_calls", `{"names":["os.system"]}`},
		{"io_cases", `{}`},
		{"io_cases", `{"cases":[{"stdin":"1"}]}`},
		{"io_cases", `{"cases":[{"stdin":"1","stdout":"1","extra":1}]}`},
		{"exit_code", `{"expected":-1}`},
		{"exit_code", `{"expected":0,"extra":1}`},
		{"no_such_kind", `{}`},
	}
	for _, tt := range tests {
		if err := Validate(tt.kind, json.RawMessage(tt.cfg)); err == nil {
			t.Errorf("%s %s: want error", tt.kind, tt.cfg)
		}
	}
}

func TestParseReport(t *testing.T) {
	if _, err := ParseReport(`{"version":1}` + "\n"); err != nil {
		t.Errorf("valid: %v", err)
	}
	for _, bad := range []string{"", "not json", `{"version":2}`, "{\"version\":1}\n{\"version\":1}"} {
		if _, err := ParseReport(bad); err == nil {
			t.Errorf("%q: want error", bad)
		}
	}
}

func TestWeightedScore(t *testing.T) {
	r := exercise.NewResult([]exercise.CheckResult{{Passed: true}, {Passed: false, Weight: 3}})
	if r.Status != exercise.StatusFail || r.Score != 0.25 {
		t.Fatalf("got %s %v, want fail 0.25", r.Status, r.Score)
	}
}

func TestSameRows(t *testing.T) {
	a := [][]any{{"x", 1.0}, {"y", 2.0}}
	tests := []struct {
		name    string
		got     [][]any
		ordered bool
		want    bool
	}{
		{"same order", [][]any{{"x", 1.0}, {"y", 2.0}}, true, true},
		{"other order, unordered", [][]any{{"y", 2.0}, {"x", 1.0}}, false, true},
		{"other order, ordered", [][]any{{"y", 2.0}, {"x", 1.0}}, true, false},
		{"int vs float", [][]any{{"x", 1}, {"y", 2.0000000000001}}, true, true},
		{"string vs number", [][]any{{"x", "1"}, {"y", 2.0}}, true, false},
		{"missing row", [][]any{{"x", 1.0}}, false, false},
		{"duplicate instead", [][]any{{"x", 1.0}, {"x", 1.0}}, false, false},
	}
	for _, tt := range tests {
		if got := SameRows(tt.got, a, tt.ordered); got != tt.want {
			t.Errorf("%s: got %v", tt.name, got)
		}
	}
}

func TestSQLResultEquals(t *testing.T) {
	rep := &Report{Version: 1, SQL: &SQLPhase{Results: []SQLResult{
		{Learner: SQLTable{Columns: []string{"n"}, Rows: [][]any{{1.0}}}},
		{Learner: SQLTable{Columns: []string{"n"}, Rows: [][]any{{2.0}}}, Solution: &SQLTable{Columns: []string{"n"}, Rows: [][]any{{3.0}}}},
		{Learner: SQLTable{Error: "no such column: x"}, Solution: &SQLTable{Columns: []string{"n"}}},
	}}}
	if rows := run(t, "sql_result_equals", `{"expected":{"columns":["n"],"rows":[[1]]}}`, rep); !rows[0].Passed {
		t.Errorf("visible match: %+v", rows[0])
	}
	if rows := run(t, "sql_result_equals", `{"expected":{"columns":["n"],"rows":[[9]]}}`, rep); rows[0].Passed || rows[0].Got != "[[1]]" || rows[0].Expected != "[[9]]" {
		t.Errorf("visible mismatch: %+v", rows[0])
	}
	if rows := run(t, "sql_result_equals", `{"dataset":1,"visible":false}`, rep); rows[0].Passed || rows[0].Feedback != "Wrong result" || rows[0].Got != "" {
		t.Errorf("hidden mismatch must not show values: %+v", rows[0])
	}
	if rows := run(t, "sql_result_equals", `{"dataset":2,"visible":false}`, rep); rows[0].Feedback != "no such column: x" {
		t.Errorf("sql error: %+v", rows[0])
	}
}
