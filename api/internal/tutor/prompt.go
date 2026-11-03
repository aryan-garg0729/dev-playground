package tutor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// systemPrompt sets the tutor's behaviour (SPEC "AI tutor").
const systemPrompt = `You are the CodeQuest tutor, helping a learner with one coding exercise.
Teach, don't solve:
- Answer with questions and small nudges that help the learner find the fix themselves.
- Point to the line or the concept involved, not the fix.
- Never write the solution. Show at most 3 consecutive lines of code, and only to illustrate a concept.
- If the learner asks for the answer, refuse in one short sentence and give the next hint instead.
- Be brief: a few sentences. Use Markdown; put code in backticks.
- Use only the exercise, the learner's code and their latest result below. Do not invent requirements.`

// Exercise is what the tutor knows about the lesson: only what the learner sees.
type Exercise struct {
	Title    string
	Type     string          // code | sql | mcq
	Public   json.RawMessage // the exercise's public_view
	Code     string          // the learner's current answer (code or SQL), if any
	Result   json.RawMessage // the learner's latest Result, if any
	History  []Message       // earlier messages, oldest first (at most HistoryLen)
	Question string          // the learner's new message
}

// publicView is the part of a public_view the tutor uses; fields differ by type.
type publicView struct {
	Prompt   string `json:"prompt"`
	Question string `json:"question"`
	Language string `json:"language"`
	Setup    string `json:"setup"`
}

type resultView struct {
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Checks []struct {
		Name     string `json:"name"`
		Passed   bool   `json:"passed"`
		Feedback string `json:"feedback"`
		Got      string `json:"got"`
		Expected string `json:"expected"`
		Line     int    `json:"line"`
	} `json:"checks"`
}

// BuildRequest turns an exercise and a question into a model request.
func BuildRequest(e Exercise) Request {
	var pv publicView
	_ = json.Unmarshal(e.Public, &pv) // a malformed view just means less context

	var b strings.Builder
	fmt.Fprintf(&b, "Exercise: %s (%s)\n\n", e.Title, e.Type)
	if p := firstNonEmpty(pv.Prompt, pv.Question); p != "" {
		fmt.Fprintf(&b, "## Task\n%s\n\n", strings.TrimSpace(p))
	}
	if pv.Setup != "" {
		fmt.Fprintf(&b, "## Database setup\n```sql\n%s\n```\n\n", clip(pv.Setup, 4000))
	}
	lang := firstNonEmpty(pv.Language, e.Type)
	if strings.TrimSpace(e.Code) != "" {
		fmt.Fprintf(&b, "## Learner's current code\n```%s\n%s\n```\n\n", lang, numbered(clip(e.Code, 12000)))
	}
	b.WriteString(describeResult(e.Result))

	sys := systemPrompt + "\n\n" + b.String()
	msgs := append([]Message(nil), e.History...)
	msgs = append(msgs, Message{Role: "user", Content: e.Question})
	return Request{System: sys, Messages: msgs, MaxTokens: MaxOutputTokens}
}

func describeResult(raw json.RawMessage) string {
	var r resultView
	if len(raw) == 0 || json.Unmarshal(raw, &r) != nil || r.Status == "" {
		return "## Latest result\nThe learner has not submitted yet.\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "## Latest result: %s\n", r.Status)
	for _, c := range r.Checks {
		if c.Passed {
			fmt.Fprintf(&b, "- PASS %s\n", c.Name)
			continue
		}
		fmt.Fprintf(&b, "- FAIL %s", c.Name)
		if c.Line > 0 {
			fmt.Fprintf(&b, " (line %d)", c.Line)
		}
		if c.Got != "" || c.Expected != "" {
			fmt.Fprintf(&b, ": got %s, expected %s", c.Got, c.Expected)
		}
		if c.Feedback != "" {
			fmt.Fprintf(&b, " — %s", clip(c.Feedback, 600))
		}
		b.WriteString("\n")
	}
	if s := strings.TrimSpace(r.Stdout); s != "" {
		fmt.Fprintf(&b, "Output:\n```\n%s\n```\n", clip(s, 2000))
	}
	return b.String()
}

// numbered prefixes lines with their numbers so the tutor can point at one.
func numbered(code string) string {
	lines := strings.Split(strings.TrimRight(code, "\n"), "\n")
	for i, l := range lines {
		lines[i] = fmt.Sprintf("%3d | %s", i+1, l)
	}
	return strings.Join(lines, "\n")
}

func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "\n…(truncated)"
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if strings.TrimSpace(s) != "" {
			return s
		}
	}
	return ""
}
