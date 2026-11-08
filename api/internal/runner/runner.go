// Package runner executes untrusted code for server-side grading.
//
// Two implementations run code: Docker (one hardened container per job) and
// wasm (a WASI module per job, package runner/wasm). The queue runner forwards
// jobs to cmd/worker over Redis, on one queue per kind. Exercise types get the runner from the
// context (NewContext / FromContext), so plugins stay free of global state.
package runner

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Limits bound one job. Zero values take the defaults below.
type Limits struct {
	TimeoutMs   int `json:"timeout_ms"`   // per harness phase; the hard kill adds a grace period
	MemoryMB    int `json:"memory_mb"`    // memory and memory+swap
	Pids        int `json:"pids"`         // max processes in the container
	OutputBytes int `json:"output_bytes"` // stdout kept; stderr keeps at most 16 KB
}

// Job is one sandboxed run. Files reach the container on stdin (the harness
// writes them to its tmpfs), so the root filesystem stays read-only and no
// host paths are mounted.
type Job struct {
	Image  string `json:"image"`
	Stdin  []byte `json:"stdin"`
	Limits Limits `json:"limits"`
}

// Output is what the container produced.
type Output struct {
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ExitCode  int    `json:"exit_code"`
	RuntimeMs int    `json:"runtime_ms"`
	TimedOut  bool   `json:"timed_out"` // killed by the hard wall-clock limit
}

// Runner runs jobs.
type Runner interface {
	Run(ctx context.Context, job Job) (Output, error)
}

// Defaults and caps.
const (
	DefaultTimeoutMs   = 3000
	DefaultMemoryMB    = 128
	DefaultPids        = 64
	DefaultOutputBytes = 256 << 10
	MaxTimeoutMs       = 10000
	MaxMemoryMB        = 512
	stderrBytes        = 16 << 10
)

// Images maps a manifest runner name to its server image. Client runners map
// to the server runner that must give the same result (pyodide → docker-python).
var Images = map[string]string{
	"docker-python": "codequest/python:3.12",
	"pyodide":       "codequest/python:3.12",
}

// ErrUnavailable means no runner could take the job (no worker, queue full, ...).
var ErrUnavailable = errors.New("runner unavailable")

// WasmPrefix starts the image name of a wasm job: "wasm:<language>".
const WasmPrefix = "wasm:"

// WasmLanguages are the lesson languages the wasm runner grades (package
// runner/wasm maps each to its module; a test keeps the two in step).
var WasmLanguages = []string{"go", "js"}

// IsWasm reports whether image names a wasm language rather than a Docker image.
func IsWasm(image string) bool { return strings.HasPrefix(image, WasmPrefix) }

// WasmImage is the job image for a wasm language.
func WasmImage(language string) string { return WasmPrefix + language }

// Route sends wasm jobs to Wasm and every other job to Docker. A nil side is
// unavailable.
type Route struct {
	Wasm, Docker Runner
}

// Run implements Runner.
func (r Route) Run(ctx context.Context, job Job) (Output, error) {
	next := r.Docker
	if IsWasm(job.Image) {
		next = r.Wasm
	}
	if next == nil {
		return Output{}, fmt.Errorf("%w: no runner for %s", ErrUnavailable, job.Image)
	}
	return next.Run(ctx, job)
}

// WithDefaults fills zero limits with the defaults and applies the caps.
func (l Limits) WithDefaults() Limits {
	if l.TimeoutMs <= 0 {
		l.TimeoutMs = DefaultTimeoutMs
	}
	if l.MemoryMB <= 0 {
		l.MemoryMB = DefaultMemoryMB
	}
	if l.Pids <= 0 {
		l.Pids = DefaultPids
	}
	if l.OutputBytes <= 0 {
		l.OutputBytes = DefaultOutputBytes
	}
	l.TimeoutMs = min(l.TimeoutMs, MaxTimeoutMs)
	l.MemoryMB = min(l.MemoryMB, MaxMemoryMB)
	return l
}

type ctxKey struct{}

// NewContext returns ctx carrying r.
func NewContext(ctx context.Context, r Runner) context.Context {
	return context.WithValue(ctx, ctxKey{}, r)
}

// FromContext returns the runner set by NewContext.
func FromContext(ctx context.Context) (Runner, bool) {
	r, ok := ctx.Value(ctxKey{}).(Runner)
	return r, ok && r != nil
}
