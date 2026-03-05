package ops

import (
	"fmt"
	"strings"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"gorm.io/gorm"
)

// TableStatus holds health info for a table across primary and backup
type TableStatus struct {
	Name         string `json:"name"`
	PrimaryCount int64  `json:"primary_count"`
	BackupCount  int64  `json:"backup_count"`
	InSync       bool   `json:"in_sync"`
	Difference   int64  `json:"difference"`
}

func splitQualifiedTable(table string) (string, string) {
	if strings.Contains(table, ".") {
		parts := strings.SplitN(table, ".", 2)
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return "public", strings.TrimSpace(table)
}

// getEffectiveSyncTables returns the list of tables to sync.
// If no tables are configured, it auto-discovers schema and excludes logs/sessions by default.
func (api *ManagementAPI) getEffectiveSyncTables() []string {
	config := api.manager.GetConfig()
	cfg := config.Database.Sync.TablesToSync
	if len(cfg) > 0 {
		return cfg
	}
	// Auto-discover when not configured
	discovered, err := api.DiscoverSchema()
	if err != nil {
		appLogger.Warnf("Schema discovery failed; using empty table list: %v", err)
		return []string{}
	}
	// Default excludes: logs and sessions tables
	var filtered []string
	for _, t := range discovered {
		lt := strings.ToLower(t)
		if strings.Contains(lt, "log") || strings.Contains(lt, "session") {
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered
}

// PruneExtras removes rows from backup that do not exist in primary for the given tables.
// It processes tables in reverse dependency order to satisfy FK constraints.
// Returns a map of table -> deleted row count.
func (api *ManagementAPI) PruneExtras(tables []string) (map[string]int64, error) {
	if api.manager.GetPrimaryDB() == nil || api.manager.GetBackupDB() == nil {
		return nil, fmt.Errorf("both primary and backup databases must be available for prune")
	}

	// Determine processing order: reverse of dependency order
	ordered := api.getSyncOrder(tables)
	for i, j := 0, len(ordered)-1; i < j; i, j = i+1, j-1 {
		ordered[i], ordered[j] = ordered[j], ordered[i]
	}

	results := make(map[string]int64)

	for _, table := range ordered {
		// Skip non-existent tables gracefully
		if !api.tableExists(table) {
			continue
		}

		// Resolve primary key dynamically to support non-standard PK columns (e.g., file_id).
		pkColumn, err := api.getPrimaryKeyColumnFromDB(table)
		if err != nil || pkColumn == "" {
			appLogger.Warnf("Skipping prune for table with unknown PK (table=%s): %v", table, err)
			continue
		}

		// Fetch PK sets from primary and backup.
		primaryIDs, err := getAllIDs(api.manager.GetPrimaryDB(), table, pkColumn)
		if err != nil {
			appLogger.Warnf("Failed to fetch primary IDs for prune (table=%s): %v", table, err)
			continue
		}
		backupIDs, err := getAllIDs(api.manager.GetBackupDB(), table, pkColumn)
		if err != nil {
			appLogger.Warnf("Failed to fetch backup IDs for prune (table=%s): %v", table, err)
			continue
		}

		if len(backupIDs) == 0 {
			results[table] = 0
			continue
		}

		primarySet := make(map[string]struct{}, len(primaryIDs))
		for _, id := range primaryIDs {
			primarySet[id] = struct{}{}
		}

		var extras []string
		for _, id := range backupIDs {
			if _, ok := primarySet[id]; !ok {
				extras = append(extras, id)
			}
		}

		if len(extras) == 0 {
			results[table] = 0
			continue
		}

		// Delete extras in batches
		const batchSize = 500
		var totalDeleted int64
		for start := 0; start < len(extras); start += batchSize {
			end := start + batchSize
			if end > len(extras) {
				end = len(extras)
			}
			batch := extras[start:end]
			tx := api.manager.GetBackupDB().Table(table).Where(pkColumn+" IN ?", batch).Delete(nil)
			if tx.Error != nil {
				appLogger.Warnf("Failed deleting extras (table=%s): %v", table, tx.Error)
				continue
			}
			totalDeleted += tx.RowsAffected
		}

		results[table] = totalDeleted
		appLogger.Infof("Pruned extras from table %s (deleted=%d)", table, totalDeleted)
	}

	return results, nil
}

// getAllIDs fetches all PK values from a table.
func getAllIDs(dbConn *gorm.DB, table string, pkColumn string) ([]string, error) {
	rows, err := dbConn.Raw(fmt.Sprintf("SELECT %s FROM %s", pkColumn, table)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var v interface{}
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		ids = append(ids, asUUIDString(v))
	}
	return ids, nil
}

// DependencyAwareSync performs ordered sync using the manager's internal logic
// Consolidated here to avoid spawning external commands.
func (api *ManagementAPI) DependencyAwareSync(commit bool) error {
	if api.manager.GetPrimaryDB() == nil || api.manager.GetBackupDB() == nil {
		return fmt.Errorf("both primary and backup databases must be available for sync")
	}

	// Build effective table list, ordered by dependency
	tables := api.getSyncOrder(api.getEffectiveSyncTables())
	for _, t := range tables {
		if api.isTableExcluded(t) {
			continue
		}
		// performSyncTable always writes from primary->backup with FK repair
		if err := api.performSyncTable(t); err != nil {
			// Continue to next table; report at the end
			appLogger.Errorf("Dependency-aware sync table %s failed: %v", t, err)
		}
	}

	// When committing, enforce authoritative deletions (backup-only extras) to eliminate drift
	if commit {
		statuses, err := api.SyncHealthCheck()
		if err != nil {
			return err
		}
		var drifted []string
		for _, s := range statuses {
			if !s.InSync && s.Difference < 0 { // backup has more rows than primary
				drifted = append(drifted, s.Name)
			}
		}
		if len(drifted) > 0 {
			// Prune in reverse dependency order to satisfy FKs
			deleted, err := api.PruneExtras(drifted)
			if err != nil {
				appLogger.Warnf("PruneExtras encountered errors: %v", err)
			}
			for table, n := range deleted {
				appLogger.Infof("Authoritative prune applied (table=%s, deleted=%d)", table, n)
			}
		}
	}
	return nil
}

// SyncHealthCheck returns per-table counts and whether they match
func (api *ManagementAPI) SyncHealthCheck() ([]TableStatus, error) {
	if api.manager.GetPrimaryDB() == nil || api.manager.GetBackupDB() == nil {
		return nil, fmt.Errorf("both primary and backup databases must be available for health check")
	}

	var results []TableStatus
	tables := api.getSyncOrder(api.getEffectiveSyncTables())
	for _, t := range tables {
		if api.isTableExcluded(t) {
			continue
		}

		var pc, bc int64
		if err := api.manager.GetPrimaryDB().Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", t)).Scan(&pc).Error; err != nil {
			return nil, fmt.Errorf("count primary %s: %w", t, err)
		}
		if err := api.manager.GetBackupDB().Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", t)).Scan(&bc).Error; err != nil {
			return nil, fmt.Errorf("count backup %s: %w", t, err)
		}

		results = append(results, TableStatus{
			Name:         t,
			PrimaryCount: pc,
			BackupCount:  bc,
			InSync:       pc == bc,
			Difference:   pc - bc,
		})
	}
	return results, nil
}

// DiscoverSchema lists base tables in non-system schemas present in primary DB.
// For public schema, table names are returned as plain table names for backward compatibility.
// For non-public schemas, names are schema-qualified (schema.table).
func (api *ManagementAPI) DiscoverSchema() ([]string, error) {
	if api.manager.GetPrimaryDB() == nil {
		return nil, fmt.Errorf("primary database must be available")
	}
	rows, err := api.manager.GetPrimaryDB().Raw(`
		SELECT table_schema, table_name 
		FROM information_schema.tables 
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema') AND table_type = 'BASE TABLE'
		ORDER BY table_schema, table_name`).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var schemaName, name string
		if err := rows.Scan(&schemaName, &name); err != nil {
			return nil, err
		}
		if schemaName == "public" {
			tables = append(tables, name)
		} else {
			tables = append(tables, schemaName+"."+name)
		}
	}
	return tables, nil
}

// SyncUsers consolidates user-related tables synchronization
func (api *ManagementAPI) SyncUsers() error {
	userTables := []string{"users", "teams", "roles", "user_roles", "user_teams"}
	for _, t := range api.getSyncOrder(userTables) {
		// Skip non-existent roles table gracefully if not configured
		if !api.tableExists(t) {
			continue
		}
		if err := api.performSyncTable(t); err != nil {
			appLogger.Warnf("User sync table %s failed: %v", t, err)
		}
	}
	return nil
}

// tableExists checks if a table exists in primary database
func (api *ManagementAPI) tableExists(table string) bool {
	schemaName, rawTable := splitQualifiedTable(table)
	var count int64
	api.manager.GetPrimaryDB().Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2`, strings.ToLower(schemaName), strings.ToLower(rawTable)).Scan(&count)
	return count > 0
}

func (api *ManagementAPI) getPrimaryKeyColumnFromDB(table string) (string, error) {
	schemaName, rawTable := splitQualifiedTable(table)

	rows, err := api.manager.GetPrimaryDB().Raw(`
		SELECT kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
		  ON tc.constraint_name = kcu.constraint_name
		 AND tc.table_schema = kcu.table_schema
		WHERE tc.table_schema = $1
		  AND tc.table_name = $2
		  AND tc.constraint_type = 'PRIMARY KEY'
		ORDER BY kcu.ordinal_position
	`, schemaName, rawTable).Rows()
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return "", err
		}
		return col, nil
	}

	return "", nil
}
