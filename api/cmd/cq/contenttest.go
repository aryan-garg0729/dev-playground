package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"

	"github.com/aryangarg/codequest/api/internal/content"
	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/runner"
)

// testCase is one answer to grade and what it must produce.
type testCase struct {
	ex       exercise.Exercise
	typ      exercise.Type
	label    string // "solution" or the wrong answer's name
	answer   json.RawMessage
	wantPass bool
}

type caseResult struct {
	tc  testCase
	res exercise.Result
	err error
}

// contentTest grades every exercise's reference answer (must pass) and each
// wrong answer (must not pass). It prints one line per case and returns the
// number of failures.
func contentTest(ctx context.Context, b *content.Bundle, run runner.Runner, out io.Writer) int {
	var cases []testCase
	failed := 0
	for _, c := range b.Courses {
		for _, ch := range c.Chapters {
			for _, comp := range ch.Exercises {
				ex := comp.Exercise
				typ, _ := exercise.Get(ex.Type) // Load already checked the type
				t, ok := typ.(exercise.Testable)
				if !ok {
					continue
				}
				ref, err := t.ReferenceAnswer(ex.Spec)
				if err != nil {
					sayf(out, "FAIL %s: reference answer: %v\n", ex.ID, err)
					failed++
					continue
				}
				cases = append(cases, testCase{ex: ex, typ: typ, label: "solution", answer: ref, wantPass: true})
				wrong, err := t.WrongAnswers(ex.Spec)
				if err != nil {
					sayf(out, "FAIL %s: wrong answers: %v\n", ex.ID, err)
					failed++
					continue
				}
				names := make([]string, 0, len(wrong))
				for n := range wrong {
					names = append(names, n)
				}
				slices.Sort(names)
				for _, n := range names {
					cases = append(cases, testCase{ex: ex, typ: typ, label: n, answer: wrong[n]})
				}
			}
		}
	}

	ctx = runner.NewContext(ctx, run)
	results := make([]caseResult, len(cases))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for i, tc := range cases {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			sub := exercise.Submission{ExerciseID: tc.ex.ID, Version: tc.ex.Version, Answer: tc.answer}
			res, err := tc.typ.Grade(ctx, tc.ex, sub)
			results[i] = caseResult{tc: tc, res: res, err: err}
		})
	}
	wg.Wait()

	for _, r := range results {
		id := fmt.Sprintf("%s [%s]", r.tc.ex.ID, r.tc.label)
		switch {
		case r.err != nil:
			sayf(out, "FAIL %s: %v\n", id, r.err)
			failed++
		case r.tc.wantPass && r.res.Status != exercise.StatusPass:
			sayf(out, "FAIL %s: reference answer got %s\n%s", id, r.res.Status, failing(r.res))
			failed++
		case !r.tc.wantPass && r.res.Status == exercise.StatusPass:
			sayf(out, "FAIL %s: wrong answer passed every check\n", id)
			failed++
		case countMismatch(r) != "":
			sayf(out, "FAIL %s: %s\n", id, countMismatch(r))
			failed++
		default:
			sayf(out, "ok   %s: %s\n", id, r.res.Status)
		}
	}
	return failed
}

// countMismatch checks a public view's declared check counts (a boss's HP)
// against the rows grading produced.
func countMismatch(r caseResult) string {
	pv, err := r.tc.typ.PublicView(r.tc.ex)
	if err != nil {
		return "public view: " + err.Error()
	}
	var v struct {
		Checks *struct{ Visible, Hidden int } `json:"check_counts"`
	}
	if json.Unmarshal(pv, &v) != nil || v.Checks == nil {
		return ""
	}
	vis, hid := 0, 0
	for _, c := range r.res.Checks {
		if c.Visible {
			vis++
		} else {
			hid++
		}
	}
	if vis != v.Checks.Visible || hid != v.Checks.Hidden {
		return fmt.Sprintf("public view declares %d visible + %d hidden checks, grading produced %d + %d",
			v.Checks.Visible, v.Checks.Hidden, vis, hid)
	}
	return ""
}

func failing(r exercise.Result) string {
	var sb strings.Builder
	for _, c := range r.Checks {
		if !c.Passed {
			fmt.Fprintf(&sb, "       - %s: %s", c.Name, c.Feedback)
			if c.Got != "" || c.Expected != "" {
				fmt.Fprintf(&sb, " (got %s, expected %s)", c.Got, c.Expected)
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
