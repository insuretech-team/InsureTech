package repository

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/env"
	"gorm.io/gorm"
)

var (
	fraudTestDBOnce sync.Once
	fraudTestDB     *gorm.DB
	fraudTestDBErr  error
)

// testFraudDB returns a gorm DB for live-testing against the real Neon/DO Postgres.
func testFraudDB(t *testing.T) *gorm.DB {
	t.Helper()

	fraudTestDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}
		fraudTestDBErr = db.InitializeManagerForService(configPath)
		if fraudTestDBErr != nil {
			return
		}
		fraudTestDB = db.GetDB()
	})

	if fraudTestDBErr != nil {
		t.Skipf("skipping live DB test: %v", fraudTestDBErr)
	}
	if fraudTestDB == nil {
		t.Skip("skipping live DB test: db is nil")
	}
	return fraudTestDB
}

func newLiveID(prefix string) string {
	return uuid.NewString()
}

func cleanupFraudRule(t *testing.T, db *gorm.DB, ruleID string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(
		`DELETE FROM insurance_schema.fraud_rules WHERE fraud_rule_id = ?`, ruleID,
	).Error
}

func cleanupFraudAlert(t *testing.T, db *gorm.DB, alertID string) {
	t.Helper()
	// Cases reference alerts, so delete cases first.
	_ = db.WithContext(context.Background()).Exec(
		`DELETE FROM insurance_schema.fraud_cases WHERE fraud_alert_id = ?`, alertID,
	).Error
	_ = db.WithContext(context.Background()).Exec(
		`DELETE FROM insurance_schema.fraud_alerts WHERE alert_id = ?`, alertID,
	).Error
}

func cleanupFraudCase(t *testing.T, db *gorm.DB, caseID string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(
		`DELETE FROM insurance_schema.fraud_cases WHERE case_id = ?`, caseID,
	).Error
}

func liveUserID(t *testing.T, db *gorm.DB) string {
	t.Helper()

	var userID string
	err := db.WithContext(context.Background()).Raw(`SELECT user_id FROM authn_schema.users LIMIT 1`).Scan(&userID).Error
	if err != nil || userID == "" {
		t.Skipf("skipping live DB test: unable to resolve auth user id (%v)", err)
	}
	return userID
}
