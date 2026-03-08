package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
)

func main() {
	var (
		configPath string
		target     string
		sqlText    string
		sqlFile    string
	)

	flag.StringVar(&configPath, "config", "database.yaml", "Path to database config file")
	flag.StringVar(&target, "target", "primary", "Target database: primary|backup|both")
	flag.StringVar(&sqlText, "sql", "", "SQL to execute")
	flag.StringVar(&sqlFile, "sql-file", "", "Path to SQL file")
	flag.Parse()

	sqlText = strings.TrimSpace(sqlText)
	sqlFile = strings.TrimSpace(sqlFile)
	target = strings.ToLower(strings.TrimSpace(target))

	if sqlText == "" && sqlFile == "" {
		fatalf("either --sql or --sql-file is required")
	}
	if sqlText != "" && sqlFile != "" {
		fatalf("use either --sql or --sql-file, not both")
	}
	if sqlFile != "" {
		b, err := os.ReadFile(sqlFile)
		if err != nil {
			fatalf("failed to read sql file: %v", err)
		}
		sqlText = string(b)
	}

	// Keep output mostly clean for scripts.
	cfg := appLogger.NoFileConfig()
	cfg.Level = "error"
	cfg.Format = "text"
	cfg.Output = "console"
	cfg.Verbose = false
	_ = appLogger.Initialize(cfg)

	// Load env from repo root if present (.env / .env.local).
	_ = env.Load()

	resolvedPath, err := config.ResolveConfigPath(configPath)
	if err != nil {
		fatalf("failed to resolve config path '%s': %v", configPath, err)
	}

	dbConfig, err := db.LoadDatabaseConfig(resolvedPath)
	if err != nil {
		fatalf("failed to load database config: %v", err)
	}

	manager := db.NewDatabaseManager(dbConfig)
	var connectTargets []db.DatabaseType
	switch target {
	case "primary":
		connectTargets = []db.DatabaseType{db.Primary}
	case "backup":
		connectTargets = []db.DatabaseType{db.Backup}
	case "both":
		connectTargets = []db.DatabaseType{db.Primary, db.Backup}
	default:
		fatalf("invalid --target '%s' (expected primary|backup|both)", target)
	}

	if err := manager.ConnectWithoutMigrationsTargets(connectTargets...); err != nil {
		fatalf("failed to connect to database(s): %v", err)
	}
	defer func() { _ = manager.Close() }()

	db.Manager = manager

	if target == "both" {
		runAndPrint("primary", sqlText)
		runAndPrint("backup", sqlText)
		return
	}

	runAndPrint(target, sqlText)
}

func runAndPrint(target, sqlText string) {
	result, err := db.ExecuteSQL(sqlText, target)
	if err != nil {
		fatalf("sql execution failed for %s: %v", target, err)
	}

	fmt.Printf("=== %s ===\n", strings.ToUpper(target))
	fmt.Println(result)
}

func fatalf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
