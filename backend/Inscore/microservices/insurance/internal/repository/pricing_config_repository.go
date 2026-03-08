package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	productsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/products/entity/v1"
)

type PricingConfigRepository struct {
	db *gorm.DB
}

func NewPricingConfigRepository(db *gorm.DB) *PricingConfigRepository {
	return &PricingConfigRepository{db: db}
}

func (r *PricingConfigRepository) Create(ctx context.Context, config *productsv1.PricingConfig) (*productsv1.PricingConfig, error) {
	if config.PricingConfigId == "" {
		return nil, fmt.Errorf("pricing_config_id is required")
	}
	
	if config.ProductId == "" {
		return nil, fmt.Errorf("product_id is required")
	}
	
	// Serialize rules to JSON
	rulesJSON, err := json.Marshal(config.Rules)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules: %w", err)
	}
	
	// Handle timestamps
	var effectiveFrom, effectiveTo interface{}
	if config.EffectiveFrom != nil {
		effectiveFrom = config.EffectiveFrom.AsTime()
	}
	if config.EffectiveTo != nil {
		effectiveTo = config.EffectiveTo.AsTime()
	}
	
	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.pricing_configs
			(pricing_config_id, product_id, rules, effective_from, effective_to,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		config.PricingConfigId,
		config.ProductId,
		rulesJSON,
		effectiveFrom,
		effectiveTo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert pricing config: %w", err)
	}

	return r.GetByID(ctx, config.PricingConfigId)
}

func (r *PricingConfigRepository) GetByID(ctx context.Context, configID string) (*productsv1.PricingConfig, error) {
	var (
		pc            productsv1.PricingConfig
		rulesJSON     []byte
		effectiveFrom time.Time
		effectiveTo   sql.NullTime
		createdAt     time.Time
		updatedAt     time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT pricing_config_id, product_id, rules, 
		       effective_from, effective_to,
		       created_at, updated_at
		FROM insurance_schema.pricing_configs
		WHERE pricing_config_id = $1`,
		configID,
	).Row().Scan(
		&pc.PricingConfigId,
		&pc.ProductId,
		&rulesJSON,
		&effectiveFrom,
		&effectiveTo,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get pricing config: %w", err)
	}

	// Deserialize rules from JSON
	var rules []*productsv1.PricingRule
	if len(rulesJSON) > 0 {
		err = json.Unmarshal(rulesJSON, &rules)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
		}
		pc.Rules = rules
	}

	if !effectiveFrom.IsZero() {
		pc.EffectiveFrom = timestamppb.New(effectiveFrom)
	}
	if effectiveTo.Valid {
		pc.EffectiveTo = timestamppb.New(effectiveTo.Time)
	}
	if !createdAt.IsZero() {
		pc.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		pc.UpdatedAt = timestamppb.New(updatedAt)
	}

	return &pc, nil
}
