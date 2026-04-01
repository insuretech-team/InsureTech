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
	_ "github.com/newage-saint/insuretech/backend/inscore/db"
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

func testNotificationDB(t *testing.T) *gorm.DB {
	t.Helper()

	testDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())

		if err := env.Load(); err != nil {
			logger.Warnf("Warning: couldn't load .env: %v", err)
		}

		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
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
		t.Skipf("skipping live DB test: failed to init DB: %v", testDBErr)
	}
	if testDB == nil {
		t.Skip("skipping live DB test: DB is nil")
	}

	return testDB
}

func genValidMobile() string {
	prefixes := []string{"13", "14", "15", "16", "17", "18", "19"}
	seed := time.Now().UnixNano()
	prefix := prefixes[seed%int64(len(prefixes))]
	suffix := seed % 100000000
	return fmt.Sprintf("+880%s%08d", prefix, suffix)
}

func cleanupNotificationUser(ctx context.Context, t *testing.T, dbConn *gorm.DB, userID string) {
	t.Helper()
	if dbConn == nil || strings.TrimSpace(userID) == "" {
		return
	}

	_ = dbConn.WithContext(ctx).Exec(`DELETE FROM notification_schema.push_device_tokens WHERE user_id = ?`, userID).Error
	_ = dbConn.WithContext(ctx).Exec(`DELETE FROM notification_schema.notifications WHERE recipient_id = ?`, userID).Error
	_ = dbConn.WithContext(ctx).Exec(`DELETE FROM authn_schema.users WHERE user_id = ?`, userID).Error
}

func cleanupNotificationTemplate(ctx context.Context, t *testing.T, dbConn *gorm.DB, templateID string) {
	t.Helper()
	if dbConn == nil || strings.TrimSpace(templateID) == "" {
		return
	}

	_ = dbConn.WithContext(ctx).Exec(`DELETE FROM notification_schema.notification_templates WHERE template_id = ?`, templateID).Error
}

func cleanupWebhookSubscription(ctx context.Context, t *testing.T, dbConn *gorm.DB, subscriptionID string) {
	t.Helper()
	if dbConn == nil || strings.TrimSpace(subscriptionID) == "" {
		return
	}

	_ = dbConn.WithContext(ctx).Exec(`DELETE FROM notification_schema.webhook_delivery_attempts WHERE subscription_id = ?`, subscriptionID).Error
	_ = dbConn.WithContext(ctx).Exec(`DELETE FROM notification_schema.webhook_subscriptions WHERE subscription_id = ?`, subscriptionID).Error
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

func isTextType(dt string) bool {
	switch strings.ToLower(dt) {
	case "character varying", "character", "text", "uuid":
		return true
	default:
		return false
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

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
		cols = append(cols, "status")
		if isNumericType(dt) {
			vals = append(vals, status)
		} else if isTextType(dt) {
			vals = append(vals, "USER_STATUS_ACTIVE")
		} else {
			cols = cols[:len(cols)-1]
		}
	}
	if dt, ok := columnDataType(t, dbConn, "authn_schema", "users", "user_type"); ok {
		cols = append(cols, "user_type")
		if isNumericType(dt) {
			vals = append(vals, 1)
		} else if isTextType(dt) {
			vals = append(vals, "USER_TYPE_B2C_CUSTOMER")
		} else {
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
