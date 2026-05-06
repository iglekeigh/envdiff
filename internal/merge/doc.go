// Package merge provides multi-file .env merging with configurable
// conflict resolution strategies.
//
// Given multiple parsed env maps, Merge combines them into a single
// unified map. When the same key appears in more than one source,
// the chosen Strategy determines the outcome:
//
//   - StrategyFirst: the value from the earliest source wins.
//   - StrategyLast:  the value from the latest source wins.
//   - StrategyError: an error is returned immediately on conflict.
//
// The Result also carries a Conflicts map that records which sources
// defined each conflicting key, enabling callers to report or audit
// merge decisions.
package merge
