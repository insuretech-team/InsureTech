package services

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/joho/godotenv"
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

// cleanupTestData removes all test data from tables
func cleanupTestData(ctx context.Context) error {
	// Delete in correct order due to foreign keys
	tables := []string{
		"webrtc.peer_sessions",
		"webrtc.room_sessions",
		"public.webrtc_tracks",
		"webrtc.peers",
		"webrtc.rooms",
	}

	for _, table := range tables {
		if err := testDB.WithContext(ctx).Exec(fmt.Sprintf("DELETE FROM %s WHERE true", table)).Error; err != nil {
			return fmt.Errorf("failed to cleanup %s: %w", table, err)
		}
	}

	return nil
}

// getTestDB returns the test database connection
func getTestDB() *gorm.DB {
	return testDB
}
