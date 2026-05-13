// Package dedup provides utilities for detecting and removing duplicate
// keys across one or more env maps, with configurable conflict strategies.
package dedup

import "fmt"

// Strategy controls how duplicate keys are handled.
type Strategy int

const (
	// StrategyKeepFirst retains the first occurrence of a duplicate key.
	StrategyKeepFirst Strategy = iota
	// StrategyKeepLast retains the last occurrence of a duplicate key.
	StrategyKeepLast
	// StrategyError returns an error when a duplicate key is found.
	StrategyError
)

// Duplicate records a key that appeared more than once and the values seen.
type Duplicate struct {
	Key    string
	Values []string
}

// Result holds the deduplicated env map and any duplicates that were found.
type Result struct {
	Env        map[string]string
	Duplicates []Duplicate
}

// Apply deduplicates keys across the provided ordered list of env maps.
// Each map is treated as a separate source; keys appearing in multiple
// sources are resolved according to the chosen Strategy.
func Apply(sources []map[string]string, strategy Strategy) (*Result, error) {
	seen := make(map[string][]string)
	order := []string{}

	for _, src := range sources {
		for k, v := range src {
			if _, exists := seen[k]; !exists {
				order = append(order, k)
			}
			seen[k] = append(seen[k], v)
		}
	}

	result := &Result{
		Env: make(map[string]string),
	}

	for _, k := range order {
		vals := seen[k]
		if len(vals) == 1 {
			result.Env[k] = vals[0]
			continue
		}

		// Duplicate detected
		result.Duplicates = append(result.Duplicates, Duplicate{Key: k, Values: vals})

		switch strategy {
		case StrategyKeepFirst:
			result.Env[k] = vals[0]
		case StrategyKeepLast:
			result.Env[k] = vals[len(vals)-1]
		case StrategyError:
			return nil, fmt.Errorf("duplicate key %q found with values %v", k, vals)
		}
	}

	return result, nil
}
