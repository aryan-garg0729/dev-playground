package content

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/aryangarg/codequest/api/internal/exercise/all"
)

// writeTree writes files (path -> content) under a temp dir and returns it.
func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, body := range files {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const mcqYAML = `id: %s
type: mcq
version: 1
title: Q
xp: 5
spec:
  question: prompt.md
  options: [{ id: a, text: "A" }, { id: b, text: "B" }]
  answer: [a]
  explain: Because.
`

func ex(id string) string { return strings.Replace(mcqYAML, "%s", id, 1) }

func baseTree() map[string]string {
	return map[string]string{
		"c/course.yaml":            "id: c\ntitle: C\nchapters: [ch2, ch1]\n",
		"c/ch1/chapter.yaml":       "id: c.one\ntitle: One\nexercises: [b, a]\n",
		"c/ch1/a/exercise.yaml":    ex("c.one.a"),
		"c/ch1/a/prompt.md":        "?",
		"c/ch1/b/exercise.yaml":    ex("c.one.b"),
		"c/ch1/b/prompt.md":        "?",
		"c/ch2/chapter.yaml":       "id: c.two\ntitle: Two\nexercises: [x]\n",
		"c/ch2/x/exercise.yaml":    ex("c.two.x"),
		"c/ch2/x/prompt.md":        "?",
		"c/ch2/x/wrong/ignored.md": "subfolders inside an exercise are its own business",
	}
}

func TestLoadKeepsManifestOrder(t *testing.T) {
	b, err := Load(writeTree(t, baseTree()))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, ch := range b.Courses[0].Chapters {
		got = append(got, ch.ID)
		for _, e := range ch.Exercises {
			got = append(got, e.ID)
		}
	}
	want := "c.two c.two.x c.one c.one.b c.one.a"
	if strings.Join(got, " ") != want {
		t.Fatalf("order = %v, want %s (manifest order, not folder order)", got, want)
	}
	if e := b.Courses[0].Chapters[0].Exercises[0]; e.Difficulty != 1 || e.Tags == nil {
		t.Fatalf("defaults not applied: %+v", e.Exercise)
	}
}

func TestLoadSingleCourseRoot(t *testing.T) {
	root := writeTree(t, baseTree())
	b, err := Load(filepath.Join(root, "c"))
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Courses) != 1 || b.Courses[0].ID != "c" {
		t.Fatalf("got %+v", b.Courses)
	}
}

func TestLoadProblems(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(m map[string]string)
		wantPath string
		wantMsg  string
	}{
		{"unknown course field", func(m map[string]string) { m["c/course.yaml"] += "colour: red\n" },
			"c/course.yaml", "colour"},
		{"bad course id", func(m map[string]string) { m["c/course.yaml"] = "id: C_1\ntitle: C\nchapters: [ch1, ch2]\n" },
			"c/course.yaml", "slug"},
		{"listed folder missing", func(m map[string]string) { m["c/course.yaml"] = "id: c\ntitle: C\nchapters: [ch1, ch2, ch3]\n" },
			"c/course.yaml", `"ch3" has no`},
		{"unlisted exercise folder", func(m map[string]string) {
			m["c/ch1/z/exercise.yaml"] = ex("c.one.z")
			m["c/ch1/z/prompt.md"] = "?"
		}, "c/ch1/z", "not listed"},
		{"duplicate exercise id", func(m map[string]string) { m["c/ch2/x/exercise.yaml"] = ex("c.one.a") },
			"c/ch1/a/exercise.yaml", "already used"}, // ch2 loads first (manifest order)
		{"missing title", func(m map[string]string) {
			m["c/ch1/a/exercise.yaml"] = strings.Replace(ex("c.one.a"), "title: Q\n", "", 1)
		}, "c/ch1/a/exercise.yaml", "title is required"},
		{"missing xp", func(m map[string]string) {
			m["c/ch1/a/exercise.yaml"] = strings.Replace(ex("c.one.a"), "xp: 5\n", "", 1)
		}, "c/ch1/a/exercise.yaml", "xp is required"},
		{"bad difficulty", func(m map[string]string) { m["c/ch1/a/exercise.yaml"] = ex("c.one.a") + "difficulty: 9\n" },
			"c/ch1/a/exercise.yaml", "difficulty must be 1-5"},
		{"unknown type", func(m map[string]string) {
			m["c/ch1/a/exercise.yaml"] = strings.Replace(ex("c.one.a"), "type: mcq", "type: essay", 1)
		}, "c/ch1/a/exercise.yaml", `type "essay" is not registered`},
		{"mcq boss", func(m map[string]string) { m["c/ch1/a/exercise.yaml"] = ex("c.one.a") + "boss: true\n" },
			"c/ch1/a/exercise.yaml", "only code and project_cli exercises can be bosses"},
		{"spec error", func(m map[string]string) {
			m["c/ch1/a/exercise.yaml"] = strings.Replace(ex("c.one.a"), "answer: [a]", "answer: [q]", 1)
		}, "c/ch1/a/exercise.yaml", `spec: answer: "q" is not an option id`},
		{"empty manifest", func(m map[string]string) { m["c/ch1/chapter.yaml"] = "" },
			"c/ch1/chapter.yaml", "empty"},
		{"folder without course.yaml", func(m map[string]string) { m["stray/readme.md"] = "hi" },
			"stray", "no course.yaml"},
		{"bad icon", func(m map[string]string) { m["c/course.yaml"] += "icon: Python\n" },
			"c/course.yaml", "icon"},
		{"bad level", func(m map[string]string) { m["c/course.yaml"] += "level: expert\n" },
			"c/course.yaml", "level"},
		{"too many tags", func(m map[string]string) { m["c/course.yaml"] += "tags: [a, b, c, d, e]\n" },
			"c/course.yaml", "at most 4"},
		{"duplicate tag", func(m map[string]string) { m["c/course.yaml"] += "tags: [AI, ai]\n" },
			"c/course.yaml", "listed twice"},
		{"path with unknown course", func(m map[string]string) {
			m["paths.yaml"] = "paths:\n  - { id: p, title: P, courses: [{ course: nope }] }\n"
		}, "paths.yaml: paths[0]", `"nope" is not a course`},
		{"path with bad goal", func(m map[string]string) {
			m["paths.yaml"] = "paths:\n  - { id: p, title: P, goal: fame, courses: [{ course: c }] }\n"
		}, "paths.yaml: paths[0]", "goal"},
		{"two paths with one goal", func(m map[string]string) {
			m["paths.yaml"] = "paths:\n  - { id: p, title: P, goal: ai, courses: [{ course: c }] }\n  - { id: q, title: Q, goal: ai, courses: [{ course: c }] }\n"
		}, "paths.yaml: paths[1]", "already belongs"},
		{"unknown path field", func(m map[string]string) {
			m["paths.yaml"] = "paths:\n  - { id: p, title: P, colour: red, courses: [{ course: c }] }\n"
		}, "paths.yaml", "colour"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := baseTree()
			tt.mutate(files)
			root := writeTree(t, files)
			_, err := Load(root)
			var verr *ValidationError
			if !errors.As(err, &verr) {
				t.Fatalf("want ValidationError, got %v", err)
			}
			want := filepath.Join(root, tt.wantPath)
			for _, p := range verr.Problems {
				if p.Path == want && strings.Contains(p.Reason, tt.wantMsg) {
					return
				}
			}
			t.Fatalf("no problem at %s containing %q; got:\n%v", tt.wantPath, tt.wantMsg, verr)
		})
	}
}

func TestLoadCatalogAndPaths(t *testing.T) {
	files := baseTree()
	files["c/course.yaml"] += "icon: py\nlevel: intermediate\ntags: [Python, AI]\ntopics: Loops and lists\n"
	files["paths.yaml"] = "paths:\n  - id: p\n    title: P\n    goal: ai\n    courses:\n      - { course: c, note: First steps }\n"
	b, err := Load(writeTree(t, files))
	if err != nil {
		t.Fatal(err)
	}
	c := b.Courses[0]
	if c.Icon != "py" || c.Level != "intermediate" || strings.Join(c.Tags, ",") != "Python,AI" || c.Topics != "Loops and lists" {
		t.Fatalf("catalog fields: %+v", c)
	}
	if len(b.Paths) != 1 || b.Paths[0].Goal != "ai" || b.Paths[0].Steps[0] != (PathStep{Course: "c", Note: "First steps"}) {
		t.Fatalf("paths: %+v", b.Paths)
	}

	plain, err := Load(writeTree(t, baseTree()))
	if err != nil {
		t.Fatal(err)
	}
	if c := plain.Courses[0]; c.Icon != "c" || c.Level != "beginner" || c.Tags == nil || len(plain.Paths) != 0 {
		t.Fatalf("defaults: %+v, paths %+v", c, plain.Paths)
	}
}

func TestFixtures(t *testing.T) {
	if _, err := Load("testdata/valid"); err != nil {
		t.Fatalf("valid fixture: %v", err)
	}
	_, err := Load("testdata/invalid")
	var verr *ValidationError
	if !errors.As(err, &verr) || len(verr.Problems) != 4 {
		t.Fatalf("invalid fixture: want 4 problems, got %v", err)
	}
}

func TestSampleContentIsValid(t *testing.T) {
	b, err := Load("../../../content")
	if err != nil {
		t.Fatalf("content/ must always validate: %v", err)
	}
	if len(b.Courses) == 0 {
		t.Fatal("no sample courses")
	}
}
