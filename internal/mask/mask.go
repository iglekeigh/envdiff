// Package mask provides utilities for partially masking sensitive values
// in .env files, preserving a configurable number of characters for identification.
package mask

import (
	"strings"
)

// Style controls how a value is masked.
type Style int

const (
	// StyleFull replaces the entire value with asterisks.
	StyleFull Style = iota
	// StylePrefix reveals the first N characters.
	StylePrefix
	// StyleSuffix reveals the last N characters.
	StyleSuffix
)

// Options configures masking behaviour.
type Options struct {
	Style   Style
	Reveal  int    // number of characters to reveal (for Prefix/Suffix styles)
	MaskChar string // character to use for masking; defaults to "*"
}

var defaultOptions = Options{
	Style:    StyleSuffix,
	Reveal:   4,
	MaskChar: "*",
}

// Apply masks a single value according to opts.
func Apply(value string, opts Options) string {
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}
	if len(value) == 0 {
		return value
	}
	reveal := opts.Reveal
	if reveal < 0 {
		reveal = 0
	}
	switch opts.Style {
	case StyleFull:
		return strings.Repeat(opts.MaskChar, len(value))
	case StylePrefix:
		if reveal >= len(value) {
			return value
		}
		return value[:reveal] + strings.Repeat(opts.MaskChar, len(value)-reveal)
	case StyleSuffix:
		if reveal >= len(value) {
			return value
		}
		return strings.Repeat(opts.MaskChar, len(value)-reveal) + value[len(value)-reveal:]
	default:
		return strings.Repeat(opts.MaskChar, len(value))
	}
}

// ApplyDefault masks a value using default options (suffix, reveal 4 chars).
func ApplyDefault(value string) string {
	return Apply(value, defaultOptions)
}

// MaskEnv returns a new map with sensitive keys masked according to opts.
// The isSensitive function determines which keys should be masked.
func MaskEnv(env map[string]string, isSensitive func(string) bool, opts Options) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k) {
			result[k] = Apply(v, opts)
		} else {
			result[k] = v
		}
	}
	return result
}
