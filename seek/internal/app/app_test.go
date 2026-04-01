package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"seek/internal/config"
	"seek/internal/indexer"
)

func TestResolveRootUsesSavedDefaultRoot(t *testing.T) {
	appData := t.TempDir()
	t.Setenv("APPDATA", appData)
	t.Setenv("XDG_CONFIG_HOME", appData)

	root := t.TempDir()
	other := t.TempDir()

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousWD)
	})
	if err := os.Chdir(other); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	if err := saveDefaultRoot(root); err != nil {
		t.Fatalf("saveDefaultRoot returned error: %v", err)
	}

	resolved, err := resolveRoot("")
	if err != nil {
		t.Fatalf("resolveRoot returned error: %v", err)
	}
	if resolved != root {
		t.Fatalf("expected saved default root %q, got %q", root, resolved)
	}
}

func TestResolveRootNormalizesRelativeSavedDefaultRoot(t *testing.T) {
	appData := t.TempDir()
	t.Setenv("APPDATA", appData)
	t.Setenv("XDG_CONFIG_HOME", appData)

	base := t.TempDir()
	root := filepath.Join(base, "project")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousWD)
	})
	if err := os.Chdir(base); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	if err := config.Save(config.Config{DefaultRoot: "project"}); err != nil {
		t.Fatalf("config.Save returned error: %v", err)
	}

	resolved, err := resolveRoot("")
	if err != nil {
		t.Fatalf("resolveRoot returned error: %v", err)
	}
	if resolved != root {
		t.Fatalf("expected normalized root %q, got %q", root, resolved)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load returned error: %v", err)
	}
	if cfg.DefaultRoot != root {
		t.Fatalf("expected config default root to be rewritten to %q, got %q", root, cfg.DefaultRoot)
	}
}

func TestOpenStoreForReadRejectsEmptySQLiteFile(t *testing.T) {
	root := t.TempDir()
	dbPath := filepath.Join(root, ".seek", "index.sqlite")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(dbPath, nil, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	_, _, err := openStoreForRead(root, "")
	if err == nil {
		t.Fatal("expected openStoreForRead to reject an unbuilt database")
	}
	if !strings.Contains(err.Error(), "no index found at") {
		t.Fatalf("expected helpful no-index error, got %v", err)
	}
}

func TestRunSearchUsesSavedDefaultRootWhenRootOmitted(t *testing.T) {
	appData := t.TempDir()
	t.Setenv("APPDATA", appData)
	t.Setenv("XDG_CONFIG_HOME", appData)

	root := t.TempDir()
	sourceFile := filepath.Join(root, "main.go")
	content := "package main\n\nconst sessionToken = \"session_token\"\n"
	if err := os.WriteFile(sourceFile, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if err := saveDefaultRoot(root); err != nil {
		t.Fatalf("saveDefaultRoot returned error: %v", err)
	}

	_, err := indexer.Build(indexer.BuildOptions{Root: root})
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	other := t.TempDir()
	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousWD)
	})
	if err := os.Chdir(other); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	var stdout bytes.Buffer
	if err := runSearch([]string{"--term", "session_token"}, &stdout); err != nil {
		t.Fatalf("runSearch returned error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, root) {
		t.Fatalf("expected output to reference saved default root %q, got %q", root, output)
	}
	if !strings.Contains(output, "session_token") {
		t.Fatalf("expected output to contain the searched term, got %q", output)
	}
}

func TestHandleInteractiveCommandTypeAcceptsCSVWithSpaces(t *testing.T) {
	root := ""
	mode := "search"
	typeArg := ""
	limit := 20
	contextLines := 0
	var stdout bytes.Buffer

	done, err := handleInteractiveCommand(":type go, ts, yaml", &root, &mode, &typeArg, &limit, &contextLines, &stdout)
	if err != nil {
		t.Fatalf("handleInteractiveCommand returned error: %v", err)
	}
	if done {
		t.Fatal("expected :type command to continue interactive mode")
	}
	if typeArg != "go,ts,yaml" {
		t.Fatalf("expected normalized type csv, got %q", typeArg)
	}
}
