package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	_ "github.com/newage-saint/insuretech/backend/inscore/db" // init() → proto_enum + proto_timestamp serializers
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/env"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	testDBOnce sync.Once
	testDB     *gorm.DB
	testDBErr  error
)

// testWorkflowDB returns a live *gorm.DB ready for workflow_schema tests.
// Mirrors authn's testAuthnDB() exactly — same database.yaml path convention.
func testWorkflowDB(t *testing.T) *gorm.DB {
	t.Helper()

	testDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())

		if err := env.Load(); err != nil {
			logger.Warnf("Warning: couldn't load .env: %v", err)
		}

		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			// Relative to: backend/inscore/microservices/workflow/internal/repository
			configPath = "../../../../database.yaml"
		}

		testDBErr = db.InitializeManagerForService(configPath)
		if testDBErr != nil {
			return
		}

		// Register proto_timestamp serializer (proto_enum is registered by db init())
		schema.RegisterSerializer("proto_timestamp", db.ProtoTimestampSerializer{})

		testDB = db.GetDB()
		if testDB != nil {
			testDB = testDB.Debug()
		}
	})

	if testDBErr != nil {
		t.Fatalf("failed to init test db: %v", testDBErr)
	}
	if testDB == nil {
		t.Fatalf("test db is nil")
	}
	return testDB
}

// ─── CLEANUP HELPERS ─────────────────────────────────────────────────────────

// cleanupWorkflowInstance deletes a workflow instance and all its tasks (CASCADE).
// Tasks are deleted by FK CASCADE, so only the instance row needs deleting.
func cleanupWorkflowInstance(ctx context.Context, t *testing.T, dbConn *gorm.DB, instanceID string) {
	t.Helper()
	if instanceID == "" {
		return
	}
	// Tasks will cascade-delete with the instance
	_ = dbConn.Table(tableInstances).
		Where("instance_id = ?", instanceID).
		Delete(map[string]any{}).Error
}

// cleanupWorkflowDefinition deletes a workflow definition by ID.
// Only call after all instances referencing it have been deleted.
func cleanupWorkflowDefinition(ctx context.Context, t *testing.T, dbConn *gorm.DB, definitionID string) {
	t.Helper()
	if definitionID == "" {
		return
	}
	_ = dbConn.Table(tableDefinitions).
		Where("workflow_definition_id = ?", definitionID).
		Delete(map[string]any{}).Error
}

// cleanupWorkflowDefinitionByName deletes a workflow definition by name.
func cleanupWorkflowDefinitionByName(ctx context.Context, t *testing.T, dbConn *gorm.DB, name string) {
	t.Helper()
	if name == "" {
		return
	}
	_ = dbConn.Table(tableDefinitions).
		Where("name = ?", name).
		Delete(map[string]any{}).Error
}

// ─── INSERT HELPERS ───────────────────────────────────────────────────────────

// insertTestUser inserts a minimal authn_schema.users row needed for FK constraints.
// workflow_instances.initiated_by and workflow_tasks.assigned_to reference authn_schema.users.
func insertTestUser(t *testing.T, dbConn *gorm.DB) string {
	t.Helper()
	userID := uuid.NewString()
	err := dbConn.Exec(
		`INSERT INTO authn_schema.users (user_id, mobile_number, password_hash, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT DO NOTHING`,
		userID,
		fmt.Sprintf("+88017%08d", time.Now().UnixNano()%100000000),
		"$2a$10$test_hash_not_real_bcrypt_placeholder_for_tests_only",
		"USER_STATUS_ACTIVE",
		time.Now().UTC(),
		time.Now().UTC(),
	).Error
	if err != nil {
		t.Fatalf("insertTestUser: %v", err)
	}
	return userID
}

// cleanupTestUser removes the test user and any workflow rows referencing it.
func cleanupTestUser(ctx context.Context, t *testing.T, dbConn *gorm.DB, userID string) {
	t.Helper()
	if userID == "" {
		return
	}
	// Tasks first (assigned_to FK SET NULL — no cascade, but best-effort)
	_ = dbConn.Table(tableTasks).
		Where("assigned_to = ?", userID).
		Updates(map[string]any{"assigned_to": nil}).Error
	// Instances (initiated_by FK RESTRICT — delete instances first)
	_ = dbConn.Table(tableInstances).
		Where("initiated_by = ?", userID).
		Delete(map[string]any{}).Error
	_ = dbConn.Table("authn_schema.users").
		Where("user_id = ?", userID).
		Delete(map[string]any{}).Error
}
