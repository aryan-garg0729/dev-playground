// Package projectcli is the `project_cli` exercise type: lessons whose steps
// (shell commands and HTTP requests) run on the learner's machine through the
// `codequest` CLI. The CLI reports each step's raw output; Grade judges that
// output against the manifest's tests on the server.
//
// The output is reported by the client, so a determined learner can forge a
// pass. This is the one accepted exception to "never trust a client-reported
// pass" (docs/SPEC.md, decision 2026-11-04). The server still owns the verdict:
// it never trusts client-computed variables or pass flags, only raw output.
package projectcli

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
)

// Limits on manifests and reported output.
const (
	MaxSteps        = 20
	MaxOutputBytes  = 16 << 10 // per step: stdout or response body
	MaxHeaderBytes  = 4 << 10  // per step: all reported response headers
	MaxHeaders      = 50
	maxPromptBytes  = 64 << 10
	maxAnswerFile   = 100 << 10
	maxCommandBytes = 4 << 10
)

// Spec is the compiled project_cli spec. Every field except Solution and
// Wrong is public: the CLI prints the tests it will be judged on.
type Spec struct {
	Prompt    string   `json:"prompt"`
	BaseURL   string   `json:"base_url,omitempty"`
	AllowedOS []string `json:"allowed_os,omitempty"`
	Steps     []Step   `json:"steps"`

	Solution []byte            `json:"-"` // solution.json, cq test only
	Wrong    map[string][]byte `json:"-"` // wrong/*.json, cq test only
}

// Step is one command or request. Exactly one of CLICommand and HTTPRequest is set.
type Step struct {
	Description string       `json:"description" yaml:"description"`
	CLICommand  *CLICommand  `json:"cli_command,omitempty" yaml:"cli_command"`
	HTTPRequest *HTTPRequest `json:"http_request,omitempty" yaml:"http_request"`
}

// CLICommand runs Command in the learner's shell.
type CLICommand struct {
	Command  string    `json:"command" yaml:"command"`
	Captures []Capture `json:"captures,omitempty" yaml:"captures"`
	// Every project_cli test is visible (the CLI prints them all), hence the
	// JSON name: "tests" means hidden tests too elsewhere and never reaches the browser.
	Tests []CLITest `json:"visible_tests" yaml:"tests"`
}

// HTTPRequest sends Request from the learner's machine.
type HTTPRequest struct {
	Request  Request    `json:"request" yaml:"request"`
	Captures []Capture  `json:"captures,omitempty" yaml:"captures"`
	Tests    []HTTPTest `json:"visible_tests" yaml:"tests"` // all visible; see CLICommand.Tests
}

// Request is an HTTP request; ${var} is interpolated in URL, header values
// and string values of BodyJSON.
type Request struct {
	Method          string            `json:"method" yaml:"method"`
	URL             string            `json:"url" yaml:"url"`
	Headers         map[string]string `json:"headers,omitempty" yaml:"headers"`
	BodyJSON        any               `json:"body_json,omitempty" yaml:"body_json"`
	FollowRedirects bool              `json:"follow_redirects,omitempty" yaml:"follow_redirects"`
}

// Capture stores a value from a step's output as ${Name} for later steps
// and this step's tests. Regex captures from stdout or the body (first
// submatch, or the whole match); JSONPath reads the JSON body (HTTP only).
type Capture struct {
	Name     string `json:"name" yaml:"name"`
	Regex    string `json:"regex,omitempty" yaml:"regex"`
	JSONPath string `json:"json_path,omitempty" yaml:"json_path"`
}

// CLITest is one assertion on a command; exactly one field is set.
type CLITest struct {
	ExitCode           *int     `json:"exit_code,omitempty" yaml:"exit_code"`
	StdoutContainsAll  []string `json:"stdout_contains_all,omitempty" yaml:"stdout_contains_all"`
	StdoutContainsNone []string `json:"stdout_contains_none,omitempty" yaml:"stdout_contains_none"`
	StdoutLinesGT      *int     `json:"stdout_lines_gt,omitempty" yaml:"stdout_lines_gt"`
}

// HTTPTest is one assertion on a response; exactly one field is set.
type HTTPTest struct {
	StatusCode   *int        `json:"status_code,omitempty" yaml:"status_code"`
	BodyContains *string     `json:"body_contains,omitempty" yaml:"body_contains"`
	JSONBody     *JSONTest   `json:"json_body,omitempty" yaml:"json_body"`
	Header       *HeaderTest `json:"header,omitempty" yaml:"header"`
}

// JSONTest compares the value at Path in the JSON body with Value.
type JSONTest struct {
	Path  string `json:"path" yaml:"path"`
	Op    string `json:"op" yaml:"op"` // eq | ne | gt | lt | contains
	Value any    `json:"value" yaml:"value"`
}

// HeaderTest checks that response header Name (case-insensitive) contains Contains.
type HeaderTest struct {
	Name     string `json:"name" yaml:"name"`
	Contains string `json:"contains" yaml:"contains"`
}

// Answer is what the CLI reports: one result per step, in order.
type Answer struct {
	// BaseURL is the base URL the learner ran against (their override or the
	// default); used only to interpolate ${baseURL} in expected values.
	BaseURL string       `json:"base_url,omitempty"`
	Steps   []StepResult `json:"steps"`
}

// StepResult is one step's raw output. Exactly one of CLI, HTTP and Skipped is set.
type StepResult struct {
	CLI     *CLIResult  `json:"cli,omitempty"`
	HTTP    *HTTPResult `json:"http,omitempty"`
	Skipped bool        `json:"skipped,omitempty"` // not run: it needed a value an earlier step did not capture
}

// CLIResult is a command's exit code and combined stdout.
type CLIResult struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
}

// HTTPResult is a response, or Error when no response arrived (Status 0).
type HTTPResult struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body"`
	Error   string            `json:"error,omitempty"`
}

var (
	varRe         = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)
	captureNameRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,31}$`)
	allowedOS     = []string{"linux", "darwin", "windows"}
	methods       = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead}
	jsonOps       = []string{"eq", "ne", "gt", "lt", "contains"}
)

// BaseURLVar is the variable that holds the lesson's base URL.
const BaseURLVar = "baseURL"

// Vars lists the ${var} names used in s.
func Vars(s string) []string {
	var out []string
	for _, m := range varRe.FindAllStringSubmatch(s, -1) {
		out = append(out, m[1])
	}
	return out
}

// Interpolate replaces ${var} with its value; unknown variables are left as is.
func Interpolate(s string, vars map[string]string) string {
	return varRe.ReplaceAllStringFunc(s, func(m string) string {
		if v, ok := vars[m[2:len(m)-1]]; ok {
			return v
		}
		return m
	})
}

// InterpolateJSON interpolates every string inside a decoded JSON value.
func InterpolateJSON(v any, vars map[string]string) any {
	switch t := v.(type) {
	case string:
		return Interpolate(t, vars)
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, x := range t {
			out[k] = InterpolateJSON(x, vars)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, x := range t {
			out[i] = InterpolateJSON(x, vars)
		}
		return out
	default:
		return v
	}
}

// Inputs are the strings a step needs before it can run (command, URL,
// headers, body), so their variables must come from earlier steps.
func (st Step) Inputs() []string {
	if st.CLICommand != nil {
		return []string{st.CLICommand.Command}
	}
	r := st.HTTPRequest.Request
	in := []string{r.URL}
	for _, v := range r.Headers {
		in = append(in, v)
	}
	collectStrings(r.BodyJSON, &in)
	return in
}

// Captures returns the step's captures.
func (st Step) Captures() []Capture {
	if st.CLICommand != nil {
		return st.CLICommand.Captures
	}
	return st.HTTPRequest.Captures
}

// Missing lists the input variables of st that vars does not define.
func (st Step) Missing(vars map[string]string) []string {
	var miss []string
	for _, s := range st.Inputs() {
		for _, v := range Vars(s) {
			if _, ok := vars[v]; !ok && !slices.Contains(miss, v) {
				miss = append(miss, v)
			}
		}
	}
	return miss
}

func collectStrings(v any, out *[]string) {
	switch t := v.(type) {
	case string:
		*out = append(*out, t)
	case map[string]any:
		for _, x := range t {
			collectStrings(x, out)
		}
	case []any:
		for _, x := range t {
			collectStrings(x, out)
		}
	}
}

// validate checks a spec read from a manifest or private_spec.
func (s Spec) validate() error {
	if strings.TrimSpace(s.Prompt) == "" {
		return errors.New("prompt: file is empty")
	}
	if s.BaseURL != "" && !strings.HasPrefix(s.BaseURL, "http://") && !strings.HasPrefix(s.BaseURL, "https://") {
		return fmt.Errorf("base_url: %q must start with http:// or https://", s.BaseURL)
	}
	for _, o := range s.AllowedOS {
		if !slices.Contains(allowedOS, o) {
			return fmt.Errorf("allowed_os: %q is not one of %v", o, allowedOS)
		}
	}
	if len(s.Steps) == 0 || len(s.Steps) > MaxSteps {
		return fmt.Errorf("steps: need 1 to %d steps, got %d", MaxSteps, len(s.Steps))
	}
	defined := map[string]bool{}
	if s.BaseURL != "" {
		defined[BaseURLVar] = true
	}
	for i, st := range s.Steps {
		if err := st.validate(defined); err != nil {
			return fmt.Errorf("steps[%d]: %w", i, err)
		}
	}
	return nil
}

func (st Step) validate(defined map[string]bool) error {
	if strings.TrimSpace(st.Description) == "" {
		return errors.New("description: required")
	}
	if (st.CLICommand == nil) == (st.HTTPRequest == nil) {
		return errors.New("set exactly one of cli_command and http_request")
	}
	if st.CLICommand != nil {
		c := st.CLICommand
		if strings.TrimSpace(c.Command) == "" || len(c.Command) > maxCommandBytes {
			return fmt.Errorf("cli_command.command: required, at most %d bytes", maxCommandBytes)
		}
	} else {
		r := st.HTTPRequest.Request
		if !slices.Contains(methods, r.Method) {
			return fmt.Errorf("http_request.request.method: %q is not one of %v", r.Method, methods)
		}
		u := r.URL
		if !strings.HasPrefix(u, "${"+BaseURLVar+"}") && !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
			return fmt.Errorf("http_request.request.url: %q must start with ${baseURL}, http:// or https://", u)
		}
	}
	for _, in := range st.Inputs() {
		for _, v := range Vars(in) {
			if !defined[v] {
				return fmt.Errorf("uses ${%s}, which no earlier step captures", v)
			}
		}
	}
	// Captures become available to this step's tests and later steps.
	for j, c := range st.Captures() {
		if !captureNameRe.MatchString(c.Name) || c.Name == BaseURLVar {
			return fmt.Errorf("captures[%d].name: %q is not a valid name", j, c.Name)
		}
		if defined[c.Name] {
			return fmt.Errorf("captures[%d].name: %q is already defined", j, c.Name)
		}
		switch {
		case c.Regex != "" && c.JSONPath == "":
			if _, err := regexp.Compile(c.Regex); err != nil {
				return fmt.Errorf("captures[%d].regex: %w", j, err)
			}
		case c.JSONPath != "" && c.Regex == "" && st.HTTPRequest != nil:
			if _, err := parsePath(c.JSONPath); err != nil {
				return fmt.Errorf("captures[%d].json_path: %w", j, err)
			}
		default:
			return fmt.Errorf("captures[%d]: set exactly one of regex and json_path (json_path is for http_request only)", j)
		}
		defined[c.Name] = true
	}
	var expected []string
	if st.CLICommand != nil {
		if len(st.CLICommand.Tests) == 0 {
			return errors.New("cli_command.tests: need at least one test")
		}
		for j, t := range st.CLICommand.Tests {
			if err := t.validate(); err != nil {
				return fmt.Errorf("cli_command.tests[%d]: %w", j, err)
			}
			expected = append(expected, t.StdoutContainsAll...)
			expected = append(expected, t.StdoutContainsNone...)
		}
	} else {
		if len(st.HTTPRequest.Tests) == 0 {
			return errors.New("http_request.tests: need at least one test")
		}
		for j, t := range st.HTTPRequest.Tests {
			if err := t.validate(); err != nil {
				return fmt.Errorf("http_request.tests[%d]: %w", j, err)
			}
			if t.BodyContains != nil {
				expected = append(expected, *t.BodyContains)
			}
			if t.Header != nil {
				expected = append(expected, t.Header.Contains)
			}
			if t.JSONBody != nil {
				collectStrings(t.JSONBody.Value, &expected)
			}
		}
	}
	for _, e := range expected {
		for _, v := range Vars(e) {
			if !defined[v] {
				return fmt.Errorf("a test uses ${%s}, which is not captured by this or an earlier step", v)
			}
		}
	}
	return nil
}

func (t CLITest) validate() error {
	n := 0
	if t.ExitCode != nil {
		n++
	}
	if t.StdoutContainsAll != nil {
		n++
		if len(t.StdoutContainsAll) == 0 {
			return errors.New("stdout_contains_all: list at least one string")
		}
	}
	if t.StdoutContainsNone != nil {
		n++
		if len(t.StdoutContainsNone) == 0 {
			return errors.New("stdout_contains_none: list at least one string")
		}
	}
	if t.StdoutLinesGT != nil {
		n++
	}
	if n != 1 {
		return errors.New("set exactly one of exit_code, stdout_contains_all, stdout_contains_none, stdout_lines_gt")
	}
	return nil
}

func (t HTTPTest) validate() error {
	n := 0
	if t.StatusCode != nil {
		n++
		if *t.StatusCode < 100 || *t.StatusCode > 599 {
			return fmt.Errorf("status_code: %d is not an HTTP status", *t.StatusCode)
		}
	}
	if t.BodyContains != nil {
		n++
	}
	if t.JSONBody != nil {
		n++
		if _, err := parsePath(t.JSONBody.Path); err != nil {
			return fmt.Errorf("json_body.path: %w", err)
		}
		if !slices.Contains(jsonOps, t.JSONBody.Op) {
			return fmt.Errorf("json_body.op: %q is not one of %v", t.JSONBody.Op, jsonOps)
		}
		if t.JSONBody.Op == "gt" || t.JSONBody.Op == "lt" {
			if _, ok := toFloat(normalize(t.JSONBody.Value)); !ok {
				return fmt.Errorf("json_body.value: %s needs a number", t.JSONBody.Op)
			}
		}
	}
	if t.Header != nil {
		n++
		if t.Header.Name == "" {
			return errors.New("header.name: required")
		}
	}
	if n != 1 {
		return errors.New("set exactly one of status_code, body_contains, json_body, header")
	}
	return nil
}
