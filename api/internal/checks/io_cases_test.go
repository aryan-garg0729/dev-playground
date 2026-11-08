package checks

import (
	"strings"
	"testing"
)

const ioCfg = `{"name":"Sample","cases":[
	{"stdin":"1\n","stdout":"one\n"},
	{"stdin":"2\n","stdout":"two\n","name":"n = 2"},
	{"stdin":"3\n","stdout":"three\n","visible":false}]}`

func TestIOCases(t *testing.T) {
	r := &Report{Version: 1, Cases: []Phase{
		{Stdout: "one  \n\n"},           // trailing space and blank lines are ignored
		{Stdout: "TWO\n"},               // wrong, visible: shows got/expected
		{Stdout: "four\n", ExitCode: 0}, // wrong, hidden: shows nothing
	}}
	rows := run(t, "io_cases", ioCfg, r)
	if len(rows) != 3 || rows[0].Name != "Sample 1" || rows[1].Name != "n = 2" || rows[2].Name != "Sample 3" {
		t.Fatalf("rows: %+v", rows)
	}
	if !rows[0].Passed || rows[1].Passed || rows[1].Got != "TWO" || rows[1].Expected != "two" {
		t.Fatalf("visible rows: %+v", rows[:2])
	}
	if rows[2].Passed || rows[2].Visible || rows[2].Got != "" || rows[2].Expected != "" || strings.Contains(rows[2].Feedback, "three") {
		t.Fatalf("hidden row leaks: %+v", rows[2])
	}
}

func TestIOCasesFailures(t *testing.T) {
	rows := run(t, "io_cases", ioCfg, &Report{Version: 1, Cases: []Phase{
		{TimedOut: true}, {ExitCode: 2, Stderr: "boom\npanic: index out of range\n"},
	}})
	if rows[0].Feedback != "Did not finish: time limit exceeded" ||
		rows[1].Feedback != "Your program stopped with an error: panic: index out of range" ||
		rows[2].Feedback != "Did not run" {
		t.Fatalf("%+v", rows)
	}
	rows = run(t, "io_cases", ioCfg, &Report{Version: 1, SyntaxError: &SyntaxError{Line: 4, Message: "undefined: x"}})
	for _, r := range rows {
		if r.Passed || r.Line != 4 {
			t.Fatalf("syntax error row: %+v", r)
		}
	}
	stdins, visible, err := IOCaseStdins([]byte(ioCfg))
	if err != nil || strings.Join(stdins, "") != "1\n2\n3\n" || visible[0] != true || visible[2] != false {
		t.Fatalf("stdins %q visible %v err %v", stdins, visible, err)
	}
}

func TestExitCode(t *testing.T) {
	for _, tt := range []struct {
		r    *Report
		pass bool
		fb   string
	}{
		{&Report{Version: 1, Program: &Phase{ExitCode: 2}}, true, ""},
		{&Report{Version: 1, Program: &Phase{ExitCode: 0}}, false, "Exited with 0, expected 2"},
		{&Report{Version: 1, Program: &Phase{TimedOut: true, ExitCode: 124}}, false, "Did not finish: time limit exceeded"},
		{&Report{Version: 1}, false, "Did not run"},
	} {
		rows := run(t, "exit_code", `{"expected":2}`, tt.r)
		if rows[0].Passed != tt.pass || rows[0].Feedback != tt.fb || rows[0].Name != "Exit code" {
			t.Errorf("%+v", rows[0])
		}
	}
}
