package scope

import (
	"testing"
)

func TestScope_AddsPrefixToAllKeys(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	res, err := Scope(env, Options{Prefix: "DB"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := res.Env["DB_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", v)
	}
	if v, ok := res.Env["DB_PORT"]; !ok || v != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", v)
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(res.Applied))
	}
}

func TestScope_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := Scope(map[string]string{"KEY": "val"}, Options{})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestScope_CustomSeparator(t *testing.T) {
	env := map[string]string{"NAME": "alice"}
	res, err := Scope(env, Options{Prefix: "APP", Separator: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["APP.NAME"]; !ok {
		t.Error("expected key APP.NAME")
	}
}

func TestUnscope_RemovesPrefixFromMatchingKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	res, err := Unscope(env, Options{Prefix: "DB"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := res.Env["HOST"]; !ok || v != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", v)
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected no skipped keys, got %v", res.Skipped)
	}
}

func TestUnscope_NonMatchingKeys_AreSkipped(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "APP_NAME": "envdiff"}
	res, err := Unscope(env, Options{Prefix: "DB"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME in skipped, got %v", res.Skipped)
	}
	if _, ok := res.Env["HOST"]; !ok {
		t.Error("expected HOST in result")
	}
}

func TestUnscope_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := Unscope(map[string]string{"KEY": "val"}, Options{})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestScope_Unscope_RoundTrip(t *testing.T) {
	original := map[string]string{"HOST": "localhost", "PORT": "5432", "NAME": "mydb"}
	scoped, err := Scope(original, Options{Prefix: "DB"})
	if err != nil {
		t.Fatalf("scope error: %v", err)
	}
	restored, err := Unscope(scoped.Env, Options{Prefix: "DB"})
	if err != nil {
		t.Fatalf("unscope error: %v", err)
	}
	for k, v := range original {
		if got := restored.Env[k]; got != v {
			t.Errorf("round-trip mismatch for %q: want %q got %q", k, v, got)
		}
	}
}
