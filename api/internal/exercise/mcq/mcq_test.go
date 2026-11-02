package mcq

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"testing/fstest"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

func node(t *testing.T, src string) yaml.Node {
	t.Helper()
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(src), &doc); err != nil {
		t.Fatal(err)
	}
	return *doc.Content[0]
}

var dir = fstest.MapFS{"prompt.md": {Data: []byte("What does `range(5)` yield?")}}

const validSpec = `
question: prompt.md
multiple: false
options:
  - { id: a, text: "0 to 4" }
  - { id: b, text: "1 to 5" }
answer: [a]
explain: range(5) starts at 0 and stops before 5.
`

func parse(t *testing.T, src string) (exercise.Exercise, error) {
	t.Helper()
	typ, _ := exercise.Get("mcq")
	s, err := typ.ParseSpec(node(t, src), dir)
	return exercise.Exercise{ID: "t.q", Type: "mcq", Spec: s}, err
}

func TestParseSpecErrors(t *testing.T) {
	tests := []struct {
		name, src, want string
	}{
		{"unknown field", validSpec + "colour: red\n", "colour"},
		{"missing question", strings.Replace(validSpec, "question: prompt.md", "", 1), "question: required"},
		{"missing file", strings.Replace(validSpec, "prompt.md", "nope.md", 1), "nope.md"},
		{"escaping path", strings.Replace(validSpec, "prompt.md", "../secret.md", 1), "relative path"},
		{"answer not an option", strings.Replace(validSpec, "answer: [a]", "answer: [z]", 1), `"z" is not an option`},
		{"two answers single", strings.Replace(validSpec, "answer: [a]", "answer: [a, b]", 1), "multiple is false"},
		{"duplicate ids", strings.Replace(validSpec, "id: b", "id: a", 1), "duplicate"},
		{"missing explain", strings.Replace(validSpec, "explain: range(5) starts at 0 and stops before 5.", "", 1), "explain"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(t, tt.src)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func TestPublicViewHidesAnswerAndExplain(t *testing.T) {
	ex, err := parse(t, validSpec)
	if err != nil {
		t.Fatal(err)
	}
	typ, _ := exercise.Get("mcq")
	pv, err := typ.PublicView(ex)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(pv, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"answer", "explain"} {
		if _, ok := m[k]; ok {
			t.Errorf("public view leaks %q: %s", k, pv)
		}
	}
	if strings.Contains(string(pv), "stops before 5") {
		t.Errorf("public view leaks explanation text: %s", pv)
	}
}

func TestLoadSpecRoundTrip(t *testing.T) {
	ex, err := parse(t, validSpec)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(ex.Spec)
	typ, _ := exercise.Get("mcq")
	s, err := typ.LoadSpec(raw)
	if err != nil {
		t.Fatal(err)
	}
	if s.(Spec).Explain != ex.Spec.(Spec).Explain {
		t.Fatal("round trip lost data")
	}
}

func TestGrade(t *testing.T) {
	ex, err := parse(t, validSpec)
	if err != nil {
		t.Fatal(err)
	}
	typ, _ := exercise.Get("mcq")
	tests := []struct {
		name       string
		answer     string
		wantStatus string
		invalid    bool
	}{
		{"correct", `{"choices":["a"]}`, exercise.StatusPass, false},
		{"wrong", `{"choices":["b"]}`, exercise.StatusFail, false},
		{"two on single", `{"choices":["a","b"]}`, "", true},
		{"unknown option", `{"choices":["z"]}`, "", true},
		{"empty", `{"choices":[]}`, "", true},
		{"malformed", `[1]`, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := typ.Grade(context.Background(), ex, exercise.Submission{Answer: json.RawMessage(tt.answer)})
			if tt.invalid {
				if !errors.Is(err, exercise.ErrInvalidAnswer) {
					t.Fatalf("got %v, want ErrInvalidAnswer", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if r.Status != tt.wantStatus {
				t.Fatalf("status %s, want %s", r.Status, tt.wantStatus)
			}
			if r.Checks[0].Feedback == "" {
				t.Fatal("explanation should be returned after submit")
			}
		})
	}
}
