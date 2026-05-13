// Package pin provides functionality to lock (pin) specific environment
// variable keys to their current values and detect drift when those values
// change in a later snapshot.
//
// Typical workflow:
//
//  1. Call Apply to capture the current values of selected keys into a
//     pinned map that can be persisted (e.g. via the snapshot package).
//  2. Later, call Check with the live env and the saved pinned map to
//     identify which keys have drifted from their pinned values.
//  3. Call Apply with Options.Release to remove keys from the pinned set.
package pin
