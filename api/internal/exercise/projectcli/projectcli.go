package projectcli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// TypeName is the manifest `type:` of these lessons.
const TypeName = "project_cli"

// manifest is the spec block as written in exercise.yaml.
type manifest struct {
	Prompt    string   `yaml:"prompt"`
	BaseURL   string   `yaml:"base_url"`
	AllowedOS []string `yaml:"allowed_os"`
	Steps     []Step   `yaml:"steps"`
	Solution  string   `yaml:"solution"`
}

type cliType struct{}

func init() { exercise.Register(cliType{}) }

func (cliType) Name() string { return TypeName }

func readFile(dir fs.FS, field, name string, limit int) ([]byte, error) {
	if name == "" {
		return nil, fmt.Errorf("%s: required", field)
	}
	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("%s: %q must be a relative path inside the exercise folder", field, name)
	}
	b, err := fs.ReadFile(dir, name)
	if err != nil {
		return nil, fmt.Errorf("%s: read %s: %w", field, name, err)
	}
	if len(b) > limit {
		return nil, fmt.Errorf("%s: %s is larger than %d bytes", field, name, limit)
	}
	if !utf8.Valid(b) {
		return nil, fmt.Errorf("%s: %s is not UTF-8", field, name)
	}
	return b, nil
}

// decodeRecorded strictly decodes a recorded answer file and checks its shape.
func decodeRecorded(s Spec, name string, b []byte) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var a Answer
	if err := dec.Decode(&a); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if _, err := decodeAnswer(s, b); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

func (cliType) ParseSpec(raw yaml.Node, dir fs.FS) (exercise.Spec, error) {
	var m manifest
	if err := exercise.DecodeStrict(raw, &m); err != nil {
		return nil, err
	}
	prompt, err := readFile(dir, "prompt", m.Prompt, maxPromptBytes)
	if err != nil {
		return nil, err
	}
	s := Spec{Prompt: string(prompt), BaseURL: m.BaseURL, AllowedOS: m.AllowedOS, Steps: m.Steps}
	for i := range s.Steps {
		if r := s.Steps[i].HTTPRequest; r != nil {
			r.Request.BodyJSON = normalize(r.Request.BodyJSON)
		}
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	if s.Solution, err = readFile(dir, "solution", m.Solution, maxAnswerFile); err != nil {
		return nil, err
	}
	if err := decodeRecorded(s, "solution", s.Solution); err != nil {
		return nil, err
	}
	if s.Wrong, err = readWrong(dir, s); err != nil {
		return nil, err
	}
	return s, nil
}

// readWrong loads wrong/*.json: recorded answers that cq test expects to fail.
func readWrong(dir fs.FS, s Spec) (map[string][]byte, error) {
	entries, err := fs.ReadDir(dir, "wrong")
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("wrong/: %w", err)
	}
	out := map[string][]byte{}
	for _, e := range entries {
		name := path.Join("wrong", e.Name())
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			return nil, fmt.Errorf("%s: wrong/ may hold only .json files", name)
		}
		b, err := readFile(dir, "wrong", name, maxAnswerFile)
		if err != nil {
			return nil, err
		}
		if err := decodeRecorded(s, name, b); err != nil {
			return nil, err
		}
		out[name] = b
	}
	return out, nil
}

func (cliType) LoadSpec(private json.RawMessage) (exercise.Spec, error) {
	var s Spec
	if err := json.Unmarshal(private, &s); err != nil {
		return nil, fmt.Errorf("project_cli: load spec: %w", err)
	}
	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("project_cli: load spec: %w", err)
	}
	return s, nil
}

func specOf(ex exercise.Exercise) (Spec, error) {
	s, ok := ex.Spec.(Spec)
	if !ok {
		return Spec{}, fmt.Errorf("project_cli: exercise %s has spec %T", ex.ID, ex.Spec)
	}
	return s, nil
}

// PublicView is the whole spec minus solution.json and wrong/ (the CLI needs
// the steps to run them and the tests to print them), plus the lesson id the
// web page shows in `codequest submit <id>`.
type PublicView struct {
	LessonID string `json:"lesson_id"`
	Spec
}

func (cliType) PublicView(ex exercise.Exercise) (json.RawMessage, error) {
	s, err := specOf(ex)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(PublicView{LessonID: ex.ID, Spec: Spec{Prompt: s.Prompt, BaseURL: s.BaseURL, AllowedOS: s.AllowedOS, Steps: s.Steps}})
	if err != nil {
		return nil, fmt.Errorf("project_cli: public view: %w", err)
	}
	return b, nil
}

// Grade judges the reported output. It runs nothing.
func (cliType) Grade(_ context.Context, ex exercise.Exercise, sub exercise.Submission) (exercise.Result, error) {
	s, err := specOf(ex)
	if err != nil {
		return exercise.Result{}, err
	}
	a, err := decodeAnswer(s, sub.Answer)
	if err != nil {
		return exercise.Result{}, err
	}
	return grade(s, a), nil
}

// ReferenceAnswer is solution.json.
func (cliType) ReferenceAnswer(sp exercise.Spec) (json.RawMessage, error) {
	s, ok := sp.(Spec)
	if !ok {
		return nil, fmt.Errorf("project_cli: spec %T", sp)
	}
	return s.Solution, nil
}

// WrongAnswers are wrong/*.json.
func (cliType) WrongAnswers(sp exercise.Spec) (map[string]json.RawMessage, error) {
	s, ok := sp.(Spec)
	if !ok {
		return nil, fmt.Errorf("project_cli: spec %T", sp)
	}
	out := make(map[string]json.RawMessage, len(s.Wrong))
	for k, v := range s.Wrong {
		out[k] = v
	}
	return out, nil
}
