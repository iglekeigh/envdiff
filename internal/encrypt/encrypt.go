// Package encrypt provides symmetric encryption and decryption of env file values
// using AES-GCM with a passphrase-derived key.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const nonceSize = 12

// ErrDecryptFailed is returned when decryption fails due to a bad key or corrupted data.
var ErrDecryptFailed = errors.New("decrypt: authentication failed; wrong key or corrupted data")

// deriveKey produces a 32-byte AES key from the given passphrase using SHA-256.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// EncryptValue encrypts a plaintext string using AES-256-GCM and returns a
// base64-encoded ciphertext string.
func EncryptValue(plaintext, passphrase string) (string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("encrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: create GCM: %w", err)
	}
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt: generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptValue decrypts a base64-encoded ciphertext string produced by EncryptValue.
func DecryptValue(encoded, passphrase string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("encrypt: base64 decode: %w", err)
	}
	if len(data) < nonceSize {
		return "", fmt.Errorf("encrypt: ciphertext too short")
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("encrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: create GCM: %w", err)
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrDecryptFailed
	}
	return string(plaintext), nil
}

// EncryptEnv encrypts every value in env using the given passphrase and returns
// a new map with encrypted values. Keys are left unchanged.
func EncryptEnv(env map[string]string, passphrase string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		enc, err := EncryptValue(v, passphrase)
		if err != nil {
			return nil, fmt.Errorf("encrypt: key %q: %w", k, err)
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptEnv decrypts every value in env using the given passphrase and returns
// a new map with plaintext values.
func DecryptEnv(env map[string]string, passphrase string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		dec, err := DecryptValue(v, passphrase)
		if err != nil {
			return nil, fmt.Errorf("encrypt: key %q: %w", k, err)
		}
		out[k] = dec
	}
	return out, nil
}
