// Package mask provides partial and full masking of sensitive environment
// variable values. Unlike full redaction (which replaces values with a
// placeholder token), masking preserves a configurable number of characters
// to assist with debugging and identification while still protecting secrets.
//
// Three masking styles are supported:
//
//   - StyleFull: replaces the entire value with mask characters.
//   - StylePrefix: reveals the first N characters, masks the rest.
//   - StyleSuffix: reveals the last N characters, masks the rest (default).
//
// Use MaskEnv to apply masking across an entire env map, driven by a
// caller-supplied sensitivity predicate.
package mask
