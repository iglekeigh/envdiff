// Package convert provides utilities for converting env maps between
// different key naming conventions (e.g. camelCase, snake_case, SCREAMING_SNAKE).
package convert

import (
	"strings"
	"unicode"
)

// Style represents a key naming convention.
type Style string

const (
	StyleScreamingSnake Style = "screaming_snake" // MY_VAR
	StyleSnake          Style = "snake"           // my_var
	StyleCamel          Style = "camel"           // myVar
	StylePascal         Style = "pascal"          // MyVar
)

// Convert transforms all keys in the env map to the target Style.
// Values are left unchanged. Unrecognised styles return an error.
func Convert(env map[string]string, target Style) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		words := splitKey(k)
		var newKey string
		switch target {
		case StyleScreamingSnake:
			newKey = strings.ToUpper(strings.Join(words, "_"))
		case StyleSnake:
			newKey = strings.ToLower(strings.Join(words, "_"))
		case StyleCamel:
			newKey = toCamel(words, false)
		case StylePascal:
			newKey = toCamel(words, true)
		default:
			return nil, &UnknownStyleError{Style: string(target)}
		}
		out[newKey] = v
	}
	return out, nil
}

// UnknownStyleError is returned when an unrecognised Style is requested.
type UnknownStyleError struct {
	Style string
}

func (e *UnknownStyleError) Error() string {
	return "convert: unknown style: " + e.Style
}

// splitKey breaks a key into lowercase words by splitting on underscores,
// hyphens, and camelCase boundaries.
func splitKey(key string) []string {
	// normalise separators first
	key = strings.ReplaceAll(key, "-", "_")
	var words []string
	for _, part := range strings.Split(key, "_") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		words = append(words, splitCamel(part)...)
	}
	return words
}

// splitCamel splits a camelCase or PascalCase token into lowercase words.
func splitCamel(s string) []string {
	var words []string
	start := 0
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			words = append(words, strings.ToLower(s[start:i]))
			start = i
		}
	}
	words = append(words, strings.ToLower(s[start:]))
	return words
}

func toCamel(words []string, upperFirst bool) string {
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	for i, w := range words {
		if i == 0 && !upperFirst {
			b.WriteString(strings.ToLower(w))
		} else {
			if len(w) == 0 {
				continue
			}
			b.WriteString(strings.ToUpper(w[:1]) + strings.ToLower(w[1:]))
		}
	}
	return b.String()
}
