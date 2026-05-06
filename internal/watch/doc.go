// Package watch provides file-change detection for .env files.
//
// It computes SHA-256 checksums of watched files and compares them
// across snapshots to detect modifications. This is useful for
// triggering re-diffs or re-validation when environment files change
// on disk during development or CI workflows.
//
// Basic usage:
//
//	w := watch.New(".env", ".env.production")
//	_ = w.Snapshot()          // record current state
//	// ... time passes ...
//	changed, _ := w.Changed() // returns paths that differ
package watch
