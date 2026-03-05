package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type purchaseOrderRow struct {
	PurchaseOrderID     string    `gorm:"column:purchase_order_id"`
	PurchaseOrderNumber string    `gorm:"column:purchase_order_number"`
	BusinessID          string    `gorm:"column:business_id"`
	DepartmentID        string    `gorm:"column:department_id"`
	ProductID           string    `gorm:"column:product_id"`
	PlanID              string    `gorm:"column:plan_id"`
	InsuranceCategory   string    `gorm:"column:insurance_category"`
	EmployeeCount       int32     `gorm:"column:employee_count"`
	NumberOfDependents  int32     `gorm:"column:number_of_dependents"`
	CoverageAmount      []byte    `gorm:"column:coverage_amount"`
	EstimatedPremium    []byte    `gorm:"column:estimated_premium"`
	Status              string    `gorm:"column:status"`
	RequestedBy         string    `gorm:"column:requested_by"`
	Notes               string    `gorm:"column:notes"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

func purchaseOrderSelectColumns() string {
	return strings.Join([]string{
		"purchase_order_id",
		"purchase_order_number",
		"business_id",
		"department_id",
		"product_id",
		"plan_id",
		"insurance_category",
		"employee_count",
		"number_of_dependents",
		"coverage_amount",
		"estimated_premium",
		"status",
		"requested_by",
		"notes",
		"created_at",
		"updated_at",
	}, ", ")
}

func (r *PortalRepository) ListPurchaseOrders(
	ctx context.Context,
	pageSize, offset int,
	businessID string,
	status b2bv1.PurchaseOrderStatus,
) ([]*b2bv1.PurchaseOrder, int64, error) {
	q := r.db.WithContext(ctx).Table("b2b_schema.purchase_orders")
	if businessID != "" {
		q = q.Where("business_id = ?", businessID)
	}
	if dbStatus := purchaseOrderStatusToDB(status); dbStatus != "" {
		q = q.Where("status = ?", dbStatus)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []purchaseOrderRow
	if err := q.
		Select(purchaseOrderSelectColumns()).
		Order("created_at DESC, purchase_order_number DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]*b2bv1.PurchaseOrder, 0, len(rows))
	for _, row := range rows {
		item, err := mapPurchaseOrderRow(row)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (r *PortalRepository) GetPurchaseOrder(ctx context.Context, purchaseOrderID string) (*b2bv1.PurchaseOrder, error) {
	var row purchaseOrderRow
	err := r.db.WithContext(ctx).
		Table("b2b_schema.purchase_orders").
		Select(purchaseOrderSelectColumns()).
		Where("purchase_order_id = ?", purchaseOrderID).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return mapPurchaseOrderRow(row)
}

func (r *PortalRepository) CreatePurchaseOrder(
	ctx context.Context,
	input domain.PurchaseOrderCreateInput,
) (*b2bv1.PurchaseOrder, error) {
	coverageAmount, err := mustMarshalMoney(input.CoverageAmount)
	if err != nil {
		return nil, fmt.Errorf("marshal coverage_amount: %w", err)
	}
	estimatedPremium, err := mustMarshalMoney(input.EstimatedPremium)
	if err != nil {
		return nil, fmt.Errorf("marshal estimated_premium: %w", err)
	}

	values := map[string]any{
		"purchase_order_id":     input.PurchaseOrderID,
		"purchase_order_number": input.PurchaseOrderNumber,
		"business_id":           input.BusinessID,
		"department_id":         input.DepartmentID,
		"product_id":            input.ProductID,
		"plan_id":               input.PlanID,
		"insurance_category":    input.InsuranceCategory.String(),
		"employee_count":        input.EmployeeCount,
		"number_of_dependents":  input.NumberOfDependents,
		"coverage_amount":       string(coverageAmount),
		"estimated_premium":     string(estimatedPremium),
		"status":                purchaseOrderStatusToDB(input.Status),
		"requested_by":          input.RequestedBy,
		"notes":                 input.Notes,
	}

	if err := r.db.WithContext(ctx).
		Table("b2b_schema.purchase_orders").
		Create(values).Error; err != nil {
		return nil, err
	}

	return r.GetPurchaseOrder(ctx, input.PurchaseOrderID)
}

func mapPurchaseOrderRow(row purchaseOrderRow) (*b2bv1.PurchaseOrder, error) {
	coverageAmount, err := parseMoney(row.CoverageAmount)
	if err != nil {
		return nil, fmt.Errorf("parse coverage_amount for purchase order %s: %w", row.PurchaseOrderID, err)
	}
	estimatedPremium, err := parseMoney(row.EstimatedPremium)
	if err != nil {
		return nil, fmt.Errorf("parse estimated_premium for purchase order %s: %w", row.PurchaseOrderID, err)
	}

	item := &b2bv1.PurchaseOrder{
		PurchaseOrderId:     row.PurchaseOrderID,
		PurchaseOrderNumber: row.PurchaseOrderNumber,
		BusinessId:          row.BusinessID,
		DepartmentId:        row.DepartmentID,
		ProductId:           row.ProductID,
		PlanId:              row.PlanID,
		InsuranceCategory:   parseInsuranceType(row.InsuranceCategory),
		EmployeeCount:       row.EmployeeCount,
		NumberOfDependents:  row.NumberOfDependents,
		CoverageAmount:      coverageAmount,
		EstimatedPremium:    estimatedPremium,
		Status:              parsePurchaseOrderStatus(row.Status),
		RequestedBy:         row.RequestedBy,
		Notes:               row.Notes,
	}

	if !row.CreatedAt.IsZero() {
		item.CreatedAt = timestamppb.New(row.CreatedAt)
	}
	if !row.UpdatedAt.IsZero() {
		item.UpdatedAt = timestamppb.New(row.UpdatedAt)
	}

	return item, nil
}
