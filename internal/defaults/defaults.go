// Package defaults provides utilities for applying default values to env maps.
// Keys present in the defaults map but absent from the target are added;
// keys already present in the target are left unchanged.
package defaults

import "fmt"

// Result holds the outcome of applying defaults to an env map.
type Result struct {
	// Env is the resulting merged map.
	Env map[string]string
	// Applied contains keys that were added from the defaults map.
	Applied []string
	// Skipped contains keys that already existed in the target and were left unchanged.
	Skipped []string
}

// Apply merges defaultEnv into target, adding only keys that are missing.
// It does not mutate either input map.
// Returns an error if either map is nil.
func Apply(target, defaultEnv map[string]string) (*Result, error) {
	if target == nil {
		return nil, fmt.Errorf("defaults: target map must not be nil")
	}
	if defaultEnv == nil {
		return nil, fmt.Errorf("defaults: default map must not be nil")
	}

	out := make(map[string]string, len(target))
	for k, v := range target {
		out[k] = v
	}

	var applied, skipped []string
	for k, v := range defaultEnv {
		if _, exists := out[k]; exists {
			skipped = append(skipped, k)
		} else {
			out[k] = v
			applied = append(applied, k)
		}
	}

	sortStrings(applied)
	sortStrings(skipped)

	return &Result{
		Env:     out,
		Applied: applied,
		Skipped: skipped,
	}, nil
}

// sortStrings sorts a string slice in-place using a simple insertion sort
// to avoid importing "sort" just for this small helper.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
