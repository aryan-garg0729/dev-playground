package checks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// ExitCodeConfig configures exit_code: the entry, run as a program, must exit
// with Expected (for example os.Exit(2) on bad input).
//
//   - kind: exit_code
//     expected: 2
//     name: Exits with 2   # default "Exit code"
type ExitCodeConfig struct {
	common
	Expected int `json:"expected"`
}

type exitCode struct{}

func init() { Register(exitCode{}) }

func (exitCode) Kind() string { return "exit_code" }

func (exitCode) Validate(raw json.RawMessage) error {
	var c ExitCodeConfig
	if err := decodeStrict(raw, &c); err != nil {
		return err
	}
	if c.Expected < 0 || c.Expected > 125 {
		return fmt.Errorf("expected: must be 0-125")
	}
	return c.validate()
}

func (exitCode) Run(_ context.Context, in CheckInput) ([]exercise.CheckResult, error) {
	var c ExitCodeConfig
	if err := decodeStrict(in.Config, &c); err != nil {
		return nil, err
	}
	name := c.Name
	if name == "" {
		name = "Exit code"
	}
	row := exercise.CheckResult{Name: name, Visible: c.visible(true), Weight: c.Weight}
	r := in.Report
	switch {
	case r == nil || r.Program == nil:
		row.Feedback = "Did not run"
	case r.SyntaxError != nil:
		row.Feedback = fmt.Sprintf("%s (line %d)", r.SyntaxError.Message, r.SyntaxError.Line)
		row.Line = r.SyntaxError.Line
	case r.Program.TimedOut:
		row.Feedback = "Did not finish: time limit exceeded"
	default:
		row.Passed = r.Program.ExitCode == c.Expected
		if !row.Passed {
			row.Feedback = fmt.Sprintf("Exited with %d, expected %d", r.Program.ExitCode, c.Expected)
		}
	}
	return []exercise.CheckResult{row}, nil
}
