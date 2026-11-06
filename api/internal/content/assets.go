package content

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Lesson images (Phase 9). An exercise folder may hold a flat assets/ folder of
// images. Everything in it is public: it is served to the browser as-is, so it
// is checked strictly and every file must be referenced.

// AssetsDir is the exercise subfolder whose files are public images.
const AssetsDir = "assets"

const (
	maxAssetBytes         = 1 << 20 // per file
	maxExerciseAssetBytes = 5 << 20 // per exercise
)

// Asset is one validated, content-addressed image.
type Asset struct {
	Name   string // path relative to the exercise folder, e.g. "assets/loop.svg"
	SHA256 string // hex
	MIME   string
	Ext    string // canonical extension without the dot: png, jpg, webp, gif, svg
	Data   []byte
}

// File is the name the API serves the asset under: "<sha256>.<ext>".
func (a Asset) File() string { return a.SHA256 + "." + a.Ext }

// extMIME maps allowed extensions to the type the bytes must sniff as.
var extMIME = map[string]struct{ mime, ext string }{
	".png":  {"image/png", "png"},
	".jpg":  {"image/jpeg", "jpg"},
	".jpeg": {"image/jpeg", "jpg"},
	".webp": {"image/webp", "webp"},
	".gif":  {"image/gif", "gif"},
	".svg":  {"image/svg+xml", "svg"},
}

// MIMEForExt returns the served type for a canonical extension, if allowed.
func MIMEForExt(ext string) (string, bool) {
	for _, v := range extMIME {
		if v.ext == ext {
			return v.mime, true
		}
	}
	return "", false
}

var serverOnlyNameRe = regexp.MustCompile(`(?i)^(solution|tests_hidden)(\.|$)`)

// loadAssets reads and validates dir/assets/. Problems are reported against
// each file's path; a missing folder means no assets.
func (l *loader) loadAssets(dir string) map[string]Asset {
	out := map[string]Asset{}
	adir := filepath.Join(dir, AssetsDir)
	info, err := os.Lstat(adir)
	if os.IsNotExist(err) {
		return out
	}
	if err != nil {
		l.problem(adir, "%v", err)
		return out
	}
	if !info.IsDir() {
		l.problem(adir, "must be a folder of images")
		return out
	}
	entries, err := os.ReadDir(adir)
	if err != nil {
		l.problem(adir, "%v", err)
		return out
	}
	total := 0
	for _, e := range entries {
		p := filepath.Join(adir, e.Name())
		name := path.Join(AssetsDir, e.Name())
		if strings.HasPrefix(e.Name(), ".") {
			continue // hidden files (.DS_Store) are ignored, like hidden folders
		}
		if serverOnlyNameRe.MatchString(e.Name()) || strings.EqualFold(e.Name(), "wrong") {
			l.problem(p, "server-only files (solution.*, tests_hidden.*, wrong/) must never be in assets/: everything there is public")
			continue
		}
		if e.Type()&os.ModeSymlink != 0 {
			l.problem(p, "symlinks are not allowed in assets/")
			continue
		}
		if e.IsDir() {
			l.problem(p, "assets/ must be flat (no subfolders)")
			continue
		}
		if !e.Type().IsRegular() {
			l.problem(p, "not a regular file")
			continue
		}
		kind, ok := extMIME[strings.ToLower(filepath.Ext(e.Name()))]
		if !ok {
			l.problem(p, "unsupported image type (allowed: png, jpg, jpeg, webp, gif, svg)")
			continue
		}
		data, err := os.ReadFile(p)
		if err != nil {
			l.problem(p, "%v", err)
			continue
		}
		if len(data) > maxAssetBytes {
			l.problem(p, "is %d bytes; images must be at most %d bytes", len(data), maxAssetBytes)
			continue
		}
		total += len(data)
		if reason := checkImageBytes(kind.mime, data); reason != "" {
			l.problem(p, "%s", reason)
			continue
		}
		sum := sha256.Sum256(data)
		out[name] = Asset{Name: name, SHA256: hex.EncodeToString(sum[:]), MIME: kind.mime, Ext: kind.ext, Data: data}
	}
	if total > maxExerciseAssetBytes {
		l.problem(adir, "images total %d bytes; an exercise may have at most %d", total, maxExerciseAssetBytes)
	}
	return out
}

// checkImageBytes confirms the content matches the extension; "" when fine.
func checkImageBytes(mime string, data []byte) string {
	if mime == "image/svg+xml" {
		return checkSVG(data)
	}
	if got := http.DetectContentType(data); got != mime {
		return fmt.Sprintf("content is %s, not %s (the extension does not match the file)", got, mime)
	}
	return ""
}

var (
	svgOpenRe     = regexp.MustCompile(`(?is)^\s*(<\?xml[^>]*>\s*)?(<!--.*?-->\s*)*(<!doctype svg[^>]*>\s*)?(<!--.*?-->\s*)*<svg[\s>]`)
	svgEventRe    = regexp.MustCompile(`(?i)\son[a-z]+\s*=`)
	svgHrefRe     = regexp.MustCompile(`(?i)\b(?:xlink:)?href\s*=\s*["']\s*([^"']*)["']`)
	svgURLRe      = regexp.MustCompile(`(?i)url\(\s*['"]?\s*([^'")\s]*)`)
	svgForbidden  = []string{"<script", "<foreignobject", "<iframe", "<embed", "<object", "<!entity", "javascript:", "@import"}
	svgOKHrefKind = func(v string) bool { return strings.HasPrefix(strings.TrimSpace(v), "#") }
)

// checkSVG rejects SVGs that could run script or load anything from outside.
// This is defence in depth: the API also serves SVGs under a sandbox CSP and
// the app only shows them through <img>, where scripts never run.
func checkSVG(data []byte) string {
	if !svgOpenRe.Match(data) {
		return "not an SVG document (expected an <svg> root element)"
	}
	low := bytes.ToLower(data)
	for _, f := range svgForbidden {
		if bytes.Contains(low, []byte(f)) {
			return fmt.Sprintf("SVG contains %q, which is not allowed", f)
		}
	}
	if svgEventRe.Match(data) {
		return "SVG contains an event handler attribute (on...=), which is not allowed"
	}
	for _, m := range svgHrefRe.FindAllSubmatch(data, -1) {
		if !svgOKHrefKind(string(m[1])) {
			return fmt.Sprintf("SVG links to %q; only #fragment links inside the file are allowed", m[1])
		}
	}
	for _, m := range svgURLRe.FindAllSubmatch(data, -1) {
		if !svgOKHrefKind(string(m[1])) {
			return fmt.Sprintf("SVG uses url(%s); only url(#id) inside the file is allowed", m[1])
		}
	}
	return ""
}

// mdImageRe finds markdown images: ![alt](target "title") and ![alt](<target>).
// Alt text may hold one level of brackets, e.g. ![nums [1, 2, 3]](assets/a.svg).
var mdImageRe = regexp.MustCompile(`!\[(?:[^\[\]]|\[[^\[\]]*\])*\]\(\s*(<[^>]*>|[^)\s]+)`)

// checkImageRefs makes every image a lesson uses point at one of its assets,
// and every asset be used. Uses come from markdown images in the exercise's
// public .md files (not solution.md) and from spec strings that start with "assets/" (e.g. mcq option images).
func (l *loader) checkImageRefs(dir string, spec yaml.Node, assets map[string]Asset) {
	used := map[string]bool{}
	mds, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	sort.Strings(mds)
	for _, md := range mds {
		if serverOnlyNameRe.MatchString(filepath.Base(md)) {
			continue // server-only markdown never reaches the browser, so it cannot use public images
		}
		body, err := os.ReadFile(md)
		if err != nil {
			l.problem(md, "%v", err)
			continue
		}
		for _, m := range mdImageRe.FindAllSubmatch(body, -1) {
			target := strings.TrimSuffix(strings.TrimPrefix(string(m[1]), "<"), ">")
			l.checkRef(md, target, assets, used)
		}
	}
	for _, s := range specStrings(&spec) {
		if strings.HasPrefix(s, AssetsDir+"/") {
			l.checkRef(filepath.Join(dir, "exercise.yaml"), s, assets, used)
		}
	}
	names := make([]string, 0, len(assets))
	for n := range assets {
		names = append(names, n)
	}
	slices.Sort(names)
	for _, n := range names {
		if !used[n] {
			l.problem(filepath.Join(dir, filepath.FromSlash(n)), "not used by the lesson; remove it (everything in assets/ is public)")
		}
	}
}

func (l *loader) checkRef(file, target string, assets map[string]Asset, used map[string]bool) {
	if strings.Contains(target, ":") || strings.HasPrefix(target, "//") {
		l.problem(file, "image %q: external images are not allowed; put the file in %s/ and link it as %s/<name>", target, AssetsDir, AssetsDir)
		return
	}
	clean := path.Clean(strings.TrimPrefix(target, "./"))
	if _, ok := assets[clean]; !ok {
		l.problem(file, "image %q: no such file in %s/", target, AssetsDir)
		return
	}
	used[clean] = true
}

// specStrings returns every scalar string value under n.
func specStrings(n *yaml.Node) []string {
	if n == nil {
		return nil
	}
	if n.Kind == yaml.ScalarNode {
		return []string{n.Value}
	}
	var out []string
	for i, c := range n.Content {
		// In a mapping, skip keys (even indexes); only values can be paths.
		if n.Kind == yaml.MappingNode && i%2 == 0 {
			continue
		}
		out = append(out, specStrings(c)...)
	}
	return out
}

// AssetMap is the exercise's public name -> served file map ("assets/x.svg" -> "<sha>.svg").
func AssetMap(assets []Asset) map[string]string {
	m := make(map[string]string, len(assets))
	for _, a := range assets {
		m[a.Name] = a.File()
	}
	return m
}
