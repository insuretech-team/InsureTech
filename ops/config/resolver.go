package config

import (
	"os"
	"path/filepath"
)

// FindProjectRoot walks up the directory tree to find the project root
// It looks for go.mod as the indicator of the project root
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if go.mod exists in current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding go.mod
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

// ResolvePath resolves a path relative to the project root
// Example: ResolvePath("inscore/configs/database.yaml")
func ResolvePath(relativePath string) (string, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, relativePath), nil
}

// ResolveConfigPath resolves a config file path
// It tries multiple locations in order:
// 1. Absolute path if provided
// 2. Relative to current directory
// 3. In inscore/configs/ (relative to current directory)
// 4. In inscore/configs/ (relative to project root via go.mod)
// 5. Relative to project root
// 6. Docker fallback paths
func ResolveConfigPath(configPath string) (string, error) {
	// If empty, use default
	if configPath == "" {
		configPath = "database.yaml"
	}

	// Try absolute path
	if filepath.IsAbs(configPath) {
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	// Try relative to current directory
	if _, err := os.Stat(configPath); err == nil {
		absPath, _ := filepath.Abs(configPath)
		return absPath, nil
	}

	// Try in inscore/configs/ relative to current directory (Docker scenario)
	fullPath := filepath.Join("inscore", "configs", configPath)
	if _, err := os.Stat(fullPath); err == nil {
		absPath, _ := filepath.Abs(fullPath)
		return absPath, nil
	}

	// Find project root via go.mod
	root, err := FindProjectRoot()
	if err == nil {
		// Try in backend/inscore/configs/ (new structure)
		fullPath = filepath.Join(root, "backend", "inscore", "configs", configPath)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}

		// Try in inscore/configs/ (legacy structure)
		fullPath = filepath.Join(root, "inscore", "configs", configPath)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}

		// Try in project root
		fullPath = filepath.Join(root, configPath)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}

	// Docker fallback: try /app/inscore/configs/ (absolute path in container)
	fullPath = filepath.Join("/app", "inscore", "configs", configPath)
	if _, err := os.Stat(fullPath); err == nil {
		return fullPath, nil
	}

	// Last resort: try ../configs/ (one level up from binary)
	fullPath = filepath.Join("..", "configs", configPath)
	if _, err := os.Stat(fullPath); err == nil {
		absPath, _ := filepath.Abs(fullPath)
		return absPath, nil
	}

	if configPath != "database.yaml" {
		return ResolveConfigPath("database.yaml")
	}

	// Not found anywhere
	return "", os.ErrNotExist
}
