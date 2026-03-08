package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	policyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/policy/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type PolicyRiderRepository struct {
	db *gorm.DB
}

func NewPolicyRiderRepository(db *gorm.DB) *PolicyRiderRepository {
	return &PolicyRiderRepository{db: db}
}

func (r *PolicyRiderRepository) Create(ctx context.Context, rider *policyv1.Rider) (*policyv1.Rider, error) {
	if rider.RiderId == "" {
		return nil, fmt.Errorf("rider_id is required")
	}
	
	premiumAmount := int64(0)
	premiumCurrency := "BDT"
	if rider.PremiumAmount != nil {
		premiumAmount = rider.PremiumAmount.Amount
		premiumCurrency = rider.PremiumAmount.Currency
	}
	
	coverageAmount := int64(0)
	coverageCurrency := "BDT"
	if rider.CoverageAmount != nil {
		coverageAmount = rider.CoverageAmount.Amount
		coverageCurrency = rider.CoverageAmount.Currency
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.policy_riders
			(rider_id, policy_id, rider_name, premium_amount, coverage_amount,
			 premium_currency, coverage_currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`,
		rider.RiderId,
		rider.PolicyId,
		rider.RiderName,
		premiumAmount,
		coverageAmount,
		premiumCurrency,
		coverageCurrency,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert policy rider: %w", err)
	}

	return r.GetByID(ctx, rider.RiderId)
}

func (r *PolicyRiderRepository) GetByID(ctx context.Context, riderID string) (*policyv1.Rider, error) {
	var (
		rider            policyv1.Rider
		premiumAmount    int64
		coverageAmount   int64
		premiumCurrency  string
		coverageCurrency string
		createdAt        time.Time
		updatedAt        time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT rider_id, policy_id, rider_name, premium_amount, coverage_amount,
		       premium_currency, coverage_currency, created_at, updated_at
		FROM insurance_schema.policy_riders
		WHERE rider_id = $1`,
		riderID,
	).Row().Scan(
		&rider.RiderId,
		&rider.PolicyId,
		&rider.RiderName,
		&premiumAmount,
		&coverageAmount,
		&premiumCurrency,
		&coverageCurrency,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get policy rider: %w", err)
	}

	rider.PremiumAmount = &commonv1.Money{Amount: premiumAmount, Currency: premiumCurrency}
	rider.CoverageAmount = &commonv1.Money{Amount: coverageAmount, Currency: coverageCurrency}
	rider.PremiumCurrency = premiumCurrency
	rider.CoverageCurrency = coverageCurrency

	if !createdAt.IsZero() {
		rider.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		rider.UpdatedAt = timestamppb.New(updatedAt)
	}

	return &rider, nil
}

func (r *PolicyRiderRepository) Update(ctx context.Context, rider *policyv1.Rider) (*policyv1.Rider, error) {
	premiumAmount := int64(0)
	premiumCurrency := "BDT"
	if rider.PremiumAmount != nil {
		premiumAmount = rider.PremiumAmount.Amount
		premiumCurrency = rider.PremiumAmount.Currency
	}
	
	coverageAmount := int64(0)
	coverageCurrency := "BDT"
	if rider.CoverageAmount != nil {
		coverageAmount = rider.CoverageAmount.Amount
		coverageCurrency = rider.CoverageAmount.Currency
	}
	
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.policy_riders
		SET policy_id = $2, rider_name = $3, premium_amount = $4, coverage_amount = $5,
		    premium_currency = $6, coverage_currency = $7, updated_at = NOW()
		WHERE rider_id = $1`,
		rider.RiderId, rider.PolicyId, rider.RiderName,
		premiumAmount, coverageAmount, premiumCurrency, coverageCurrency,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update policy rider: %w", err)
	}

	return r.GetByID(ctx, rider.RiderId)
}

func (r *PolicyRiderRepository) Delete(ctx context.Context, riderID string) error {
	return r.db.WithContext(ctx).Exec(`DELETE FROM insurance_schema.policy_riders WHERE rider_id = $1`, riderID).Error
}

func (r *PolicyRiderRepository) ListByPolicyID(ctx context.Context, policyID string) ([]*policyv1.Rider, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT rider_id, policy_id, rider_name, premium_amount, coverage_amount,
		       premium_currency, coverage_currency, created_at, updated_at
		FROM insurance_schema.policy_riders WHERE policy_id = $1 ORDER BY created_at DESC`, policyID).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list policy riders: %w", err)
	}
	defer rows.Close()

	riders := make([]*policyv1.Rider, 0)
	for rows.Next() {
		var (
			rider            policyv1.Rider
			premiumAmount    int64
			coverageAmount   int64
			premiumCurrency  string
			coverageCurrency string
			createdAt        time.Time
			updatedAt        time.Time
		)

		err := rows.Scan(&rider.RiderId, &rider.PolicyId, &rider.RiderName,
			&premiumAmount, &coverageAmount, &premiumCurrency, &coverageCurrency,
			&createdAt, &updatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan policy rider: %w", err)
		}

		rider.PremiumAmount = &commonv1.Money{Amount: premiumAmount, Currency: premiumCurrency}
		rider.CoverageAmount = &commonv1.Money{Amount: coverageAmount, Currency: coverageCurrency}
		rider.PremiumCurrency = premiumCurrency
		rider.CoverageCurrency = coverageCurrency

		if !createdAt.IsZero() {
			rider.CreatedAt = timestamppb.New(createdAt)
		}
		if !updatedAt.IsZero() {
			rider.UpdatedAt = timestamppb.New(updatedAt)
		}

		riders = append(riders, &rider)
	}

	return riders, nil
}
