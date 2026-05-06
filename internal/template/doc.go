// Package template provides utilities for comparing a .env file against
// a .env.template (or .env.example) file.
//
// A template file defines the expected keys for an environment. Values in
// the template may include an inline comment (after '#') to describe the
// purpose of each key.
//
// Use Check to identify missing or extra keys, and GenerateTemplate to
// produce a template skeleton from an existing env map.
package template
