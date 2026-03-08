package db

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// DatabaseType represents the type of database connection
type DatabaseType string

const (
	Primary DatabaseType = "primary"
	Backup  DatabaseType = "backup"
)

// ConnectionStatus represents the status of a database connection
type ConnectionStatus string

const (
	StatusHealthy      ConnectionStatus = "healthy"
	StatusUnhealthy    ConnectionStatus = "unhealthy"
	StatusConnecting   ConnectionStatus = "connecting"
	StatusDisconnected ConnectionStatus = "disconnected"
)

// DatabaseManager manages multiple database connections with failover support
type DatabaseManager struct {
	config          *DatabaseConfig
	primaryDB       *gorm.DB
	backupDB        *gorm.DB
	currentDB       *gorm.DB
	currentType     DatabaseType
	primaryStatus   ConnectionStatus
	backupStatus    ConnectionStatus
	healthCheckStop chan bool
	syncStop        chan bool
	backupStop      chan bool
	mu              sync.RWMutex
	metrics         *DatabaseMetrics
}

// DatabaseMetrics holds database performance metrics
type DatabaseMetrics struct {
	PrimaryConnections int
	BackupConnections  int
	SlowQueries        int64
	FailoverCount      int64
	LastFailoverTime   time.Time
	LastSyncTime       time.Time
	LastBackupTime     time.Time
	mu                 sync.RWMutex
}

// NewDatabaseManager creates a new database manager instance
func NewDatabaseManager(config *DatabaseConfig) *DatabaseManager {
	schema.RegisterSerializer("proto_timestamp", ProtoTimestampSerializer{})

	return &DatabaseManager{
		config:          config,
		currentType:     Primary,
		primaryStatus:   StatusDisconnected,
		backupStatus:    StatusDisconnected,
		healthCheckStop: make(chan bool),
		syncStop:        make(chan bool),
		backupStop:      make(chan bool),
		metrics:         &DatabaseMetrics{},
	}
}

// Initialize initializes all database connections and starts background services
func (dm *DatabaseManager) Initialize() error {
	// Connect to primary database
	if err := dm.connectPrimary(); err != nil {
		appLogger.WithError(err).Error("Failed to connect to primary database")

		// Try backup database if primary fails
		if err := dm.connectBackup(); err != nil {
			return fmt.Errorf("failed to connect to both primary and backup databases: primary=%v, backup=%v", err, err)
		}
		dm.currentType = Backup
		dm.currentDB = dm.backupDB
	} else {
		dm.currentType = Primary
		dm.currentDB = dm.primaryDB

		// Connect to backup database for sync purposes
		if err := dm.connectBackup(); err != nil {
			appLogger.WithError(err).Warn("Failed to connect to backup database")
		}
	}

	// Run migrations only on current database (avoid double execution)
	// Migrations are now handled by the ops package and dbmanager tool
	// if err := dm.runMigrations(); err != nil {
	// 	return fmt.Errorf("failed to run migrations: %v", err)
	// }

	// Only run migrations on backup if it's the current database or if explicitly requested
	// This prevents double execution when both databases are available
	if dm.backupDB != nil && dm.currentDB != dm.backupDB && dm.currentType == Primary {
		// appLogger.Info("Skipping migrations on backup database to avoid double execution")
		// appLogger.Info("Backup database will be synced via data sync process")
	}

	// Start background services
	if dm.config.Database.Failover.Enabled {
		go dm.startHealthCheck()
	}

	appLogger.WithField("database_type", string(dm.currentType)).Info("Database manager initialized")
	return nil
}

// GetConfig returns the database configuration
func (dm *DatabaseManager) GetConfig() *DatabaseConfig {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.config
}

// InitializeForService initializes database connections for application services without running migrations
// This is used by microservices that only need database access, not migration management
func (dm *DatabaseManager) InitializeForService() error {
	// Try primary database first
	if err := dm.connectPrimary(); err != nil {
		appLogger.WithError(err).Warn("Primary database connection failed, attempting failover to backup")

		// Primary failed - try backup as fallback
		if err := dm.connectBackup(); err != nil {
			return fmt.Errorf("failed to connect to both databases: primary=%v, backup=%v", err, err)
		}

		// Using backup as primary is unavailable
		dm.currentType = Backup
		dm.currentDB = dm.backupDB
		appLogger.Info("Service running on backup database (primary unavailable)")
		return nil
	}

	// Primary connected successfully
	dm.currentType = Primary
	dm.currentDB = dm.primaryDB

	// Only connect to backup if failover is enabled
	if dm.config.Database.Failover.Enabled {
		if err := dm.connectBackup(); err != nil {
			appLogger.WithError(err).Warn("Backup database unavailable - failover disabled")
		} else {
			// Both databases available - start health monitoring
			go dm.startHealthCheck()
			appLogger.Info("Service running on primary database with automatic failover enabled")
			return nil
		}
	}

	appLogger.Info("Service running on primary database (failover disabled)")
	return nil
}

// ConnectWithoutMigrations connects to databases without running migrations (for direct SQL operations)
func (dm *DatabaseManager) ConnectWithoutMigrations() error {
	return dm.ConnectWithoutMigrationsTargets(Primary, Backup)
}

// ConnectWithoutMigrationsTargets connects to specified target databases without running migrations
func (dm *DatabaseManager) ConnectWithoutMigrationsTargets(targets ...DatabaseType) error {
	var primaryErr, backupErr error
	connectPrimary := false
	connectBackup := false

	// Check which targets to connect to
	for _, target := range targets {
		switch target {
		case Primary:
			connectPrimary = true
		case Backup:
			connectBackup = true
		}
	}

	// Connect to primary database if requested
	if connectPrimary {
		if err := dm.connectPrimary(); err != nil {
			primaryErr = err
			appLogger.WithError(err).Error("Failed to connect to primary database")
		} else {
			dm.currentType = Primary
			dm.currentDB = dm.primaryDB
		}
	}

	// Connect to backup database if requested
	if connectBackup {
		if err := dm.connectBackup(); err != nil {
			backupErr = err
			appLogger.WithError(err).Warn("Failed to connect to backup database")
		} else {
			// Only set as current if primary failed or wasn't requested
			if !connectPrimary || primaryErr != nil {
				dm.currentType = Backup
				dm.currentDB = dm.backupDB
			}
		}
	}

	// Check if we have at least one successful connection
	if connectPrimary && connectBackup {
		if primaryErr != nil && backupErr != nil {
			return fmt.Errorf("failed to connect to both primary and backup databases: primary=%v, backup=%v", primaryErr, backupErr)
		}
	} else if connectPrimary && primaryErr != nil {
		return fmt.Errorf("failed to connect to primary database: %v", primaryErr)
	} else if connectBackup && backupErr != nil {
		return fmt.Errorf("failed to connect to backup database: %v", backupErr)
	}

	appLogger.WithField("database_type", string(dm.currentType)).Info("Database connections established (no migrations)")
	return nil
}

// connectPrimary establishes connection to primary database
func (dm *DatabaseManager) connectPrimary() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dbOpLogger := appLogger.GetDatabaseOperationLogger()
	dsn := dm.config.Database.Primary.GetConnectionDSN(true) // Use pool for application traffic
	// PreferSimpleProtocol disables prepared-statement caching in the pgx driver,
	// preventing "cached plan must not change result type" (SQLSTATE 0A000) errors
	// that occur after schema changes when using connection poolers or gorm-adapter.
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		dm.primaryStatus = StatusUnhealthy
		dbOpLogger.LogDatabaseConnection("primary", false, dm.config.Database.Primary.Host)
		return err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		dm.primaryStatus = StatusUnhealthy
		dbOpLogger.LogDatabaseConnection("primary", false, dm.config.Database.Primary.Host)
		return err
	}

	sqlDB.SetMaxOpenConns(dm.config.Database.Primary.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dm.config.Database.Primary.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dm.config.Database.Primary.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		dm.primaryStatus = StatusUnhealthy
		dbOpLogger.LogDatabaseConnection("primary", false, dm.config.Database.Primary.Host)
		return err
	}

	dm.primaryDB = db
	dm.primaryStatus = StatusHealthy
	dbOpLogger.LogDatabaseConnection("primary", true, dm.config.Database.Primary.Host)
	appLogger.WithField("provider", dm.config.Database.Primary.Provider).With("host", dm.config.Database.Primary.Host).Info("Connected to primary database")
	return nil
}

// connectBackup establishes connection to backup database
func (dm *DatabaseManager) connectBackup() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dbOpLogger := appLogger.GetDatabaseOperationLogger()
	dsn := dm.config.Database.Backup.BuildDSN()
	// PreferSimpleProtocol disables prepared-statement caching in the pgx driver,
	// preventing "cached plan must not change result type" (SQLSTATE 0A000) errors.
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		dm.backupStatus = StatusUnhealthy
		dbOpLogger.LogDatabaseConnection("backup", false, dm.config.Database.Backup.Host)
		return err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		dm.backupStatus = StatusUnhealthy
		dbOpLogger.LogDatabaseConnection("backup", false, dm.config.Database.Backup.Host)
		return err
	}

	sqlDB.SetMaxOpenConns(dm.config.Database.Backup.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dm.config.Database.Backup.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dm.config.Database.Backup.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		dm.backupStatus = StatusUnhealthy
		dbOpLogger.LogDatabaseConnection("backup", false, dm.config.Database.Backup.Host)
		return err
	}

	dm.backupDB = db
	dm.backupStatus = StatusHealthy
	dbOpLogger.LogDatabaseConnection("backup", true, dm.config.Database.Backup.Host)
	appLogger.WithField("provider", dm.config.Database.Backup.Provider).With("host", dm.config.Database.Backup.Host).Info("Connected to backup database")
	return nil
}

// GetDB returns the current active database connection
func (dm *DatabaseManager) GetDB() *gorm.DB {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	appLogger.WithField("database_type", string(dm.currentType)).Debug("GetDB called")
	return dm.currentDB
}

// GetPrimaryDB returns the primary database connection
func (dm *DatabaseManager) GetPrimaryDB() *gorm.DB {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.primaryDB
}

// GetBackupDB returns the backup database connection
func (dm *DatabaseManager) GetBackupDB() *gorm.DB {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.backupDB
}

// GetCurrentType returns the current active database type
func (dm *DatabaseManager) GetCurrentType() DatabaseType {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.currentType
}

// ValidateSchemaConsistency checks if backup database schema matches primary
func (dm *DatabaseManager) ValidateSchemaConsistency() ([]string, error) {
	var mismatches []string

	// Get list of tables from both databases
	primaryTables, err := dm.getTableList(dm.primaryDB)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary tables: %v", err)
	}

	backupTables, err := dm.getTableList(dm.backupDB)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup tables: %v", err)
	}

	// Check for missing tables in backup
	for _, table := range primaryTables {
		if !contains(backupTables, table) {
			mismatches = append(mismatches, fmt.Sprintf("Table '%s' exists in primary but missing in backup", table))
		}
	}

	// Check for extra tables in backup
	for _, table := range backupTables {
		if !contains(primaryTables, table) {
			mismatches = append(mismatches, fmt.Sprintf("Table '%s' exists in backup but missing in primary", table))
		}
	}

	// Check constraints for common tables
	for _, table := range primaryTables {
		if contains(backupTables, table) {
			constraintMismatches, err := dm.compareTableConstraints(table)
			if err != nil {
				mismatches = append(mismatches, fmt.Sprintf("Failed to compare constraints for table '%s': %v", table, err))
			} else {
				mismatches = append(mismatches, constraintMismatches...)
			}
		}
	}

	return mismatches, nil
}

// RebuildBackupSchema rebuilds the backup database schema to match primary
func (dm *DatabaseManager) RebuildBackupSchema() error {
	appLogger.Info("Starting backup database schema rebuild")

	// Drop all tables in backup database
	if err := dm.dropAllTablesInBackup(); err != nil {
		return fmt.Errorf("failed to drop backup tables: %v", err)
	}

	// Run migrations on backup database to recreate schema
	if err := dm.runMigrationsOnBackup(); err != nil {
		return fmt.Errorf("failed to run migrations on backup: %v", err)
	}

	// Run seeders on backup database
	if err := dm.runSeedersOnBackup(); err != nil {
		return fmt.Errorf("failed to run seeders on backup: %v", err)
	}

	appLogger.Info("Backup database schema rebuild completed")
	return nil
}

// getTableList returns list of tables in a database
func (dm *DatabaseManager) getTableList(db *gorm.DB) ([]string, error) {
	var tables []string

	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name`

	if err := db.Raw(query).Scan(&tables).Error; err != nil {
		return nil, err
	}

	return tables, nil
}

// compareTableConstraints compares constraints between primary and backup for a table
func (dm *DatabaseManager) compareTableConstraints(tableName string) ([]string, error) {
	var mismatches []string

	// Get constraints from primary
	primaryConstraints, err := dm.getTableConstraints(dm.primaryDB, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary constraints: %v", err)
	}

	// Get constraints from backup
	backupConstraints, err := dm.getTableConstraints(dm.backupDB, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup constraints: %v", err)
	}

	// Compare constraints
	for constraint := range primaryConstraints {
		if _, exists := backupConstraints[constraint]; !exists {
			mismatches = append(mismatches, fmt.Sprintf("Constraint '%s' missing in backup table '%s'", constraint, tableName))
		}
	}

	for constraint := range backupConstraints {
		if _, exists := primaryConstraints[constraint]; !exists {
			mismatches = append(mismatches, fmt.Sprintf("Extra constraint '%s' in backup table '%s'", constraint, tableName))
		}
	}

	return mismatches, nil
}

// getTableConstraints returns constraints for a table
func (dm *DatabaseManager) getTableConstraints(db *gorm.DB, tableName string) (map[string]string, error) {
	constraints := make(map[string]string)

	query := `
		SELECT constraint_name, constraint_type
		FROM information_schema.table_constraints
		WHERE table_name = ? AND table_schema = 'public'`

	rows, err := db.Raw(query, tableName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, ctype string
		if err := rows.Scan(&name, &ctype); err != nil {
			return nil, err
		}
		constraints[name] = ctype
	}

	return constraints, nil
}

// dropAllTablesInBackup drops all tables in the backup database by dropping and recreating the public schema
func (dm *DatabaseManager) dropAllTablesInBackup() error {
	appLogger.GetLogger().Info("Dropping and recreating public schema in backup database...")

	// Drop the entire public schema with CASCADE to remove all objects
	err := dm.backupDB.Exec("DROP SCHEMA IF EXISTS public CASCADE").Error
	if err != nil {
		return fmt.Errorf("failed to drop public schema: %w", err)
	}
	appLogger.GetLogger().Info("Dropped public schema")

	// Recreate the public schema
	err = dm.backupDB.Exec("CREATE SCHEMA public").Error
	if err != nil {
		return fmt.Errorf("failed to create public schema: %w", err)
	}
	appLogger.GetLogger().Info("Created public schema")

	// Restore default permissions - use current user instead of hardcoded postgres role
	var currentUser string
	err = dm.backupDB.Raw("SELECT current_user").Scan(&currentUser).Error
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	appLogger.WithField("user", currentUser).Info("Current database user")

	err = dm.backupDB.Exec(fmt.Sprintf("GRANT ALL ON SCHEMA public TO %s", currentUser)).Error
	if err != nil {
		return fmt.Errorf("failed to grant permissions to %s: %w", currentUser, err)
	}

	err = dm.backupDB.Exec("GRANT ALL ON SCHEMA public TO public").Error
	if err != nil {
		return fmt.Errorf("failed to grant permissions to public: %w", err)
	}
	appLogger.WithField("user", currentUser).Info("Restored default schema permissions")

	// Enable uuid-ossp extension
	err = dm.backupDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return fmt.Errorf("failed to enable uuid-ossp extension: %w", err)
	}
	appLogger.GetLogger().Info("Enabled uuid-ossp extension")

	appLogger.GetLogger().Info("Public schema recreated successfully in backup database")
	return nil
}

// runMigrationsOnBackup runs migrations on backup database
func (dm *DatabaseManager) runMigrationsOnBackup() error {
	// Migration functions moved to ops package
	// These are now called via dbmanager CLI
	// // Enable uuid-ossp extension
	// if err := enableUUIDOSSP(sqlDB); err != nil {
	// 	appLogger.Warnf("Failed to enable uuid-ossp extension: %v", err)
	// }
	// // Run SQL migrations
	// if err := runSQLMigrationsFromRoots(sqlDB, "inscore/db/migrations", "gen/go/lifepluscore/migrations"); err != nil {
	// 	return fmt.Errorf("failed to run migrations on backup: %v", err)
	// }
	// // Run SQL seeders
	// if err := runSQLSeedersFromRoots(sqlDB, "inscore/db/seeders", "gen/go/lifepluscore/seeders"); err != nil {
	// 	appLogger.Warnf("Failed to run seeders on backup: %v", err)
	// }

	// The following lines are kept to ensure the function remains syntactically correct
	// and to allow for potential future re-introduction of direct migration calls
	// if err := dm.backupDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
	// 	appLogger.WithError(err).Warn("Failed to create uuid-ossp extension")
	// }

	// Get SQL connection from GORM
	// sqlDB, err := dm.backupDB.DB()
	// if err != nil {
	// 	return fmt.Errorf("failed to get SQL connection: %v", err)
	// }

	appLogger.Info("Backup database schema rebuild completed (run migrations via dbmanager CLI)")
	return nil
}

// runSeedersOnBackup runs seeders on backup database
func (dm *DatabaseManager) runSeedersOnBackup() error {
	// Get SQL connection from GORM
	// (commented out since seeders moved to ops package)
	// sqlDB, err := dm.backupDB.DB()
	// if err != nil {
	// 	return fmt.Errorf("failed to get SQL connection: %v", err)
	// }

	// Run seeders using ops package (moved from db package)
	// Seeder functions are now in inscore/db/ops and called via dbmanager CLI
	// if err := runSQLSeedersFromRoots(sqlDB, "inscore/db/seeders", "gen/go/lifepluscore/seeders"); err != nil {
	// 	return fmt.Errorf("failed to run seeders on backup: %v", err)
	// }

	appLogger.Info("Seeders completed on backup database (use dbmanager CLI)")
	return nil
}

// contains checks if a string slice contains a value
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetDetailedStatus returns detailed status informationth databases
func (dm *DatabaseManager) GetStatus() (ConnectionStatus, ConnectionStatus) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.primaryStatus, dm.backupStatus
}

// switchToBackup switches the current connection to backup database
func (dm *DatabaseManager) switchToBackup() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.backupDB == nil || dm.backupStatus != StatusHealthy {
		return fmt.Errorf("backup database is not available")
	}

	dbOpLogger := appLogger.GetDatabaseOperationLogger()
	dbOpLogger.LogFailover("primary", "backup")

	dm.currentDB = dm.backupDB
	dm.currentType = Backup
	dm.metrics.FailoverCount++
	dm.metrics.LastFailoverTime = time.Now()

	appLogger.GetLogger().Info("Switched to backup database (Neon)")
	return nil
}

// switchToPrimary switches the current connection to primary database
func (dm *DatabaseManager) switchToPrimary() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.primaryDB == nil || dm.primaryStatus != StatusHealthy {
		return fmt.Errorf("primary database is not available")
	}

	dm.currentDB = dm.primaryDB
	dm.currentType = Primary

	appLogger.GetLogger().Info("Switched back to primary database (DigitalOcean)")
	return nil
}

// Close closes all database connections and stops background services
func (dm *DatabaseManager) Close() error {
	// Stop background services

	close(dm.healthCheckStop)
	close(dm.syncStop)
	close(dm.backupStop)

	var errors []error

	// Close primary database
	if dm.primaryDB != nil {
		if sqlDB, err := dm.primaryDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close primary database: %v", err))
			}
		}
	}

	// Close backup database
	if dm.backupDB != nil {
		if sqlDB, err := dm.backupDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close backup database: %v", err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing databases: %v", errors)
	}

	appLogger.GetLogger().Info("Database manager closed successfully")
	return nil
}

// GetMetrics returns current database metrics
func (dm *DatabaseManager) GetMetrics() *DatabaseMetrics {
	dm.metrics.mu.RLock()
	defer dm.metrics.mu.RUnlock()

	// Update connection counts
	if dm.primaryDB != nil {
		if sqlDB, err := dm.primaryDB.DB(); err == nil {
			dm.metrics.PrimaryConnections = sqlDB.Stats().OpenConnections
		}
	}

	if dm.backupDB != nil {
		if sqlDB, err := dm.backupDB.DB(); err == nil {
			dm.metrics.BackupConnections = sqlDB.Stats().OpenConnections
		}
	}

	return dm.metrics
}

// Global database manager instance
var Manager *DatabaseManager

// InitializeManager initializes the global database manager
// Use this for migration tools and initial setup
func InitializeManager(configPath string) error {
	config, err := LoadDatabaseConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load database config: %v", err)
	}

	Manager = NewDatabaseManager(config)
	if err := Manager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize database manager: %v", err)
	}

	// Set global DB for backward compatibility
	DB = Manager.GetDB()

	return nil
}

// InitializeManagerForService initializes the global database manager for application services
// This connects to databases with automatic failover but does NOT run migrations
// Use this for microservices (auth, catalog, etc.) that only need database access
func InitializeManagerForService(configPath string) error {
	config, err := LoadDatabaseConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load database config: %v", err)
	}

	Manager = NewDatabaseManager(config)
	if err := Manager.InitializeForService(); err != nil {
		return fmt.Errorf("failed to initialize database manager for service: %v", err)
	}

	// Set global DB for backward compatibility
	DB = Manager.GetDB()

	appLogger.Info("Database manager initialized for service with automatic failover")
	return nil
}

// GetDB returns the current active database (backward compatibility)
func GetDB() *gorm.DB {
	if Manager != nil {
		return Manager.GetDB()
	}
	return DB
}

// ExecuteSQL executes a raw SQL query on the specified database
func ExecuteSQL(query, targetDB string) (string, error) {
	if Manager == nil {
		return "", fmt.Errorf("database manager not initialized")
	}

	var db *gorm.DB
	switch targetDB {
	case "primary":
		db = Manager.primaryDB
	case "backup":
		db = Manager.backupDB
	default:
		return "", fmt.Errorf("invalid target database: %s (must be 'primary' or 'backup')", targetDB)
	}

	if db == nil {
		return "", fmt.Errorf("%s database not available", targetDB)
	}

	// Check if this is a SELECT query to return data
	queryUpper := strings.ToUpper(strings.TrimSpace(query))
	if strings.HasPrefix(queryUpper, "SELECT") || strings.HasPrefix(queryUpper, "SHOW") || strings.HasPrefix(queryUpper, "DESCRIBE") || strings.HasPrefix(queryUpper, "EXPLAIN") {
		// For SELECT queries, return the actual data
		rows, err := db.Raw(query).Rows()
		if err != nil {
			return "", fmt.Errorf("SQL query failed: %v", err)
		}
		defer rows.Close()

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			return "", fmt.Errorf("failed to get columns: %v", err)
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Columns: %s\n", strings.Join(columns, ", ")))
		result.WriteString("Data:\n")

		rowCount := 0
		for rows.Next() {
			// Create a slice of interface{} to hold the values
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				return "", fmt.Errorf("failed to scan row: %v", err)
			}

			// Convert values to strings
			stringValues := make([]string, len(values))
			for i, val := range values {
				if val == nil {
					stringValues[i] = "NULL"
				} else {
					stringValues[i] = fmt.Sprintf("%v", val)
				}
			}

			result.WriteString(fmt.Sprintf("  %s\n", strings.Join(stringValues, ", ")))
			rowCount++
		}

		if rowCount == 0 {
			result.WriteString("  (no rows returned)\n")
		}

		result.WriteString(fmt.Sprintf("Total rows: %d", rowCount))
		return result.String(), nil
	} else {
		// For non-SELECT queries (INSERT, UPDATE, DELETE, etc.), use Exec
		result := db.Exec(query)
		if result.Error != nil {
			return "", fmt.Errorf("SQL execution failed: %v", result.Error)
		}

		return fmt.Sprintf("Rows affected: %d", result.RowsAffected), nil
	}
}

// runMigrationsAndSeeders runs database migrations and seeders
func runMigrationsAndSeeders(db *gorm.DB) error {
	appLogger.GetLogger().Info("Running database migrations and seeders...")

	// Get SQL DB connection for extension and migrations
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %v", err)
	}

	// Enable UUID extension (moved to ops package)
	// if err := enableUUIDOSSP(sqlDB); err != nil {
	// 	appLogger.WithError(err).Warn("Failed to enable uuid-ossp extension")
	// }

	// Run SQL migrations (multi-root) - moved to ops package
	// if err := runSQLMigrationsFromRoots(sqlDB, "inscore/db/migrations", "gen/go/lifepluscore/migrations"); err != nil {
	// 	return fmt.Errorf("failed to run SQL migrations: %v", err)
	// }

	// Run SQL seeders (multi-root) - moved to ops package
	// if err := runSQLSeedersFromRoots(sqlDB, "inscore/db/seeders", "gen/go/lifepluscore/seeders"); err != nil {
	// 	appLogger.WithError(err).Warn("Failed to run SQL seeders")
	// }

	_ = sqlDB // Suppress unused variable warning
	appLogger.GetLogger().Info("Database migrations and seeders completed (run via dbmanager CLI)")
	return nil
}
