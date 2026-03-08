package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	productsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/products/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type RiderRepository struct {
	db *gorm.DB
}

func NewRiderRepository(db *gorm.DB) *RiderRepository {
	return &RiderRepository{db: db}
}

func (r *RiderRepository) Create(ctx context.Context, rider *productsv1.Rider) (*productsv1.Rider, error) {
	if rider.RiderId == "" {
		return nil, fmt.Errorf("rider_id is required")
	}
	
	if rider.ProductId == "" {
		return nil, fmt.Errorf("product_id is required")
	}
	
	// Extract Money values
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
		INSERT INTO insurance_schema.product_riders
			(rider_id, product_id, rider_name, description,
			 premium_amount, coverage_amount, is_mandatory,
			 premium_currency, coverage_currency,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`,
		rider.RiderId,
		rider.ProductId,
		rider.RiderName,
		rider.Description,
		premiumAmount,
		coverageAmount,
		rider.IsMandatory,
		premiumCurrency,
		coverageCurrency,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert rider: %w", err)
	}

	return r.GetByID(ctx, rider.RiderId)
}

func (r *RiderRepository) GetByID(ctx context.Context, riderID string) (*productsv1.Rider, error) {
	var (
		rider            productsv1.Rider
		premiumAmount    int64
		coverageAmount   int64
		premiumCurrency  string
		coverageCurrency string
		createdAt        time.Time
		updatedAt        time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT rider_id, product_id, rider_name, 
		       COALESCE(description, '') as description,
		       premium_amount, coverage_amount, is_mandatory,
		       premium_currency, coverage_currency,
		       created_at, updated_at
		FROM insurance_schema.product_riders
		WHERE rider_id = $1`,
		riderID,
	).Row().Scan(
		&rider.RiderId,
		&rider.ProductId,
		&rider.RiderName,
		&rider.Description,
		&premiumAmount,
		&coverageAmount,
		&rider.IsMandatory,
		&premiumCurrency,
		&coverageCurrency,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get rider: %w", err)
	}

	// Set Money fields
	rider.PremiumAmount = &commonv1.Money{
		Amount:   premiumAmount,
		Currency: premiumCurrency,
	}
	rider.CoverageAmount = &commonv1.Money{
		Amount:   coverageAmount,
		Currency: coverageCurrency,
	}
	
	// Set currency companion fields
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

func (r *RiderRepository) ListByProductID(ctx context.Context, productID string) ([]*productsv1.Rider, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT rider_id, product_id, rider_name, 
		       COALESCE(description, '') as description,
		       premium_amount, coverage_amount, is_mandatory,
		       premium_currency, coverage_currency,
		       created_at, updated_at
		FROM insurance_schema.product_riders
		WHERE product_id = $1
		ORDER BY created_at DESC`,
		productID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list riders: %w", err)
	}
	defer rows.Close()

	riders := make([]*productsv1.Rider, 0)
	for rows.Next() {
		var (
			rider            productsv1.Rider
			premiumAmount    int64
			coverageAmount   int64
			premiumCurrency  string
			coverageCurrency string
			createdAt        time.Time
			updatedAt        time.Time
		)

		err := rows.Scan(
			&rider.RiderId,
			&rider.ProductId,
			&rider.RiderName,
			&rider.Description,
			&premiumAmount,
			&coverageAmount,
			&rider.IsMandatory,
			&premiumCurrency,
			&coverageCurrency,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan rider: %w", err)
		}

		// Set Money fields
		rider.PremiumAmount = &commonv1.Money{
			Amount:   premiumAmount,
			Currency: premiumCurrency,
		}
		rider.CoverageAmount = &commonv1.Money{
			Amount:   coverageAmount,
			Currency: coverageCurrency,
		}
		
		// Set currency companion fields
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
