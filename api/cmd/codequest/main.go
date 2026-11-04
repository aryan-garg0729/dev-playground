// Command codequest is the learner CLI: it signs in to CodeQuest, runs
// command-line lessons (project_cli) on this machine and submits their
// output for grading. `cq` is the separate content-authoring tool.
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
)

var version = "dev" // set with -ldflags "-X main.version=..."

const usage = `codequest — run and submit CodeQuest command-line lessons

Usage:
  codequest login                 sign in (paste a code from the website)
  codequest logout                sign out and revoke this machine's token
  codequest status                your XP, level, streak and league
  codequest run [lesson-id]       run a lesson's steps here; no grading
  codequest submit [lesson-id]    run the steps and submit the output for grading
  codequest config [key [value]]  show or set base_url, api_url, app_url (--reset to clear)
  codequest version

With no lesson id, run and submit use your next unfinished command-line lesson.

Flags for run and submit:
  -y          do not ask before running the lesson's commands
  -s          (run) submit after running, same as submit
  --no-color  plain output (also NO_COLOR=1)
`

// app holds the CLI's I/O so commands can be tested.
type app struct {
	cfg  *config
	out  io.Writer
	err  io.Writer
	in   *bufio.Reader
	open func(url string) error // opens the browser; nil to skip
}

func main() { os.Exit(realMain()) }

func realMain() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cfg, err := loadConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "codequest:", err)
		return 1
	}
	a := &app{cfg: cfg, out: os.Stdout, err: os.Stderr, in: bufio.NewReader(os.Stdin), open: openBrowser}
	return a.run(ctx, os.Args[1:])
}

// run executes one command and returns the process exit code: 0 ok (or a
// passing submit), 1 an error or a failing submit, 2 bad usage.
func (a *app) run(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.out, usage)
		return 2
	}
	var err error
	switch cmd, rest := args[0], args[1:]; cmd {
	case "login", "signin":
		err = a.login(ctx)
	case "logout":
		err = a.logout(ctx)
	case "status", "whoami":
		err = a.status(ctx)
	case "run":
		err = a.lesson(ctx, rest, false)
	case "submit":
		err = a.lesson(ctx, rest, true)
	case "config":
		err = a.config(rest)
	case "version", "--version":
		_, _ = fmt.Fprintln(a.out, "codequest", version)
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, usage)
	default:
		_, _ = fmt.Fprintf(a.err, "codequest: unknown command %q\n\n%s", cmd, usage)
		return 2
	}
	var u usageError
	switch {
	case err == nil:
		return 0
	case errors.Is(err, errFailed):
		return 1
	case errors.As(err, &u):
		_, _ = fmt.Fprintf(a.err, "codequest: %v\n", err)
		return 2
	default:
		_, _ = fmt.Fprintf(a.err, "codequest: %v\n", err)
		return 1
	}
}

// errFailed is a submit the server graded as not passing (already printed).
var errFailed = errors.New("lesson not passed")

type usageError struct{ msg string }

func (e usageError) Error() string { return e.msg }

// parseArgs parses flags placed before or after positional arguments.
func parseArgs(fs *flag.FlagSet, args []string) ([]string, error) {
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, usageError{err.Error()}
		}
		if fs.NArg() == 0 {
			return pos, nil
		}
		pos = append(pos, fs.Arg(0))
		args = fs.Args()[1:]
	}
}

// ask prints a question and reads one line.
func (a *app) ask(question string) (string, error) {
	_, _ = fmt.Fprint(a.out, question)
	line, err := a.in.ReadString('\n')
	if err != nil && (line == "" || !errors.Is(err, io.EOF)) {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func (a *app) config(args []string) error {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.SetOutput(a.err)
	reset := fs.Bool("reset", false, "clear the value")
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	c := a.cfg
	fields := map[string]*string{"base_url": &c.BaseURL, "api_url": &c.APIURL, "app_url": &c.AppURL}
	if len(pos) == 0 {
		_, _ = fmt.Fprintf(a.out, "config file: %s\napi_url:  %s\napp_url:  %s\nbase_url: %s\nsigned in: %v\n",
			c.path, c.apiURL(), c.appURL(), orNone(c.BaseURL), c.Token != "")
		return nil
	}
	f, ok := fields[pos[0]]
	if !ok {
		return usageError{fmt.Sprintf("unknown config key %q (base_url, api_url, app_url)", pos[0])}
	}
	switch {
	case *reset && len(pos) == 1:
		*f = ""
	case !*reset && len(pos) == 2:
		v := strings.TrimSuffix(pos[1], "/")
		if pos[0] == "base_url" {
			if !strings.HasPrefix(v, "http://") && !strings.HasPrefix(v, "https://") {
				return usageError{"base_url must start with http:// or https://"}
			}
		} else if err := checkServiceURL(v); err != nil {
			return err
		}
		if pos[0] == "api_url" && v != c.APIURL && c.Token != "" {
			c.Token, c.Name, c.Email = "", "", ""
			_, _ = fmt.Fprintln(a.out, "Signed out: a token is valid only for the API that issued it.")
		}
		*f = v
	case len(pos) == 1:
		_, _ = fmt.Fprintln(a.out, orNone(*f))
		return nil
	default:
		return usageError{"usage: codequest config <key> <value> | codequest config <key> --reset"}
	}
	if err := c.save(); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(a.out, "%s = %s\n", pos[0], orNone(*f))
	return nil
}

func orNone(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}
