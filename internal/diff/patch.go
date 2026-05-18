package diff

import "fmt"

// PatchMode controls how a patch is applied to a base env map.
type PatchMode int

const (
	// PatchAddOnly only adds missing keys from the diff result; existing keys are untouched.
	PatchAddOnly PatchMode = iota
	// PatchAddAndUpdate adds missing keys and updates changed keys.
	PatchAddAndUpdate
	// PatchFull applies all changes: adds missing, updates changed, and removes extra keys.
	PatchFull
)

// PatchResult summarises what was applied during a patch operation.
type PatchResult struct {
	Added   []string
	Updated []string
	Removed []string
	Skipped []string
}

// HasChanges returns true when at least one key was modified.
func (p PatchResult) HasChanges() bool {
	return len(p.Added)+len(p.Updated)+len(p.Removed) > 0
}

// Summary returns a human-readable one-liner describing the patch.
func (p PatchResult) Summary() string {
	return fmt.Sprintf("patch: +%d ~%d -%d (skipped %d)",
		len(p.Added), len(p.Updated), len(p.Removed), len(p.Skipped))
}

// Patch applies the differences described by result onto base according to mode.
// base is mutated in place and a PatchResult is returned.
func Patch(base map[string]string, result Result, mode PatchMode) (PatchResult, error) {
	if base == nil {
		return PatchResult{}, fmt.Errorf("patch: base map must not be nil")
	}

	var pr PatchResult

	for _, entry := range result.Entries {
		switch entry.Status {
		case StatusMissing:
			// Key exists in other but not in base — add it.
			base[entry.Key] = entry.OtherValue
			pr.Added = append(pr.Added, entry.Key)

		case StatusChanged:
			if mode == PatchAddOnly {
				pr.Skipped = append(pr.Skipped, entry.Key)
				continue
			}
			base[entry.Key] = entry.OtherValue
			pr.Updated = append(pr.Updated, entry.Key)

		case StatusExtra:
			if mode != PatchFull {
				pr.Skipped = append(pr.Skipped, entry.Key)
				continue
			}
			delete(base, entry.Key)
			pr.Removed = append(pr.Removed, entry.Key)
		}
	}

	return pr, nil
}
