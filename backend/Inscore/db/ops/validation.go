package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// CheckProtoFreshness verifies that all generated .pb.go files are newer than their source .proto files.
// This prevents running migrations or logic with stale generated code.
func CheckProtoFreshness() error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	protoRoot := filepath.Join(projectRoot, "proto", "insuretech")
	genRoot := filepath.Join(projectRoot, "gen", "go", "insuretech")

	appLogger.Info("🛡️ Validating Proto Freshness...")

	staleFiles := []string{}

	err = filepath.Walk(protoRoot, func(protoPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(protoPath, ".proto") {
			return nil
		}

		// Calculate relative path from proto root
		// e.g. "authn/entity/v1/user.proto"
		relPath, err := filepath.Rel(protoRoot, protoPath)
		if err != nil {
			return err
		}

		// Calculate expected generated file path
		// e.g. "gen/go/insuretech/authn/entity/v1/user.pb.go"
		pbPath := filepath.Join(genRoot, strings.Replace(relPath, ".proto", ".pb.go", 1))

		pbInfo, err := os.Stat(pbPath)
		if os.IsNotExist(err) {
			staleFiles = append(staleFiles, fmt.Sprintf("%s (missing .pb.go)", relPath))
			return nil
		}
		if err != nil {
			return err
		}

		// Safety check: if proto is strictly newer than generated file
		if info.ModTime().After(pbInfo.ModTime()) {
			staleFiles = append(staleFiles, fmt.Sprintf("%s (proto newer than gen code)", relPath))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk proto directory: %w", err)
	}

	if len(staleFiles) > 0 {
		appLogger.Error("❌ Local Environment Integrity Check Failed!")
		for _, f := range staleFiles {
			appLogger.Errorf("  - Stale: %s", f)
		}
		return fmt.Errorf("found %d stale or missing generated files. Please run 'buf generate' to refresh your code", len(staleFiles))
	}

	appLogger.Info("✅ Proto Freshness Verified")
	return nil
}
