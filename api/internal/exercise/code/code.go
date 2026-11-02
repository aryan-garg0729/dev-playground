// Package code is the "write a program" exercise type. Grading runs the
// learner's file plus every test file in a sandboxed runner (see
// internal/runner) and turns the harness report into check rows.
package code

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/checks"
	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/runner"
)

const (
	maxPromptBytes = 64 << 10
	maxFileBytes   = 50 << 10 // one learner or test file
	stdoutLimit    = 16 << 10 // Result.Stdout
)

// fileNameRe is a Python module file directly in the exercise folder.
var fileNameRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*\.py$`)

// Example is one "Example n:" block under the problem statement.
type Example struct {
	Input       string `json:"input" yaml:"input"`
	Output      string `json:"output" yaml:"output"`
	Explanation string `json:"explanation,omitempty" yaml:"explanation"`
}

// Limits are the per-exercise runner limits.
type Limits struct {
	TimeoutMs int `json:"timeout_ms" yaml:"timeout_ms"`
	MemoryMB  int `json:"memory_mb" yaml:"memory_mb"`
}

// Check is one manifest check: its kind plus the rest of its fields as JSON.
type Check struct {
	Kind   string          `json:"kind"`
	Config json.RawMessage `json:"config"`
}

// TestFile is a test file's name and source.
type TestFile struct {
	File   string `json:"file"`
	Source string `json:"source"`
}

// Spec is the compiled code spec, stored in private_spec. Solution and Wrong are
// read from the content folder by cq build/test and are never stored.
type Spec struct {
	Language    string     `json:"language"`
	Runner      string     `json:"runner"`
	Prompt      string     `json:"prompt"` // markdown
	Entry       string     `json:"entry"`
	Starter     string     `json:"starter"`
	Limits      Limits     `json:"limits"`
	Examples    []Example  `json:"examples"`
	Constraints []string   `json:"constraints"`
	Hints       []string   `json:"hints"`
	Checks      []Check    `json:"checks"`
	Tests       []TestFile `json:"tests"` // every unit_tests file, hidden ones included

	Solution string            `json:"-"`
	Wrong    map[string]string `json:"-"` // "wrong/off_by_one.py" -> source
}

type manifest struct {
	Language    string      `yaml:"language"`
	Runner      string      `yaml:"runner"`
	Prompt      string      `yaml:"prompt"`
	Entry       string      `yaml:"entry"`
	Starter     string      `yaml:"starter"`
	Solution    string      `yaml:"solution"`
	Limits      Limits      `yaml:"limits"`
	Examples    []Example   `yaml:"examples"`
	Constraints []string    `yaml:"constraints"`
	Hints       []string    `yaml:"hints"`
	Checks      []yaml.Node `yaml:"checks"`
}

// PublicView is what the browser receives: no solution, no hidden tests.
type PublicView struct {
	Language     string     `json:"language"`
	Runner       string     `json:"runner"`
	Prompt       string     `json:"prompt"`
	Entry        string     `json:"entry"`
	Starter      string     `json:"starter"`
	TimeoutMs    int        `json:"timeout_ms"`
	Examples     []Example  `json:"examples"`
	Constraints  []string   `json:"constraints"`
	Hints        []string   `json:"hints"`
	VisibleTests []TestFile `json:"visible_tests"`
}

// Answer is the learner's answer shape: exactly one file, the entry.
type Answer struct {
	Files map[string]string `json:"files"`
}

type codeType struct{}

func init() { exercise.Register(codeType{}) }

func (codeType) Name() string { return "code" }

func (codeType) ParseSpec(raw yaml.Node, dir fs.FS) (exercise.Spec, error) {
	var m manifest
	if err := exercise.DecodeStrict(raw, &m); err != nil {
		return nil, err
	}
	s := Spec{
		Language: m.Language, Runner: m.Runner, Entry: m.Entry, Limits: m.Limits,
		Examples: m.Examples, Constraints: m.Constraints, Hints: m.Hints,
	}
	if s.Limits.TimeoutMs == 0 {
		s.Limits.TimeoutMs = runner.DefaultTimeoutMs
	}
	if s.Limits.MemoryMB == 0 {
		s.Limits.MemoryMB = runner.DefaultMemoryMB
	}
	var err error
	if s.Prompt, err = readFile(dir, "prompt", m.Prompt, maxPromptBytes); err != nil {
		return nil, err
	}
	if s.Starter, err = readFile(dir, "starter", m.Starter, maxFileBytes); err != nil {
		return nil, err
	}
	if s.Solution, err = readFile(dir, "solution", m.Solution, maxFileBytes); err != nil {
		return nil, err
	}
	for i, node := range m.Checks {
		c, err := parseCheck(node)
		if err != nil {
			return nil, fmt.Errorf("checks[%d]: %w", i, err)
		}
		s.Checks = append(s.Checks, c)
		if c.Kind == "unit_tests" {
			var cfg checks.UnitTestsConfig
			_ = json.Unmarshal(c.Config, &cfg) // validated by parseCheck
			if slices.ContainsFunc(s.Tests, func(t TestFile) bool { return t.File == cfg.File }) {
				return nil, fmt.Errorf("checks[%d]: %s is used by two unit_tests checks", i, cfg.File)
			}
			src, err := readFile(dir, fmt.Sprintf("checks[%d].file", i), cfg.File, maxFileBytes)
			if err != nil {
				return nil, err
			}
			if !strings.Contains(src, "def test_") {
				return nil, fmt.Errorf("checks[%d]: %s has no test_* functions", i, cfg.File)
			}
			s.Tests = append(s.Tests, TestFile{File: cfg.File, Source: src})
		}
	}
	if s.Wrong, err = readWrong(dir); err != nil {
		return nil, err
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func parseCheck(node yaml.Node) (Check, error) {
	var fields map[string]any
	if err := node.Decode(&fields); err != nil {
		return Check{}, err
	}
	kind, _ := fields["kind"].(string)
	if kind == "" {
		return Check{}, errors.New("kind: required")
	}
	delete(fields, "kind")
	if fields == nil {
		fields = map[string]any{}
	}
	cfg, err := json.Marshal(fields)
	if err != nil {
		return Check{}, fmt.Errorf("%s: %w", kind, err)
	}
	if kind == "exact_choice" {
		return Check{}, errors.New("exact_choice is for mcq exercises")
	}
	if err := checks.Validate(kind, cfg); err != nil {
		return Check{}, fmt.Errorf("%s: %w", kind, err)
	}
	return Check{Kind: kind, Config: cfg}, nil
}

func readFile(dir fs.FS, field, name string, limit int) (string, error) {
	if name == "" {
		return "", fmt.Errorf("%s: required", field)
	}
	if !fs.ValidPath(name) {
		return "", fmt.Errorf("%s: %q must be a relative path inside the exercise folder", field, name)
	}
	b, err := fs.ReadFile(dir, name)
	if err != nil {
		return "", fmt.Errorf("%s: read %s: %w", field, name, err)
	}
	if len(b) > limit {
		return "", fmt.Errorf("%s: %s is larger than %d bytes", field, name, limit)
	}
	if !utf8.Valid(b) {
		return "", fmt.Errorf("%s: %s is not UTF-8", field, name)
	}
	return string(b), nil
}

// readWrong loads wrong/*.py: known-bad solutions that cq test expects to fail.
func readWrong(dir fs.FS) (map[string]string, error) {
	entries, err := fs.ReadDir(dir, "wrong")
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("wrong/: %w", err)
	}
	out := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		name := path.Join("wrong", e.Name())
		if path.Ext(name) != ".py" {
			return nil, fmt.Errorf("%s: wrong/ may hold only .py files", name)
		}
		src, err := readFile(dir, "wrong", name, maxFileBytes)
		if err != nil {
			return nil, err
		}
		out[name] = src
	}
	return out, nil
}

func (s Spec) validate() error {
	if s.Language != "python" {
		return fmt.Errorf("language: %q is not supported (python)", s.Language)
	}
	if _, ok := runner.Images[s.Runner]; !ok {
		return fmt.Errorf("runner: %q is not supported (docker-python, pyodide)", s.Runner)
	}
	if !fileNameRe.MatchString(s.Entry) || s.Entry == "cq.py" {
		return fmt.Errorf("entry: %q must be a Python file name such as main.py", s.Entry)
	}
	if s.Limits.TimeoutMs < 100 || s.Limits.TimeoutMs > runner.MaxTimeoutMs {
		return fmt.Errorf("limits.timeout_ms: must be 100-%d", runner.MaxTimeoutMs)
	}
	if s.Limits.MemoryMB < 64 || s.Limits.MemoryMB > runner.MaxMemoryMB {
		return fmt.Errorf("limits.memory_mb: must be 64-%d", runner.MaxMemoryMB)
	}
	if strings.TrimSpace(s.Prompt) == "" {
		return errors.New("prompt: file is empty")
	}
	if len(s.Checks) == 0 {
		return errors.New("checks: need at least one")
	}
	for i, ex := range s.Examples {
		if ex.Input == "" || ex.Output == "" {
			return fmt.Errorf("examples[%d]: input and output are required", i)
		}
	}
	for _, t := range s.Tests {
		if !fileNameRe.MatchString(t.File) || t.File == s.Entry || t.File == "cq.py" {
			return fmt.Errorf("checks: test file %q must be a Python file name other than the entry", t.File)
		}
	}
	for _, c := range s.Checks {
		if err := checks.Validate(c.Kind, c.Config); err != nil {
			return fmt.Errorf("checks: %s: %w", c.Kind, err)
		}
	}
	return nil
}

func (codeType) LoadSpec(private json.RawMessage) (exercise.Spec, error) {
	var s Spec
	if err := json.Unmarshal(private, &s); err != nil {
		return nil, fmt.Errorf("code: load spec: %w", err)
	}
	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("code: load spec: %w", err)
	}
	return s, nil
}

func specOf(ex exercise.Exercise) (Spec, error) {
	s, ok := ex.Spec.(Spec)
	if !ok {
		return Spec{}, fmt.Errorf("code: exercise %s has spec %T", ex.ID, ex.Spec)
	}
	return s, nil
}

// visibleTests returns the test files of unit_tests checks marked visible.
func (s Spec) visibleTests() []TestFile {
	out := []TestFile{}
	for _, c := range s.Checks {
		if c.Kind != "unit_tests" {
			continue
		}
		var cfg checks.UnitTestsConfig
		if json.Unmarshal(c.Config, &cfg) != nil || cfg.Visible == nil || !*cfg.Visible {
			continue
		}
		for _, t := range s.Tests {
			if t.File == cfg.File {
				out = append(out, t)
			}
		}
	}
	return out
}

func orEmpty[T any](v []T) []T {
	if v == nil {
		return []T{}
	}
	return v
}

func (codeType) PublicView(ex exercise.Exercise) (json.RawMessage, error) {
	s, err := specOf(ex)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(PublicView{
		Language: s.Language, Runner: s.Runner, Prompt: s.Prompt, Entry: s.Entry, Starter: s.Starter,
		TimeoutMs: s.Limits.TimeoutMs, Examples: orEmpty(s.Examples), Constraints: orEmpty(s.Constraints),
		Hints: orEmpty(s.Hints), VisibleTests: s.visibleTests(),
	})
	if err != nil {
		return nil, fmt.Errorf("code: public view: %w", err)
	}
	return b, nil
}

// harnessJob is the harness's stdin (runners/python/harness/harness.py).
type harnessJob struct {
	Files       map[string]string `json:"files"`
	Entry       string            `json:"entry"`
	Tests       []TestFile        `json:"tests"`
	RunEntry    bool              `json:"run_entry"`
	Forbidden   []string          `json:"forbidden"`
	TimeoutMs   int               `json:"timeout_ms"`
	OutputLimit int               `json:"output_limit"`
}

func (s Spec) parseAnswer(raw json.RawMessage) (string, error) {
	var a Answer
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&a); err != nil {
		return "", fmt.Errorf("%w: %w", exercise.ErrInvalidAnswer, err)
	}
	src, ok := a.Files[s.Entry]
	if !ok || len(a.Files) != 1 {
		return "", fmt.Errorf("%w: files must contain exactly %s", exercise.ErrInvalidAnswer, s.Entry)
	}
	if len(src) > maxFileBytes {
		return "", fmt.Errorf("%w: %s is larger than %d bytes", exercise.ErrInvalidAnswer, s.Entry, maxFileBytes)
	}
	if strings.ContainsRune(src, 0) {
		return "", fmt.Errorf("%w: %s contains a NUL byte", exercise.ErrInvalidAnswer, s.Entry)
	}
	return src, nil
}

func (codeType) Grade(ctx context.Context, ex exercise.Exercise, sub exercise.Submission) (exercise.Result, error) {
	s, err := specOf(ex)
	if err != nil {
		return exercise.Result{}, err
	}
	src, err := s.parseAnswer(sub.Answer)
	if err != nil {
		return exercise.Result{}, err
	}
	run, ok := runner.FromContext(ctx)
	if !ok {
		return exercise.Result{}, errors.New("code: no runner in context")
	}

	job := harnessJob{
		Files: map[string]string{s.Entry: src}, Entry: s.Entry, Tests: orEmpty(s.Tests),
		Forbidden: []string{}, TimeoutMs: s.Limits.TimeoutMs, OutputLimit: stdoutLimit,
	}
	for _, c := range s.Checks {
		switch c.Kind {
		case "stdout_equals", "stdout_regex":
			job.RunEntry = true
		case "forbidden_calls":
			names, err := checks.ForbiddenNames(c.Config)
			if err != nil {
				return exercise.Result{}, fmt.Errorf("code: %w", err)
			}
			job.Forbidden = append(job.Forbidden, names...)
		}
	}
	stdin, err := json.Marshal(job)
	if err != nil {
		return exercise.Result{}, fmt.Errorf("code: encode job: %w", err)
	}
	out, err := run.Run(ctx, runner.Job{
		Image: runner.Images[s.Runner], Stdin: stdin,
		Limits: runner.Limits{TimeoutMs: s.Limits.TimeoutMs, MemoryMB: s.Limits.MemoryMB},
	})
	if err != nil {
		return exercise.Result{}, fmt.Errorf("code: run: %w", err)
	}
	if out.TimedOut {
		return single(exercise.StatusTimeout, "Time limit",
			fmt.Sprintf("Your code ran longer than %d ms", s.Limits.TimeoutMs), out.RuntimeMs), nil
	}
	report, err := checks.ParseReport(out.Stdout)
	if err != nil {
		if out.ExitCode == 137 {
			return single(exercise.StatusError, "Memory limit",
				fmt.Sprintf("Your code used more than %d MB of memory", s.Limits.MemoryMB), out.RuntimeMs), nil
		}
		return exercise.Result{}, fmt.Errorf("code: %w (exit %d)", err, out.ExitCode)
	}

	var rows []exercise.CheckResult
	for _, c := range s.Checks {
		r, err := checks.Run(ctx, c.Kind, checks.CheckInput{Config: c.Config, Answer: sub.Answer, Report: report})
		if err != nil {
			return exercise.Result{}, fmt.Errorf("code: %w", err)
		}
		rows = append(rows, r...)
	}
	res := exercise.NewResult(rows)
	if report.TimedOut {
		res.Status = exercise.StatusTimeout
	}
	switch {
	case report.Program != nil:
		res.Stdout, res.RuntimeMs = report.Program.Stdout, report.Program.RuntimeMs
	case report.Tests != nil:
		res.Stdout, res.RuntimeMs = report.Tests.Stdout, report.Tests.RuntimeMs
	}
	if len(res.Stdout) > stdoutLimit {
		res.Stdout = res.Stdout[:stdoutLimit]
	}
	return res, nil
}

func single(status, name, feedback string, ms int) exercise.Result {
	return exercise.Result{
		Status: status, RuntimeMs: ms,
		Checks: []exercise.CheckResult{{Name: name, Passed: false, Visible: true, Feedback: feedback}},
	}
}

// ReferenceAnswer is the solution file submitted as the entry.
func (codeType) ReferenceAnswer(sp exercise.Spec) (json.RawMessage, error) {
	s, ok := sp.(Spec)
	if !ok {
		return nil, fmt.Errorf("code: spec %T", sp)
	}
	return json.Marshal(Answer{Files: map[string]string{s.Entry: s.Solution}})
}

// WrongAnswers submits each wrong/ file as the entry.
func (codeType) WrongAnswers(sp exercise.Spec) (map[string]json.RawMessage, error) {
	s, ok := sp.(Spec)
	if !ok {
		return nil, fmt.Errorf("code: spec %T", sp)
	}
	out := map[string]json.RawMessage{}
	for name, src := range s.Wrong {
		b, err := json.Marshal(Answer{Files: map[string]string{s.Entry: src}})
		if err != nil {
			return nil, err
		}
		out[name] = b
	}
	return out, nil
}
