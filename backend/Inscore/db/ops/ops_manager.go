package ops

import (
	"fmt"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"gorm.io/gorm"
)

// ManagementAPI provides database management operations
type ManagementAPI struct {
	manager *db.DatabaseManager
}

// NewManagementAPI creates a new management API instance
func NewManagementAPI(manager *db.DatabaseManager) *ManagementAPI {
	return &ManagementAPI{
		manager: manager,
	}
}

// LintMigrationFiles checks for forbidden SQL patterns (Rule 2).
// It does not require a DB connection.
func (api *ManagementAPI) LintMigrationFiles() error {
	// We pass nil for DB because linter only needs filesystem and proto registry
	umm := NewUnifiedMigrationManager(nil)
	return umm.LintMigrationFiles()
}

// MigrateOnly runs migrations and seeders ONLY on the specified database.
// Uses the new UnifiedMigrationManager for proto -> SQL -> seeders flow.
func (api *ManagementAPI) MigrateOnly(targetDB db.DatabaseType, prune bool, strict bool) error {
	appLogger.Infof("Running unified migrations on %s database", targetDB)

	var targetGormDB *gorm.DB
	var targetName string
	if targetDB == db.Primary {
		targetGormDB = api.manager.GetPrimaryDB()
		targetName = "primary"
	} else {
		targetGormDB = api.manager.GetBackupDB()
		targetName = "backup"
	}

	if targetGormDB == nil {
		return fmt.Errorf("target database is not available")
	}

	// Get SQL DB for migrations
	sqlDB, err := targetGormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %v", err)
	}

	// Create unified migration manager (auto-discovers schemas from proto)
	umm := NewUnifiedMigrationManager(sqlDB)
	umm.SetPruneColumns(prune)
	umm.SetStrictMode(strict)

	// Initialize migration system (creates extensions, schemas, metadata table)
	if err := umm.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize migration system: %v", err)
	}

	// Run complete migration flow: Proto -> SQL Migrations -> Seeders
	if err := umm.RunAll(); err != nil {
		return fmt.Errorf("migration flow failed: %v", err)
	}

	appLogger.Infof("✅ %s database migrations completed successfully", targetName)
	return nil
}

// DatabaseOperations provides high-level database operations
type DatabaseOperations struct {
	// Connection management
	SwitchToPrimary func() error
	SwitchToBackup  func() error
	ForceFailover   func() error

	// Synchronization
	SyncNow   func() error
	SyncTable func(tableName string) error

	// Backup and restore
	BackupNow         func() error
	RestoreFromBackup func(backupPath string, targetDB db.DatabaseType) error
	ListBackups       func() ([]BackupInfo, error)

	// Status and monitoring
	GetStatus  func() map[string]interface{}
	GetMetrics func() *db.DatabaseMetrics
	IsHealthy  func() bool
}

// GetOperations returns database operations interface
func (api *ManagementAPI) GetOperations() *DatabaseOperations {
	return &DatabaseOperations{
		SwitchToPrimary:   api.manager.ForceSwitchBack,
		SwitchToBackup:    api.manager.ForceFailover,
		ForceFailover:     api.manager.ForceFailover,
		SyncNow:           api.SyncNow,
		SyncTable:         api.SyncTable,
		BackupNow:         api.BackupNow,
		RestoreFromBackup: api.RestoreFromBackup,
		ListBackups:       api.ListBackups,
		GetStatus:         api.manager.GetDetailedStatus,
		GetMetrics:        api.manager.GetMetrics,
		IsHealthy:         api.manager.IsHealthy,
	}
}

// CopyDatabase copies data from source to target database
func (api *ManagementAPI) CopyDatabase(sourceDB, targetDB db.DatabaseType) error {
	appLogger.Infof("Starting database copy from %s to %s", sourceDB, targetDB)

	var sourceGormDB, targetGormDB *gorm.DB

	if sourceDB == db.Primary {
		sourceGormDB = api.manager.GetPrimaryDB()
	} else {
		sourceGormDB = api.manager.GetBackupDB()
	}

	if targetDB == db.Primary {
		targetGormDB = api.manager.GetPrimaryDB()
	} else {
		targetGormDB = api.manager.GetBackupDB()
	}

	if sourceGormDB == nil || targetGormDB == nil {
		return fmt.Errorf("source or target database is not available")
	}

	// Get list of tables to copy
	tables, err := api.getTableList(sourceGormDB)
	if err != nil {
		return fmt.Errorf("failed to get table list: %v", err)
	}

	for _, tableName := range tables {
		if api.isTableExcluded(tableName) {
			appLogger.Infof("Skipping excluded table: %s", tableName)
			continue
		}

		if err := api.copyTable(sourceGormDB, targetGormDB, tableName); err != nil {
			appLogger.Warnf("Failed to copy table %s: %v", tableName, err)
			continue
		}

		appLogger.Infof("Successfully copied table: %s", tableName)
	}

	appLogger.Infof("Database copy completed from %s to %s", sourceDB, targetDB)
	return nil
}

// copyTable copies a single table from source to target
func (api *ManagementAPI) copyTable(sourceDB, targetDB *gorm.DB, tableName string) error {
	// Get table structure
	columns, err := api.getTableColumnsFromDB(sourceDB, tableName)
	if err != nil {
		return fmt.Errorf("failed to get table columns: %v", err)
	}

	// Truncate target table
	if err := targetDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)).Error; err != nil {
		appLogger.Warnf("Failed to truncate table %s: %v", tableName, err)
	}

	// Copy data in batches
	batchSize := api.manager.GetConfig().Database.Sync.BatchSize
	offset := 0

	for {
		// Get batch from source
		query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at LIMIT %d OFFSET %d",
			tableName, batchSize, offset)

		rows, err := sourceDB.Raw(query).Rows()
		if err != nil {
			return fmt.Errorf("failed to query source table: %v", err)
		}

		// Check if we have any rows
		hasRows := false
		var insertData []map[string]interface{}

		for rows.Next() {
			hasRows = true

			// Get column values
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan row: %v", err)
			}

			// Create map for this row
			rowData := make(map[string]interface{})
			for i, col := range columns {
				rowData[col] = values[i]
			}
			insertData = append(insertData, rowData)
		}
		rows.Close()

		if !hasRows {
			break // No more data
		}

		// Insert batch into target
		if len(insertData) > 0 {
			if err := targetDB.Table(tableName).Create(&insertData).Error; err != nil {
				return fmt.Errorf("failed to insert batch: %v", err)
			}
		}

		offset += batchSize
	}

	return nil
}

// getTableList returns list of tables in the database
func (api *ManagementAPI) getTableList(db *gorm.DB) ([]string, error) {
	var tables []string

	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name`

	rows, err := db.Raw(query).Rows()
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

// getTableColumnsFromDB gets column names for a table from database
func (api *ManagementAPI) getTableColumnsFromDB(db *gorm.DB, tableName string) ([]string, error) {
	var columns []string

	query := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = ? AND table_schema = 'public'
		ORDER BY ordinal_position`

	rows, err := db.Raw(query, tableName).Rows()
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

// MigrateToDatabase runs migrations on a specific database and optionally syncs to the other
// This ensures both primary and backup databases stay in sync
func (api *ManagementAPI) MigrateToDatabase(targetDB db.DatabaseType, prune bool, strict bool) error {
	appLogger.Infof("Running migrations on %s database", targetDB)

	// Run migrations on target database
	if err := api.MigrateOnly(targetDB, prune, strict); err != nil {
		return fmt.Errorf("failed to migrate %s database: %w", targetDB, err)
	}

	// Determine the other database
	var otherDB db.DatabaseType
	var otherGormDB *gorm.DB
	var otherDBName string

	if targetDB == db.Primary {
		otherDB = db.Backup
		otherGormDB = api.manager.GetBackupDB()
		otherDBName = "backup"
	} else {
		otherDB = db.Primary
		otherGormDB = api.manager.GetPrimaryDB()
		otherDBName = "primary"
	}

	// If the other database is available, migrate it too for consistency
	if otherGormDB != nil {
		appLogger.Infof("Also running migrations on %s database for consistency...", otherDBName)

		if err := api.MigrateOnly(otherDB, prune, strict); err != nil {
			// Don't fail the entire operation if secondary database fails
			appLogger.Warnf("Failed to migrate %s database: %v", otherDBName, err)
			appLogger.Warn("Primary migration succeeded, but secondary database may be out of sync")
		} else {
			appLogger.Infof("✅ %s database migrations completed", otherDBName)
		}
	} else {
		appLogger.Infof("Skipping %s database (not available or not configured)", otherDBName)
	}

	appLogger.Infof("✅ Migration operation completed for %s database", targetDB)
	return nil
}

// CompareSchemas compares schemas between primary and backup databases
func (api *ManagementAPI) CompareSchemas() (map[string]interface{}, error) {
	if api.manager.GetPrimaryDB() == nil || api.manager.GetBackupDB() == nil {
		return nil, fmt.Errorf("both databases must be available for schema comparison")
	}

	primaryTables, err := api.getTableList(api.manager.GetPrimaryDB())
	if err != nil {
		return nil, fmt.Errorf("failed to get primary tables: %v", err)
	}

	backupTables, err := api.getTableList(api.manager.GetBackupDB())
	if err != nil {
		return nil, fmt.Errorf("failed to get backup tables: %v", err)
	}

	// Find differences
	primaryOnly := []string{}
	backupOnly := []string{}
	common := []string{}

	primaryMap := make(map[string]bool)
	for _, table := range primaryTables {
		primaryMap[table] = true
	}

	backupMap := make(map[string]bool)
	for _, table := range backupTables {
		backupMap[table] = true
	}

	// Find tables only in primary
	for _, table := range primaryTables {
		if !backupMap[table] {
			primaryOnly = append(primaryOnly, table)
		} else {
			common = append(common, table)
		}
	}

	// Find tables only in backup
	for _, table := range backupTables {
		if !primaryMap[table] {
			backupOnly = append(backupOnly, table)
		}
	}

	return map[string]interface{}{
		"primary_only_tables": primaryOnly,
		"backup_only_tables":  backupOnly,
		"common_tables":       common,
		"primary_table_count": len(primaryTables),
		"backup_table_count":  len(backupTables),
		"schemas_match":       len(primaryOnly) == 0 && len(backupOnly) == 0,
	}, nil
}

// GetDatabaseSizes returns size information for both databases
func (api *ManagementAPI) GetDatabaseSizes() (map[string]interface{}, error) {
	sizes := make(map[string]interface{})

	if api.manager.GetPrimaryDB() != nil {
		primarySize, err := api.getDatabaseSize(api.manager.GetPrimaryDB())
		if err != nil {
			appLogger.Warnf("Failed to get primary database size: %v", err)
		} else {
			sizes["primary_size"] = primarySize
		}
	}

	if api.manager.GetBackupDB() != nil {
		backupSize, err := api.getDatabaseSize(api.manager.GetBackupDB())
		if err != nil {
			appLogger.Warnf("Failed to get backup database size: %v", err)
		} else {
			sizes["backup_size"] = backupSize
		}
	}

	return sizes, nil
}

// getDatabaseSize gets the size of a database
func (api *ManagementAPI) getDatabaseSize(db *gorm.DB) (map[string]interface{}, error) {
	var totalSize int64
	err := db.Raw("SELECT pg_database_size(current_database())").Scan(&totalSize).Error
	if err != nil {
		return nil, err
	}

	// Get table sizes
	query := `
		SELECT 
			schemaname,
			tablename,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
			pg_total_relation_size(schemaname||'.'||tablename) as size_bytes
		FROM pg_tables 
		WHERE schemaname = 'public'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC`

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []map[string]interface{}
	for rows.Next() {
		var schema, table, sizeStr string
		var sizeBytes int64

		if err := rows.Scan(&schema, &table, &sizeStr, &sizeBytes); err != nil {
			continue
		}

		tables = append(tables, map[string]interface{}{
			"schema":     schema,
			"table":      table,
			"size":       sizeStr,
			"size_bytes": sizeBytes,
		})
	}

	return map[string]interface{}{
		"total_size_bytes": totalSize,
		"total_size":       formatBytes(totalSize),
		"tables":           tables,
	}, nil
}

// formatBytes formats bytes into human readable format
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
