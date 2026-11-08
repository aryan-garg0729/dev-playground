package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// IOCasesConfig configures io_cases: the program is run once per case with the
// case's stdin, and its output must equal the case's stdout (compared like
// stdout_equals). Each case is one result row. Hidden cases show no got/expected.
//
//   - kind: io_cases
//     name: Sample          # row name prefix, default "Case"; rows "Sample 1", ...
//     visible: true         # default for the cases, default true
//     cases:
//   - stdin: "3\n"
//     stdout: "1\n2\nFizz\n"
//     name: n = 3           # optional row name
//     visible: false        # optional, overrides the check's
//
// The runner reports the runs in Report.Cases in case order. An exercise has at
// most one io_cases check.
type IOCasesConfig struct {
	common
	Cases []IOCase `json:"cases"`
}

// IOCase is one stdin and the output it must produce.
type IOCase struct {
	Name    string `json:"name"`
	Stdin   string `json:"stdin"`
	Stdout  string `json:"stdout"`
	Visible *bool  `json:"visible"`
}

// MaxIOCases is the most cases one io_cases check may have.
const MaxIOCases = 50

type ioCases struct{}

func init() { Register(ioCases{}) }

func (ioCases) Kind() string { return "io_cases" }

func (ioCases) Validate(raw json.RawMessage) error {
	var c IOCasesConfig
	if err := decodeStrict(raw, &c); err != nil {
		return err
	}
	if len(c.Cases) == 0 || len(c.Cases) > MaxIOCases {
		return fmt.Errorf("cases: need 1-%d", MaxIOCases)
	}
	for i, k := range c.Cases {
		if k.Stdout == "" {
			return fmt.Errorf("cases[%d].stdout: required", i)
		}
	}
	return c.validate()
}

// IOCaseStdins returns each case's stdin, in order, and whether it is visible.
func IOCaseStdins(raw json.RawMessage) (stdins []string, visible []bool, err error) {
	var c IOCasesConfig
	if err := decodeStrict(raw, &c); err != nil {
		return nil, nil, err
	}
	for _, k := range c.Cases {
		stdins = append(stdins, k.Stdin)
		visible = append(visible, k.visible(c.common))
	}
	return stdins, visible, nil
}

func (k IOCase) visible(c common) bool {
	if k.Visible != nil {
		return *k.Visible
	}
	return c.visible(true)
}

func (ioCases) Run(_ context.Context, in CheckInput) ([]exercise.CheckResult, error) {
	var c IOCasesConfig
	if err := decodeStrict(in.Config, &c); err != nil {
		return nil, err
	}
	if in.Report != nil && in.Report.SyntaxError == nil && len(in.Report.Cases) > len(c.Cases) {
		return nil, errors.New("io_cases: report has more runs than cases")
	}
	prefix := c.Name
	if prefix == "" {
		prefix = "Case"
	}
	rows := make([]exercise.CheckResult, 0, len(c.Cases))
	for i, k := range c.Cases {
		name := k.Name
		if name == "" {
			name = fmt.Sprintf("%s %d", prefix, i+1)
		}
		row := exercise.CheckResult{Name: name, Visible: k.visible(c.common), Weight: c.Weight}
		r := in.Report
		switch {
		case r == nil:
			row.Feedback = "Did not run"
		case r.SyntaxError != nil:
			row.Feedback = fmt.Sprintf("%s (line %d)", r.SyntaxError.Message, r.SyntaxError.Line)
			row.Line = r.SyntaxError.Line
		case i >= len(r.Cases):
			row.Feedback = "Did not run"
		case r.Cases[i].TimedOut:
			row.Feedback = "Did not finish: time limit exceeded"
		case r.Cases[i].ExitCode != 0:
			row.Feedback = "Your program stopped with an error: " + lastLine(r.Cases[i].Stderr)
		default:
			got, want := NormalizeOutput(r.Cases[i].Stdout), NormalizeOutput(k.Stdout)
			row.Passed = got == want
			if !row.Passed {
				row.Feedback = "Output differs"
				if row.Visible {
					row.Got, row.Expected = clip(got), clip(want)
				}
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}
