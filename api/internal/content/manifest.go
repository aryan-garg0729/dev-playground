// Package content loads and validates course folders and compiles them into a
// Bundle that cq build writes to Postgres. The API never parses YAML.
package content

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

var (
	courseIDRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	dottedIDRe = regexp.MustCompile(`^[a-z0-9]+([.-][a-z0-9]+)*$`)
)

// courseManifest is course.yaml.
type courseManifest struct {
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Order       int      `yaml:"order"`
	Published   *bool    `yaml:"published"`
	Chapters    []string `yaml:"chapters"` // chapter folder names, in order
	Icon        string   `yaml:"icon"`     // 1-4 lowercase letters or digits shown in the course tile ("py")
	Level       string   `yaml:"level"`    // beginner (default) | intermediate | advanced
	Tags        []string `yaml:"tags"`     // catalog filters ("Python", "AI"), at most 4
	Topics      string   `yaml:"topics"`   // one short line ("Chat models, prompts, LCEL")
}

// pathsManifest is the optional paths.yaml at the root of a folder of courses.
type pathsManifest struct {
	Paths []pathManifest `yaml:"paths"`
}

type pathManifest struct {
	ID          string         `yaml:"id"`
	Title       string         `yaml:"title"`
	Description string         `yaml:"description"`
	Goal        string         `yaml:"goal"` // onboarding goal that makes this "your path"
	Courses     []pathStepYAML `yaml:"courses"`
}

type pathStepYAML struct {
	Course string `yaml:"course"`
	Note   string `yaml:"note"`
}

// chapterManifest is chapter.yaml.
type chapterManifest struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	Exercises []string `yaml:"exercises"` // exercise folder names, in order
}

// exerciseHeader is the common header of exercise.yaml.
type exerciseHeader struct {
	ID         string    `yaml:"id"`
	Type       string    `yaml:"type"`
	Version    int       `yaml:"version"`
	Title      string    `yaml:"title"`
	XP         *int      `yaml:"xp"`
	Difficulty int       `yaml:"difficulty"`
	Tags       []string  `yaml:"tags"`
	Free       bool      `yaml:"free"`
	Boss       bool      `yaml:"boss"`
	Spec       yaml.Node `yaml:"spec"`
}

// Bundle is the compiled content tree.
type Bundle struct {
	Courses []Course
	Paths   []Path
}

// Path is a compiled career path: courses in order.
type Path struct {
	ID, Title, Description, Goal string
	Steps                        []PathStep
}

// PathStep is one course of a path, with a short note shown under it.
type PathStep struct {
	Course string `json:"course"`
	Note   string `json:"note"`
}

// Course is a compiled course.
type Course struct {
	ID, Title, Description string
	Order                  int
	Published              bool
	Chapters               []Chapter
	Icon, Level, Topics    string
	Tags                   []string
}

// Chapter is a compiled chapter.
type Chapter struct {
	ID, Title string
	Exercises []Compiled
}

// Compiled is an exercise ready to store: header, public view and private spec.
type Compiled struct {
	exercise.Exercise
	PublicView  json.RawMessage
	PrivateSpec json.RawMessage
	Assets      []Asset // public images from assets/, sorted by name
}

// Problem is one validation failure: a file path and a reason.
type Problem struct {
	Path   string
	Reason string
}

func (p Problem) String() string { return p.Path + ": " + p.Reason }

// ValidationError collects every problem found in a content tree.
type ValidationError struct {
	Problems []Problem
}

func (e *ValidationError) Error() string {
	lines := make([]string, len(e.Problems))
	for i, p := range e.Problems {
		lines[i] = p.String()
	}
	return fmt.Sprintf("%d content problem(s):\n%s", len(e.Problems), strings.Join(lines, "\n"))
}
