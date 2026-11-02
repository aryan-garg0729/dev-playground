package exercise

import (
	"bytes"
	"fmt"
	"sort"

	"go.yaml.in/yaml/v3"
)

var registry = map[string]Type{}

// Register adds a type plugin. Called from each plugin's init().
func Register(t Type) {
	if _, dup := registry[t.Name()]; dup {
		panic("exercise: type registered twice: " + t.Name())
	}
	registry[t.Name()] = t
}

// Get returns the plugin for a type name.
func Get(name string) (Type, bool) { t, ok := registry[name]; return t, ok }

// Names lists registered type names, sorted.
func Names() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// DecodeStrict decodes a YAML node into v, rejecting unknown fields.
// Plugins use it so a typo in a manifest fails validation instead of being ignored.
func DecodeStrict(node yaml.Node, v any) error {
	if node.Kind == 0 {
		return fmt.Errorf("spec is missing")
	}
	raw, err := yaml.Marshal(&node)
	if err != nil {
		return fmt.Errorf("re-encode spec: %w", err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}
