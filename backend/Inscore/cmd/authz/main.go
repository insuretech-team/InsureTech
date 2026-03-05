package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
)

// This is a wrapper that delegates to the actual authz service implementation.
// The authz service has complex internal package dependencies that cannot be
// imported from outside the microservices/authz package due to Go's internal package rules.
//
// This wrapper delegates to the actual authz microservice implementation.
//
// This wrapper:
//   1. Loads .env
//   2. Delegates to the actual authz microservice
//
// Usage: go run ./backend/inscore/cmd/authz/main.go

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

	actualMain := filepath.Join(projectRoot, "backend", "inscore", "microservices", "authz", "cmd", "server", "main.go")

	logger.Info("Starting AuthZ service...")
	logger.Infof("Delegating to: %s", actualMain)

	cmd := exec.Command("go", "run", actualMain)
	cmd.Dir = projectRoot
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		logger.Errorf("AuthZ service failed: %v", err)
		os.Exit(1)
	}
}
