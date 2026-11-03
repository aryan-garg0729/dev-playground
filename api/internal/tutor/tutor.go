// Package tutor is the AI tutor: it builds the hint request for an exercise
// and streams a model's answer. The model lives behind Provider so the vendor
// can change. The tutor never sees solution files or hidden test source: it is
// given only what the learner already sees (the prompt, their code, their
// latest Result) and the conversation so far.
//
// For now the tutor runs only in development, on a local Ollama model
// (decision 2026-11-03); production refuses to start with a provider set.
package tutor

import (
	"context"
	"errors"
)

// MaxOutputTokens caps each reply (SPEC "AI tutor").
const MaxOutputTokens = 400

// HistoryLen is how many earlier tutor messages are sent with a question.
const HistoryLen = 6

// Message is one turn of the conversation.
type Message struct {
	Role    string `json:"role"` // user | assistant
	Content string `json:"content"`
}

// Request is what a Provider is asked.
type Request struct {
	System    string
	Messages  []Message
	MaxTokens int
}

// Usage is the token count of one reply.
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// Provider streams a model reply, calling onDelta with each piece of text in
// order. An error from onDelta stops the stream and is returned.
type Provider interface {
	Model() string
	Stream(ctx context.Context, req Request, onDelta func(string) error) (Usage, error)
}

// ErrProvider wraps failures of the model backend (unreachable, bad status).
var ErrProvider = errors.New("tutor provider")
