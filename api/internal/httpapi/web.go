package httpapi

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// WithWeb serves the built web app (web/dist) from dir in front of the API, so
// app and API share one origin and the session cookie stays first-party (on
// hosts like Render, two *.onrender.com services are different sites).
// /api/... goes to api. A file that exists is served as is; Vite's hashed
// /assets/ files are cached for a year and a missing one is a 404. Any other
// path is a client-side route and gets index.html.
func WithWeb(api http.Handler, dir string) (http.Handler, error) {
	index, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		return nil, err
	}
	files := http.FileServer(http.Dir(dir))
	started := time.Now()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
			api.ServeHTTP(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "same-origin")
		h.Set("X-Frame-Options", "DENY")
		clean := path.Clean("/" + r.URL.Path)
		if st, err := os.Stat(filepath.Join(dir, filepath.FromSlash(clean))); err == nil && !st.IsDir() {
			if strings.HasPrefix(clean, "/assets/") {
				h.Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			files.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(clean, "/assets/") {
			http.NotFound(w, r)
			return
		}
		h.Set("Cache-Control", "no-cache")
		h.Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", started, bytes.NewReader(index))
	}), nil
}
