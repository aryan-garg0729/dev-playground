package httpapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aryangarg/codequest/api/internal/httpapi"
)

// The web build is served next to the API: files as is, client routes get
// index.html, /api goes to the API, and a missing hashed asset is a 404.
func TestWithWeb(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"index.html": "<html>app</html>", "assets/app-1.js": "js", "favicon.svg": "svg"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	api := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, "api") })
	h, err := httpapi.WithWeb(api, dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		method, path string
		code         int
		body         string
	}{
		{"GET", "/api/v1/healthz", 200, "api"},
		{"POST", "/api/v1/shop/peek/buy", 200, "api"},
		{"GET", "/", 200, "<html>app</html>"},
		{"GET", "/courses/git", 200, "<html>app</html>"},
		{"GET", "/assets/app-1.js", 200, "js"},
		{"GET", "/assets/gone.js", 404, ""},
		{"GET", "/favicon.svg", 200, "svg"},
		{"GET", "/../../etc/passwd", 200, "<html>app</html>"},
		{"POST", "/courses", 405, ""},
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(c.method, c.path, nil))
		if rec.Code != c.code || (c.body != "" && rec.Body.String() != c.body) {
			t.Errorf("%s %s = %d %q, want %d %q", c.method, c.path, rec.Code, rec.Body.String(), c.code, c.body)
		}
	}
	if _, err := httpapi.WithWeb(api, t.TempDir()); err == nil {
		t.Error("a directory without index.html was accepted")
	}
}
