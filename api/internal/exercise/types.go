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

// NewResult builds a Result from check rows: pass only when every check passes,
// score is the passing fraction. Weighted scoring arrives with code checks (Phase 2).
func NewResult(checks []CheckResult) Result {
	r := Result{Status: StatusPass, Checks: checks}
	if len(checks) == 0 {
		r.Status = StatusError
		return r
	}
	passed := 0
	for _, c := range checks {
		if c.Passed {
			passed++
		} else {
			r.Status = StatusFail
		}
	}
	r.Score = float64(passed) / float64(len(checks))
	return r
}
