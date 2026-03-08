package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/lib/pq"

	productsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/products/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *productsv1.Product) (*productsv1.Product, error) {
	if product.ProductId == "" {
		return nil, fmt.Errorf("product_id is required")
	}
	
	if product.CreatedBy == "" {
		return nil, fmt.Errorf("created_by is required")
	}
	
	// Extract Money values (stored as paisa/cents in BIGINT)
	basePremium := int64(0)
	basePremiumCurrency := "BDT"
	if product.BasePremium != nil {
		basePremium = product.BasePremium.Amount
		basePremiumCurrency = product.BasePremium.Currency
	}
	
	minSumInsured := int64(0)
	minSumInsuredCurrency := "BDT"
	if product.MinSumInsured != nil {
		minSumInsured = product.MinSumInsured.Amount
		minSumInsuredCurrency = product.MinSumInsured.Currency
	}
	
	maxSumInsured := int64(0)
	maxSumInsuredCurrency := "BDT"
	if product.MaxSumInsured != nil {
		maxSumInsured = product.MaxSumInsured.Amount
		maxSumInsuredCurrency = product.MaxSumInsured.Currency
	}
	
	// Handle product_attributes JSONB
	var productAttrs interface{}
	if product.ProductAttributes != "" {
		productAttrs = product.ProductAttributes
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.products
			(product_id, product_code, product_name, category, description, 
			 base_premium, min_sum_insured, max_sum_insured,
			 base_premium_currency, min_sum_insured_currency, max_sum_insured_currency,
			 min_tenure_months, max_tenure_months, exclusions, status, 
			 created_at, updated_at, created_by, product_attributes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW(), $16, $17)`,
		product.ProductId,
		product.ProductCode,
		product.ProductName,
		strings.ToUpper(product.Category.String()),
		product.Description,
		basePremium,
		minSumInsured,
		maxSumInsured,
		basePremiumCurrency,
		minSumInsuredCurrency,
		maxSumInsuredCurrency,
		product.MinTenureMonths,
		product.MaxTenureMonths,
		pq.Array(product.Exclusions),
		strings.ToUpper(product.Status.String()),
		product.CreatedBy,
		productAttrs,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}

	return r.GetByID(ctx, product.ProductId)
}

func (r *ProductRepository) GetByID(ctx context.Context, productID string) (*productsv1.Product, error) {
	var (
		p                     productsv1.Product
		categoryStr           sql.NullString
		statusStr             sql.NullString
		createdBy             sql.NullString
		basePremium           int64
		minSumInsured         int64
		maxSumInsured         int64
		basePremiumCurrency   string
		minSumInsuredCurrency string
		maxSumInsuredCurrency string
		exclusions            pq.StringArray
		productAttrs          sql.NullString
		createdAt             time.Time
		updatedAt             time.Time
		deletedAt             sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT product_id, product_code, product_name, category, 
		       COALESCE(description, '') as description,
		       base_premium, min_sum_insured, max_sum_insured,
		       base_premium_currency, min_sum_insured_currency, max_sum_insured_currency,
		       min_tenure_months, max_tenure_months, 
		       COALESCE(exclusions, '{}') as exclusions,
		       status, created_at, updated_at, created_by,
		       product_attributes, deleted_at
		FROM insurance_schema.products
		WHERE product_id = $1 AND deleted_at IS NULL`,
		productID,
	).Row().Scan(
		&p.ProductId,
		&p.ProductCode,
		&p.ProductName,
		&categoryStr,
		&p.Description,
		&basePremium,
		&minSumInsured,
		&maxSumInsured,
		&basePremiumCurrency,
		&minSumInsuredCurrency,
		&maxSumInsuredCurrency,
		&p.MinTenureMonths,
		&p.MaxTenureMonths,
		&exclusions,
		&statusStr,
		&createdAt,
		&updatedAt,
		&createdBy,
		&productAttrs,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Set Money fields
	p.BasePremium = &commonv1.Money{
		Amount:   basePremium,
		Currency: basePremiumCurrency,
	}
	p.MinSumInsured = &commonv1.Money{
		Amount:   minSumInsured,
		Currency: minSumInsuredCurrency,
	}
	p.MaxSumInsured = &commonv1.Money{
		Amount:   maxSumInsured,
		Currency: maxSumInsuredCurrency,
	}
	
	// Set currency companion fields
	p.BasePremiumCurrency = basePremiumCurrency
	p.MinSumInsuredCurrency = minSumInsuredCurrency
	p.MaxSumInsuredCurrency = maxSumInsuredCurrency

	// Set exclusions array
	p.Exclusions = exclusions

	// Set product_attributes
	if productAttrs.Valid {
		p.ProductAttributes = productAttrs.String
	}

	// Set created_by if valid
	if createdBy.Valid {
		p.CreatedBy = createdBy.String
	}

	// Parse category enum
	if categoryStr.Valid {
		k := strings.ToUpper(categoryStr.String)
		if v, ok := productsv1.ProductCategory_value[k]; ok {
			p.Category = productsv1.ProductCategory(v)
		}
	}

	// Parse status enum
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := productsv1.ProductStatus_value[k]; ok {
			p.Status = productsv1.ProductStatus(v)
		}
	}

	if !createdAt.IsZero() {
		p.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		p.UpdatedAt = timestamppb.New(updatedAt)
	}
	if deletedAt.Valid {
		p.DeletedAt = timestamppb.New(deletedAt.Time)
	}

	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *productsv1.Product) (*productsv1.Product, error) {
	// Extract Money values
	basePremium := int64(0)
	basePremiumCurrency := "BDT"
	if product.BasePremium != nil {
		basePremium = product.BasePremium.Amount
		basePremiumCurrency = product.BasePremium.Currency
	}
	
	minSumInsured := int64(0)
	minSumInsuredCurrency := "BDT"
	if product.MinSumInsured != nil {
		minSumInsured = product.MinSumInsured.Amount
		minSumInsuredCurrency = product.MinSumInsured.Currency
	}
	
	maxSumInsured := int64(0)
	maxSumInsuredCurrency := "BDT"
	if product.MaxSumInsured != nil {
		maxSumInsured = product.MaxSumInsured.Amount
		maxSumInsuredCurrency = product.MaxSumInsured.Currency
	}
	
	// Handle product_attributes JSONB
	var productAttrs interface{}
	if product.ProductAttributes != "" {
		productAttrs = product.ProductAttributes
	}
	
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.products
		SET product_code = $2,
		    product_name = $3,
		    category = $4,
		    description = $5,
		    base_premium = $6,
		    min_sum_insured = $7,
		    max_sum_insured = $8,
		    base_premium_currency = $9,
		    min_sum_insured_currency = $10,
		    max_sum_insured_currency = $11,
		    min_tenure_months = $12,
		    max_tenure_months = $13,
		    exclusions = $14,
		    status = $15,
		    product_attributes = $16,
		    updated_at = NOW()
		WHERE product_id = $1 AND deleted_at IS NULL`,
		product.ProductId,
		product.ProductCode,
		product.ProductName,
		strings.ToUpper(product.Category.String()),
		product.Description,
		basePremium,
		minSumInsured,
		maxSumInsured,
		basePremiumCurrency,
		minSumInsuredCurrency,
		maxSumInsuredCurrency,
		product.MinTenureMonths,
		product.MaxTenureMonths,
		pq.Array(product.Exclusions),
		strings.ToUpper(product.Status.String()),
		productAttrs,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return r.GetByID(ctx, product.ProductId)
}

func (r *ProductRepository) Delete(ctx context.Context, productID string) error {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.products
		SET deleted_at = NOW()
		WHERE product_id = $1 AND deleted_at IS NULL`,
		productID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context, tenantID string, page, pageSize int) ([]*productsv1.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM insurance_schema.products WHERE deleted_at IS NULL`
	if tenantID != "" {
		countQuery += ` AND tenant_id = $1`
		err := r.db.WithContext(ctx).Raw(countQuery, tenantID).Scan(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count products: %w", err)
		}
	} else {
		err := r.db.WithContext(ctx).Raw(countQuery).Scan(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count products: %w", err)
		}
	}

	// Get products
	query := `
		SELECT product_id, product_code, product_name, category, 
		       COALESCE(description, '') as description,
		       base_premium, min_sum_insured, max_sum_insured,
		       base_premium_currency, min_sum_insured_currency, max_sum_insured_currency,
		       min_tenure_months, max_tenure_months, 
		       COALESCE(exclusions, '{}') as exclusions,
		       status, created_at, updated_at, created_by,
		       product_attributes, deleted_at
		FROM insurance_schema.products
		WHERE deleted_at IS NULL`

	if tenantID != "" {
		query += ` AND tenant_id = $1`
		query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, pageSize, offset)
	} else {
		query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, pageSize, offset)
	}

	var rows *sql.Rows
	var err error
	if tenantID != "" {
		rows, err = r.db.WithContext(ctx).Raw(query, tenantID).Rows()
	} else {
		rows, err = r.db.WithContext(ctx).Raw(query).Rows()
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := make([]*productsv1.Product, 0)
	for rows.Next() {
		var (
			p                     productsv1.Product
			categoryStr           sql.NullString
			statusStr             sql.NullString
			createdBy             sql.NullString
			basePremium           int64
			minSumInsured         int64
			maxSumInsured         int64
			basePremiumCurrency   string
			minSumInsuredCurrency string
			maxSumInsuredCurrency string
			exclusions            pq.StringArray
			productAttrs          sql.NullString
			createdAt             time.Time
			updatedAt             time.Time
			deletedAt             sql.NullTime
		)

		err := rows.Scan(
			&p.ProductId,
			&p.ProductCode,
			&p.ProductName,
			&categoryStr,
			&p.Description,
			&basePremium,
			&minSumInsured,
			&maxSumInsured,
			&basePremiumCurrency,
			&minSumInsuredCurrency,
			&maxSumInsuredCurrency,
			&p.MinTenureMonths,
			&p.MaxTenureMonths,
			&exclusions,
			&statusStr,
			&createdAt,
			&updatedAt,
			&createdBy,
			&productAttrs,
			&deletedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}

		// Set Money fields
		p.BasePremium = &commonv1.Money{
			Amount:   basePremium,
			Currency: basePremiumCurrency,
		}
		p.MinSumInsured = &commonv1.Money{
			Amount:   minSumInsured,
			Currency: minSumInsuredCurrency,
		}
		p.MaxSumInsured = &commonv1.Money{
			Amount:   maxSumInsured,
			Currency: maxSumInsuredCurrency,
		}
		
		// Set currency companion fields
		p.BasePremiumCurrency = basePremiumCurrency
		p.MinSumInsuredCurrency = minSumInsuredCurrency
		p.MaxSumInsuredCurrency = maxSumInsuredCurrency

		// Set exclusions array
		p.Exclusions = exclusions

		// Set product_attributes
		if productAttrs.Valid {
			p.ProductAttributes = productAttrs.String
		}

		// Set created_by if valid
		if createdBy.Valid {
			p.CreatedBy = createdBy.String
		}

		// Parse category enum
		if categoryStr.Valid {
			k := strings.ToUpper(categoryStr.String)
			if v, ok := productsv1.ProductCategory_value[k]; ok {
				p.Category = productsv1.ProductCategory(v)
			}
		}

		// Parse status enum
		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := productsv1.ProductStatus_value[k]; ok {
				p.Status = productsv1.ProductStatus(v)
			}
		}

		if !createdAt.IsZero() {
			p.CreatedAt = timestamppb.New(createdAt)
		}
		if !updatedAt.IsZero() {
			p.UpdatedAt = timestamppb.New(updatedAt)
		}
		if deletedAt.Valid {
			p.DeletedAt = timestamppb.New(deletedAt.Time)
		}

		products = append(products, &p)
	}

	return products, total, nil
}
