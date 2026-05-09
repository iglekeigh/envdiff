// Package chain implements ordered chaining of multiple .env files.
//
// It allows users to layer environment configurations — for example, applying
// a base .env followed by a .env.local override — with explicit control over
// which value wins when a key appears in multiple files.
//
// Supported strategies:
//
//   - StrategyFirst: the first file to define a key wins (base takes precedence).
//   - StrategyLast: the last file to define a key wins (overrides take precedence).
//
// The Result type records not only the resolved environment but also the source
// file for each key and the full list of files that defined it, enabling
// transparent auditing of where each value came from.
package chain
