package config

import (
	"strings"
	"testing"
)

func TestLoadTutor(t *testing.T) {
	prod := map[string]string{"APP_ENV": "production", "SMTP_ADDR": "smtp:25", "APP_URL": "https://app.example"}
	tests := []struct {
		name    string
		env     map[string]string
		wantErr string
		check   func(Config) bool
	}{
		{name: "off by default", env: nil, check: func(c Config) bool {
			return c.Tutor == "" && c.TutorDailyLimit == 5 && c.TutorModel == "qwen2.5-coder:3b"
		}},
		{name: "ollama in development", env: map[string]string{"TUTOR_PROVIDER": "ollama", "OLLAMA_URL": "http://o:1", "TUTOR_DAILY_LIMIT": "9"},
			check: func(c Config) bool {
				return c.Tutor == "ollama" && c.OllamaURL == "http://o:1" && c.TutorDailyLimit == 9
			}},
		{name: "unknown provider", env: map[string]string{"TUTOR_PROVIDER": "openai"}, wantErr: "TUTOR_PROVIDER"},
		{name: "bad limit", env: map[string]string{"TUTOR_DAILY_LIMIT": "0"}, wantErr: "TUTOR_DAILY_LIMIT"},
		{name: "production without tutor", env: prod, check: func(c Config) bool { return c.Production && c.Tutor == "" }},
		{name: "production refuses the tutor", env: merge(prod, map[string]string{"TUTOR_PROVIDER": "ollama"}), wantErr: "development-only"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, k := range []string{"APP_ENV", "SMTP_ADDR", "APP_URL", "DEV_AUTH", "TUTOR_PROVIDER", "OLLAMA_URL", "TUTOR_MODEL", "TUTOR_DAILY_LIMIT", "GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET"} {
				t.Setenv(k, "")
			}
			t.Setenv("DATABASE_URL", "postgres://x")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			c, err := Load()
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want one mentioning %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !tt.check(c) {
				t.Fatalf("config = %+v", c)
			}
		})
	}
}

func merge(a, b map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
