// Package filter provides utilities for filtering environment variable maps
// by key prefix, suffix, or pattern matching.
package filter

import (
	"regexp"
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	// Prefix retains only keys that start with this string.
	Prefix string

	// Suffix retains only keys that end with this string.
	Suffix string

	// Pattern retains only keys matching this regular expression.
	Pattern string

	// Invert reverses the filter, excluding matched keys instead of retaining them.
	Invert bool
}

// Apply returns a new map containing only the entries from env that satisfy
// the filter criteria defined in opts. If no criteria are set, the original
// map is returned unchanged (as a copy).
func Apply(env map[string]string, opts Options) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		matched := matches(k, opts.Prefix, opts.Suffix, re)
		if opts.Invert {
			matched = !matched
		}
		if matched {
			result[k] = v
		}
	}
	return result, nil
}

// matches reports whether key satisfies all non-empty filter criteria.
func matches(key, prefix, suffix string, re *regexp.Regexp) bool {
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}
	if suffix != "" && !strings.HasSuffix(key, suffix) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	return true
}
