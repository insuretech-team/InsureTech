package main

import (
	"database/sql"
	"fmt"
	"strings"
	"text/tabwriter"
	"os"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"gorm.io/gorm"
)

// ColumnInfo represents detailed column information
type ColumnInfo struct {
	Name         string
	DataType     string
	IsNullable   string
	ColumnDefault sql.NullString
	CharMaxLength sql.NullInt64
	IsPrimaryKey bool
	IsUnique     bool
}

// ForeignKeyInfo represents foreign key constraint information
type ForeignKeyInfo struct {
	ConstraintName   string
	ColumnName       string
	ForeignTableName string
	ForeignColumnName string
	OnUpdate         string
	OnDelete         string
}

// IndexInfo represents index information
type IndexInfo struct {
	IndexName  string
	ColumnName string
	IsUnique   bool
	IndexType  string
}

// TableInfo represents comprehensive table information
type TableInfo struct {
	SchemaName    string
	TableName     string
	RowCount      int64
	TableSize     string
	IndexSize     string
	TotalSize     string
	Columns       []ColumnInfo
	ForeignKeys   []ForeignKeyInfo
	Indexes       []IndexInfo
	InheritsFrom  sql.NullString
	Description   sql.NullString
}

// handlePrintSchema prints detailed schema information
func handlePrintSchema(schemaName, targetDB string) {
	// Normalize and validate target
	targetDB = strings.ToLower(strings.TrimSpace(targetDB))
	if targetDB == "" {
		targetDB = "primary"
	}
	if targetDB != "primary" && targetDB != "backup" {
		fmt.Printf("⚠️  Invalid target '%s', using 'primary'\n\n", targetDB)
		targetDB = "primary"
	}
	
	targetGorm := getDBConnection(targetDB)
	if targetGorm == nil {
		fmt.Printf("❌ Error: Could not get connection to %s database\n", targetDB)
		fmt.Println("\n💡 Tip: Check your database configuration in database.yaml")
		return
	}

	sqlDB, err := targetGorm.DB()
	if err != nil {
		fmt.Printf("❌ Error getting SQL DB: %v\n", err)
		return
	}

	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                   SCHEMA INFORMATION - %s                    ║\n", strings.ToUpper(targetDB))
	fmt.Printf("╚═══════════════════════════════════════════════════════════════════════╝\n\n")

	// Get schemas to display
	var schemas []string
	var query string
	
	schemaName = strings.TrimSpace(schemaName)
	
	if schemaName != "" {
		// Check if schema exists
		var exists bool
		err := sqlDB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM information_schema.schemata 
				WHERE schema_name = $1
			)
		`, schemaName).Scan(&exists)
		
		if err != nil {
			fmt.Printf("❌ Error checking schema existence: %v\n", err)
			return
		}
		
		if !exists {
			fmt.Printf("❌ Schema '%s' does not exist\n\n", schemaName)
			fmt.Println("💡 Available schemas:")
			
			// Show available schemas
			availQuery := `
				SELECT schema_name 
				FROM information_schema.schemata 
				WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
				ORDER BY schema_name
			`
			rows, err := sqlDB.Query(availQuery)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var s string
					if rows.Scan(&s) == nil {
						fmt.Printf("   • %s\n", s)
					}
				}
			}
			return
		}
		
		schemas = []string{schemaName}
	} else {
		// Get all user schemas (excluding system schemas)
		query = `
			SELECT schema_name 
			FROM information_schema.schemata 
			WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
			ORDER BY schema_name
		`
		rows, err := sqlDB.Query(query)
		if err != nil {
			fmt.Printf("❌ Error querying schemas: %v\n", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var schema string
			if err := rows.Scan(&schema); err != nil {
				fmt.Printf("⚠️  Error scanning schema: %v\n", err)
				continue
			}
			schemas = append(schemas, schema)
		}
		
		if len(schemas) == 0 {
			fmt.Println("⚠️  No user schemas found in database")
			return
		}
	}

	// Print each schema
	for _, schema := range schemas {
		printSchemaDetails(sqlDB, schema)
	}
}

// printSchemaDetails prints detailed information for a single schema
func printSchemaDetails(sqlDB *sql.DB, schemaName string) {
	fmt.Printf("📁 SCHEMA: %s\n", schemaName)
	fmt.Printf("═══════════════════════════════════════════════════════════════════════\n\n")

	// Get tables in schema
	query := `
		SELECT 
			t.table_name,
			COALESCE(pg_total_relation_size(quote_ident(t.table_schema)||'.'||quote_ident(t.table_name))::text, '0') as total_size,
			obj_description((quote_ident(t.table_schema)||'.'||quote_ident(t.table_name))::regclass, 'pg_class') as description
		FROM information_schema.tables t
		WHERE t.table_schema = $1
		AND t.table_type = 'BASE TABLE'
		ORDER BY t.table_name
	`

	rows, err := sqlDB.Query(query, schemaName)
	if err != nil {
		fmt.Printf("Error querying tables: %v\n", err)
		return
	}
	defer rows.Close()

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TABLE\tSIZE\tDESCRIPTION")
	fmt.Fprintln(tw, "─────\t────\t───────────")

	tableCount := 0
	for rows.Next() {
		var tableName string
		var totalSize string
		var description sql.NullString

		if err := rows.Scan(&tableName, &totalSize, &description); err != nil {
			fmt.Printf("Error scanning table: %v\n", err)
			continue
		}

		desc := ""
		if description.Valid {
			desc = description.String
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\n", tableName, formatSize(totalSize), desc)
		tableCount++
	}

	tw.Flush()
	fmt.Printf("\nTotal tables: %d\n\n", tableCount)
}

// handlePrintTable prints detailed information about a specific table or all tables in a schema
func handlePrintTable(tableName, targetDB string) {
	// Normalize and validate target
	targetDB = strings.ToLower(strings.TrimSpace(targetDB))
	if targetDB == "" {
		targetDB = "primary"
	}
	if targetDB != "primary" && targetDB != "backup" {
		fmt.Printf("⚠️  Invalid target '%s', using 'primary'\n\n", targetDB)
		targetDB = "primary"
	}
	
	targetGorm := getDBConnection(targetDB)
	if targetGorm == nil {
		fmt.Printf("❌ Error: Could not get connection to %s database\n", targetDB)
		fmt.Println("\n💡 Tip: Check your database configuration in database.yaml")
		return
	}

	sqlDB, err := targetGorm.DB()
	if err != nil {
		fmt.Printf("❌ Error getting SQL DB: %v\n", err)
		return
	}

	// If no table name provided, this is an error
	if tableName == "" {
		fmt.Println("❌ Error: Table name is required for print-table command")
		fmt.Println("\n📖 Usage Examples:")
		fmt.Println("   dbmanager print-table --table=users")
		fmt.Println("   dbmanager print-table --table=auth.users")
		fmt.Println("   dbmanager print-table --table=auth.customers --target=backup")
		fmt.Println("\n💡 Tip: Use 'dbmanager print-schema' to see available tables")
		return
	}

	// Parse schema and table name
	var schemaName, tblName string
	if strings.Contains(tableName, ".") {
		parts := strings.SplitN(tableName, ".", 2)
		schemaName = parts[0]
		tblName = parts[1]
	} else {
		schemaName = "public"
		tblName = tableName
	}

	// Get table information
	tableInfo, err := getTableInfo(sqlDB, schemaName, tblName)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		fmt.Println("\n💡 Available options:")
		fmt.Println("   • Use 'dbmanager print-schema' to list all tables")
		fmt.Println("   • Use 'dbmanager print-schema --schema=" + schemaName + "' to list tables in " + schemaName)
		fmt.Println("   • Check table name spelling and schema")
		return
	}

	printTableInfo(tableInfo)
}

// getTableInfo retrieves comprehensive table information
func getTableInfo(sqlDB *sql.DB, schemaName, tableName string) (*TableInfo, error) {
	info := &TableInfo{
		SchemaName: schemaName,
		TableName:  tableName,
	}

	// Check if table exists
	var exists bool
	err := sqlDB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = $1 AND table_name = $2
		)
	`, schemaName, tableName).Scan(&exists)
	
	if err != nil {
		return nil, err
	}
	
	if !exists {
		return nil, fmt.Errorf("table %s.%s does not exist", schemaName, tableName)
	}

	// Get row count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", 
		quoteIdentifier(schemaName), quoteIdentifier(tableName))
	err = sqlDB.QueryRow(countQuery).Scan(&info.RowCount)
	if err != nil {
		// Non-fatal, some tables might not allow count
		info.RowCount = -1
	}

	// Get table sizes
	sizeQuery := `
		SELECT
			pg_size_pretty(pg_table_size(c.oid)) as table_size,
			pg_size_pretty(pg_indexes_size(c.oid)) as index_size,
			pg_size_pretty(pg_total_relation_size(c.oid)) as total_size,
			obj_description(c.oid, 'pg_class') as description
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1 AND c.relname = $2
	`
	err = sqlDB.QueryRow(sizeQuery, schemaName, tableName).Scan(
		&info.TableSize, &info.IndexSize, &info.TotalSize, &info.Description)
	if err != nil {
		// Set defaults if query fails
		info.TableSize = "N/A"
		info.IndexSize = "N/A"
		info.TotalSize = "N/A"
	}

	// Get inheritance info
	inheritQuery := `
		SELECT p.relname
		FROM pg_inherits
		JOIN pg_class c ON c.oid = pg_inherits.inhrelid
		JOIN pg_class p ON p.oid = pg_inherits.inhparent
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1 AND c.relname = $2
	`
	sqlDB.QueryRow(inheritQuery, schemaName, tableName).Scan(&info.InheritsFrom)

	// Get columns
	info.Columns, err = getDetailedTableColumns(sqlDB, schemaName, tableName)
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %v", err)
	}

	// Get foreign keys
	info.ForeignKeys, err = getTableForeignKeys(sqlDB, schemaName, tableName)
	if err != nil {
		return nil, fmt.Errorf("error getting foreign keys: %v", err)
	}

	// Get indexes
	info.Indexes, err = getTableIndexes(sqlDB, schemaName, tableName)
	if err != nil {
		return nil, fmt.Errorf("error getting indexes: %v", err)
	}

	return info, nil
}

// getDetailedTableColumns retrieves column information
func getDetailedTableColumns(sqlDB *sql.DB, schemaName, tableName string) ([]ColumnInfo, error) {
	query := `
		SELECT 
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			c.character_maximum_length,
			COALESCE(
				(SELECT true 
				 FROM information_schema.table_constraints tc
				 JOIN information_schema.key_column_usage kcu 
				   ON tc.constraint_name = kcu.constraint_name
				 WHERE tc.table_schema = c.table_schema
				   AND tc.table_name = c.table_name
				   AND kcu.column_name = c.column_name
				   AND tc.constraint_type = 'PRIMARY KEY'
				 LIMIT 1), false
			) as is_primary_key,
			COALESCE(
				(SELECT true 
				 FROM information_schema.table_constraints tc
				 JOIN information_schema.key_column_usage kcu 
				   ON tc.constraint_name = kcu.constraint_name
				 WHERE tc.table_schema = c.table_schema
				   AND tc.table_name = c.table_name
				   AND kcu.column_name = c.column_name
				   AND tc.constraint_type = 'UNIQUE'
				 LIMIT 1), false
			) as is_unique
		FROM information_schema.columns c
		WHERE c.table_schema = $1 AND c.table_name = $2
		ORDER BY c.ordinal_position
	`

	rows, err := sqlDB.Query(query, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		err := rows.Scan(
			&col.Name,
			&col.DataType,
			&col.IsNullable,
			&col.ColumnDefault,
			&col.CharMaxLength,
			&col.IsPrimaryKey,
			&col.IsUnique,
		)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

// getTableForeignKeys retrieves foreign key information
func getTableForeignKeys(sqlDB *sql.DB, schemaName, tableName string) ([]ForeignKeyInfo, error) {
	query := `
		SELECT
			tc.constraint_name,
			kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			rc.update_rule,
			rc.delete_rule
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
		  ON tc.constraint_name = kcu.constraint_name
		  AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
		  ON ccu.constraint_name = tc.constraint_name
		  AND ccu.table_schema = tc.table_schema
		JOIN information_schema.referential_constraints AS rc
		  ON rc.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
		  AND tc.table_schema = $1
		  AND tc.table_name = $2
		ORDER BY tc.constraint_name, kcu.ordinal_position
	`

	rows, err := sqlDB.Query(query, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fks []ForeignKeyInfo
	for rows.Next() {
		var fk ForeignKeyInfo
		err := rows.Scan(
			&fk.ConstraintName,
			&fk.ColumnName,
			&fk.ForeignTableName,
			&fk.ForeignColumnName,
			&fk.OnUpdate,
			&fk.OnDelete,
		)
		if err != nil {
			return nil, err
		}
		fks = append(fks, fk)
	}

	return fks, nil
}

// getTableIndexes retrieves index information
func getTableIndexes(sqlDB *sql.DB, schemaName, tableName string) ([]IndexInfo, error) {
	query := `
		SELECT
			i.relname as index_name,
			a.attname as column_name,
			ix.indisunique as is_unique,
			am.amname as index_type
		FROM pg_class t
		JOIN pg_namespace n ON n.oid = t.relnamespace
		JOIN pg_index ix ON t.oid = ix.indrelid
		JOIN pg_class i ON i.oid = ix.indexrelid
		JOIN pg_am am ON i.relam = am.oid
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
		WHERE n.nspname = $1 AND t.relname = $2
		ORDER BY i.relname, a.attnum
	`

	rows, err := sqlDB.Query(query, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []IndexInfo
	for rows.Next() {
		var idx IndexInfo
		err := rows.Scan(
			&idx.IndexName,
			&idx.ColumnName,
			&idx.IsUnique,
			&idx.IndexType,
		)
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, idx)
	}

	return indexes, nil
}

// printTableInfo prints formatted table information
func printTableInfo(info *TableInfo) {
	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                      TABLE INFORMATION                                ║\n")
	fmt.Printf("╚═══════════════════════════════════════════════════════════════════════╝\n\n")

	// Table metadata
	fmt.Printf("📋 Table: %s.%s\n", info.SchemaName, info.TableName)
	if info.Description.Valid && info.Description.String != "" {
		fmt.Printf("   Description: %s\n", info.Description.String)
	}
	if info.InheritsFrom.Valid && info.InheritsFrom.String != "" {
		fmt.Printf("   🔗 Inherits from: %s\n", info.InheritsFrom.String)
	}
	fmt.Println()

	// Statistics
	fmt.Printf("📊 Statistics:\n")
	if info.RowCount >= 0 {
		fmt.Printf("   Rows: %d\n", info.RowCount)
	} else {
		fmt.Printf("   Rows: N/A (cannot count)\n")
	}
	fmt.Printf("   Table Size: %s\n", info.TableSize)
	fmt.Printf("   Index Size: %s\n", info.IndexSize)
	fmt.Printf("   Total Size: %s\n", info.TotalSize)
	fmt.Println()

	// Columns
	fmt.Printf("🗂️  Columns (%d):\n", len(info.Columns))
	fmt.Printf("───────────────────────────────────────────────────────────────────────\n")
	
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tTYPE\tNULL\tDEFAULT\tKEY")
	fmt.Fprintln(tw, "────\t────\t────\t───────\t───")
	
	for _, col := range info.Columns {
		typeStr := col.DataType
		if col.CharMaxLength.Valid && col.CharMaxLength.Int64 > 0 {
			typeStr = fmt.Sprintf("%s(%d)", col.DataType, col.CharMaxLength.Int64)
		}
		
		nullStr := "YES"
		if col.IsNullable == "NO" {
			nullStr = "NO"
		}
		
		defaultStr := ""
		if col.ColumnDefault.Valid {
			defaultStr = col.ColumnDefault.String
			if len(defaultStr) > 30 {
				defaultStr = defaultStr[:27] + "..."
			}
		}
		
		keyStr := ""
		if col.IsPrimaryKey {
			keyStr = "PK"
		} else if col.IsUnique {
			keyStr = "UQ"
		}
		
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", 
			col.Name, typeStr, nullStr, defaultStr, keyStr)
	}
	tw.Flush()
	fmt.Println()

	// Foreign Keys
	if len(info.ForeignKeys) > 0 {
		fmt.Printf("🔗 Foreign Keys (%d):\n", len(info.ForeignKeys))
		fmt.Printf("───────────────────────────────────────────────────────────────────────\n")
		
		for _, fk := range info.ForeignKeys {
			fmt.Printf("   %s (%s) → %s (%s)\n", 
				fk.ColumnName, fk.ConstraintName,
				fk.ForeignTableName, fk.ForeignColumnName)
			fmt.Printf("      ON UPDATE %s, ON DELETE %s\n", fk.OnUpdate, fk.OnDelete)
		}
		fmt.Println()
	}

	// Indexes
	if len(info.Indexes) > 0 {
		fmt.Printf("📇 Indexes (%d):\n", len(info.Indexes))
		fmt.Printf("───────────────────────────────────────────────────────────────────────\n")
		
		// Group indexes by name
		indexMap := make(map[string][]string)
		indexUnique := make(map[string]bool)
		indexType := make(map[string]string)
		
		for _, idx := range info.Indexes {
			indexMap[idx.IndexName] = append(indexMap[idx.IndexName], idx.ColumnName)
			indexUnique[idx.IndexName] = idx.IsUnique
			indexType[idx.IndexName] = idx.IndexType
		}
		
		for idxName, cols := range indexMap {
			uniqueStr := ""
			if indexUnique[idxName] {
				uniqueStr = " UNIQUE"
			}
			fmt.Printf("   %s%s (%s) on (%s)\n", 
				idxName, uniqueStr, indexType[idxName], strings.Join(cols, ", "))
		}
		fmt.Println()
	}
}

// handlePrintTables prints detailed information for all tables in a schema or all schemas
func handlePrintTables(schemaName, targetDB string) {
	// Normalize and validate target
	targetDB = strings.ToLower(strings.TrimSpace(targetDB))
	if targetDB == "" {
		targetDB = "primary"
	}
	if targetDB != "primary" && targetDB != "backup" {
		fmt.Printf("⚠️  Invalid target '%s', using 'primary'\n\n", targetDB)
		targetDB = "primary"
	}
	
	targetGorm := getDBConnection(targetDB)
	if targetGorm == nil {
		fmt.Printf("❌ Error: Could not get connection to %s database\n", targetDB)
		fmt.Println("\n💡 Tip: Check your database configuration in database.yaml")
		return
	}

	sqlDB, err := targetGorm.DB()
	if err != nil {
		fmt.Printf("❌ Error getting SQL DB: %v\n", err)
		return
	}

	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║              DETAILED TABLE INFORMATION - %s                  ║\n", strings.ToUpper(targetDB))
	fmt.Printf("╚═══════════════════════════════════════════════════════════════════════╝\n")

	// Get schemas to process
	var schemas []string
	schemaName = strings.TrimSpace(schemaName)
	
	if schemaName != "" {
		// Check if schema exists
		var exists bool
		err := sqlDB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM information_schema.schemata 
				WHERE schema_name = $1
			)
		`, schemaName).Scan(&exists)
		
		if err != nil {
			fmt.Printf("❌ Error checking schema existence: %v\n", err)
			return
		}
		
		if !exists {
			fmt.Printf("❌ Schema '%s' does not exist\n\n", schemaName)
			fmt.Println("💡 Available schemas:")
			
			availQuery := `
				SELECT schema_name 
				FROM information_schema.schemata 
				WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
				ORDER BY schema_name
			`
			rows, err := sqlDB.Query(availQuery)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var s string
					if rows.Scan(&s) == nil {
						fmt.Printf("   • %s\n", s)
					}
				}
			}
			return
		}
		
		schemas = []string{schemaName}
	} else {
		// Get all user schemas
		query := `
			SELECT schema_name 
			FROM information_schema.schemata 
			WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
			ORDER BY schema_name
		`
		rows, err := sqlDB.Query(query)
		if err != nil {
			fmt.Printf("❌ Error querying schemas: %v\n", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var schema string
			if err := rows.Scan(&schema); err != nil {
				continue
			}
			schemas = append(schemas, schema)
		}
		
		if len(schemas) == 0 {
			fmt.Println("⚠️  No user schemas found in database")
			return
		}
	}

	// Process each schema
	totalTables := 0
	for _, schema := range schemas {
		// Get all tables in schema
		query := `
			SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = $1 AND table_type = 'BASE TABLE'
			ORDER BY table_name
		`
		rows, err := sqlDB.Query(query, schema)
		if err != nil {
			fmt.Printf("⚠️  Error querying tables in schema %s: %v\n", schema, err)
			continue
		}

		var tables []string
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				continue
			}
			tables = append(tables, table)
		}
		rows.Close()

		if len(tables) == 0 {
			fmt.Printf("\n📁 Schema: %s (no tables)\n", schema)
			continue
		}

		fmt.Printf("\n📁 Schema: %s (%d tables)\n", schema, len(tables))
		fmt.Printf("═══════════════════════════════════════════════════════════════════════\n")

		// Print each table
		for i, table := range tables {
			tableInfo, err := getTableInfo(sqlDB, schema, table)
			if err != nil {
				fmt.Printf("⚠️  Skipping %s.%s: %v\n", schema, table, err)
				continue
			}

			printTableInfo(tableInfo)
			totalTables++

			// Add separator between tables
			if i < len(tables)-1 {
				fmt.Printf("\n───────────────────────────────────────────────────────────────────────\n")
			}
		}
	}

	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  Total tables processed: %-48d ║\n", totalTables)
	fmt.Printf("╚═══════════════════════════════════════════════════════════════════════╝\n")
}

// handlePrintAll prints all schemas and tables
func handlePrintAll(targetDB string) {
	// Normalize and validate target
	targetDB = strings.ToLower(strings.TrimSpace(targetDB))
	if targetDB == "" {
		targetDB = "primary"
	}
	if targetDB != "primary" && targetDB != "backup" {
		fmt.Printf("⚠️  Invalid target '%s', using 'primary'\n\n", targetDB)
		targetDB = "primary"
	}
	
	targetGorm := getDBConnection(targetDB)
	if targetGorm == nil {
		fmt.Printf("❌ Error: Could not get connection to %s database\n", targetDB)
		fmt.Println("\n💡 Tip: Check your database configuration in database.yaml")
		return
	}

	sqlDB, err := targetGorm.DB()
	if err != nil {
		fmt.Printf("❌ Error getting SQL DB: %v\n", err)
		return
	}

	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║              DATABASE OVERVIEW - %s                       ║\n", strings.ToUpper(targetDB))
	fmt.Printf("╚═══════════════════════════════════════════════════════════════════════╝\n\n")

	// Get database size
	var dbSize string
	err = sqlDB.QueryRow("SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&dbSize)
	if err == nil {
		fmt.Printf("💾 Database Size: %s\n\n", dbSize)
	}

	// Get all schemas
	query := `
		SELECT schema_name 
		FROM information_schema.schemata 
		WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
		ORDER BY schema_name
	`
	rows, err := sqlDB.Query(query)
	if err != nil {
		fmt.Printf("Error querying schemas: %v\n", err)
		return
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			continue
		}
		schemas = append(schemas, schema)
	}

	// Print summary
	totalTables := 0
	for _, schema := range schemas {
		var count int
		countQuery := `
			SELECT COUNT(*) 
			FROM information_schema.tables 
			WHERE table_schema = $1 AND table_type = 'BASE TABLE'
		`
		sqlDB.QueryRow(countQuery, schema).Scan(&count)
		totalTables += count
	}

	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   Schemas: %d\n", len(schemas))
	fmt.Printf("   Tables: %d\n\n", totalTables)

	// Print each schema
	for _, schema := range schemas {
		printSchemaDetails(sqlDB, schema)
	}
}

// Helper functions

func getDBConnection(targetDB string) *gorm.DB {
	if db.Manager == nil {
		return nil
	}

	switch targetDB {
	case "primary":
		return db.Manager.GetPrimaryDB()
	case "backup":
		return db.Manager.GetBackupDB()
	default:
		return db.Manager.GetPrimaryDB()
	}
}

func quoteIdentifier(name string) string {
	return fmt.Sprintf("\"%s\"", strings.ReplaceAll(name, "\"", "\"\""))
}

func formatSize(bytesStr string) string {
	// If already formatted (contains letters), return as-is
	if strings.ContainsAny(bytesStr, "kMGTPEZY") {
		return bytesStr
	}
	
	// Otherwise it should be a number - parse and format
	var bytes int64
	fmt.Sscanf(bytesStr, "%d", &bytes)
	
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

// handlePrintTableData prints the actual data/entries from a table
func handlePrintTableData(tableName, targetDB string, limit int) {
	// Normalize and validate target
	targetDB = strings.ToLower(strings.TrimSpace(targetDB))
	if targetDB == "" {
		targetDB = "primary"
	}
	if targetDB != "primary" && targetDB != "backup" {
		fmt.Printf("⚠️  Invalid target '%s', using 'primary'\n\n", targetDB)
		targetDB = "primary"
	}
	
	targetGorm := getDBConnection(targetDB)
	if targetGorm == nil {
		fmt.Printf("❌ Error: Could not get connection to %s database\n", targetDB)
		fmt.Println("\n💡 Tip: Check your database configuration in database.yaml")
		return
	}

	sqlDB, err := targetGorm.DB()
	if err != nil {
		fmt.Printf("❌ Error getting SQL DB: %v\n", err)
		return
	}

	// If no table name provided, this is an error
	if tableName == "" {
		fmt.Println("❌ Error: Table name is required for print-table-data command")
		fmt.Println("\n📖 Usage Examples:")
		fmt.Println("   dbmanager print-table-data --table=users")
		fmt.Println("   dbmanager print-table-data --table=auth.users")
		fmt.Println("   dbmanager print-table-data --table=auth.users --limit=50")
		fmt.Println("   dbmanager print-table-data --table=auth.customers --target=backup")
		fmt.Println("\n💡 Tip: Use 'dbmanager print-schema' to see available tables")
		return
	}

	// Parse schema and table name
	var schemaName, tblName string
	if strings.Contains(tableName, ".") {
		parts := strings.SplitN(tableName, ".", 2)
		schemaName = parts[0]
		tblName = parts[1]
	} else {
		schemaName = "public"
		tblName = tableName
	}

	// Check if table exists
	var exists bool
	err = sqlDB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = $1 AND table_name = $2
		)
	`, schemaName, tblName).Scan(&exists)
	
	if err != nil {
		fmt.Printf("❌ Error checking table existence: %v\n", err)
		return
	}
	
	if !exists {
		fmt.Printf("❌ Error: Table %s.%s does not exist\n\n", schemaName, tblName)
		fmt.Println("💡 Available options:")
		fmt.Println("   • Use 'dbmanager print-schema' to list all tables")
		fmt.Println("   • Use 'dbmanager print-schema --schema=" + schemaName + "' to list tables in " + schemaName)
		fmt.Println("   • Check table name spelling and schema")
		return
	}

	// Get column information
	columns, err := getDetailedTableColumns(sqlDB, schemaName, tblName)
	if err != nil {
		fmt.Printf("❌ Error getting columns: %v\n", err)
		return
	}

	if len(columns) == 0 {
		fmt.Println("⚠️  Table has no columns")
		return
	}

	// Get row count
	var rowCount int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", 
		quoteIdentifier(schemaName), quoteIdentifier(tblName))
	err = sqlDB.QueryRow(countQuery).Scan(&rowCount)
	if err != nil {
		fmt.Printf("❌ Error counting rows: %v\n", err)
		return
	}

	// Print header
	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                      TABLE DATA - %s                    ║\n", strings.ToUpper(targetDB))
	fmt.Printf("╚═══════════════════════════════════════════════════════════════════════╝\n\n")
	fmt.Printf("📋 Table: %s.%s\n", schemaName, tblName)
	fmt.Printf("📊 Total Rows: %d\n", rowCount)
	
	if rowCount == 0 {
		fmt.Println("\n⚠️  Table is empty (no data)")
		return
	}

	// Apply limit
	displayLimit := limit
	if displayLimit <= 0 || displayLimit > 1000 {
		displayLimit = 100 // Default limit
	}
	
	if rowCount > int64(displayLimit) {
		fmt.Printf("📄 Showing: First %d rows (use --limit to adjust)\n\n", displayLimit)
	} else {
		fmt.Printf("📄 Showing: All %d rows\n\n", rowCount)
	}

	// Build column names for query
	var columnNames []string
	for _, col := range columns {
		columnNames = append(columnNames, quoteIdentifier(col.Name))
	}

	// Query data
	dataQuery := fmt.Sprintf("SELECT %s FROM %s.%s LIMIT %d",
		strings.Join(columnNames, ", "),
		quoteIdentifier(schemaName),
		quoteIdentifier(tblName),
		displayLimit)

	rows, err := sqlDB.Query(dataQuery)
	if err != nil {
		fmt.Printf("❌ Error querying data: %v\n", err)
		return
	}
	defer rows.Close()

	// Prepare scan destinations
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Print table header
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	// Print column names
	headerParts := make([]string, len(columns))
	for i, col := range columns {
		headerParts[i] = col.Name
	}
	fmt.Fprintln(tw, strings.Join(headerParts, "\t"))
	
	// Print separator
	separatorParts := make([]string, len(columns))
	for i := range columns {
		separatorParts[i] = strings.Repeat("─", len(columns[i].Name))
	}
	fmt.Fprintln(tw, strings.Join(separatorParts, "\t"))

	// Print data rows
	rowNum := 0
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			fmt.Printf("\n⚠️  Error scanning row %d: %v\n", rowNum+1, err)
			continue
		}

		rowData := make([]string, len(columns))
		for i, val := range values {
			rowData[i] = formatValue(val)
		}
		
		fmt.Fprintln(tw, strings.Join(rowData, "\t"))
		rowNum++
	}

	tw.Flush()

	if err := rows.Err(); err != nil {
		fmt.Printf("\n⚠️  Error during row iteration: %v\n", err)
	}

	fmt.Printf("\n✅ Displayed %d rows\n", rowNum)
	
	if rowCount > int64(displayLimit) {
		fmt.Printf("\n💡 Tip: Use --limit=%d to see more rows (max 1000)\n", min(int(rowCount), 1000))
	}
}

// formatValue formats a database value for display
func formatValue(val interface{}) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case []byte:
		str := string(v)
		if len(str) > 50 {
			return str[:47] + "..."
		}
		return str
	case string:
		if len(v) > 50 {
			return v[:47] + "..."
		}
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		str := fmt.Sprintf("%v", v)
		if len(str) > 50 {
			return str[:47] + "..."
		}
		return str
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

