// Package merge provides functionality to merge multiple .env files
// into a single unified environment map, with configurable conflict resolution.
package merge

import "fmt"

// Strategy defines how key conflicts are handled during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error if the same key appears in multiple files.
	StrategyError
)

// Result holds the merged environment map and metadata about the merge.
type Result struct {
	Env       map[string]string
	Conflicts map[string][]string // key -> list of files that defined it
}

// Merge combines multiple env maps into one according to the given strategy.
// The sources slice contains env maps in order; labels provides optional
// human-readable names for each source (used in conflict reporting).
func Merge(sources []map[string]string, labels []string, strategy Strategy) (*Result, error) {
	result := &Result{
		Env:       make(map[string]string),
		Conflicts: make(map[string][]string),
	}

	for i, src := range sources {
		label := sourceLabel(labels, i)
		for k, v := range src {
			existing, exists := result.Env[k]
			if !exists {
				result.Env[k] = v
				result.Conflicts[k] = []string{label}
				continue
			}
			result.Conflicts[k] = append(result.Conflicts[k], label)
			switch strategy {
			case StrategyFirst:
				// keep existing, do nothing
				_ = existing
			case StrategyLast:
				result.Env[k] = v
			case StrategyError:
				return nil, fmt.Errorf("merge conflict: key %q defined in multiple sources", k)
			}
		}
	}

	// Remove non-conflicted keys from Conflicts map
	for k, srcs := range result.Conflicts {
		if len(srcs) <= 1 {
			delete(result.Conflicts, k)
		}
	}

	return result, nil
}

func sourceLabel(labels []string, i int) string {
	if i < len(labels) && labels[i] != "" {
		return labels[i]
	}
	return fmt.Sprintf("source%d", i+1)
}
