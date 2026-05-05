// Package reconcile provides tools to merge two .env file maps into a
// single reconciled output. It supports multiple conflict resolution
// strategies:
//
//   - PreferBase: on key conflicts, the base value is retained.
//   - PreferOther: on key conflicts, the incoming (other) value wins.
//   - ErrorOnConflict: returns an error if conflicting values are found.
//
// Keys present only in other are always added to the merged result.
// Keys present only in base are always retained.
//
// Reconcile is typically used after Compare (internal/diff) to produce
// a merged .env file suitable for a target environment.
package reconcile
