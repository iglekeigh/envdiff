// Package export provides functionality to write env maps to files
// in multiple formats including .env, JSON, and shell script.
//
// Supported formats:
//   - env:   KEY=VALUE pairs, one per line
//   - json:  JSON object with string values
//   - shell: export KEY="VALUE" statements for sourcing in shell scripts
//
// The export format can be specified explicitly or inferred from the
// file extension of the output path.
package export
