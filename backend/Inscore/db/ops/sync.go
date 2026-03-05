package ops

import (
	"fmt"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// StartSync starts the database synchronization routine
func (api *ManagementAPI) StartSync() {
	config := api.manager.GetConfig()
	ticker := time.NewTicker(config.Database.Sync.Interval)
	defer ticker.Stop()

	appLogger.Infof("Starting database sync (interval=%s)", config.Database.Sync.Interval)

	for {
		select {
		case <-ticker.C:
			if err := api.performSync(); err != nil {
				appLogger.Errorf("Sync failed: %v", err)
			}
			// case <-dm.syncStop: // Need a way to stop this
			// 	appLogger.Info("Database sync stopped")
			// 	return
		}
	}
}

// performSync synchronizes data between primary and backup databases
func (api *ManagementAPI) performSync() error {
	if api.manager.GetPrimaryDB() == nil || api.manager.GetBackupDB() == nil {
		return fmt.Errorf("both primary and backup databases must be available for sync")
	}

	primaryStatus, backupStatus := api.manager.GetStatus()
	if primaryStatus != db.StatusHealthy || backupStatus != db.StatusHealthy {
		return fmt.Errorf("both databases must be healthy for sync")
	}

	appLogger.Info("Starting database synchronization...")
	start := time.Now()

	// Sync tables in dependency order to avoid foreign key violations
	config := api.manager.GetConfig()
	syncOrder := api.getSyncOrder(config.Database.Sync.TablesToSync)

	for _, tableName := range syncOrder {
		if api.isTableExcluded(tableName) {
			continue
		}

		if err := api.performSyncTable(tableName); err != nil {
			appLogger.Errorf("Failed to sync table %s: %v", tableName, err)
			continue
		}
	}

	// dm.metrics.mu.Lock()
	// dm.metrics.LastSyncTime = time.Now()
	// dm.metrics.mu.Unlock()

	duration := time.Since(start)
	appLogger.Infof("Database sync completed (duration=%s)", duration)
	return nil
}

// performSyncTable performs the actual synchronization of a table
func (api *ManagementAPI) performSyncTable(tableName string) error {
	appLogger.Infof("Starting table sync: %s", tableName)

	// Get table columns
	columns, err := api.getTableColumns(tableName)
	if err != nil {
		return fmt.Errorf("failed to get columns for table %s: %v", tableName, err)
	}

	// Check if table has a primary key column that we can use for sync
	primaryKeyColumn := api.getPrimaryKeyColumn(tableName, columns)
	if primaryKeyColumn == "" {
		appLogger.Infof("Skipping table (no identifiable primary key): %s", tableName)
		return nil
	}

	// Get total row count
	var totalRows int64
	if err := api.manager.GetPrimaryDB().Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&totalRows).Error; err != nil {
		return fmt.Errorf("failed to get row count: %v", err)
	}

	if totalRows == 0 {
		appLogger.Infof("Table is empty, skipping sync: %s", tableName)
		return nil
	}

	// Sync in batches (use configured batch size when available)
	config := api.manager.GetConfig()
	batchSize := config.Database.Sync.BatchSize
	if batchSize <= 0 {
		batchSize = 1000
	}
	rowCount := int64(0)
	for offset := 0; offset < int(totalRows); offset += batchSize {
		batchStart := time.Now()
		appLogger.Debug(fmt.Sprintf("Syncing batch table=%s offset=%d limit=%d total_rows=%d", tableName, offset, batchSize, totalRows))

		if err := api.syncTableBatch(tableName, columns, offset, batchSize); err != nil {
			appLogger.Errorf("Failed to sync batch for table %s (offset=%d, limit=%d): %v", tableName, offset, batchSize, err)
			// Continue with next batch instead of failing completely
			continue
		}
		appLogger.Debug(fmt.Sprintf("Batch synced table=%s offset=%d limit=%d duration=%s", tableName, offset, batchSize, time.Since(batchStart)))
		// Accumulate by min(batchSize, remaining)
		remaining := int64(totalRows) - rowCount
		if remaining >= int64(batchSize) {
			rowCount += int64(batchSize)
		} else {
			rowCount += remaining
		}
	}

	appLogger.Infof("Table sync completed: %s (rows=%d)", tableName, rowCount)
	return nil
}

// syncTableBatch synchronizes a batch of rows from a table
func (api *ManagementAPI) syncTableBatch(tableName string, columns []string, offset, limit int) error {
	// Build column list for SELECT
	columnList := strings.Join(columns, ", ")

	// Determine ordering column
	// Special-case: process_events has index on timestamp DESC
	orderBy := ""
	if tableName == "process_events" {
		// Prefer indexed timestamp for efficient keyset-like pagination with OFFSET
		for _, c := range columns {
			if c == "timestamp" {
				orderBy = "timestamp"
				break
			}
		}
	}
	// Otherwise prefer created_at, then id
	for _, c := range columns {
		if orderBy == "" && c == "created_at" {
			orderBy = "created_at"
			break
		}
	}
	if orderBy == "" {
		for _, c := range columns {
			if c == "id" {
				orderBy = "id"
				break
			}
		}
	}

	// Get data from primary database
	var query string
	if orderBy != "" {
		query = fmt.Sprintf("SELECT %s FROM %s ORDER BY %s LIMIT %d OFFSET %d", columnList, tableName, orderBy, limit, offset)
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s LIMIT %d OFFSET %d", columnList, tableName, limit, offset)
	}

	rows, err := api.manager.GetPrimaryDB().Raw(query).Rows()
	if err != nil {
		return fmt.Errorf("failed to query primary database: %v", err)
	}
	defer rows.Close()

	// Process rows individually to handle foreign key violations gracefully
	successCount := 0
	errorCount := 0

	for rows.Next() {
		// Scan row values
		rowValues := make([]interface{}, len(columns))
		rowPointers := make([]interface{}, len(columns))
		for i := range rowValues {
			rowPointers[i] = &rowValues[i]
		}

		if err := rows.Scan(rowPointers...); err != nil {
			appLogger.Errorf("Failed to scan row for table %s: %v", tableName, err)
			errorCount++
			continue
		}

		// Build single row upsert
		placeholders := make([]string, len(columns))
		for i := range placeholders {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}

		// Get the appropriate ON CONFLICT clause for this table
		conflictClause := api.getConflictClause(tableName, columns)

		// Execute single row upsert
		upsertQuery := fmt.Sprintf(`
			INSERT INTO %s (%s) VALUES (%s) 
			%s`,
			tableName,
			columnList,
			strings.Join(placeholders, ", "),
			conflictClause,
		)

		if err := api.manager.GetBackupDB().Exec(upsertQuery, rowValues...).Error; err != nil {
			// Check if this is a foreign key violation
			if api.isForeignKeyViolation(err) {
				// Try to sync the missing referenced record first
				if api.syncMissingReferences(tableName, rowValues, columns) {
					// Retry the insert after syncing references
					if retryErr := api.manager.GetBackupDB().Exec(upsertQuery, rowValues...).Error; retryErr == nil {
						successCount++
						continue
					}
				}
			}

			// Log the error but continue with next row
			appLogger.Warnf("Failed to sync individual row, skipping (table=%s, id=%v): %v", tableName, rowValues[0], err)
			errorCount++
			continue
		}
		successCount++
	}

	if errorCount > 0 {
		appLogger.Warnf("Batch sync completed with errors (table=%s, success=%d, errors=%d)", tableName, successCount, errorCount)
	}

	return nil
}

// getConflictClause returns appropriate ON CONFLICT clause for upserts
func (api *ManagementAPI) getConflictClause(tableName string, columns []string) string {
	switch tableName {
	case "process_events":
		// Event store is append-only. Use natural unique constraint to ensure idempotency.
		// No updated_at column in this table; avoid timestamp-based update predicate.
		return "ON CONFLICT (aggregate_id, event_version) DO NOTHING"
	case "approval_templates":
		// Handle unique constraint on (name, entity_type)
		return "ON CONFLICT (name, entity_type) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= " + tableName + ".updated_at"
	case "approval_rules":
		// Handle unique constraint on (name, entity_type)
		return "ON CONFLICT (name, entity_type) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= " + tableName + ".updated_at"
	case "users":
		// Handle email unique constraint - sync newer records based on updated_at
		return "ON CONFLICT (email) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= users.updated_at"
	case "teams":
		// Handle unique constraint on name - sync newer records
		return "ON CONFLICT (name) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= teams.updated_at"
	case "user_teams":
		// Handle composite key (user_id, team_id) - allow updates for newer records
		return "ON CONFLICT (user_id, team_id) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= user_teams.updated_at"
	case "companies", "clients", "suppliers":
		// Business entities - sync newer records based on updated_at
		return "ON CONFLICT (id) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= " + tableName + ".updated_at"
	case "enquiries", "styles", "draft_costs":
		// Core business data - always sync newer records
		return "ON CONFLICT (id) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= " + tableName + ".updated_at"
	case "fabric_constructions", "enhanced_fabric_constructions", "fabric_weave_knits", "fabric_specifications":
		// Reference data - sync newer records
		return "ON CONFLICT (id) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= " + tableName + ".updated_at"
	case "files", "storage_schema.files":
		// Storage metadata table uses file_id as PK.
		return "ON CONFLICT (file_id) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= " + tableName + ".updated_at"
	default:
		// Default approach: update all columns on ID conflict with timestamp check
		return "ON CONFLICT (id) DO UPDATE SET " + api.buildUpdateClause(columns) + " WHERE EXCLUDED.updated_at >= " + tableName + ".updated_at"
	}
}

// getSyncOrder returns tables in dependency order to avoid foreign key violations
func (api *ManagementAPI) getSyncOrder(tables []string) []string {
	// Define dependency order - parent tables first, then child tables
	dependencyOrder := []string{
		// Foundation tables (no dependencies)
		"addresses", "contacts", "companies", "clients",

		// User management tables - users and teams first, then dependent tables
		"users", "teams", "user_roles", "user_teams",

		// Enlistments must come before suppliers (suppliers references enlistments)
		"enlistments",

		// Suppliers (depends on enlistments)
		"suppliers",

		// Fabric foundation tables
		"fabric_constructions", "enhanced_fabric_constructions", "fabric_weave_knits",

		// Fabric specification tables (depends on constructions and weave_knits)
		"fabric_specifications",

		// Enquiry tables must come before styles (styles references enquiries)
		"enquiries",

		// Style tables (depends on fabric_specifications and enquiries)
		"styles", "style_colors", "style_images", "style_sizes", "style_measurements", "style_comments",

		// Draft costs (depends on styles and enquiries)
		"draft_costs",

		// Other dependent tables
		"company_versions", "supplier_products", "product_types",
		"fabrics", "trims", "packaging", "processings", "fibers", "wash",
		"negotiations", "style_negotiations", "enquiry_conversations",
		"order_execution", "order_apparels", "samples", "sample_images",

		// Approval system tables
		"approval_templates", "approval_rules", "approval_workflows",

		// System tables
		"document_versions", "attachments", "audit_logs", "approvals", "logos",

		// Workflow / Process engine tables (parents first)
		// Processes are the root for tasks and process events
		"processes",
		// Instances/workflow execution state (may reference processes)
		"workflow_instances",
		// Tasks depend on processes
		"tasks",
		// Assignments depend on tasks and users
		"task_assignments",
		// Milestones depend on processes; templates are independent
		"process_milestones", "milestone_templates",
		// Events typically reference processes
		"process_events",
		// Workflow errors can reference workflow_instances or processes
		"workflow_errors",
		// Notification system (subscriptions may reference users)
		"notification_templates", "notification_subscriptions", "notifications",
	}

	// Filter to only include tables that are in the sync list
	var orderedTables []string
	tableSet := make(map[string]bool)
	for _, table := range tables {
		tableSet[table] = true
	}

	// Add tables in dependency order
	for _, table := range dependencyOrder {
		if tableSet[table] {
			orderedTables = append(orderedTables, table)
			delete(tableSet, table)
		}
	}

	// Add any remaining tables not in the dependency order
	for table := range tableSet {
		orderedTables = append(orderedTables, table)
	}

	return orderedTables
}

// getPrimaryKeyColumn returns the primary key column name for sync operations
func (api *ManagementAPI) getPrimaryKeyColumn(tableName string, columns []string) string {
	// Check for standard 'id' column first
	for _, col := range columns {
		if col == "id" {
			return "id"
		}
	}

	// Handle special cases with non-standard primary keys
	switch tableName {
	case "users":
		for _, col := range columns {
			if col == "user_id" {
				return "user_id"
			}
		}
	case "teams":
		for _, col := range columns {
			if col == "team_id" {
				return "team_id"
			}
		}
	case "user_teams":
		// Composite key table - use user_id for ordering, but handle specially
		for _, col := range columns {
			if col == "user_id" {
				return "user_id"
			}
		}
	case "files", "storage_schema.files":
		for _, col := range columns {
			if col == "file_id" {
				return "file_id"
			}
		}
	}

	// Fallback for common *_id primary key naming.
	for _, col := range columns {
		if strings.HasSuffix(col, "_id") {
			return col
		}
	}

	return ""
}

// isForeignKeyViolation checks if the error is a foreign key constraint violation
func (api *ManagementAPI) isForeignKeyViolation(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "violates foreign key constraint") ||
		strings.Contains(errStr, "SQLSTATE 23503")
}

// recordExists checks if a record exists in the backup database
func (api *ManagementAPI) recordExists(tableName, id string) bool {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = $1", tableName)
	api.manager.GetBackupDB().Raw(query, id).Scan(&count)
	return count > 0
}

// syncSingleRecord syncs a single record by ID from primary to backup
func (api *ManagementAPI) syncSingleRecord(tableName, id string) bool {
	// Get table columns
	columns, err := api.getTableColumns(tableName)
	if err != nil {
		appLogger.Errorf("Failed to get columns for single record sync (table=%s): %v", tableName, err)
		return false
	}

	// Get the record from primary database
	columnList := strings.Join(columns, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", columnList, tableName)

	rows, err := api.manager.GetPrimaryDB().Raw(query, id).Rows()
	if err != nil {
		appLogger.Errorf("Failed to query single record from primary (table=%s, id=%v): %v", tableName, id, err)
		return false
	}
	defer rows.Close()

	if !rows.Next() {
		appLogger.Warnf("Record not found in primary database (table=%s, id=%v)", tableName, id)
		return false
	}

	// Scan the record
	rowValues := make([]interface{}, len(columns))
	rowPointers := make([]interface{}, len(columns))
	for i := range rowValues {
		rowPointers[i] = &rowValues[i]
	}

	if err := rows.Scan(rowPointers...); err != nil {
		appLogger.Errorf("Failed to scan single record (table=%s, id=%v): %v", tableName, id, err)
		return false
	}

	// Build and execute insert
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	conflictClause := api.getConflictClause(tableName, columns)
	upsertQuery := fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES (%s) 
		%s`,
		tableName,
		columnList,
		strings.Join(placeholders, ", "),
		conflictClause,
	)

	if err := api.manager.GetBackupDB().Exec(upsertQuery, rowValues...).Error; err != nil {
		appLogger.Errorf("Failed to insert single record (table=%s, id=%v): %v", tableName, id, err)
		return false
	}

	appLogger.Infof("Successfully synced missing reference record (table=%s, id=%v)", tableName, id)

	return true
}

// getTableColumns retrieves column names for a table
func (api *ManagementAPI) getTableColumns(tableName string) ([]string, error) {
	var columns []string

	schemaName := "public"
	rawTable := tableName
	if strings.Contains(tableName, ".") {
		parts := strings.SplitN(tableName, ".", 2)
		schemaName = strings.TrimSpace(parts[0])
		rawTable = strings.TrimSpace(parts[1])
	}

	query := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = ? AND table_schema = ?
		ORDER BY ordinal_position`

	rows, err := api.manager.GetPrimaryDB().Raw(query, rawTable, schemaName).Rows()
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

// buildUpdateClause builds the UPDATE clause for upsert operations
func (api *ManagementAPI) buildUpdateClause(columns []string) string {
	var updates []string
	for _, col := range columns {
		if col != "id" && col != "created_at" { // Don't update ID and created_at
			updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
		}
	}
	return strings.Join(updates, ", ")
}

// asUUIDString safely converts scanned DB value (string or []byte) to string UUID
func asUUIDString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// isTableExcluded checks if a table should be excluded from sync
func (api *ManagementAPI) isTableExcluded(tableName string) bool {
	config := api.manager.GetConfig()
	for _, excluded := range config.Database.Sync.ExcludeTables {
		if excluded == tableName {
			return true
		}
	}
	return false
}

// SyncNow manually triggers an immediate synchronization
func (api *ManagementAPI) SyncNow() error {
	appLogger.Info("Manual sync triggered")
	return api.performSync()
}

// SyncTable manually synchronizes a specific table
func (api *ManagementAPI) SyncTable(tableName string) error {
	appLogger.Infof("Manual sync triggered for table: %s", tableName)

	if api.manager.GetPrimaryDB() == nil || api.manager.GetBackupDB() == nil {
		return fmt.Errorf("both primary and backup databases must be available")
	}

	return api.performSyncTable(tableName)
}

// GetSyncStatus returns the current synchronization status
func (api *ManagementAPI) GetSyncStatus() map[string]interface{} {
	// dm.metrics.mu.RLock()
	// defer dm.metrics.mu.RUnlock()

	config := api.manager.GetConfig()

	return map[string]interface{}{
		"sync_enabled":  config.Database.Sync.Enabled,
		"sync_interval": config.Database.Sync.Interval.String(),
		// "last_sync_time":  dm.metrics.LastSyncTime,
		"tables_to_sync":  config.Database.Sync.TablesToSync,
		"excluded_tables": config.Database.Sync.ExcludeTables,
		"batch_size":      config.Database.Sync.BatchSize,
	}
}

// syncMissingReferences tries to sync referenced parent records when an FK violation occurs.
// For now, this is a stub that returns false to skip retry; refine as needed.
func (api *ManagementAPI) syncMissingReferences(tableName string, rowValues []interface{}, columns []string) bool {
	return false
}
