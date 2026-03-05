package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	_ "github.com/newage-saint/insuretech/backend/inscore/db" // ensure init() fires schema registry
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

// genValidMobile returns a digit-only Bangladeshi number in +8801XXXXXXXXX format.
// Live DB enforces chk_users_mobile_number.
func genValidMobile() string {
	n := time.Now().UnixNano() % 1_000_000_000 // 9 digits
	return fmt.Sprintf("+8801%09d", n)
}

func genValidNID() string {
	// 10 digits is commonly accepted; adjust if DB constraint changes.
	n := time.Now().UnixNano() % 10_000_000_000
	return fmt.Sprintf("%010d", n)
}

// testAuthnDB returns a gorm DB scoped to the authn schema.
// It uses the same approach as `test_all_authn.go`, but is test-friendly and reusable.
func testAuthnDB(t *testing.T) *gorm.DB {
	t.Helper()

	testDBOnce.Do(func() {
		// Ensure logger is initialized for tests.
		_ = logger.Initialize(logger.NoFileConfig())

		if err := env.Load(); err != nil {
			logger.Warnf("Warning: couldn't load .env: %v", err)
		}

		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			// Relative to this package: backend/inscore/microservices/authn/internal/repository
			configPath = "../../../../database.yaml"
		}

		testDBErr = db.InitializeManagerForService(configPath)
		if testDBErr != nil {
			return
		}

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

	// Always scope to schema tables explicitly.
	return testDB
}

func cleanupAuthnUser(ctx context.Context, t *testing.T, dbConn *gorm.DB, mobile string, userID string) {
	// NOTE: some live DBs may not match the latest proto-generated schema.
	// Cleanup should be best-effort and avoid relying on columns that may not exist.

	t.Helper()
	if dbConn == nil {
		return
	}

	// Delete in dependency order
	if userID != "" {
		_ = dbConn.Table("authn_schema.otps").Where("user_id = ?", userID).Delete(map[string]any{}).Error
		_ = dbConn.Table("authn_schema.sessions").Where("user_id = ?", userID).Delete(map[string]any{}).Error
		_ = dbConn.Table("authn_schema.users").Where("user_id = ?", userID).Delete(map[string]any{}).Error
		return
	}

	// Mobile-based cleanup is best-effort. If user schema differs, we skip it.
	if mobile != "" {
		// Attempt to delete sessions/users by mobile only if the column exists.
		var n int
		err := dbConn.Raw(`select count(1) from information_schema.columns where table_schema='authn_schema' and table_name='users' and column_name='mobile_number'`).Scan(&n).Error
		if err == nil && n > 0 {
			_ = dbConn.Exec(`delete from authn_schema.sessions where user_id in (select user_id from authn_schema.users where mobile_number = ?)`, mobile).Error
			_ = dbConn.Exec(`delete from authn_schema.otps where user_id in (select user_id from authn_schema.users where mobile_number = ?)`, mobile).Error
			_ = dbConn.Exec(`delete from authn_schema.users where mobile_number = ?`, mobile).Error
		}
	}
}

func columnExists(t *testing.T, dbConn *gorm.DB, schemaName, tableName, columnName string) bool {
	t.Helper()
	_, ok := columnDataType(t, dbConn, schemaName, tableName, columnName)
	return ok
}

func columnDataType(t *testing.T, dbConn *gorm.DB, schemaName, tableName, columnName string) (string, bool) {
	t.Helper()
	var dt string
	err := dbConn.Raw(
		`select data_type from information_schema.columns where table_schema=? and table_name=? and column_name=?`,
		schemaName, tableName, columnName,
	).Scan(&dt).Error
	requireNoError(t, err)
	if dt == "" {
		return "", false
	}
	return dt, true
}

func isNumericType(dt string) bool {
	switch strings.ToLower(dt) {
	case "smallint", "integer", "bigint", "numeric", "double precision", "real":
		return true
	default:
		return false
	}
}

func isBoolType(dt string) bool {
	return strings.ToLower(dt) == "boolean"
}

func isTextType(dt string) bool {
	switch strings.ToLower(dt) {
	case "character varying", "character", "text", "uuid":
		return true
	default:
		return false
	}
}

// legacy helper retained for compatibility
func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// insertUserMinimal inserts a user row using only columns that exist in the live DB.
func insertUserMinimal(t *testing.T, dbConn *gorm.DB, userID, mobile, email, passwordHash string, status int32) {
	t.Helper()

	cols := []string{"user_id"}
	vals := []any{userID}

	if mobile != "" && columnExists(t, dbConn, "authn_schema", "users", "mobile_number") {
		cols = append(cols, "mobile_number")
		vals = append(vals, mobile)
	}
	if email != "" && columnExists(t, dbConn, "authn_schema", "users", "email") {
		cols = append(cols, "email")
		vals = append(vals, email)
	}
	if passwordHash != "" && columnExists(t, dbConn, "authn_schema", "users", "password_hash") {
		cols = append(cols, "password_hash")
		vals = append(vals, passwordHash)
	}
	if dt, ok := columnDataType(t, dbConn, "authn_schema", "users", "status"); ok {
		// Live DB may represent status as INT or VARCHAR. Use a compatible value.
		cols = append(cols, "status")
		if isNumericType(dt) {
			vals = append(vals, status)
		} else if isTextType(dt) {
			// best-effort: store enum name
			vals = append(vals, "USER_STATUS_ACTIVE")
		} else {
			// unknown type: skip inserting status
			cols = cols[:len(cols)-1]
		}
	}
	if columnExists(t, dbConn, "authn_schema", "users", "created_at") {
		cols = append(cols, "created_at")
		vals = append(vals, time.Now())
	}
	if columnExists(t, dbConn, "authn_schema", "users", "updated_at") {
		cols = append(cols, "updated_at")
		vals = append(vals, time.Now())
	}

	placeholders := make([]string, 0, len(cols))
	for range cols {
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf(
		"insert into authn_schema.users (%s) values (%s)",
		strings.Join(cols, ","),
		strings.Join(placeholders, ","),
	)
	requireNoError(t, dbConn.Exec(query, vals...).Error)
}

// insertSessionMinimal inserts a session row using only columns that exist in the live DB.
// It also attempts to satisfy NOT NULL constraints by populating common required columns when present.
func insertSessionMinimal(t *testing.T, dbConn *gorm.DB, sessionID, userID, sessionType string, isActive bool, expiresAt time.Time) {
	t.Helper()

	cols := []string{"session_id", "user_id"}
	vals := []any{sessionID, userID}

	// session_type
	if dt, ok := columnDataType(t, dbConn, "authn_schema", "sessions", "session_type"); ok {
		cols = append(cols, "session_type")
		if isNumericType(dt) {
			// best-effort numeric mapping: 1
			vals = append(vals, 1)
		} else {
			normalized := strings.TrimPrefix(strings.ToUpper(strings.TrimSpace(sessionType)), "SESSION_TYPE_")
			if normalized == "" {
				normalized = "JWT"
			}
			vals = append(vals, normalized)
		}
	}

	// device_type is often NOT NULL
	if dt, ok := columnDataType(t, dbConn, "authn_schema", "sessions", "device_type"); ok {
		cols = append(cols, "device_type")
		if isNumericType(dt) {
			vals = append(vals, 1)
		} else {
			vals = append(vals, "WEB")
		}
	}

	// device_id is often NOT NULL
	if _, ok := columnDataType(t, dbConn, "authn_schema", "sessions", "device_id"); ok {
		cols = append(cols, "device_id")
		vals = append(vals, "test-device")
	}

	// ip_address/user_agent are sometimes NOT NULL
	if _, ok := columnDataType(t, dbConn, "authn_schema", "sessions", "ip_address"); ok {
		cols = append(cols, "ip_address")
		vals = append(vals, "127.0.0.1")
	}
	if _, ok := columnDataType(t, dbConn, "authn_schema", "sessions", "user_agent"); ok {
		cols = append(cols, "user_agent")
		vals = append(vals, "test")
	}

	if columnExists(t, dbConn, "authn_schema", "sessions", "is_active") {
		cols = append(cols, "is_active")
		vals = append(vals, isActive)
	}
	if columnExists(t, dbConn, "authn_schema", "sessions", "expires_at") {
		cols = append(cols, "expires_at")
		vals = append(vals, expiresAt)
	}
	if columnExists(t, dbConn, "authn_schema", "sessions", "created_at") {
		cols = append(cols, "created_at")
		vals = append(vals, time.Now())
	}
	if columnExists(t, dbConn, "authn_schema", "sessions", "last_activity_at") {
		cols = append(cols, "last_activity_at")
		vals = append(vals, time.Now())
	}

	placeholders := make([]string, 0, len(cols))
	for range cols {
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf(
		"insert into authn_schema.sessions (%s) values (%s)",
		strings.Join(cols, ","),
		strings.Join(placeholders, ","),
	)
	requireNoError(t, dbConn.Exec(query, vals...).Error)
}
