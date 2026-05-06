// Package encrypt provides AES-256-GCM encryption and decryption for .env file
// values. It is intended to allow envdiff to store and compare encrypted
// snapshots of sensitive environment files without exposing plaintext secrets.
//
// Keys are derived from a user-supplied passphrase using SHA-256. Each call to
// EncryptValue generates a fresh random nonce, so identical plaintexts produce
// different ciphertexts.
//
// Typical usage:
//
//	enc, err := encrypt.EncryptEnv(env, passphrase)
//	// store enc safely ...
//	dec, err := encrypt.DecryptEnv(enc, passphrase)
package encrypt
