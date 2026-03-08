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

type GracePeriodRepository struct {
	db *gorm.DB
}

func NewGracePeriodRepository(db *gorm.DB) *GracePeriodRepository {
	return &GracePeriodRepository{db: db}
}

func (r *GracePeriodRepository) Create(ctx context.Context, gracePeriod *renewalv1.GracePeriod) (*renewalv1.GracePeriod, error) {
	if gracePeriod.Id == "" {
		return nil, fmt.Errorf("grace_period_id is required")
	}

	// Handle timestamps
	var startDate time.Time
	if gracePeriod.StartDate != nil {
		startDate = gracePeriod.StartDate.AsTime()
	}

	var endDate time.Time
	if gracePeriod.EndDate != nil {
		endDate = gracePeriod.EndDate.AsTime()
	}

	var revivedAt sql.NullTime
	if gracePeriod.RevivedAt != nil {
		revivedAt = sql.NullTime{Time: gracePeriod.RevivedAt.AsTime(), Valid: true}
	}

	var auditInfo interface{}
	if gracePeriod.AuditInfo != nil {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.grace_periods
			(grace_period_id, policy_id, start_date, end_date, days_remaining,
			 status, coverage_active, revived_at, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		gracePeriod.Id,
		gracePeriod.PolicyId,
		startDate,
		endDate,
		gracePeriod.DaysRemaining,
		strings.ToUpper(gracePeriod.Status.String()),
		gracePeriod.CoverageActive,
		revivedAt,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert grace period: %w", err)
	}

	return r.GetByID(ctx, gracePeriod.Id)
}

func (r *GracePeriodRepository) GetByID(ctx context.Context, gracePeriodID string) (*renewalv1.GracePeriod, error) {
	var (
		gp        renewalv1.GracePeriod
		statusStr string
		startDate time.Time
		endDate   time.Time
		revivedAt sql.NullTime
		auditInfo sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT grace_period_id, policy_id, start_date, end_date, days_remaining,
		       status, coverage_active, revived_at, audit_info
		FROM insurance_schema.grace_periods
		WHERE grace_period_id = $1`,
		gracePeriodID,
	).Row().Scan(
		&gp.Id,
		&gp.PolicyId,
		&startDate,
		&endDate,
		&gp.DaysRemaining,
		&statusStr,
		&gp.CoverageActive,
		&revivedAt,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get grace period: %w", err)
	}

	// Parse enum
	k := strings.ToUpper(statusStr)
	if v, ok := renewalv1.GracePeriodStatus_value[k]; ok {
		gp.Status = renewalv1.GracePeriodStatus(v)
	}

	// Set timestamps
	if !startDate.IsZero() {
		gp.StartDate = timestamppb.New(startDate)
	}
	if !endDate.IsZero() {
		gp.EndDate = timestamppb.New(endDate)
	}
	if revivedAt.Valid {
		gp.RevivedAt = timestamppb.New(revivedAt.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		gp.AuditInfo = &commonv1.AuditInfo{}
	}

	return &gp, nil
}

func (r *GracePeriodRepository) Update(ctx context.Context, gracePeriod *renewalv1.GracePeriod) (*renewalv1.GracePeriod, error) {
	// Handle timestamps
	var startDate time.Time
	if gracePeriod.StartDate != nil {
		startDate = gracePeriod.StartDate.AsTime()
	}

	var endDate time.Time
	if gracePeriod.EndDate != nil {
		endDate = gracePeriod.EndDate.AsTime()
	}

	var revivedAt sql.NullTime
	if gracePeriod.RevivedAt != nil {
		revivedAt = sql.NullTime{Time: gracePeriod.RevivedAt.AsTime(), Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.grace_periods
		SET policy_id = $2,
		    start_date = $3,
		    end_date = $4,
		    days_remaining = $5,
		    status = $6,
		    coverage_active = $7,
		    revived_at = $8
		WHERE grace_period_id = $1`,
		gracePeriod.Id,
		gracePeriod.PolicyId,
		startDate,
		endDate,
		gracePeriod.DaysRemaining,
		strings.ToUpper(gracePeriod.Status.String()),
		gracePeriod.CoverageActive,
		revivedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update grace period: %w", err)
	}

	return r.GetByID(ctx, gracePeriod.Id)
}

func (r *GracePeriodRepository) Delete(ctx context.Context, gracePeriodID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.grace_periods
		WHERE grace_period_id = $1`,
		gracePeriodID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete grace period: %w", err)
	}

	return nil
}

func (r *GracePeriodRepository) GetByPolicyID(ctx context.Context, policyID string) (*renewalv1.GracePeriod, error) {
	var (
		gp        renewalv1.GracePeriod
		statusStr string
		startDate time.Time
		endDate   time.Time
		revivedAt sql.NullTime
		auditInfo sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT grace_period_id, policy_id, start_date, end_date, days_remaining,
		       status, coverage_active, revived_at, audit_info
		FROM insurance_schema.grace_periods
		WHERE policy_id = $1
		ORDER BY start_date DESC
		LIMIT 1`,
		policyID,
	).Row().Scan(
		&gp.Id,
		&gp.PolicyId,
		&startDate,
		&endDate,
		&gp.DaysRemaining,
		&statusStr,
		&gp.CoverageActive,
		&revivedAt,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get grace period by policy_id: %w", err)
	}

	// Parse enum
	k := strings.ToUpper(statusStr)
	if v, ok := renewalv1.GracePeriodStatus_value[k]; ok {
		gp.Status = renewalv1.GracePeriodStatus(v)
	}

	// Set timestamps
	if !startDate.IsZero() {
		gp.StartDate = timestamppb.New(startDate)
	}
	if !endDate.IsZero() {
		gp.EndDate = timestamppb.New(endDate)
	}
	if revivedAt.Valid {
		gp.RevivedAt = timestamppb.New(revivedAt.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		gp.AuditInfo = &commonv1.AuditInfo{}
	}

	return &gp, nil
}

func (r *GracePeriodRepository) ListActive(ctx context.Context) ([]*renewalv1.GracePeriod, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT grace_period_id, policy_id, start_date, end_date, days_remaining,
		       status, coverage_active, revived_at, audit_info
		FROM insurance_schema.grace_periods
		WHERE status = 'GRACE_PERIOD_STATUS_ACTIVE'
		ORDER BY end_date ASC`).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list active grace periods: %w", err)
	}
	defer rows.Close()

	gracePeriods := make([]*renewalv1.GracePeriod, 0)
	for rows.Next() {
		var (
			gp        renewalv1.GracePeriod
			statusStr string
			startDate time.Time
			endDate   time.Time
			revivedAt sql.NullTime
			auditInfo sql.NullString
		)

		err := rows.Scan(
			&gp.Id,
			&gp.PolicyId,
			&startDate,
			&endDate,
			&gp.DaysRemaining,
			&statusStr,
			&gp.CoverageActive,
			&revivedAt,
			&auditInfo,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan grace period: %w", err)
		}

		// Parse enum
		k := strings.ToUpper(statusStr)
		if v, ok := renewalv1.GracePeriodStatus_value[k]; ok {
			gp.Status = renewalv1.GracePeriodStatus(v)
		}

		// Set timestamps
		if !startDate.IsZero() {
			gp.StartDate = timestamppb.New(startDate)
		}
		if !endDate.IsZero() {
			gp.EndDate = timestamppb.New(endDate)
		}
		if revivedAt.Valid {
			gp.RevivedAt = timestamppb.New(revivedAt.Time)
		}

		// Set audit info
		if auditInfo.Valid {
			gp.AuditInfo = &commonv1.AuditInfo{}
		}

		gracePeriods = append(gracePeriods, &gp)
	}

	return gracePeriods, nil
}
