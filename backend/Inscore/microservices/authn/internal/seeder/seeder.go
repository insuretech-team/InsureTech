package seeder

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/sms"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
)

// SeedAdminUser bootstraps a default SYSTEM_USER account on first deploy.
//
// Env vars:
// - ADMIN_EMAIL
// - ADMIN_MOBILE (Bangladesh format, e.g. +8801XXXXXXXXX)
// - ADMIN_PASSWORD
// - ADMIN_PASSWARD (legacy typo fallback; supported for backward compatibility)
//
// Behavior:
// - If any env var is missing, seeding is skipped.
// - If the user exists by email, we update password + status + email_verified.
// - If it does not exist, we create it.
func SeedAdminUser(ctx context.Context, db *gorm.DB) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminMobileRaw := os.Getenv("ADMIN_MOBILE")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		legacyPassword := os.Getenv("ADMIN_PASSWARD")
		if legacyPassword != "" {
			adminPassword = legacyPassword
			appLogger.Warn("Admin seeder: using legacy env ADMIN_PASSWARD; prefer ADMIN_PASSWORD")
		}
	}

	if adminEmail == "" || adminMobileRaw == "" || adminPassword == "" {
		appLogger.Info("Admin seeder: skipped (ADMIN_EMAIL / ADMIN_MOBILE / ADMIN_PASSWORD or ADMIN_PASSWARD not set)")
		return nil
	}
	adminMobile, err := normalizeAdminMobile(adminMobileRaw)
	if err != nil {
		return fmt.Errorf("admin seeder: invalid ADMIN_MOBILE %q: %w", adminMobileRaw, err)
	}
	if db == nil {
		return nil
	}

	userRepo := repository.NewUserRepository(db)

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Check if exists
	existing, err := userRepo.GetByEmail(ctx, adminEmail)
	if err == nil && existing != nil {
		// Best-effort updates (repository methods are safer than full update due to schema drift).
		_ = userRepo.UpdatePassword(ctx, existing.UserId, string(hash))
		_ = userRepo.UpdateStatus(ctx, existing.UserId, authnentityv1.UserStatus_USER_STATUS_ACTIVE)
		_ = userRepo.UpdateEmailVerified(ctx, existing.UserId)
		if existing.MobileNumber != adminMobile {
			_ = db.WithContext(ctx).
				Table("authn_schema.users").
				Where("user_id = ?", existing.UserId).
				Updates(map[string]any{
					"mobile_number": adminMobile,
					"updated_at":    time.Now(),
				}).Error
		}
		appLogger.Infof("Admin seeder: admin user already exists, ensured active (email=%s user_id=%s)", adminEmail, existing.UserId)
		return nil
	}

	admin := &authnentityv1.User{
		UserId:             uuid.NewString(),
		MobileNumber:       adminMobile,
		Email:              adminEmail,
		PasswordHash:       string(hash),
		Status:             authnentityv1.UserStatus_USER_STATUS_ACTIVE,
		UserType:           authnentityv1.UserType_USER_TYPE_SYSTEM_USER,
		EmailVerified:      true,
		EmailVerifiedAt:    timestamppb.Now(),
		EmailLoginAttempts: 0,
		CreatedAt:          timestamppb.Now(),
		UpdatedAt:          timestamppb.Now(),
	}

	if err := userRepo.CreateFull(ctx, admin); err != nil {
		return err
	}

	appLogger.Infof("Admin seeder: seeded admin user (email=%s user_id=%s)", adminEmail, admin.UserId)
	return nil
}

// SeedB2bAdminUser bootstraps a default B2B_ORG_ADMIN account.
func SeedB2bAdminUser(ctx context.Context, db *gorm.DB) error {
	adminEmail := os.Getenv("B2B_ADMIN")
	adminMobileRaw := os.Getenv("B2B_ADMIN_MOBILE")
	adminPassword := os.Getenv("B2B_ADMIN_PASSWARD")

	if adminEmail == "" || adminMobileRaw == "" || adminPassword == "" {
		appLogger.Info("B2B Admin seeder: skipped (B2B_ADMIN / B2B_ADMIN_MOBILE / B2B_ADMIN_PASSWARD not set)")
		return nil
	}
	adminMobile, err := normalizeAdminMobile(adminMobileRaw)
	if err != nil {
		return fmt.Errorf("b2b admin seeder: invalid B2B_ADMIN_MOBILE %q: %w", adminMobileRaw, err)
	}
	if db == nil {
		return nil
	}

	userRepo := repository.NewUserRepository(db)

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Check if exists
	existing, err := userRepo.GetByEmail(ctx, adminEmail)
	if err == nil && existing != nil {
		// Best-effort updates
		_ = userRepo.UpdatePassword(ctx, existing.UserId, string(hash))
		_ = userRepo.UpdateStatus(ctx, existing.UserId, authnentityv1.UserStatus_USER_STATUS_ACTIVE)
		_ = userRepo.UpdateEmailVerified(ctx, existing.UserId)
		if existing.MobileNumber != adminMobile {
			_ = db.WithContext(ctx).
				Table("authn_schema.users").
				Where("user_id = ?", existing.UserId).
				Updates(map[string]any{
					"mobile_number": adminMobile,
					"updated_at":    time.Now(),
				}).Error
		}
		appLogger.Infof("B2B Admin seeder: b2b admin user already exists, ensured active (email=%s user_id=%s)", adminEmail, existing.UserId)
		return nil
	}

	admin := &authnentityv1.User{
		UserId:             uuid.NewString(),
		MobileNumber:       adminMobile,
		Email:              adminEmail,
		PasswordHash:       string(hash),
		Status:             authnentityv1.UserStatus_USER_STATUS_ACTIVE,
		UserType:           authnentityv1.UserType_USER_TYPE_B2B_ORG_ADMIN,
		EmailVerified:      true,
		EmailVerifiedAt:    timestamppb.Now(),
		EmailLoginAttempts: 0,
		CreatedAt:          timestamppb.Now(),
		UpdatedAt:          timestamppb.Now(),
	}

	if err := userRepo.CreateFull(ctx, admin); err != nil {
		return err
	}

	appLogger.Infof("B2B Admin seeder: seeded b2b admin user (email=%s user_id=%s)", adminEmail, admin.UserId)
	return nil
}

func normalizeAdminMobile(mobile string) (string, error) {
	normalized, err := sms.NormalizePhoneNumber(mobile)
	if err != nil {
		return "", err
	}
	// DB constraint requires leading '+' in E.164-like format:
	// ^\+8801[0-9]{9}$
	return "+" + normalized, nil
}

// SeedDocumentTypes seeds the default document types (idempotent upsert by code).
func SeedDocumentTypes(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return nil
	}
	docTypeRepo := repository.NewDocumentTypeRepository(db)

	types := []struct {
		code        string
		name        string
		description string
	}{
		{"NID", "National ID Card", "Bangladesh National Identity Card"},
		{"PASSPORT", "Passport", "International travel passport"},
		{"BIRTH_CERTIFICATE", "Birth Certificate", "Official birth certificate"},
		{"DRIVING_LICENSE", "Driving License", "Bangladesh driving license"},
		{"TIN_CERTIFICATE", "TIN Certificate", "Tax Identification Number certificate"},
	}

	for _, t := range types {
		existing, err := docTypeRepo.GetByCode(ctx, t.code)
		if err == nil && existing != nil {
			continue // already exists
		}
		dt := &authnentityv1.DocumentType{
			DocumentTypeId: uuid.NewString(),
			Code:           t.code,
			Name:           t.name,
			Description:    t.description,
			IsActive:       true,
		}
		if err := docTypeRepo.Create(ctx, dt); err != nil {
			appLogger.Warnf("Document type seeder: failed to seed %s: %v", t.code, err)
		} else {
			appLogger.Infof("Document type seeder: seeded %s", t.code)
		}
	}
	return nil
}
