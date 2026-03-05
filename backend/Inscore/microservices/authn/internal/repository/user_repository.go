package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// getOne retrieves a single User using the proto-generated struct directly.
// user.pb.go has GORM tags on every field (injected by scripts/inject_gorm_tags.go),
// including serializer:proto_enum for UserStatus/UserType and serializer:proto_timestamp
// for all *timestamppb.Timestamp fields — so GORM handles all conversions natively.
func (r *UserRepository) getOne(ctx context.Context, where string, args ...any) (*authnentityv1.User, error) {
	query := `SELECT user_id, mobile_number, email, password_hash, status, user_type, created_at, updated_at, wallet_balance,
	                 email_verified, email_verified_at, email_login_attempts, login_attempts, last_login_at, last_login_session_type,
	                 totp_enabled, totp_secret_enc, locked_until, email_locked_until, notification_preference,
	                 preferred_language, biometric_token_enc
	            FROM authn_schema.users
	           WHERE ` + where + ` LIMIT 1`
	row := r.db.WithContext(ctx).Raw(query, args...).Row()
	if err := row.Err(); err != nil {
		return nil, err
	}

	var (
		u                      authnentityv1.User
		mobile, email          sql.NullString
		passwordHash           sql.NullString
		statusStr, userTypeStr sql.NullString
		createdAt, updatedAt   sql.NullTime
		walletBalance          sql.NullInt64
		emailVerified          sql.NullBool
		emailVerifiedAt        sql.NullTime
		emailLoginAttempts     sql.NullInt64
		loginAttempts          sql.NullInt64
		lastLoginAt            sql.NullTime
		lastLoginSessionType   sql.NullString
		totpEnabled            sql.NullBool
		totpSecretEnc          sql.NullString
		lockedUntil            sql.NullTime
		emailLockedUntil       sql.NullTime
		notificationPreference sql.NullString
		preferredLanguage      sql.NullString
		biometricTokenEnc      sql.NullString
	)
	if err := row.Scan(
		&u.UserId,
		&mobile,
		&email,
		&passwordHash,
		&statusStr,
		&userTypeStr,
		&createdAt,
		&updatedAt,
		&walletBalance,
		&emailVerified,
		&emailVerifiedAt,
		&emailLoginAttempts,
		&loginAttempts,
		&lastLoginAt,
		&lastLoginSessionType,
		&totpEnabled,
		&totpSecretEnc,
		&lockedUntil,
		&emailLockedUntil,
		&notificationPreference,
		&preferredLanguage,
		&biometricTokenEnc,
	); err != nil {
		return nil, err
	}

	if mobile.Valid {
		u.MobileNumber = mobile.String
	}
	if email.Valid {
		u.Email = email.String
	}
	if passwordHash.Valid {
		u.PasswordHash = passwordHash.String
	}
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := authnentityv1.UserStatus_value[k]; ok {
			u.Status = authnentityv1.UserStatus(v)
		}
	}
	if userTypeStr.Valid {
		k := strings.ToUpper(userTypeStr.String)
		if v, ok := authnentityv1.UserType_value[k]; ok {
			u.UserType = authnentityv1.UserType(v)
		}
	}
	if createdAt.Valid {
		u.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		u.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	u.WalletBalance = &commonv1.Money{
		Amount:   walletBalance.Int64,
		Currency: "BDT",
	}
	u.EmailVerified = emailVerified.Valid && emailVerified.Bool
	if emailVerifiedAt.Valid {
		u.EmailVerifiedAt = timestamppb.New(emailVerifiedAt.Time)
	}
	if emailLoginAttempts.Valid {
		u.EmailLoginAttempts = int32(emailLoginAttempts.Int64)
	}
	if loginAttempts.Valid {
		u.LoginAttempts = int32(loginAttempts.Int64)
	}
	if lastLoginAt.Valid {
		u.LastLoginAt = timestamppb.New(lastLoginAt.Time)
	}
	if lastLoginSessionType.Valid {
		u.LastLoginSessionType = lastLoginSessionType.String
	}
	u.TotpEnabled = totpEnabled.Valid && totpEnabled.Bool
	if totpSecretEnc.Valid {
		u.TotpSecretEnc = totpSecretEnc.String
	}
	if lockedUntil.Valid {
		u.LockedUntil = timestamppb.New(lockedUntil.Time)
	}
	if emailLockedUntil.Valid {
		u.EmailLockedUntil = timestamppb.New(emailLockedUntil.Time)
	}
	if notificationPreference.Valid {
		u.NotificationPreference = notificationPreference.String
	}
	if preferredLanguage.Valid {
		u.PreferredLanguage = preferredLanguage.String
	}
	if biometricTokenEnc.Valid {
		u.BiometricTokenEnc = biometricTokenEnc.String
	}

	return &u, nil
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, mobile, passwordHash, email string, status authnentityv1.UserStatus) (*authnentityv1.User, error) {
	now := time.Now()
	userID := uuid.New().String()
	if status == authnentityv1.UserStatus_USER_STATUS_UNSPECIFIED {
		status = authnentityv1.UserStatus_USER_STATUS_ACTIVE
	}

	// Use map insert (not struct insert) to avoid GORM reflection on proto-only fields like WalletBalance.
	values := map[string]any{
		"user_id":               userID,
		"mobile_number":         mobile,
		"password_hash":         passwordHash,
		"email":                 email,
		"status":                status.String(),
		"user_type":             authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER.String(),
		"wallet_balance":        int64(0), // paisa
		"active_policies_count": 0,
		"pending_claims_count":  0,
		"biometric_token_enc":   "",
		"created_at":            now,
		"updated_at":            now,
	}
	if err := r.db.WithContext(ctx).Table("authn_schema.users").Create(values).Error; err != nil {
		return nil, err
	}

	return &authnentityv1.User{
		UserId:       userID,
		MobileNumber: mobile,
		PasswordHash: passwordHash,
		Email:        email,
		Status:       status,
		CreatedAt:    timestamppb.New(now),
		UpdatedAt:    timestamppb.New(now),
		WalletBalance: &commonv1.Money{
			Amount:   0,
			Currency: "BDT",
		},
		ActivePoliciesCount: 0,
		PendingClaimsCount:  0,
		BiometricTokenEnc:   "",
	}, nil
}

// CreateFull creates a user from a fully populated User entity (for email-based registration)
func (r *UserRepository) CreateFull(ctx context.Context, user *authnentityv1.User) error {
	if user.UserId == "" {
		user.UserId = uuid.New().String()
	}
	now := time.Now()
	user.CreatedAt = timestamppb.New(now)
	user.UpdatedAt = timestamppb.New(now)
	if user.WalletBalance == nil {
		user.WalletBalance = &commonv1.Money{Amount: 0, Currency: "BDT"}
	}
	if user.Status == authnentityv1.UserStatus_USER_STATUS_UNSPECIFIED {
		user.Status = authnentityv1.UserStatus_USER_STATUS_ACTIVE
	}
	if user.UserType == authnentityv1.UserType_USER_TYPE_UNSPECIFIED {
		user.UserType = authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER
	}

	values := map[string]any{
		"user_id":                 user.UserId,
		"mobile_number":           user.MobileNumber,
		"email":                   user.Email,
		"password_hash":           user.PasswordHash,
		"status":                  user.Status.String(),
		"user_type":               user.UserType.String(),
		"wallet_balance":          user.WalletBalance.Amount,
		"active_policies_count":   user.ActivePoliciesCount,
		"pending_claims_count":    user.PendingClaimsCount,
		"biometric_token_enc":     user.BiometricTokenEnc,
		"email_verified":          user.EmailVerified,
		"email_login_attempts":    user.EmailLoginAttempts,
		"login_attempts":          user.LoginAttempts,
		"notification_preference": user.NotificationPreference,
		"preferred_language":      user.PreferredLanguage,
		"created_at":              now,
		"updated_at":              now,
	}
	return r.db.WithContext(ctx).Table("authn_schema.users").Create(values).Error
}

// Create creates a new user in the database
func (r *UserRepository) createLegacyUnused(ctx context.Context, mobile, passwordHash, email string, status authnentityv1.UserStatus) (*authnentityv1.User, error) {
	user := &authnentityv1.User{
		UserId:       uuid.New().String(),
		MobileNumber: mobile,
		PasswordHash: passwordHash,
		Email:        email,
		Status:       status,
		CreatedAt:    timestamppb.Now(),
		UpdatedAt:    timestamppb.Now(),
		WalletBalance: &commonv1.Money{
			Amount:   0,
			Currency: "BDT",
		},
		ActivePoliciesCount: 0,
		PendingClaimsCount:  0,
		BiometricTokenEnc:   "",
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetByMobileNumber retrieves a user by mobile number
func (r *UserRepository) GetByMobileNumber(ctx context.Context, mobile string) (*authnentityv1.User, error) {
	return r.getOne(ctx, "mobile_number = ?", mobile)
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*authnentityv1.User, error) {
	return r.getOne(ctx, "user_id = ?", id)
}

// GetByEmail retrieves a user by email address
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*authnentityv1.User, error) {
	return r.getOne(ctx, "email = ? AND deleted_at IS NULL", email)
}

// UpdatePassword updates user's password hash
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	result := r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"password_hash": passwordHash,
			"updated_at":    time.Now(),
		})
	return result.Error
}

// UpdateEmailVerified marks the user's email as verified
func (r *UserRepository) UpdateEmailVerified(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"email_verified":    true,
			"email_verified_at": time.Now(),
			"updated_at":        time.Now(),
		}).Error
}

// UpdateStatus updates the user's account status
func (r *UserRepository) UpdateStatus(ctx context.Context, userID string, status authnentityv1.UserStatus) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// IncrementEmailLoginAttempts increments the failed email login counter
// Returns the new attempt count
func (r *UserRepository) IncrementEmailLoginAttempts(ctx context.Context, userID string) (int32, error) {
	result := r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		UpdateColumn("email_login_attempts", gorm.Expr("email_login_attempts + 1"))
	if result.Error != nil {
		return 0, result.Error
	}
	// Fetch updated count
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.EmailLoginAttempts, nil
}

// LockEmailAuth sets email_locked_until for the user (30 min lockout)
func (r *UserRepository) LockEmailAuth(ctx context.Context, userID string, lockDuration time.Duration) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"email_locked_until": time.Now().Add(lockDuration),
			"updated_at":         time.Now(),
		}).Error
}

// ResetEmailLoginAttempts resets the email login attempt counter after successful login
func (r *UserRepository) ResetEmailLoginAttempts(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"email_login_attempts": 0,
			"email_locked_until":   nil,
			"updated_at":           time.Now(),
		}).Error
}

// GetByBiometricTokenIdx retrieves a user by the HMAC blind index of their biometric token.
func (r *UserRepository) GetByBiometricTokenIdx(ctx context.Context, tokenIdx string) (*authnentityv1.User, error) {
	return r.getOne(ctx, "biometric_token_idx = ? AND deleted_at IS NULL", tokenIdx)
}

// UpdateLastLogin updates last_login_at and last_login_session_type
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID, sessionType string) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at":           time.Now(),
			"last_login_session_type": sessionType,
			"updated_at":              time.Now(),
		}).Error
}

// UpdateTOTPSecret stores an encrypted TOTP secret for the user.
// Pass an empty string to clear the secret (on DisableTOTP).
func (r *UserRepository) UpdateTOTPSecret(ctx context.Context, userID, encSecret string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE authn_schema.users SET totp_secret_enc = NULLIF(?, ''), updated_at = NOW() WHERE user_id = ?`,
		encSecret, userID,
	).Error
}

// SetTOTPEnabled sets or clears the totp_enabled flag for a user.
func (r *UserRepository) SetTOTPEnabled(ctx context.Context, userID string, enabled bool) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE authn_schema.users SET totp_enabled = ?, updated_at = NOW() WHERE user_id = ?`,
		enabled, userID,
	).Error
}

// IncrementLoginAttempts increments the failed mobile login counter and returns the new count.
func (r *UserRepository) IncrementLoginAttempts(ctx context.Context, userID string) (int32, error) {
	result := r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		UpdateColumn("login_attempts", gorm.Expr("login_attempts + 1"))
	if result.Error != nil {
		return 0, result.Error
	}
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.LoginAttempts, nil
}

// LockAccount sets locked_until for the user (account lockout after too many failures).
func (r *UserRepository) LockAccount(ctx context.Context, userID string, lockDuration time.Duration) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"locked_until": time.Now().Add(lockDuration),
			"updated_at":   time.Now(),
		}).Error
}

// ResetLoginAttempts resets the mobile login attempt counter after successful login.
func (r *UserRepository) ResetLoginAttempts(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.users").
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"login_attempts": 0,
			"locked_until":   nil,
			"updated_at":     time.Now(),
		}).Error
}

// UpdateNotificationPreferences updates notification_preference and preferred_language for a user.
func (r *UserRepository) UpdateNotificationPreferences(ctx context.Context, userID, notificationPreference, preferredLanguage string) error {
	upd := map[string]any{"updated_at": "NOW()"}
	if notificationPreference != "" {
		upd["notification_preference"] = notificationPreference
	}
	if preferredLanguage != "" {
		upd["preferred_language"] = preferredLanguage
	}
	return r.db.WithContext(ctx).Table("authn_schema.users").Where("user_id = ?", userID).Updates(upd).Error
}
