// Package main implements dbops — the InsureTech database operations entrypoint.
//
// Design rationale: The gateway is a stateless HTTP→gRPC proxy and must NEVER
// own database connections. All DB lifecycle work (migrate, sync, backup) is
// delegated here, run as isolated jobs or sidecars before/alongside services.
//
// Usage:
//
//	# Run once at startup — blocks until done, exits 0 (success) or 1 (failure)
//	go run ./backend/inscore/cmd/dbops migrate --target=both
//
//	# Long-running background sidecar — sync + auto-backup loops
//	go run ./backend/inscore/cmd/dbops sidecar
//
//	# Individual sidecar modes
//	go run ./backend/inscore/cmd/dbops sync-watch   # continuous sync loop
//	go run ./backend/inscore/cmd/dbops backup-watch # continuous backup loop
//
//	# One-shot ops
//	go run ./backend/inscore/cmd/dbops sync-now
//	go run ./backend/inscore/cmd/dbops backup-now
//	go run ./backend/inscore/cmd/dbops status
//	go run ./backend/inscore/cmd/dbops validate
//	go run ./backend/inscore/cmd/dbops lint
//
// Docker Compose usage (recommended):
//
//	db-migrate:
//	  command: go run ./backend/inscore/cmd/dbops migrate --target=both
//	  depends_on: { postgres: { condition: service_healthy } }
//	  restart: "no"
//
//	gateway:
//	  depends_on: { db-migrate: { condition: service_completed_successfully } }
//
//	db-sidecar:
//	  command: go run ./backend/inscore/cmd/dbops sidecar
//	  restart: unless-stopped
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/db/ops"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	opsconfig "github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/spf13/cobra"

	// Register all proto entity packages so the UnifiedMigrationManager can
	// discover schemas from protobuf descriptors (same pattern as dbmanager).
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

// managementAPI is the global ops handle, set during PreRunE after DB init.
var managementAPI *ops.ManagementAPI

func main() {
	// ── Logger: console-only, no file output ──────────────────────────────────
	logCfg := appLogger.NoFileConfig()
	logCfg.Level = "info"
	_ = appLogger.Initialize(logCfg)

	// ── Load .env (walks up to project root automatically) ───────────────────
	if err := env.Load(); err != nil {
		appLogger.Warn("No .env file found; using environment variables only")
	}

	if err := newRootCmd().Execute(); err != nil {
		appLogger.Errorf("dbops: %v", err)
		os.Exit(1)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Root command
// ─────────────────────────────────────────────────────────────────────────────

func newRootCmd() *cobra.Command {
	var cfgPath string

	root := &cobra.Command{
		Use:   "dbops",
		Short: "InsureTech DB operations — migrate, sync, backup, status",
		Long: `dbops is the database operations entrypoint for InsureTech.

It is designed to run OUTSIDE the gateway so that the API gateway stays
stateless. Run 'migrate' once at startup (as a Docker init-container or
Kubernetes Job), and 'sidecar' as a long-running companion service for
continuous sync and backup.`,
		Example: strings.Join([]string{
			"  dbops migrate --target=both          # startup job (blocks, exits 0/1)",
			"  dbops sidecar                        # continuous sync+backup daemon",
			"  dbops sync-watch                     # sync loop only",
			"  dbops backup-watch                   # backup loop only",
			"  dbops sync-now                       # immediate one-shot sync",
			"  dbops backup-now                     # immediate one-shot backup",
			"  dbops status                         # print DB health to stdout",
			"  dbops validate                       # schema consistency check",
			"  dbops lint                           # lint migration SQL files",
		}, "\n"),
		SilenceErrors: true,
		SilenceUsage:  false,
	}

	root.PersistentFlags().StringVar(&cfgPath, "config", "database.yaml",
		"Path to database configuration file (resolved via project root)")

	// ── Sub-commands ──────────────────────────────────────────────────────────
	root.AddCommand(
		newMigrateCmd(&cfgPath),
		newSidecarCmd(&cfgPath),
		newSyncWatchCmd(&cfgPath),
		newBackupWatchCmd(&cfgPath),
		newSyncNowCmd(&cfgPath),
		newBackupNowCmd(&cfgPath),
		newStatusCmd(&cfgPath),
		newValidateCmd(&cfgPath),
		newLintCmd(),
	)

	return root
}

// ─────────────────────────────────────────────────────────────────────────────
// migrate — run once at startup, exits 0 or 1
// ─────────────────────────────────────────────────────────────────────────────

func newMigrateCmd(cfgPath *string) *cobra.Command {
	var target string
	var prune bool
	var strict bool

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run DB migrations (blocking — exits 0 on success, 1 on failure)",
		Long: `Runs the full migration flow (proto → SQL → seeders) against the selected
database(s). Designed to run ONCE at startup before any microservice starts.
Blocks until complete and exits 0 (success) or 1 (failure) so Docker Compose
and Kubernetes can gate dependent services on successful completion.

Primary is always migrated first (authoritative source of truth).
Backup migration failure is non-fatal — a warning is logged and the job
still exits 0 so the gateway is not blocked by backup issues.`,
		Example: strings.Join([]string{
			"  dbops migrate                        # both DBs (default)",
			"  dbops migrate --target=primary       # primary only",
			"  dbops migrate --target=backup        # backup only",
			"  dbops migrate --target=both --prune  # prune zombie columns",
			"  dbops migrate --strict               # fail on schema drift",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initDBTargets(*cfgPath, target)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate(target, prune, strict)
		},
	}

	cmd.Flags().StringVar(&target, "target", "both",
		"Target database: primary | backup | both")
	cmd.Flags().BoolVar(&prune, "prune", false,
		"Remove columns absent from proto definitions")
	cmd.Flags().BoolVar(&strict, "strict", false,
		"Exit non-zero on any schema drift detected")

	return cmd
}

func runMigrate(target string, prune, strict bool) error {
	t := strings.ToLower(strings.TrimSpace(target))
	start := time.Now()
	appLogger.Infof("▶ dbops migrate  target=%s prune=%v strict=%v", t, prune, strict)

	switch t {
	case "primary":
		if err := managementAPI.MigrateOnly(db.Primary, prune, strict); err != nil {
			return fmt.Errorf("primary migration failed: %w", err)
		}

	case "backup":
		if err := managementAPI.MigrateOnly(db.Backup, prune, strict); err != nil {
			return fmt.Errorf("backup migration failed: %w", err)
		}

	case "both", "all", "":
		// Primary first — authoritative source of truth
		appLogger.Info("  [1/2] Migrating primary database...")
		if err := managementAPI.MigrateOnly(db.Primary, prune, strict); err != nil {
			return fmt.Errorf("primary migration failed: %w", err)
		}
		// Backup second — schema-identical to primary but non-fatal on failure
		appLogger.Info("  [2/2] Migrating backup database...")
		if err := managementAPI.MigrateOnly(db.Backup, prune, strict); err != nil {
			appLogger.Warnf("  ⚠ backup migration failed (non-fatal): %v", err)
			appLogger.Warn("  Backup DB may be out of sync — run 'dbops migrate --target=backup' manually.")
		}

	default:
		return fmt.Errorf("unknown --target %q — must be primary | backup | both", t)
	}

	appLogger.Infof("✅ migrate complete  duration=%s", time.Since(start).Round(time.Millisecond))
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// sidecar — long-running: sync + backup loops (SIGTERM-aware)
// ─────────────────────────────────────────────────────────────────────────────

func newSidecarCmd(cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sidecar",
		Short: "Run sync + backup loops as a long-running sidecar process",
		Long: `Starts both the primary→backup sync loop and the auto-backup loop as
concurrent goroutines. Blocks until SIGTERM or SIGINT is received, then
exits cleanly. Does NOT run migrations — use 'dbops migrate' for that.

Intervals are read from database.yaml. Designed to run as a Docker Compose
service (restart: unless-stopped) or a Kubernetes Deployment sidecar.`,
		Example: "  dbops sidecar",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initDBTargets(*cfgPath, "both")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSidecar()
		},
	}
	return cmd
}

func runSidecar() error {
	appLogger.Info("▶ dbops sidecar starting (sync + backup loops)...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		appLogger.Info("  [sync-watch] goroutine started")
		managementAPI.StartSync()
	}()

	go func() {
		appLogger.Info("  [backup-watch] goroutine started")
		managementAPI.StartAutoBackup()
	}()

	appLogger.Info("✅ sidecar running — send SIGTERM or SIGINT to stop")
	<-quit

	appLogger.Info("🛑 sidecar shutting down...")
	if db.Manager != nil {
		_ = db.Manager.Close()
	}
	appLogger.Info("sidecar stopped cleanly")
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// sync-watch — sync loop only (SIGTERM-aware)
// ─────────────────────────────────────────────────────────────────────────────

func newSyncWatchCmd(cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-watch",
		Short: "Run the database sync loop only (continuous, SIGTERM-aware)",
		Long: `Starts the primary→backup sync loop only. Interval is read from
database.yaml (Database.Sync.Interval). Blocks until SIGTERM/SIGINT.`,
		Example: "  dbops sync-watch",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initDBTargets(*cfgPath, "both")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			appLogger.Info("▶ dbops sync-watch started")
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
			go managementAPI.StartSync()
			<-quit
			appLogger.Info("🛑 sync-watch stopping...")
			if db.Manager != nil {
				_ = db.Manager.Close()
			}
			return nil
		},
	}
	return cmd
}

// ─────────────────────────────────────────────────────────────────────────────
// backup-watch — backup loop only (SIGTERM-aware)
// ─────────────────────────────────────────────────────────────────────────────

func newBackupWatchCmd(cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup-watch",
		Short: "Run the auto-backup loop only (continuous, SIGTERM-aware)",
		Long: `Starts the automatic backup loop only. Interval is read from database.yaml
(Database.BackupSettings.BackupInterval). Blocks until SIGTERM/SIGINT.`,
		Example: "  dbops backup-watch",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initDBTargets(*cfgPath, "both")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			appLogger.Info("▶ dbops backup-watch started")
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
			go managementAPI.StartAutoBackup()
			<-quit
			appLogger.Info("🛑 backup-watch stopping...")
			if db.Manager != nil {
				_ = db.Manager.Close()
			}
			return nil
		},
	}
	return cmd
}

// ─────────────────────────────────────────────────────────────────────────────
// sync-now — one-shot immediate sync
// ─────────────────────────────────────────────────────────────────────────────

func newSyncNowCmd(cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-now",
		Short: "Run a one-shot immediate sync (primary → backup) and exit",
		Long: `Performs a single synchronisation pass from primary to backup and exits.
Useful for manual operations or CI pipeline steps.`,
		Example: "  dbops sync-now",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initDBTargets(*cfgPath, "both")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()
			appLogger.Info("▶ dbops sync-now")
			if err := managementAPI.SyncNow(); err != nil {
				return fmt.Errorf("sync failed: %w", err)
			}
			appLogger.Infof("✅ sync-now complete  duration=%s", time.Since(start).Round(time.Millisecond))
			return nil
		},
	}
	return cmd
}

// ─────────────────────────────────────────────────────────────────────────────
// backup-now — one-shot immediate backup
// ─────────────────────────────────────────────────────────────────────────────

func newBackupNowCmd(cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup-now",
		Short: "Run a one-shot immediate backup of the primary database and exit",
		Long: `Triggers a single pg_dump backup of the primary database and exits.
The backup is written to the path configured in database.yaml
(Database.BackupSettings.BackupPath).`,
		Example: "  dbops backup-now",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initDBTargets(*cfgPath, "primary")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()
			appLogger.Info("▶ dbops backup-now")
			if err := managementAPI.BackupNow(); err != nil {
				return fmt.Errorf("backup failed: %w", err)
			}
			appLogger.Infof("✅ backup-now complete  duration=%s", time.Since(start).Round(time.Millisecond))
			return nil
		},
	}
	return cmd
}

// ─────────────────────────────────────────────────────────────────────────────
// status — print DB health (human or JSON)
// ─────────────────────────────────────────────────────────────────────────────

func newStatusCmd(cfgPath *string) *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Print database connection status and metrics",
		Long: `Connects to both databases and prints their live status. Use --json for
machine-readable output suitable for health probes and monitoring scripts.`,
		Example: strings.Join([]string{
			"  dbops status",
			"  dbops status --json",
		}, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initDBTargets(*cfgPath, "both")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(jsonOut)
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output status as JSON")
	return cmd
}

func runStatus(jsonOut bool) error {
	if db.Manager == nil {
		return fmt.Errorf("database manager not initialized")
	}

	primaryStatus, backupStatus := db.Manager.GetStatus()
	healthy := db.Manager.IsHealthy()
	metrics := db.Manager.GetMetrics()

	if jsonOut {
		out := map[string]interface{}{
			"healthy":            healthy,
			"primary_status":     string(primaryStatus),
			"backup_status":      string(backupStatus),
			"primary_conns":      metrics.PrimaryConnections,
			"backup_conns":       metrics.BackupConnections,
			"failover_count":     metrics.FailoverCount,
			"last_failover_time": metrics.LastFailoverTime,
			"last_sync_time":     metrics.LastSyncTime,
			"last_backup_time":   metrics.LastBackupTime,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	// Human-readable output
	statusIcon := func(s db.ConnectionStatus) string {
		if s == db.StatusHealthy {
			return "✅"
		}
		return "❌"
	}

	fmt.Println("════════════════════════════════════")
	fmt.Println("  dbops status")
	fmt.Println("════════════════════════════════════")
	fmt.Printf("  Overall healthy   : %v\n", healthy)
	fmt.Printf("  Primary DB        : %s %s  (open conns: %d)\n",
		statusIcon(primaryStatus), primaryStatus, metrics.PrimaryConnections)
	fmt.Printf("  Backup DB         : %s %s  (open conns: %d)\n",
		statusIcon(backupStatus), backupStatus, metrics.BackupConnections)
	fmt.Printf("  Failover count    : %d\n", metrics.FailoverCount)
	if !metrics.LastFailoverTime.IsZero() {
		fmt.Printf("  Last failover     : %s\n", metrics.LastFailoverTime.Format(time.RFC3339))
	}
	if !metrics.LastSyncTime.IsZero() {
		fmt.Printf("  Last sync         : %s\n", metrics.LastSyncTime.Format(time.RFC3339))
	}
	if !metrics.LastBackupTime.IsZero() {
		fmt.Printf("  Last backup       : %s\n", metrics.LastBackupTime.Format(time.RFC3339))
	}
	fmt.Println("════════════════════════════════════")
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// validate — schema consistency check between primary and backup
// ─────────────────────────────────────────────────────────────────────────────

func newValidateCmd(cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate schema consistency between primary and backup databases",
		Long: `Compares the schema of primary and backup databases and reports any
mismatches (missing tables, missing columns, constraint differences).
Exits 0 if schemas match, 1 if mismatches are found.

Tip: run 'dbops migrate --target=backup' or 'dbmanager rebuild-backup'
to fix schema drift in the backup database.`,
		Example: "  dbops validate",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initDBTargets(*cfgPath, "both")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate()
		},
	}
	return cmd
}

func runValidate() error {
	appLogger.Info("▶ dbops validate")

	mismatches, err := db.Manager.ValidateSchemaConsistency()
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	if len(mismatches) == 0 {
		appLogger.Info("✅ Schema is consistent between primary and backup")
		return nil
	}

	appLogger.Errorf("❌ %d schema mismatch(es) found:", len(mismatches))
	for _, m := range mismatches {
		appLogger.Errorf("   • %s", m)
	}
	appLogger.Info("Tip: run 'dbops migrate --target=backup' or 'dbmanager rebuild-backup' to fix")
	return fmt.Errorf("schema validation failed: %d mismatch(es)", len(mismatches))
}

// ─────────────────────────────────────────────────────────────────────────────
// lint — lint migration SQL files (no DB connection required)
// ─────────────────────────────────────────────────────────────────────────────

func newLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint migration SQL files for forbidden patterns (no DB required)",
		Long: `Checks all migration SQL files for forbidden patterns such as raw DROP TABLE
without IF EXISTS, missing rollback annotations, and other common mistakes.
Does NOT require a database connection. Exits 0 if clean, 1 if issues found.

Safe to run in CI without any database credentials.`,
		Example: "  dbops lint",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Lint does not need a DB — pass a nil manager
			lintAPI := ops.NewManagementAPI(nil)
			appLogger.Info("▶ dbops lint")
			if err := lintAPI.LintMigrationFiles(); err != nil {
				return fmt.Errorf("lint failed: %w", err)
			}
			appLogger.Info("✅ All migration files passed lint checks")
			return nil
		},
	}
	return cmd
}

// ─────────────────────────────────────────────────────────────────────────────
// DB initialisation helper
// ─────────────────────────────────────────────────────────────────────────────

// initDBTargets connects to the specified database target(s) without running
// migrations, then sets the global db.Manager and managementAPI.
// target must be "primary", "backup", or "both".
func initDBTargets(cfgPath, target string) error {
	resolvedPath, err := opsconfig.ResolveConfigPath(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to resolve config path %q: %w", cfgPath, err)
	}

	dbConfig, err := db.LoadDatabaseConfig(resolvedPath)
	if err != nil {
		return fmt.Errorf("failed to load database config from %q: %w", resolvedPath, err)
	}

	manager := db.NewDatabaseManager(dbConfig)

	t := strings.ToLower(strings.TrimSpace(target))
	switch t {
	case "primary":
		if err := manager.ConnectWithoutMigrationsTargets(db.Primary); err != nil {
			return fmt.Errorf("failed to connect to primary DB: %w", err)
		}
	case "backup":
		if err := manager.ConnectWithoutMigrationsTargets(db.Backup); err != nil {
			return fmt.Errorf("failed to connect to backup DB: %w", err)
		}
	default: // "both", "all", ""
		if err := manager.ConnectWithoutMigrations(); err != nil {
			return fmt.Errorf("failed to connect to databases: %w", err)
		}
	}

	// Set globals — db.Manager is used by status, validate, and sidecar shutdown
	db.Manager = manager
	managementAPI = ops.NewManagementAPI(manager)
	return nil
}
