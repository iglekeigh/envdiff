package convert_test

import (
	"testing"

	"github.com/user/envdiff/internal/convert"
)

func TestConvert_ScreamingSnake_FromSnake(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "db_port": "5432"}
	got, err := convert.Convert(env, convert.StyleScreamingSnake)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, got, "DB_HOST", "localhost")
	assertKey(t, got, "DB_PORT", "5432")
}

func TestConvert_Snake_FromScreamingSnake(t *testing.T) {
	env := map[string]string{"API_KEY": "secret", "BASE_URL": "http://x"}
	got, err := convert.Convert(env, convert.StyleSnake)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, got, "api_key", "secret")
	assertKey(t, got, "base_url", "http://x")
}

func TestConvert_Camel_FromScreamingSnake(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	got, err := convert.Convert(env, convert.StyleCamel)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, got, "dbHost", "localhost")
}

func TestConvert_Pascal_FromSnake(t *testing.T) {
	env := map[string]string{"my_var_name": "val"}
	got, err := convert.Convert(env, convert.StylePascal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, got, "MyVarName", "val")
}

func TestConvert_CamelInput_ToSnake(t *testing.T) {
	env := map[string]string{"dbHost": "localhost"}
	got, err := convert.Convert(env, convert.StyleSnake)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, got, "db_host", "localhost")
}

func TestConvert_UnknownStyle_ReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	_, err := convert.Convert(env, convert.Style("kebab"))
	if err == nil {
		t.Fatal("expected error for unknown style, got nil")
	}
}

func TestConvert_EmptyEnv_ReturnsEmpty(t *testing.T) {
	got, err := convert.Convert(map[string]string{}, convert.StyleSnake)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestConvert_ValuesUnchanged(t *testing.T) {
	env := map[string]string{"MY_KEY": "some value with spaces"}
	got, err := convert.Convert(env, convert.StyleCamel)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, got, "myKey", "some value with spaces")
}

func assertKey(t *testing.T, m map[string]string, key, want string) {
	t.Helper()
	v, ok := m[key]
	if !ok {
		t.Errorf("key %q not found in result %v", key, m)
		return
	}
	if v != want {
		t.Errorf("key %q: got %q, want %q", key, v, want)
	}
}
