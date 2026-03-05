package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

func init() {
	schema.RegisterSerializer("proto_timestamp", ProtoTimestampSerializer{})
}

// DB is the global database connection (for backward compatibility)
var DB *gorm.DB

func InitDB(dsn string) error {
	dbOpLogger := appLogger.GetDatabaseOperationLogger()

	// Check if we should use the new database manager
	configPath := "configs/database.yaml"
	if _, err := os.Stat(configPath); err == nil {
		// Use new database manager
		appLogger.Info("Using enhanced database manager with multi-database support")
		dbOpLogger.LogServerStart("multi-database-manager")
		return InitializeManager(configPath)
	}

	// Fallback to legacy single database connection
	appLogger.Info("Using legacy single database connection")
	dbOpLogger.LogServerStart("legacy-single-database")
	return initLegacyDB(dsn)
}

// initLegacyDB initializes a single database connection (legacy mode)
func initLegacyDB(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable SQL query logging to prevent sensitive data exposure
	})
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	DB = db // Assign to global DB variable

	// Get SQL DB connection for extension and migrations
	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Fatalf("Failed to get SQL DB from GORM: %v", err)
		return err
	}

	// Configure connection pool to avoid prepared statement collisions
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Enable UUID extension (moved to ops package)
	// err = enableUUIDOSSP(sqlDB)
	// if err != nil {
	// 	appLogger.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	// 	return err
	// }

	// Register models with GORM for association resolution (no-op for ERP SQL migrations)
	if err := registerModels(db); err != nil {
		appLogger.Fatalf("Failed to register models: %v", err)
		return err
	}

	// Run SQL migrations/seeders frominsuretech folders recursively
	// (migrations now handled by ops package and dbmanager CLI)
	// if err := runSQLMigrationsFromDir(sqlDB, "lpc/db/migrations"); err != nil {
	// 	appLogger.Fatalf("Failed to run SQL migrations: %v", err)
	// 	return err
	// }

	// Run proto-driven migrations (ensure tables/columns from protobuf descriptors)
	// (migrations now handled by ops package and dbmanager CLI)
	// if err := runProtoMigrations(sqlDB); err != nil {
	// 	appLogger.Fatalf("Failed to run proto migrations: %v", err)
	// 	return err
	// }

	// if err := runSQLSeedersFromDir(sqlDB, "lpc/db/seeders"); err != nil {
	// 	appLogger.Warnf("Failed to run seeders: %v", err)
	// }

	appLogger.Info("Legacy database initialization completed (migrations should be run via dbmanager CLI)")
	appLogger.Info("Database connection established")
	return nil
}

// registerModels is intentionally a no-op; ERP schema managed via SQL migrations
func registerModels(db *gorm.DB) error { return nil }

// readAppliedFromLog parses internal/db/migration.log and returns set of names with the given prefix ("M:" for migrations, "S:" for seeders)
func readAppliedFromLog(path, prefix string) map[string]bool {
	m := make(map[string]bool)
	b, err := os.ReadFile(path)
	if err != nil {
		return m
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[1])
				if name != "" {
					m[name] = true
				}
			}
		}
	}
	return m
}

// appendToMigrationLog appends a timestamped entry to internal/db/migration.log
func appendToMigrationLog(path, prefix, name string) {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	ts := time.Now().Format(time.RFC3339)
	_, _ = f.WriteString(fmt.Sprintf("%s %s %s\n", prefix, name, ts))
}
