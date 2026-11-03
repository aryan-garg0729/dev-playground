package auth

import "testing"

func TestSafeNext(t *testing.T) {
	tests := map[string]string{
		"":                     "/",
		"/":                    "/",
		"/learn/a/b":           "/learn/a/b",
		"/courses?x=1":         "/courses?x=1",
		"//evil.example":       "/",
		"/\\evil.example":      "/",
		"https://evil.example": "/",
		"javascript:alert(1)":  "/",
		"learn":                "/",
		"/api/v1/auth/logout":  "/",
		"/a\r\nSet-Cookie: x":  "/",
	}
	for in, want := range tests {
		if got := SafeNext(in); got != want {
			t.Errorf("SafeNext(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalizeEmail(t *testing.T) {
	ok := []struct{ in, want string }{
		{"A@Example.COM", "a@example.com"},
		{"\tb@x.io  ", "b@x.io"}, // surrounding space is trimmed
		{"first.last+tag@sub.example.org", "first.last+tag@sub.example.org"},
	}
	for _, tt := range ok {
		if got, err := NormalizeEmail(tt.in); err != nil || got != tt.want {
			t.Errorf("NormalizeEmail(%q) = %q, %v", tt.in, got, err)
		}
	}
	for _, bad := range []string{"", "a", "a@b", "Name <a@example.com>", "a@example.com, b@example.com", "a@example.com\r\nBcc: x@y.z"} {
		if got, err := NormalizeEmail(bad); err == nil {
			t.Errorf("NormalizeEmail(%q) = %q, want error", bad, got)
		}
	}
}

func TestTokens(t *testing.T) {
	a, ha := NewToken()
	b, _ := NewToken()
	if a == b || len(a) < 40 || Hash(a) != ha || ha == a {
		t.Fatalf("tokens: %q %q hash %q", a, b, ha)
	}
}
