package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// This is a wrapper that delegates to the actual conference service implementation.
// The conference service is a standalone program in microservices/conference.
//
// This wrapper simply executes the actual main.go from microservices/conference
// providing a unified entry point: go run backend/inscore/cmd/conference/main.go

func main() {
	// Get the project root directory
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get executable path: %v\n", err)
		os.Exit(1)
	}
	
	projectRoot := filepath.Join(filepath.Dir(exePath), "..", "..", "..", "..")
	actualMain := filepath.Join(projectRoot, "backend", "inscore", "microservices", "conference", "main.go")
	
	// Check if running with go run (development mode)
	if _, err := os.Stat(actualMain); err != nil {
		// Try relative path for go run
		actualMain = filepath.Join("backend", "inscore", "microservices", "conference", "main.go")
	}
	
	fmt.Println("Starting Conference service...")
	fmt.Printf("Delegating to: %s\n", actualMain)
	
	// Execute the actual conference service
	cmd := exec.Command("go", "run", actualMain)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Conference service failed: %v\n", err)
		os.Exit(1)
	}
}
