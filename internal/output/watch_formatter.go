package output

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/envdiff/internal/watch"
)

// WriteWatchStatus writes a human-readable summary of file watch states to w.
func WriteWatchStatus(w io.Writer, states []watch.FileState) {
	if len(states) == 0 {
		fmt.Fprintln(w, "watch: no files tracked")
		return
	}

	sorted := make([]watch.FileState, len(states))
	copy(sorted, states)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	fmt.Fprintf(w, "watch: tracking %d file(s)\n", len(sorted))
	for _, s := range sorted {
		fmt.Fprintf(w, "  %-40s  sha256:%.12s  %s\n",
			s.Path, s.Checksum, s.ModTime.Format(time.RFC3339))
	}
}

// WriteChangedFiles writes the list of changed file paths to w.
func WriteChangedFiles(w io.Writer, changed []string) {
	if len(changed) == 0 {
		fmt.Fprintln(w, "watch: no changes detected")
		return
	}

	sorted := make([]string, len(changed))
	copy(sorted, changed)
	sort.Strings(sorted)

	fmt.Fprintf(w, "watch: %d file(s) changed\n", len(sorted))
	for _, p := range sorted {
		fmt.Fprintf(w, "  changed: %s\n", p)
	}
}
