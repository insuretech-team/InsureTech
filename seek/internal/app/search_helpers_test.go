package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSeekIgnore(t *testing.T) {
	root := t.TempDir()
	content := "# comment\n\ndocs/\n*.tmp\nnested/file.txt\n"
	if err := os.WriteFile(filepath.Join(root, ".seekignore"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	patterns, err := loadSeekIgnore(root)
	if err != nil {
		t.Fatalf("loadSeekIgnore returned error: %v", err)
	}

	if len(patterns) == 0 {
		t.Fatal("expected patterns to be loaded")
	}
}

func TestMergeIncludeExtensionsFromTypes(t *testing.T) {
	extensions, err := mergeIncludeExtensions(".go,.md", "ts,go")
	if err != nil {
		t.Fatalf("mergeIncludeExtensions returned error: %v", err)
	}

	foundGo := false
	foundTs := false
	for _, ext := range extensions {
		if ext == ".go" {
			foundGo = true
		}
		if ext == ".ts" {
			foundTs = true
		}
	}

	if !foundGo || !foundTs {
		t.Fatalf("expected merged extensions to include .go and .ts, got %v", extensions)
	}
}
