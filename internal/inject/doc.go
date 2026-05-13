// Package inject provides utilities for injecting environment variables
// from a parsed env map into the current OS process environment or into
// another map.
//
// It supports:
//   - Prefix-based filtering to inject only matching keys
//   - Optional prefix stripping before injection
//   - Overwrite control to protect existing values
//
// Example usage:
//
//	result, err := inject.IntoOS(env, inject.Options{
//		Prefix:      "APP_",
//		StripPrefix: true,
//		Overwrite:   false,
//	})
package inject
