package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/aryangarg/codequest/api/internal/exercise/projectcli"
)

// lesson is `run` (submit=false) and `submit`.
func (a *app) lesson(ctx context.Context, args []string, submit bool) error {
	name := "run"
	if submit {
		name = "submit"
	}
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(a.err)
	yes := fs.Bool("y", false, "do not ask before running commands")
	noColor := fs.Bool("no-color", false, "plain output")
	var thenSubmit *bool
	if !submit {
		thenSubmit = fs.Bool("s", false, "submit after running")
	}
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	if len(pos) > 1 {
		return usageError{"usage: codequest " + name + " [lesson-id]"}
	}
	if thenSubmit != nil && *thenSubmit {
		submit = true
	}
	u := newUI(a.out, *noColor)
	cl, err := a.client()
	if err != nil {
		return err
	}

	id := ""
	if len(pos) == 1 {
		id = pos[0]
	} else {
		var next nextResponse
		if err := cl.do(ctx, "GET", "/cli/next", nil, &next); err != nil {
			if isAPIError(err, "not_found") {
				return errors.New("you have no command-line lessons left; pass a lesson id to redo one")
			}
			return err
		}
		id = next.ExerciseID
		u.printf("Next lesson: %s %s\n", u.bold(fmt.Sprintf("%d.%d %s", next.ChapterNo, next.LessonNo, next.Title)), u.dim("("+next.CourseTitle+" · "+id+")"))
	}
	var ex exerciseResponse
	if err := cl.do(ctx, "GET", "/exercises/"+id, nil, &ex); err != nil {
		if isAPIError(err, "not_found") {
			return fmt.Errorf("no lesson %q", id)
		}
		return err
	}
	if ex.Type != projectcli.TypeName {
		return fmt.Errorf("%s is a %s lesson; solve it in the browser", id, ex.Type)
	}
	var spec projectcli.Spec
	if err := json.Unmarshal(ex.PublicView, &spec); err != nil || len(spec.Steps) == 0 {
		return fmt.Errorf("this CLI cannot read lesson %s; update codequest", id)
	}
	if len(spec.AllowedOS) > 0 && !slices.Contains(spec.AllowedOS, runtime.GOOS) {
		return fmt.Errorf("this lesson runs on %s, not %s (on Windows, use WSL)", strings.Join(spec.AllowedOS, ", "), runtime.GOOS)
	}
	if len(pos) == 1 {
		u.printf("Lesson: %s %s\n", u.bold(ex.Title), u.dim("("+ex.Chapter.Title+" · "+id+")"))
	}
	baseURL := a.cfg.BaseURL
	if spec.BaseURL != "" && baseURL != "" {
		u.printf("%s\n", u.dim("Using base_url "+baseURL+" (reset with `codequest config base_url --reset`)"))
	}

	if !*yes {
		ok, err := a.confirm(u, spec)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("canceled")
		}
	}

	sr := newStepRunner()
	sr.onStep = func(i int, st projectcli.Step, vars map[string]string) { u.stepHeading(i, st, vars) }
	sr.onResult = func(_ int, st projectcli.Step, r projectcli.StepResult, vars map[string]string) {
		u.stepOutput(st, r, vars)
	}
	answer := sr.run(ctx, spec, baseURL)
	if !submit {
		u.printf("\nRun %s to grade it.\n", u.bold("codequest submit "+id))
		return nil
	}

	u.printf("\nSubmitting…\n")
	var res submitResponse
	err = cl.do(ctx, "POST", "/exercises/"+id+"/submit", map[string]any{"exercise_id": id, "version": ex.Version, "answer": answer}, &res)
	if isAPIError(err, "stale_version") {
		return errors.New("the lesson was just updated; run the command again")
	}
	if err != nil {
		return err
	}
	for wait := 0; res.pending(); wait++ {
		if wait >= 60 {
			return fmt.Errorf("grading is taking long; check the result in the browser (submission %s)", res.SubmissionID)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
		if err := cl.do(ctx, "GET", "/submissions/"+res.SubmissionID, nil, &res); err != nil {
			return err
		}
	}
	var gr gradeResult
	if err := json.Unmarshal(res.Result, &gr); err != nil {
		return fmt.Errorf("unexpected result: %w", err)
	}
	u.results(gr)
	u.outcome(res.Outcome)
	if gr.Status != "pass" {
		return errFailed
	}
	return nil
}

// confirm lists what the lesson will run and asks first: lesson commands run
// in the learner's own shell.
func (a *app) confirm(u *ui, s projectcli.Spec) (bool, error) {
	u.printf("\nThis lesson runs on your machine, in %s:\n", u.bold(cwd()))
	for i, st := range s.Steps {
		if st.CLICommand != nil {
			u.printf("  %d. %s %s\n", i+1, u.dim("$"), st.CLICommand.Command)
		} else {
			u.printf("  %d. %s %s %s\n", i+1, u.dim(">"), st.HTTPRequest.Request.Method, st.HTTPRequest.Request.URL)
		}
	}
	ans, err := a.ask("Continue? [Y/n] ")
	if err != nil {
		return false, errors.New("could not read an answer; pass -y to run without asking")
	}
	switch strings.ToLower(ans) {
	case "", "y", "yes":
		return true, nil
	}
	return false, nil
}

func cwd() string {
	d, err := os.Getwd()
	if err != nil {
		return "the current directory"
	}
	return d
}
