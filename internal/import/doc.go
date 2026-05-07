// Package importer provides utilities for importing environment variable
// definitions from external file formats into the envdiff ecosystem.
//
// Supported formats:
//   - .env / dotenv (default)
//   - JSON (flat key-value objects; nested values are skipped)
//
// The format is inferred automatically from the file extension when not
// explicitly specified. Non-string JSON values (numbers, booleans, null)
// are coerced to strings; nested objects and arrays are skipped and
// recorded in the Result.Skipped field.
package importer
