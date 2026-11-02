// Command cq is the content CLI.
//
//	cq validate <dir>   validate a content tree; print every problem as "path: reason"
//	cq build <dir>      validate, then write the compiled bundle to DATABASE_URL
//	cq test <dir>       content CI: every reference answer passes, every wrong/ answer fails
//	cq migrate up|down  apply or roll back (one step) database migrations
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/db"
	"github.com/aryangarg/codequest/api/internal/config"
	"github.com/aryangarg/codequest/api/internal/content"
	_ "github.com/aryangarg/codequest/api/internal/exercise/all" // register exercise types
	"github.com/aryangarg/codequest/api/internal/runner"
)

const usage = `usage:
  cq validate <dir>
  cq build <dir>
  cq test <dir>
  cq migrate up|down`

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 {
		say(stderr, usage)
		return 2
	}
	switch args[0] {
	case "validate":
		b, ok := load(args[1], stderr)
		if !ok {
			return 1
		}
		sayf(stdout, "ok: %s\n", summary(b))
		return 0
	case "build":
		b, ok := load(args[1], stderr)
		if !ok {
			return 1
		}
		dsn, err := config.DatabaseURL()
		if err != nil {
			say(stderr, err)
			return 1
		}
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			say(stderr, "connect database:", err)
			return 1
		}
		defer pool.Close()
		st, err := content.Build(ctx, pool, b)
		if err != nil {
			say(stderr, err)
			return 1
		}
		sayf(stdout, "built: %d course(s), %d chapter(s), %d exercise(s)\n", st.Courses, st.Chapters, st.Exercises)
		return 0
	case "test":
		b, ok := load(args[1], stderr)
		if !ok {
			return 1
		}
		bin := os.Getenv("DOCKER_BIN")
		if bin == "" {
			bin = "docker"
		}
		failed := contentTest(ctx, b, runner.Docker{Binary: bin}, stdout)
		if failed > 0 {
			sayf(stderr, "%d content check(s) failed\n", failed)
			return 1
		}
		sayf(stdout, "ok: %s, every reference passes and every wrong answer fails\n", summary(b))
		return 0
	case "migrate":
		dsn, err := config.DatabaseURL()
		if err != nil {
			say(stderr, err)
			return 1
		}
		switch args[1] {
		case "up":
			err = db.Up(dsn)
		case "down":
			err = db.Down(dsn)
		default:
			say(stderr, usage)
			return 2
		}
		if err != nil {
			say(stderr, err)
			return 1
		}
		sayf(stdout, "migrate %s: ok\n", args[1])
		return 0
	default:
		say(stderr, usage)
		return 2
	}
}

func load(dir string, stderr io.Writer) (*content.Bundle, bool) {
	b, err := content.Load(dir)
	if verr, ok := errors.AsType[*content.ValidationError](err); ok {
		for _, p := range verr.Problems {
			say(stderr, p.String())
		}
		sayf(stderr, "%d problem(s) found\n", len(verr.Problems))
		return nil, false
	}
	if err != nil {
		say(stderr, err)
		return nil, false
	}
	return b, true
}

func summary(b *content.Bundle) string {
	ch, ex := 0, 0
	for _, c := range b.Courses {
		ch += len(c.Chapters)
		for _, h := range c.Chapters {
			ex += len(h.Exercises)
		}
	}
	return fmt.Sprintf("%d course(s), %d chapter(s), %d exercise(s)", len(b.Courses), ch, ex)
}

// say and sayf write CLI output; a failed write to stdout/stderr has nowhere to be reported.
func say(w io.Writer, a ...any) { _, _ = fmt.Fprintln(w, a...) }

func sayf(w io.Writer, format string, a ...any) { _, _ = fmt.Fprintf(w, format, a...) }
