package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// UnitTestsConfig configures unit_tests: run the test_* functions of File
// against the learner's code, one CheckResult per test.
//
//   - kind: unit_tests
//     file: tests_hidden.py   # required, in the exercise folder
//     visible: false          # default false; visible tests show got/expected
//     weight: 2               # per test, default 1
type UnitTestsConfig struct {
	common
	File string `json:"file"`
}

type unitTests struct{}

func init() { Register(unitTests{}) }

func (unitTests) Kind() string { return "unit_tests" }

func (unitTests) Validate(raw json.RawMessage) error {
	var c UnitTestsConfig
	if err := decodeStrict(raw, &c); err != nil {
		return err
	}
	if c.File == "" {
		return errors.New("file: required")
	}
	if path.Base(c.File) != c.File || !strings.HasSuffix(c.File, ".py") {
		return fmt.Errorf("file: %q must be a .py file directly in the exercise folder", c.File)
	}
	if c.Name != "" {
		return errors.New("name: not used by unit_tests (each test's docstring is its name)")
	}
	return c.validate()
}

func (unitTests) Run(_ context.Context, in CheckInput) ([]exercise.CheckResult, error) {
	var c UnitTestsConfig
	if err := decodeStrict(in.Config, &c); err != nil {
		return nil, err
	}
	visible := c.visible(false)
	tf := in.Report.File(c.File)
	if tf == nil {
		return nil, fmt.Errorf("unit_tests: no results for %s in the report", c.File)
	}
	var syntax *SyntaxError
	if in.Report != nil {
		syntax = in.Report.SyntaxError
	}
	timedOut := in.Report != nil && in.Report.TimedOut
	rows := make([]exercise.CheckResult, 0, len(tf.Cases))
	for i, tc := range tf.Cases {
		row := exercise.CheckResult{Passed: tc.Ran && tc.Passed, Visible: visible, Weight: c.Weight}
		switch {
		case tc.Label != "" && (visible || !looksLikeCode(tc.Label)):
			row.Name = tc.Label
		case visible:
			row.Name = strings.TrimPrefix(tc.Name, "test_")
		default:
			row.Name = fmt.Sprintf("Hidden test %d", i+1)
		}
		if !row.Passed {
			row.Feedback, row.Line = caseFeedback(tc, visible, syntax, timedOut)
			if visible && tc.Ran && tc.Error == "" {
				row.Got, row.Expected = tc.Got, tc.Expected
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// caseFeedback explains a failed test. Hidden tests never reveal values,
// assertion internals or exception messages, only what kind of thing went wrong.
func caseFeedback(tc TestCase, visible bool, syntax *SyntaxError, timedOut bool) (string, int) {
	switch {
	case syntax != nil:
		return fmt.Sprintf("%s (line %d)", syntax.Message, syntax.Line), syntax.Line
	case !tc.Ran && timedOut:
		return "Did not finish: time limit exceeded", 0
	case !tc.Ran:
		return "Did not run: your code stopped the tests early", 0
	case tc.Error != "":
		if visible {
			return tc.Error, tc.Line
		}
		return "Your code raised " + errorType(tc.Error), tc.Line
	case tc.Got != "" || tc.Expected != "":
		if tc.Message != "" {
			return tc.Message, tc.Line
		}
		if visible {
			return "", tc.Line // got/expected carry it
		}
		return "Wrong result", tc.Line
	case tc.Message != "":
		return tc.Message, tc.Line
	default:
		return "Assertion failed", tc.Line
	}
}

// looksLikeCode guards hidden test labels: authors should label hidden tests in
// words ("handles an empty list"), never with the call and its input.
func looksLikeCode(label string) bool {
	return strings.ContainsAny(label, "()[]{}=")
}
