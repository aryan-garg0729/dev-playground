// Package wasm runs learner programs as WASI modules under wazero, in-process.
//
// Each language is an interpreter compiled to WASI (runners/wasm). A job gets a
// fresh module instance per phase with: the learner's entry file as the only
// file (read-only, at /work), no environment, no sockets (WASI preview 1 has
// none), no subprocesses, a memory cap and a deadline. The job's stdin and the
// report this package writes match the Python harness's (checks.Report), so
// exercise types grade both the same way.
//
// The interpreters exit 3 for a compile error with stderr starting
// "<line>:<col>: <message>"; this package turns that into Report.SyntaxError.
package wasm

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing/fstest"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"

	"github.com/aryangarg/codequest/api/internal/checks"
	"github.com/aryangarg/codequest/api/internal/runner"
)

// Languages maps a lesson language to its module file in the module directory.
// Adding a language is a module in runners/wasm plus an entry here.
var Languages = map[string]string{
	"go": "go.wasm",
	"js": "js.wasm",
}

// Job is the stdin of a wasm job; the field names match the Python harness's.
type Job struct {
	Files       map[string]string `json:"files"`
	Entry       string            `json:"entry"`
	RunEntry    bool              `json:"run_entry"`
	Stdin       string            `json:"stdin"` // the program phase's stdin
	Cases       []string          `json:"cases"` // one more program run per stdin
	Tests       []json.RawMessage `json:"tests"`
	Forbidden   []string          `json:"forbidden"`
	TimeoutMs   int               `json:"timeout_ms"`
	OutputLimit int               `json:"output_limit"`
}

const (
	stderrLimit  = 16 << 10
	maxCases     = 50
	syntaxExit   = 3
	pagesPerMB   = 16 // 64 KiB wasm pages
	workDir      = "/work"
	maxJobStdin  = 4 << 20
	defaultLimit = 16 << 10
)

var syntaxRe = regexp.MustCompile(`^(\d+):(\d+): (.+)`)

// Runner implements runner.Runner for "wasm:<language>" images.
type Runner struct {
	dir   string
	cache wazero.CompilationCache
	bins  map[string][]byte // language -> module bytes

	mu  sync.Mutex
	rts map[int]*instanceSet // by memory limit in MB
}

type instanceSet struct {
	rt   wazero.Runtime
	mods map[string]wazero.CompiledModule
}

// New loads every module in Languages from dir and compiles them once (cached
// on disk in cacheDir when it is set, so restarts are fast). It fails if a
// module is missing or does not compile.
func New(ctx context.Context, dir, cacheDir string) (*Runner, error) {
	r := &Runner{dir: dir, bins: map[string][]byte{}, rts: map[int]*instanceSet{}}
	r.cache = wazero.NewCompilationCache()
	if cacheDir != "" {
		c, err := wazero.NewCompilationCacheWithDir(cacheDir)
		if err != nil {
			return nil, fmt.Errorf("wasm: cache dir: %w", err)
		}
		r.cache = c
	}
	for lang, file := range Languages {
		b, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return nil, fmt.Errorf("wasm: %s module: %w (build it with `make wasm-modules`)", lang, err)
		}
		r.bins[lang] = b
	}
	if _, err := r.set(ctx, runner.DefaultMemoryMB); err != nil {
		return nil, err
	}
	return r, nil
}

// Close releases every runtime.
func (r *Runner) Close(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var errs []error
	for _, s := range r.rts {
		errs = append(errs, s.rt.Close(ctx))
	}
	errs = append(errs, r.cache.Close(ctx))
	return errors.Join(errs...)
}

// set returns the runtime whose memory cap is mb, compiling the modules into it
// on first use (from the compilation cache after the first runtime).
func (r *Runner) set(ctx context.Context, mb int) (*instanceSet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.rts[mb]; ok {
		return s, nil
	}
	cfg := wazero.NewRuntimeConfig().
		WithCompilationCache(r.cache).
		WithCloseOnContextDone(true).
		WithMemoryLimitPages(uint32(mb * pagesPerMB))
	rt := wazero.NewRuntimeWithConfig(ctx, cfg)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil {
		_ = rt.Close(ctx)
		return nil, fmt.Errorf("wasm: wasi: %w", err)
	}
	s := &instanceSet{rt: rt, mods: map[string]wazero.CompiledModule{}}
	for lang, b := range r.bins {
		m, err := rt.CompileModule(ctx, b)
		if err != nil {
			_ = rt.Close(ctx)
			return nil, fmt.Errorf("wasm: compile %s: %w", lang, err)
		}
		s.mods[lang] = m
	}
	r.rts[mb] = s
	return s, nil
}

// Run implements runner.Runner. Output.Stdout is one checks.Report JSON line.
func (r *Runner) Run(ctx context.Context, job runner.Job) (runner.Output, error) {
	lang, ok := strings.CutPrefix(job.Image, runner.WasmPrefix)
	if !ok {
		return runner.Output{}, fmt.Errorf("wasm: %q is not a wasm image", job.Image)
	}
	if _, ok := r.bins[lang]; !ok {
		return runner.Output{}, fmt.Errorf("wasm: unknown language %q", lang)
	}
	if len(job.Stdin) > maxJobStdin {
		return runner.Output{}, errors.New("wasm: job too large")
	}
	var j Job
	if err := json.Unmarshal(job.Stdin, &j); err != nil {
		return runner.Output{}, fmt.Errorf("wasm: decode job: %w", err)
	}
	src, ok := j.Files[j.Entry]
	switch {
	case !ok || len(j.Files) != 1:
		return runner.Output{}, errors.New("wasm: files must hold exactly the entry")
	case strings.ContainsAny(j.Entry, "/\\") || j.Entry == "" || j.Entry[0] == '.':
		return runner.Output{}, fmt.Errorf("wasm: bad entry %q", j.Entry)
	case len(j.Tests) > 0 || len(j.Forbidden) > 0:
		return runner.Output{}, errors.New("wasm: unit tests and forbidden calls are Python-only")
	case len(j.Cases) > maxCases:
		return runner.Output{}, fmt.Errorf("wasm: at most %d cases", maxCases)
	}
	l := job.Limits.WithDefaults()
	if j.TimeoutMs <= 0 || j.TimeoutMs > l.TimeoutMs {
		j.TimeoutMs = l.TimeoutMs
	}
	if j.OutputLimit <= 0 || j.OutputLimit > l.OutputBytes {
		j.OutputLimit = defaultLimit
	}
	set, err := r.set(ctx, l.MemoryMB)
	if err != nil {
		return runner.Output{}, err
	}
	p := phaseRunner{mod: set.mods[lang], rt: set.rt, lang: lang, entry: j.Entry, src: []byte(src),
		timeout: time.Duration(j.TimeoutMs) * time.Millisecond, outLimit: j.OutputLimit, memoryMB: l.MemoryMB}

	// Every phase shares one budget, the same 2 x timeout the queue waits for
	// (as the Python harness's two phases do), however many cases there are.
	jobCtx, cancel := context.WithTimeout(ctx, 2*p.timeout)
	defer cancel()
	start := time.Now()
	rep := checks.Report{Version: 1, Forbidden: []checks.Forbidden{}}
	stdins := []string{}
	if j.RunEntry {
		stdins = append(stdins, j.Stdin)
	}
	stdins = append(stdins, j.Cases...)
	for i, in := range stdins {
		ph, err := p.run(jobCtx, in)
		if err != nil {
			return runner.Output{}, err
		}
		if se := syntaxError(j.Entry, ph); se != nil {
			// The same source fails to compile on every run: report it once.
			rep.SyntaxError = se
			break
		}
		rep.TimedOut = rep.TimedOut || ph.TimedOut
		rep.RuntimeMs += ph.RuntimeMs
		if j.RunEntry && i == 0 {
			rep.Program = ph
		} else {
			rep.Cases = append(rep.Cases, *ph)
		}
	}
	b, err := json.Marshal(rep)
	if err != nil {
		return runner.Output{}, fmt.Errorf("wasm: encode report: %w", err)
	}
	return runner.Output{Stdout: string(b), RuntimeMs: int(time.Since(start).Milliseconds())}, nil
}

type phaseRunner struct {
	mod      wazero.CompiledModule
	rt       wazero.Runtime
	lang     string
	entry    string
	src      []byte
	timeout  time.Duration
	outLimit int
	memoryMB int
}

// run starts one fresh instance of the interpreter on the entry file.
func (p phaseRunner) run(ctx context.Context, stdin string) (*checks.Phase, error) {
	stdout, stderr := &capped{max: p.outLimit}, &capped{max: stderrLimit}
	fsys := fstest.MapFS{p.entry: {Data: slices.Clone(p.src), Mode: 0o444}}
	cfg := wazero.NewModuleConfig().
		WithName(""). // anonymous: instances never clash
		WithArgs(p.lang, workDir+"/"+p.entry).
		WithStdin(strings.NewReader(stdin)).
		WithStdout(stdout).
		WithStderr(stderr).
		WithFSConfig(wazero.NewFSConfig().WithFSMount(fsys, workDir)).
		WithSysWalltime().
		WithSysNanotime().
		WithSysNanosleep().
		WithRandSource(rand.Reader)
	if ctx.Err() != nil { // the job's budget is spent
		return &checks.Phase{TimedOut: true, ExitCode: 124}, nil
	}
	runCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	start := time.Now()
	mod, err := p.rt.InstantiateModule(runCtx, p.mod, cfg)
	ph := &checks.Phase{RuntimeMs: int(time.Since(start).Milliseconds())}
	if mod != nil {
		_ = mod.Close(context.WithoutCancel(ctx))
	}
	var exit *sys.ExitError
	switch {
	case err == nil:
	case errors.As(err, &exit) && exit.ExitCode() == sys.ExitCodeDeadlineExceeded:
		ph.TimedOut, ph.ExitCode = true, 124
	case errors.As(err, &exit) && exit.ExitCode() == sys.ExitCodeContextCanceled:
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			ph.TimedOut, ph.ExitCode = true, 124
			break
		}
		return nil, fmt.Errorf("wasm: run canceled: %w", ctx.Err())
	case errors.As(err, &exit):
		ph.ExitCode = int(exit.ExitCode())
	default:
		// A trap (for example a failed memory grow) ends the instance.
		ph.ExitCode = 1
		stderr.WriteString("runtime error: " + trapMessage(err) + "\n")
	}
	ph.Stdout, ph.Stderr = stdout.String(), stderr.String()
	if strings.Contains(ph.Stderr, "fatal error: out of memory") || strings.HasPrefix(ph.Stderr, "runtime error: out of memory") {
		ph.Stderr = fmt.Sprintf("Memory limit: your program used more than %d MB\n", p.memoryMB)
	}
	return ph, nil
}

func trapMessage(err error) string {
	msg := err.Error()
	if strings.Contains(msg, "out of memory") {
		return "out of memory"
	}
	if i := strings.IndexByte(msg, '\n'); i > 0 {
		msg = msg[:i]
	}
	return msg
}

func syntaxError(entry string, ph *checks.Phase) *checks.SyntaxError {
	if ph.ExitCode != syntaxExit {
		return nil
	}
	m := syntaxRe.FindStringSubmatch(ph.Stderr)
	if m == nil {
		return nil
	}
	line, _ := strconv.Atoi(m[1])
	return &checks.SyntaxError{File: entry, Line: line, Message: strings.TrimSpace(m[3])}
}

// capped keeps the first max bytes written and silently drops the rest, so a
// flood of output cannot grow memory. Writes never fail.
type capped struct {
	buf bytes.Buffer
	max int
}

func (c *capped) Write(p []byte) (int, error) {
	if room := c.max - c.buf.Len(); room > 0 {
		c.buf.Write(p[:min(len(p), room)])
	}
	return len(p), nil
}

func (c *capped) WriteString(s string) { _, _ = c.Write([]byte(s)) }

func (c *capped) String() string { return c.buf.String() }
