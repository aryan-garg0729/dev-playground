package main

// TEMPORARY helper for authoring content: runs a lesson's steps in a folder
// and prints the recorded answer JSON. Delete before committing (copied from the codequest-course-writer skill).
//   CQ_REC_YAML=... CQ_REC_DIR=... go test ./cmd/codequest -run TestZZRecord -v

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/aryangarg/codequest/api/internal/exercise/projectcli"
	"go.yaml.in/yaml/v3"
)

func TestZZRecord(t *testing.T) {
	path := os.Getenv("CQ_REC_YAML")
	if path == "" {
		t.Skip("no CQ_REC_YAML")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Spec struct {
			BaseURL string             `yaml:"base_url"`
			Steps   []projectcli.Step `yaml:"steps"`
		} `yaml:"spec"`
	}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	sr := newStepRunner()
	sr.dir = os.Getenv("CQ_REC_DIR")
	a := sr.run(context.Background(), projectcli.Spec{BaseURL: doc.Spec.BaseURL, Steps: doc.Spec.Steps}, "")
	type cli struct {
		ExitCode int    `json:"exit_code"`
		Stdout   string `json:"stdout"`
	}
	type ht struct {
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	}
	type st struct {
		CLI  *cli `json:"cli,omitempty"`
		HTTP *ht  `json:"http,omitempty"`
	}
	out := struct {
		Steps []st `json:"steps"`
	}{}
	for _, r := range a.Steps {
		switch {
		case r.CLI != nil:
			out.Steps = append(out.Steps, st{CLI: &cli{r.CLI.ExitCode, r.CLI.Stdout}})
		case r.HTTP != nil:
			h := map[string]string{}
			if ct, ok := r.HTTP.Headers["Content-Type"]; ok {
				h["Content-Type"] = ct
			}
			out.Steps = append(out.Steps, st{HTTP: &ht{r.HTTP.Status, h, r.HTTP.Body}})
		default:
			out.Steps = append(out.Steps, st{CLI: &cli{-1, "[skipped]"}})
		}
	}
	b, _ := json.Marshal(out)
	if err := os.WriteFile(os.Getenv("CQ_REC_OUT"), append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
