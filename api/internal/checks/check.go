// Package checks holds reusable check primitives. Each kind is registered by
// name and can be used by any exercise type.
package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// CheckInput is everything a check may look at. Config is the check's own
// manifest entry; Answer is the learner's answer. Runner output is added with
// the code type in Phase 2.
type CheckInput struct {
	Config json.RawMessage
	Answer json.RawMessage
}

// Check is one check kind.
type Check interface {
	Kind() string
	Run(ctx context.Context, in CheckInput) ([]exercise.CheckResult, error)
}

var registry = map[string]Check{}

// Register adds a check kind. Called from init().
func Register(c Check) {
	if _, dup := registry[c.Kind()]; dup {
		panic("checks: kind registered twice: " + c.Kind())
	}
	registry[c.Kind()] = c
}

// Get returns a check kind.
func Get(kind string) (Check, bool) { c, ok := registry[kind]; return c, ok }

// Kinds lists registered kinds, sorted.
func Kinds() []string {
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Run looks up kind and runs it.
func Run(ctx context.Context, kind string, in CheckInput) ([]exercise.CheckResult, error) {
	c, ok := Get(kind)
	if !ok {
		return nil, fmt.Errorf("unknown check kind %q", kind)
	}
	res, err := c.Run(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("check %s: %w", kind, err)
	}
	return res, nil
}
