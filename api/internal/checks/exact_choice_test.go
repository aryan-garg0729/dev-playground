package checks

import (
	"context"
	"encoding/json"
	"testing"
)

func TestExactChoice(t *testing.T) {
	cfg, _ := json.Marshal(ExactChoiceConfig{Answer: []string{"a", "c"}, Feedback: "because"})
	tests := []struct {
		name    string
		answer  string
		passed  bool
		wantErr bool
	}{
		{"exact set passes", `{"choices":["a","c"]}`, true, false},
		{"order ignored", `{"choices":["c","a"]}`, true, false},
		{"duplicates ignored", `{"choices":["a","c","a"]}`, true, false},
		{"subset fails", `{"choices":["a"]}`, false, false},
		{"superset fails", `{"choices":["a","b","c"]}`, false, false},
		{"empty fails", `{"choices":[]}`, false, false},
		{"bad json errors", `{"choices":`, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Run(context.Background(), "exact_choice", CheckInput{Config: cfg, Answer: json.RawMessage(tt.answer)})
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(res) != 1 || res[0].Passed != tt.passed {
				t.Fatalf("got %+v, want passed=%v", res, tt.passed)
			}
			if res[0].Feedback != "because" || !res[0].Visible {
				t.Fatalf("feedback/visible not carried: %+v", res[0])
			}
		})
	}
}

func TestUnknownKind(t *testing.T) {
	if _, err := Run(context.Background(), "nope", CheckInput{}); err == nil {
		t.Fatal("want error for unknown kind")
	}
}
