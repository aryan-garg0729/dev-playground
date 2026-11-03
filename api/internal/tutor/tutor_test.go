package tutor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/aryangarg/codequest/api/internal/testredis"
)

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		name     string
		ex       Exercise
		want     []string
		wantNot  []string
		messages int
	}{
		{
			name: "code with a failing result and history",
			ex: Exercise{
				Title: "Sum a list", Type: "code",
				Public:   json.RawMessage(`{"prompt":"Return the sum of xs.","language":"python","hints":["loop"]}`),
				Code:     "def total(xs):\n    return 0",
				Result:   json.RawMessage(`{"status":"fail","checks":[{"name":"adds numbers","passed":false,"got":"0","expected":"6","line":2},{"name":"empty","passed":true}],"stdout":""}`),
				History:  []Message{{Role: "user", Content: "hi"}, {Role: "assistant", Content: "hello"}},
				Question: "why does it fail?",
			},
			want:     []string{"Return the sum of xs.", "  2 |     return 0", "FAIL adds numbers (line 2): got 0, expected 6", "PASS empty", "```python"},
			messages: 3,
		},
		{
			name:     "no submission yet",
			ex:       Exercise{Title: "Pick", Type: "mcq", Public: json.RawMessage(`{"question":"Which one?"}`), Question: "help"},
			want:     []string{"Which one?", "has not submitted yet"},
			wantNot:  []string{"Learner's current code"},
			messages: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := BuildRequest(tt.ex)
			for _, s := range tt.want {
				if !strings.Contains(r.System, s) {
					t.Errorf("system prompt missing %q:\n%s", s, r.System)
				}
			}
			for _, s := range tt.wantNot {
				if strings.Contains(r.System, s) {
					t.Errorf("system prompt has %q", s)
				}
			}
			if len(r.Messages) != tt.messages || r.Messages[len(r.Messages)-1].Content != tt.ex.Question {
				t.Errorf("messages = %+v", r.Messages)
			}
			if r.MaxTokens != MaxOutputTokens {
				t.Errorf("max tokens = %d", r.MaxTokens)
			}
		})
	}
}

func TestOllamaStream(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			http.NotFound(w, r)
			return
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		for _, c := range []string{"What ", "does ", "line 2 return?"} {
			_, _ = fmt.Fprintf(w, `{"message":{"role":"assistant","content":%q},"done":false}`+"\n", c)
		}
		_, _ = fmt.Fprintln(w, `{"message":{"content":""},"done":true,"prompt_eval_count":120,"eval_count":9}`)
	}))
	defer srv.Close()

	o := &Ollama{URL: srv.URL, Name: "m"}
	var out strings.Builder
	u, err := o.Stream(context.Background(), Request{System: "sys", Messages: []Message{{Role: "user", Content: "q"}}, MaxTokens: 400},
		func(s string) error { out.WriteString(s); return nil })
	if err != nil {
		t.Fatal(err)
	}
	if out.String() != "What does line 2 return?" || u != (Usage{InputTokens: 120, OutputTokens: 9}) {
		t.Fatalf("got %q %+v", out.String(), u)
	}
	msgs, _ := got["messages"].([]any)
	if got["model"] != "m" || got["stream"] != true || len(msgs) != 2 {
		t.Fatalf("request = %v", got)
	}
	if opts, _ := got["options"].(map[string]any); opts["num_predict"] != float64(400) {
		t.Fatalf("options = %v", got["options"])
	}
}

func TestOllamaErrors(t *testing.T) {
	for name, h := range map[string]http.HandlerFunc{
		"bad status":     func(w http.ResponseWriter, _ *http.Request) { http.Error(w, "boom", http.StatusInternalServerError) },
		"error in body":  func(w http.ResponseWriter, _ *http.Request) { _, _ = fmt.Fprintln(w, `{"error":"model not found"}`) },
		"ends too early": func(w http.ResponseWriter, _ *http.Request) { _, _ = fmt.Fprintln(w, `{"message":{"content":"hi"}}`) },
	} {
		srv := httptest.NewServer(h)
		o := &Ollama{URL: srv.URL, Name: "m"}
		if _, err := o.Stream(context.Background(), Request{}, func(string) error { return nil }); !errors.Is(err, ErrProvider) {
			t.Errorf("%s: err = %v, want ErrProvider", name, err)
		}
		srv.Close()
	}
	o := &Ollama{Name: "m"}
	o.URL = "http://127.0.0.1:1" // nothing listens
	if _, err := o.Stream(context.Background(), Request{}, func(string) error { return nil }); !errors.Is(err, ErrProvider) {
		t.Fatalf("unreachable: err = %v", err)
	}
}

func TestQuota(t *testing.T) {
	rdb, prefix := testredis.New(t)
	q := &Quota{Redis: rdb, Prefix: prefix, Limit: 2}
	ctx := context.Background()
	user := uuid.New()
	day := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)

	for want := 1; want >= 0; want-- {
		ok, left, err := q.Take(ctx, user, day)
		if err != nil || !ok || left != want {
			t.Fatalf("take: ok=%v left=%d err=%v, want %d left", ok, left, err, want)
		}
	}
	if ok, _, _ := q.Take(ctx, user, day); ok {
		t.Fatal("third hint must be refused")
	}
	if n, _ := q.Used(ctx, user, day); n != 2 {
		t.Fatalf("used = %d after a refused take, want 2", n)
	}
	if err := q.Refund(ctx, user, day); err != nil {
		t.Fatal(err)
	}
	if ok, _, _ := q.Take(ctx, user, day); !ok {
		t.Fatal("a refunded hint can be taken again")
	}
	if ok, _, _ := q.Take(ctx, user, day.AddDate(0, 0, 1)); !ok {
		t.Fatal("the next day has its own limit")
	}
}
