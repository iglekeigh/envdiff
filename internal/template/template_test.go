package template

import (
	"testing"
)

func TestCheck_NoIssues(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	tmpl := map[string]string{"HOST": "# hostname", "PORT": "# port number"}

	result := Check(env, tmpl)

	if result.HasIssues() {
		t.Fatalf("expected no issues, got missing=%v extra=%v", result.Missing, result.Extra)
	}
}

func TestCheck_MissingKey(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	tmpl := map[string]string{"HOST": "# hostname", "PORT": "# port"}

	result := Check(env, tmpl)

	if len(result.Missing) != 1 {
		t.Fatalf("expected 1 missing key, got %d", len(result.Missing))
	}
	if result.Missing[0].Key != "PORT" {
		t.Errorf("expected missing key PORT, got %s", result.Missing[0].Key)
	}
	if result.Missing[0].Comment != "port" {
		t.Errorf("expected comment 'port', got %q", result.Missing[0].Comment)
	}
}

func TestCheck_ExtraKey(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "DEBUG": "true"}
	tmpl := map[string]string{"HOST": "# hostname"}

	result := Check(env, tmpl)

	if len(result.Extra) != 1 {
		t.Fatalf("expected 1 extra key, got %d", len(result.Extra))
	}
	if result.Extra[0] != "DEBUG" {
		t.Errorf("expected extra key DEBUG, got %s", result.Extra[0])
	}
}

func TestCheck_EmptyEnv(t *testing.T) {
	env := map[string]string{}
	tmpl := map[string]string{"HOST": "# hostname", "PORT": "# port"}

	result := Check(env, tmpl)

	if len(result.Missing) != 2 {
		t.Errorf("expected 2 missing keys, got %d", len(result.Missing))
	}
}

func TestExtractComment_WithComment(t *testing.T) {
	got := extractComment("somevalue# my description")
	if got != "my description" {
		t.Errorf("expected 'my description', got %q", got)
	}
}

func TestExtractComment_NoComment(t *testing.T) {
	got := extractComment("plainvalue")
	if got != "" {
		t.Errorf("expected empty comment, got %q", got)
	}
}

func TestGenerateTemplate_PopulatesKeys(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "SECRET": ""}
	tmpl := GenerateTemplate(env)

	if _, ok := tmpl["HOST"]; !ok {
		t.Error("expected HOST in template")
	}
	if tmpl["SECRET"] != "# required" {
		t.Errorf("expected '# required' for empty key, got %q", tmpl["SECRET"])
	}
}
