package projectcli

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

const (
	maxResultStdout = 16 << 10
	maxRepr         = 200
	maxBaseURL      = 2 << 10
)

// short truncates s for Got/Expected.
func short(s string) string {
	if utf8.RuneCountInString(s) <= maxRepr {
		return s
	}
	r := []rune(s)
	return string(r[:maxRepr]) + "…"
}

// ApplyCaptures adds the step's captured values from its result to vars.
// The CLI uses it to interpolate later steps; the server uses it to rebuild
// the same values from the reported output (never from client variables).
// A capture that does not match leaves its variable undefined.
func ApplyCaptures(st Step, res StepResult, vars map[string]string) {
	var text string
	var body any
	bodyParsed := false
	switch {
	case st.CLICommand != nil && res.CLI != nil:
		text = res.CLI.Stdout
	case st.HTTPRequest != nil && res.HTTP != nil:
		text = res.HTTP.Body
	default:
		return
	}
	for _, c := range st.Captures() {
		if c.Regex != "" {
			re, err := regexp.Compile(c.Regex)
			if err != nil {
				continue
			}
			m := re.FindStringSubmatch(text)
			switch {
			case m == nil:
			case len(m) > 1:
				vars[c.Name] = m[1]
			default:
				vars[c.Name] = m[0]
			}
			continue
		}
		if !bodyParsed {
			bodyParsed = true
			if err := json.Unmarshal([]byte(text), &body); err != nil {
				body = nil
			}
		}
		if body == nil {
			continue
		}
		v, ok := lookup(body, c.JSONPath)
		if !ok {
			continue
		}
		if s, isStr := v.(string); isStr {
			vars[c.Name] = s
		} else {
			vars[c.Name] = canon(v)
		}
	}
}

// decodeAnswer parses and bounds-checks a reported answer against s.
func decodeAnswer(s Spec, raw json.RawMessage) (Answer, error) {
	var a Answer
	if err := json.Unmarshal(raw, &a); err != nil {
		return a, fmt.Errorf("%w: %w", exercise.ErrInvalidAnswer, err)
	}
	if len(a.BaseURL) > maxBaseURL {
		return a, fmt.Errorf("%w: base_url is too long", exercise.ErrInvalidAnswer)
	}
	if len(a.Steps) != len(s.Steps) {
		return a, fmt.Errorf("%w: reported %d steps, the lesson has %d", exercise.ErrInvalidAnswer, len(a.Steps), len(s.Steps))
	}
	for i, r := range a.Steps {
		st := s.Steps[i]
		n := 0
		if r.CLI != nil {
			n++
		}
		if r.HTTP != nil {
			n++
		}
		if r.Skipped {
			n++
		}
		if n != 1 {
			return a, fmt.Errorf("%w: steps[%d]: set exactly one of cli, http and skipped", exercise.ErrInvalidAnswer, i)
		}
		switch {
		case r.CLI != nil:
			if st.CLICommand == nil {
				return a, fmt.Errorf("%w: steps[%d] is an http_request step", exercise.ErrInvalidAnswer, i)
			}
			if len(r.CLI.Stdout) > MaxOutputBytes {
				return a, fmt.Errorf("%w: steps[%d].cli.stdout is larger than %d bytes", exercise.ErrInvalidAnswer, i, MaxOutputBytes)
			}
		case r.HTTP != nil:
			if st.HTTPRequest == nil {
				return a, fmt.Errorf("%w: steps[%d] is a cli_command step", exercise.ErrInvalidAnswer, i)
			}
			if len(r.HTTP.Body) > MaxOutputBytes {
				return a, fmt.Errorf("%w: steps[%d].http.body is larger than %d bytes", exercise.ErrInvalidAnswer, i, MaxOutputBytes)
			}
			if r.HTTP.Status != 0 && (r.HTTP.Status < 100 || r.HTTP.Status > 599) {
				return a, fmt.Errorf("%w: steps[%d].http.status is not an HTTP status", exercise.ErrInvalidAnswer, i)
			}
			if len(r.HTTP.Headers) > MaxHeaders {
				return a, fmt.Errorf("%w: steps[%d].http.headers: more than %d", exercise.ErrInvalidAnswer, i, MaxHeaders)
			}
			size := len(r.HTTP.Error)
			for k, v := range r.HTTP.Headers {
				size += len(k) + len(v)
			}
			if size > MaxHeaderBytes {
				return a, fmt.Errorf("%w: steps[%d].http: headers and error are larger than %d bytes", exercise.ErrInvalidAnswer, i, MaxHeaderBytes)
			}
		}
	}
	return a, nil
}

// grade judges a decoded answer. One visible check row per test.
func grade(s Spec, a Answer) exercise.Result {
	vars := map[string]string{}
	if s.BaseURL != "" {
		vars[BaseURLVar] = strings.TrimSuffix(s.BaseURL, "/")
		if a.BaseURL != "" {
			vars[BaseURLVar] = strings.TrimSuffix(a.BaseURL, "/")
		}
	}
	var rows []exercise.CheckResult
	var transcript strings.Builder
	for i, st := range s.Steps {
		r := a.Steps[i]
		label := fmt.Sprintf("Step %d", i+1)
		names := testNames(st, vars)
		fail := func(feedback string) {
			for _, n := range names {
				rows = append(rows, exercise.CheckResult{Name: label + ": " + n, Visible: true, Feedback: feedback})
			}
		}
		if miss := st.Missing(vars); len(miss) > 0 {
			fail(fmt.Sprintf("Needs ${%s}, which an earlier step did not capture", miss[0]))
			continue
		}
		if r.Skipped {
			fail("This step was not run")
			continue
		}
		ApplyCaptures(st, r, vars)
		if st.CLICommand != nil {
			fmt.Fprintf(&transcript, "$ %s\n%s", Interpolate(st.CLICommand.Command, vars), r.CLI.Stdout)
			for j, t := range st.CLICommand.Tests {
				row := cliCheck(t, *r.CLI, vars)
				row.Name = label + ": " + names[j]
				rows = append(rows, row)
			}
			continue
		}
		req := st.HTTPRequest.Request
		fmt.Fprintf(&transcript, "> %s %s\n", req.Method, Interpolate(req.URL, vars))
		if r.HTTP.Error != "" {
			fmt.Fprintf(&transcript, "%s\n", r.HTTP.Error)
		} else {
			fmt.Fprintf(&transcript, "< %d\n%s", r.HTTP.Status, r.HTTP.Body)
		}
		for j, t := range st.HTTPRequest.Tests {
			row := httpCheck(t, *r.HTTP, vars)
			row.Name = label + ": " + names[j]
			rows = append(rows, row)
		}
		if !strings.HasSuffix(transcript.String(), "\n") {
			transcript.WriteByte('\n')
		}
	}
	res := exercise.NewResult(rows)
	out := transcript.String()
	if len(out) > maxResultStdout {
		out = out[:maxResultStdout]
		for !utf8.ValidString(out) {
			out = out[:len(out)-1]
		}
	}
	res.Stdout = out
	return res
}

// testNames describes each test of a step, e.g. "exit code is 0".
func testNames(st Step, vars map[string]string) []string {
	var out []string
	if st.CLICommand != nil {
		for _, t := range st.CLICommand.Tests {
			out = append(out, DescribeCLITest(t, vars))
		}
		return out
	}
	for _, t := range st.HTTPRequest.Tests {
		out = append(out, DescribeHTTPTest(t, vars))
	}
	return out
}

func quoteAll(ss []string, vars map[string]string) string {
	q := make([]string, len(ss))
	for i, s := range ss {
		q[i] = strconv.Quote(Interpolate(s, vars))
	}
	return strings.Join(q, ", ")
}

// DescribeCLITest is the human-readable form of a test (CLI and check names).
func DescribeCLITest(t CLITest, vars map[string]string) string {
	switch {
	case t.ExitCode != nil:
		return fmt.Sprintf("exit code is %d", *t.ExitCode)
	case t.StdoutContainsAll != nil:
		return "output contains " + quoteAll(t.StdoutContainsAll, vars)
	case t.StdoutContainsNone != nil:
		return "output does not contain " + quoteAll(t.StdoutContainsNone, vars)
	default:
		return fmt.Sprintf("output has more than %d lines", *t.StdoutLinesGT)
	}
}

// DescribeHTTPTest is the human-readable form of a test.
func DescribeHTTPTest(t HTTPTest, vars map[string]string) string {
	switch {
	case t.StatusCode != nil:
		return fmt.Sprintf("status code is %d", *t.StatusCode)
	case t.BodyContains != nil:
		return "body contains " + strconv.Quote(Interpolate(*t.BodyContains, vars))
	case t.JSONBody != nil:
		ops := map[string]string{"eq": "is", "ne": "is not", "gt": "is greater than", "lt": "is less than", "contains": "contains"}
		return fmt.Sprintf("JSON %s %s %s", t.JSONBody.Path, ops[t.JSONBody.Op], repr(InterpolateJSON(normalize(t.JSONBody.Value), vars)))
	default:
		return fmt.Sprintf("header %s contains %q", t.Header.Name, Interpolate(t.Header.Contains, vars))
	}
}

func lineCount(s string) int {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func cliCheck(t CLITest, r CLIResult, vars map[string]string) exercise.CheckResult {
	row := exercise.CheckResult{Visible: true}
	switch {
	case t.ExitCode != nil:
		row.Passed = r.ExitCode == *t.ExitCode
		row.Got, row.Expected = strconv.Itoa(r.ExitCode), strconv.Itoa(*t.ExitCode)
	case t.StdoutContainsAll != nil:
		row.Passed = true
		for _, w := range t.StdoutContainsAll {
			w = Interpolate(w, vars)
			if !strings.Contains(r.Stdout, w) {
				row.Passed = false
				row.Feedback = fmt.Sprintf("Output is missing %q", w)
				break
			}
		}
		row.Got = short(r.Stdout)
	case t.StdoutContainsNone != nil:
		row.Passed = true
		for _, w := range t.StdoutContainsNone {
			w = Interpolate(w, vars)
			if strings.Contains(r.Stdout, w) {
				row.Passed = false
				row.Feedback = fmt.Sprintf("Output should not contain %q", w)
				break
			}
		}
		row.Got = short(r.Stdout)
	default:
		n := lineCount(r.Stdout)
		row.Passed = n > *t.StdoutLinesGT
		row.Got, row.Expected = fmt.Sprintf("%d lines", n), fmt.Sprintf("more than %d lines", *t.StdoutLinesGT)
	}
	return clearIfPassed(row)
}

func httpCheck(t HTTPTest, r HTTPResult, vars map[string]string) exercise.CheckResult {
	row := exercise.CheckResult{Visible: true}
	if r.Error != "" {
		row.Feedback = "No response: " + short(r.Error)
		return row
	}
	switch {
	case t.StatusCode != nil:
		row.Passed = r.Status == *t.StatusCode
		row.Got, row.Expected = strconv.Itoa(r.Status), strconv.Itoa(*t.StatusCode)
	case t.BodyContains != nil:
		w := Interpolate(*t.BodyContains, vars)
		row.Passed = strings.Contains(r.Body, w)
		row.Got, row.Expected = short(r.Body), w
	case t.JSONBody != nil:
		want := InterpolateJSON(normalize(t.JSONBody.Value), vars)
		var body any
		if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
			row.Feedback = "The response body is not JSON"
			row.Got = short(r.Body)
			return row
		}
		got, ok := lookup(body, t.JSONBody.Path)
		if !ok {
			row.Feedback = fmt.Sprintf("The JSON body has no %s", t.JSONBody.Path)
			return row
		}
		row.Passed = compare(t.JSONBody.Op, got, want)
		row.Got, row.Expected = repr(got), repr(want)
	default:
		w := Interpolate(t.Header.Contains, vars)
		var got string
		found := false
		for k, v := range r.Headers {
			if strings.EqualFold(k, t.Header.Name) {
				got, found = v, true
				break
			}
		}
		row.Passed = found && strings.Contains(got, w)
		if !found {
			row.Feedback = fmt.Sprintf("No %s header", t.Header.Name)
		}
		row.Got, row.Expected = short(got), w
	}
	return clearIfPassed(row)
}

// clearIfPassed drops Got/Expected on passing rows (they are for failures).
func clearIfPassed(r exercise.CheckResult) exercise.CheckResult {
	if r.Passed {
		r.Got, r.Expected, r.Feedback = "", "", ""
	}
	return r
}
