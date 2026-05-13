package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/pin"
)

// WritePinResult writes a human-readable summary of a pin operation.
func WritePinResult(w io.Writer, res pin.PinResult) {
	keys := sortedKeys(res.Pinned)
	if len(keys) > 0 {
		fmt.Fprintf(w, "Pinned keys (%d):\n", len(keys))
		for _, k := range keys {
			fmt.Fprintf(w, "  %s=%s\n", k, res.Pinned[k])
		}
	} else {
		fmt.Fprintln(w, "No keys currently pinned.")
	}

	if len(res.Skipped) > 0 {
		sort.Strings(res.Skipped)
		fmt.Fprintf(w, "Skipped (not found): %s\n", strings.Join(res.Skipped, ", "))
	}

	if len(res.Released) > 0 {
		sort.Strings(res.Released)
		fmt.Fprintf(w, "Released: %s\n", strings.Join(res.Released, ", "))
	}
}

// WritePinDrift writes the list of keys that have drifted from their pinned values.
func WritePinDrift(w io.Writer, drifted []string) {
	if len(drifted) == 0 {
		fmt.Fprintln(w, "OK: all pinned keys match current values.")
		return
	}
	sort.Strings(drifted)
	fmt.Fprintf(w, "DRIFT detected in %d key(s):\n", len(drifted))
	for _, k := range drifted {
		fmt.Fprintf(w, "  ! %s\n", k)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
