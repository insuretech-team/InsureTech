package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestWatchTreeIgnoresSeekDatabasePaths(t *testing.T) {
	root := t.TempDir()

	tree, err := newWatchTree(root, false, nil, nil)
	if err != nil {
		t.Fatalf("newWatchTree returned error: %v", err)
	}
	defer tree.Close()

	path := filepath.Join(root, ".seek", "index.sqlite")
	if tree.shouldHandlePath(path) {
		t.Fatalf("expected %q to be ignored", path)
	}
}

func TestWatchTreeHonorsExcludeDirsAndHiddenRules(t *testing.T) {
	root := t.TempDir()

	tree, err := newWatchTree(root, false, []string{"vendor"}, []string{"docs/**"})
	if err != nil {
		t.Fatalf("newWatchTree returned error: %v", err)
	}
	defer tree.Close()

	cases := []struct {
		path string
		want bool
	}{
		{filepath.Join(root, "vendor", "pkg", "file.go"), false},
		{filepath.Join(root, ".hidden", "file.go"), false},
		{filepath.Join(root, "docs", "readme.md"), false},
		{filepath.Join(root, "backend", "main.go"), true},
	}

	for _, tc := range cases {
		if got := tree.shouldHandlePath(tc.path); got != tc.want {
			t.Fatalf("shouldHandlePath(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestRootSignatureTracksTopLevelChangesOnly(t *testing.T) {
	root := t.TempDir()

	tree, err := newWatchTree(root, false, nil, nil)
	if err != nil {
		t.Fatalf("newWatchTree returned error: %v", err)
	}
	defer tree.Close()

	before, err := tree.RootSignature()
	if err != nil {
		t.Fatalf("RootSignature returned error: %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, "top.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	afterTopLevel, err := tree.RootSignature()
	if err != nil {
		t.Fatalf("RootSignature returned error: %v", err)
	}
	if before == afterTopLevel {
		t.Fatal("expected top-level root signature to change after new root file")
	}

	nestedDir := filepath.Join(root, "nested")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	baselineNested, err := tree.RootSignature()
	if err != nil {
		t.Fatalf("RootSignature returned error: %v", err)
	}

	if err := os.WriteFile(filepath.Join(nestedDir, "inner.txt"), []byte("data"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	afterNested, err := tree.RootSignature()
	if err != nil {
		t.Fatalf("RootSignature returned error: %v", err)
	}
	if baselineNested != afterNested {
		t.Fatal("expected nested-only change to not alter root-only signature")
	}
}

func TestDescribeWatchEventIncludesOperationAndPath(t *testing.T) {
	message := describeWatchEvent(fsnotify.Event{
		Name: filepath.Join("E:\\repo", "file.go"),
		Op:   fsnotify.Write,
	})
	if !strings.Contains(message, "updated") {
		t.Fatalf("expected updated message, got %q", message)
	}
	if !strings.Contains(message, "file.go") {
		t.Fatalf("expected path in message, got %q", message)
	}
}
