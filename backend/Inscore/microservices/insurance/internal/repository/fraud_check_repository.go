package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	claimsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/claims/entity/v1"
)

type FraudCheckRepository struct {
	db *gorm.DB
}

func NewFraudCheckRepository(db *gorm.DB) *FraudCheckRepository {
	return &FraudCheckRepository{db: db}
}

func (r *FraudCheckRepository) Create(ctx context.Context, check *claimsv1.FraudCheckResult) (*claimsv1.FraudCheckResult, error) {
	if check.FraudCheckId == "" {
		return nil, fmt.Errorf("fraud_check_id is required")
	}
	
	var reviewedAt interface{}
	if check.ReviewedAt != nil {
		reviewedAt = check.ReviewedAt.AsTime()
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.fraud_checks
			(fraud_check_id, claim_id, fraud_score, risk_factors, flagged,
			 reviewed_by, reviewed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
		check.FraudCheckId, check.ClaimId, check.FraudScore,
		pq.Array(check.RiskFactors), check.Flagged, check.ReviewedBy, reviewedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert fraud check: %w", err)
	}

	return r.GetByID(ctx, check.FraudCheckId)
}

func (r *FraudCheckRepository) GetByID(ctx context.Context, fraudCheckID string) (*claimsv1.FraudCheckResult, error) {
	var (
		f           claimsv1.FraudCheckResult
		riskFactors pq.StringArray
		reviewedBy  sql.NullString
		reviewedAt  sql.NullTime
		createdAt   time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT fraud_check_id, claim_id, fraud_score, 
		       COALESCE(risk_factors, '{}') as risk_factors,
		       flagged, reviewed_by, reviewed_at, created_at
		FROM insurance_schema.fraud_checks
		WHERE fraud_check_id = $1`,
		fraudCheckID,
	).Row().Scan(
		&f.FraudCheckId, &f.ClaimId, &f.FraudScore,
		&riskFactors, &f.Flagged, &reviewedBy, &reviewedAt, &createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get fraud check: %w", err)
	}

	f.RiskFactors = riskFactors
	
	if reviewedBy.Valid {
		f.ReviewedBy = reviewedBy.String
	}
	if reviewedAt.Valid {
		f.ReviewedAt = timestamppb.New(reviewedAt.Time)
	}
	if !createdAt.IsZero() {
		f.CreatedAt = timestamppb.New(createdAt)
	}

	return &f, nil
}

func (r *FraudCheckRepository) GetByClaimID(ctx context.Context, claimID string) (*claimsv1.FraudCheckResult, error) {
	var (
		f           claimsv1.FraudCheckResult
		riskFactors pq.StringArray
		reviewedBy  sql.NullString
		reviewedAt  sql.NullTime
		createdAt   time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT fraud_check_id, claim_id, fraud_score, 
		       COALESCE(risk_factors, '{}') as risk_factors,
		       flagged, reviewed_by, reviewed_at, created_at
		FROM insurance_schema.fraud_checks
		WHERE claim_id = $1`,
		claimID,
	).Row().Scan(
		&f.FraudCheckId, &f.ClaimId, &f.FraudScore,
		&riskFactors, &f.Flagged, &reviewedBy, &reviewedAt, &createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get fraud check by claim: %w", err)
	}

	f.RiskFactors = riskFactors
	
	if reviewedBy.Valid {
		f.ReviewedBy = reviewedBy.String
	}
	if reviewedAt.Valid {
		f.ReviewedAt = timestamppb.New(reviewedAt.Time)
	}
	if !createdAt.IsZero() {
		f.CreatedAt = timestamppb.New(createdAt)
	}

	return &f, nil
}

func (r *FraudCheckRepository) Update(ctx context.Context, check *claimsv1.FraudCheckResult) (*claimsv1.FraudCheckResult, error) {
	var reviewedAt interface{}
	if check.ReviewedAt != nil {
		reviewedAt = check.ReviewedAt.AsTime()
	}
	
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.fraud_checks
		SET claim_id = $2, fraud_score = $3, risk_factors = $4, flagged = $5,
		    reviewed_by = $6, reviewed_at = $7
		WHERE fraud_check_id = $1`,
		check.FraudCheckId, check.ClaimId, check.FraudScore,
		pq.Array(check.RiskFactors), check.Flagged, check.ReviewedBy, reviewedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update fraud check: %w", err)
	}

	return r.GetByID(ctx, check.FraudCheckId)
}

func (r *FraudCheckRepository) Delete(ctx context.Context, fraudCheckID string) error {
	return r.db.WithContext(ctx).Exec(`DELETE FROM insurance_schema.fraud_checks WHERE fraud_check_id = $1`, fraudCheckID).Error
}

func (r *FraudCheckRepository) ListFlagged(ctx context.Context, page, pageSize int) ([]*claimsv1.FraudCheckResult, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	err := r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM insurance_schema.fraud_checks WHERE flagged = true`).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count flagged fraud checks: %w", err)
	}

	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT fraud_check_id, claim_id, fraud_score, 
		       COALESCE(risk_factors, '{}') as risk_factors,
		       flagged, reviewed_by, reviewed_at, created_at
		FROM insurance_schema.fraud_checks
		WHERE flagged = true
		ORDER BY fraud_score DESC, created_at DESC
		LIMIT $1 OFFSET $2`,
		pageSize, offset,
	).Rows()

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list flagged fraud checks: %w", err)
	}
	defer rows.Close()

	checks := make([]*claimsv1.FraudCheckResult, 0)
	for rows.Next() {
		var (
			f           claimsv1.FraudCheckResult
			riskFactors pq.StringArray
			reviewedBy  sql.NullString
			reviewedAt  sql.NullTime
			createdAt   time.Time
		)

		err := rows.Scan(
			&f.FraudCheckId, &f.ClaimId, &f.FraudScore,
			&riskFactors, &f.Flagged, &reviewedBy, &reviewedAt, &createdAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fraud check: %w", err)
		}

		f.RiskFactors = riskFactors
		
		if reviewedBy.Valid {
			f.ReviewedBy = reviewedBy.String
		}
		if reviewedAt.Valid {
			f.ReviewedAt = timestamppb.New(reviewedAt.Time)
		}
		if !createdAt.IsZero() {
			f.CreatedAt = timestamppb.New(createdAt)
		}

		checks = append(checks, &f)
	}

	return checks, total, nil
}
