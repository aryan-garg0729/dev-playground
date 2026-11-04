package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aryangarg/codequest/api/internal/exercise/projectcli"
)

// stepRunner executes a lesson's steps on this machine.
type stepRunner struct {
	dir         string        // working directory for commands ("" = current)
	cmdTimeout  time.Duration // per command
	httpTimeout time.Duration // per request
	// onStep is called before each step runs and onResult after (for output).
	onStep   func(i int, st projectcli.Step, vars map[string]string)
	onResult func(i int, st projectcli.Step, r projectcli.StepResult, vars map[string]string)
}

func newStepRunner() *stepRunner {
	return &stepRunner{cmdTimeout: 30 * time.Second, httpTimeout: 10 * time.Second}
}

// run executes every step in order and returns the raw output to report.
// Captures are applied only to interpolate later steps; the server rebuilds
// them from this output on its own.
func (sr *stepRunner) run(ctx context.Context, s projectcli.Spec, baseURL string) projectcli.Answer {
	vars := map[string]string{}
	a := projectcli.Answer{Steps: make([]projectcli.StepResult, len(s.Steps))}
	if s.BaseURL != "" {
		if baseURL == "" {
			baseURL = s.BaseURL
		}
		baseURL = strings.TrimSuffix(baseURL, "/")
		vars[projectcli.BaseURLVar] = baseURL
		a.BaseURL = baseURL
	}
	for i, st := range s.Steps {
		if sr.onStep != nil {
			sr.onStep(i, st, vars)
		}
		var r projectcli.StepResult
		switch {
		case len(st.Missing(vars)) > 0:
			r.Skipped = true
		case st.CLICommand != nil:
			r.CLI = sr.runCommand(ctx, projectcli.Interpolate(st.CLICommand.Command, vars))
		default:
			r.HTTP = sr.runRequest(ctx, st.HTTPRequest.Request, vars)
		}
		projectcli.ApplyCaptures(st, r, vars)
		a.Steps[i] = r
		if sr.onResult != nil {
			sr.onResult(i, st, r, vars)
		}
	}
	return a
}

// clip makes s valid UTF-8 and at most n bytes. Sanitizing first matters:
// JSON encoding would otherwise expand each invalid byte to three.
func clip(s string, n int) string {
	s = strings.ToValidUTF8(s, "�")
	if len(s) <= n {
		return s
	}
	s = s[:n]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// boundedBuffer keeps the first limit bytes and discards the rest.
type boundedBuffer struct {
	buf       bytes.Buffer
	limit     int
	truncated bool
}

func (b *boundedBuffer) Write(p []byte) (int, error) {
	if room := b.limit - b.buf.Len(); room > 0 {
		if len(p) > room {
			b.buf.Write(p[:room])
			b.truncated = true
		} else {
			b.buf.Write(p)
		}
	} else if len(p) > 0 {
		b.truncated = true
	}
	return len(p), nil
}

func shellCommand(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/C", command)
	}
	return exec.CommandContext(ctx, "sh", "-c", command)
}

const truncNote = "\n[output truncated]\n"

// runCommand runs command in the learner's shell. Output is stdout and
// stderr combined (a traceback goes to stderr), capped at MaxOutputBytes.
func (sr *stepRunner) runCommand(ctx context.Context, command string) *projectcli.CLIResult {
	ctx, cancel := context.WithTimeout(ctx, sr.cmdTimeout)
	defer cancel()
	cmd := shellCommand(ctx, command)
	cmd.Dir = sr.dir
	cmd.Stdin = nil
	out := &boundedBuffer{limit: projectcli.MaxOutputBytes - len(truncNote) - 64}
	cmd.Stdout, cmd.Stderr = out, out
	cmd.Env = append(os.Environ(), "CODEQUEST=1")
	// A background child holding the pipes must not hang the CLI.
	cmd.WaitDelay = 2 * time.Second
	err := cmd.Run()
	text := out.buf.String()
	if out.truncated {
		text += truncNote
	}
	code := 0
	var exitErr *exec.ExitError
	switch {
	case ctx.Err() == context.DeadlineExceeded:
		code = -1
		text += "\n[killed: the command ran longer than " + sr.cmdTimeout.String() + "]\n"
	case errors.As(err, &exitErr):
		code = exitErr.ExitCode()
	case err != nil:
		code = 127
		text += err.Error() + "\n"
	}
	return &projectcli.CLIResult{ExitCode: code, Stdout: clip(text, projectcli.MaxOutputBytes)}
}

// runRequest sends a lesson request. It uses its own client, so the
// CodeQuest token is never attached, whatever the URL.
func (sr *stepRunner) runRequest(ctx context.Context, rq projectcli.Request, vars map[string]string) *projectcli.HTTPResult {
	fail := func(err error) *projectcli.HTTPResult {
		return &projectcli.HTTPResult{Error: clip(err.Error(), 1024)}
	}
	var body io.Reader
	if rq.BodyJSON != nil {
		b, err := json.Marshal(projectcli.InterpolateJSON(rq.BodyJSON, vars))
		if err != nil {
			return fail(err)
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, rq.Method, projectcli.Interpolate(rq.URL, vars), body)
	if err != nil {
		return fail(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "codequest-cli/"+version)
	for k, v := range rq.Headers {
		req.Header.Set(k, projectcli.Interpolate(v, vars))
	}
	client := &http.Client{Timeout: sr.httpTimeout}
	if !rq.FollowRedirects {
		client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	}
	resp, err := client.Do(req)
	if err != nil {
		return fail(err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, projectcli.MaxOutputBytes))
	res := &projectcli.HTTPResult{Status: resp.StatusCode, Body: clip(string(b), projectcli.MaxOutputBytes), Headers: map[string]string{}}
	keys := make([]string, 0, len(resp.Header))
	for k := range resp.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	size := 0
	for _, k := range keys {
		v := clip(strings.Join(resp.Header.Values(k), ", "), 512)
		if len(res.Headers) == projectcli.MaxHeaders || size+len(k)+len(v) > projectcli.MaxHeaderBytes {
			break
		}
		res.Headers[k] = v
		size += len(k) + len(v)
	}
	return res
}
