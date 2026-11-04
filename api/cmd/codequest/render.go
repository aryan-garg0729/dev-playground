package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aryangarg/codequest/api/internal/exercise/projectcli"
)

// ui writes plain text, with ANSI colors when out is a terminal and neither
// NO_COLOR nor --no-color is set.
type ui struct {
	out   io.Writer
	color bool
}

func newUI(out io.Writer, noColor bool) *ui {
	color := !noColor && os.Getenv("NO_COLOR") == ""
	if f, ok := out.(*os.File); !ok {
		color = false
	} else if fi, err := f.Stat(); err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		color = false
	}
	return &ui{out: out, color: color}
}

func (u *ui) paint(code, s string) string {
	if !u.color {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func (u *ui) bold(s string) string  { return u.paint("1", s) }
func (u *ui) dim(s string) string   { return u.paint("2", s) }
func (u *ui) green(s string) string { return u.paint("32", s) }
func (u *ui) red(s string) string   { return u.paint("31", s) }
func (u *ui) cyan(s string) string  { return u.paint("36", s) }

func (u *ui) printf(format string, a ...any) { _, _ = fmt.Fprintf(u.out, format, a...) }

// indent prefixes every line of s.
func indent(s, prefix string) string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return ""
	}
	return prefix + strings.ReplaceAll(s, "\n", "\n"+prefix) + "\n"
}

// stepHeading prints "1. Description" and the command or request.
func (u *ui) stepHeading(i int, st projectcli.Step, vars map[string]string) {
	u.printf("\n%s %s\n", u.bold(fmt.Sprintf("%d.", i+1)), u.bold(st.Description))
	if st.CLICommand != nil {
		u.printf("   %s %s\n", u.dim("$"), u.cyan(projectcli.Interpolate(st.CLICommand.Command, vars)))
		return
	}
	r := st.HTTPRequest.Request
	u.printf("   %s %s %s\n", u.dim(">"), u.cyan(r.Method), u.cyan(projectcli.Interpolate(r.URL, vars)))
}

// stepOutput prints a step's raw output and the tests it will be judged on,
// with no verdict (grading happens on the server).
func (u *ui) stepOutput(st projectcli.Step, r projectcli.StepResult, vars map[string]string) {
	const pad = "   │ "
	switch {
	case r.Skipped:
		u.printf("   %s\n", u.dim("skipped: it needs ${"+strings.Join(st.Missing(vars), "}, ${")+"} from an earlier step"))
	case r.CLI != nil:
		u.printf("%s", u.dim(indent(r.CLI.Stdout, pad)))
		u.printf("   %s\n", u.dim(fmt.Sprintf("exit code %d", r.CLI.ExitCode)))
	case r.HTTP != nil && r.HTTP.Error != "":
		u.printf("   %s\n", u.red("no response: "+r.HTTP.Error))
	case r.HTTP != nil:
		u.printf("   %s\n", u.dim(fmt.Sprintf("< %d", r.HTTP.Status)))
		body := r.HTTP.Body
		if len(body) > 2000 {
			body = clip(body, 2000) + "\n…"
		}
		u.printf("%s", u.dim(indent(body, pad)))
	}
	var tests []string
	if st.CLICommand != nil {
		for _, t := range st.CLICommand.Tests {
			tests = append(tests, projectcli.DescribeCLITest(t, vars))
		}
	} else {
		for _, t := range st.HTTPRequest.Tests {
			tests = append(tests, projectcli.DescribeHTTPTest(t, vars))
		}
	}
	for _, t := range tests {
		u.printf("   %s %s\n", u.dim("•"), t)
	}
}

// results prints the server's verdict.
func (u *ui) results(r gradeResult) {
	u.printf("\n")
	for _, c := range r.Checks {
		if c.Passed {
			u.printf("%s %s\n", u.green("✓"), c.Name)
			continue
		}
		u.printf("%s %s\n", u.red("✗"), c.Name)
		if c.Feedback != "" {
			u.printf("    %s\n", c.Feedback)
		}
		if c.Expected != "" {
			u.printf("    expected: %s\n", c.Expected)
		}
		if c.Got != "" {
			u.printf("    got:      %s\n", strings.ReplaceAll(strings.TrimRight(c.Got, "\n"), "\n", "\n              "))
		}
	}
	passed := 0
	for _, c := range r.Checks {
		if c.Passed {
			passed++
		}
	}
	u.printf("\n")
	switch r.Status {
	case "pass":
		u.printf("%s %d/%d checks passed\n", u.green(u.bold("PASS")), passed, len(r.Checks))
	case "fail":
		u.printf("%s %d/%d checks passed\n", u.red(u.bold("FAIL")), passed, len(r.Checks))
	default:
		u.printf("%s grading did not finish (%s); this attempt does not count\n", u.red(u.bold("ERROR")), r.Status)
	}
}

func (u *ui) outcome(o *outcome) {
	if o == nil {
		return
	}
	if o.XPAwarded > 0 {
		u.printf("%s  %s\n", u.green(fmt.Sprintf("+%d XP", o.XPAwarded)), u.dim(fmt.Sprintf("(total %d)", o.TotalXP)))
	}
	if o.LevelUp {
		u.printf("Level up! You are now level %d\n", o.Level)
	}
	if o.StreakExtended {
		u.printf("Streak: %d days\n", o.StreakDays)
	}
	if o.Rewards == nil {
		return
	}
	for _, b := range o.Rewards.Bytes {
		u.printf("+%d Bytes  %s\n", b.Amount, u.dim(b.Label))
	}
	for _, m := range o.Rewards.Medals {
		u.printf("Medal earned: %s (%s)\n", m.Name, m.Tier)
	}
	for _, q := range o.Rewards.QuestsDone {
		u.printf("Quest complete: %s\n", q)
	}
	if o.Rewards.ChapterComplete {
		u.printf("Chapter complete!\n")
	}
}
