package repository

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func genAuthzValidMobile() string {
	// Keep deterministic format for DB checks while making collisions practically impossible.
	// Build 9 digits from UUID hex characters to satisfy +8801XXXXXXXXX format.
	u := strings.ReplaceAll(uuid.NewString(), "-", "")
	digits := make([]byte, 0, 9)
	for i := 0; i < len(u) && len(digits) < 9; i++ {
		c := u[i]
		if c >= '0' && c <= '9' {
			digits = append(digits, c)
		}
	}
	for len(digits) < 9 {
		digits = append(digits, '0')
	}
	return "+8801" + string(digits[:9])
}

func insertAuthnUserMinimal(t *testing.T, db *gorm.DB, userID string) string {
	t.Helper()

	mobile := genAuthzValidMobile()
	err := db.WithContext(context.Background()).Exec(
		`INSERT INTO authn_schema.users
		   (user_id, mobile_number, password_hash, status, user_type, created_at, updated_at)
		 VALUES (?, ?, 'test-hash', 'USER_STATUS_ACTIVE', 'USER_TYPE_B2C_CUSTOMER', NOW(), NOW())`,
		userID, mobile,
	).Error
	require.NoError(t, err)
	return mobile
}

func cleanupAuthnUserByID(t *testing.T, db *gorm.DB, userID string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authz_schema.user_roles WHERE user_id = ?`, userID).Error
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authz_schema.access_decision_audits WHERE user_id = ?`, userID).Error
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authn_schema.sessions WHERE user_id = ?`, userID).Error
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authn_schema.otps WHERE user_id = ?`, userID).Error
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authn_schema.users WHERE user_id = ?`, userID).Error
}

func cleanupRoleByID(t *testing.T, db *gorm.DB, roleID string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authz_schema.user_roles WHERE role_id = ?`, roleID).Error
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authz_schema.roles WHERE role_id = ?`, roleID).Error
}

func cleanupPolicyByID(t *testing.T, db *gorm.DB, policyID string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authz_schema.policy_rules WHERE policy_id = ?`, policyID).Error
}

func cleanupAuditByID(t *testing.T, db *gorm.DB, auditID string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authz_schema.access_decision_audits WHERE audit_id = ?`, auditID).Error
}

func cleanupTokenConfigByKID(t *testing.T, db *gorm.DB, kid string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authz_schema.token_configs WHERE kid = ?`, kid).Error
}

func cleanupCasbinByDomainPrefix(t *testing.T, db *gorm.DB, domainPrefix string) {
	t.Helper()
	_ = db.WithContext(context.Background()).Exec(`DELETE FROM authz_schema.casbin_rules WHERE v1 LIKE ?`, domainPrefix+"%").Error
}

func newLiveID(prefix string) string {
	return prefix + "_" + uuid.New().String()[:8]
}
