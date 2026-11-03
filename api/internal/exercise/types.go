// Package exercise defines the exercise type plugin contract and the grading
// contracts (Submission, Result) shared by every type.
package exercise

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"

	"go.yaml.in/yaml/v3"
)

// Spec is a type-specific, validated exercise spec. It must round-trip through
// encoding/json: cq build stores json.Marshal(spec) in exercises.private_spec and
// the API rebuilds it with Type.LoadSpec.
type Spec any

// Exercise is a compiled exercise: the common manifest header plus its spec.
type Exercise struct {
	ID         string
	Type       string
	Version    int
	Title      string
	XP         int
	Difficulty int
	Tags       []string
	Free       bool
	Boss       bool
	Spec       Spec
}

// Type is one exercise type plugin ("code", "mcq", "sql").
type Type interface {
	Name() string
	// ParseSpec validates the manifest's spec block and loads any files it names from dir.
	ParseSpec(raw yaml.Node, dir fs.FS) (Spec, error)
	// LoadSpec rebuilds a Spec from the JSON stored in exercises.private_spec.
	LoadSpec(private json.RawMessage) (Spec, error)
	// PublicView returns what the browser may see. It must never include answers,
	// solutions, hidden tests or wrong/ files.
	PublicView(ex Exercise) (json.RawMessage, error)
	Grade(ctx context.Context, ex Exercise, sub Submission) (Result, error)
}

// Submission is what the browser POSTs to /exercises/{id}/submit.
type Submission struct {
	ExerciseID string          `json:"exercise_id"`
	Version    int             `json:"version"`
	Answer     json.RawMessage `json:"answer"`
}

// Result is the only authority on pass/fail. Game, tutor and analytics code
// depend on Result, never on the exercise type.
type Result struct {
	Status    string        `json:"status"` // pass | fail | error | timeout
	Score     float64       `json:"score"`  // 0..1, weighted by checks
	Checks    []CheckResult `json:"checks"`
	Stdout    string        `json:"stdout"`
	RuntimeMs int           `json:"runtime_ms"`
}

// CheckResult is one row of a Result.
type CheckResult struct {
	Name     string `json:"name"`
	Passed   bool   `json:"passed"`
	Visible  bool   `json:"visible"`
	Feedback string `json:"feedback,omitempty"` // hidden checks: message only, never test source
	// Got and Expected are short value reprs, set only on failed visible checks.
	Got      string `json:"got,omitempty"`
	Expected string `json:"expected,omitempty"`
	// Line is the learner's source line the failure points at (1-based), if known.
	Line int `json:"line,omitempty"`
	// Weight scales this row in Score; 0 counts as 1. Not sent to the browser.
	Weight float64 `json:"-"`
}

// Result statuses.
const (
	StatusPass    = "pass"
	StatusFail    = "fail"
	StatusError   = "error"
	StatusTimeout = "timeout"
)

// ErrInvalidAnswer means the answer does not match the type's answer shape.
// The API maps it to 400; it is not a failed attempt.
var ErrInvalidAnswer = errors.New("invalid answer")

// NewResult builds a Result from check rows: pass only when every check passes;
// score is the weighted passing fraction. No rows is an error.
func NewResult(checks []CheckResult) Result {
	r := Result{Status: StatusPass, Checks: checks}
	if len(checks) == 0 {
		r.Status = StatusError
		return r
	}
	var passed, total float64
	for _, c := range checks {
		w := c.Weight
		if w <= 0 {
			w = 1
		}
		total += w
		if c.Passed {
			passed += w
		} else {
			r.Status = StatusFail
		}
	}
	r.Score = passed / total
	return r
}

// Testable is implemented by types whose exercises carry reference answers.
// cq test grades the reference answer (must pass) and every wrong answer
// (each must fail). Specs come straight from ParseSpec, so fields that are not
// stored in private_spec (solution, wrong/) are still present.
type Testable interface {
	ReferenceAnswer(s Spec) (json.RawMessage, error)
	// WrongAnswers maps a name (e.g. "wrong/off_by_one.py") to an answer.
	WrongAnswers(s Spec) (map[string]json.RawMessage, error)
}

// Run modes: where the Run button executes an exercise (manifest `run:`).
const (
	RunBrowser = "browser" // sandbox iframe (Pyodide, sql.js)
	RunServer  = "server"  // POST /exercises/{id}/run on the grading worker, ungraded
	RunOff     = "off"     // no Run button
)

// ErrRunUnavailable means the exercise does not run on the server.
var ErrRunUnavailable = errors.New("run is not available on the server for this exercise")

// Table is a query result shown to the learner.
type Table struct {
	Columns   []string `json:"columns"`
	Rows      [][]any  `json:"rows"`
	Truncated bool     `json:"truncated"`
}

// RunOutput is an ungraded run: visible checks only, nothing stored.
type RunOutput struct {
	Stdout    string        `json:"stdout"`
	Stderr    string        `json:"stderr"`
	Checks    []CheckResult `json:"checks"`
	RuntimeMs int           `json:"runtime_ms"`
	TimedOut  bool          `json:"timed_out"`
	Table     *Table        `json:"table,omitempty"`
	Error     string        `json:"error,omitempty"`
}

// Runnable is implemented by types whose exercises can Run on the server
// (`run: server`). Run must use only public inputs: visible tests and data.
type Runnable interface {
	Run(ctx context.Context, ex Exercise, answer json.RawMessage) (RunOutput, error)
}

// ValidRunMode reports whether m is a known run mode.
func ValidRunMode(m string) bool { return m == RunBrowser || m == RunServer || m == RunOff }
