package checks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Report is the harness's JSON report (runners/python/harness/harness.py).
// The harness is the only writer of the container's stdout, so learner output
// appears here only inside the Stdout fields.
type Report struct {
	Version     int          `json:"version"`
	SyntaxError *SyntaxError `json:"syntax_error"`
	Forbidden   []Forbidden  `json:"forbidden"`
	Program     *Phase       `json:"program"` // nil unless a stdout check asked for it
	Tests       *TestPhase   `json:"tests"`
	TimedOut    bool         `json:"timed_out"`
	RuntimeMs   int          `json:"runtime_ms"`
}

// SyntaxError is the first syntax error in the learner's files.
type SyntaxError struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// Forbidden is one call to a banned name found by the static scan.
type Forbidden struct {
	Name string `json:"name"`
	File string `json:"file"`
	Line int    `json:"line"`
}

// Phase is one child process run (the program, or the tests).
type Phase struct {
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ExitCode  int    `json:"exit_code"`
	TimedOut  bool   `json:"timed_out"`
	RuntimeMs int    `json:"runtime_ms"`
}

// TestPhase is the tests run plus per-test results.
type TestPhase struct {
	Phase
	Files     []TestFile `json:"files"`
	Completed bool       `json:"completed"`
}

// TestFile groups the cases of one test file, in source order.
type TestFile struct {
	File  string     `json:"file"`
	Cases []TestCase `json:"cases"`
}

// TestCase is one test function's outcome.
type TestCase struct {
	Name     string `json:"name"`
	Label    string `json:"label"` // first docstring line
	Ran      bool   `json:"ran"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message"` // assertion message
	Error    string `json:"error"`   // "TypeError: ..." when the code raised
	Got      string `json:"got"`
	Expected string `json:"expected"`
	Line     int    `json:"line"`
}

// ParseReport decodes the single JSON line the harness writes.
func ParseReport(stdout string) (*Report, error) {
	line := strings.TrimSpace(stdout)
	if line == "" {
		return nil, errors.New("harness wrote no report")
	}
	if strings.ContainsRune(line, '\n') {
		return nil, errors.New("harness report is not a single line")
	}
	dec := json.NewDecoder(bytes.NewReader([]byte(line)))
	var r Report
	if err := dec.Decode(&r); err != nil {
		return nil, fmt.Errorf("decode harness report: %w", err)
	}
	if r.Version != 1 {
		return nil, fmt.Errorf("harness report version %d, want 1", r.Version)
	}
	return &r, nil
}

// File returns the test results for one test file, or nil.
func (r *Report) File(name string) *TestFile {
	if r == nil || r.Tests == nil {
		return nil
	}
	for i := range r.Tests.Files {
		if r.Tests.Files[i].File == name {
			return &r.Tests.Files[i]
		}
	}
	return nil
}

// common is the fields every code check accepts.
type common struct {
	Name    string  `json:"name"`
	Visible *bool   `json:"visible"`
	Weight  float64 `json:"weight"`
}

func (c common) visible(def bool) bool {
	if c.Visible == nil {
		return def
	}
	return *c.Visible
}

func (c common) validate() error {
	if c.Weight < 0 {
		return errors.New("weight: must be ≥ 0")
	}
	return nil
}

// decodeStrict decodes a check config, rejecting unknown fields.
func decodeStrict(raw json.RawMessage, v any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("check config: %w", err)
	}
	return nil
}

func errorType(e string) string {
	if i := strings.IndexByte(e, ':'); i > 0 {
		return e[:i]
	}
	return e
}

func lastLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.LastIndexByte(s, '\n'); i >= 0 {
		s = s[i+1:]
	}
	return s
}
