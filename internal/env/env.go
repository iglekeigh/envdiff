// Package env provides utilities for resolving environment variable values
// from multiple sources with priority ordering.
package env

import (
	"fmt"
	"os"
	"strings"
)

// Source represents a named source of environment variables.
type Source struct {
	Label  string
	Values map[string]string
}

// ResolveResult holds the outcome of resolving a single key.
type ResolveResult struct {
	Key        string
	Value      string
	SourceLabel string
	Found      bool
}

// Resolve looks up a key across sources in priority order (first match wins).
// If no source contains the key, it falls back to the OS environment.
func Resolve(key string, sources []Source, fallbackToOS bool) ResolveResult {
	for _, src := range sources {
		if val, ok := src.Values[key]; ok {
			return ResolveResult{
				Key:         key,
				Value:       val,
				SourceLabel: src.Label,
				Found:       true,
			}
		}
	}
	if fallbackToOS {
		if val, ok := os.LookupEnv(key); ok {
			return ResolveResult{
				Key:         key,
				Value:       val,
				SourceLabel: "os",
				Found:       true,
			}
		}
	}
	return ResolveResult{Key: key, Found: false}
}

// ResolveAll resolves all unique keys found across all sources.
// Each key is resolved using priority order (first source wins).
func ResolveAll(sources []Source, fallbackToOS bool) map[string]ResolveResult {
	seen := make(map[string]struct{})
	for _, src := range sources {
		for k := range src.Values {
			seen[k] = struct{}{}
		}
	}
	results := make(map[string]ResolveResult, len(seen))
	for key := range seen {
		results[key] = Resolve(key, sources, fallbackToOS)
	}
	return results
}

// Expand replaces ${VAR} or $VAR references in a string using the provided map.
func Expand(s string, env map[string]string) (string, error) {
	var missing []string
	result := os.Expand(s, func(key string) string {
		if val, ok := env[key]; ok {
			return val
		}
		missing = append(missing, key)
		return ""
	})
	if len(missing) > 0 {
		return result, fmt.Errorf("unresolved variables: %s", strings.Join(missing, ", "))
	}
	return result, nil
}
