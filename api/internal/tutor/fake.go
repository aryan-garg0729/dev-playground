package tutor

import (
	"context"
	"strings"
	"sync"
)

// Fake is a scripted Provider for tests: it streams Reply word by word and
// records each request.
type Fake struct {
	Reply string
	Err   error // returned before any text, when set

	mu   sync.Mutex
	reqs []Request
}

// Model returns "fake".
func (f *Fake) Model() string { return "fake" }

// Stream replays Reply.
func (f *Fake) Stream(_ context.Context, req Request, onDelta func(string) error) (Usage, error) {
	f.mu.Lock()
	f.reqs = append(f.reqs, req)
	f.mu.Unlock()
	if f.Err != nil {
		return Usage{}, f.Err
	}
	words := strings.SplitAfter(f.Reply, " ")
	for _, w := range words {
		if err := onDelta(w); err != nil {
			return Usage{}, err
		}
	}
	return Usage{InputTokens: 10, OutputTokens: len(words)}, nil
}

// Requests returns the requests seen so far.
func (f *Fake) Requests() []Request {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]Request(nil), f.reqs...)
}
