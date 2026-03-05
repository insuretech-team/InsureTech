package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// LintMigrationFiles checks SQL migration files for forbidden patterns.
// Rule 2: No hardcoded SQL should add columns (Proto must be source of truth).
// Rule 3: No ALTER TABLE for structural changes in migrations.
func (umm *UnifiedMigrationManager) LintMigrationFiles() error {
	appLogger.Info("🔍 Linting SQL migration files for safety violation...")

	if !umm.dirExists(umm.migrationRoot) {
		return nil
	}

	// Forbidden patterns that imply structural changes which should be in Proto
	forbiddenPatterns := []struct {
		pattern *regexp.Regexp
		msg     string
	}{
		{
			pattern: regexp.MustCompile(`(?i)\bALTER\s+TABLE\s+.*\s+ADD\s+COLUMN\b`),
			msg:     "Do not ADD COLUMN via SQL. Add field to .proto instead.",
		},
		{
			pattern: regexp.MustCompile(`(?i)\bALTER\s+TABLE\s+.*\s+DROP\s+COLUMN\b`),
			msg:     "Do not DROP COLUMN via SQL. Remove field from .proto and use --prune (future).",
		},
		{
			pattern: regexp.MustCompile(`(?i)\bCREATE\s+TABLE\b`),
			msg:     "Do not CREATE TABLE via SQL. Define message in .proto instead.",
		},
	}

	violations := []string{}

	err := filepath.Walk(umm.migrationRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".sql") {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(contentBytes)

		// Check violations
		for _, rule := range forbiddenPatterns {
			if rule.pattern.MatchString(content) {
				relPath, _ := filepath.Rel(umm.migrationRoot, path)
				violations = append(violations, fmt.Sprintf("%s: %s", relPath, rule.msg))
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("linter failed to walk migrations: %w", err)
	}

	if len(violations) > 0 {
		appLogger.Error("❌ Migration Safety Check Failed!")
		for _, v := range violations {
			appLogger.Errorf("  - %s", v)
		}
		return fmt.Errorf("found %d safety violations in SQL migrations. Please move structural changes to Protobuf", len(violations))
	}

	appLogger.Info("✅ SQL Migrations Safe")
	return nil
}
