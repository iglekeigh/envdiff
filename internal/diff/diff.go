// Package diff provides functionality to compare two sets of environment
// variables and report differences between them.
package diff

// Status represents the type of difference for a key.
type Status string

const (
	// StatusMissing indicates the key exists in base but not in target.
	StatusMissing Status = "missing"
	// StatusExtra indicates the key exists in target but not in base.
	StatusExtra Status = "extra"
	// StatusChanged indicates the key exists in both but values differ.
	StatusChanged Status = "changed"
	// StatusMatch indicates the key exists in both with equal values.
	StatusMatch Status = "match"
)

// Entry represents a single diff result for one environment key.
type Entry struct {
	Key       string
	Status    Status
	BaseValue string
	TargetValue string
}

// Result holds all diff entries produced by comparing two env maps.
type Result struct {
	Entries []Entry
}

// HasDifferences returns true if any entry is not a match.
func (r *Result) HasDifferences() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Compare compares base and target env maps and returns a Result.
// Keys are evaluated from the union of both maps.
func Compare(base, target map[string]string) *Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, bv := range base {
		seen[k] = true
		if tv, ok := target[k]; !ok {
			entries = append(entries, Entry{
				Key:       k,
				Status:    StatusMissing,
				BaseValue: bv,
			})
		} else if bv != tv {
			entries = append(entries, Entry{
				Key:         k,
				Status:      StatusChanged,
				BaseValue:   bv,
				TargetValue: tv,
			})
		} else {
			entries = append(entries, Entry{
				Key:         k,
				Status:      StatusMatch,
				BaseValue:   bv,
				TargetValue: tv,
			})
		}
	}

	for k, tv := range target {
		if !seen[k] {
			entries = append(entries, Entry{
				Key:         k,
				Status:      StatusExtra,
				TargetValue: tv,
			})
		}
	}

	return &Result{Entries: entries}
}
