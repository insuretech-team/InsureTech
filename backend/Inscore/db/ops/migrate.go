package ops

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

// MigrationMetadata tracks migration execution history in PostgreSQL
type MigrationMetadata struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Type        string    `db:"type"` // 'proto', 'migration', or 'seeder'
	Schema      string    `db:"schema"`
	AppliedAt   time.Time `db:"applied_at"`
	Checksum    string    `db:"checksum"`
	ExecutionMS int64     `db:"execution_ms"`
	Status      string    `db:"status"` // 'success' or 'failed'
	ErrorMsg    string    `db:"error_msg,omitempty"`
}

// UnifiedMigrationManager handles all migration types: proto, SQL migrations, and seeders
type UnifiedMigrationManager struct {
	db             *sql.DB
	metadataSchema string
	metadataTable  string
	enabledSchemas []string
	migrationRoot  string
	seederRoot     string
	pruneColumns   bool // Rule 6: If true, drop columns not in proto
	strictMode     bool // Rule 6 Phase 3: Fail if zombies/drift detected

	appliedCache     map[string]bool // Cache for isApplied checks (schema:type:name)
	cacheMutex       sync.RWMutex    // Mutex for appliedCache
	cacheInitialized bool            // Flag to check if cache is initialized

	// Schema Snapshot Cache (Batch Metadata Loading)
	metadataMutex   sync.RWMutex
	schemaSnapshots map[string]*SchemaSnapshot // Map[schema] -> Snapshot
}

// SchemaSnapshot holds all metadata for a schema in memory
// This allows O(1) existence checks instead of DB queries
type SchemaSnapshot struct {
	Tables      map[string]bool                            // TableName -> Exists
	Columns     map[string]map[string]existingColumnDetail // TableName -> ColumnName -> Detail
	Constraints map[string]ConstraintDetail                // ConstraintName -> Detail
	Indexes     map[string]bool                            // IndexName -> Exists
}

type ConstraintDetail struct {
	Name       string
	Type       string // PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK
	Table      string
	Definition string // For checks
	FK         *ForeignKeyDetail
}

type ForeignKeyDetail struct {
	RefTable   string
	RefColumn  string
	DeleteRule string
	UpdateRule string
}

// NewUnifiedMigrationManager creates a new unified migration manager
func NewUnifiedMigrationManager(db *sql.DB) *UnifiedMigrationManager {
	// Find project root to build absolute paths
	projectRoot, err := findProjectRoot()
	if err != nil {
		appLogger.Warnf("Could not find project root, using relative paths: %v", err)
		projectRoot = "."
	}

	migrationRoot := filepath.Join(projectRoot, "backend", "inscore", "db", "migrations")
	seederRoot := filepath.Join(projectRoot, "backend", "inscore", "db", "seeds")

	umm := &UnifiedMigrationManager{
		db:             db,
		metadataSchema: "public",
		metadataTable:  "schema_migrations",
		enabledSchemas: []string{"public"},
		migrationRoot:  migrationRoot,
		seederRoot:     seederRoot,
	}

	// Auto-discover schemas from proto files
	umm.discoverSchemasFromProto()

	appLogger.Infof("Migration root: %s", migrationRoot)
	appLogger.Infof("Seeder root: %s", seederRoot)

	return umm
}

// LoadAppliedMigrations pre-fetches all applied migrations for a schema into memory
// This drastically reduces DB queries during migration checks
func (umm *UnifiedMigrationManager) LoadAppliedMigrations(schema string) error {
	fullTableName := umm.getFullTableName(umm.metadataSchema, umm.metadataTable)

	query := fmt.Sprintf("SELECT name, type, schema FROM %s WHERE schema = $1 AND status = 'success'", fullTableName)
	rows, err := umm.db.Query(query, schema)
	if err != nil {
		return fmt.Errorf("failed to load applied migrations for schema %s: %w", schema, err)
	}
	defer rows.Close()

	umm.cacheMutex.Lock()
	defer umm.cacheMutex.Unlock()

	if umm.appliedCache == nil {
		umm.appliedCache = make(map[string]bool)
	}

	for rows.Next() {
		var name, migType, sch string
		if err := rows.Scan(&name, &migType, &sch); err != nil {
			return err
		}
		// Key format: schema:type:name (normalized)
		key := fmt.Sprintf("%s:%s:%s", sch, migType, name)
		umm.appliedCache[key] = true
	}

	umm.cacheInitialized = true
	return nil
}

// LoadSchemaSnapshot loads all schema metadata into memory in one go
func (umm *UnifiedMigrationManager) LoadSchemaSnapshot(schema string) error {
	snapshot := &SchemaSnapshot{
		Tables:      make(map[string]bool),
		Columns:     make(map[string]map[string]existingColumnDetail),
		Constraints: make(map[string]ConstraintDetail),
		Indexes:     make(map[string]bool),
	}

	// 1. Load Tables
	tableQuery := `SELECT table_name FROM information_schema.tables WHERE table_schema = $1`
	rows, err := umm.db.Query(tableQuery, schema)
	if err != nil {
		return fmt.Errorf("snapshot tables failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil {
			snapshot.Tables[t] = true
		}
	}
	rows.Close()

	// 2. Load Columns
	colQuery := `
		SELECT table_name, column_name, data_type, is_nullable, column_default, character_maximum_length 
		FROM information_schema.columns 
		WHERE table_schema = $1`
	rows, err = umm.db.Query(colQuery, schema)
	if err != nil {
		return fmt.Errorf("snapshot columns failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t, col string
		var det existingColumnDetail
		var isNullableStr string
		if err := rows.Scan(&t, &col, &det.DataType, &isNullableStr, &det.DefaultValue, &det.CharMaxLength); err == nil {
			det.Name = col
			det.IsNullable = (isNullableStr == "YES")
			if snapshot.Columns[t] == nil {
				snapshot.Columns[t] = make(map[string]existingColumnDetail)
			}
			snapshot.Columns[t][strings.ToLower(col)] = det
		}
	}
	rows.Close()

	// 3. Load Constraints (PK, FK, Unique, Check)
	// Simplified: Load referential constraints for FK details
	// Complex join to get everything
	fkQuery := `
		SELECT 
			tc.constraint_name, tc.constraint_type, tc.table_name,
			kcu.column_name, 
			ccu.table_name AS ref_table,
			ccu.column_name AS ref_column,
			rc.delete_rule, rc.update_rule
		FROM information_schema.table_constraints tc
		LEFT JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
		LEFT JOIN information_schema.referential_constraints rc
			ON tc.constraint_name = rc.constraint_name AND tc.constraint_schema = rc.constraint_schema
		LEFT JOIN information_schema.constraint_column_usage ccu
			ON rc.unique_constraint_name = ccu.constraint_name AND rc.unique_constraint_schema = ccu.constraint_schema
		WHERE tc.table_schema = $1
	`
	rows, err = umm.db.Query(fkQuery, schema)
	if err != nil {
		return fmt.Errorf("snapshot constraints failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cName, cType, tName string
		var colName, refTable, refCol, delRule, updRule sql.NullString

		if err := rows.Scan(&cName, &cType, &tName, &colName, &refTable, &refCol, &delRule, &updRule); err == nil {
			cd := ConstraintDetail{
				Name:  cName,
				Type:  cType,
				Table: tName,
			}
			if cType == "FOREIGN KEY" {
				cd.FK = &ForeignKeyDetail{
					RefTable:   refTable.String,
					RefColumn:  refCol.String,
					DeleteRule: delRule.String,
					UpdateRule: updRule.String,
				}
			}
			snapshot.Constraints[cName] = cd
		}
	}
	rows.Close()

	// 4. Load Indexes
	idxQuery := `SELECT indexname FROM pg_indexes WHERE schemaname = $1`
	rows, err = umm.db.Query(idxQuery, schema)
	if err != nil {
		return fmt.Errorf("snapshot indexes failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var idx string
		if err := rows.Scan(&idx); err == nil {
			snapshot.Indexes[idx] = true
		}
	}
	rows.Close()

	umm.metadataMutex.Lock()
	if umm.schemaSnapshots == nil {
		umm.schemaSnapshots = make(map[string]*SchemaSnapshot)
	}
	umm.schemaSnapshots[schema] = snapshot
	umm.metadataMutex.Unlock()

	return nil
}

// SetPruneColumns enables/disables zombie column pruning
func (umm *UnifiedMigrationManager) SetPruneColumns(prune bool) {
	umm.pruneColumns = prune
}

// SetStrictMode enables/disables strict mode (fail on drift)
func (umm *UnifiedMigrationManager) SetStrictMode(strict bool) {
	umm.strictMode = strict
}

// findProjectRoot walks up the directory tree to find go.mod
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}
		dir = parent
	}
}

// discoverSchemasFromProto scans all proto files and extracts unique schemas
func (umm *UnifiedMigrationManager) discoverSchemasFromProto() {
	schemaSet := make(map[string]bool)
	schemaSet["public"] = true

	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if strings.HasPrefix(string(fd.Package()), "insuretech") {
			md := fd.Messages()
			for i := 0; i < md.Len(); i++ {
				msg := md.Get(i)
				tableOpts := umm.getTableOptions(msg)
				if tableOpts != nil {
					schema := umm.getSchemaFromProto(tableOpts)
					if schema != "" {
						schemaSet[schema] = true
					}
				}
			}
		}
		return true
	})

	var schemas []string
	for schema := range schemaSet {
		schemas = append(schemas, schema)
	}
	sort.Strings(schemas)

	if len(schemas) > 1 {
		umm.enabledSchemas = schemas
		appLogger.Infof("✓ Discovered %d schemas from proto: %v", len(schemas), schemas)
	}
}

// Initialize prepares the database for migrations
func (umm *UnifiedMigrationManager) Initialize() error {
	appLogger.Info("🚀 Initializing unified migration system...")

	if err := umm.enableExtensions(); err != nil {
		return fmt.Errorf("failed to enable extensions: %w", err)
	}

	if err := umm.ensureSchemas(umm.enabledSchemas...); err != nil {
		return fmt.Errorf("failed to ensure schemas: %w", err)
	}

	if err := umm.createMetadataTable(); err != nil {
		return fmt.Errorf("failed to create metadata table: %w", err)
	}

	appLogger.Info("✓ Unified migration system initialized")
	return nil
}

// enableExtensions enables required PostgreSQL extensions
func (umm *UnifiedMigrationManager) enableExtensions() error {
	extensions := []string{
		"uuid-ossp",
		"pg_trgm",
		"btree_gin",
	}

	for _, ext := range extensions {
		query := fmt.Sprintf(`CREATE EXTENSION IF NOT EXISTS "%s" WITH SCHEMA public`, ext)
		if _, err := umm.db.Exec(query); err != nil {
			if ext == "uuid-ossp" {
				return fmt.Errorf("failed to enable required extension %s: %w", ext, err)
			}
			appLogger.Warnf("Optional extension %s not available: %v", ext, err)
			continue
		}
		appLogger.Infof("✓ Extension enabled: %s", ext)
	}

	return nil
}

// ensureSchemas creates database schemas if they don't exist
func (umm *UnifiedMigrationManager) ensureSchemas(schemas ...string) error {
	for _, schema := range schemas {
		if schema == "" || schema == "public" {
			continue
		}

		query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", umm.quoteIdentifier(schema))
		if _, err := umm.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create schema %s: %w", schema, err)
		}

		appLogger.Infof("✓ Schema ensured: %s", schema)

		commentQuery := fmt.Sprintf(
			"COMMENT ON SCHEMA %s IS 'Schema for %s domain entities'",
			umm.quoteIdentifier(schema),
			schema,
		)
		if _, err := umm.db.Exec(commentQuery); err != nil {
			appLogger.Warnf("Failed to add comment to schema %s: %v", schema, err)
		}
	}

	return nil
}

// createMetadataTable creates the migration metadata tracking table
func (umm *UnifiedMigrationManager) createMetadataTable() error {
	fullTableName := umm.getFullTableName(umm.metadataSchema, umm.metadataTable)

	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name TEXT NOT NULL,
			type TEXT NOT NULL CHECK (type IN ('proto', 'migration', 'seeder')),
			schema TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			checksum TEXT NOT NULL,
			execution_ms BIGINT NOT NULL DEFAULT 0,
			status TEXT NOT NULL CHECK (status IN ('success', 'failed')),
			error_msg TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			UNIQUE(name, type, schema)
		)`, fullTableName)

	if _, err := umm.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create metadata table: %w", err)
	}

	indexes := []string{
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_type ON %s(type)", umm.metadataTable, fullTableName),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_schema ON %s(schema)", umm.metadataTable, fullTableName),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_status ON %s(status)", umm.metadataTable, fullTableName),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_applied_at ON %s(applied_at)", umm.metadataTable, fullTableName),
	}

	for _, idx := range indexes {
		if _, err := umm.db.Exec(idx); err != nil {
			appLogger.Warnf("Failed to create index: %v", err)
		}
	}

	comment := fmt.Sprintf(
		"COMMENT ON TABLE %s IS 'Tracks proto, migration, and seeder execution history'",
		fullTableName,
	)
	if _, err := umm.db.Exec(comment); err != nil {
		appLogger.Warnf("Failed to add table comment: %v", err)
	}

	appLogger.Infof("✓ Metadata table created: %s", fullTableName)
	return nil
}

// RunAll executes the complete migration flow: Proto -> SQL Migrations -> Seeders
func (umm *UnifiedMigrationManager) RunAll() error {
	appLogger.Info("==========================================")
	appLogger.Info("🔄 Starting Unified Migration Flow")
	appLogger.Info("==========================================")

	// Phase 1: Proto-driven table creation
	appLogger.Info("\n📦 PHASE 1: Proto-Driven Table Creation")
	appLogger.Info("------------------------------------------")
	if err := umm.RunProtoMigrations(); err != nil {
		return fmt.Errorf("proto migrations failed: %w", err)
	}

	// Phase 2: Custom SQL migrations for advanced features
	appLogger.Info("\n🔧 PHASE 2: Custom SQL Migrations")
	appLogger.Info("------------------------------------------")
	if err := umm.RunSQLMigrations(); err != nil {
		return fmt.Errorf("SQL migrations failed: %w", err)
	}

	// Phase 3: Seed data
	appLogger.Info("\n🌱 PHASE 3: Seeding Data")
	appLogger.Info("------------------------------------------")
	if err := umm.RunSeeders(); err != nil {
		return fmt.Errorf("seeders failed: %w", err)
	}

	appLogger.Info("\n==========================================")
	appLogger.Info("✅ All Migrations Completed Successfully!")
	appLogger.Info("==========================================")

	return nil
}

// TableMetadata holds all schema information extracted from a proto message
// This allows parallel discovery before sequential DB execution
type TableMetadata struct {
	Msg               protoreflect.MessageDescriptor
	Schema            string
	TableName         string
	Columns           []columnDef
	ForeignKeys       []ForeignKeyInfo
	Indexes           []*commonv1.IndexOptions
	UniqueConstraints []string
	TableOpts         *commonv1.TableOptions
	MigrationOrder    int32
}

// AnalyzeTable extracts all schema metadata from a proto message
// This is CPU-bound and safe to run in parallel
func (umm *UnifiedMigrationManager) AnalyzeTable(msg protoreflect.MessageDescriptor) *TableMetadata {
	tableOpts := umm.getTableOptions(msg)
	if tableOpts == nil {
		return nil
	}

	schema := umm.getSchemaFromProto(tableOpts)
	tableName := umm.getTableNameFromProto(msg, tableOpts)

	migOrder := int32(1000)
	if tableOpts.MigrationOrder > 0 {
		migOrder = tableOpts.MigrationOrder
	}

	appLogger.Infof("[DEBUG] AnalyzeTable: %s.%s (order: %d)", schema, tableName, migOrder)
	
	cols := umm.collectColumns(msg)
	fks := umm.getForeignKeys(msg)
	indexes := umm.getIndexes(msg)
	uniques := umm.getUniqueConstraints(msg)
	
	appLogger.Infof("[DEBUG] AnalyzeTable %s.%s: cols=%d, fks=%d, indexes=%d, uniques=%d", 
		schema, tableName, len(cols), len(fks), len(indexes), len(uniques))

	return &TableMetadata{
		Msg:               msg,
		Schema:            schema,
		TableName:         tableName,
		Columns:           cols,
		ForeignKeys:       fks,
		Indexes:           indexes,
		UniqueConstraints: uniques,
		TableOpts:         tableOpts,
		MigrationOrder:    migOrder,
	}
}

// RunProtoMigrations executes proto-driven schema creation (Phase 1)
func (umm *UnifiedMigrationManager) RunProtoMigrations() error {
	appLogger.Info("Running proto-driven schema-aware migrations (Parallel Discovery + Snapshot)...")

	// 1. Collect all messages
	var messages []protoreflect.MessageDescriptor
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if strings.HasPrefix(string(fd.Package()), "insuretech") {
			md := fd.Messages()
			for i := 0; i < md.Len(); i++ {
				messages = append(messages, md.Get(i))
			}
		}
		return true
	})

	appLogger.Infof("Found %d proto messages in insuretech package. Starting parallel analysis...", len(messages))

	// 2. Parallel Analysis (Map Phase)
	numWorkers := runtime.NumCPU() * 2
	if numWorkers < 4 {
		numWorkers = 4
	}

	jobs := make(chan protoreflect.MessageDescriptor, len(messages))
	results := make(chan *TableMetadata, len(messages))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range jobs {
				// AnalyzeTable: CPU bound
				meta := umm.AnalyzeTable(msg)
				results <- meta
			}
		}()
	}

	// Send jobs
	for _, msg := range messages {
		jobs <- msg
	}
	close(jobs)
	wg.Wait()
	close(results)

	// 3. Collect Results & Unique Schemas
	var allTables []*TableMetadata
	uniqueSchemas := make(map[string]bool)

	for meta := range results {
		if meta != nil {
			allTables = append(allTables, meta)
			uniqueSchemas[meta.Schema] = true
		}
	}

	// 4. Batch Load Schema Snapshots (I/O Phase)
	appLogger.Infof("Loading metadata snapshots for %d schemas...", len(uniqueSchemas))
	for schema := range uniqueSchemas {
		if err := umm.LoadSchemaSnapshot(schema); err != nil {
			return fmt.Errorf("failed to load snapshot for schema %s: %w", schema, err)
		}
	}

	// 5. Sort by order
	sort.Slice(allTables, func(i, j int) bool {
		return allTables[i].MigrationOrder < allTables[j].MigrationOrder
	})

	appLogger.Infof("Analyzed %d tables. Starting DB execution...", len(allTables))

	// 6. Execution (Sequential DB Ops using In-Memory Snapshots)
	// Phase 1A: Create all tables without foreign keys
	appLogger.Info("Phase 1A: Creating tables...")
	tablesCreated := 0
	var tablesWithForeignKeys []*TableMetadata

	for _, meta := range allTables {
		schema := meta.Schema
		tableName := meta.TableName
		protoName := string(meta.Msg.FullName())

		appLogger.Infof("Processing: %s.%s (order: %d)", schema, tableName, meta.MigrationOrder)

		// Ensure schema exists (Snapshot doesn't track schema existence, DB does)
		if err := umm.ensureSchemas(schema); err != nil {
			return fmt.Errorf("failed to ensure schema %s: %w", schema, err)
		}

		// Check if table exists (Using Snapshot)
		exists, err := umm.tableExistsInSchema(schema, tableName)
		if err != nil {
			return fmt.Errorf("failed to check table existence: %w", err)
		}

		if !exists {
			// Create table
			if err := umm.createTableWithConstraints(schema, tableName, meta.Columns, meta.UniqueConstraints); err != nil {
				return fmt.Errorf("failed to create table %s.%s: %w", schema, tableName, err)
			}

			// Add comments
			if err := umm.addTableComment(schema, tableName, umm.getTableComment(meta.TableOpts)); err != nil {
				appLogger.Warnf("Failed to add table comment: %v", err)
			}
			if err := umm.addColumnComments(schema, tableName, meta.Msg); err != nil {
				appLogger.Warnf("Failed to add column comments: %v", err)
			}

			// Add constraints & indexes
			if err := umm.addCheckConstraints(schema, tableName, meta.Columns); err != nil {
				appLogger.Warnf("Failed to add CHECK constraints: %v", err)
			}
			if len(meta.Indexes) > 0 {
				if err := umm.addIndexes(schema, tableName, meta.Indexes); err != nil {
					appLogger.Warnf("Failed to add indexes: %v", err)
				}
			}

			tablesCreated++
			appLogger.Infof("✓ Created table: %s.%s", schema, tableName)

			// Record migration
			checksum := umm.calculateChecksum([]byte(protoName))
			if err := umm.recordMigration(protoName, "proto", schema, checksum, 0, "success", ""); err != nil {
				appLogger.Warnf("Failed to record proto migration: %v", err)
			}
		} else {
			// Table exists, sync columns (Snapshot based)
			if err := umm.syncTableColumns(meta.Msg, schema, tableName, meta.Columns); err != nil {
				appLogger.Warnf("Failed to sync columns: %v", err)
			}
			if err := umm.addCheckConstraints(schema, tableName, meta.Columns); err != nil {
				appLogger.Warnf("Failed to sync CHECK constraints: %v", err)
			}
			if len(meta.Indexes) > 0 {
				if err := umm.addIndexes(schema, tableName, meta.Indexes); err != nil {
					appLogger.Warnf("Failed to add missing indexes: %v", err)
				}
			}
			appLogger.Infof("✓ Table exists, verified: %s.%s", schema, tableName)
		}

		if len(meta.ForeignKeys) > 0 {
			tablesWithForeignKeys = append(tablesWithForeignKeys, meta)
			appLogger.Infof("  Queued %d FKs for %s.%s", len(meta.ForeignKeys), schema, tableName)
		}
	}

	appLogger.Infof("Phase 1A complete: %d tables created", tablesCreated)
	appLogger.Infof("Tables with foreign keys queued: %d", len(tablesWithForeignKeys))

	// Phase 1B: Add all foreign keys (Snapshot based)
	appLogger.Info("Phase 1B: Syncing foreign key constraints...")
	foreignKeysAdded := 0
	for _, meta := range tablesWithForeignKeys {
		if err := umm.syncForeignKeys(meta.Schema, meta.TableName, meta.ForeignKeys); err != nil {
			appLogger.Warnf("Failed to sync foreign keys for %s.%s: %v", meta.Schema, meta.TableName, err)
		} else {
			foreignKeysAdded += len(meta.ForeignKeys)
		}
	}

	appLogger.Infof("Phase 1B complete: %d foreign keys added", foreignKeysAdded)
	appLogger.Info("✓ Proto-driven migrations completed")

	return nil
}

// RunSQLMigrations executes custom SQL migrations (Phase 2)
func (umm *UnifiedMigrationManager) RunSQLMigrations() error {
	if !umm.dirExists(umm.migrationRoot) {
		appLogger.Infof("Migration root does not exist: %s", umm.migrationRoot)
		return nil
	}

	totalApplied := 0

	// Process each schema directory
	for _, schema := range umm.enabledSchemas {
		schemaDir := filepath.Join(umm.migrationRoot, schema)
		if !umm.dirExists(schemaDir) {
			continue
		}

		files, err := umm.collectMigrationFiles(schemaDir)
		if err != nil {
			return fmt.Errorf("failed to collect migrations from %s: %w", schemaDir, err)
		}

		if len(files) == 0 {
			continue
		}

		appLogger.Infof("Processing migrations for schema: %s (%d files)", schema, len(files))

		// Pre-load applied migrations for this schema (Cache Optimization)
		if err := umm.LoadAppliedMigrations(schema); err != nil {
			appLogger.Warnf("Failed to pre-load applied migrations for schema %s: %v", schema, err)
		}

		for _, file := range files {
			if umm.isApplied(file.RelPath, "migration", schema) {
				appLogger.Infof("⊘ Already applied: %s", file.RelPath)
				continue
			}

			if err := umm.applyMigration(file, schema); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", file.RelPath, err)
			}

			totalApplied++
			appLogger.Infof("✓ Applied migration: %s [%s]", file.RelPath, schema)
		}
	}

	appLogger.Infof("✓ Total SQL migrations applied: %d", totalApplied)
	return nil
}

// RunSeeders executes seed data scripts (Phase 3)
func (umm *UnifiedMigrationManager) RunSeeders() error {
	if !umm.dirExists(umm.seederRoot) {
		appLogger.Infof("Seeder root does not exist: %s", umm.seederRoot)
		return nil
	}

	totalApplied := 0

	for _, schema := range umm.enabledSchemas {
		schemaDir := filepath.Join(umm.seederRoot, schema)
		if !umm.dirExists(schemaDir) {
			continue
		}

		files, err := umm.collectSeederFiles(schemaDir)
		if err != nil {
			return fmt.Errorf("failed to collect seeders from %s: %w", schemaDir, err)
		}

		if len(files) == 0 {
			continue
		}

		appLogger.Infof("Processing seeders for schema: %s (%d files)", schema, len(files))

		// Pre-load applied migrations for this schema (Cache Optimization)
		if err := umm.LoadAppliedMigrations(schema); err != nil {
			appLogger.Warnf("Failed to pre-load applied seeders for schema %s: %v", schema, err)
		}

		for _, file := range files {
			if umm.isApplied(file.RelPath, "seeder", schema) {
				appLogger.Infof("⊘ Already applied: %s", file.RelPath)
				continue
			}

			if err := umm.applySeeder(file, schema); err != nil {
				appLogger.Warnf("Failed to apply seeder %s: %v (continuing...)", file.RelPath, err)
				continue
			}

			totalApplied++
			appLogger.Infof("✓ Applied seeder: %s [%s]", file.RelPath, schema)
		}
	}

	appLogger.Infof("✓ Total seeders applied: %d", totalApplied)
	return nil
}

// MigrationFile represents a migration or seeder file
type MigrationFile struct {
	AbsPath string
	RelPath string
	Content []byte
}

// collectMigrationFiles collects .up.sql migration files from a directory
func (umm *UnifiedMigrationManager) collectMigrationFiles(dir string) ([]MigrationFile, error) {
	var files []MigrationFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".up.sql") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = info.Name()
		}

		files = append(files, MigrationFile{
			AbsPath: path,
			RelPath: filepath.ToSlash(relPath),
			Content: content,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].RelPath < files[j].RelPath
	})

	return files, nil
}

// collectSeederFiles collects .sql seeder files from a directory
func (umm *UnifiedMigrationManager) collectSeederFiles(dir string) ([]MigrationFile, error) {
	var files []MigrationFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".sql" || strings.HasSuffix(info.Name(), ".up.sql") || strings.HasSuffix(info.Name(), ".down.sql") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = info.Name()
		}

		files = append(files, MigrationFile{
			AbsPath: path,
			RelPath: filepath.ToSlash(relPath),
			Content: content,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].RelPath < files[j].RelPath
	})

	return files, nil
}

// applyMigration applies a migration file with transactional safety
func (umm *UnifiedMigrationManager) applyMigration(file MigrationFile, schema string) error {
	startTime := time.Now()
	checksum := umm.calculateChecksum(file.Content)

	tx, err := umm.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(file.Content)); err != nil {
		umm.recordMigration(file.RelPath, "migration", schema, checksum, time.Since(startTime).Milliseconds(), "failed", err.Error())
		return fmt.Errorf("migration execution failed: %w", err)
	}

	if err := umm.recordMigrationTx(tx, file.RelPath, "migration", schema, checksum, time.Since(startTime).Milliseconds(), "success", ""); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// applySeeder applies a seeder file (no transaction for partial success)
func (umm *UnifiedMigrationManager) applySeeder(file MigrationFile, schema string) error {
	startTime := time.Now()
	checksum := umm.calculateChecksum(file.Content)

	if _, err := umm.db.Exec(string(file.Content)); err != nil {
		umm.recordMigration(file.RelPath, "seeder", schema, checksum, time.Since(startTime).Milliseconds(), "failed", err.Error())
		return fmt.Errorf("seeder execution failed: %w", err)
	}

	if err := umm.recordMigration(file.RelPath, "seeder", schema, checksum, time.Since(startTime).Milliseconds(), "success", ""); err != nil {
		return fmt.Errorf("failed to record seeder: %w", err)
	}

	return nil
}

// isApplied checks if a migration/seeder has been applied successfully
func (umm *UnifiedMigrationManager) isApplied(name, migType, schema string) bool {
	fullTableName := umm.getFullTableName(umm.metadataSchema, umm.metadataTable)

	var count int
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE name = $1 AND type = $2 AND schema = $3 AND status = 'success'",
		fullTableName,
	)

	err := umm.db.QueryRow(query, name, migType, schema).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// recordMigration records migration execution in metadata table
func (umm *UnifiedMigrationManager) recordMigration(name, migType, schema, checksum string, executionMS int64, status, errorMsg string) error {
	fullTableName := umm.getFullTableName(umm.metadataSchema, umm.metadataTable)

	query := fmt.Sprintf(`
		INSERT INTO %s (name, type, schema, checksum, execution_ms, status, error_msg)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (name, type, schema) 
		DO UPDATE SET 
			applied_at = now(),
			checksum = EXCLUDED.checksum,
			execution_ms = EXCLUDED.execution_ms,
			status = EXCLUDED.status,
			error_msg = EXCLUDED.error_msg,
			updated_at = now()
	`, fullTableName)

	_, err := umm.db.Exec(query, name, migType, schema, checksum, executionMS, status, errorMsg)
	return err
}

// recordMigrationTx records migration within a transaction
func (umm *UnifiedMigrationManager) recordMigrationTx(tx *sql.Tx, name, migType, schema, checksum string, executionMS int64, status, errorMsg string) error {
	fullTableName := umm.getFullTableName(umm.metadataSchema, umm.metadataTable)

	query := fmt.Sprintf(`
		INSERT INTO %s (name, type, schema, checksum, execution_ms, status, error_msg)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (name, type, schema) 
		DO UPDATE SET 
			applied_at = now(),
			checksum = EXCLUDED.checksum,
			execution_ms = EXCLUDED.execution_ms,
			status = EXCLUDED.status,
			error_msg = EXCLUDED.error_msg,
			updated_at = now()
	`, fullTableName)

	_, err := tx.Exec(query, name, migType, schema, checksum, executionMS, status, errorMsg)
	return err
}

// GetMigrationStatus returns the status of all migrations
func (umm *UnifiedMigrationManager) GetMigrationStatus() ([]MigrationMetadata, error) {
	fullTableName := umm.getFullTableName(umm.metadataSchema, umm.metadataTable)

	query := fmt.Sprintf(`
		SELECT id, name, type, schema, applied_at, checksum, execution_ms, status, COALESCE(error_msg, '') as error_msg
		FROM %s
		ORDER BY schema, type, applied_at DESC
	`, fullTableName)

	rows, err := umm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migration status: %w", err)
	}
	defer rows.Close()

	var migrations []MigrationMetadata
	for rows.Next() {
		var m MigrationMetadata
		if err := rows.Scan(&m.ID, &m.Name, &m.Type, &m.Schema, &m.AppliedAt, &m.Checksum, &m.ExecutionMS, &m.Status, &m.ErrorMsg); err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		migrations = append(migrations, m)
	}

	return migrations, rows.Err()
}

// calculateChecksum calculates SHA256 checksum
func (umm *UnifiedMigrationManager) calculateChecksum(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}

// dirExists checks if a directory exists
func (umm *UnifiedMigrationManager) dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// getFullTableName returns schema.table format
func (umm *UnifiedMigrationManager) getFullTableName(schema, table string) string {
	if schema == "" {
		schema = "public"
	}
	return fmt.Sprintf("%s.%s", umm.quoteIdentifier(schema), umm.quoteIdentifier(table))
}

// quoteIdentifier quotes SQL identifiers
func (umm *UnifiedMigrationManager) quoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(name, `"`, `""`))
}

// ==================== PROTO HELPER METHODS ====================

// getTableOptions extracts TableOptions from proto message options
func (umm *UnifiedMigrationManager) getTableOptions(md protoreflect.MessageDescriptor) *commonv1.TableOptions {
	opts := md.Options()
	if opts == nil {
		return nil
	}

	if proto.HasExtension(opts, commonv1.E_Table) {
		val := proto.GetExtension(opts, commonv1.E_Table)
		if tableOpts, ok := val.(*commonv1.TableOptions); ok {
			return tableOpts
		}
	}

	return nil
}

// getSchemaFromProto extracts schema name from TableOptions
func (umm *UnifiedMigrationManager) getSchemaFromProto(opts *commonv1.TableOptions) string {
	if opts == nil {
		return "public"
	}
	schema := opts.SchemaName
	if schema == "" {
		return "public"
	}
	return schema
}

// getTableNameFromProto gets table name from proto message
func (umm *UnifiedMigrationManager) getTableNameFromProto(md protoreflect.MessageDescriptor, opts *commonv1.TableOptions) string {
	if opts != nil && opts.TableName != "" {
		return opts.TableName
	}
	base := umm.snakeCase(string(md.Name()))
	return umm.pluralize(base)
}

// getTableComment extracts table comment from TableOptions
func (umm *UnifiedMigrationManager) getTableComment(opts *commonv1.TableOptions) string {
	if opts != nil && opts.Comment != "" {
		return opts.Comment
	}
	return ""
}

// existingColumnDetail captures rich metadata for type comparison
type existingColumnDetail struct {
	Name          string
	DataType      string
	IsNullable    bool
	DefaultValue  sql.NullString
	CharMaxLength sql.NullInt64
}

// tableExistsInSchema checks if table exists in specific schema
func (umm *UnifiedMigrationManager) tableExistsInSchema(schema, table string) (bool, error) {
	// 1. Check Snapshot
	umm.cacheMutex.RLock()
	snap, ok := umm.schemaSnapshots[schema]
	umm.cacheMutex.RUnlock()

	if ok {
		return snap.Tables[table], nil
	}

	// 2. Fallback to DB
	var n int
	err := umm.db.QueryRow(
		"SELECT count(*) FROM information_schema.tables WHERE table_schema=$1 AND table_name=$2",
		schema, table,
	).Scan(&n)
	return n > 0, err
}

// getExistingColumnsInSchema gets existing columns with rich metadata for a table
func (umm *UnifiedMigrationManager) getExistingColumnsInSchema(schema, table string) (map[string]existingColumnDetail, error) {
	// 1. Check Snapshot
	umm.cacheMutex.RLock()
	snap, ok := umm.schemaSnapshots[schema]
	umm.cacheMutex.RUnlock()

	if ok {
		if cols, found := snap.Columns[table]; found {
			// Return a copy to avoid mutation issues if any (though we only read)
			return cols, nil
		}
		// If table known to exist but no columns found, return empty map
		if snap.Tables[table] {
			return make(map[string]existingColumnDetail), nil
		}
		// If table not in snapshot, maybe it doesn't exist, fallback to DB to be safe or return error?
		// For now, consistent behavior: if snapshot exists, trust it.
		return nil, fmt.Errorf("table %s not found in schema snapshot", table)
	}

	// 2. Fallback to DB
	query := `
		SELECT 
			column_name, 
			data_type, 
			is_nullable, 
			column_default, 
			character_maximum_length
		FROM information_schema.columns 
		WHERE table_schema=$1 AND table_name=$2
	`
	rows, err := umm.db.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[string]existingColumnDetail)
	for rows.Next() {
		var c existingColumnDetail
		var isNullableStr string

		if err := rows.Scan(
			&c.Name,
			&c.DataType,
			&isNullableStr,
			&c.DefaultValue,
			&c.CharMaxLength,
		); err != nil {
			return nil, err
		}

		c.IsNullable = (isNullableStr == "YES")
		m[strings.ToLower(c.Name)] = c
	}
	return m, rows.Err()
}

// columnDef represents a column definition
type columnDef struct {
	name            string
	typ             string
	defaultExp      string
	primary         bool
	notNull         bool
	unique          bool
	checkConstraint string // CHECK constraint expression from proto
	comment         string // Column comment from proto
}

// collectColumns extracts column definitions from proto message
func (umm *UnifiedMigrationManager) collectColumns(md protoreflect.MessageDescriptor) []columnDef {
	var cols []columnDef
	fields := md.Fields()
	seenColumns := make(map[string]bool) // Track column names to detect duplicates

	appLogger.Infof("  [DEBUG] collectColumns: Starting for message %s with %d fields", md.Name(), fields.Len())

	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		colName := umm.snakeCase(string(f.Name()))

		// Check for explicit column_name in options
		fOpts := f.Options()
		if fOpts != nil && proto.HasExtension(fOpts, commonv1.E_Column) {
			if colOpts, ok := proto.GetExtension(fOpts, commonv1.E_Column).(*commonv1.ColumnOptions); ok && colOpts != nil {
				if colOpts.ColumnName != "" {
					colName = colOpts.ColumnName
				}
			}
		}
		
		// DEBUG: Log every field being processed
		appLogger.Infof("    [DEBUG] Field %d: proto_name=%s, column_name=%s, field_number=%d", i, f.Name(), colName, f.Number())
		
		// Check for duplicate column names
		if seenColumns[colName] {
			appLogger.Warnf("    [WARN] DUPLICATE COLUMN DETECTED AND SKIPPED: %s already exists in message %s!", colName, md.FullName())
			appLogger.Warnf("    [WARN] Duplicate is field %d (%s) with field_number=%d trying to use column name %s", i, f.Name(), f.Number(), colName)
			continue // Skip duplicate
		}
		
		c := columnDef{name: colName}

		c.typ, c.defaultExp = umm.mapFieldType(f, colName)
		if c.typ == "" {
			appLogger.Infof("    [DEBUG] Skipping field %s - no type mapped", f.Name())
			continue
		}

		// Check field options - these OVERRIDE auto-detected values
		if fOpts != nil && proto.HasExtension(fOpts, commonv1.E_Column) {
			if colOpts, ok := proto.GetExtension(fOpts, commonv1.E_Column).(*commonv1.ColumnOptions); ok && colOpts != nil {
				c.notNull = colOpts.NotNull
				c.unique = colOpts.Unique
				if colOpts.PrimaryKey {
					c.primary = true
					appLogger.Infof("    [DEBUG] Column %s marked as PRIMARY KEY", colName)
				}
				// If sql_type is specified, use it and clear auto-default
				if colOpts.SqlType != "" {
					c.typ = colOpts.SqlType
					// Reset default when type is overridden, unless explicitly set
					if colOpts.DefaultValue == "" {
						c.defaultExp = "NULL"
					}
				}
				// If default_value is explicitly set, use it
				if colOpts.DefaultValue != "" {
					c.defaultExp = colOpts.DefaultValue
				}
				// Extract CHECK constraint from proto
				if colOpts.CheckConstraint != "" {
					c.checkConstraint = colOpts.CheckConstraint
				}
				// Extract column comment from proto
				if colOpts.Comment != "" {
					c.comment = colOpts.Comment
				}
				// Log if this column has a foreign key
				if colOpts.ForeignKey != nil {
					appLogger.Infof("    [DEBUG] Column %s has FK to %s.%s", colName, colOpts.ForeignKey.ReferencesTable, colOpts.ForeignKey.ReferencesColumn)
				}
			}
		}

		cols = append(cols, c)
		seenColumns[colName] = true
		appLogger.Infof("    [DEBUG] Added column: %s %s (primary=%v, notNull=%v, unique=%v)", colName, c.typ, c.primary, c.notNull, c.unique)
	}
	
	appLogger.Infof("  [DEBUG] Total columns collected: %d, Unique columns: %d", len(cols), len(seenColumns))
	
	// Double-check for duplicates in final list
	finalCheck := make(map[string]int)
	for _, col := range cols {
		finalCheck[col.name]++
	}
	for name, count := range finalCheck {
		if count > 1 {
			appLogger.Errorf("  [ERROR] FINAL CHECK: Column %s appears %d times in cols array!", name, count)
		}
	}
	
	return cols
}

// mapFieldType maps proto field types to SQL types
func (umm *UnifiedMigrationManager) mapFieldType(f protoreflect.FieldDescriptor, col string) (string, string) {
	if f.Cardinality() == protoreflect.Repeated {
		if f.Kind() == protoreflect.StringKind {
			return "TEXT[]", "'{}'::text[]"
		}
		return "JSONB", "'{}'::jsonb"
	}

	switch f.Kind() {
	case protoreflect.BoolKind:
		return "BOOLEAN", "NULL"
	case protoreflect.Int32Kind:
		return "INTEGER", "NULL"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "BIGINT", "NULL"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "DOUBLE PRECISION", "NULL"
	case protoreflect.EnumKind:
		return "INTEGER", "NULL"
	case protoreflect.StringKind:
		return "TEXT", "NULL"
	case protoreflect.MessageKind:
		m := f.Message()
		full := string(m.FullName())
		switch full {
		case "google.protobuf.Timestamp":
			if strings.HasSuffix(col, "created_at") || strings.HasSuffix(col, "updated_at") {
				return "TIMESTAMPTZ", "now()"
			}
			return "TIMESTAMPTZ", "NULL"
		case "google.protobuf.Struct":
			return "JSONB", "'{}'::jsonb"
		case "insuretech.common.v1.UUID":
			// For UUID type, always use NULL as default (unless explicitly set)
			// Primary key UUIDs will have default set via column options
			return "UUID", "NULL"
		default:
			return "JSONB", "'{}'::jsonb"
		}
	}
	return "", ""
}

// getUniqueConstraints extracts unique constraints from proto message
func (umm *UnifiedMigrationManager) getUniqueConstraints(md protoreflect.MessageDescriptor) []string {
	var uniqueFields []string
	fields := md.Fields()

	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		fOpts := f.Options()

		if fOpts != nil && proto.HasExtension(fOpts, commonv1.E_Column) {
			if colOpts, ok := proto.GetExtension(fOpts, commonv1.E_Column).(*commonv1.ColumnOptions); ok && colOpts != nil {
				if colOpts.Unique {
					uniqueFields = append(uniqueFields, umm.snakeCase(string(f.Name())))
				}
			}
		}
	}

	return uniqueFields
}

// getForeignKeys extracts foreign keys from proto message
// ForeignKeyInfo wraps ForeignKey with the column name it applies to
type ForeignKeyInfo struct {
	FK         *commonv1.ForeignKey
	ColumnName string
}

func (umm *UnifiedMigrationManager) getForeignKeys(md protoreflect.MessageDescriptor) []ForeignKeyInfo {
	var foreignKeys []ForeignKeyInfo
	fields := md.Fields()

	tableName := string(md.Name())

	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		fOpts := f.Options()

		if fOpts != nil && proto.HasExtension(fOpts, commonv1.E_Column) {
			if colOpts, ok := proto.GetExtension(fOpts, commonv1.E_Column).(*commonv1.ColumnOptions); ok && colOpts != nil {
				if colOpts.ForeignKey != nil {
					colName := umm.snakeCase(string(f.Name()))

					if colOpts.ColumnName != "" {
						colName = colOpts.ColumnName
					}

					foreignKeys = append(foreignKeys, ForeignKeyInfo{
						FK:         colOpts.ForeignKey,
						ColumnName: colName,
					})

					appLogger.Infof("  Found FK in %s.%s -> %s.%s", tableName, colName, colOpts.ForeignKey.ReferencesTable, colOpts.ForeignKey.ReferencesColumn)
				}
			}
		}
	}

	if len(foreignKeys) > 0 {
		appLogger.Infof("Table %s has %d foreign keys", tableName, len(foreignKeys))
	}

	return foreignKeys
}

// getIndexes extracts index definitions from proto message fields
// Note: In the new proto structure, indexes are defined per-column in IndexOptions
func (umm *UnifiedMigrationManager) getIndexes(md protoreflect.MessageDescriptor) []*commonv1.IndexOptions {
	var indexes []*commonv1.IndexOptions
	fields := md.Fields()

	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		fOpts := f.Options()

		if fOpts != nil && proto.HasExtension(fOpts, commonv1.E_Column) {
			if colOpts, ok := proto.GetExtension(fOpts, commonv1.E_Column).(*commonv1.ColumnOptions); ok && colOpts != nil {
				if colOpts.Index != nil {
					indexes = append(indexes, colOpts.Index)
				}
			}
		}
	}

	return indexes
}

// ==================== TABLE CREATION METHODS ====================

// createTableWithConstraints creates a table with columns and constraints
func (umm *UnifiedMigrationManager) createTableWithConstraints(schema, tableName string, cols []columnDef, uniqueConstraints []string) error {
	fullTableName := umm.getFullTableName(schema, tableName)

	var parts []string
	var primaryKeys []string
	
	appLogger.Infof("  [DEBUG] Creating table %s with %d columns", fullTableName, len(cols))
	
	// Check for duplicate column names in input
	colNames := make(map[string]int)
	for _, c := range cols {
		colNames[c.name]++
	}
	for name, count := range colNames {
		if count > 1 {
			appLogger.Errorf("  [ERROR] createTableWithConstraints: Column %s appears %d times in cols array!", name, count)
		}
	}

	for idx, c := range cols {
		colDef := fmt.Sprintf("%s %s", umm.quoteIdentifier(c.name), c.typ)
		
		appLogger.Infof("  [DEBUG] Processing column %d/%d: %s (primary=%v, notNull=%v, unique=%v)", idx+1, len(cols), c.name, c.primary, c.notNull, c.unique)

		if c.primary {
			primaryKeys = append(primaryKeys, c.name)
		}

		if c.notNull && !c.primary {
			colDef += " NOT NULL"
		}

		if c.defaultExp != "" && c.defaultExp != "NULL" {
			colDef += fmt.Sprintf(" DEFAULT %s", c.defaultExp)
		}

		parts = append(parts, colDef)
		appLogger.Infof("  [DEBUG] Column def: %s", colDef)
	}

	// Add primary key constraint
	if len(primaryKeys) > 0 {
		pkCols := make([]string, len(primaryKeys))
		for i, pk := range primaryKeys {
			pkCols[i] = umm.quoteIdentifier(pk)
		}
		pkConstraint := fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkCols, ", "))
		parts = append(parts, pkConstraint)
		appLogger.Infof("  [DEBUG] Adding PK constraint: %s", pkConstraint)
	}

	// Add unique constraints
	for _, uc := range uniqueConstraints {
		parts = append(parts, fmt.Sprintf("UNIQUE (%s)", umm.quoteIdentifier(uc)))
	}

	query := fmt.Sprintf("CREATE TABLE %s (\n\t%s\n)", fullTableName, strings.Join(parts, ",\n\t"))
	
	appLogger.Infof("  [DEBUG] Final CREATE TABLE SQL:\n%s", query)
	appLogger.Infof("  [DEBUG] Total parts in CREATE TABLE: %d (columns: %d, PK: %d, unique: %d)", len(parts), len(cols), len(primaryKeys), len(uniqueConstraints))

	if _, err := umm.db.Exec(query); err != nil {
		appLogger.Errorf("  [ERROR] CREATE TABLE failed with error: %v", err)
		appLogger.Errorf("  [ERROR] SQL was:\n%s", query)
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// syncTableColumns ensures DB columns match Proto definition (Add, Alter, Warn)
func (umm *UnifiedMigrationManager) syncTableColumns(md protoreflect.MessageDescriptor, schema, tableName string, cols []columnDef) error {
	existing, err := umm.getExistingColumnsInSchema(schema, tableName)
	if err != nil {
		return fmt.Errorf("failed to get existing columns: %w", err)
	}

	// 1. Add missing or Update existing
	for _, c := range cols {
		existingCol, ok := existing[strings.ToLower(c.name)]
		if !ok {
			// Missing: Add it
			if err := umm.addColumn(schema, tableName, c); err != nil {
				return fmt.Errorf("failed to add column %s: %w", c.name, err)
			}
			appLogger.Infof("  ✓ Added column: %s.%s.%s", schema, tableName, c.name)
		} else {
			// Existing: Check for Type Drift (Rule 8)
			if err := umm.checkAndAlterColumnType(schema, tableName, c, existingCol); err != nil {
				appLogger.Warnf("  ! Failed to check/alter column %s: %v", c.name, err)
			}
		}
	}

	// 2. Detect Deprecated/Zombie Columns (Rule 6)
	for existingName := range existing {
		found := false
		for _, c := range cols {
			if strings.EqualFold(c.name, existingName) {
				found = true
				break
			}
		}
		if !found {
			if umm.pruneColumns {
				// PRUNE MODE: Drop the column
				appLogger.Warnf("  🗑️  Pruning zombie column: %s.%s.%s", schema, tableName, existingName)
				fullTableName := umm.getFullTableName(schema, tableName)
				query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", fullTableName, umm.quoteIdentifier(existingName))
				if _, err := umm.db.Exec(query); err != nil {
					appLogger.Warnf("  ! Failed to prune column %s: %v", existingName, err)
				} else {
					appLogger.Infof("  ✓ Pruned column: %s", existingName)
				}
			} else {
				// STANDARD MODE: Warn only
				if umm.strictMode {
					return fmt.Errorf("strict mode violation: zombie column found: %s.%s.%s", schema, tableName, existingName)
				}
				appLogger.Warnf("  ⚠️  Zombie column found: %s.%s.%s (In DB but not in Proto). Use --prune to remove.", schema, tableName, existingName)
			}
		}
	}

	return nil
}

// isSerialPseudoType returns true for SERIAL/BIGSERIAL/SMALLSERIAL pseudo-types.
// PostgreSQL does not allow these in ALTER COLUMN ... TYPE statements; they are
// only valid in CREATE TABLE. The underlying storage type (integer/bigint/smallint)
// is always what information_schema.columns reports, so a mismatch against these
// is a false positive — no ALTER is ever needed or valid.
func isSerialPseudoType(t string) bool {
	switch strings.ToUpper(strings.TrimSpace(t)) {
	case "SERIAL", "BIGSERIAL", "SMALLSERIAL":
		return true
	}
	return false
}

// checkAndAlterColumnType compares Proto type vs DB type and alters if mismatched
func (umm *UnifiedMigrationManager) checkAndAlterColumnType(schema, table string, desired columnDef, actual existingColumnDetail) error {
	// Basic compatibility check
	if umm.isTypeCompatible(desired.typ, actual.DataType) {
		return nil
	}

	// Serial pseudo-types (BIGSERIAL, SERIAL, SMALLSERIAL) are stored as their
	// underlying integer types in Postgres. ALTER COLUMN ... TYPE BIGSERIAL is
	// invalid SQL (error 42704). Skip silently — the column is already correct.
	if isSerialPseudoType(desired.typ) {
		appLogger.Infof("  ✓ Skipping serial pseudo-type: %s.%s (Proto: %s stored as DB: %s — compatible)",
			table, desired.name, desired.typ, actual.DataType)
		return nil
	}

	appLogger.Warnf("  ≠ Type Mismatch: %s.%s (Proto: %s, DB: %s). Fixing...", table, desired.name, desired.typ, actual.DataType)

	fullTableName := umm.getFullTableName(schema, table)

	// Construct ALTER statement with USING clause for safe casting
	query := fmt.Sprintf(
		"ALTER TABLE %s ALTER COLUMN %s TYPE %s USING %s::%s",
		fullTableName,
		umm.quoteIdentifier(desired.name),
		desired.typ,
		umm.quoteIdentifier(desired.name),
		desired.typ,
	)

	if _, err := umm.db.Exec(query); err != nil {
		return fmt.Errorf("ALTER TYPE failed: %w", err)
	}

	appLogger.Infof("  ✓ Altered column type: %s.%s -> %s", table, desired.name, desired.typ)
	return nil
}

// isTypeCompatible checks if SQL types are effectively the same
func (umm *UnifiedMigrationManager) isTypeCompatible(protoType, dbType string) bool {
	p := strings.ToUpper(protoType)
	d := strings.ToUpper(dbType)

	if p == d {
		return true
	}

	// Common aliases mapping
	// Proto uses UPPER CASE standard SQL (e.g. TIMESTAMPTZ, DOUBLE PRECISION)
	// Postgres information_schema uses lower case specific names (e.g. timestamp with time zone, double precision)

	// Normalize aliases
	switch p {
	case "TIMESTAMPTZ":
		if d == "TIMESTAMP WITH TIME ZONE" {
			return true
		}
	case "INT", "INTEGER":
		if d == "INTEGER" || d == "INT4" {
			return true
		}
	// SERIAL types: Postgres stores them as their underlying integer type in
	// information_schema.columns (BIGSERIAL→bigint, SERIAL→integer, SMALLSERIAL→smallint).
	// They are compatible — no ALTER needed.
	case "BIGSERIAL":
		if d == "BIGINT" || d == "INT8" {
			return true
		}
	case "SERIAL":
		if d == "INTEGER" || d == "INT4" || d == "INT" {
			return true
		}
	case "SMALLSERIAL":
		if d == "SMALLINT" || d == "INT2" {
			return true
		}
	case "BIGINT":
		if d == "BIGINT" || d == "INT8" {
			return true
		}
	case "BOOLEAN":
		if d == "BOOLEAN" || d == "BOOL" {
			return true
		}
	case "DOUBLE PRECISION":
		if d == "DOUBLE PRECISION" || d == "FLOAT8" {
			return true
		}
	case "TEXT", "STRING":
		if d == "TEXT" || d == "VARCHAR" || d == "CHARACTER VARYING" {
			return true
		}
	case "JSONB":
		if d == "JSONB" {
			return true
		}
	case "UUID":
		if d == "UUID" {
			return true
		}
	}

	// VARCHAR(N) vs character varying
	if strings.HasPrefix(p, "VARCHAR") || strings.HasPrefix(p, "CHAR") {
		if d == "CHARACTER VARYING" || strings.Contains(d, "CHARACTER VARYING") {
			return true
		}
	}

	// DECIMAL/NUMERIC
	if strings.HasPrefix(p, "DECIMAL") || strings.HasPrefix(p, "NUMERIC") {
		if d == "NUMERIC" || strings.HasPrefix(d, "NUMERIC") {
			return true
		}
	}
	// Arrays
	if strings.HasSuffix(p, "[]") && strings.HasPrefix(d, "ARRAY") {
		return true // Loose check for arrays for now
	}

	return false
}

// addColumn adds a single column to a table
func (umm *UnifiedMigrationManager) addColumn(schema, tableName string, c columnDef) error {
	fullTableName := umm.getFullTableName(schema, tableName)

	colDef := fmt.Sprintf("%s %s", umm.quoteIdentifier(c.name), c.typ)

	if c.notNull {
		colDef += " NOT NULL"
	}

	if c.defaultExp != "" && c.defaultExp != "NULL" {
		colDef += fmt.Sprintf(" DEFAULT %s", c.defaultExp)
	}

	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", fullTableName, colDef)

	if _, err := umm.db.Exec(query); err != nil {
		return err
	}

	return nil
}

// syncForeignKeys ensures foreign key constraints match Proto definition (Add, Update)
func (umm *UnifiedMigrationManager) syncForeignKeys(schema, tableName string, foreignKeys []ForeignKeyInfo) error {
	fullTableName := umm.getFullTableName(schema, tableName)

	// Check Snapshot availability
	umm.cacheMutex.RLock()
	snap, snapOk := umm.schemaSnapshots[schema]
	umm.cacheMutex.RUnlock()

	for _, fkInfo := range foreignKeys {
		fk := fkInfo.FK
		colName := fkInfo.ColumnName

		if colName == "" {
			continue
		}

		// Generate constraint name (Standardized)
		constraintName := fmt.Sprintf("fk_%s_%s", tableName, colName)

		// Parse reference table (may include schema)
		refTable := fk.ReferencesTable
		refSchema := fk.ReferencesSchema

		// If no schema specified, use the same schema as the current table
		if refSchema == "" {
			refSchema = schema
		}

		// Build fully qualified table name
		refTable = umm.getFullTableName(refSchema, refTable)

		refColumn := fk.ReferencesColumn
		if refColumn == "" {
			refColumn = "id"
		}

		// Map ReferentialAction enum to SQL string
		onDelete := "NO ACTION"
		switch fk.OnDelete {
		case commonv1.ReferentialAction_REFERENTIAL_ACTION_CASCADE:
			onDelete = "CASCADE"
		case commonv1.ReferentialAction_REFERENTIAL_ACTION_SET_NULL:
			onDelete = "SET NULL"
		case commonv1.ReferentialAction_REFERENTIAL_ACTION_SET_DEFAULT:
			onDelete = "SET DEFAULT"
		case commonv1.ReferentialAction_REFERENTIAL_ACTION_RESTRICT:
			onDelete = "RESTRICT"
		}

		onUpdate := "NO ACTION"
		switch fk.OnUpdate {
		case commonv1.ReferentialAction_REFERENTIAL_ACTION_CASCADE:
			onUpdate = "CASCADE"
		case commonv1.ReferentialAction_REFERENTIAL_ACTION_SET_NULL:
			onUpdate = "SET NULL"
		case commonv1.ReferentialAction_REFERENTIAL_ACTION_SET_DEFAULT:
			onUpdate = "SET DEFAULT"
		case commonv1.ReferentialAction_REFERENTIAL_ACTION_RESTRICT:
			onUpdate = "RESTRICT"
		}

		var dbDeleteRule, dbUpdateRule string
		var constraintExists bool

		if snapOk {
			// Fast Check via Snapshot
			if cd, found := snap.Constraints[constraintName]; found {
				constraintExists = true
				if cd.FK != nil {
					dbDeleteRule = cd.FK.DeleteRule
					dbUpdateRule = cd.FK.UpdateRule
				}
			}
		} else {
			// Slow Check via DB
			checkQuery := `
				SELECT 
					rc.delete_rule,
					rc.update_rule
				FROM information_schema.referential_constraints rc
				JOIN information_schema.table_constraints tc ON rc.constraint_name = tc.constraint_name
				WHERE tc.constraint_name = $1 AND tc.table_schema = $2 AND tc.table_name = $3
			`
			err := umm.db.QueryRow(checkQuery, constraintName, schema, tableName).Scan(&dbDeleteRule, &dbUpdateRule)
			if err == nil {
				constraintExists = true
			}
		}

		if constraintExists {
			// Compare rules (Postgres returns rules matching our string generation usually, but let's be careful)
			// Postgres: CASCADE, SET NULL, SET DEFAULT, RESTRICT, NO ACTION

			if dbDeleteRule == onDelete && dbUpdateRule == onUpdate {
				continue // Matches exactly
			}

			appLogger.Warnf("  ≠ FK Mismatch %s: (Proto: %s/%s, DB: %s/%s). Recreating...", constraintName, onDelete, onUpdate, dbDeleteRule, dbUpdateRule)

			// Drop existing constraint
			dropQuery := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", fullTableName, umm.quoteIdentifier(constraintName))
			if _, err := umm.db.Exec(dropQuery); err != nil {
				return fmt.Errorf("failed to drop mismatched constraint %s: %w", constraintName, err)
			}
		}

		// Build FK query
		query := fmt.Sprintf(
			"ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE %s ON UPDATE %s",
			fullTableName,
			umm.quoteIdentifier(constraintName), // fk_tablename_columnname
			umm.quoteIdentifier(colName),        // actual column name (e.g., user_id)
			refTable,                            // referenced table (with schema if needed)
			umm.quoteIdentifier(refColumn),      // referenced column (usually 'id')
			onDelete,
			onUpdate,
		)

		if !constraintExists {
			appLogger.Infof("  Executing FK: %s.%s -> %s", tableName, colName, refTable)
		}

		if _, err := umm.db.Exec(query); err != nil {
			// Log warning but don't fail - some FKs may reference tables that don't exist yet or are in wrong schema
			appLogger.Warnf("  Failed to add/update FK %s.%s -> %s: %v", tableName, colName, refTable, err)
			continue
		}

		if constraintExists {
			appLogger.Infof("  ✓ Updated FK: %s", constraintName)
		}
	}

	return nil
}

// addIndexes adds indexes to a table
func (umm *UnifiedMigrationManager) addIndexes(schema, tableName string, indexes []*commonv1.IndexOptions) error {
	fullTableName := umm.getFullTableName(schema, tableName)

	// Check Snapshot availability
	umm.cacheMutex.RLock()
	snap, snapOk := umm.schemaSnapshots[schema]
	umm.cacheMutex.RUnlock()

	for _, idx := range indexes {
		indexName := idx.IndexName
		if indexName == "" {
			continue // Skip if no index name specified
		}

		// Check if index already exists
		var exists bool

		if snapOk {
			// Fast Check via Snapshot
			_, exists = snap.Indexes[indexName]
		} else {
			// Slow Check via DB
			checkQuery := `
				SELECT EXISTS (
					SELECT 1 FROM pg_indexes 
					WHERE schemaname = $1 AND tablename = $2 AND indexname = $3
				)
			`
			if err := umm.db.QueryRow(checkQuery, schema, tableName, indexName).Scan(&exists); err != nil {
				return fmt.Errorf("failed to check index existence: %w", err)
			}
		}

		if exists {
			continue
		}

		unique := ""
		if idx.Unique {
			unique = "UNIQUE"
		}

		// Map IndexType enum to SQL method
		method := "BTREE"
		switch idx.IndexType {
		case commonv1.IndexType_INDEX_TYPE_BTREE:
			method = "BTREE"
		case commonv1.IndexType_INDEX_TYPE_HASH:
			method = "HASH"
		case commonv1.IndexType_INDEX_TYPE_GIN:
			method = "GIN"
		case commonv1.IndexType_INDEX_TYPE_GIST:
			method = "GIST"
		case commonv1.IndexType_INDEX_TYPE_BRIN:
			method = "BRIN"
		}

		// Override with custom method if specified
		if idx.IndexMethod != "" {
			method = idx.IndexMethod
		}

		// Build column list (single column or composite)
		var quotedCols []string
		if len(idx.CompositeFields) > 0 {
			quotedCols = make([]string, len(idx.CompositeFields))
			for i, col := range idx.CompositeFields {
				quotedCols[i] = umm.quoteIdentifier(col)
			}
		} else {
			// Single column index - this should be handled elsewhere
			continue
		}

		where := ""
		if idx.WhereClause != "" {
			where = fmt.Sprintf(" WHERE %s", idx.WhereClause)
		}

		query := fmt.Sprintf(
			"CREATE %s INDEX %s ON %s USING %s (%s)%s",
			unique,
			umm.quoteIdentifier(indexName),
			fullTableName,
			method,
			strings.Join(quotedCols, ", "),
			where,
		)

		if _, err := umm.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create index %s: %w", indexName, err)
		}
	}

	return nil
}

// addTableComment adds a comment to a table
func (umm *UnifiedMigrationManager) addTableComment(schema, tableName, comment string) error {
	if comment == "" {
		return nil
	}

	fullTableName := umm.getFullTableName(schema, tableName)
	query := fmt.Sprintf("COMMENT ON TABLE %s IS %s", fullTableName, umm.quoteLiteral(comment))

	if _, err := umm.db.Exec(query); err != nil {
		return err
	}

	return nil
}

// addColumnComments adds comments to columns
func (umm *UnifiedMigrationManager) addColumnComments(schema, tableName string, md protoreflect.MessageDescriptor) error {
	fullTableName := umm.getFullTableName(schema, tableName)
	fields := md.Fields()

	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		colName := umm.snakeCase(string(f.Name()))

		fOpts := f.Options()
		if fOpts != nil && proto.HasExtension(fOpts, commonv1.E_Column) {
			if colOpts, ok := proto.GetExtension(fOpts, commonv1.E_Column).(*commonv1.ColumnOptions); ok && colOpts != nil {
				if colOpts.Comment != "" {
					query := fmt.Sprintf(
						"COMMENT ON COLUMN %s.%s IS %s",
						fullTableName,
						umm.quoteIdentifier(colName),
						umm.quoteLiteral(colOpts.Comment),
					)

					if _, err := umm.db.Exec(query); err != nil {
						appLogger.Warnf("Failed to add comment to column %s: %v", colName, err)
					}
				}
			}
		}
	}

	return nil
}

// addCheckConstraints adds CHECK constraints from proto column options
func (umm *UnifiedMigrationManager) addCheckConstraints(schema, tableName string, cols []columnDef) error {
	fullTableName := umm.getFullTableName(schema, tableName)

	for _, c := range cols {
		if c.checkConstraint == "" {
			continue
		}

		constraintName := fmt.Sprintf("chk_%s_%s", tableName, c.name)

		// Check if constraint already exists
		var exists bool
		checkQuery := `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint pc
				JOIN pg_class rel ON rel.oid = pc.conrelid
				JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
				WHERE nsp.nspname = $1 AND rel.relname = $2 AND pc.conname = $3
			)`
		if err := umm.db.QueryRow(checkQuery, schema, tableName, constraintName).Scan(&exists); err != nil {
			appLogger.Warnf("Failed to check constraint existence for %s: %v", constraintName, err)
			continue
		}

		if exists {
			continue // Constraint already exists
		}

		// Add the CHECK constraint
		query := fmt.Sprintf(
			"ALTER TABLE %s ADD CONSTRAINT %s CHECK (%s)",
			fullTableName,
			umm.quoteIdentifier(constraintName),
			c.checkConstraint,
		)

		if _, err := umm.db.Exec(query); err != nil {
			appLogger.Warnf("Failed to add CHECK constraint %s: %v", constraintName, err)
		} else {
			appLogger.Infof("  ✓ Added CHECK constraint: %s", constraintName)
		}
	}

	return nil
}

// ==================== UTILITY METHODS ====================

// snakeCase converts CamelCase to snake_case
func (umm *UnifiedMigrationManager) snakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			if i+1 < len(s) && unicode.IsLower(rune(s[i+1])) {
				result = append(result, '_')
			} else if i > 0 && unicode.IsLower(rune(s[i-1])) {
				result = append(result, '_')
			}
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// pluralize converts singular noun to plural (simple English rules)
func (umm *UnifiedMigrationManager) pluralize(s string) string {
	if s == "" {
		return s
	}

	// Irregular plurals
	irregulars := map[string]string{
		"person": "people",
		"man":    "men",
		"woman":  "women",
		"child":  "children",
		"foot":   "feet",
		"tooth":  "teeth",
		"mouse":  "mice",
		"goose":  "geese",
	}

	if plural, ok := irregulars[s]; ok {
		return plural
	}

	// Already plural or ends with s
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "ss") {
		return s + "es"
	}

	// Ends with y preceded by consonant
	if len(s) >= 2 && s[len(s)-1] == 'y' {
		prev := rune(s[len(s)-2])
		if !strings.ContainsRune("aeiou", prev) {
			return s[:len(s)-1] + "ies"
		}
	}

	// Ends with specific patterns
	if strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") ||
		strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") {
		return s + "es"
	}

	if strings.HasSuffix(s, "fe") {
		return s[:len(s)-2] + "ves"
	}

	if strings.HasSuffix(s, "f") {
		return s[:len(s)-1] + "ves"
	}

	if strings.HasSuffix(s, "o") {
		return s + "es"
	}

	// Default: add 's'
	return s + "s"
}

// quoteLiteral quotes a string literal for SQL
func (umm *UnifiedMigrationManager) quoteLiteral(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "''"))
}
