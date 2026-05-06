package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/snapshot"
)

// WriteSnapshotInfo writes metadata and env contents of a snapshot to w.
func WriteSnapshotInfo(w io.Writer, s *snapshot.Snapshot, redactSensitive bool, sensitive map[string]bool) {
	fmt.Fprintf(w, "Snapshot: %s\n", s.Label)
	fmt.Fprintf(w, "Created:  %s\n", s.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Keys:     %d\n", len(s.Env))
	fmt.Fprintln(w, "---")

	keys := make([]string, 0, len(s.Env))
	for k := range s.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := s.Env[k]
		if redactSensitive && sensitive[k] {
			v = "[REDACTED]"
		}
		fmt.Fprintf(w, "  %s=%s\n", k, v)
	}
}

// WriteSnapshotDiff writes a human-readable comparison between two snapshots.
func WriteSnapshotDiff(w io.Writer, base, other *snapshot.Snapshot) {
	fmt.Fprintf(w, "Comparing snapshots: %q → %q\n", base.Label, other.Label)
	fmt.Fprintln(w, "---")

	allKeys := map[string]struct{}{}
	for k := range base.Env {
		allKeys[k] = struct{}{}
	}
	for k := range other.Env {
		allKeys[k] = struct{}{}
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	changed := 0
	for _, k := range keys {
		bv, bOk := base.Env[k]
		ov, oOk := other.Env[k]
		switch {
		case bOk && !oOk:
			fmt.Fprintf(w, "  - %s (removed)\n", k)
			changed++
		case !bOk && oOk:
			fmt.Fprintf(w, "  + %s=%s (added)\n", k, ov)
			changed++
		case bv != ov:
			fmt.Fprintf(w, "  ~ %s: %q → %q\n", k, bv, ov)
			changed++
		}
	}

	if changed == 0 {
		fmt.Fprintln(w, "  No differences found.")
	} else {
		fmt.Fprintf(w, "---\n%d change(s) detected.\n", changed)
	}
}
