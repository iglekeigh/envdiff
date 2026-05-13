// Package clone provides functionality to duplicate an env map into a new map
// with optional key transformations such as adding or stripping a prefix.
//
// Supported strategies:
//   - StrategyExact: copy keys unchanged
//   - StrategyAddPrefix: prepend a prefix to every key
//   - StrategyStripPrefix: remove a prefix from matching keys
//
// When OnlyMatching is true, keys that do not match the transformation
// criteria (e.g. lack the expected prefix) are omitted from the result
// and recorded in Result.Skipped.
package clone
