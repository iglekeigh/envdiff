package clone

import (
	"testing"
)

func TestClone_Exact_CopiesAllKeys(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2"}
	res, err := Clone(src, Options{Strategy: StrategyExact})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "1" || res.Env["BAR"] != "2" {
		t.Errorf("expected exact copy, got %v", res.Env)
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected no skipped keys")
	}
}

func TestClone_AddPrefix_PrependsPrefixToAllKeys(t *testing.T) {
	src := map[string]string{"HOST": "localhost", "PORT": "5432"}
	res, err := Clone(src, Options{Strategy: StrategyAddPrefix, Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", res.Env)
	}
	if res.Env["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %v", res.Env)
	}
	if res.Mapped["HOST"] != "DB_HOST" {
		t.Errorf("expected mapping HOST->DB_HOST")
	}
}

func TestClone_StripPrefix_RemovesMatchingPrefix(t *testing.T) {
	src := map[string]string{"DB_HOST": "localhost", "APP_PORT": "8080"}
	res, err := Clone(src, Options{Strategy: StrategyStripPrefix, Prefix: "DB_", OnlyMatching: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %v", res.Env)
	}
	if _, exists := res.Env["APP_PORT"]; exists {
		t.Errorf("expected APP_PORT to be skipped")
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "APP_PORT" {
		t.Errorf("expected APP_PORT in skipped, got %v", res.Skipped)
	}
}

func TestClone_StripPrefix_NoOnlyMatching_KeepsNonMatchingKeys(t *testing.T) {
	src := map[string]string{"DB_HOST": "localhost", "PORT": "8080"}
	res, err := Clone(src, Options{Strategy: StrategyStripPrefix, Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost")
	}
	if res.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080 (non-matching kept)")
	}
}

func TestClone_AddPrefix_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := Clone(map[string]string{"A": "1"}, Options{Strategy: StrategyAddPrefix, Prefix: ""})
	if err == nil {
		t.Fatal("expected error for empty prefix with StrategyAddPrefix")
	}
}

func TestClone_StripPrefix_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := Clone(map[string]string{"A": "1"}, Options{Strategy: StrategyStripPrefix, Prefix: ""})
	if err == nil {
		t.Fatal("expected error for empty prefix with StrategyStripPrefix")
	}
}

func TestClone_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	_, err := Clone(src, Options{Strategy: StrategyAddPrefix, Prefix: "X_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := src["X_KEY"]; ok {
		t.Error("Clone mutated the source map")
	}
}
