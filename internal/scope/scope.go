// Package scope provides utilities for scoping env maps to a named
// namespace prefix, allowing isolation of keys per environment or service.
package scope

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the output of a scope operation.
type Result struct {
	// Env is the resulting scoped or unscoped env map.
	Env map[string]string
	// Applied is the list of keys that were transformed.
	Applied []string
	// Skipped is the list of keys that did not match during unscope.
	Skipped []string
}

// Options controls scoping behaviour.
type Options struct {
	// Prefix is the namespace prefix, e.g. "APP" produces "APP_KEY".
	Prefix string
	// Separator between prefix and key; defaults to "_".
	Separator string
}

func (o *Options) sep() string {
	if o.Separator == "" {
		return "_"
	}
	return o.Separator
}

// Scope adds a prefix to every key in env.
func Scope(env map[string]string, opts Options) (Result, error) {
	if opts.Prefix == "" {
		return Result{}, fmt.Errorf("scope: prefix must not be empty")
	}
	sep := opts.sep()
	out := make(map[string]string, len(env))
	applied := make([]string, 0, len(env))

	keys := sortedKeys(env)
	for _, k := range keys {
		newKey := opts.Prefix + sep + k
		out[newKey] = env[k]
		applied = append(applied, newKey)
	}
	return Result{Env: out, Applied: applied}, nil
}

// Unscope removes a prefix from keys that carry it; keys without the prefix
// are collected in Skipped.
func Unscope(env map[string]string, opts Options) (Result, error) {
	if opts.Prefix == "" {
		return Result{}, fmt.Errorf("scope: prefix must not be empty")
	}
	sep := opts.sep()
	prefixWithSep := opts.Prefix + sep
	out := make(map[string]string, len(env))
	applied := []string{}
	skipped := []string{}

	keys := sortedKeys(env)
	for _, k := range keys {
		if strings.HasPrefix(k, prefixWithSep) {
			newKey := strings.TrimPrefix(k, prefixWithSep)
			out[newKey] = env[k]
			applied = append(applied, k)
		} else {
			skipped = append(skipped, k)
		}
	}
	return Result{Env: out, Applied: applied, Skipped: skipped}, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
