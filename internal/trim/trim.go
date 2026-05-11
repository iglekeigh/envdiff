// Package trim provides utilities for cleaning .env maps by removing
// leading/trailing whitespace from keys and values, stripping blank entries,
// and deduplicating keys (last-write wins).
package trim

import (
	"strings"
)

// Options controls which trim operations are applied.
type Options struct {
	// TrimKeys removes surrounding whitespace from keys.
	TrimKeys bool
	// TrimValues removes surrounding whitespace from values.
	TrimValues bool
	// RemoveEmpty drops entries whose value is empty after trimming.
	RemoveEmpty bool
}

// DefaultOptions returns a sensible default: trim both keys and values,
// but keep empty-value entries.
func DefaultOptions() Options {
	return Options{
		TrimKeys:    true,
		TrimValues:  true,
		RemoveEmpty: false,
	}
}

// Result holds the cleaned env map and a summary of changes made.
type Result struct {
	Env          map[string]string
	TrimmedKeys  []string
	TrimmedValues []string
	RemovedKeys  []string
}

// Apply cleans the provided env map according to opts and returns a Result.
// The original map is never mutated.
func Apply(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	result := Result{
		Env: out,
	}

	for k, v := range env {
		newKey := k
		if opts.TrimKeys {
			newKey = strings.TrimSpace(k)
			if newKey != k {
				result.TrimmedKeys = append(result.TrimmedKeys, k)
			}
		}

		newVal := v
		if opts.TrimValues {
			newVal = strings.TrimSpace(v)
			if newVal != v {
				result.TrimmedValues = append(result.TrimmedValues, newKey)
			}
		}

		if opts.RemoveEmpty && newVal == "" {
			result.RemovedKeys = append(result.RemovedKeys, newKey)
			continue
		}

		out[newKey] = newVal
	}

	return result
}
