package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
)

// ─── SQL ──────────────────────────────────────────────────────────────────────

const catalogCols = `
	pp.product_id,
	p.product_name,
	pp.plan_id,
	pp.plan_name,
	p.category AS insurance_category,
	COALESCE(pp.premium_amount::TEXT, 'null') AS premium_amount
`

// ─── Scanner ──────────────────────────────────────────────────────────────────

func scanCatalogPlan(row interface{ Scan(...any) error }) (*domain.CatalogPlan, error) {
	var (
		productID            string
		productName          string
		planID               string
		planName             string
		insuranceCategoryStr sql.NullString
		premiumJSON          sql.NullString
	)

	if err := row.Scan(
		&productID,
		&productName,
		&planID,
		&planName,
		&insuranceCategoryStr,
		&premiumJSON,
	); err != nil {
		return nil, err
	}

	premium := scanMoney(premiumJSON)

	return &domain.CatalogPlan{
		ProductID:         productID,
		ProductName:       productName,
		PlanID:            planID,
		PlanName:          planName,
		InsuranceCategory: parseInsuranceType(insuranceCategoryStr.String),
		PremiumAmount:     premium,
	}, nil
}

// ─── Queries ──────────────────────────────────────────────────────────────────

func (r *PortalRepository) ListCatalogPlans(ctx context.Context) ([]*domain.CatalogPlan, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM insurance_schema.product_plans AS pp
		JOIN insurance_schema.products AS p ON p.product_id = pp.product_id
		WHERE p.status IN ('ACTIVE', 'PRODUCT_STATUS_ACTIVE', '2')
		ORDER BY p.product_name ASC, pp.plan_name ASC`,
		catalogCols,
	)

	rows, err := r.db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.CatalogPlan
	for rows.Next() {
		item, err := scanCatalogPlan(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PortalRepository) GetCatalogPlansByPlanIDs(ctx context.Context, planIDs []string) (map[string]*domain.CatalogPlan, error) {
	result := make(map[string]*domain.CatalogPlan)
	if len(planIDs) == 0 {
		return result, nil
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM insurance_schema.product_plans AS pp
		JOIN insurance_schema.products AS p ON p.product_id = pp.product_id
		WHERE pp.plan_id = ANY($1)`,
		catalogCols,
	)

	rows, err := r.db.WithContext(ctx).Raw(query, planIDs).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item, err := scanCatalogPlan(rows)
		if err != nil {
			return nil, err
		}
		result[item.PlanID] = item
	}
	return result, rows.Err()
}
