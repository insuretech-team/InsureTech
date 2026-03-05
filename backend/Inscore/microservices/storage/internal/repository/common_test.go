package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/ops/config"
	"gorm.io/gorm"
)

var testDB *gorm.DB

// TestMain sets up the test database connection
func TestMain(m *testing.M) {
	// Initialize database connection
	if err := setupTestDB(); err != nil {
		fmt.Printf("Failed to setup test database: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestDB()

	os.Exit(code)
}

// setupTestDB initializes the test database connection
func setupTestDB() error {
	// Load environment variables from .env file
	envPath, err := config.ResolvePath(".env")
	if err != nil {
		return fmt.Errorf("failed to resolve .env file: %w", err)
	}

	// Load .env file (ignore error if file doesn't exist - env vars might be set already)
	_ = godotenv.Load(envPath)

	// Use dynamic config resolver to find database.yaml
	configPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		return fmt.Errorf("failed to resolve config file: %w", err)
	}

	// Initialize database manager from config
	if err := db.InitializeManagerForService(configPath); err != nil {
		return fmt.Errorf("failed to initialize database manager: %w", err)
	}

	// Get database connection using GetDB function
	testDB = db.GetDB()
	if testDB == nil {
		return fmt.Errorf("failed to get database connection")
	}

	// Verify connection
	sqlDB, err := testDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create storage schema if it doesn't exist
	if err := testDB.Exec("CREATE SCHEMA IF NOT EXISTS storage_schema").Error; err != nil {
		return fmt.Errorf("failed to create storage_schema: %w", err)
	}

	// Create files table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS storage_schema.files (
		file_id UUID PRIMARY KEY,
		tenant_id UUID NOT NULL,
		filename VARCHAR(255) NOT NULL,
		content_type VARCHAR(100) NOT NULL,
		size_bytes BIGINT NOT NULL,
		storage_key TEXT NOT NULL UNIQUE,
		bucket VARCHAR(100) NOT NULL,
		url TEXT NOT NULL,
		cdn_url TEXT,
		file_type VARCHAR(50) NOT NULL,
		reference_id UUID,
		reference_type TEXT,
		is_public BOOLEAN DEFAULT false,
		expires_at TIMESTAMPTZ,
		uploaded_by UUID NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_files_tenant_type ON storage_schema.files(tenant_id, file_type);
	CREATE INDEX IF NOT EXISTS idx_files_reference ON storage_schema.files(reference_id);
	`

	if err := testDB.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("failed to create files table: %w", err)
	}

	fmt.Println("✅ Test database connection established")
	return nil
}

// cleanupTestDB closes the database connection
func cleanupTestDB() {
	if testDB != nil {
		sqlDB, _ := testDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

// cleanupTestData removes all test data from storage tables
func cleanupTestData(ctx context.Context, tenantID string) error {
	// Delete storage files for specific tenant
	tables := []string{
		"storage_schema.files",
	}

	for _, table := range tables {
		if err := testDB.WithContext(ctx).Exec(fmt.Sprintf("DELETE FROM %s WHERE tenant_id = $1", table), tenantID).Error; err != nil {
			return fmt.Errorf("failed to cleanup %s: %w", table, err)
		}
	}

	return nil
}

// getTestDB returns the test database connection
func getTestDB() *gorm.DB {
	return testDB
}
