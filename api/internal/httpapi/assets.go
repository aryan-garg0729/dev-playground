package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/content"
	"github.com/aryangarg/codequest/api/internal/store"
)

// AssetsPath is the URL prefix lesson images are served under.
const AssetsPath = "/api/v1/assets/"

var assetFileRe = regexp.MustCompile(`^([0-9a-f]{64})\.([a-z]{3,4})$`)

// getAsset serves a lesson image by content hash. Bodies never change for a
// hash, so they cache forever. Even though cq build rejects scriptable SVGs,
// every response carries a sandbox CSP and nosniff, so an image opened directly
// (not through <img>) can never run script on the API's origin.
func (a *api) getAsset(w http.ResponseWriter, r *http.Request) {
	m := assetFileRe.FindStringSubmatch(chi.URLParam(r, "file"))
	if m == nil {
		writeError(w, http.StatusNotFound, "not_found", "asset not found")
		return
	}
	sha, ext := m[1], m[2]
	want, ok := content.MIMEForExt(ext)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "asset not found")
		return
	}
	etag := `"` + sha + `"`
	if r.Header.Get("If-None-Match") == etag {
		setAssetHeaders(w, etag)
		w.WriteHeader(http.StatusNotModified)
		return
	}
	row, err := store.New(a.Pool).GetContentAsset(r.Context(), sha)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && row.Mime != want) {
		writeError(w, http.StatusNotFound, "not_found", "asset not found")
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	setAssetHeaders(w, etag)
	w.Header().Set("Content-Type", row.Mime)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(row.Bytes)
}

func setAssetHeaders(w http.ResponseWriter, etag string) {
	h := w.Header()
	h.Set("Cache-Control", "public, max-age=31536000, immutable")
	h.Set("ETag", etag)
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")
	h.Set("Cross-Origin-Resource-Policy", "same-site")
	h.Set("Content-Disposition", "inline")
}

// assetURLs turns the stored name -> file map into name -> URL for the browser.
func assetURLs(raw json.RawMessage) (map[string]string, error) {
	files := map[string]string{}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &files); err != nil {
			return nil, err
		}
	}
	out := make(map[string]string, len(files))
	for name, file := range files {
		out[name] = AssetsPath + file
	}
	return out, nil
}
