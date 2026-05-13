package pin

import (
	"fmt"
	"sort"
	"strings"
)

// PinResult holds the outcome of a pin operation.
type PinResult struct {
	Pinned   map[string]string // keys locked to their current values
	Skipped  []string          // keys not found in the source env
	Released []string          // keys removed from the pinned set
}

// Options controls pin behaviour.
type Options struct {
	// Keys to pin. If empty, all keys in Env are pinned.
	Keys []string
	// Release removes the given keys from Pinned instead of adding them.
	Release bool
}

// Apply pins (or releases) keys in env according to opts.
// existing is the previously pinned map; it is not mutated.
func Apply(env map[string]string, existing map[string]string, opts Options) (PinResult, error) {
	if env == nil {
		return PinResult{}, fmt.Errorf("pin: env must not be nil")
	}

	pinned := copyMap(existing)
	result := PinResult{Pinned: pinned}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range env {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if opts.Release {
			if _, ok := pinned[k]; ok {
				delete(pinned, k)
				result.Released = append(result.Released, k)
			}
			continue
		}
		v, ok := env[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		pinned[k] = v
	}

	return result, nil
}

// Check returns keys whose current value in env differs from the pinned value.
func Check(env map[string]string, pinned map[string]string) []string {
	var drifted []string
	for k, pinnedVal := range pinned {
		current, ok := env[k]
		if !ok || current != pinnedVal {
			drifted = append(drifted, k)
		}
	}
	sort.Strings(drifted)
	return drifted
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
