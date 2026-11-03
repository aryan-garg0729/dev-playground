// Package sql is the "write a query" exercise type. The learner writes one
// SQLite statement; it runs on fresh in-memory databases in the grading
// container (never inside the API process), and sql_result_equals compares the
// rows. In the browser, Run uses sql.js on the public setup data.
package sql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/checks"
	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/runner"
)

const (
	maxPromptBytes = 64 << 10
	maxSetupBytes  = 1 << 20
	maxQueryBytes  = 16 << 10
	image          = "docker-python" // server runner for sqljs exercises
)

// Check is one compiled check.
type Check struct {
	Kind   string          `json:"kind"`
	Config json.RawMessage `json:"config"`
}

// Spec is the compiled sql spec (private_spec). Datasets[0] is the public
// setup; further datasets come from hidden checks' setup files and are server-only.
// Solution is stored because hidden checks compare against it at grade time.
type Spec struct {
	Prompt    string            `json:"prompt"`
	Starter   string            `json:"starter"`
	Run       string            `json:"run"`
	TimeoutMs int               `json:"timeout_ms"`
	Datasets  []string          `json:"datasets"`
	Solution  string            `json:"solution"`
	Checks    []Check           `json:"checks"`
	Wrong     map[string]string `json:"-"`
}

type manifest struct {
	Prompt   string `yaml:"prompt"`
	Setup    string `yaml:"setup"`
	Starter  string `yaml:"starter"`
	Solution string `yaml:"solution"`
	Run      string `yaml:"run"`
	Limits   struct {
		TimeoutMs int `yaml:"timeout_ms"`
	} `yaml:"limits"`
	Checks []yaml.Node `yaml:"checks"`
}

// VisibleCheck is a visible check as the browser sees it (sandbox SqlCheck).
type VisibleCheck struct {
	Name     string              `json:"name"`
	Ordered  bool                `json:"ordered"`
	Expected checks.ExpectedRows `json:"expected"`
}

// PublicView is what the browser receives. The setup data is public; hidden
// datasets and the solution are not.
type PublicView struct {
	Prompt        string         `json:"prompt"`
	Setup         string         `json:"setup"`
	Starter       string         `json:"starter"`
	Run           string         `json:"run"`
	TimeoutMs     int            `json:"timeout_ms"`
	VisibleChecks []VisibleCheck `json:"visible_checks"`
}

// Answer is the learner's answer shape.
type Answer struct {
	Query string `json:"query"`
}

type sqlType struct{}

func init() { exercise.Register(sqlType{}) }

func (sqlType) Name() string { return "sql" }

func readFile(dir fs.FS, field, name string, limit int, required bool) (string, error) {
	if name == "" {
		if required {
			return "", fmt.Errorf("%s: required", field)
		}
		return "", nil
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

func (sqlType) ParseSpec(raw yaml.Node, dir fs.FS) (exercise.Spec, error) {
	var m manifest
	if err := exercise.DecodeStrict(raw, &m); err != nil {
		return nil, err
	}
	s := Spec{Run: m.Run, TimeoutMs: m.Limits.TimeoutMs}
	if s.Run == "" {
		s.Run = exercise.RunBrowser
	}
	if s.TimeoutMs == 0 {
		s.TimeoutMs = 2000
	}
	var err error
	if s.Prompt, err = readFile(dir, "prompt", m.Prompt, maxPromptBytes, true); err != nil {
		return nil, err
	}
	setup, err := readFile(dir, "setup", m.Setup, maxSetupBytes, true)
	if err != nil {
		return nil, err
	}
	s.Datasets = []string{setup}
	if s.Starter, err = readFile(dir, "starter", m.Starter, maxQueryBytes, false); err != nil {
		return nil, err
	}
	if s.Solution, err = readFile(dir, "solution", m.Solution, maxQueryBytes, true); err != nil {
		return nil, err
	}
	files := map[string]int{}
	for i, node := range m.Checks {
		var fields map[string]any
		if err := node.Decode(&fields); err != nil {
			return nil, fmt.Errorf("checks[%d]: %w", i, err)
		}
		if kind, _ := fields["kind"].(string); kind != "sql_result_equals" {
			return nil, fmt.Errorf("checks[%d]: kind %q is not supported for sql (sql_result_equals)", i, fields["kind"])
		}
		delete(fields, "kind")
		if _, ok := fields["dataset"]; ok {
			return nil, fmt.Errorf("checks[%d]: dataset is set by cq build; use setup", i)
		}
		cfgJSON, err := json.Marshal(fields)
		if err != nil {
			return nil, fmt.Errorf("checks[%d]: %w", i, err)
		}
		if err := checks.Validate("sql_result_equals", cfgJSON); err != nil {
			return nil, fmt.Errorf("checks[%d]: %w", i, err)
		}
		var cfg checks.SQLResultEqualsConfig
		_ = json.Unmarshal(cfgJSON, &cfg)
		if cfg.Setup != "" {
			if cfg.Visible == nil || *cfg.Visible {
				return nil, fmt.Errorf("checks[%d]: a check with its own setup must be hidden (visible: false); the browser only has the public setup", i)
			}
			idx, ok := files[cfg.Setup]
			if !ok {
				data, err := readFile(dir, fmt.Sprintf("checks[%d].setup", i), cfg.Setup, maxSetupBytes, true)
				if err != nil {
					return nil, err
				}
				s.Datasets = append(s.Datasets, data)
				idx = len(s.Datasets) - 1
				files[cfg.Setup] = idx
			}
			cfg.Dataset, cfg.Setup = idx, ""
		}
		compiled, err := json.Marshal(cfg)
		if err != nil {
			return nil, err
		}
		s.Checks = append(s.Checks, Check{Kind: "sql_result_equals", Config: compiled})
	}
	if s.Wrong, err = readWrong(dir); err != nil {
		return nil, err
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return s, nil
}

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
		if path.Ext(name) != ".sql" {
			return nil, fmt.Errorf("%s: wrong/ may hold only .sql files", name)
		}
		src, err := readFile(dir, "wrong", name, maxQueryBytes, true)
		if err != nil {
			return nil, err
		}
		out[name] = src
	}
	return out, nil
}

func (s Spec) validate() error {
	if !exercise.ValidRunMode(s.Run) {
		return fmt.Errorf("run: %q must be browser, server or off", s.Run)
	}
	if s.TimeoutMs < 100 || s.TimeoutMs > runner.MaxTimeoutMs {
		return fmt.Errorf("limits.timeout_ms: must be 100-%d", runner.MaxTimeoutMs)
	}
	if strings.TrimSpace(s.Prompt) == "" {
		return errors.New("prompt: file is empty")
	}
	if len(s.Datasets) == 0 || strings.TrimSpace(s.Datasets[0]) == "" {
		return errors.New("setup: file is empty")
	}
	if strings.TrimSpace(s.Solution) == "" {
		return errors.New("solution: file is empty")
	}
	if len(s.Checks) == 0 {
		return errors.New("checks: need at least one")
	}
	for i, c := range s.Checks {
		if err := checks.Validate(c.Kind, c.Config); err != nil {
			return fmt.Errorf("checks[%d]: %w", i, err)
		}
		var cfg checks.SQLResultEqualsConfig
		if err := json.Unmarshal(c.Config, &cfg); err != nil || cfg.Dataset >= len(s.Datasets) {
			return fmt.Errorf("checks[%d]: bad dataset", i)
		}
	}
	return nil
}

func (sqlType) LoadSpec(private json.RawMessage) (exercise.Spec, error) {
	var s Spec
	if err := json.Unmarshal(private, &s); err != nil {
		return nil, fmt.Errorf("sql: load spec: %w", err)
	}
	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("sql: load spec: %w", err)
	}
	return s, nil
}

func specOf(ex exercise.Exercise) (Spec, error) {
	s, ok := ex.Spec.(Spec)
	if !ok {
		return Spec{}, fmt.Errorf("sql: exercise %s has spec %T", ex.ID, ex.Spec)
	}
	return s, nil
}

// visible returns the visible checks (all on dataset 0, all with expected rows).
func (s Spec) visible() ([]VisibleCheck, []Check) {
	views, raw := []VisibleCheck{}, []Check{}
	for _, c := range s.Checks {
		var cfg checks.SQLResultEqualsConfig
		if json.Unmarshal(c.Config, &cfg) != nil || (cfg.Visible != nil && !*cfg.Visible) || cfg.Expected == nil {
			continue
		}
		name := cfg.Name
		if name == "" {
			name = "Returns the right rows"
		}
		views = append(views, VisibleCheck{Name: name, Ordered: cfg.Ordered, Expected: *cfg.Expected})
		raw = append(raw, c)
	}
	return views, raw
}

func (sqlType) PublicView(ex exercise.Exercise) (json.RawMessage, error) {
	s, err := specOf(ex)
	if err != nil {
		return nil, err
	}
	views, _ := s.visible()
	b, err := json.Marshal(PublicView{
		Prompt: s.Prompt, Setup: s.Datasets[0], Starter: s.Starter, Run: s.Run,
		TimeoutMs: s.TimeoutMs, VisibleChecks: views,
	})
	if err != nil {
		return nil, fmt.Errorf("sql: public view: %w", err)
	}
	return b, nil
}

func parseAnswer(raw json.RawMessage) (string, error) {
	var a Answer
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&a); err != nil {
		return "", fmt.Errorf("%w: %w", exercise.ErrInvalidAnswer, err)
	}
	if len(a.Query) > maxQueryBytes {
		return "", fmt.Errorf("%w: query is larger than %d bytes", exercise.ErrInvalidAnswer, maxQueryBytes)
	}
	if strings.ContainsRune(a.Query, 0) {
		return "", fmt.Errorf("%w: query contains a NUL byte", exercise.ErrInvalidAnswer)
	}
	return a.Query, nil
}

type dataset struct {
	Setup    string `json:"setup"`
	Solution bool   `json:"solution"`
}

type harnessJob struct {
	Mode      string    `json:"mode"`
	Query     string    `json:"query"`
	Solution  string    `json:"solution"`
	Datasets  []dataset `json:"datasets"`
	TimeoutMs int       `json:"timeout_ms"`
}

// execute runs query over the first n datasets; withSolution also runs the reference.
func (s Spec) execute(ctx context.Context, query string, n int, withSolution bool) (*checks.Report, int, error) {
	run, ok := runner.FromContext(ctx)
	if !ok {
		return nil, 0, errors.New("sql: no runner in context")
	}
	job := harnessJob{Mode: "sql", Query: query, TimeoutMs: s.TimeoutMs}
	if withSolution {
		job.Solution = s.Solution
	}
	for _, d := range s.Datasets[:n] {
		job.Datasets = append(job.Datasets, dataset{Setup: d, Solution: withSolution})
	}
	stdin, err := json.Marshal(job)
	if err != nil {
		return nil, 0, fmt.Errorf("sql: encode job: %w", err)
	}
	out, err := run.Run(ctx, runner.Job{Image: runner.Images[image], Stdin: stdin, Limits: runner.Limits{TimeoutMs: s.TimeoutMs}})
	if err != nil {
		return nil, 0, fmt.Errorf("sql: run: %w", err)
	}
	if out.TimedOut {
		return &checks.Report{Version: 1, TimedOut: true}, out.RuntimeMs, nil
	}
	report, err := checks.ParseReport(out.Stdout)
	if err != nil {
		return nil, 0, fmt.Errorf("sql: %w (exit %d)", err, out.ExitCode)
	}
	if report.Error != "" {
		return nil, 0, fmt.Errorf("sql: harness: %s", report.Error)
	}
	return report, report.RuntimeMs, nil
}

func (sqlType) Grade(ctx context.Context, ex exercise.Exercise, sub exercise.Submission) (exercise.Result, error) {
	s, err := specOf(ex)
	if err != nil {
		return exercise.Result{}, err
	}
	query, err := parseAnswer(sub.Answer)
	if err != nil {
		return exercise.Result{}, err
	}
	report, ms, err := s.execute(ctx, query, len(s.Datasets), true)
	if err != nil {
		return exercise.Result{}, err
	}
	var rows []exercise.CheckResult
	for _, c := range s.Checks {
		r, err := checks.Run(ctx, c.Kind, checks.CheckInput{Config: c.Config, Answer: sub.Answer, Report: report})
		if err != nil {
			return exercise.Result{}, fmt.Errorf("sql: %w", err)
		}
		rows = append(rows, r...)
	}
	res := exercise.NewResult(rows)
	res.RuntimeMs = ms
	if report.TimedOut {
		res.Status = exercise.StatusTimeout
	}
	return res, nil
}

// Run executes the query on the public setup only and returns the table plus
// the visible checks.
func (sqlType) Run(ctx context.Context, ex exercise.Exercise, answer json.RawMessage) (exercise.RunOutput, error) {
	s, err := specOf(ex)
	if err != nil {
		return exercise.RunOutput{}, err
	}
	if s.Run != exercise.RunServer {
		return exercise.RunOutput{}, exercise.ErrRunUnavailable
	}
	query, err := parseAnswer(answer)
	if err != nil {
		return exercise.RunOutput{}, err
	}
	report, ms, err := s.execute(ctx, query, 1, false)
	if err != nil {
		return exercise.RunOutput{}, err
	}
	out := exercise.RunOutput{Checks: []exercise.CheckResult{}, RuntimeMs: ms, TimedOut: report.TimedOut}
	if report.TimedOut {
		out.Error = fmt.Sprintf("Stopped after %d ms", s.TimeoutMs)
		return out, nil
	}
	learner := report.SQL.Results[0].Learner
	if learner.Error != "" {
		out.Error = learner.Error
	} else {
		out.Table = &exercise.Table{Columns: learner.Columns, Rows: learner.Rows, Truncated: learner.Truncated}
	}
	_, visible := s.visible()
	for _, c := range visible {
		r, err := checks.Run(ctx, c.Kind, checks.CheckInput{Config: c.Config, Answer: answer, Report: report})
		if err != nil {
			return exercise.RunOutput{}, fmt.Errorf("sql: %w", err)
		}
		out.Checks = append(out.Checks, r...)
	}
	return out, nil
}

// ReferenceAnswer is the solution query.
func (sqlType) ReferenceAnswer(sp exercise.Spec) (json.RawMessage, error) {
	s, ok := sp.(Spec)
	if !ok {
		return nil, fmt.Errorf("sql: spec %T", sp)
	}
	return json.Marshal(Answer{Query: s.Solution})
}

// WrongAnswers submits each wrong/*.sql file.
func (sqlType) WrongAnswers(sp exercise.Spec) (map[string]json.RawMessage, error) {
	s, ok := sp.(Spec)
	if !ok {
		return nil, fmt.Errorf("sql: spec %T", sp)
	}
	out := map[string]json.RawMessage{}
	for name, q := range s.Wrong {
		b, err := json.Marshal(Answer{Query: q})
		if err != nil {
			return nil, err
		}
		out[name] = b
	}
	return out, nil
}
