package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// ExactChoiceConfig configures exact_choice: the selected option ids must equal
// the answer set (order ignored). Feedback is shown after any submit.
type ExactChoiceConfig struct {
	Name     string   `json:"name"`
	Answer   []string `json:"answer"`
	Feedback string   `json:"feedback,omitempty"`
}

// ExactChoiceAnswer is the mcq answer shape.
type ExactChoiceAnswer struct {
	Choices []string `json:"choices"`
}

type exactChoice struct{}

func init() { Register(exactChoice{}) }

func (exactChoice) Kind() string { return "exact_choice" }

func (exactChoice) Run(_ context.Context, in CheckInput) ([]exercise.CheckResult, error) {
	var cfg ExactChoiceConfig
	if err := json.Unmarshal(in.Config, &cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	var ans ExactChoiceAnswer
	if err := json.Unmarshal(in.Answer, &ans); err != nil {
		return nil, fmt.Errorf("answer: %w", err)
	}
	name := cfg.Name
	if name == "" {
		name = "Correct answer"
	}
	return []exercise.CheckResult{{
		Name:     name,
		Passed:   sameSet(ans.Choices, cfg.Answer),
		Visible:  true,
		Feedback: cfg.Feedback,
	}}, nil
}

func sameSet(a, b []string) bool {
	a, b = slices.Clone(a), slices.Clone(b)
	slices.Sort(a)
	slices.Sort(b)
	return slices.Equal(slices.Compact(a), slices.Compact(b))
}
