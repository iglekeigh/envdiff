package template

import (
	"fmt"
	"sort"
	"strings"
)

// MissingKey represents a key that is required by the template but absent in the env map.
type MissingKey struct {
	Key     string
	Comment string
}

// Result holds the outcome of a template check.
type Result struct {
	Missing  []MissingKey
	Extra    []string
	Required []string
}

// HasIssues returns true if there are missing or extra keys.
func (r Result) HasIssues() bool {
	return len(r.Missing) > 0 || len(r.Extra) > 0
}

// Check compares an env map against a template map.
// Template values may contain a comment after a '#' describing the key's purpose.
// Keys present in env but not in the template are considered extra.
func Check(env map[string]string, tmpl map[string]string) Result {
	required := make([]string, 0, len(tmpl))
	for k := range tmpl {
		required = append(required, k)
	}
	sort.Strings(required)

	var missing []MissingKey
	for _, k := range required {
		if _, ok := env[k]; !ok {
			comment := extractComment(tmpl[k])
			missing = append(missing, MissingKey{Key: k, Comment: comment})
		}
	}

	var extra []string
	for k := range env {
		if _, ok := tmpl[k]; !ok {
			extra = append(extra, k)
		}
	}
	sort.Strings(extra)

	return Result{
		Missing:  missing,
		Extra:    extra,
		Required: required,
	}
}

// extractComment parses an optional inline comment from a template value.
// Format: "<default># <comment>" or just "# <comment>" for required fields.
func extractComment(val string) string {
	if idx := strings.Index(val, "#"); idx >= 0 {
		return strings.TrimSpace(val[idx+1:])
	}
	return ""
}

// GenerateTemplate builds a template map from an existing env map,
// annotating each key with a placeholder comment.
func GenerateTemplate(env map[string]string) map[string]string {
	tmpl := make(map[string]string, len(env))
	for k, v := range env {
		if v == "" {
			tmpl[k] = fmt.Sprintf("# required")
		} else {
			tmpl[k] = v
		}
	}
	return tmpl
}
