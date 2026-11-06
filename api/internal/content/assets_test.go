package content

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A 1x1 PNG.
var pngBytes = string([]byte{
	0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0, 0, 0, 0x0d, 'I', 'H', 'D', 'R',
	0, 0, 0, 1, 0, 0, 0, 1, 8, 6, 0, 0, 0, 0x1f, 0x15, 0xc4, 0x89,
})

const goodSVG = `<?xml version="1.0"?>
<!-- drawn in Excalidraw -->
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><defs><marker id="m"/></defs>
<path d="M0 0h10" marker-end="url(#m)"/><use href="#m"/></svg>`

// withImages is the base tree plus one exercise that uses images.
func withImages(extra map[string]string) map[string]string {
	m := baseTree()
	m["c/ch1/a/prompt.md"] = "Look:\n\n![a loop over [1, 2]](assets/loop.svg)\n\n![dot](<./assets/dot.png> \"title\")\n"
	m["c/ch1/a/assets/loop.svg"] = goodSVG
	m["c/ch1/a/assets/dot.png"] = pngBytes
	for k, v := range extra {
		m[k] = v
	}
	return m
}

func TestLoadAssets(t *testing.T) {
	b, err := Load(writeTree(t, withImages(nil)))
	if err != nil {
		t.Fatal(err)
	}
	var got []Asset
	for _, ch := range b.Courses[0].Chapters {
		for _, e := range ch.Exercises {
			if e.ID == "c.one.a" {
				got = e.Assets
			}
		}
	}
	if len(got) != 2 || got[0].Name != "assets/dot.png" || got[1].Name != "assets/loop.svg" {
		t.Fatalf("assets = %+v", got)
	}
	if got[0].MIME != "image/png" || got[1].MIME != "image/svg+xml" || len(got[0].SHA256) != 64 {
		t.Fatalf("bad metadata: %+v %+v", got[0], got[1])
	}
	m := AssetMap(got)
	if m["assets/loop.svg"] != got[1].SHA256+".svg" {
		t.Fatalf("map = %v", m)
	}
}

func TestMCQOptionImageCountsAsUse(t *testing.T) {
	files := baseTree()
	files["c/ch1/a/exercise.yaml"] = strings.Replace(ex("c.one.a"), `{ id: a, text: "A" }`, `{ id: a, text: "A", image: assets/dot.png }`, 1)
	files["c/ch1/a/assets/dot.png"] = pngBytes
	if _, err := Load(writeTree(t, files)); err != nil {
		t.Fatal(err)
	}
}

func TestAssetProblems(t *testing.T) {
	big := strings.Repeat("x", maxAssetBytes+1)
	tests := []struct {
		name     string
		files    map[string]string
		wantPath string
		wantMsg  string
	}{
		{"missing image", map[string]string{"c/ch1/a/prompt.md": "![x](assets/nope.png) ![a](assets/loop.svg) ![d](assets/dot.png)"},
			"c/ch1/a/prompt.md", "no such file"},
		{"external image", map[string]string{"c/ch1/a/prompt.md": "![x](https://example.com/a.png) ![a](assets/loop.svg) ![d](assets/dot.png)"},
			"c/ch1/a/prompt.md", "external images are not allowed"},
		{"protocol-relative image", map[string]string{"c/ch1/a/prompt.md": "![x](//cdn.example.com/a.png) ![a](assets/loop.svg) ![d](assets/dot.png)"},
			"c/ch1/a/prompt.md", "external images are not allowed"},
		{"unused asset", map[string]string{"c/ch1/a/assets/extra.png": pngBytes},
			"c/ch1/a/assets/extra.png", "not used"},
		{"over-size file", map[string]string{"c/ch1/a/assets/dot.png": "\x89PNG\r\n\x1a\n" + big},
			"c/ch1/a/assets/dot.png", "at most"},
		{"disallowed type", map[string]string{"c/ch1/a/assets/notes.txt": "hi"},
			"c/ch1/a/assets/notes.txt", "unsupported image type"},
		{"extension lies", map[string]string{"c/ch1/a/assets/dot.png": "GIF89a......"},
			"c/ch1/a/assets/dot.png", "not image/png"},
		{"svg with script", map[string]string{"c/ch1/a/assets/loop.svg": `<svg><script>alert(1)</script></svg>`},
			"c/ch1/a/assets/loop.svg", "<script"},
		{"svg with handler", map[string]string{"c/ch1/a/assets/loop.svg": `<svg onload="x()"></svg>`},
			"c/ch1/a/assets/loop.svg", "event handler"},
		{"svg external link", map[string]string{"c/ch1/a/assets/loop.svg": `<svg><image href="https://x.test/a.png"/></svg>`},
			"c/ch1/a/assets/loop.svg", "only #fragment"},
		{"svg external url()", map[string]string{"c/ch1/a/assets/loop.svg": `<svg><rect style="fill:url(https://x.test/a)"/></svg>`},
			"c/ch1/a/assets/loop.svg", "url(#id)"},
		{"svg foreignObject", map[string]string{"c/ch1/a/assets/loop.svg": `<svg><foreignObject><div/></foreignObject></svg>`},
			"c/ch1/a/assets/loop.svg", "foreignobject"},
		{"not an svg", map[string]string{"c/ch1/a/assets/loop.svg": `<html><body/></html>`},
			"c/ch1/a/assets/loop.svg", "not an SVG"},
		{"server-only name", map[string]string{"c/ch1/a/assets/solution.py": "print(1)"},
			"c/ch1/a/assets/solution.py", "server-only"},
		{"hidden tests name", map[string]string{"c/ch1/a/assets/tests_hidden.png": pngBytes},
			"c/ch1/a/assets/tests_hidden.png", "server-only"},
		{"wrong folder", map[string]string{"c/ch1/a/assets/wrong/a.png": pngBytes},
			"c/ch1/a/assets/wrong", "server-only"},
		{"subfolder", map[string]string{"c/ch1/a/assets/more/a.png": pngBytes},
			"c/ch1/a/assets/more", "flat"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := writeTree(t, withImages(tt.files))
			_, err := Load(root)
			var verr *ValidationError
			if !errors.As(err, &verr) {
				t.Fatalf("want ValidationError, got %v", err)
			}
			want := filepath.Join(root, tt.wantPath)
			for _, p := range verr.Problems {
				if p.Path == want && strings.Contains(p.Reason, tt.wantMsg) {
					return
				}
			}
			t.Fatalf("no problem at %s containing %q; got:\n%v", tt.wantPath, tt.wantMsg, verr)
		})
	}
}

func TestAssetSymlinkRejected(t *testing.T) {
	root := writeTree(t, withImages(map[string]string{"secret.png": pngBytes}))
	link := filepath.Join(root, "c/ch1/a/assets/link.png")
	if err := os.Symlink(filepath.Join(root, "secret.png"), link); err != nil {
		t.Skip("symlinks unavailable:", err)
	}
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "symlinks are not allowed") {
		t.Fatalf("got %v", err)
	}
}

func TestExerciseAssetTotalLimit(t *testing.T) {
	extra := map[string]string{}
	prompt := "![a](assets/loop.svg) ![d](assets/dot.png)"
	for _, n := range []string{"a", "b", "c", "d", "e", "f"} {
		extra["c/ch1/a/assets/"+n+".svg"] = `<svg>` + strings.Repeat(" ", maxAssetBytes-20) + `</svg>`
		prompt += " ![x](assets/" + n + ".svg)"
	}
	extra["c/ch1/a/prompt.md"] = prompt
	_, err := Load(writeTree(t, withImages(extra)))
	if err == nil || !strings.Contains(err.Error(), "an exercise may have at most") {
		t.Fatalf("got %v", err)
	}
}
