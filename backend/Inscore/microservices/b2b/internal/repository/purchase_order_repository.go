package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ─── SQL ──────────────────────────────────────────────────────────────────────

const purchaseOrderCols = `
	purchase_order_id,
	purchase_order_number,
	business_id,
	department_id,
	product_id,
	plan_id,
	insurance_category,
	COALESCE(employee_count, 0) AS employee_count,
	COALESCE(number_of_dependents, 0) AS number_of_dependents,
	COALESCE(coverage_amount::TEXT, 'null') AS coverage_amount,
	COALESCE(estimated_premium::TEXT, 'null') AS estimated_premium,
	status,
	COALESCE(requested_by::TEXT, '') AS requested_by,
	COALESCE(notes, '') AS notes,
	created_at,
	updated_at
`

// ─── Scanner ──────────────────────────────────────────────────────────────────

func scanPurchaseOrder(row interface{ Scan(...any) error }) (*b2bv1.PurchaseOrder, error) {
	var (
		o                    b2bv1.PurchaseOrder
		insuranceCategoryStr sql.NullString
		coverageJSON         sql.NullString
		estimatedPremiumJSON sql.NullString
		statusStr            sql.NullString
		createdAt            time.Time
		updatedAt            time.Time
	)

	if err := row.Scan(
		&o.PurchaseOrderId,
		&o.PurchaseOrderNumber,
		&o.BusinessId,
		&o.DepartmentId,
		&o.ProductId,
		&o.PlanId,
		&insuranceCategoryStr,
		&o.EmployeeCount,
		&o.NumberOfDependents,
		&coverageJSON,
		&estimatedPremiumJSON,
		&statusStr,
		&o.RequestedBy,
		&o.Notes,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if insuranceCategoryStr.Valid {
		o.InsuranceCategory = parseInsuranceType(insuranceCategoryStr.String)
	}
	o.CoverageAmount = scanMoney(coverageJSON)
	o.EstimatedPremium = scanMoney(estimatedPremiumJSON)

	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := b2bv1.PurchaseOrderStatus_value[k]; ok {
			o.Status = b2bv1.PurchaseOrderStatus(v)
		} else if v, ok := b2bv1.PurchaseOrderStatus_value["PURCHASE_ORDER_STATUS_"+k]; ok {
			o.Status = b2bv1.PurchaseOrderStatus(v)
		}
	}

	if !createdAt.IsZero() { o.CreatedAt = timestamppb.New(createdAt) }
	if !updatedAt.IsZero() { o.UpdatedAt = timestamppb.New(updatedAt) }
	return &o, nil
}

// ─── Queries ──────────────────────────────────────────────────────────────────

func (r *PortalRepository) GetPurchaseOrder(ctx context.Context, purchaseOrderID string) (*b2bv1.PurchaseOrder, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM b2b_schema.purchase_orders WHERE purchase_order_id = $1 AND deleted_at IS NULL LIMIT 1`,
		purchaseOrderCols,
	)
	row := r.db.WithContext(ctx).Raw(query, purchaseOrderID).Row()
	return scanPurchaseOrder(row)
}

func (r *PortalRepository) ListPurchaseOrders(
	ctx context.Context,
	pageSize, offset int,
	businessID string,
	status b2bv1.PurchaseOrderStatus,
) ([]*b2bv1.PurchaseOrder, int64, error) {
	where := "deleted_at IS NULL"
	args := []interface{}{}
	idx := 1

	if businessID != "" {
		where += fmt.Sprintf(" AND business_id = $%d", idx)
		args = append(args, businessID)
		idx++
	}
	if status != b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED {
		where += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, purchaseOrderStatusStr(status))
		idx++
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	if err := r.db.WithContext(ctx).Raw(
		fmt.Sprintf("SELECT COUNT(*) FROM b2b_schema.purchase_orders WHERE %s", where),
		countArgs...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(
		`SELECT %s FROM b2b_schema.purchase_orders WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		purchaseOrderCols, where, idx, idx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*b2bv1.PurchaseOrder
	for rows.Next() {
		o, err := scanPurchaseOrder(rows)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}
	return orders, total, rows.Err()
}

func (r *PortalRepository) CreatePurchaseOrder(ctx context.Context, input domain.PurchaseOrderCreateInput) (*b2bv1.PurchaseOrder, error) {
	id := input.PurchaseOrderID
	if id == "" {
		id = newUUID()
	}

	coverageJSONBytes, err := marshalMoney(input.CoverageAmount)
	if err != nil {
		return nil, fmt.Errorf("marshal coverage_amount: %w", err)
	}
	premiumJSONBytes, err := marshalMoney(input.EstimatedPremium)
	if err != nil {
		return nil, fmt.Errorf("marshal estimated_premium: %w", err)
	}
	// Pass as string (not []byte) so PostgreSQL receives valid JSON text for JSONB columns.
	coverageJSON := string(coverageJSONBytes)
	premiumJSON := string(premiumJSONBytes)

	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO b2b_schema.purchase_orders (
			purchase_order_id, purchase_order_number, business_id, department_id,
			product_id, plan_id, insurance_category,
			employee_count, number_of_dependents,
			coverage_amount, estimated_premium,
			status, requested_by, notes
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9,
			$10, $11,
			$12, $13, $14
		)`,
		id,
		input.PurchaseOrderNumber,
		input.BusinessID,
		input.DepartmentID,
		input.ProductID,
		input.PlanID,
		input.InsuranceCategory.String(),
		input.EmployeeCount,
		input.NumberOfDependents,
		coverageJSON,
		premiumJSON,
		purchaseOrderStatusStr(input.Status),
		nullableStr(input.RequestedBy),
		nullableStr(input.Notes),
	).Error; err != nil {
		return nil, fmt.Errorf("insert purchase_order: %w", err)
	}

	return r.GetPurchaseOrder(ctx, id)
}
