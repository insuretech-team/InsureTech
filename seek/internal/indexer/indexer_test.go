package indexer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildSearchAndGrep(t *testing.T) {
	root := t.TempDir()

	writeTestFile(t, filepath.Join(root, "backend", "service.go"), "package backend\n\nfunc PaymentStatus() string {\n\treturn \"policy claim\"\n}\n")
	writeTestFile(t, filepath.Join(root, "api", "spec.yaml"), "summary: Policy lookup\n")
	writeTestFile(t, filepath.Join(root, "bin", "skip.exe"), "not really binary")

	dbPath := filepath.Join(root, ".seek", "index.sqlite")
	stats, err := Build(BuildOptions{
		Root:   root,
		DBPath: dbPath,
	})
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if stats.IndexedFiles != 2 {
		t.Fatalf("expected 2 indexed files, got %d", stats.IndexedFiles)
	}

	store, err := OpenStore(dbPath)
	if err != nil {
		t.Fatalf("OpenStore returned error: %v", err)
	}
	defer store.Close()

	searchResults, err := store.SearchTerm("policy", nil, 10)
	if err != nil {
		t.Fatalf("SearchTerm returned error: %v", err)
	}
	if len(searchResults) != 2 {
		t.Fatalf("expected 2 search results, got %d", len(searchResults))
	}

	grepResults, err := store.Grep(GrepOptions{
		Pattern: "PaymentStatus",
		Limit:   10,
	})
	if err != nil {
		t.Fatalf("Grep returned error: %v", err)
	}
	if len(grepResults) != 1 {
		t.Fatalf("expected 1 grep result, got %d", len(grepResults))
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
