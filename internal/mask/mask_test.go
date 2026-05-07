package mask_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/mask"
)

func TestApply_StyleFull_ReplacesAll(t *testing.T) {
	result := mask.Apply("mysecret", mask.Options{Style: mask.StyleFull, MaskChar: "*"})
	if result != "********" {
		t.Errorf("expected ********, got %q", result)
	}
}

func TestApply_StylePrefix_RevealsStart(t *testing.T) {
	result := mask.Apply("abcdefgh", mask.Options{Style: mask.StylePrefix, Reveal: 3, MaskChar: "*"})
	if result != "abc*****" {
		t.Errorf("expected abc*****, got %q", result)
	}
}

func TestApply_StyleSuffix_RevealsEnd(t *testing.T) {
	result := mask.Apply("abcdefgh", mask.Options{Style: mask.StyleSuffix, Reveal: 3, MaskChar: "*"})
	if result != "*****fgh" {
		t.Errorf("expected *****fgh, got %q", result)
	}
}

func TestApply_EmptyValue_ReturnsEmpty(t *testing.T) {
	result := mask.Apply("", mask.Options{Style: mask.StyleFull, MaskChar: "*"})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestApply_RevealExceedsLength_ReturnsOriginal(t *testing.T) {
	result := mask.Apply("hi", mask.Options{Style: mask.StylePrefix, Reveal: 10, MaskChar: "*"})
	if result != "hi" {
		t.Errorf("expected original value, got %q", result)
	}
}

func TestApply_DefaultMaskChar_UsesAsterisk(t *testing.T) {
	result := mask.Apply("secret", mask.Options{Style: mask.StyleFull})
	if !strings.ContainsRune(result, '*') {
		t.Errorf("expected asterisks in masked output, got %q", result)
	}
}

func TestApplyDefault_MasksSuffix(t *testing.T) {
	result := mask.ApplyDefault("supersecretvalue")
	// last 4 chars revealed: "alue"
	if !strings.HasSuffix(result, "alue") {
		t.Errorf("expected suffix 'alue' to be revealed, got %q", result)
	}
	if !strings.HasPrefix(result, "************") {
		t.Errorf("expected leading asterisks, got %q", result)
	}
}

func TestMaskEnv_OnlySensitiveKeysMasked(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "topsecret",
		"APP_NAME": "myapp",
	}
	isSensitive := func(k string) bool { return k == "API_KEY" }
	result := mask.MaskEnv(env, isSensitive, mask.Options{Style: mask.StyleFull, MaskChar: "*"})

	if result["API_KEY"] != "---------" && result["API_KEY"] == "topsecret" {
		t.Error("expected API_KEY to be masked")
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", result["APP_NAME"])
	}
}

func TestMaskEnv_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"SECRET": "abc123"}
	isSensitive := func(k string) bool { return true }
	_ = mask.MaskEnv(env, isSensitive, mask.Options{Style: mask.StyleFull, MaskChar: "*"})
	if env["SECRET"] != "abc123" {
		t.Error("original map was mutated")
	}
}
