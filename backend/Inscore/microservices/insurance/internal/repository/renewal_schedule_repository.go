package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	renewalv1 "github.com/newage-saint/insuretech/gen/go/insuretech/renewal/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type RenewalScheduleRepository struct {
	db *gorm.DB
}

func NewRenewalScheduleRepository(db *gorm.DB) *RenewalScheduleRepository {
	return &RenewalScheduleRepository{db: db}
}

func (r *RenewalScheduleRepository) Create(ctx context.Context, schedule *renewalv1.RenewalSchedule) (*renewalv1.RenewalSchedule, error) {
	if schedule.Id == "" {
		return nil, fmt.Errorf("schedule_id is required")
	}

	// Extract Money values
	renewalPremium := int64(0)
	renewalPremiumCurrency := "BDT"
	if schedule.RenewalPremium != nil {
		renewalPremium = schedule.RenewalPremium.Amount
		renewalPremiumCurrency = schedule.RenewalPremium.Currency
	}

	// Handle timestamps
	var renewalDueDate time.Time
	if schedule.RenewalDueDate != nil {
		renewalDueDate = schedule.RenewalDueDate.AsTime()
	}

	var gracePeriodEnd sql.NullTime
	if schedule.GracePeriodEnd != nil {
		gracePeriodEnd = sql.NullTime{Time: schedule.GracePeriodEnd.AsTime(), Valid: true}
	}

	var renewedAt sql.NullTime
	if schedule.RenewedAt != nil {
		renewedAt = sql.NullTime{Time: schedule.RenewedAt.AsTime(), Valid: true}
	}

	var renewedPolicyID sql.NullString
	if schedule.RenewedPolicyId != "" {
		renewedPolicyID = sql.NullString{String: schedule.RenewedPolicyId, Valid: true}
	}

	var auditInfo interface{}
	if schedule.AuditInfo != nil {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.renewal_schedules
			(schedule_id, policy_id, renewal_due_date, renewal_premium, renewal_premium_currency,
			 renewal_type, status, grace_period_days, grace_period_end, renewed_at,
			 renewed_policy_id, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		schedule.Id,
		schedule.PolicyId,
		renewalDueDate,
		renewalPremium,
		renewalPremiumCurrency,
		strings.ToUpper(schedule.RenewalType.String()),
		strings.ToUpper(schedule.Status.String()),
		schedule.GracePeriodDays,
		gracePeriodEnd,
		renewedAt,
		renewedPolicyID,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert renewal schedule: %w", err)
	}

	return r.GetByID(ctx, schedule.Id)
}

func (r *RenewalScheduleRepository) GetByID(ctx context.Context, scheduleID string) (*renewalv1.RenewalSchedule, error) {
	var (
		s                      renewalv1.RenewalSchedule
		renewalPremium         int64
		renewalPremiumCurrency string
		renewalTypeStr         sql.NullString
		statusStr              sql.NullString
		renewalDueDate         time.Time
		gracePeriodEnd         sql.NullTime
		renewedAt              sql.NullTime
		renewedPolicyID        sql.NullString
		auditInfo              sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT schedule_id, policy_id, renewal_due_date, renewal_premium, renewal_premium_currency,
		       renewal_type, status, grace_period_days, grace_period_end, renewed_at,
		       renewed_policy_id, audit_info
		FROM insurance_schema.renewal_schedules
		WHERE schedule_id = $1`,
		scheduleID,
	).Row().Scan(
		&s.Id,
		&s.PolicyId,
		&renewalDueDate,
		&renewalPremium,
		&renewalPremiumCurrency,
		&renewalTypeStr,
		&statusStr,
		&s.GracePeriodDays,
		&gracePeriodEnd,
		&renewedAt,
		&renewedPolicyID,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get renewal schedule: %w", err)
	}

	// Set Money field
	s.RenewalPremium = &commonv1.Money{Amount: renewalPremium, Currency: renewalPremiumCurrency}

	// Set optional fields
	if renewedPolicyID.Valid {
		s.RenewedPolicyId = renewedPolicyID.String
	}

	// Parse enums
	if renewalTypeStr.Valid {
		k := strings.ToUpper(renewalTypeStr.String)
		if v, ok := renewalv1.RenewalType_value[k]; ok {
			s.RenewalType = renewalv1.RenewalType(v)
		}
	}
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := renewalv1.RenewalStatus_value[k]; ok {
			s.Status = renewalv1.RenewalStatus(v)
		}
	}

	// Set timestamps
	if !renewalDueDate.IsZero() {
		s.RenewalDueDate = timestamppb.New(renewalDueDate)
	}
	if gracePeriodEnd.Valid {
		s.GracePeriodEnd = timestamppb.New(gracePeriodEnd.Time)
	}
	if renewedAt.Valid {
		s.RenewedAt = timestamppb.New(renewedAt.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		s.AuditInfo = &commonv1.AuditInfo{}
	}

	return &s, nil
}

func (r *RenewalScheduleRepository) Update(ctx context.Context, schedule *renewalv1.RenewalSchedule) (*renewalv1.RenewalSchedule, error) {
	// Extract Money values
	renewalPremium := int64(0)
	renewalPremiumCurrency := "BDT"
	if schedule.RenewalPremium != nil {
		renewalPremium = schedule.RenewalPremium.Amount
		renewalPremiumCurrency = schedule.RenewalPremium.Currency
	}

	// Handle timestamps
	var renewalDueDate time.Time
	if schedule.RenewalDueDate != nil {
		renewalDueDate = schedule.RenewalDueDate.AsTime()
	}

	var gracePeriodEnd sql.NullTime
	if schedule.GracePeriodEnd != nil {
		gracePeriodEnd = sql.NullTime{Time: schedule.GracePeriodEnd.AsTime(), Valid: true}
	}

	var renewedAt sql.NullTime
	if schedule.RenewedAt != nil {
		renewedAt = sql.NullTime{Time: schedule.RenewedAt.AsTime(), Valid: true}
	}

	var renewedPolicyID sql.NullString
	if schedule.RenewedPolicyId != "" {
		renewedPolicyID = sql.NullString{String: schedule.RenewedPolicyId, Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.renewal_schedules
		SET policy_id = $2,
		    renewal_due_date = $3,
		    renewal_premium = $4,
		    renewal_premium_currency = $5,
		    renewal_type = $6,
		    status = $7,
		    grace_period_days = $8,
		    grace_period_end = $9,
		    renewed_at = $10,
		    renewed_policy_id = $11
		WHERE schedule_id = $1`,
		schedule.Id,
		schedule.PolicyId,
		renewalDueDate,
		renewalPremium,
		renewalPremiumCurrency,
		strings.ToUpper(schedule.RenewalType.String()),
		strings.ToUpper(schedule.Status.String()),
		schedule.GracePeriodDays,
		gracePeriodEnd,
		renewedAt,
		renewedPolicyID,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update renewal schedule: %w", err)
	}

	return r.GetByID(ctx, schedule.Id)
}

func (r *RenewalScheduleRepository) Delete(ctx context.Context, scheduleID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.renewal_schedules
		WHERE schedule_id = $1`,
		scheduleID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete renewal schedule: %w", err)
	}

	return nil
}

func (r *RenewalScheduleRepository) ListByPolicyID(ctx context.Context, policyID string) ([]*renewalv1.RenewalSchedule, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT schedule_id, policy_id, renewal_due_date, renewal_premium, renewal_premium_currency,
		       renewal_type, status, grace_period_days, grace_period_end, renewed_at,
		       renewed_policy_id, audit_info
		FROM insurance_schema.renewal_schedules
		WHERE policy_id = $1
		ORDER BY renewal_due_date DESC`,
		policyID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list renewal schedules: %w", err)
	}
	defer rows.Close()

	schedules := make([]*renewalv1.RenewalSchedule, 0)
	for rows.Next() {
		var (
			s                      renewalv1.RenewalSchedule
			renewalPremium         int64
			renewalPremiumCurrency string
			renewalTypeStr         sql.NullString
			statusStr              sql.NullString
			renewalDueDate         time.Time
			gracePeriodEnd         sql.NullTime
			renewedAt              sql.NullTime
			renewedPolicyID        sql.NullString
			auditInfo              sql.NullString
		)

		err := rows.Scan(
			&s.Id,
			&s.PolicyId,
			&renewalDueDate,
			&renewalPremium,
			&renewalPremiumCurrency,
			&renewalTypeStr,
			&statusStr,
			&s.GracePeriodDays,
			&gracePeriodEnd,
			&renewedAt,
			&renewedPolicyID,
			&auditInfo,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan renewal schedule: %w", err)
		}

		// Set Money field
		s.RenewalPremium = &commonv1.Money{Amount: renewalPremium, Currency: renewalPremiumCurrency}

		// Set optional fields
		if renewedPolicyID.Valid {
			s.RenewedPolicyId = renewedPolicyID.String
		}

		// Parse enums
		if renewalTypeStr.Valid {
			k := strings.ToUpper(renewalTypeStr.String)
			if v, ok := renewalv1.RenewalType_value[k]; ok {
				s.RenewalType = renewalv1.RenewalType(v)
			}
		}
		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := renewalv1.RenewalStatus_value[k]; ok {
				s.Status = renewalv1.RenewalStatus(v)
			}
		}

		// Set timestamps
		if !renewalDueDate.IsZero() {
			s.RenewalDueDate = timestamppb.New(renewalDueDate)
		}
		if gracePeriodEnd.Valid {
			s.GracePeriodEnd = timestamppb.New(gracePeriodEnd.Time)
		}
		if renewedAt.Valid {
			s.RenewedAt = timestamppb.New(renewedAt.Time)
		}

		// Set audit info
		if auditInfo.Valid {
			s.AuditInfo = &commonv1.AuditInfo{}
		}

		schedules = append(schedules, &s)
	}

	return schedules, nil
}
