package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Format represents a supported export format.
type Format string

const (
	FormatEnv  Format = "env"
	FormatJSON Format = "json"
	FormatShell Format = "shell"
)

// Result holds the outcome of an export operation.
type Result struct {
	Path    string
	Format  Format
	Count   int
	Skipped []string
}

// Export writes env to the given file path in the specified format.
func Export(env map[string]string, path string, format Format) (Result, error) {
	if format == "" {
		format = inferFormat(path)
	}

	var content string
	var skipped []string
	var err error

	switch format {
	case FormatEnv:
		content = toEnvFormat(env)
	case FormatJSON:
		content, err = toJSONFormat(env)
		if err != nil {
			return Result{}, fmt.Errorf("export: json marshal failed: %w", err)
		}
	case FormatShell:
		content, skipped = toShellFormat(env)
	default:
		return Result{}, fmt.Errorf("export: unsupported format %q", format)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return Result{}, fmt.Errorf("export: failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return Result{}, fmt.Errorf("export: failed to write file: %w", err)
	}

	return Result{
		Path:    path,
		Format:  format,
		Count:   len(env) - len(skipped),
		Skipped: skipped,
	}, nil
}

func inferFormat(path string) Format {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return FormatJSON
	case ".sh":
		return FormatShell
	default:
		return FormatEnv
	}
}

func toEnvFormat(env map[string]string) string {
	keys := sortedKeys(env)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String()
}

func toJSONFormat(env map[string]string) (string, error) {
	b, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func toShellFormat(env map[string]string) (string, []string) {
	keys := sortedKeys(env)
	var sb strings.Builder
	var skipped []string
	for _, k := range keys {
		if strings.ContainsAny(k, " -.") {
			skipped = append(skipped, k)
			continue
		}
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String(), skipped
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
