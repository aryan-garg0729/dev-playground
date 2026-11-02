package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// ForbiddenCallsConfig configures forbidden_calls: the learner's code must not
// call any of Names (a static scan of the source: `sum(...)` and `x.sum(...)`
// both count). It is a teaching aid, not a sandbox.
//
//   - kind: forbidden_calls
//     names: [sum]
//     message: Try it with a loop instead of sum()
type ForbiddenCallsConfig struct {
	common
	Names   []string `json:"names"`
	Message string   `json:"message"`
}

var identRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type forbiddenCalls struct{}

func init() { Register(forbiddenCalls{}) }

func (forbiddenCalls) Kind() string { return "forbidden_calls" }

func (forbiddenCalls) Validate(raw json.RawMessage) error {
	var c ForbiddenCallsConfig
	if err := decodeStrict(raw, &c); err != nil {
		return err
	}
	if len(c.Names) == 0 {
		return errors.New("names: required")
	}
	for _, n := range c.Names {
		if !identRe.MatchString(n) {
			return fmt.Errorf("names: %q is not an identifier", n)
		}
	}
	return c.validate()
}

// ForbiddenNames returns the names a forbidden_calls config bans (for the harness job).
func ForbiddenNames(raw json.RawMessage) ([]string, error) {
	var c ForbiddenCallsConfig
	if err := decodeStrict(raw, &c); err != nil {
		return nil, err
	}
	return c.Names, nil
}

func (forbiddenCalls) Run(_ context.Context, in CheckInput) ([]exercise.CheckResult, error) {
	var c ForbiddenCallsConfig
	if err := decodeStrict(in.Config, &c); err != nil {
		return nil, err
	}
	name := c.Name
	if name == "" {
		calls := make([]string, len(c.Names))
		for i, n := range c.Names {
			calls[i] = n + "()"
		}
		name = "Does not call " + strings.Join(calls, ", ")
	}
	row := exercise.CheckResult{Name: name, Passed: true, Visible: true, Weight: c.Weight}
	if in.Report == nil {
		return nil, errors.New("forbidden_calls: no report")
	}
	if in.Report.SyntaxError != nil {
		row.Passed = false
		row.Feedback = "Could not check: your code has a syntax error"
		return []exercise.CheckResult{row}, nil
	}
	for _, f := range in.Report.Forbidden {
		if slices.Contains(c.Names, f.Name) {
			row.Passed = false
			msg := c.Message
			if msg == "" {
				msg = f.Name + "() is not allowed here"
			}
			row.Feedback = fmt.Sprintf("%s (line %d)", msg, f.Line)
			row.Line = f.Line
			break
		}
	}
	return []exercise.CheckResult{row}, nil
}
