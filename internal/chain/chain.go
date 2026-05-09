// Package chain provides support for chaining multiple .env files together,
// applying them in order with a defined override strategy.
package chain

import (
	"fmt"

	"github.com/user/envdiff/internal/envfile"
)

// Strategy controls how values are resolved when multiple files define the same key.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
)

// Result holds the merged environment and metadata about the chain operation.
type Result struct {
	// Env is the final resolved key-value map.
	Env map[string]string
	// Sources maps each key to the file path it was resolved from.
	Sources map[string]string
	// Overrides maps each key to all file paths that defined it (in order).
	Overrides map[string][]string
}

// Apply loads and chains the given env files in order using the provided strategy.
// Files are processed left-to-right; earlier files take precedence with StrategyFirst,
// later files take precedence with StrategyLast.
func Apply(paths []string, strategy Strategy) (*Result, error) {
	if len(paths) == 0 {
		return &Result{
			Env:       make(map[string]string),
			Sources:   make(map[string]string),
			Overrides: make(map[string][]string),
		}, nil
	}

	result := &Result{
		Env:       make(map[string]string),
		Sources:   make(map[string]string),
		Overrides: make(map[string][]string),
	}

	for _, path := range paths {
		env, err := envfile.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("chain: failed to parse %q: %w", path, err)
		}

		for k, v := range env {
			result.Overrides[k] = append(result.Overrides[k], path)

			_, exists := result.Env[k]
			switch {
			case !exists:
				result.Env[k] = v
				result.Sources[k] = path
			case strategy == StrategyLast:
				result.Env[k] = v
				result.Sources[k] = path
			}
		}
	}

	return result, nil
}
