package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// StdoutEqualsConfig configures stdout_equals: the entry file, run as a
// program, must print exactly Expected (trailing whitespace on each line and
// trailing blank lines are ignored).
//
//   - kind: stdout_equals
//     expected: "Hello, Ada!\n"
//     name: Prints the greeting   # default "Program output"
//     visible: true               # default true
type StdoutEqualsConfig struct {
	common
	Expected string `json:"expected"`
}

// StdoutRegexConfig configures stdout_regex: the program's output must match
// Pattern (Go RE2 syntax, unanchored; use ^...$ with (?m) for whole lines).
type StdoutRegexConfig struct {
	common
	Pattern string `json:"pattern"`
	Message string `json:"message"` // shown when it fails
}

type stdoutEquals struct{}
type stdoutRegex struct{}

func init() {
	Register(stdoutEquals{})
	Register(stdoutRegex{})
}

func (stdoutEquals) Kind() string { return "stdout_equals" }
func (stdoutRegex) Kind() string  { return "stdout_regex" }

func (stdoutEquals) Validate(raw json.RawMessage) error {
	var c StdoutEqualsConfig
	if err := decodeStrict(raw, &c); err != nil {
		return err
	}
	if c.Expected == "" {
		return errors.New("expected: required")
	}
	return c.validate()
}

func (stdoutRegex) Validate(raw json.RawMessage) error {
	var c StdoutRegexConfig
	if err := decodeStrict(raw, &c); err != nil {
		return err
	}
	if c.Pattern == "" {
		return errors.New("pattern: required")
	}
	if _, err := regexp.Compile(c.Pattern); err != nil {
		return fmt.Errorf("pattern: %w", err)
	}
	return c.validate()
}

// NormalizeOutput trims trailing whitespace from every line and drops trailing
// blank lines.
func NormalizeOutput(s string) string {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}

// programRow starts a row for a stdout check, failing it when the program did
// not run to completion. ok is false when the row is already decided.
func programRow(name string, c common, r *Report) (exercise.CheckResult, *Phase, bool) {
	if name == "" {
		name = "Program output"
	}
	row := exercise.CheckResult{Name: name, Visible: c.visible(true), Weight: c.Weight}
	switch {
	case r == nil:
		row.Feedback = "Did not run"
	case r.SyntaxError != nil:
		row.Feedback = fmt.Sprintf("%s (line %d)", r.SyntaxError.Message, r.SyntaxError.Line)
		row.Line = r.SyntaxError.Line
	case r.Program == nil:
		row.Feedback = "Did not run"
	case r.Program.TimedOut:
		row.Feedback = "Did not finish: time limit exceeded"
	case r.Program.ExitCode != 0:
		row.Feedback = "Your program stopped with an error: " + lastLine(r.Program.Stderr)
	default:
		return row, r.Program, true
	}
	return row, nil, false
}

func (stdoutEquals) Run(_ context.Context, in CheckInput) ([]exercise.CheckResult, error) {
	var c StdoutEqualsConfig
	if err := decodeStrict(in.Config, &c); err != nil {
		return nil, err
	}
	row, prog, ok := programRow(c.Name, c.common, in.Report)
	if ok {
		got, want := NormalizeOutput(prog.Stdout), NormalizeOutput(c.Expected)
		row.Passed = got == want
		if !row.Passed {
			row.Feedback = "Output differs"
			if row.Visible {
				row.Got, row.Expected = clip(got), clip(want)
			}
		}
	}
	return []exercise.CheckResult{row}, nil
}

func (stdoutRegex) Run(_ context.Context, in CheckInput) ([]exercise.CheckResult, error) {
	var c StdoutRegexConfig
	if err := decodeStrict(in.Config, &c); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(c.Pattern)
	if err != nil {
		return nil, fmt.Errorf("stdout_regex: %w", err)
	}
	row, prog, ok := programRow(c.Name, c.common, in.Report)
	if ok {
		row.Passed = re.MatchString(NormalizeOutput(prog.Stdout))
		if !row.Passed {
			row.Feedback = c.Message
			if row.Feedback == "" {
				row.Feedback = "Output does not look right"
			}
		}
	}
	return []exercise.CheckResult{row}, nil
}

func clip(s string) string {
	const n = 200
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
