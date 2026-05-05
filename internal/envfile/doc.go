// Package envfile provides utilities for parsing and representing .env files.
//
// A .env file is a plain-text file containing key=value pairs, one per line.
// Lines beginning with '#' are treated as comments and ignored. Inline comments
// following a value (separated by ' #') are captured in the Entry.Comment field.
// Values may optionally be wrapped in double quotes, which are stripped on parse.
//
// Example usage:
//
//	env, err := envfile.Parse(".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range env.Entries {
//		fmt.Printf("%s = %s\n", entry.Key, entry.Value)
//	}
package envfile
