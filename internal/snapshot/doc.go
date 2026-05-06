// Package snapshot provides functionality for saving and loading point-in-time
// snapshots of environment variable maps.
//
// Snapshots are stored as JSON files and include a label and timestamp
// alongside the captured key-value pairs. They can be used to compare
// an env file against a previously recorded baseline, enabling drift
// detection across deployments or time periods.
package snapshot
