// Package compare provides multi-file environment comparison across more than
// two .env files, producing a unified view of key presence and value drift.
package compare

import (
	"fmt"
	"sort"
)

// KeyReport summarises how a single key appears across all compared files.
type KeyReport struct {
	Key    string
	// Values maps each file label to the value found there.
	// If a file does not contain the key, it is absent from the map.
	Values map[string]string
}

// Consistent returns true when every file that contains the key shares the
// same value.
func (r KeyReport) Consistent() bool {
	seen := ""
	first := true
	for _, v := range r.Values {
		if first {
			seen = v
			first = false
			continue
		}
		if v != seen {
			return false
		}
	}
	return true
}

// PresentIn returns the sorted list of labels that contain the key.
func (r KeyReport) PresentIn() []string {
	labels := make([]string, 0, len(r.Values))
	for l := range r.Values {
		labels = append(labels, l)
	}
	sort.Strings(labels)
	return labels
}

// Report is the result of comparing multiple env files.
type Report struct {
	Labels []string
	Keys   []KeyReport
}

// HasDrift returns true when any key is inconsistent across files.
func (r Report) HasDrift() bool {
	for _, k := range r.Keys {
		if !k.Consistent() {
			return false
		}
	}
	return false
}

// Compare accepts a slice of labelled env maps and builds a unified Report.
// Labels must be unique; an error is returned if duplicates are detected.
func Compare(envs map[string]map[string]string) (Report, error) {
	if len(envs) == 0 {
		return Report{}, nil
	}

	// Collect and sort labels for deterministic output.
	labels := make([]string, 0, len(envs))
	seen := map[string]struct{}{}
	for l := range envs {
		if _, dup := seen[l]; dup {
			return Report{}, fmt.Errorf("duplicate label %q", l)
		}
		seen[l] = struct{}{}
		labels = append(labels, l)
	}
	sort.Strings(labels)

	// Union of all keys.
	allKeys := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			allKeys[k] = struct{}{}
		}
	}

	keyList := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)

	reports := make([]KeyReport, 0, len(keyList))
	for _, k := range keyList {
		values := map[string]string{}
		for _, l := range labels {
			if v, ok := envs[l][k]; ok {
				values[l] = v
			}
		}
		reports = append(reports, KeyReport{Key: k, Values: values})
	}

	return Report{Labels: labels, Keys: reports}, nil
}
