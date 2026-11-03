package tutor

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Ollama talks to a local Ollama server's /api/chat (streamed NDJSON).
type Ollama struct {
	URL    string // e.g. http://localhost:11434
	Name   string // model, e.g. qwen2.5-coder:3b
	Client *http.Client
}

// Model returns the model name.
func (o *Ollama) Model() string { return o.Name }

type ollamaChunk struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done            bool   `json:"done"`
	PromptEvalCount int    `json:"prompt_eval_count"`
	EvalCount       int    `json:"eval_count"`
	Error           string `json:"error"`
}

// Stream sends the request and forwards each content piece.
func (o *Ollama) Stream(ctx context.Context, req Request, onDelta func(string) error) (Usage, error) {
	msgs := make([]Message, 0, len(req.Messages)+1)
	if req.System != "" {
		msgs = append(msgs, Message{Role: "system", Content: req.System})
	}
	msgs = append(msgs, req.Messages...)
	body, err := json.Marshal(map[string]any{
		"model":    o.Name,
		"messages": msgs,
		"stream":   true,
		"options":  map[string]any{"num_predict": req.MaxTokens, "temperature": 0.3},
	})
	if err != nil {
		return Usage{}, fmt.Errorf("ollama: encode: %w", err)
	}
	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(o.URL, "/")+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return Usage{}, fmt.Errorf("%w: %w", ErrProvider, err)
	}
	hreq.Header.Set("Content-Type", "application/json")
	client := o.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(hreq)
	if err != nil {
		return Usage{}, fmt.Errorf("%w: ollama: %w", ErrProvider, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return Usage{}, fmt.Errorf("%w: ollama: %s: %s", ErrProvider, resp.Status, strings.TrimSpace(string(msg)))
	}
	var u Usage
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 64<<10), 1<<20)
	for sc.Scan() {
		line := sc.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var c ollamaChunk
		if err := json.Unmarshal(line, &c); err != nil {
			return u, fmt.Errorf("%w: ollama: bad chunk: %w", ErrProvider, err)
		}
		if c.Error != "" {
			return u, fmt.Errorf("%w: ollama: %s", ErrProvider, c.Error)
		}
		if c.Message.Content != "" {
			if err := onDelta(c.Message.Content); err != nil {
				return u, err
			}
		}
		if c.Done {
			u.InputTokens, u.OutputTokens = c.PromptEvalCount, c.EvalCount
			return u, nil
		}
	}
	if err := sc.Err(); err != nil {
		return u, fmt.Errorf("%w: ollama: read: %w", ErrProvider, err)
	}
	return u, fmt.Errorf("%w: ollama: stream ended early", ErrProvider)
}
