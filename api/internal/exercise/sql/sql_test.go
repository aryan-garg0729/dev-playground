package sql_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"testing/fstest"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/exercise"
	_ "github.com/aryangarg/codequest/api/internal/exercise/sql"
	"github.com/aryangarg/codequest/api/internal/runner"
)

const specYAML = `
prompt: prompt.md
setup: setup.sql
solution: solution.sql
run: server
checks:
  - kind: sql_result_equals
    name: totals
    expected: { columns: [name, total], rows: [[Ada, 35], [Ravi, 40]] }
  - kind: sql_result_equals
    setup: hidden.sql
    visible: false
`

const setup = `CREATE TABLE o (name TEXT, amount INT);
INSERT INTO o VALUES ('Ada', 30), ('Ada', 5), ('Ravi', 40);`

func files() fstest.MapFS {
	return fstest.MapFS{
		"prompt.md":     {Data: []byte("Total per name.")},
		"setup.sql":     {Data: []byte(setup)},
		"solution.sql":  {Data: []byte("SELECT name, SUM(amount) AS total FROM o GROUP BY name -- SOLUTION_SECRET")},
		"hidden.sql":    {Data: []byte("CREATE TABLE o (name TEXT, amount INT); INSERT INTO o VALUES ('HIDDEN_SECRET', 1);")},
		"wrong/max.sql": {Data: []byte("SELECT name, MAX(amount) FROM o GROUP BY name")},
	}
}

func parse(t *testing.T, src string, fsys fstest.MapFS) (exercise.Type, exercise.Spec, error) {
	t.Helper()
	typ, _ := exercise.Get("sql")
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(src), &node); err != nil {
		t.Fatal(err)
	}
	spec, err := typ.ParseSpec(*node.Content[0], fsys)
	return typ, spec, err
}

func mustEx(t *testing.T) (exercise.Type, exercise.Exercise) {
	t.Helper()
	typ, spec, err := parse(t, specYAML, files())
	if err != nil {
		t.Fatal(err)
	}
	return typ, exercise.Exercise{ID: "t.sql", Type: "sql", Version: 1, Spec: spec}
}

func TestParseRejects(t *testing.T) {
	for name, src := range map[string]string{
		"visible without expected": strings.Replace(specYAML, "    expected: { columns: [name, total], rows: [[Ada, 35], [Ravi, 40]] }\n", "", 1),
		"visible with own setup":   strings.Replace(specYAML, "    visible: false\n", "", 1),
		"other check kind":         strings.Replace(specYAML, "kind: sql_result_equals\n    name: totals", "kind: unit_tests\n    name: totals", 1),
		"bad run mode":             strings.Replace(specYAML, "run: server", "run: cloud", 1),
		"ragged expected":          strings.Replace(specYAML, "[Ravi, 40]", "[Ravi]", 1),
		"unknown field":            specYAML + "bogus: 1\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, err := parse(t, src, files()); err == nil {
				t.Fatal("want error")
			}
		})
	}
}

func TestPublicViewIsPublicOnly(t *testing.T) {
	typ, ex := mustEx(t)
	pv, err := typ.PublicView(ex)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"SOLUTION_SECRET", "HIDDEN_SECRET", "hidden.sql"} {
		if strings.Contains(string(pv), secret) {
			t.Errorf("public view leaks %q: %s", secret, pv)
		}
	}
	if !strings.Contains(string(pv), `"visible_checks":[{"name":"totals"`) || !strings.Contains(string(pv), "INSERT INTO o") {
		t.Errorf("public view = %s", pv)
	}
}

func TestInvalidAnswers(t *testing.T) {
	typ, ex := mustEx(t)
	ctx := runner.NewContext(context.Background(), runner.Docker{})
	for _, a := range []string{`{"query":1}`, `{"query":"x","extra":1}`, `{"query":"` + strings.Repeat("x", 17<<10) + `"}`} {
		if _, err := typ.Grade(ctx, ex, exercise.Submission{Answer: json.RawMessage(a)}); !errors.Is(err, exercise.ErrInvalidAnswer) {
			t.Errorf("%.40s: %v", a, err)
		}
	}
}

func docker(t *testing.T) runner.Runner {
	t.Helper()
	if err := exec.Command("docker", "image", "inspect", runner.Images["docker-python"]).Run(); err != nil {
		if os.Getenv("CQ_REQUIRE_DOCKER") == "1" {
			t.Fatalf("runner image missing: %v", err)
		}
		t.Skip("docker runner image not available")
	}
	return runner.Docker{}
}

func TestGradeAndRunInSandbox(t *testing.T) {
	r := docker(t)
	typ, ex := mustEx(t)
	ctx := runner.NewContext(context.Background(), r)
	grade := func(q string) exercise.Result {
		ans, _ := json.Marshal(map[string]string{"query": q})
		res, err := typ.Grade(ctx, ex, exercise.Submission{Answer: ans})
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	if res := grade("SELECT name, SUM(amount) AS total FROM o GROUP BY name"); res.Status != exercise.StatusPass {
		t.Fatalf("solution: %+v", res)
	}
	res := grade("SELECT 'Ada', 35 UNION ALL SELECT 'Ravi', 40")
	if res.Status != exercise.StatusFail || !res.Checks[0].Passed || res.Checks[1].Passed || res.Checks[1].Feedback != "Wrong result" {
		t.Fatalf("hard-coded rows must fail only the hidden check: %+v", res)
	}
	if strings.Contains(res.Checks[1].Feedback+res.Checks[1].Got, "HIDDEN") {
		t.Fatal("hidden check leaks data")
	}
	if res := grade("ATTACH DATABASE '/work/x.db' AS x"); res.Status != exercise.StatusFail {
		t.Fatalf("attach: %+v", res)
	}
	if res := grade("WITH RECURSIVE c(x) AS (SELECT 1 UNION ALL SELECT x + 1 FROM c) SELECT count(*) FROM c"); res.Status != exercise.StatusTimeout {
		t.Fatalf("runaway query: %+v", res)
	}

	runner, _ := typ.(exercise.Runnable)
	out, err := runner.Run(ctx, ex, json.RawMessage(`{"query":"SELECT name, SUM(amount) AS total FROM o GROUP BY name"}`))
	if err != nil {
		t.Fatal(err)
	}
	if out.Table == nil || len(out.Table.Rows) != 2 || len(out.Checks) != 1 || !out.Checks[0].Passed {
		t.Fatalf("run: %+v", out)
	}
}
