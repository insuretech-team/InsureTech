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
	partnerTestDBOnce sync.Once
	partnerTestDB     *gorm.DB
	partnerTestDBErr  error
)

// testPartnerDB returns a gorm DB for live-testing against the real Neon/DO Postgres.
func testPartnerDB(t *testing.T) *gorm.DB {
	t.Helper()

	partnerTestDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}
		partnerTestDBErr = db.InitializeManagerForService(configPath)
		if partnerTestDBErr != nil {
			return
		}
		partnerTestDB = db.GetDB()
	})

	if partnerTestDBErr != nil {
		t.Skipf("skipping live DB test: %v", partnerTestDBErr)
	}
	if partnerTestDB == nil {
		t.Skip("skipping live DB test: db is nil")
	}
	return partnerTestDB
}

func newPartnerLiveID(prefix string) string {
	return prefix + "_" + uuid.New().String()[:8]
}

func cleanupPartner(t *testing.T, db *gorm.DB, partnerID string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(
		`DELETE FROM partner_schema.commissions WHERE partner_id = ?`, partnerID,
	).Error
	_ = db.WithContext(context.Background()).Exec(
		`DELETE FROM partner_schema.agents WHERE partner_id = ?`, partnerID,
	).Error
	_ = db.WithContext(context.Background()).Exec(
		`DELETE FROM partner_schema.partners WHERE partner_id = ?`, partnerID,
	).Error
}

func cleanupCommission(t *testing.T, db *gorm.DB, commissionID string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(
		`DELETE FROM partner_schema.commissions WHERE commission_id = ?`, commissionID,
	).Error
}
