package content

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/aryangarg/codequest/api/internal/testdb"
)

func TestBuildIsIdempotentAndUnpublishes(t *testing.T) {
	pool, _ := testdb.New(t)
	ctx := context.Background()
	files := baseTree()
	files["c/course.yaml"] += "icon: py\ntags: [Python]\n"
	files["paths.yaml"] = "paths:\n  - { id: p, title: P, goal: ai, courses: [{ course: c, note: Start here }] }\n"
	root := writeTree(t, files)

	q := func(sql string) int {
		t.Helper()
		var n int
		if err := pool.QueryRow(ctx, sql).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}

	for range 2 { // building twice must give the same rows
		b, err := Load(root)
		if err != nil {
			t.Fatal(err)
		}
		st, err := Build(ctx, pool, b)
		if err != nil {
			t.Fatal(err)
		}
		if st.Exercises != 3 || st.Chapters != 2 || st.Courses != 1 {
			t.Fatalf("stats %+v", st)
		}
	}
	if n := q(`SELECT count(*) FROM exercises WHERE published`); n != 3 {
		t.Fatalf("published exercises = %d, want 3", n)
	}
	if n := q(`SELECT order_index FROM exercises WHERE id = 'c.one.a'`); n != 1 {
		t.Fatalf("c.one.a order_index = %d, want 1 (listed second)", n)
	}
	if n := q(`SELECT count(*) FROM courses WHERE id = 'c' AND icon = 'py' AND tags = '{Python}' AND level = 'beginner'`); n != 1 {
		t.Fatal("course catalog fields not stored")
	}
	if n := q(`SELECT count(*) FROM career_paths WHERE published AND goal = 'ai' AND steps->0->>'note' = 'Start here'`); n != 1 {
		t.Fatal("career path not stored")
	}

	// Drop paths.yaml: the path is unpublished, not deleted.
	if err := os.Remove(filepath.Join(root, "paths.yaml")); err != nil {
		t.Fatal(err)
	}

	// Remove exercise b from the chapter: it is unpublished, not deleted.
	if err := os.RemoveAll(filepath.Join(root, "c/ch1/b")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "c/ch1/chapter.yaml"), []byte("id: c.one\ntitle: One\nexercises: [a]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	b, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Build(ctx, pool, b); err != nil {
		t.Fatal(err)
	}
	if n := q(`SELECT count(*) FROM exercises WHERE id = 'c.one.b' AND NOT published`); n != 1 {
		t.Fatal("removed exercise should remain as an unpublished row")
	}
	if n := q(`SELECT count(*) FROM career_paths WHERE id = 'p' AND NOT published`); n != 1 {
		t.Fatal("removed path should remain as an unpublished row")
	}
	if n := q(`SELECT order_index FROM exercises WHERE id = 'c.one.a'`); n != 0 {
		t.Fatalf("c.one.a order_index = %d, want 0", n)
	}
}
