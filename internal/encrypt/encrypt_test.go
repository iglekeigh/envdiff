package encrypt_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/encrypt"
)

const testPass = "super-secret-passphrase"

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "my-secret-value"
	enc, err := encrypt.EncryptValue(plaintext, testPass)
	if err != nil {
		t.Fatalf("EncryptValue: unexpected error: %v", err)
	}
	if enc == plaintext {
		t.Fatal("encrypted value should differ from plaintext")
	}
	dec, err := encrypt.DecryptValue(enc, testPass)
	if err != nil {
		t.Fatalf("DecryptValue: unexpected error: %v", err)
	}
	if dec != plaintext {
		t.Errorf("want %q, got %q", plaintext, dec)
	}
}

func TestEncryptValue_DifferentNonceEachCall(t *testing.T) {
	enc1, _ := encrypt.EncryptValue("value", testPass)
	enc2, _ := encrypt.EncryptValue("value", testPass)
	if enc1 == enc2 {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestDecryptValue_WrongPassphrase_ReturnsError(t *testing.T) {
	enc, err := encrypt.EncryptValue("secret", testPass)
	if err != nil {
		t.Fatalf("EncryptValue: %v", err)
	}
	_, err = encrypt.DecryptValue(enc, "wrong-passphrase")
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
	if err != encrypt.ErrDecryptFailed {
		t.Errorf("expected ErrDecryptFailed, got %v", err)
	}
}

func TestDecryptValue_InvalidBase64_ReturnsError(t *testing.T) {
	_, err := encrypt.DecryptValue("!!!not-base64!!!", testPass)
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecryptValue_TooShort_ReturnsError(t *testing.T) {
	// valid base64 but too short to contain nonce
	_, err := encrypt.DecryptValue("dG9vc2hvcnQ=", testPass)
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}

func TestEncryptEnv_RoundTrip(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
		"HOST":        "localhost",
	}
	enc, err := encrypt.EncryptEnv(env, testPass)
	if err != nil {
		t.Fatalf("EncryptEnv: %v", err)
	}
	for k, v := range enc {
		if v == env[k] {
			t.Errorf("key %q: encrypted value should differ from plaintext", k)
		}
	}
	dec, err := encrypt.DecryptEnv(enc, testPass)
	if err != nil {
		t.Fatalf("DecryptEnv: %v", err)
	}
	for k, want := range env {
		if got := dec[k]; got != want {
			t.Errorf("key %q: want %q, got %q", k, want, got)
		}
	}
}

func TestEncryptEnv_WrongPassphrase_ReturnsError(t *testing.T) {
	env := map[string]string{"SECRET": "value"}
	enc, _ := encrypt.EncryptEnv(env, testPass)
	_, err := encrypt.DecryptEnv(enc, "bad-pass")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
	if !strings.Contains(err.Error(), "SECRET") {
		t.Errorf("error should mention the offending key, got: %v", err)
	}
}
