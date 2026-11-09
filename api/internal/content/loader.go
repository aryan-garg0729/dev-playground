package content

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/aryangarg/codequest/api/internal/exercise"
)

// Load reads and validates a content tree. root is either a folder of courses
// (each subfolder holds a course.yaml) or a single course folder. On any
// problem it returns a *ValidationError listing every path and reason.
func Load(root string) (*Bundle, error) {
	l := &loader{seen: map[string]string{}}
	b := &Bundle{}

	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("content root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("content root %s is not a directory", root)
	}

	if fileExists(filepath.Join(root, "course.yaml")) {
		if c, ok := l.course(root); ok {
			b.Courses = append(b.Courses, c)
		}
	} else {
		dirs, err := subdirs(root)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", root, err)
		}
		if len(dirs) == 0 {
			l.problem(root, "no course folders found (each needs a course.yaml)")
		}
		for _, d := range dirs {
			dir := filepath.Join(root, d)
			if !fileExists(filepath.Join(dir, "course.yaml")) {
				l.problem(dir, "folder has no course.yaml")
				continue
			}
			if c, ok := l.course(dir); ok {
				b.Courses = append(b.Courses, c)
			}
		}
	}

	if fileExists(filepath.Join(root, "paths.yaml")) && !fileExists(filepath.Join(root, "course.yaml")) {
		b.Paths = l.paths(filepath.Join(root, "paths.yaml"), b.Courses)
	}

	if len(l.problems) > 0 {
		return nil, &ValidationError{Problems: l.problems}
	}
	slices.SortStableFunc(b.Courses, func(a, c Course) int {
		if a.Order != c.Order {
			return a.Order - c.Order
		}
		return strings.Compare(a.ID, c.ID)
	})
	return b, nil
}

type loader struct {
	problems []Problem
	seen     map[string]string // "kind:id" -> path that declared it
}

func (l *loader) problem(path, format string, args ...any) {
	l.problems = append(l.problems, Problem{Path: path, Reason: fmt.Sprintf(format, args...)})
}

// claim records a globally unique id; it reports a duplicate and returns false.
func (l *loader) claim(kind, id, path string) bool {
	key := kind + ":" + id
	if prev, dup := l.seen[key]; dup {
		l.problem(path, "%s id %q is already used by %s", kind, id, prev)
		return false
	}
	l.seen[key] = path
	return true
}

func (l *loader) course(dir string) (Course, bool) {
	path := filepath.Join(dir, "course.yaml")
	var m courseManifest
	if err := readStrict(path, &m); err != nil {
		l.problem(path, "%v", err)
		return Course{}, false
	}
	ok := true
	if !courseIDRe.MatchString(m.ID) {
		l.problem(path, "id %q must be a lowercase slug like python-basics", m.ID)
		ok = false
	} else if !l.claim("course", m.ID, path) {
		ok = false
	}
	if strings.TrimSpace(m.Title) == "" {
		l.problem(path, "title is required")
		ok = false
	}
	c := Course{ID: m.ID, Title: m.Title, Description: strings.TrimSpace(m.Description), Order: m.Order, Published: m.Published == nil || *m.Published}
	ok = l.catalog(path, m, &c) && ok

	folders, listOK := l.checkListing(dir, path, "chapters", m.Chapters, "chapter.yaml")
	ok = ok && listOK
	for _, f := range folders {
		ch, chOK := l.chapter(filepath.Join(dir, f))
		ok = ok && chOK
		c.Chapters = append(c.Chapters, ch)
	}
	return c, ok
}

func (l *loader) chapter(dir string) (Chapter, bool) {
	path := filepath.Join(dir, "chapter.yaml")
	var m chapterManifest
	if err := readStrict(path, &m); err != nil {
		l.problem(path, "%v", err)
		return Chapter{}, false
	}
	ok := true
	if !dottedIDRe.MatchString(m.ID) {
		l.problem(path, "id %q must be lowercase words joined by . or - (e.g. lc.chat-models)", m.ID)
		ok = false
	} else if !l.claim("chapter", m.ID, path) {
		ok = false
	}
	if strings.TrimSpace(m.Title) == "" {
		l.problem(path, "title is required")
		ok = false
	}
	ch := Chapter{ID: m.ID, Title: m.Title}
	folders, listOK := l.checkListing(dir, path, "exercises", m.Exercises, "exercise.yaml")
	ok = ok && listOK
	for _, f := range folders {
		ex, exOK := l.exercise(filepath.Join(dir, f))
		ok = ok && exOK
		ch.Exercises = append(ch.Exercises, ex)
	}
	return ch, ok
}

// checkListing verifies a manifest's ordered folder list against the folders on
// disk: every listed folder must exist with its manifest, and every folder on
// disk must be listed. It returns the listed folders that exist.
func (l *loader) checkListing(dir, manifestPath, field string, listed []string, childManifest string) ([]string, bool) {
	ok := true
	if len(listed) == 0 {
		l.problem(manifestPath, "%s: list at least one folder", field)
		return nil, false
	}
	var found []string
	seen := map[string]bool{}
	for _, f := range listed {
		if seen[f] {
			l.problem(manifestPath, "%s: %q is listed twice", field, f)
			ok = false
			continue
		}
		seen[f] = true
		if f == "" || f != filepath.Base(f) || strings.HasPrefix(f, ".") {
			l.problem(manifestPath, "%s: %q must be a folder name in %s", field, f, dir)
			ok = false
			continue
		}
		if !fileExists(filepath.Join(dir, f, childManifest)) {
			l.problem(manifestPath, "%s: %q has no %s", field, f, filepath.Join(dir, f, childManifest))
			ok = false
			continue
		}
		found = append(found, f)
	}
	onDisk, err := subdirs(dir)
	if err != nil {
		l.problem(dir, "read folder: %v", err)
		return found, false
	}
	for _, d := range onDisk {
		if !seen[d] {
			l.problem(filepath.Join(dir, d), "folder is not listed in %s %s", filepath.Base(manifestPath), field)
			ok = false
		}
	}
	return found, ok
}

func (l *loader) exercise(dir string) (Compiled, bool) {
	path := filepath.Join(dir, "exercise.yaml")
	var h exerciseHeader
	if err := readStrict(path, &h); err != nil {
		l.problem(path, "%v", err)
		return Compiled{}, false
	}
	n := len(l.problems)
	if !dottedIDRe.MatchString(h.ID) {
		l.problem(path, "id %q must be lowercase words joined by . or - (e.g. py.loops.sum-list)", h.ID)
	} else {
		l.claim("exercise", h.ID, path)
	}
	if strings.TrimSpace(h.Title) == "" {
		l.problem(path, "title is required")
	}
	if h.Version < 1 {
		l.problem(path, "version must be 1 or higher")
	}
	if h.XP == nil || *h.XP < 0 {
		l.problem(path, "xp is required and must be 0 or more")
	}
	if h.Difficulty == 0 {
		h.Difficulty = 1
	}
	if h.Difficulty < 1 || h.Difficulty > 5 {
		l.problem(path, "difficulty must be 1-5, got %d", h.Difficulty)
	}
	typ, known := exercise.Get(h.Type)
	if !known {
		l.problem(path, "type %q is not registered (known: %s)", h.Type, strings.Join(exercise.Names(), ", "))
	}
	if h.Boss && h.Type != "code" && h.Type != "project_cli" {
		l.problem(path, "boss: only code and project_cli exercises can be bosses")
	}
	if !known {
		return Compiled{}, false
	}

	spec, err := typ.ParseSpec(h.Spec, os.DirFS(dir))
	if err != nil {
		l.problem(path, "spec: %v", err)
		return Compiled{}, false
	}
	assets := l.loadAssets(dir)
	l.checkImageRefs(dir, h.Spec, assets)
	if len(l.problems) > n {
		return Compiled{}, false
	}

	xp := 0
	if h.XP != nil {
		xp = *h.XP
	}
	tags := h.Tags
	if tags == nil {
		tags = []string{}
	}
	ex := exercise.Exercise{
		ID: h.ID, Type: h.Type, Version: h.Version, Title: h.Title, XP: xp,
		Difficulty: h.Difficulty, Tags: tags, Free: h.Free, Boss: h.Boss, Spec: spec,
	}
	pub, err := typ.PublicView(ex)
	if err != nil {
		l.problem(path, "public view: %v", err)
		return Compiled{}, false
	}
	priv, err := json.Marshal(spec)
	if err != nil {
		l.problem(path, "encode spec: %v", err)
		return Compiled{}, false
	}
	if _, err := typ.LoadSpec(priv); err != nil {
		l.problem(path, "spec does not round-trip through JSON: %v", err)
		return Compiled{}, false
	}
	list := make([]Asset, 0, len(assets))
	for _, a := range assets {
		list = append(list, a)
	}
	slices.SortFunc(list, func(a, b Asset) int { return strings.Compare(a.Name, b.Name) })
	return Compiled{Exercise: ex, PublicView: pub, PrivateSpec: priv, Assets: list}, true
}

// readStrict decodes one YAML document into v, rejecting unknown fields.
func readStrict(path string, v any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("file is empty")
		}
		return err
	}
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// subdirs lists non-hidden subfolders, sorted.
func subdirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

var (
	iconRe     = regexp.MustCompile(`^[a-z0-9]{1,4}$`)
	courseTags = 4
	maxTopics  = 60
	levels     = map[string]bool{"beginner": true, "intermediate": true, "advanced": true}
	// pathGoals are the onboarding goals a path can belong to (users.goal).
	pathGoals = map[string]bool{"zero": true, "backend": true, "interviews": true, "ai": true}
)

// catalog validates the course's catalog fields (icon, level, tags, topics)
// and copies them into c, filling defaults.
func (l *loader) catalog(path string, m courseManifest, c *Course) bool {
	ok := true
	c.Icon = m.Icon
	if c.Icon == "" {
		c.Icon = strings.ReplaceAll(m.ID, "-", "")
		c.Icon = c.Icon[:min(2, len(c.Icon))]
	} else if !iconRe.MatchString(m.Icon) {
		l.problem(path, "icon %q must be 1-4 lowercase letters or digits", m.Icon)
		ok = false
	}
	c.Level = m.Level
	if c.Level == "" {
		c.Level = "beginner"
	} else if !levels[m.Level] {
		l.problem(path, "level %q must be beginner, intermediate or advanced", m.Level)
		ok = false
	}
	if len(m.Tags) > courseTags {
		l.problem(path, "tags: at most %d", courseTags)
		ok = false
	}
	c.Tags = []string{}
	seen := map[string]bool{}
	for _, t := range m.Tags {
		t = strings.TrimSpace(t)
		switch {
		case t == "" || len(t) > 20:
			l.problem(path, "tags: %q must be 1-20 characters", t)
			ok = false
		case seen[strings.ToLower(t)]:
			l.problem(path, "tags: %q is listed twice", t)
			ok = false
		default:
			seen[strings.ToLower(t)] = true
			c.Tags = append(c.Tags, t)
		}
	}
	c.Topics = strings.TrimSpace(m.Topics)
	if len(c.Topics) > maxTopics {
		l.problem(path, "topics: at most %d characters", maxTopics)
		ok = false
	}
	return ok
}

// paths reads paths.yaml: career paths made of the courses in this tree.
func (l *loader) paths(path string, courses []Course) []Path {
	var m pathsManifest
	if err := readStrict(path, &m); err != nil {
		l.problem(path, "%v", err)
		return nil
	}
	known := map[string]bool{}
	for _, c := range courses {
		known[c.ID] = c.Published
	}
	var out []Path
	goals := map[string]string{}
	for i, pm := range m.Paths {
		where := fmt.Sprintf("%s: paths[%d]", path, i)
		ok := true
		if !courseIDRe.MatchString(pm.ID) {
			l.problem(where, "id %q must be a lowercase slug like backend-developer", pm.ID)
			ok = false
		} else if !l.claim("path", pm.ID, where) {
			ok = false
		}
		if strings.TrimSpace(pm.Title) == "" {
			l.problem(where, "title is required")
			ok = false
		}
		if pm.Goal != "" {
			if !pathGoals[pm.Goal] {
				l.problem(where, "goal %q must be one of zero, backend, interviews, ai", pm.Goal)
				ok = false
			} else if prev, dup := goals[pm.Goal]; dup {
				l.problem(where, "goal %q already belongs to path %q", pm.Goal, prev)
				ok = false
			}
			goals[pm.Goal] = pm.ID
		}
		if len(pm.Courses) == 0 {
			l.problem(where, "courses: list at least one course")
			ok = false
		}
		p := Path{ID: pm.ID, Title: pm.Title, Description: strings.TrimSpace(pm.Description), Goal: pm.Goal}
		inPath := map[string]bool{}
		for _, st := range pm.Courses {
			published, exists := known[st.Course]
			switch {
			case !exists:
				l.problem(where, "courses: %q is not a course in this folder", st.Course)
				ok = false
			case !published:
				l.problem(where, "courses: %q is not published", st.Course)
				ok = false
			case inPath[st.Course]:
				l.problem(where, "courses: %q is listed twice", st.Course)
				ok = false
			case len(strings.TrimSpace(st.Note)) > maxTopics:
				l.problem(where, "courses: note for %q is over %d characters", st.Course, maxTopics)
				ok = false
			}
			inPath[st.Course] = true
			p.Steps = append(p.Steps, PathStep{Course: st.Course, Note: strings.TrimSpace(st.Note)})
		}
		if ok {
			out = append(out, p)
		}
	}
	return out
}
