package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// hardGrace is added to the harness's own budget (two phases of TimeoutMs)
// before the container is killed from outside.
const hardGrace = 5 * time.Second

// Docker runs jobs with the docker CLI.
//
// Every container gets: no network, a read-only root filesystem with small
// tmpfs mounts for /work and /tmp, a non-root user, no capabilities, no new
// privileges, memory (and swap) and pids limits, one CPU, and a hard wall-clock
// kill. Stdout and stderr are capped.
type Docker struct {
	Binary string // docker CLI path; "docker" when empty
}

// Args returns the docker run arguments for a job (exported for tests).
func (d Docker) Args(name string, job Job) []string {
	l := job.Limits.WithDefaults()
	mem := strconv.Itoa(l.MemoryMB) + "m"
	return []string{
		"run", "--rm", "-i", "--name", name,
		"--network", "none",
		"--read-only",
		"--tmpfs", "/work:rw,nosuid,nodev,size=16m,uid=10001,gid=10001,mode=0700",
		"--tmpfs", "/tmp:rw,nosuid,nodev,noexec,size=8m",
		"--user", "10001:10001",
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
		"--memory", mem, "--memory-swap", mem,
		"--pids-limit", strconv.Itoa(l.Pids),
		"--cpus", "1",
		"--ulimit", "nofile=256:256",
		"--ulimit", "fsize=16777216:16777216",
		"--log-driver", "none",
		"--label", "codequest.runner=1",
		job.Image,
	}
}

// Run starts the container, feeds job.Stdin and waits for it to exit or for
// the hard limit, whichever comes first.
func (d Docker) Run(ctx context.Context, job Job) (Output, error) {
	if job.Image == "" {
		return Output{}, errors.New("runner: job has no image")
	}
	bin := d.Binary
	if bin == "" {
		bin = "docker"
	}
	l := job.Limits.WithDefaults()
	hard := 2*time.Duration(l.TimeoutMs)*time.Millisecond + hardGrace
	runCtx, cancel := context.WithTimeout(ctx, hard)
	defer cancel()

	name := "cq-run-" + uuid.NewString()
	cmd := exec.CommandContext(runCtx, bin, d.Args(name, job)...)
	cmd.Stdin = bytes.NewReader(job.Stdin)
	stdout := &capped{limit: l.OutputBytes}
	stderr := &capped{limit: stderrBytes}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	// Killing the CLI does not stop the container; remove it by name.
	cmd.Cancel = func() error {
		rm, cancelRm := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelRm()
		_ = exec.CommandContext(rm, bin, "rm", "-f", name).Run()
		return cmd.Process.Kill()
	}
	cmd.WaitDelay = 2 * time.Second

	start := time.Now()
	err := cmd.Run()
	out := Output{
		Stdout: stdout.String(), Stderr: stderr.String(),
		RuntimeMs: int(time.Since(start).Milliseconds()),
	}
	if runCtx.Err() != nil {
		if ctx.Err() != nil {
			return out, fmt.Errorf("runner: %w", ctx.Err())
		}
		out.TimedOut = true
		out.ExitCode = -1
		return out, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		out.ExitCode = exitErr.ExitCode()
		// 125-127 are docker's own failures (bad image, cannot start), not the code's.
		if out.ExitCode >= 125 && out.ExitCode <= 127 {
			return out, fmt.Errorf("%w: docker exited %d: %s", ErrUnavailable, out.ExitCode, lastLine(out.Stderr))
		}
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("%w: %w", ErrUnavailable, err)
	}
	return out, nil
}

// capped is an io.Writer that keeps the first limit bytes and drops the rest.
type capped struct {
	buf   bytes.Buffer
	limit int
}

func (c *capped) Write(p []byte) (int, error) {
	if room := c.limit - c.buf.Len(); room > 0 {
		c.buf.Write(p[:min(room, len(p))])
	}
	return len(p), nil
}

func (c *capped) String() string { return c.buf.String() }

func lastLine(s string) string {
	b := bytes.TrimSpace([]byte(s))
	if i := bytes.LastIndexByte(b, '\n'); i >= 0 {
		b = b[i+1:]
	}
	return string(b)
}
