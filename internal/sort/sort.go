// Package sort provides utilities for ordering and ranking .env keys
// by various criteria such as alphabetical order, key length, or group prefix.
package sort

import (
	"sort"
	"strings"
)

// Strategy defines how keys should be sorted.
type Strategy string

const (
	StrategyAlpha  Strategy = "alpha"   // alphabetical by key name
	StrategyLength Strategy = "length"  // ascending key length
	StrategyGroup  Strategy = "group"   // group by common prefix (underscore-delimited)
)

// Options configures sorting behaviour.
type Options struct {
	Strategy  Strategy
	Descending bool
}

// DefaultOptions returns sensible defaults.
var DefaultOptions = Options{
	Strategy:  StrategyAlpha,
	Descending: false,
}

// Apply returns a new map with the keys sorted according to opts.
// Because maps are unordered, Apply returns an ordered slice of keys
// alongside the original map so callers can iterate deterministically.
func Apply(env map[string]string, opts Options) ([]string, map[string]string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch opts.Strategy {
	case StrategyLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) != len(keys[j]) {
				return len(keys[i]) < len(keys[j])
			}
			return keys[i] < keys[j]
		})
	case StrategyGroup:
		sort.Slice(keys, func(i, j int) bool {
			gi := groupPrefix(keys[i])
			gj := groupPrefix(keys[j])
			if gi != gj {
				return gi < gj
			}
			return keys[i] < keys[j]
		})
	default: // StrategyAlpha
		sort.Strings(keys)
	}

	if opts.Descending {
		for l, r := 0, len(keys)-1; l < r; l, r = l+1, r-1 {
			keys[l], keys[r] = keys[r], keys[l]
		}
	}

	// Return a copy of the map so callers cannot mutate the original.
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	return keys, out
}

// groupPrefix returns the first underscore-delimited segment of a key,
// which is used as the group identifier in StrategyGroup.
func groupPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
