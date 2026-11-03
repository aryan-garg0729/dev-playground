package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"path"
	"strings"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// SQLResultEqualsConfig configures sql_result_equals: the learner's single
// statement, run on a dataset, must return the expected rows. Visible checks
// list Expected (shown to the learner and checked in the browser too); hidden
// checks compare with the reference solution on their own dataset.
//
//   - kind: sql_result_equals
//     name: Lists each customer's total
//     ordered: true                 # row order matters (default false)
//     expected: { columns: [name, total], rows: [[Ada, 30]] }
//   - kind: sql_result_equals
//     setup: hidden_data.sql        # server-only data; default: the spec's setup
//     visible: false
type SQLResultEqualsConfig struct {
	common
	Ordered  bool          `json:"ordered"`
	Setup    string        `json:"setup,omitempty"`    // manifest only: a .sql file
	Dataset  int           `json:"dataset"`            // compiled: index into the job's datasets
	Expected *ExpectedRows `json:"expected,omitempty"` // required when visible
}

// ExpectedRows is an expected result set.
type ExpectedRows struct {
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
}

type sqlResultEquals struct{}

func init() { Register(sqlResultEquals{}) }

func (sqlResultEquals) Kind() string { return "sql_result_equals" }

func (sqlResultEquals) Validate(raw json.RawMessage) error {
	var c SQLResultEqualsConfig
	if err := decodeStrict(raw, &c); err != nil {
		return err
	}
	if c.Setup != "" && (path.Base(c.Setup) != c.Setup || !strings.HasSuffix(c.Setup, ".sql")) {
		return fmt.Errorf("setup: %q must be a .sql file in the exercise folder", c.Setup)
	}
	if c.Dataset < 0 {
		return errors.New("dataset: must be ≥ 0")
	}
	if c.visible(true) && c.Expected == nil {
		return errors.New("expected: required on visible checks (it is shown and checked in the browser)")
	}
	if c.Expected != nil {
		if len(c.Expected.Columns) == 0 {
			return errors.New("expected.columns: required")
		}
		for i, r := range c.Expected.Rows {
			if len(r) != len(c.Expected.Columns) {
				return fmt.Errorf("expected.rows[%d]: has %d cells, want %d", i, len(r), len(c.Expected.Columns))
			}
		}
	}
	return c.validate()
}

func (sqlResultEquals) Run(_ context.Context, in CheckInput) ([]exercise.CheckResult, error) {
	var c SQLResultEqualsConfig
	if err := decodeStrict(in.Config, &c); err != nil {
		return nil, err
	}
	name := c.Name
	visible := c.visible(true)
	if name == "" {
		name = "Returns the right rows"
		if !visible {
			name = "Hidden data"
		}
	}
	row := exercise.CheckResult{Name: name, Visible: visible, Weight: c.Weight}
	r := in.Report
	switch {
	case r == nil:
		return nil, errors.New("sql_result_equals: no report")
	case r.TimedOut:
		row.Feedback = "Did not finish: time limit exceeded"
		return []exercise.CheckResult{row}, nil
	case r.SQL == nil || c.Dataset >= len(r.SQL.Results):
		return nil, fmt.Errorf("sql_result_equals: no result for dataset %d", c.Dataset)
	}
	res := r.SQL.Results[c.Dataset]
	want := c.Expected
	if want == nil {
		if res.Solution == nil || res.Solution.Error != "" {
			return nil, fmt.Errorf("sql_result_equals: reference solution failed on dataset %d", c.Dataset)
		}
		want = &ExpectedRows{Columns: res.Solution.Columns, Rows: res.Solution.Rows}
	}
	got := res.Learner
	switch {
	case got.Error != "":
		row.Feedback = got.Error
	case len(got.Columns) != len(want.Columns):
		row.Feedback = fmt.Sprintf("Expected %d column(s), got %d", len(want.Columns), len(got.Columns))
		if visible {
			row.Feedback = fmt.Sprintf("Expected %d column(s) (%s), got %d", len(want.Columns), strings.Join(want.Columns, ", "), len(got.Columns))
		}
	case !got.Truncated && SameRows(got.Rows, want.Rows, c.Ordered):
		row.Passed = true
	case visible:
		row.Got, row.Expected = previewRows(got.Rows), previewRows(want.Rows)
	default:
		row.Feedback = "Wrong result"
	}
	return []exercise.CheckResult{row}, nil
}

// SameRows compares result sets: in order when ordered, else as multisets.
// Numbers compare with a small relative tolerance (1 and 1.0 are equal).
// Must match sandbox/src/sqlcompare.ts.
func SameRows(got, want [][]any, ordered bool) bool {
	if len(got) != len(want) {
		return false
	}
	if ordered {
		for i := range got {
			if !sameRow(got[i], want[i]) {
				return false
			}
		}
		return true
	}
	left := append([][]any(nil), want...)
	for _, g := range got {
		found := -1
		for i, w := range left {
			if sameRow(g, w) {
				found = i
				break
			}
		}
		if found < 0 {
			return false
		}
		left = append(left[:found], left[found+1:]...)
	}
	return true
}

func sameRow(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !sameCell(a[i], b[i]) {
			return false
		}
	}
	return true
}

func sameCell(a, b any) bool {
	fa, aNum := number(a)
	fb, bNum := number(b)
	if aNum && bNum {
		return math.Abs(fa-fb) <= 1e-9*math.Max(1, math.Max(math.Abs(fa), math.Abs(fb)))
	}
	if aNum || bNum {
		return false
	}
	return a == b // strings or nil
}

func number(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}

func previewRows(rows [][]any) string {
	shown := rows
	if len(shown) > 5 {
		shown = shown[:5]
	}
	b, _ := json.Marshal(shown)
	s := string(b)
	if len(rows) > 5 {
		s = s[:len(s)-1] + fmt.Sprintf(", … %d rows]", len(rows))
	}
	return clip(s)
}
