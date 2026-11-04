package projectcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// A path is a small subset of jq paths: "." (the whole value), ".a.b",
// ".items[0].name", "[2]". No filters, wildcards or quoting.
type pathSeg struct {
	key   string
	index int
	isIdx bool
}

func parsePath(p string) ([]pathSeg, error) {
	if p == "" {
		return nil, errors.New("required")
	}
	if p == "." {
		return nil, nil
	}
	var segs []pathSeg
	i := 0
	for i < len(p) {
		switch p[i] {
		case '.':
			j := i + 1
			for j < len(p) && p[j] != '.' && p[j] != '[' {
				j++
			}
			key := p[i+1 : j]
			if key == "" {
				return nil, fmt.Errorf("%q: empty key", p)
			}
			segs = append(segs, pathSeg{key: key})
			i = j
		case '[':
			j := strings.IndexByte(p[i:], ']')
			if j < 0 {
				return nil, fmt.Errorf("%q: missing ]", p)
			}
			n, err := strconv.Atoi(p[i+1 : i+j])
			if err != nil || n < 0 {
				return nil, fmt.Errorf("%q: index must be a non-negative integer", p)
			}
			segs = append(segs, pathSeg{index: n, isIdx: true})
			i += j + 1
		default:
			return nil, fmt.Errorf("%q: must start with . or [", p)
		}
	}
	return segs, nil
}

// lookup returns the value at path p in v (decoded JSON).
func lookup(v any, p string) (any, bool) {
	segs, err := parsePath(p)
	if err != nil {
		return nil, false
	}
	for _, s := range segs {
		if s.isIdx {
			a, ok := v.([]any)
			if !ok || s.index >= len(a) {
				return nil, false
			}
			v = a[s.index]
			continue
		}
		m, ok := v.(map[string]any)
		if !ok {
			return nil, false
		}
		if v, ok = m[s.key]; !ok {
			return nil, false
		}
	}
	return v, true
}

// normalize converts a value (from YAML or JSON) to its encoding/json form,
// so numbers are float64 and maps are map[string]any.
func normalize(v any) any {
	b, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return v
	}
	return out
}

func toFloat(v any) (float64, bool) {
	f, ok := v.(float64)
	return f, ok
}

// repr is a short JSON form of a value for Got/Expected.
func repr(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return short(string(b))
}

// canon is the full JSON form of a normalized value, for equality.
func canon(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}

// compare applies op to got (from the body) and want (from the manifest).
func compare(op string, got, want any) bool {
	switch op {
	case "eq":
		return canon(got) == canon(want) // both normalized; json.Marshal sorts map keys
	case "ne":
		return canon(got) != canon(want)
	case "gt", "lt":
		g, ok1 := toFloat(got)
		w, ok2 := toFloat(want)
		if !ok1 || !ok2 {
			return false
		}
		if op == "gt" {
			return g > w
		}
		return g < w
	case "contains":
		switch g := got.(type) {
		case string:
			w, ok := want.(string)
			return ok && strings.Contains(g, w)
		case []any:
			for _, x := range g {
				if canon(x) == canon(want) {
					return true
				}
			}
		}
	}
	return false
}
