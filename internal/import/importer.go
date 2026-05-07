package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envdiff/internal/envfile"
)

// Format represents a supported import file format.
type Format string

const (
	FormatEnv    Format = "env"
	FormatJSON   Format = "json"
	FormatDotenv Format = "dotenv"
)

// Result holds the outcome of an import operation.
type Result struct {
	Source  string
	Format  Format
	Env     map[string]string
	Skipped []string
}

// Import reads a file and returns a parsed env map.
// The format is inferred from the file extension if not specified.
func Import(path string, format Format) (*Result, error) {
	if format == "" {
		format = inferFormat(path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("import: reading file %q: %w", path, err)
	}

	var env map[string]string
	var skipped []string

	switch format {
	case FormatEnv, FormatDotenv:
		env, err = envfile.Parse(strings.NewReader(string(data)))
		if err != nil {
			return nil, fmt.Errorf("import: parsing env file: %w", err)
		}
	case FormatJSON:
		env, skipped, err = parseJSON(data)
		if err != nil {
			return nil, fmt.Errorf("import: parsing json file: %w", err)
		}
	default:
		return nil, fmt.Errorf("import: unsupported format %q", format)
	}

	return &Result{
		Source:  path,
		Format:  format,
		Env:     env,
		Skipped: skipped,
	}, nil
}

// inferFormat determines the Format based on the file extension.
func inferFormat(path string) Format {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return FormatJSON
	default:
		return FormatEnv
	}
}
