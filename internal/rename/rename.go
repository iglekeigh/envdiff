// Package rename provides utilities for renaming keys in an env map,
// either by exact match or by applying a prefix/suffix substitution.
package rename

import (
	"fmt"
	"strings"
)

// Rule describes a single rename operation.
type Rule struct {
	// From is the exact key name or prefix/suffix pattern to match.
	From string
	// To is the replacement key name or prefix/suffix.
	To string
	// Mode controls how matching is performed.
	Mode Mode
}

// Mode controls how a Rule matches keys.
type Mode int

const (
	// ModeExact renames only the key that exactly matches Rule.From.
	ModeExact Mode = iota
	// ModePrefix renames all keys whose name starts with Rule.From,
	// replacing the prefix with Rule.To.
	ModePrefix
	// ModeSuffix renames all keys whose name ends with Rule.From,
	// replacing the suffix with Rule.To.
	ModeSuffix
)

// Result holds the outcome of a rename operation.
type Result struct {
	Env      map[string]string
	Renamed  []Change
	Conflicts []string
}

// Change records a single key rename.
type Change struct {
	OldKey string
	NewKey string
}

// Apply applies the given rules to env and returns a Result.
// If a rename would overwrite an existing key, the key is added to Conflicts
// and the original entry is left unchanged.
func Apply(env map[string]string, rules []Rule) (Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var renamed []Change
	var conflicts []string

	for _, rule := range rules {
		if rule.From == "" {
			return Result{}, fmt.Errorf("rename rule has empty From field")
		}
		for oldKey, val := range env {
			newKey, ok := applyRule(rule, oldKey)
			if !ok || newKey == oldKey {
				continue
			}
			if _, exists := out[newKey]; exists {
				conflicts = append(conflicts, newKey)
				continue
			}
			out[newKey] = val
			delete(out, oldKey)
			renamed = append(renamed, Change{OldKey: oldKey, NewKey: newKey})
		}
	}

	return Result{Env: out, Renamed: renamed, Conflicts: conflicts}, nil
}

func applyRule(rule Rule, key string) (string, bool) {
	switch rule.Mode {
	case ModeExact:
		if key == rule.From {
			return rule.To, true
		}
	case ModePrefix:
		if strings.HasPrefix(key, rule.From) {
			return rule.To + strings.TrimPrefix(key, rule.From), true
		}
	case ModeSuffix:
		if strings.HasSuffix(key, rule.From) {
			return strings.TrimSuffix(key, rule.From) + rule.To, true
		}
	}
	return "", false
}
