package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
)

// This is a wrapper that delegates to the actual authn service implementation.
// The authn service has complex internal package dependencies that cannot be
// imported from outside the microservices/authn package due to Go's internal package rules.
//
// This wrapper:
//   1. Loads .env
//   2. Optionally runs DB migrations (only when AUTHN_RUN_MIGRATIONS=true)
//   3. Starts the sync+backup sidecar as a background subprocess
//   4. Delegates to the actual authn microservice
//
// Usage: go run ./backend/inscore/cmd/authn/main.go
//
// To run migrations on startup, set the environment variable:
//
//	AUTHN_RUN_MIGRATIONS=true go run ./backend/inscore/cmd/authn/main.go

func main() {
	// Initialize logger for console output
	if err := logger.Initialize(logger.NoFileConfig()); err != nil {
		logger.Fatalf("Failed to initialize logger: %v", err)
	}

	// Resolve repository root via go.mod to keep env/config loading deterministic.
	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		logger.Errorf("Failed to resolve project root (go.mod): %v", err)
		os.Exit(1)
	}
	if err := os.Chdir(projectRoot); err != nil {
		logger.Errorf("Failed to switch working directory to project root (%s): %v", projectRoot, err)
		os.Exit(1)
	}

	// Preload .env in wrapper so all child processes inherit the same environment.
	if loadErr := env.Load(); loadErr != nil {
		logger.Warnf("Could not preload .env in wrapper: %v", loadErr)
	}

	dbopsMain := filepath.Join(projectRoot, "backend", "inscore", "cmd", "dbops", "main.go")

	// ── Step 1: Run DB migrations (optional, controlled by AUTHN_RUN_MIGRATIONS) ──
	// Migrations are skipped by default to avoid blocking service startup on every
	// restart. Set AUTHN_RUN_MIGRATIONS=true to run them explicitly (e.g. after a
	// schema change or on first deploy).
	if os.Getenv("AUTHN_RUN_MIGRATIONS") == "true" {
		logger.Info("AUTHN_RUN_MIGRATIONS=true — running DB migrations before starting AuthN...")
		migrateCmd := exec.Command("go", "run", dbopsMain, "migrate", "--target=both")
		migrateCmd.Dir = projectRoot
		migrateCmd.Env = os.Environ()
		migrateCmd.Stdout = os.Stdout
		migrateCmd.Stderr = os.Stderr
		if err := migrateCmd.Run(); err != nil {
			logger.Errorf("DB migration failed — cannot start AuthN: %v", err)
			os.Exit(1)
		}
		logger.Info("DB migrations complete.")
	} else {
		logger.Info("Skipping DB migrations (set AUTHN_RUN_MIGRATIONS=true to enable).")
	}

	// ── Step 2: Start sync+backup sidecar in background (non-blocking) ───────
	// The sidecar runs continuously alongside the authn service.
	// It is killed automatically when the authn process exits (same process group).
	logger.Info("Starting DB sync+backup sidecar in background...")
	sidecarCmd := exec.Command("go", "run", dbopsMain, "sidecar")
	sidecarCmd.Dir = projectRoot
	sidecarCmd.Env = os.Environ()
	sidecarCmd.Stdout = os.Stdout
	sidecarCmd.Stderr = os.Stderr
	if err := sidecarCmd.Start(); err != nil {
		// Non-fatal: log and continue — authn can run without sync/backup
		logger.Warnf("Could not start DB sidecar (non-fatal): %v", err)
	} else {
		logger.Infof("DB sidecar started (pid=%d)", sidecarCmd.Process.Pid)
		// Reap the sidecar asynchronously so it doesn't become a zombie
		go func() { _ = sidecarCmd.Wait() }()
	}

	// ── Step 3: Delegate to the actual authn microservice ────────────────────
	actualMain := filepath.Join(projectRoot, "backend", "inscore", "microservices", "authn", "cmd", "server", "main.go")
	logger.Info("Starting AuthN service...")
	logger.Infof("Delegating to: %s", actualMain)

	authnCmd := exec.Command("go", "run", actualMain)
	authnCmd.Dir = projectRoot
	authnCmd.Env = os.Environ()
	authnCmd.Stdout = os.Stdout
	authnCmd.Stderr = os.Stderr
	authnCmd.Stdin = os.Stdin

	if err := authnCmd.Run(); err != nil {
		logger.Errorf("AuthN service failed: %v", err)
		os.Exit(1)
	}
}
