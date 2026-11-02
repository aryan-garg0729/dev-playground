// Package mcq is the multiple-choice exercise type.
package mcq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"slices"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/checks"
	"github.com/aryangarg/codequest/api/internal/exercise"
)

const maxQuestionBytes = 64 << 10

var optionIDRe = regexp.MustCompile(`^[a-z0-9_-]{1,32}$`)

// Option is one choice.
type Option struct {
	ID   string `json:"id" yaml:"id"`
	Text string `json:"text" yaml:"text"`
}

// Spec is the compiled mcq spec. Answer and Explain are server-only until a
// submit; PublicView leaves them out.
type Spec struct {
	Question string   `json:"question"` // markdown, loaded from the file the manifest names
	Multiple bool     `json:"multiple"`
	Options  []Option `json:"options"`
	Answer   []string `json:"answer"`
	Explain  string   `json:"explain"`
}

// manifest is the spec block as written in exercise.yaml.
type manifest struct {
	Question string   `yaml:"question"`
	Multiple bool     `yaml:"multiple"`
	Options  []Option `yaml:"options"`
	Answer   []string `yaml:"answer"`
	Explain  string   `yaml:"explain"`
}

// PublicView is what the browser receives.
type PublicView struct {
	Question string   `json:"question"`
	Multiple bool     `json:"multiple"`
	Options  []Option `json:"options"`
}

// Answer is the learner's answer shape.
type Answer = checks.ExactChoiceAnswer

type mcqType struct{}

func init() { exercise.Register(mcqType{}) }

func (mcqType) Name() string { return "mcq" }

func (mcqType) ParseSpec(raw yaml.Node, dir fs.FS) (exercise.Spec, error) {
	var m manifest
	if err := exercise.DecodeStrict(raw, &m); err != nil {
		return nil, err
	}
	if m.Question == "" {
		return nil, errors.New("question: required (path to a markdown file, e.g. prompt.md)")
	}
	if !fs.ValidPath(m.Question) {
		return nil, fmt.Errorf("question: %q must be a relative path inside the exercise folder", m.Question)
	}
	q, err := fs.ReadFile(dir, m.Question)
	if err != nil {
		return nil, fmt.Errorf("question: read %s: %w", m.Question, err)
	}
	if len(q) > maxQuestionBytes {
		return nil, fmt.Errorf("question: %s is larger than %d bytes", m.Question, maxQuestionBytes)
	}
	s := Spec{Question: string(q), Multiple: m.Multiple, Options: m.Options, Answer: m.Answer, Explain: m.Explain}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s Spec) validate() error {
	if len(s.Options) < 2 {
		return errors.New("options: need at least 2")
	}
	ids := make([]string, 0, len(s.Options))
	for i, o := range s.Options {
		if !optionIDRe.MatchString(o.ID) {
			return fmt.Errorf("options[%d].id: %q must match %s", i, o.ID, optionIDRe)
		}
		if slices.Contains(ids, o.ID) {
			return fmt.Errorf("options[%d].id: duplicate id %q", i, o.ID)
		}
		if o.Text == "" {
			return fmt.Errorf("options[%d].text: required", i)
		}
		ids = append(ids, o.ID)
	}
	if len(s.Answer) == 0 {
		return errors.New("answer: required")
	}
	if !s.Multiple && len(s.Answer) != 1 {
		return fmt.Errorf("answer: has %d ids but multiple is false", len(s.Answer))
	}
	for _, a := range s.Answer {
		if !slices.Contains(ids, a) {
			return fmt.Errorf("answer: %q is not an option id", a)
		}
	}
	if s.Explain == "" {
		return errors.New("explain: required")
	}
	return nil
}

func (mcqType) LoadSpec(private json.RawMessage) (exercise.Spec, error) {
	var s Spec
	if err := json.Unmarshal(private, &s); err != nil {
		return nil, fmt.Errorf("mcq: load spec: %w", err)
	}
	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("mcq: load spec: %w", err)
	}
	return s, nil
}

func specOf(ex exercise.Exercise) (Spec, error) {
	s, ok := ex.Spec.(Spec)
	if !ok {
		return Spec{}, fmt.Errorf("mcq: exercise %s has spec %T", ex.ID, ex.Spec)
	}
	return s, nil
}

func (mcqType) PublicView(ex exercise.Exercise) (json.RawMessage, error) {
	s, err := specOf(ex)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(PublicView{Question: s.Question, Multiple: s.Multiple, Options: s.Options})
	if err != nil {
		return nil, fmt.Errorf("mcq: public view: %w", err)
	}
	return b, nil
}

func (mcqType) Grade(ctx context.Context, ex exercise.Exercise, sub exercise.Submission) (exercise.Result, error) {
	s, err := specOf(ex)
	if err != nil {
		return exercise.Result{}, err
	}
	var ans Answer
	if err := json.Unmarshal(sub.Answer, &ans); err != nil {
		return exercise.Result{}, fmt.Errorf("%w: %w", exercise.ErrInvalidAnswer, err)
	}
	if len(ans.Choices) == 0 {
		return exercise.Result{}, fmt.Errorf("%w: choose at least one option", exercise.ErrInvalidAnswer)
	}
	if !s.Multiple && len(ans.Choices) != 1 {
		return exercise.Result{}, fmt.Errorf("%w: choose exactly one option", exercise.ErrInvalidAnswer)
	}
	for _, c := range ans.Choices {
		if !slices.ContainsFunc(s.Options, func(o Option) bool { return o.ID == c }) {
			return exercise.Result{}, fmt.Errorf("%w: unknown option %q", exercise.ErrInvalidAnswer, c)
		}
	}
	cfg, err := json.Marshal(checks.ExactChoiceConfig{Answer: s.Answer, Feedback: s.Explain})
	if err != nil {
		return exercise.Result{}, fmt.Errorf("mcq: check config: %w", err)
	}
	rows, err := checks.Run(ctx, "exact_choice", checks.CheckInput{Config: cfg, Answer: sub.Answer})
	if err != nil {
		return exercise.Result{}, fmt.Errorf("mcq: grade: %w", err)
	}
	return exercise.NewResult(rows), nil
}
