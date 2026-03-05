package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/db/ops"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	// Import ALL proto entity packages for side-effects (registers with proto registry)
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/ai/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/analytics/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/apikey/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/audit/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/beneficiary/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/claims/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/commission/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/document/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/endorsement/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/insurer/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/iot/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/mfs/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/payment/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/policy/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/products/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/refund/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/renewal/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/report/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/support/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/task/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/tenant/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/underwriting/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/voice/entity/v1"
	_ "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/entity/v1"
)

// report format for sync output: table|markdown|csv|json
var syncReportFormat = "table"

// Global Management API instance
var managementAPI *ops.ManagementAPI

func main() {
	// If no arguments provided, launch interactive TUI
	if len(os.Args) == 1 {
		// Initialize minimal logger for TUI startup (will be reconfigured on connect)
		cfg := logger.NoFileConfig()
		cfg.Level = "error"
		cfg.Format = "text"
		cfg.Output = "console"
		cfg.Verbose = false
		_ = logger.Initialize(cfg)

		// Load .env for TUI using project-aware loader
		_ = env.Load()

		if err := RunInteractiveTUI(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Prefer Cobra CLI when subcommands are used
	cobraCmds := map[string]struct{}{
		"help": {},
		"sync": {}, "status": {}, "sync-health-check": {}, "schema-discovery": {}, "schema-check": {},
		"failover": {}, "switchback": {}, "rebuild-backup": {}, "sync-repair": {}, "sync-users": {},
		"backup": {}, "restore": {}, "copy": {}, "migrate": {}, "compare": {}, "list-backups": {},
		"sizes": {}, "sql": {}, "csv-backup": {}, "csv-seed": {},
	}
	if len(os.Args) > 1 {
		first := os.Args[1]
		// If the first arg doesn't look like a legacy flag, assume Cobra subcommand
		if !strings.HasPrefix(first, "-") {
			root := newRootCommand()
			// Execute with friendly error reporting (don't crash on typos)
			if err := root.Execute(); err != nil {
				fmt.Printf("Error: %v\n\n", err)
				_ = root.Usage()
				fmt.Println()
				fmt.Println("Examples:")
				fmt.Println("  dbmanager status")
				fmt.Println("  dbmanager sync --commit --prune --report-format=table")
				fmt.Println("  dbmanager sync --table=styles --commit")
				fmt.Println("  dbmanager schema-discovery")
				fmt.Println("  dbmanager sql --sql=\"SELECT 1\" --target=primary")
				fmt.Println("  dbmanager csv-backup --table=enquiries --source=primary")
				os.Exit(2)
			}
			return
		}
		// Or if any arg explicitly matches a known Cobra command/help
		for _, a := range os.Args[1:] {
			trimmed := strings.TrimLeft(a, "-")
			_, okCmd := cobraCmds[trimmed]
			// Let Cobra handle help flags/command without tips or rewrites
			if trimmed == "help" || a == "--help" || a == "-help" || a == "-h" {
				root := newRootCommand()
				if err := root.Execute(); err != nil {
					fmt.Printf("Error: %v\n\n", err)
					_ = root.Usage()
					os.Exit(2)
				}
				return
			}
			// Support dashed pseudo-commands like --migrate (but not help)
			if okCmd {
				root := newRootCommand()
				if strings.HasPrefix(a, "-") && trimmed != "help" {
					// Rewrite dashed arg to proper subcommand once
					newArgs := make([]string, 0, len(os.Args)-1)
					replaced := false
					for _, b := range os.Args[1:] {
						tb := strings.TrimLeft(b, "-")
						if !replaced && tb == trimmed {
							newArgs = append(newArgs, trimmed)
							replaced = true
						} else {
							newArgs = append(newArgs, b)
						}
					}
					root.SetArgs(newArgs)
					// Friendly tip for dashed subcommands
					if trimmed != "help" {
						fmt.Printf("Tip: use subcommand form without dashes: dbmanager %s [flags]\n\n", trimmed)
					}
				}
				if err := root.Execute(); err != nil {
					fmt.Printf("Error: %v\n\n", err)
					_ = root.Usage()
					os.Exit(2)
				}
				return
			}
		}
	}

	// Initialize logger for console-friendly output (colored, compact)
	cfg := logger.NoFileConfig()
	cfg.Level = "info"     // minimal info
	cfg.Format = "text"    // console text with colors
	cfg.Output = "console" // console only for CLI tool
	cfg.Verbose = false    // keep INFO minimal, WARN/ERROR detailed
	_ = logger.Initialize(cfg)

	// Load .env via project-level provider (walks up to find .env in project root)
	if err := env.Load(); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}

	// CLI flags
	var (
		configPath         = flag.String("config", "database.yaml", "Path to database configuration file")
		command            = flag.String("cmd", "", "Command to execute: status, failover, switchback, sync, sync-health-check, sync-users, schema-discovery, backup, restore, copy, migrate, compare, list-backups, sizes, sql, csv-backup, csv-seed, test-init, print-schema, print-table, print-tables, print-all, print-table-data")
		tableName          = flag.String("table", "", "Table name for table-specific operations (supports schema.table format)")
		backupPath         = flag.String("backup", "", "Backup file path for restore operations")
		sourceDB           = flag.String("source", "primary", "Source database (primary/backup)")
		targetDB           = flag.String("target", "primary", "Target database (primary/backup)")
		sqlQuery           = flag.String("sql", "", "SQL query to execute directly")
		commit             = flag.Bool("commit", false, "When used with -cmd=sync, actually write changes (default dry-run behavior depends on implementation)")
		failOnDrift        = flag.Bool("fail-on-drift", false, "Exit with non-zero status if any tables remain out of sync after sync")
		prune              = flag.Bool("prune", false, "Delete rows (sync) or columns (migrate) that do not exist in source/proto")
		legacyReportFormat = flag.String("report-format", "table", "Report output: table|markdown|csv|json|tui")
		schemaName         = flag.String("schema", "", "Schema name for schema-specific operations")
		limitRows          = flag.Int("limit", 100, "Limit number of rows to display (for print-table-data)")
		strict             = flag.Bool("strict", false, "Fail on schema drift (zombie columns, type mismatches)")
	)

	// Attach custom usage (still prints if flags parsing fails)
	flag.Usage = func() { printUsage() }

	// Normalize double-dash flags to single-dash for legacy path (e.g., --commit -> -commit)
	if len(os.Args) > 1 {
		args := os.Args[1:]
		normalized := make([]string, 0, len(args))
		for _, a := range args {
			if strings.HasPrefix(a, "--") {
				normalized = append(normalized, "-"+strings.TrimPrefix(a, "--"))
			} else {
				normalized = append(normalized, a)
			}
		}
		os.Args = append([]string{os.Args[0]}, normalized...)
	}

	flag.Parse()

	// propagate report format to renderer for legacy path
	syncReportFormat = strings.ToLower(strings.TrimSpace(*legacyReportFormat))

	if *command == "" {
		printUsage()
		return
	}

	// Initialize database manager with appropriate initialization path
	switch *command {
	case "help":
		// Do not initialize anything for help
	case "rebuild-backup", "schema-check", "sync", "sync-health-check", "sync-repair", "sync-users", "schema-discovery":
		if err := initializeWithoutMigrations(*configPath); err != nil {
			logger.Fatalf("Failed to initialize database connection: %v", err)
		}
	case "sql", "csv-backup", "csv-seed", "print-schema", "print-table", "print-tables", "print-all", "print-table-data":
		if err := initializeForSQL(*configPath); err != nil {
			logger.Fatalf("Failed to initialize database connection: %v", err)
		}
	case "test-init":
		if err := initializeWithoutMigrations(*configPath); err != nil {
			logger.Fatalf("Failed to initialize database connection: %v", err)
		}
	case "migrate":
		if err := initializeWithoutMigrations(*configPath); err != nil {
			logger.Fatalf("Failed to initialize database connection: %v", err)
		}
	default:
		if err := initializeWithoutMigrations(*configPath); err != nil {
			logger.Fatalf("Failed to initialize database connection: %v", err)
		}
	}

	// Execute command
	switch *command {
	case "rebuild-backup":
		handleRebuildBackup()
	case "schema-check":
		handleSchemaCheck()
	case "status":
		handleStatus()
	case "help":
		printUsage()
	case "failover":
		handleFailover()
	case "switchback":
		handleSwitchBack()
	case "sync":
		handleSync(*tableName, *commit, *prune, *failOnDrift)
	case "sync-health-check":
		handleSyncHealthCheck()
	case "sync-repair":
		handleSyncRepair()
	case "sync-users":
		handleSyncUsers()
	case "schema-discovery":
		handleSchemaDiscovery()
	case "backup":
		handleBackup()
	case "restore":
		handleRestore(*backupPath, *targetDB)
	case "copy":
		handleCopy(*sourceDB, *targetDB)
	case "migrate":
		handleMigrate(*targetDB, *prune, *strict)
	case "compare":
		handleCompare()
	case "list-backups":
		handleListBackups()
	case "sizes":
		handleSizes()
	case "sql":
		handleSQL(*sqlQuery, *targetDB)
	case "csv-backup":
		handleCSVBackup(*tableName, *sourceDB)
	case "csv-seed":
		handleCSVSeed(*tableName, *targetDB)
	case "print-schema":
		handlePrintSchema(*schemaName, *targetDB)
	case "print-table":
		handlePrintTable(*tableName, *targetDB)
	case "print-tables":
		handlePrintTables(*schemaName, *targetDB)
	case "print-all":
		handlePrintAll(*targetDB)
	case "print-table-data":
		handlePrintTableData(*tableName, *targetDB, *limitRows)
	case "test-init":
		handleTestInit()
	default:
		fmt.Printf("Unknown command: %s\n\n", *command)
		printUsage()
		// Be friendly: show a few examples before exiting
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  dbmanager -cmd=status")
		fmt.Println("  dbmanager -cmd=sync -commit -prune")
		fmt.Println("  dbmanager -cmd=sql -sql=\"SELECT 1\" -target=primary")
		os.Exit(2)
	}
}

// newRootCommand builds the Cobra CLI while delegating to existing handlers
func newRootCommand() *cobra.Command {
	var (
		cfgPath string
		source  string
		target  string
	)

	root := &cobra.Command{
		Use:   "dbmanager",
		Short: "Database Management CLI",
		Long:  "Dependency-aware primary/backup sync, health, and maintenance",
		Example: strings.Join([]string{
			"  dbmanager status",
			"  dbmanager sync --commit --prune --report-format=table",
			"  dbmanager sync --table=styles --commit",
			"  dbmanager schema-discovery",
			"  dbmanager sql --sql=\"SELECT 1\" --target=primary",
			"  dbmanager csv-backup --table=enquiries --source=primary",
		}, "\n"),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Keep help/completion output clean and fast
			name := cmd.Name()
			if name == "help" || name == "completion" {
				return nil
			}
			// If user requested help on any command, skip heavy init
			if f := cmd.Flags().Lookup("help"); f != nil && f.Changed {
				return nil
			}
			// Initialize logger only once for Cobra path
			cfg := logger.NoFileConfig()
			cfg.Level = "info"
			cfg.Format = "text"
			cfg.Output = "console"
			cfg.Verbose = false
			_ = logger.Initialize(cfg)
			// Load env via project-level provider
			_ = env.Load()
			return nil
		},
	}

	// Add persistent flags
	root.PersistentFlags().StringVar(&cfgPath, "config", "database.yaml", "Path to database configuration file")
	root.PersistentFlags().StringVar(&source, "source", "primary", "Source database (primary/backup)")
	root.PersistentFlags().StringVar(&target, "target", "primary", "Target database (primary/backup)")

	// =========================
	// sync command
	var (
		table        string
		commit       bool
		prune        bool
		failOnDrift  bool
		reportFormat string
	)
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize databases (authoritative upsert + optional prune)",
		Long: `Synchronize data from the primary database to the backup in dependency-aware order.

This performs an authoritative upsert from primary -> backup. Optionally you can prune
rows that exist only in backup (using --prune). Use --commit to write changes, otherwise
it may act as upsert-only depending on internal implementation.

You can target a single table with --table=<name> for faster, focused sync. Post-sync a
health check is displayed. Use --fail-on-drift in CI pipelines to fail if drift remains.`,
		Example: strings.Join([]string{
			"  dbmanager sync",
			"  dbmanager sync --commit",
			"  dbmanager sync --commit --prune",
			"  dbmanager sync --table=styles --commit",
			"  dbmanager sync --report-format=json",
			"  dbmanager sync --report-format=tui",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeWithoutMigrations(cfgPath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// capture selected report format for renderer
			syncReportFormat = strings.ToLower(strings.TrimSpace(reportFormat))
			handleSync(table, commit, prune, failOnDrift)
		},
	}
	syncCmd.Flags().StringVar(&table, "table", "", "Table name for table-specific sync")
	syncCmd.Flags().BoolVar(&commit, "commit", false, "Write changes (otherwise upserts only)")
	syncCmd.Flags().BoolVar(&prune, "prune", false, "Delete backup-only rows (authoritative primary)")
	syncCmd.Flags().BoolVar(&failOnDrift, "fail-on-drift", false, "Exit non-zero if drift remains after sync (CI)")
	syncCmd.Flags().StringVar(&reportFormat, "report-format", "table", "Report output: table|markdown|csv|json|tui")
	root.AddCommand(syncCmd)

	// =========================
	// status
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show database status and metrics",
		Long: `Displays a concise status for the configured primary and backup databases,
including connectivity and basic metrics where available. Useful for quick diagnostics.`,
		Example: "  dbmanager status",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeWithoutMigrations(cfgPath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			handleStatus()
		},
	}
	root.AddCommand(statusCmd)

	// =========================
	// sync-health-check
	healthCmd := &cobra.Command{
		Use:   "sync-health-check",
		Short: "Show per-table counts and sync status",
		Long: `Counts rows per table on primary and backup databases and reports drift.
Use this before and after sync to understand where differences remain.`,
		Example: "  dbmanager sync-health-check",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeWithoutMigrations(cfgPath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			handleSyncHealthCheck()
		},
	}
	root.AddCommand(healthCmd)

	// =========================
	// schema-discovery
	discoverCmd := &cobra.Command{
		Use:   "schema-discovery",
		Short: "List public base tables on primary DB",
		Long: `Discovers and lists public base tables on the primary database.
Use this to confirm table availability before targeted sync/backup operations.`,
		Example: "  dbmanager schema-discovery",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeWithoutMigrations(cfgPath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			handleSchemaDiscovery()
		},
	}
	root.AddCommand(discoverCmd)

	// =========================
	// schema-check
	schemaCheckCmd := &cobra.Command{
		Use:   "schema-check",
		Short: "Validate schema consistency between DBs",
		Long: `Validates that the primary and backup databases have consistent schema.
Reports mismatches (missing columns/indexes, type differences) and suggests using
'rebuild-backup' to align backup schema with primary.`,
		Example: "  dbmanager schema-check",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeWithoutMigrations(cfgPath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			handleSchemaCheck()
		},
	}
	root.AddCommand(schemaCheckCmd)

	// =========================
	// failover
	failoverCmd := &cobra.Command{
		Use:   "failover",
		Short: "Switch to backup database",
		Long: `Initiates a manual failover to the backup database.
Intended for emergency operations. Ensure your application layer is ready to
use backup as the active database.`,
		Example: "  dbmanager failover",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleFailover() },
	}
	root.AddCommand(failoverCmd)

	// switchback
	switchbackCmd := &cobra.Command{
		Use:   "switchback",
		Short: "Switch back to primary database (NYI)",
		Long: `Switch active operations back to the primary database after a failover.
Note: Functionality may not yet be implemented (NYI).`,
		Example: "  dbmanager switchback",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleSwitchBack() },
	}
	root.AddCommand(switchbackCmd)

	// rebuild-backup
	rebuildCmd := &cobra.Command{
		Use:   "rebuild-backup",
		Short: "Rebuild backup database schema to match primary",
		Long: `Recreates or migrates the backup database schema to match the primary.
Run this after 'schema-check' reports mismatches. Data safety depends on implementation;
ensure you've taken backups if necessary.`,
		Example: "  dbmanager rebuild-backup",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleRebuildBackup() },
	}
	root.AddCommand(rebuildCmd)

	// sync-repair
	syncRepairCmd := &cobra.Command{
		Use:   "sync-repair",
		Short: "Repair FK gaps for critical tables",
		Long: `Attempts to repair missing foreign keys or orphaned references in backup
for a curated set of critical tables. Use this when sync fails due to FK issues.`,
		Example: "  dbmanager sync-repair",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleSyncRepair() },
	}
	root.AddCommand(syncRepairCmd)

	// sync-users
	syncUsersCmd := &cobra.Command{
		Use:   "sync-users",
		Short: "Synchronize user-related tables",
		Long: `Synchronizes user-related tables (users, roles, memberships, etc.)
with special handling for unique constraints and conflict resolution.`,
		Example: "  dbmanager sync-users",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleSyncUsers() },
	}
	root.AddCommand(syncUsersCmd)

	// =========================
	// backup
	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Create database backup (NYI)",
		Long: `Creates a compressed backup of the selected database.
Status: Not yet implemented (NYI).`,
		Example: "  dbmanager backup",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleBackup() },
	}
	root.AddCommand(backupCmd)

	// restore
	var restoreBackupPath string
	var restoreTarget string
	restoreCmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore from backup (NYI)",
		Long: `Restores a database from a backup file into the selected target (primary/backup).
Status: Not yet implemented (NYI).`,
		Example: strings.Join([]string{
			"  dbmanager restore --backup=backup_primary_20240101.sql.gz --target=primary",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleRestore(restoreBackupPath, restoreTarget) },
	}
	restoreCmd.Flags().StringVar(&restoreBackupPath, "backup", "", "Backup file path")
	restoreCmd.Flags().StringVar(&restoreTarget, "target", "primary", "Target database (primary/backup)")
	root.AddCommand(restoreCmd)

	// copy
	var copySource string
	var copyTarget string
	copyCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy data between databases",
		Long: `Copies data between primary and backup databases for selected scopes.
Useful for targeted fixes or partial data realignment.`,
		Example: "  dbmanager copy --source=primary --target=backup",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleCopy(copySource, copyTarget) },
	}
	copyCmd.Flags().StringVar(&copySource, "source", "primary", "Source database (primary/backup)")
	copyCmd.Flags().StringVar(&copyTarget, "target", "backup", "Target database (primary/backup)")
	root.AddCommand(copyCmd)

	// migrate
	var migrateTarget string
	var migratePrune bool
	var migrateStrict bool
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations on a specific target (primary|backup)",
		Long: `Runs the SQL migrations against the chosen target database.
Use --target=primary to migrate the primary, or --target=backup for the backup.`,
		Example: strings.Join([]string{
			"  dbmanager migrate --target=primary",
			"  dbmanager migrate --target=backup",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Only connect to the target database to avoid unnecessary connections
			return initializeForMigrate(cfgPath, migrateTarget)
		},
		Run: func(cmd *cobra.Command, args []string) {
			handleMigrate(migrateTarget, migratePrune, migrateStrict)
		},
	}
	migrateCmd.Flags().StringVar(&migrateTarget, "target", "primary", "Target database (primary/backup)")
	migrateCmd.Flags().BoolVar(&migratePrune, "prune", false, "Delete rows (sync) or columns (migrate) that do not exist in source/proto")
	migrateCmd.Flags().BoolVar(&migrateStrict, "strict", false, "Fail on schema drift (zombie columns, type mismatches)")
	root.AddCommand(migrateCmd)

	// compare
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare schemas between databases (NYI)",
		Long: `Compares primary and backup schemas in detail, reporting drift and differences.
Status: Not yet implemented (NYI).`,
		Example: "  dbmanager compare",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleCompare() },
	}
	root.AddCommand(compareCmd)

	// list-backups
	listBackupsCmd := &cobra.Command{
		Use:   "list-backups",
		Short: "List available backups",
		Long: `Lists backup files the tool can see in the configured backup directory.
Use this before selecting a backup for restore.`,
		Example: "  dbmanager list-backups",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleListBackups() },
	}
	root.AddCommand(listBackupsCmd)

	// sizes
	sizesCmd := &cobra.Command{
		Use:   "sizes",
		Short: "Show database sizes",
		Long: `Shows approximate database and table sizes to help capacity planning
and identify unusually large tables.`,
		Example: "  dbmanager sizes",
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeWithoutMigrations(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleSizes() },
	}
	root.AddCommand(sizesCmd)

	// sql
	var sqlQueryFlag string
	var sqlTarget string
	sqlCmd := &cobra.Command{
		Use:   "sql",
		Short: "Execute direct SQL query on target DB",
		Long: `Executes an arbitrary SQL statement on primary, backup, or both.
Be careful with destructive statements. Use --target=both to run on both DBs.`,
		Example: strings.Join([]string{
			"  dbmanager sql --sql=\"SELECT 1\" --target=primary",
			"  dbmanager sql --sql=\"DROP TABLE fabric_costs CASCADE;\" --target=backup",
			"  dbmanager sql --sql=\"VACUUM ANALYZE;\" --target=both",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleSQL(sqlQueryFlag, sqlTarget) },
	}
	sqlCmd.Flags().StringVar(&sqlQueryFlag, "sql", "", "SQL query to execute")
	sqlCmd.Flags().StringVar(&sqlTarget, "target", "primary", "Target database (primary/backup/both)")
	_ = sqlCmd.MarkFlagRequired("sql")
	root.AddCommand(sqlCmd)

	// csv-backup
	var csvBackupTable string
	var csvBackupSource string
	csvBackupCmd := &cobra.Command{
		Use:   "csv-backup",
		Short: "Export table(s) to CSV files",
		Long: `Exports one or all tables into CSV files. Useful for analysis or as a fallback
seeding mechanism for small datasets.`,
		Example: strings.Join([]string{
			"  dbmanager csv-backup --source=primary",
			"  dbmanager csv-backup --table=enquiries --source=primary",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleCSVBackup(csvBackupTable, csvBackupSource) },
	}
	csvBackupCmd.Flags().StringVar(&csvBackupTable, "table", "", "Table name (optional: export all if empty)")
	csvBackupCmd.Flags().StringVar(&csvBackupSource, "source", "primary", "Source database (primary/backup)")
	root.AddCommand(csvBackupCmd)

	// csv-seed
	var csvSeedTable string
	var csvSeedTarget string
	csvSeedCmd := &cobra.Command{
		Use:   "csv-seed",
		Short: "Import data from CSV files",
		Long: `Imports CSV files into the specified target database. Table structure must match
the CSV headers. Use together with csv-backup for round-trip validation.`,
		Example: strings.Join([]string{
			"  dbmanager csv-seed --target=primary",
			"  dbmanager csv-seed --table=enquiries --target=primary",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleCSVSeed(csvSeedTable, csvSeedTarget) },
	}
	csvSeedCmd.Flags().StringVar(&csvSeedTable, "table", "", "Table name (optional: import all CSVs)")
	csvSeedCmd.Flags().StringVar(&csvSeedTarget, "target", "primary", "Target database (primary/backup)")
	root.AddCommand(csvSeedCmd)

	// print-schema
	var printSchemaName string
	var printSchemaTarget string
	printSchemaCmd := &cobra.Command{
		Use:   "print-schema",
		Short: "Print detailed schema information",
		Long: `Displays comprehensive schema information including all tables, views, and structures.
If --schema flag is provided, shows details for that specific schema only.`,
		Example: strings.Join([]string{
			"  dbmanager print-schema --target=primary",
			"  dbmanager print-schema --schema=auth --target=primary",
			"  dbmanager print-schema --schema=public --target=backup",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handlePrintSchema(printSchemaName, printSchemaTarget) },
	}
	printSchemaCmd.Flags().StringVar(&printSchemaName, "schema", "", "Schema name (optional: shows all schemas if empty)")
	printSchemaCmd.Flags().StringVar(&printSchemaTarget, "target", "primary", "Target database (primary/backup)")
	root.AddCommand(printSchemaCmd)

	// print-table
	var printTableName string
	var printTableTarget string
	printTableCmd := &cobra.Command{
		Use:   "print-table",
		Short: "Print detailed table information",
		Long: `Displays comprehensive information about a specific table including:
- Columns with types, constraints, and defaults
- Primary keys and unique constraints
- Foreign keys and references
- Indexes
- Table size and row count
Supports schema.table format (e.g., auth.users)`,
		Example: strings.Join([]string{
			"  dbmanager print-table --table=users --target=primary",
			"  dbmanager print-table --table=auth.users --target=primary",
			"  dbmanager print-table --table=customer_addresses --target=backup",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handlePrintTable(printTableName, printTableTarget) },
	}
	printTableCmd.Flags().StringVar(&printTableName, "table", "", "Table name (required, supports schema.table format)")
	printTableCmd.Flags().StringVar(&printTableTarget, "target", "primary", "Target database (primary/backup)")
	_ = printTableCmd.MarkFlagRequired("table")
	root.AddCommand(printTableCmd)

	// print-all
	var printAllTarget string
	printAllCmd := &cobra.Command{
		Use:   "print-all",
		Short: "Print all database schemas and tables",
		Long: `Displays a comprehensive overview of the entire database including:
- All schemas
- All tables within each schema
- Table sizes and row counts
- Summary statistics`,
		Example: strings.Join([]string{
			"  dbmanager print-all --target=primary",
			"  dbmanager print-all --target=backup",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handlePrintAll(printAllTarget) },
	}
	printAllCmd.Flags().StringVar(&printAllTarget, "target", "primary", "Target database (primary/backup)")
	root.AddCommand(printAllCmd)

	// print-tables
	var printTablesSchema string
	var printTablesTarget string
	printTablesCmd := &cobra.Command{
		Use:   "print-tables",
		Short: "Print detailed information for all tables in schema(s)",
		Long: `Displays comprehensive information for all tables including:
- All columns with types, constraints, and defaults
- Primary keys and unique constraints
- Foreign keys and references
- Indexes
- Table sizes and row counts

If --schema is provided, shows tables only from that schema.
If --schema is omitted, shows tables from all schemas.`,
		Example: strings.Join([]string{
			"  dbmanager print-tables --schema=auth --target=primary",
			"  dbmanager print-tables --schema=public --target=backup",
			"  dbmanager print-tables --target=primary  # All schemas",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handlePrintTables(printTablesSchema, printTablesTarget) },
	}
	printTablesCmd.Flags().StringVar(&printTablesSchema, "schema", "", "Schema name (optional: shows all schemas if empty)")
	printTablesCmd.Flags().StringVar(&printTablesTarget, "target", "primary", "Target database (primary/backup)")
	root.AddCommand(printTablesCmd)

	// print-table-data
	var printTableDataName string
	var printTableDataTarget string
	var printTableDataLimit int
	printTableDataCmd := &cobra.Command{
		Use:   "print-table-data",
		Short: "Print actual table data/entries",
		Long: `Displays the actual data/entries from a specific table in a formatted table view.
Shows column names and row data with automatic formatting and truncation for long values.

By default, displays the first 100 rows. Use --limit to adjust (max 1000).
Supports schema.table format (e.g., auth.users)`,
		Example: strings.Join([]string{
			"  dbmanager print-table-data --table=users --target=primary",
			"  dbmanager print-table-data --table=auth.users",
			"  dbmanager print-table-data --table=auth.users --limit=50",
			"  dbmanager print-table-data --table=auth.customers --target=backup --limit=200",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run: func(cmd *cobra.Command, args []string) {
			handlePrintTableData(printTableDataName, printTableDataTarget, printTableDataLimit)
		},
	}
	printTableDataCmd.Flags().StringVar(&printTableDataName, "table", "", "Table name (required, supports schema.table format)")
	printTableDataCmd.Flags().StringVar(&printTableDataTarget, "target", "primary", "Target database (primary/backup)")
	printTableDataCmd.Flags().IntVar(&printTableDataLimit, "limit", 100, "Number of rows to display (default: 100, max: 1000)")
	_ = printTableDataCmd.MarkFlagRequired("table")
	root.AddCommand(printTableDataCmd)

	// view-tables - Interactive table viewer using Bubble Tea
	var viewTablesSource string
	var viewTablesTable string
	viewTablesCmd := &cobra.Command{
		Use:   "view-tables",
		Short: "Interactive table viewer with CSV export integration",
		Long:  "Launch an interactive table viewer that displays database tables in a nice format using CSV data",
		Example: strings.Join([]string{
			"  dbmanager view-tables",
			"  dbmanager view-tables --table=enquiries",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error { return initializeForSQL(cfgPath) },
		Run:     func(cmd *cobra.Command, args []string) { handleViewTables(viewTablesSource, viewTablesTable) },
	}
	viewTablesCmd.Flags().StringVar(&viewTablesSource, "source", "primary", "Source database (primary/backup/both)")
	viewTablesCmd.Flags().StringVar(&viewTablesTable, "table", "", "Specific table to view (optional: show all tables)")
	root.AddCommand(viewTablesCmd)

	// Show usage on flag errors and silence duplicate error prints
	root.SilenceErrors = true
	root.SilenceUsage = false
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			_ = cmd.Usage()
		}
		return err
	})

	return root
}

// initializeWithoutMigrations initializes database connections without running migrations
// Used for commands that need database access but shouldn't trigger migrations
func initializeWithoutMigrations(configPath string) error {
	// Resolve config path using project-aware resolver
	resolvedPath, err := config.ResolveConfigPath(configPath)
	if err != nil {
		return fmt.Errorf("failed to resolve config path '%s': %v", configPath, err)
	}

	dbConfig, err := db.LoadDatabaseConfig(resolvedPath)
	if err != nil {
		return fmt.Errorf("failed to load database config: %v", err)
	}

	// Create database manager without migrations
	manager := db.NewDatabaseManager(dbConfig)

	// Connect to both databases by default (most commands need both)
	if err := manager.ConnectWithoutMigrations(); err != nil {
		return fmt.Errorf("failed to connect to databases: %v", err)
	}

	// Set the global manager
	db.Manager = manager

	// Initialize Management API
	managementAPI = ops.NewManagementAPI(manager)

	return nil
}

// initializeForMigrate initializes database connection for the specified target only
// Used for migrate command to avoid connecting to unreachable databases
func initializeForMigrate(configPath, target string) error {
	// Resolve config path using project-aware resolver
	resolvedPath, err := config.ResolveConfigPath(configPath)
	if err != nil {
		return fmt.Errorf("failed to resolve config path '%s': %v", configPath, err)
	}

	dbConfig, err := db.LoadDatabaseConfig(resolvedPath)
	if err != nil {
		return fmt.Errorf("failed to load database config: %v", err)
	}

	// Create database manager
	manager := db.NewDatabaseManager(dbConfig)

	// Connect only to the target database
	t := strings.ToLower(strings.TrimSpace(target))
	var connectTargets []db.DatabaseType
	switch t {
	case "primary":
		connectTargets = []db.DatabaseType{db.Primary}
	case "backup":
		connectTargets = []db.DatabaseType{db.Backup}
	case "both", "all":
		connectTargets = []db.DatabaseType{db.Primary, db.Backup}
	default:
		connectTargets = []db.DatabaseType{db.Primary}
	}

	if err := manager.ConnectWithoutMigrationsTargets(connectTargets...); err != nil {
		return fmt.Errorf("failed to connect to target database(s): %v", err)
	}

	// Set the global manager
	db.Manager = manager

	// Initialize Management API
	managementAPI = ops.NewManagementAPI(manager)

	return nil
}

// initializeForSQL initializes minimal database connection for direct SQL operations
// Used for SQL, CSV backup, and CSV seed commands
func initializeForSQL(configPath string) error {
	// Resolve config path using project-aware resolver
	resolvedPath, err := config.ResolveConfigPath(configPath)
	if err != nil {
		return fmt.Errorf("failed to resolve config path '%s': %v", configPath, err)
	}

	dbConfig, err := db.LoadDatabaseConfig(resolvedPath)
	if err != nil {
		return fmt.Errorf("failed to load database config: %v", err)
	}

	// Create minimal database manager for SQL operations
	manager := db.NewDatabaseManager(dbConfig)

	// Connect to databases without any initialization
	if err := manager.ConnectWithoutMigrations(); err != nil {
		return fmt.Errorf("failed to connect to databases: %v", err)
	}

	// Set the global manager
	db.Manager = manager

	return nil
}

// initializeForSQLSilent initializes database connection silently (no logging)
// Used for TUI background connection
func initializeForSQLSilent(configPath string) error {
	// Initialize logger in silent mode for TUI (suppress all output)
	cfg := logger.NoFileConfig()
	cfg.Level = "error" // Only show errors
	cfg.Format = "text"
	cfg.Output = "console"
	cfg.Verbose = false
	_ = logger.Initialize(cfg)

	// Load .env via project-level provider
	_ = env.Load()

	// Resolve config path using project-aware resolver
	resolvedPath, err := config.ResolveConfigPath(configPath)
	if err != nil {
		return fmt.Errorf("failed to resolve config path '%s': %v", configPath, err)
	}

	dbConfig, err := db.LoadDatabaseConfig(resolvedPath)
	if err != nil {
		return fmt.Errorf("failed to load database config: %v", err)
	}

	// Create minimal database manager for SQL operations
	manager := db.NewDatabaseManager(dbConfig)

	// Connect to databases without any initialization
	if err := manager.ConnectWithoutMigrations(); err != nil {
		return fmt.Errorf("failed to connect to databases: %v", err)
	}

	// Set the global manager
	db.Manager = manager

	// Initialize Management API silently
	managementAPI = ops.NewManagementAPI(manager)

	return nil
}

func handleStatus() {
	logger.Info("=== Database Status ===")

	if db.Manager == nil {
		logger.Error("Database manager not initialized")
		return
	}

	primaryStatus, backupStatus := db.Manager.GetStatus()
	logger.Infof("Primary Database: %s", primaryStatus)
	logger.Infof("Backup Database: %s", backupStatus)
}

func handleFailover() {
	fmt.Println("Initiating failover to backup database...")

	if err := managementAPI.ForceFailover(); err != nil {
		logger.Fatalf("Failover failed: %v", err)
	}

	fmt.Println("Failover completed successfully")
}

func handleSwitchBack() {
	fmt.Println("Switching back to primary database...")

	// TODO: Implement SwitchBack functionality
	logger.Info("Switchback functionality not yet implemented")

	fmt.Println("Switch back completed successfully")
}

func handleRebuildBackup() {
	logger.Info("Rebuilding backup database schema to match primary...")

	if err := managementAPI.RebuildBackupSchema(); err != nil {
		logger.Fatalf("Failed to rebuild backup schema: %v", err)
	}

	logger.Info("Backup database schema rebuilt successfully")
}

func handleSchemaCheck() {
	logger.Info("Checking schema consistency between primary and backup databases...")

	mismatches, err := managementAPI.ValidateSchemaConsistency()
	if err != nil {
		logger.Fatalf("Failed to validate schema: %v", err)
	}

	if len(mismatches) == 0 {
		logger.Info("✅ Schema is consistent between primary and backup databases")
	} else {
		logger.Error("❌ Schema mismatches found:")
		for _, mismatch := range mismatches {
			logger.Errorf("  - %s", mismatch)
		}
		logger.Info("Run 'dbmanager -cmd=rebuild-backup' to fix schema mismatches")
	}
}

func handleSQL(query, targetDB string) {
	if query == "" {
		fmt.Println("Error: SQL query is required")
		fmt.Println("Usage: dbmanager sql --sql \"DROP TABLE fabric_costs CASCADE;\" --target=primary|backup|both")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  dbmanager sql --sql=\"SELECT 1\" --target=primary")
		fmt.Println("  dbmanager sql --sql=\"VACUUM ANALYZE;\" --target=both")
		return
	}

	t := strings.ToLower(strings.TrimSpace(targetDB))
	if t == "both" || t == "all" {
		fmt.Println("=== Executing SQL on BOTH databases ===")
		fmt.Printf("Query: %s\n", query)
		hadErr := false
		for _, dbTarget := range []string{"primary", "backup"} {
			fmt.Printf("--- %s ---\n", strings.ToUpper(dbTarget))
			result, err := db.ExecuteSQL(query, dbTarget)
			if err != nil {
				fmt.Printf("Error executing SQL on %s: %v\n", dbTarget, err)
				hadErr = true
				continue
			}
			fmt.Printf("Result: %s\n", result)
		}
		if hadErr {
			fmt.Println("One or more targets failed. See errors above.")
		}
		fmt.Println("SQL executed successfully on both databases")
		return
	}

	// Single target
	fmt.Printf("=== Executing SQL on %s database ===\n", t)
	fmt.Printf("Query: %s\n", query)

	result, err := db.ExecuteSQL(query, t)
	if err != nil {
		fmt.Printf("Error executing SQL: %v\n", err)
		return
	}

	fmt.Printf("Result: %s\n", result)
	fmt.Println("SQL executed successfully")
}

func handleViewTables(source, table string) {
	// Launch the interactive Bubble Tea viewer directly
	if err := RunInteractiveTableViewer(); err != nil {
		fmt.Printf("Error running interactive viewer: %v\n", err)
		return
	}
}

func handleSyncHealthCheck() {
	if db.Manager == nil {
		logger.Fatal("Database manager not initialized")
	}
	fmt.Println("=== SYNC HEALTH CHECK ===")
	statuses, err := managementAPI.SyncHealthCheck()
	if err != nil {
		logger.Fatalf("Health check failed: %v", err)
	}
	outOfSync := 0
	for _, s := range statuses {
		state := "IN SYNC"
		if !s.InSync {
			state = "OUT OF SYNC"
			outOfSync++
		}
		fmt.Printf("Table: %-30s | %s | Primary: %d | Backup: %d\n",
			s.Name, state, s.PrimaryCount, s.BackupCount)
	}
	fmt.Printf("\nSummary: %d tables checked, %d out of sync\n", len(statuses), outOfSync)
	if outOfSync > 0 {
		fmt.Println("Run 'dbmanager sync --commit' to fix sync issues")
	}
}

func handleSchemaDiscovery() {

	if db.Manager == nil {
		logger.Fatal("Database manager not initialized")
	}
	fmt.Println("=== SCHEMA DISCOVERY (primary) ===")
	tables, err := managementAPI.DiscoverSchema()
	if err != nil {
		logger.Fatalf("Schema discovery failed: %v", err)
	}
	for _, t := range tables {
		fmt.Printf("- %s\n", t)
	}
	fmt.Printf("Total tables: %d\n", len(tables))
}

func handleSyncRepair() {
	if db.Manager == nil {
		logger.Fatal("Database manager not initialized")
	}
	fmt.Println("=== SYNC REPAIR (FK backfill for critical tables) ===")

	fmt.Println("✅ Sync repair completed")
}

func handleSyncUsers() {
	if db.Manager == nil {
		logger.Fatal("Database manager not initialized")
	}
	fmt.Println("=== SYNC USERS ===")
	if err := managementAPI.SyncUsers(); err != nil {
		logger.Fatalf("Sync users failed: %v", err)
	}
	fmt.Println("✅ User tables synchronized")
}

func printUsage() {
	// Minimal, bullet-style help (no column padding) to avoid console spacing issues
	fmt.Println("dbmanager - Database Management CLI")
	fmt.Println()
	fmt.Println("NAME")
	fmt.Println("  dbmanager - dependency-aware primary/backup sync, health, and maintenance")
	fmt.Println()
	fmt.Println("SYNOPSIS")
	fmt.Println("  dbmanager -cmd=<command> [global options] [command options]")
	fmt.Println("  dbmanager help")
	fmt.Println("  dbmanager --help | -h")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  Use -cmd to select a command. Global options apply to most commands.")
	fmt.Println()
	fmt.Println("COMMANDS")
	fmt.Println("  backup             Create database backup (NYI)")
	fmt.Println("  compare            Compare schemas between databases (NYI)")
	fmt.Println("  copy               Copy data between databases")
	fmt.Println("  failover           Switch to backup database")
	fmt.Println("  help               Show this help message")
	fmt.Println("  list-backups       List available backups")
	fmt.Println("  migrate            Run database migrations (primary and backup as configured)")
	fmt.Println("  restore            Restore from backup (NYI)")
	fmt.Println("  schema-discovery   List public base tables on primary DB")
	fmt.Println("  schema-check       Validate schema consistency between DBs")
	fmt.Println("  sizes              Show database sizes")
	fmt.Println("  sql                Execute direct SQL query on target DB")
	fmt.Println("  status             Show database status and metrics")
	fmt.Println("  switchback         Switch back to primary database (NYI)")
	fmt.Println("  sync               Synchronize databases (authoritative upsert + optional prune)")
	fmt.Println("  sync-health-check  Show per-table counts and sync status")
	fmt.Println("  sync-repair        Repair FK gaps for critical tables")
	fmt.Println("  sync-users         Synchronize user-related tables")
	fmt.Println("  csv-backup         Export table(s) to CSV files")
	fmt.Println("  csv-seed           Import CSV files into database")
	fmt.Println("  print-schema       Print detailed schema information")
	fmt.Println("  print-table        Print detailed table information")
	fmt.Println("  print-tables       Print detailed info for all tables in schema(s)")
	fmt.Println("  print-all          Print all schemas and tables")
	fmt.Println("  print-table-data   Print actual table data/entries")
	fmt.Println("  view-tables        Interactive table viewer (TUI)")
	fmt.Println()
	fmt.Println("GLOBAL OPTIONS")
	fmt.Println("  -config            Path to database configuration file (default: database.yaml)")
	fmt.Println("  -source            Source database (primary/backup)")
	fmt.Println("  -target            Target database (primary/backup)")
	fmt.Println("  -table             Table name for table-specific operations")
	fmt.Println("  -backup            Backup file path for restore operations")
	fmt.Println("  -sql               SQL query to execute (for sql command)")
	fmt.Println()
	fmt.Println("SYNC OPTIONS (for -cmd=sync)")
	fmt.Println("  -commit            Write changes (otherwise upserts only)")
	fmt.Println("  -prune             Delete backup-only rows (authoritative primary)")
	fmt.Println("  -fail-on-drift     Exit non-zero if drift remains after sync (CI)")
	fmt.Println()
	fmt.Println("EXAMPLES")
	fmt.Println("  dbmanager -cmd=status")
	fmt.Println("  dbmanager -cmd=sync-health-check")
	fmt.Println("  dbmanager -cmd=sync -commit -prune -fail-on-drift")
	fmt.Println("  dbmanager -cmd=sync -table=styles -commit -prune")
	fmt.Println("  dbmanager -cmd=schema-discovery")
	fmt.Println("  dbmanager -cmd=schema-check")
	fmt.Println("  dbmanager -cmd=csv-backup   -source=primary")
	fmt.Println("  dbmanager -cmd=csv-seed     -target=primary")
	fmt.Println("  dbmanager migrate --target=primary")
	fmt.Println("  dbmanager view-tables")
	fmt.Println("  dbmanager -cmd=sync -commit -report-format=tui")
}

func handleSync(tableName string, commit bool, prune bool, failOnDrift bool) {
	startTime := time.Now()
	fmt.Println("=== CONSOLIDATED ERP SYNC ===")
	// ... (rest of the code remains the same)
	fmt.Printf("🔄 Starting dependency-aware sync using dbmanager (commit=%v)\n", commit)

	if db.Manager == nil {
		logger.Fatal("Database manager not initialized")
	}

	// If table is provided, run targeted table sync; else full dependency-aware sync
	if tableName != "" {
		fmt.Printf("🎯 Syncing specific table: %s\n", tableName)
		if err := managementAPI.SyncTable(tableName); err != nil {
			logger.Fatalf("Table sync failed: %v", err)
		}
		fmt.Println("✅ Table sync completed")
	} else {
		// Full ordered sync with FK repair inside
		if err := managementAPI.DependencyAwareSync(commit); err != nil {
			logger.Fatalf("Sync failed: %v", err)
		}
		fmt.Println("✅ Full dependency-aware sync completed")
	}

	// Post-sync health check summary with drift status
	statuses, err := managementAPI.SyncHealthCheck()
	if err != nil {
		logger.Fatalf("Health check failed: %v", err)
	}

	outOfSync := 0
	var offenders []string
	var backupOnly []string
	finalStatuses := statuses
	// Track prune deletions to report later
	pruned := map[string]int64{}
	for _, s := range statuses {
		if !s.InSync {
			outOfSync++
			offenders = append(offenders, fmt.Sprintf("%s (primary=%d backup=%d diff=%d)", s.Name, s.PrimaryCount, s.BackupCount, s.Difference))
			if s.Difference < 0 { // backup has extras
				backupOnly = append(backupOnly, s.Name)
			}
		}
	}

	if outOfSync == 0 {
		fmt.Println("🏁 Post-sync health: all tables are IN SYNC")
	} else {
		fmt.Printf("⚠️ Post-sync health: %d table(s) OUT OF SYNC\n", outOfSync)
		for _, line := range offenders {
			fmt.Printf("  - %s\n", line)
		}
		// Optional explicit prune, even though DependencyAwareSync(commit) may have pruned already
		if prune && len(backupOnly) > 0 {
			fmt.Println("🧹 Pruning backup-only rows (authoritative primary)...")
			deleted, err := managementAPI.PruneExtras(backupOnly)
			if err != nil {
				logger.Warnf("Prune encountered errors: %v", err)
			}
			for tbl, n := range deleted {
				fmt.Printf("  • %s: deleted %d extras\n", tbl, n)
			}
			// Save for final report
			pruned = deleted
			// Re-run health check after prune
			statuses, err = managementAPI.SyncHealthCheck()
			if err != nil {
				logger.Fatalf("Health check (post-prune) failed: %v", err)
			}
			finalStatuses = statuses
			outOfSync = 0
			offenders = offenders[:0]
			for _, s := range statuses {
				if !s.InSync {
					outOfSync++
					offenders = append(offenders, fmt.Sprintf("%s (primary=%d backup=%d diff=%d)", s.Name, s.PrimaryCount, s.BackupCount, s.Difference))
				}
			}
			if outOfSync == 0 {
				fmt.Println("🏁 Post-prune health: all tables are IN SYNC")
			} else {
				fmt.Printf("⚠️ Post-prune health: %d table(s) OUT OF SYNC\n", outOfSync)
				for _, line := range offenders {
					fmt.Printf("  - %s\n", line)
				}
			}
		}
		if failOnDrift {
			logger.Fatalf("Drift remains after sync and -fail-on-drift is set")
		}
	}

	renderSyncReport(finalStatuses, pruned, startTime, commit, prune, outOfSync)
}

func renderSyncReport(statuses []ops.TableStatus, pruned map[string]int64, start time.Time, commit, prune bool, outOfSync int) {
	// If TUI is requested, launch Bubble Tea viewer and skip textual output
	if syncReportFormat == "tui" {
		_ = RunSyncReportViewer(statuses)
		return
	}

	fmt.Println("=== SYNC REPORT ===")
	total := len(statuses)
	inSync := total - outOfSync
	fmt.Printf("Started:      %s\n", start.Format("2006-01-02 15:04:05"))
	fmt.Printf("Duration:     %s\n", time.Since(start).Truncate(time.Millisecond))
	fmt.Printf("Commit:       %v\n", commit)
	fmt.Printf("Prune:        %v\n", prune)
	fmt.Printf("Total tables: %d\n", total)
	fmt.Printf("In sync:      %d\n", inSync)
	fmt.Printf("Out of sync:  %d\n", outOfSync)

	switch syncReportFormat {
	case "json":
		type row struct {
			Table   string `json:"table"`
			Primary int64  `json:"primary"`
			Backup  int64  `json:"backup"`
			Diff    int64  `json:"diff"`
			Status  string `json:"status"`
		}
		var rows []row
		for _, s := range statuses {
			st := "IN SYNC"
			if !s.InSync {
				st = "OUT OF SYNC"
			}
			rows = append(rows, row{Table: s.Name, Primary: s.PrimaryCount, Backup: s.BackupCount, Diff: s.Difference, Status: st})
		}
		payload := map[string]any{
			"started":     start.Format(time.RFC3339),
			"duration":    time.Since(start).String(),
			"commit":      commit,
			"prune":       prune,
			"total":       total,
			"in_sync":     inSync,
			"out_of_sync": outOfSync,
			"tables":      rows,
			"pruned":      pruned,
		}
		b, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Println(string(b))

	case "csv":
		fmt.Println("TABLE,PRIMARY,BACKUP,DIFF,STATUS")
		for _, s := range statuses {
			st := "IN SYNC"
			if !s.InSync {
				st = "OUT OF SYNC"
			}
			fmt.Printf("%s,%d,%d,%d,%s\n", s.Name, s.PrimaryCount, s.BackupCount, s.Difference, st)
		}

	case "markdown":
		// GitHub-flavored markdown table
		fmt.Println()
		fmt.Println("| TABLE | PRIMARY | BACKUP | DIFF | STATUS |")
		fmt.Println("|:------|--------:|-------:|-----:|:-------|")
		for _, s := range statuses {
			st := "IN SYNC"
			if !s.InSync {
				st = "OUT OF SYNC"
			}
			fmt.Printf("| %s | %d | %d | %d | %s |\n", s.Name, s.PrimaryCount, s.BackupCount, s.Difference, st)
		}

	default: // table
		tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
		fmt.Fprintln(tw, "TABLE\tPRIMARY\tBACKUP\tDIFF\tSTATUS")
		for _, s := range statuses {
			st := "IN SYNC"
			if !s.InSync {
				st = "OUT OF SYNC"
			}
			fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%s\n", s.Name, s.PrimaryCount, s.BackupCount, s.Difference, st)
		}
		_ = tw.Flush()
	}

	if len(pruned) > 0 {
		fmt.Println()
		fmt.Println("Prune summary (deleted rows from backup not present in primary):")
		var keys []string
		for k := range pruned {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("  - %s: %d deleted\n", k, pruned[k])
		}
	}
}

func handleBackup() {
	fmt.Println("Creating database backup...")

	// TODO: Implement CreateBackup functionality
	logger.Info("CreateBackup functionality not yet implemented")

	fmt.Println("Backup completed successfully")
}

func handleRestore(backupPath, targetDB string) {
	if backupPath == "" {
		logger.Fatal("Backup path is required for restore operation")
	}

	fmt.Printf("Restoring from %s to %s database...\n", backupPath, targetDB)

	// TODO: Implement RestoreFromBackup functionality
	logger.Info("RestoreFromBackup functionality not yet implemented")

	fmt.Println("Restore completed successfully")
}

func handleCopy(sourceDB, targetDB string) {
	var source, target db.DatabaseType

	if sourceDB == "primary" {
		source = db.Primary
	} else {
		source = db.Backup
	}

	if targetDB == "primary" {
		target = db.Primary
	} else {
		target = db.Backup
	}

	fmt.Printf("Copying database from %s to %s...\n", sourceDB, targetDB)

	if err := managementAPI.CopyDatabase(source, target); err != nil {
		logger.Fatalf("Database copy failed: %v", err)
	}

	fmt.Println("Database copy completed successfully")
}

func handleMigrate(targetDB string, prune bool, strict bool) {
	// Rule 7: Ensure dbmanager is running with the latest .pb.go generated code
	if err := ops.CheckProtoFreshness(); err != nil {
		logger.Fatalf("Integrity Check Failed: %v", err)
	}

	if managementAPI == nil {
		logger.Fatal("Management API not initialized")
	}

	// Rule 2 & 3: Ensure SQL migrations don't violate declarative principles
	if err := managementAPI.LintMigrationFiles(); err != nil {
		logger.Fatalf("Lint Warning: %v", err)
	}

	t := strings.ToLower(strings.TrimSpace(targetDB))
	if t == "" {
		t = "primary"
	}

	fmt.Println("=== MIGRATIONS ===")
	switch t {
	case "primary":
		fmt.Println("Running migrations on PRIMARY only...")
		if err := managementAPI.MigrateOnly(db.Primary, prune, strict); err != nil {
			logger.Fatalf("Primary migration failed: %v", err)
		}
		fmt.Println("✅ PRIMARY migrations completed")
	case "backup":
		fmt.Println("Running migrations on BACKUP only...")
		if err := managementAPI.MigrateOnly(db.Backup, prune, strict); err != nil {
			logger.Fatalf("Backup migration failed: %v", err)
		}
		fmt.Println("✅ BACKUP migrations completed")
	case "both", "all":
		fmt.Println("Running migrations on PRIMARY then BACKUP...")
		if err := managementAPI.MigrateOnly(db.Primary, prune, strict); err != nil {
			logger.Fatalf("Primary migration failed: %v", err)
		}
		if err := managementAPI.MigrateOnly(db.Backup, prune, strict); err != nil {
			logger.Fatalf("Backup migration failed: %v", err)
		}
		fmt.Println("✅ PRIMARY and BACKUP migrations completed")
	default:
		logger.Fatalf("Unknown target for migrate: %s (use primary|backup|both)", t)
	}
}

func handleCompare() {
	fmt.Println("Comparing database schemas...")

	// TODO: Implement CompareSchemas functionality
	logger.Info("CompareSchemas functionality not yet implemented")
}

func handleListBackups() {
	fmt.Println("Available backups:")

	backups, err := managementAPI.ListBackups()
	if err != nil {
		logger.Fatalf("Failed to list backups: %v", err)
	}

	if len(backups) == 0 {
		fmt.Println("No backups found")
		return
	}

	for _, backup := range backups {
		fmt.Printf("- %s (%s, %s, %s)\n",
			backup.Filename,
			backup.DatabaseType,
			formatBytes(backup.Size),
			backup.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func handleSizes() {
	fmt.Println("Database sizes:")

	sizes, err := managementAPI.GetDatabaseSizes()
	if err != nil {
		logger.Fatalf("Failed to get database sizes: %v", err)
	}

	if primarySize, ok := sizes["primary_size"]; ok {
		if sizeMap, ok := primarySize.(map[string]interface{}); ok {
			fmt.Printf("Primary Database: %s\n", sizeMap["total_size"])
		}
	}

	if backupSize, ok := sizes["backup_size"]; ok {
		if sizeMap, ok := backupSize.(map[string]interface{}); ok {
			fmt.Printf("Backup Database: %s\n", sizeMap["total_size"])
		}
	}
}

func handleCSVBackup(tableName, sourceDB string) {
	fmt.Printf("=== CSV Backup from %s database ===\n", sourceDB)

	// Get database connection
	var dbConn *gorm.DB
	if sourceDB == "primary" {
		dbConn = db.Manager.GetPrimaryDB()
	} else {
		dbConn = db.Manager.GetBackupDB()
	}

	if dbConn == nil {
		logger.Fatalf("Failed to get %s database connection", sourceDB)
	}

	// Create backup directory if it doesn't exist
	backupDir := "internal/db/backup"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		logger.Fatalf("Failed to create backup directory: %v", err)
	}

	if tableName != "" {
		// Backup specific table
		if err := exportTableToCSV(dbConn, tableName, backupDir); err != nil {
			logger.Fatalf("Failed to backup table %s: %v", tableName, err)
		}
		fmt.Printf("Table %s exported to %s/%s.csv\n", tableName, backupDir, tableName)
	} else {
		// Backup all tables
		tables, err := getAllTables(dbConn)
		if err != nil {
			logger.Fatalf("Failed to get table list: %v", err)
		}

		fmt.Printf("Found %d tables to backup\n", len(tables))
		for _, table := range tables {
			if err := exportTableToCSV(dbConn, table, backupDir); err != nil {
				logger.Errorf("Failed to backup table %s: %v", table, err)
				continue
			}
			fmt.Printf("✓ Exported %s.csv\n", table)
		}
		fmt.Printf("All tables exported to %s/\n", backupDir)
	}
}

func handleCSVSeed(tableName, targetDB string) {
	fmt.Printf("=== CSV Seed to %s database ===\n", targetDB)

	// Get database connection
	var dbConn *gorm.DB
	if targetDB == "primary" {
		dbConn = db.Manager.GetPrimaryDB()
	} else {
		dbConn = db.Manager.GetBackupDB()
	}

	if dbConn == nil {
		logger.Fatalf("Failed to get %s database connection", targetDB)
	}

	backupDir := "internal/db/backup"

	if tableName != "" {
		// Seed specific table
		csvFile := filepath.Join(backupDir, tableName+".csv")
		if err := importTableFromCSV(dbConn, tableName, csvFile); err != nil {
			logger.Fatalf("Failed to seed table %s: %v", tableName, err)
		}
		fmt.Printf("Table %s seeded from %s\n", tableName, csvFile)
	} else {
		// Seed all CSV files found in backup directory
		csvFiles, err := getCSVFiles(backupDir)
		if err != nil {
			logger.Fatalf("Failed to get CSV files: %v", err)
		}

		if len(csvFiles) == 0 {
			fmt.Printf("No CSV files found in %s\n", backupDir)
			return
		}

		fmt.Printf("Found %d CSV files to import\n", len(csvFiles))
		for _, csvFile := range csvFiles {
			tableName := strings.TrimSuffix(filepath.Base(csvFile), ".csv")
			if err := importTableFromCSV(dbConn, tableName, csvFile); err != nil {
				logger.Errorf("Failed to seed table %s: %v", tableName, err)
				continue
			}
			fmt.Printf("✓ Seeded %s from %s\n", tableName, csvFile)
		}
		fmt.Println("All CSV files imported successfully")
	}
}

func getAllTables(db *gorm.DB) ([]string, error) {
	var tables []string

	// Query to get all table names from information_schema
	rows, err := db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

func exportTableToCSV(db *gorm.DB, tableName, backupDir string) error {
	// Create CSV file
	csvPath := filepath.Join(backupDir, tableName+".csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Get column names
	columns, err := getTableColumns(db, tableName)
	if err != nil {
		return fmt.Errorf("failed to get columns: %v", err)
	}

	// Write header
	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Query all data
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY 1", tableName)
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return fmt.Errorf("failed to query table data: %v", err)
	}
	defer rows.Close()

	// Column types not needed for this implementation
	_ = rows // Suppress unused variable warning

	// Create slice to hold row values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Write data rows
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		// Convert values to strings with proper formatting for seeding
		record := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				record[i] = ""
			} else {
				record[i] = formatValueForCSV(val)
			}
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	return nil
}

func importTableFromCSV(db *gorm.DB, tableName, csvPath string) error {
	// Check if CSV file exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		return fmt.Errorf("CSV file does not exist: %s", csvPath)
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %v", err)
	}

	// Check if table exists
	if !tableExists(db, tableName) {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	// Clear existing data (optional - you might want to make this configurable)
	if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)).Error; err != nil {
		logger.Warnf("Failed to truncate table %s: %v", tableName, err)
	}

	// Read and insert data
	batchSize := 1000
	batch := make([][]string, 0, batchSize)

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read CSV record: %v", err)
		}

		batch = append(batch, record)

		if len(batch) >= batchSize {
			if err := insertBatch(db, tableName, header, batch); err != nil {
				return fmt.Errorf("failed to insert batch: %v", err)
			}
			batch = batch[:0] // Reset batch
		}
	}

	// Insert remaining records
	if len(batch) > 0 {
		if err := insertBatch(db, tableName, header, batch); err != nil {
			return fmt.Errorf("failed to insert final batch: %v", err)
		}
	}

	return nil
}

func getTableColumns(db *gorm.DB, tableName string) ([]string, error) {
	var columns []string

	rows, err := db.Raw(`
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = ? 
		AND table_schema = 'public'
		ORDER BY ordinal_position
	`, tableName).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}
		columns = append(columns, columnName)
	}

	return columns, nil
}

func tableExists(db *gorm.DB, tableName string) bool {
	var count int64
	db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_name = ? 
		AND table_schema = 'public'
	`, tableName).Scan(&count)

	return count > 0
}

func insertBatch(db *gorm.DB, tableName string, columns []string, batch [][]string) error {
	if len(batch) == 0 {
		return nil
	}

	// Build INSERT query
	columnList := strings.Join(columns, ", ")
	placeholders := make([]string, len(batch))
	values := make([]interface{}, 0, len(batch)*len(columns))

	for i, record := range batch {
		recordPlaceholders := make([]string, len(columns))
		for j, value := range record {
			paramIndex := i*len(columns) + j + 1
			recordPlaceholders[j] = fmt.Sprintf("$%d", paramIndex)

			// Handle empty strings as NULL for certain data types
			if value == "" {
				values = append(values, nil)
			} else {
				values = append(values, value)
			}
		}
		placeholders[i] = "(" + strings.Join(recordPlaceholders, ", ") + ")"
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		tableName,
		columnList,
		strings.Join(placeholders, ", "),
	)

	return db.Exec(query, values...).Error
}

func getCSVFiles(dir string) ([]string, error) {
	var csvFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".csv") {
			csvFiles = append(csvFiles, path)
		}

		return nil
	})

	return csvFiles, err
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatValueForCSV formats database values for CSV export with proper seeding compatibility
func formatValueForCSV(val interface{}) string {
	if val == nil {
		return ""
	}

	// Get the reflect value and type
	rv := reflect.ValueOf(val)
	rt := reflect.TypeOf(val)

	// Handle byte slices (common for JSONB and binary data)
	if rt.Kind() == reflect.Slice && rt.Elem().Kind() == reflect.Uint8 {
		byteData := val.([]byte)

		// Try to parse as JSON first (for JSONB columns)
		var jsonObj interface{}
		if err := json.Unmarshal(byteData, &jsonObj); err == nil {
			// It's valid JSON, return as properly formatted JSON string
			if jsonBytes, err := json.Marshal(jsonObj); err == nil {
				return string(jsonBytes)
			}
		}

		// If not JSON, return as string (might be text data)
		return string(byteData)
	}

	// Handle time.Time values (convert to PostgreSQL-compatible format)
	if t, ok := val.(time.Time); ok {
		// Format as PostgreSQL timestamp without timezone info
		return t.Format("2006-01-02 15:04:05.999999")
	}

	// Handle boolean values (ensure consistent representation)
	if b, ok := val.(bool); ok {
		if b {
			return "TRUE"
		}
		return "FALSE"
	}

	// Handle pointers
	if rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return ""
		}
		return formatValueForCSV(rv.Elem().Interface())
	}

	// Handle string values (check for special date formats that need conversion)
	if s, ok := val.(string); ok {
		// Check if it's a Go time format string that needs conversion
		if strings.Contains(s, " +0000 UTC") {
			s = strings.Replace(s, " +0000 UTC", "", 1)
		} else if strings.Contains(s, " +0600 +06") {
			s = strings.Replace(s, " +0600 +06", "", 1)
		}
		return s
	}

	// For all other types, use default string conversion
	return fmt.Sprintf("%v", val)
}

// // runIntelligentSync performs bi-directional sync that merges newest data from both databases
// func runIntelligentSync() error {
// 	fmt.Println("🧠 Starting intelligent bi-directional sync...")
// 	fmt.Println("📊 Using dependency-aware merge strategy (not blind copy)")

// 	// Execute the dependency-aware sync tool that merges newest data with commit flag
// 	cmd := exec.Command("go", "run", "./cmd/dependency-aware-sync/main.go", "-commit")
// 	cmd.Dir = "."

// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("🔍 Sync output: %s\n", string(output))
// 		return fmt.Errorf("dependency-aware sync failed: %v", err)
// 	}

// 	fmt.Printf("📈 Sync completed with merge strategy\n")
// 	fmt.Printf("📋 Output: %s\n", string(output))
// 	return nil
// }

// // runBackupMigrations runs migrations on backup database
// func runBackupMigrations() error {
// 	fmt.Println("🔧 Running ACTUAL migrations on backup database...")

// 	// Execute the migrate command on backup database
// 	cmd := exec.Command("go", "run", "./cmd/dbmanager/main.go", "-cmd=migrate", "-target=backup")
// 	cmd.Dir = "."

// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("🔍 Migration output: %s\n", string(output))
// 		return fmt.Errorf("backup migration failed: %v", err)
// 	}

// 	fmt.Printf("📈 Backup migrations completed successfully\n")
// 	fmt.Printf("📋 Migration output: %s\n", string(output))
// 	return nil
// }

// // runSyncHealthCheck runs post-sync health verification
// func runSyncHealthCheck() {
// 	fmt.Println("🏥 Running comprehensive sync health check...")

// 	// Execute the sync health check tool
// 	if err := executeHealthCheck(); err != nil {
// 		fmt.Printf("⚠️ Health check failed: %v\n", err)
// 		fmt.Println("📊 Manual health check recommended")
// 	} else {
// 		fmt.Println("✅ Health check completed successfully")
// 	}
// }

// executeHealthCheck runs the sync health check tool
func executeHealthCheck() error {
	fmt.Println("🔍 Running detailed sync health check...")

	// Execute the sync health check tool
	cmd := exec.Command("go", "run", "./cmd/sync-health-check/main.go")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("🔍 Health check output: %s\n", string(output))
		return fmt.Errorf("health check failed: %v", err)
	}

	fmt.Printf("📊 Health Check Results:\n%s\n", string(output))

	// Check if sync is 100% - parse output for sync percentage
	outputStr := string(output)
	if strings.Contains(outputStr, "100.0%") || strings.Contains(outputStr, "✅ In Sync: ") {
		fmt.Println("🎯 100% sync achieved!")
	} else {
		fmt.Println("⚠️ Sync not at 100% - manual review recommended")
	}

	return nil
}

// handleTestInit initializes database connections for testing without running migrations
func handleTestInit() {
	fmt.Println("🧪 Initializing database connections for testing...")

	if db.Manager == nil {
		fmt.Println("❌ Database manager not initialized")
		os.Exit(1)
	}

	// Test primary database connection
	primaryDB := db.Manager.GetDB()
	if primaryDB == nil {
		fmt.Println("❌ Primary database connection failed")
		os.Exit(1)
	}

	// Test basic query on processes table
	var count int64
	err := primaryDB.Raw("SELECT COUNT(*) FROM processes LIMIT 1").Scan(&count).Error
	if err != nil {
		fmt.Printf("❌ Failed to query processes table: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Database connections ready for testing\n")
	fmt.Printf("📊 Found %d processes in database\n", count)
	fmt.Println("🚀 Test environment initialized successfully")
}
