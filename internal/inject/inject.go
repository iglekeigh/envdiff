// Package inject provides functionality for injecting environment variables
// from a map into a target string or command environment, with optional
// override and filter support.
package inject

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Options controls how injection is performed.
type Options struct {
	// Overwrite allows injected values to overwrite existing OS env vars.
	Overwrite bool
	// Prefix filters keys to only those with the given prefix.
	Prefix string
	// StripPrefix removes the prefix from keys before injecting.
	StripPrefix bool
}

// Result holds the outcome of an injection operation.
type Result struct {
	Injected  []string
	Skipped   []string
	Overwrote []string
}

// IntoOS injects the given env map into the current process environment.
func IntoOS(env map[string]string, opts Options) (Result, error) {
	var result Result

	keys := sortedKeys(env)
	for _, k := range keys {
		v := env[k]

		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		targetKey := k
		if opts.StripPrefix && opts.Prefix != "" {
			targetKey = strings.TrimPrefix(k, opts.Prefix)
		}

		if targetKey == "" {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		existing := os.Getenv(targetKey)
		if existing != "" && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if err := os.Setenv(targetKey, v); err != nil {
			return result, fmt.Errorf("inject: failed to set %q: %w", targetKey, err)
		}

		if existing != "" {
			result.Overwrote = append(result.Overwrote, targetKey)
		} else {
			result.Injected = append(result.Injected, targetKey)
		}
	}

	return result, nil
}

// IntoMap merges the source env into the target map, respecting Options.
func IntoMap(target map[string]string, source map[string]string, opts Options) (Result, error) {
	var result Result

	keys := sortedKeys(source)
	for _, k := range keys {
		v := source[k]

		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		targetKey := k
		if opts.StripPrefix && opts.Prefix != "" {
			targetKey = strings.TrimPrefix(k, opts.Prefix)
		}

		if targetKey == "" {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		_, exists := target[targetKey]
		if exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if exists {
			result.Overwrote = append(result.Overwrote, targetKey)
		} else {
			result.Injected = append(result.Injected, targetKey)
		}
		target[targetKey] = v
	}

	return result, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
