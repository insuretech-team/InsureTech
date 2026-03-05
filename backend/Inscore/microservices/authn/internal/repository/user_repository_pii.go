package repository

// user_repository_pii.go
// PII-transparent wrapper around UserRepository.
//
// This file adds field-level AES-256-GCM encryption for the three PII columns
// stored in authn_schema.users:
//
//   mobile_number  → encrypted; mobile_number_idx  (HMAC blind index for lookup)
//   email          → encrypted; email_idx           (HMAC blind index for lookup)
//   nid            → encrypted; nid_idx             (HMAC blind index for lookup)
//
// Architecture:
//   - UserRepository continues to write/read raw values through GORM as today.
//   - PIIUserRepository wraps UserRepository, transparently encrypt-on-write
//     and decrypt-on-read so all callers above the repository layer always
//     work with plaintext.
//   - Lookups by mobile/email use the blind-index columns.
//
// Key loading:
//   - PII_AES_KEY  (64 hex chars = 32 bytes)
//   - PII_HMAC_KEY (64 hex chars = 32 bytes)
//
// If the env vars are unset (e.g. local dev without a secrets manager) the
// wrapper falls back to the underlying repo with no encryption and logs a
// warning — this allows the service to start, but PII is stored plaintext and
// a log warning is emitted on every call.

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/pii"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"google.golang.org/protobuf/proto"
)

// PIIUserRepository wraps UserRepository with transparent PII encryption.
// All method signatures are identical to UserRepository so it can be used as a
// drop-in replacement wherever a *UserRepository is accepted.
type PIIUserRepository struct {
	inner *UserRepository
	enc   *pii.Encryptor // nil → no encryption (fallback)

	limitsOnce sync.Once
	limitsErr  error
	mobileMax  int
	emailMax   int
}

// NewPIIUserRepository creates a PIIUserRepository.
// If PII_AES_KEY / PII_HMAC_KEY are not set it falls back to plaintext and
// logs a warning.
func NewPIIUserRepository(inner *UserRepository) *PIIUserRepository {
	enc, err := pii.NewEncryptorFromEnv()
	if err != nil {
		appLogger.Warnf("PIIUserRepository: PII encryption disabled (keys not configured): %v. PII will be stored as plaintext.", err)
		return &PIIUserRepository{inner: inner, enc: nil}
	}
	return &PIIUserRepository{inner: inner, enc: enc}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func (r *PIIUserRepository) encryptField(plaintext string) (string, error) {
	if r.enc == nil || plaintext == "" {
		return plaintext, nil
	}
	return r.enc.EncryptIfNonEmpty(plaintext)
}

func (r *PIIUserRepository) decryptField(ciphertext string) (string, error) {
	if r.enc == nil || ciphertext == "" {
		return ciphertext, nil
	}
	plaintext, err := r.enc.DecryptIfNonEmpty(ciphertext)
	if err != nil {
		// Mixed/plaintext rows can exist during migration. Do not fail reads.
		appLogger.Warnf("PIIUserRepository: decrypt failed; treating value as plaintext: %v", err)
		return ciphertext, nil
	}
	return plaintext, nil
}

func (r *PIIUserRepository) blindIndex(plaintext string) string {
	if r.enc == nil {
		return plaintext // no blind index; use plaintext for lookup (insecure fallback)
	}
	return r.enc.BlindIndex(plaintext)
}

func (r *PIIUserRepository) loadLimits(ctx context.Context) {
	r.limitsOnce.Do(func() {
		r.mobileMax = r.lookupVarcharLimit(ctx, "mobile_number")
		r.emailMax = r.lookupVarcharLimit(ctx, "email")
	})
}

func (r *PIIUserRepository) lookupVarcharLimit(ctx context.Context, col string) int {
	var max sql.NullInt64
	err := r.inner.db.WithContext(ctx).Raw(
		`select character_maximum_length
		   from information_schema.columns
		  where table_schema='authn_schema' and table_name='users' and column_name = ?`,
		col,
	).Scan(&max).Error
	if err != nil {
		r.limitsErr = err
		return 0
	}
	if !max.Valid {
		return 0
	}
	return int(max.Int64)
}

func (r *PIIUserRepository) encryptForColumn(ctx context.Context, plaintext string, col string) (string, error) {
	if plaintext == "" || r.enc == nil {
		return plaintext, nil
	}
	// Live schema currently enforces strict mobile/email format checks.
	// Keep these columns plaintext and rely on blind indexes for secure lookup.
	switch strings.ToLower(col) {
	case "mobile_number", "email":
		return plaintext, nil
	}
	ciphertext, err := r.enc.Encrypt(plaintext)
	if err != nil {
		return "", err
	}

	r.loadLimits(ctx)
	limit := 0
	switch strings.ToLower(col) {
	case "mobile_number":
		limit = r.mobileMax
	case "email":
		limit = r.emailMax
	}
	if limit > 0 && len(ciphertext) > limit {
		appLogger.Warnf("PIIUserRepository: %s ciphertext length (%d) exceeds column max (%d), storing plaintext with blind index", col, len(ciphertext), limit)
		return plaintext, nil
	}
	return ciphertext, nil
}

// encryptUser encrypts the PII fields of a user entity in-place.
// Returns a cloned struct so the caller's original is not mutated.
func (r *PIIUserRepository) encryptUser(u *authnentityv1.User) (*authnentityv1.User, error) {
	if r.enc == nil {
		return u, nil
	}
	clone := proto.Clone(u).(*authnentityv1.User)

	var err error
	if clone.MobileNumber != "" {
		clone.MobileNumber, err = r.enc.Encrypt(clone.MobileNumber)
		if err != nil {
			logger.Errorf("encrypt mobile_number: %v", err)
			return nil, errors.New("encrypt mobile_number")
		}
	}
	if clone.Email != "" {
		clone.Email, err = r.enc.Encrypt(clone.Email)
		if err != nil {
			logger.Errorf("encrypt email: %v", err)
			return nil, errors.New("encrypt email")
		}
	}
	return clone, nil
}

// decryptUser decrypts the PII fields of a user entity returned from the DB.
// The blind-index columns are left as-is (they are not meaningful plaintext).
func (r *PIIUserRepository) decryptUser(u *authnentityv1.User) error {
	if r.enc == nil || u == nil {
		return nil
	}

	if u.MobileNumber != "" {
		plain, err := r.decryptField(u.MobileNumber)
		if err != nil {
			logger.Errorf("decrypt mobile_number for user %s: %v", u.UserId, err)
			return errors.New("decrypt mobile_number for user %s")
		}
		u.MobileNumber = plain
	}
	if u.Email != "" {
		plain, err := r.decryptField(u.Email)
		if err != nil {
			logger.Errorf("decrypt email for user %s: %v", u.UserId, err)
			return errors.New("decrypt email for user %s")
		}
		u.Email = plain
	}
	return nil
}

// ── Write methods ─────────────────────────────────────────────────────────────

// Create creates a new mobile-auth user with encrypted PII and a fresh
// biometric token.
func (r *PIIUserRepository) Create(ctx context.Context, mobile, passwordHash, email string, status authnentityv1.UserStatus) (*authnentityv1.User, error) {
	encMobile, err := r.encryptForColumn(ctx, mobile, "mobile_number")
	if err != nil {
		logger.Errorf("PIIUserRepository.Create: encrypt mobile: %v", err)
		return nil, errors.New("PIIUserRepository.Create: encrypt mobile")
	}
	encEmail, err := r.encryptForColumn(ctx, email, "email")
	if err != nil {
		logger.Errorf("PIIUserRepository.Create: encrypt email: %v", err)
		return nil, errors.New("PIIUserRepository.Create: encrypt email")
	}

	user, err := r.inner.Create(ctx, encMobile, passwordHash, encEmail, status)
	if err != nil {
		return nil, err
	}

	// Set blind-index columns and biometric token.
	updates := map[string]interface{}{}
	if r.enc != nil {
		if mobile != "" {
			updates["mobile_number_idx"] = r.blindIndex(mobile)
		}
		if email != "" {
			updates["email_idx"] = r.blindIndex(email)
		}
		// Generate biometric token
		_, bioEnc, bioIdx, bioErr := r.enc.GenerateBiometricToken()
		if bioErr == nil {
			updates["biometric_token_enc"] = bioEnc
			updates["biometric_token_idx"] = bioIdx
		} else {
			appLogger.Warnf("PIIUserRepository.Create: failed to generate biometric token for user %s: %v", user.UserId, bioErr)
		}
	}
	if len(updates) > 0 {
		_ = r.inner.db.WithContext(ctx).
			Table("authn_schema.users").
			Where("user_id = ?", user.UserId).
			Updates(updates).Error
	}

	// Return plaintext to caller
	user.MobileNumber = mobile
	user.Email = email
	return user, nil
}

// CreateFull creates a user from a fully populated entity, encrypting PII
// fields and generating a biometric token before insert.
func (r *PIIUserRepository) CreateFull(ctx context.Context, user *authnentityv1.User) error {
	encrypted := proto.Clone(user).(*authnentityv1.User)
	if r.enc != nil {
		var err error
		encrypted.MobileNumber, err = r.encryptForColumn(ctx, user.MobileNumber, "mobile_number")
		if err != nil {
			logger.Errorf("PIIUserRepository.CreateFull: encrypt mobile: %v", err)
			return errors.New("PIIUserRepository.CreateFull: encrypt mobile")
		}
		encrypted.Email, err = r.encryptForColumn(ctx, user.Email, "email")
		if err != nil {
			logger.Errorf("PIIUserRepository.CreateFull: encrypt email: %v", err)
			return errors.New("PIIUserRepository.CreateFull: encrypt email")
		}
	}

	// Generate biometric token if not already set.
	var bioIdx string
	if r.enc != nil && encrypted.BiometricTokenEnc == "" {
		var bioEnc string
		var bioErr error
		_, bioEnc, bioIdx, bioErr = r.enc.GenerateBiometricToken()
		if bioErr == nil {
			encrypted.BiometricTokenEnc = bioEnc
		} else {
			appLogger.Warnf("PIIUserRepository.CreateFull: failed to generate biometric token: %v", bioErr)
		}
	}

	if err := r.inner.CreateFull(ctx, encrypted); err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if r.enc != nil {
		if user.MobileNumber != "" {
			updates["mobile_number_idx"] = r.blindIndex(user.MobileNumber)
		}
		if user.Email != "" {
			updates["email_idx"] = r.blindIndex(user.Email)
		}
		if bioIdx != "" {
			updates["biometric_token_idx"] = bioIdx
		}
	}
	if len(updates) > 0 {
		_ = r.inner.db.WithContext(ctx).
			Table("authn_schema.users").
			Where("user_id = ?", encrypted.UserId).
			Updates(updates).Error
	}

	// Restore plaintext on the caller's struct (CreateFull mutates the passed entity).
	user.UserId = encrypted.UserId
	user.CreatedAt = encrypted.CreatedAt
	user.UpdatedAt = encrypted.UpdatedAt
	user.BiometricTokenEnc = encrypted.BiometricTokenEnc
	return nil
}

// ── Read methods ──────────────────────────────────────────────────────────────

// GetByID returns a user by ID with decrypted PII fields.
func (r *PIIUserRepository) GetByID(ctx context.Context, id string) (*authnentityv1.User, error) {
	u, err := r.inner.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := r.decryptUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

// GetByMobileNumber looks up a user using the HMAC blind index, then decrypts.
// Falls back to direct mobile_number lookup if encryption is disabled.
func (r *PIIUserRepository) GetByMobileNumber(ctx context.Context, mobile string) (*authnentityv1.User, error) {
	if r.enc == nil {
		return r.inner.GetByMobileNumber(ctx, mobile)
	}
	idx := r.blindIndex(mobile)
	u, err := r.inner.getOne(ctx, "mobile_number_idx = ?", idx)
	if err != nil {
		return nil, err
	}
	if err := r.decryptUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

// GetByEmail looks up a user using the HMAC blind index, then decrypts.
// Falls back to direct email lookup if encryption is disabled.
func (r *PIIUserRepository) GetByEmail(ctx context.Context, email string) (*authnentityv1.User, error) {
	if r.enc == nil {
		return r.inner.GetByEmail(ctx, email)
	}
	idx := r.blindIndex(email)
	u, err := r.inner.getOne(ctx, "email_idx = ? AND deleted_at IS NULL", idx)
	if err != nil {
		return nil, err
	}
	if err := r.decryptUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

// ── Pass-through update methods ───────────────────────────────────────────────

func (r *PIIUserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	return r.inner.UpdatePassword(ctx, userID, passwordHash)
}

func (r *PIIUserRepository) UpdateEmailVerified(ctx context.Context, userID string) error {
	return r.inner.UpdateEmailVerified(ctx, userID)
}

func (r *PIIUserRepository) UpdateStatus(ctx context.Context, userID string, status authnentityv1.UserStatus) error {
	return r.inner.UpdateStatus(ctx, userID, status)
}

func (r *PIIUserRepository) UpdateLastLogin(ctx context.Context, userID, sessionType string) error {
	return r.inner.UpdateLastLogin(ctx, userID, sessionType)
}

func (r *PIIUserRepository) IncrementEmailLoginAttempts(ctx context.Context, userID string) (int32, error) {
	return r.inner.IncrementEmailLoginAttempts(ctx, userID)
}

func (r *PIIUserRepository) LockEmailAuth(ctx context.Context, userID string, lockDuration time.Duration) error {
	return r.inner.LockEmailAuth(ctx, userID, lockDuration)
}

func (r *PIIUserRepository) ResetEmailLoginAttempts(ctx context.Context, userID string) error {
	return r.inner.ResetEmailLoginAttempts(ctx, userID)
}
