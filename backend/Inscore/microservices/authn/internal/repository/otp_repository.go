package repository

import (
	"context"
	"time"

	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// OTPRepository handles OTP database operations
type OTPRepository struct {
	db *gorm.DB
}

// NewOTPRepository creates a new OTP repository
func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

// getOne retrieves a single OTP using the proto-generated struct directly.
// otp.pb.go has GORM tags on every field (injected by scripts/inject_gorm_tags.go),
// so GORM can scan all columns without an intermediate row struct.
// The ip_address column is INET in PostgreSQL; GORM scans it as text into the string field safely.
func (r *OTPRepository) getOne(ctx context.Context, where string, args ...any) (*authnentityv1.OTP, error) {
	var otp authnentityv1.OTP
	err := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where(where, args...).
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// Create creates a new OTP record
func (r *OTPRepository) Create(ctx context.Context, otp *authnentityv1.OTP) error {
	otp.CreatedAt = timestamppb.Now()
	return r.db.WithContext(ctx).Create(otp).Error
}

// GetByID retrieves an OTP by ID
func (r *OTPRepository) GetByID(ctx context.Context, otpID string) (*authnentityv1.OTP, error) {
	return r.getOne(ctx, "otp_id = ?", otpID)
}

// GetByProviderMessageID retrieves an OTP by provider message ID (for DLR tracking)
func (r *OTPRepository) GetByProviderMessageID(ctx context.Context, providerMessageID string) (*authnentityv1.OTP, error) {
	return r.getOne(ctx, "provider_message_id = ?", providerMessageID)
}

// IncrementAttempts increments the attempt count for an OTP
func (r *OTPRepository) IncrementAttempts(ctx context.Context, otpID string) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("otp_id = ?", otpID).
		UpdateColumn("attempts", gorm.Expr("attempts + ?", 1)).
		Error
}

// MarkVerified marks an OTP as verified
func (r *OTPRepository) MarkVerified(ctx context.Context, otpID string) error {
	updates := map[string]interface{}{
		"verified":    true,
		"verified_at": time.Now(),
	}
	return r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("otp_id = ?", otpID).
		Updates(updates).
		Error
}

// ExpireOTP sets expires_at to now and marks the OTP as unusable.
// We do NOT set verified=true, because that would semantically mean success.
func (r *OTPRepository) ExpireOTP(ctx context.Context, otpID string) error {
	updates := map[string]any{
		"expires_at": time.Now(),
	}
	return r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("otp_id = ?", otpID).
		Updates(updates).Error
}

// CountRecentOTPs counts OTPs sent to a recipient since a given time
func (r *OTPRepository) CountRecentOTPs(ctx context.Context, recipient string, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("recipient = ? AND created_at >= ?", recipient, since).
		Count(&count).
		Error
	return count, err
}

// GetLastOTP retrieves the most recent OTP for a recipient
func (r *OTPRepository) GetLastOTP(ctx context.Context, recipient string) (*authnentityv1.OTP, error) {
	var otp authnentityv1.OTP
	err := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("recipient = ?", recipient).
		Order("created_at DESC").
		First(&otp).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

// ListByRecipient lists all OTPs for a recipient with optional filters
func (r *OTPRepository) ListByRecipient(ctx context.Context, recipient string, verified *bool, otpType string) ([]*authnentityv1.OTP, error) {
	var otps []*authnentityv1.OTP

	query := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("recipient = ?", recipient)

	if verified != nil {
		query = query.Where("verified = ?", *verified)
	}

	if otpType != "" {
		query = query.Where("purpose = ?", otpType)
	}

	err := query.Order("created_at DESC").Find(&otps).Error
	return otps, err
}

// CleanupExpiredOTPs deletes expired OTPs (background job)
func (r *OTPRepository) CleanupExpiredOTPs(ctx context.Context, olderThan time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("expires_at < ?", olderThan).
		Delete(map[string]any{})

	return result.RowsAffected, result.Error
}

// GetPendingDLRs retrieves OTPs with pending delivery reports
func (r *OTPRepository) GetPendingDLRs(ctx context.Context, limit int) ([]*authnentityv1.OTP, error) {
	var otps []*authnentityv1.OTP

	err := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("dlr_status = ? AND channel = ?", "PENDING", "sms").
		Where("created_at >= ?", time.Now().Add(-24*time.Hour)). // Only last 24 hours
		Order("created_at ASC").
		Limit(limit).
		Find(&otps).
		Error

	return otps, err
}

// GetStatsByCarrier retrieves OTP statistics grouped by carrier
func (r *OTPRepository) GetStatsByCarrier(ctx context.Context, since time.Time) (map[string]int64, error) {
	type Result struct {
		Carrier string
		Count   int64
	}

	var results []Result
	err := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Select("carrier, COUNT(*) as count").
		Where("created_at >= ? AND channel = ?", since, "sms").
		Group("carrier").
		Find(&results).
		Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Carrier] = r.Count
	}

	return stats, nil
}

// GetDeliveryRate calculates delivery success rate
func (r *OTPRepository) GetDeliveryRate(ctx context.Context, since time.Time) (float64, error) {
	var total, delivered int64

	// Total sent
	err := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("created_at >= ? AND channel = ?", since, "sms").
		Count(&total).
		Error
	if err != nil {
		return 0, err
	}

	if total == 0 {
		return 0, nil
	}

	// Delivered count
	err = r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("created_at >= ? AND channel = ? AND dlr_status = ?", since, "sms", "DELIVERED").
		Count(&delivered).
		Error
	if err != nil {
		return 0, err
	}

	return float64(delivered) / float64(total) * 100, nil
}

// UpdateDLRStatus updates the delivery report status for an OTP record
// matched by provider_message_id (the SMS gateway's message ID).
// dlrStatus values: "DELIVERED", "FAILED", "PENDING", "REJECTED", etc.
func (r *OTPRepository) UpdateDLRStatus(ctx context.Context, providerMessageID, dlrStatus, errorCode string) error {
	updates := map[string]interface{}{
		"dlr_status":     dlrStatus,
		"dlr_updated_at": time.Now(),
	}
	if errorCode != "" {
		updates["dlr_error_code"] = errorCode
	}
	result := r.db.WithContext(ctx).
		Table("authn_schema.otps").
		Where("provider_message_id = ?", providerMessageID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Non-fatal: DLR can arrive before OTP record is persisted in some edge cases.
		return nil
	}
	return nil
}
