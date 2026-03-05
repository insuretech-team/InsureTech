package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
)

// Wrapper entrypoint that delegates to microservices/partner/cmd/server/main.go.
func main() {
	if err := logger.Initialize(logger.NoFileConfig()); err != nil {
		logger.Fatalf("Failed to initialize logger: %v", err)
	}

	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		logger.Errorf("Failed to resolve project root: %v", err)
		os.Exit(1)
	}
	if err := os.Chdir(projectRoot); err != nil {
		logger.Errorf("Failed to switch to project root (%s): %v", projectRoot, err)
		os.Exit(1)
	}

	if loadErr := env.Load(); loadErr != nil {
		logger.Warnf("Could not preload .env in wrapper: %v", loadErr)
	}

	actualMain := filepath.Join(projectRoot, "backend", "inscore", "microservices", "partner", "cmd", "server", "main.go")
	logger.Info("Starting Partner service...")
	logger.Infof("Delegating to: %s", actualMain)

	cmd := exec.Command("go", "run", actualMain)
	cmd.Dir = projectRoot
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		logger.Errorf("Partner service failed: %v", err)
		os.Exit(1)
	}
}
