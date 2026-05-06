// Package profile provides named environment profile management for envdiff.
//
// A profile is a named collection of .env file paths that can be stored,
// retrieved, and persisted to a local JSON registry file. This allows users
// to define shorthand names (e.g. "staging", "prod") for sets of env files
// they frequently compare or reconcile.
//
// Example usage:
//
//	r, err := profile.Load(".envdiff/profiles.json")
//	r.Set(profile.Profile{Name: "staging", Files: []string{".env", ".env.staging"}})
//	r.Save()
package profile
