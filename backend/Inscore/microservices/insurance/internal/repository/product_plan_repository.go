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

type ProductPlanRepository struct {
	db *gorm.DB
}

func NewProductPlanRepository(db *gorm.DB) *ProductPlanRepository {
	return &ProductPlanRepository{db: db}
}

func (r *ProductPlanRepository) Create(ctx context.Context, plan *productsv1.ProductPlan) (*productsv1.ProductPlan, error) {
	if plan.PlanId == "" {
		return nil, fmt.Errorf("plan_id is required")
	}
	
	if plan.ProductId == "" {
		return nil, fmt.Errorf("product_id is required")
	}
	
	// Extract Money values (currency not stored separately in DB)
	premiumAmount := int64(0)
	if plan.PremiumAmount != nil {
		premiumAmount = plan.PremiumAmount.Amount
	}
	
	minSumInsured := int64(0)
	if plan.MinSumInsured != nil {
		minSumInsured = plan.MinSumInsured.Amount
	}
	
	maxSumInsured := int64(0)
	if plan.MaxSumInsured != nil {
		maxSumInsured = plan.MaxSumInsured.Amount
	}
	
	// Handle attributes JSONB
	var attributes interface{}
	if plan.Attributes != "" {
		attributes = plan.Attributes
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.product_plans
			(plan_id, product_id, plan_name, plan_description,
			 premium_amount, min_sum_insured, max_sum_insured,
			 attributes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`,
		plan.PlanId,
		plan.ProductId,
		plan.PlanName,
		plan.PlanDescription,
		premiumAmount,
		minSumInsured,
		maxSumInsured,
		attributes,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert product plan: %w", err)
	}

	return r.GetByID(ctx, plan.PlanId)
}

func (r *ProductPlanRepository) GetByID(ctx context.Context, planID string) (*productsv1.ProductPlan, error) {
	var (
		p                     productsv1.ProductPlan
		premiumAmount         int64
		minSumInsured         int64
		maxSumInsured         int64
		attributes            sql.NullString
		createdAt             time.Time
		updatedAt             time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT plan_id, product_id, plan_name, 
		       COALESCE(plan_description, '') as plan_description,
		       premium_amount, min_sum_insured, max_sum_insured,
		       attributes, created_at, updated_at
		FROM insurance_schema.product_plans
		WHERE plan_id = $1`,
		planID,
	).Row().Scan(
		&p.PlanId,
		&p.ProductId,
		&p.PlanName,
		&p.PlanDescription,
		&premiumAmount,
		&minSumInsured,
		&maxSumInsured,
		&attributes,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get product plan: %w", err)
	}

	// Set Money fields (currency defaults to BDT since not stored separately)
	p.PremiumAmount = &commonv1.Money{
		Amount:   premiumAmount,
		Currency: "BDT",
	}
	p.MinSumInsured = &commonv1.Money{
		Amount:   minSumInsured,
		Currency: "BDT",
	}
	p.MaxSumInsured = &commonv1.Money{
		Amount:   maxSumInsured,
		Currency: "BDT",
	}

	// Set attributes
	if attributes.Valid {
		p.Attributes = attributes.String
	}

	if !createdAt.IsZero() {
		p.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		p.UpdatedAt = timestamppb.New(updatedAt)
	}

	return &p, nil
}

func (r *ProductPlanRepository) ListByProductID(ctx context.Context, productID string) ([]*productsv1.ProductPlan, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT plan_id, product_id, plan_name, 
		       COALESCE(plan_description, '') as plan_description,
		       premium_amount, min_sum_insured, max_sum_insured,
		       attributes, created_at, updated_at
		FROM insurance_schema.product_plans
		WHERE product_id = $1
		ORDER BY created_at DESC`,
		productID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list product plans: %w", err)
	}
	defer rows.Close()

	plans := make([]*productsv1.ProductPlan, 0)
	for rows.Next() {
		var (
			p                     productsv1.ProductPlan
			premiumAmount         int64
			minSumInsured         int64
			maxSumInsured         int64
			attributes            sql.NullString
			createdAt             time.Time
			updatedAt             time.Time
		)

		err := rows.Scan(
			&p.PlanId,
			&p.ProductId,
			&p.PlanName,
			&p.PlanDescription,
			&premiumAmount,
			&minSumInsured,
			&maxSumInsured,
			&attributes,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan product plan: %w", err)
		}

		// Set Money fields (currency defaults to BDT)
		p.PremiumAmount = &commonv1.Money{
			Amount:   premiumAmount,
			Currency: "BDT",
		}
		p.MinSumInsured = &commonv1.Money{
			Amount:   minSumInsured,
			Currency: "BDT",
		}
		p.MaxSumInsured = &commonv1.Money{
			Amount:   maxSumInsured,
			Currency: "BDT",
		}

		// Set attributes
		if attributes.Valid {
			p.Attributes = attributes.String
		}

		if !createdAt.IsZero() {
			p.CreatedAt = timestamppb.New(createdAt)
		}
		if !updatedAt.IsZero() {
			p.UpdatedAt = timestamppb.New(updatedAt)
		}

		plans = append(plans, &p)
	}

	return plans, nil
}
