package clone

import (
	"fmt"
	"sort"
)

// Strategy controls how cloned keys are transformed.
type Strategy int

const (
	// StrategyExact copies keys as-is.
	StrategyExact Strategy = iota
	// StrategyAddPrefix prepends a prefix to every key.
	StrategyAddPrefix
	// StrategyStripPrefix removes a prefix from matching keys.
	StrategyStripPrefix
)

// Options configures a Clone operation.
type Options struct {
	Strategy    Strategy
	Prefix      string
	OnlyMatching bool // when true, non-matching keys are omitted
}

// Result holds the output of a Clone operation.
type Result struct {
	Env      map[string]string
	Skipped  []string
	Mapped   map[string]string // original key -> new key
}

// Clone duplicates src into a new map, applying the given options.
func Clone(src map[string]string, opts Options) (*Result, error) {
	if opts.Strategy == StrategyAddPrefix && opts.Prefix == "" {
		return nil, fmt.Errorf("clone: StrategyAddPrefix requires a non-empty Prefix")
	}
	if opts.Strategy == StrategyStripPrefix && opts.Prefix == "" {
		return nil, fmt.Errorf("clone: StrategyStripPrefix requires a non-empty Prefix")
	}

	result := &Result{
		Env:    make(map[string]string),
		Mapped: make(map[string]string),
	}

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		newKey, ok := transformKey(k, opts)
		if !ok {
			if !opts.OnlyMatching {
				result.Env[k] = src[k]
				result.Mapped[k] = k
			} else {
				result.Skipped = append(result.Skipped, k)
			}
			continue
		}
		result.Env[newKey] = src[k]
		result.Mapped[k] = newKey
	}

	return result, nil
}

func transformKey(key string, opts Options) (string, bool) {
	switch opts.Strategy {
	case StrategyAddPrefix:
		return opts.Prefix + key, true
	case StrategyStripPrefix:
		if len(key) >= len(opts.Prefix) && key[:len(opts.Prefix)] == opts.Prefix {
			return key[len(opts.Prefix):], true
		}
		return key, false
	default:
		return key, true
	}
}
