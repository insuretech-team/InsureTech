// Package repository provides PostgreSQL data access for the orders microservice.
//
// Hybrid pattern:
//   - INSERT: map[string]any with raw Go types (string for enums, int64 for money, time.Time for timestamps)
//   - READ:   raw SQL with orderScanRow (plain struct without proto fields)
//             then convert with scanRowToProto() to construct proto Money and timestamps
package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const ordersTable = "insurance_schema.orders"

// OrderRepositoryImpl satisfies domain.OrderRepository using GORM.
type OrderRepositoryImpl struct {
	db *gorm.DB
}

var _ domain.OrderRepository = (*OrderRepositoryImpl)(nil)

func NewOrderRepository(db *gorm.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{db: db}
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// orderScanRow is a plain struct for raw SQL scanning, without proto fields.
// Covers all columns including phase-2 additions (fields 19–32).
type orderScanRow struct {
	// Original fields
	OrderId            string     `gorm:"column:order_id"`
	OrderNumber        string     `gorm:"column:order_number"`
	TenantId           string     `gorm:"column:tenant_id"`
	QuotationId        string     `gorm:"column:quotation_id"`
	CustomerId         string     `gorm:"column:customer_id"`
	ProductId          string     `gorm:"column:product_id"`
	PlanId             string     `gorm:"column:plan_id"`
	Status             string     `gorm:"column:status"`
	TotalPayable       int64      `gorm:"column:total_payable"`
	Currency           string     `gorm:"column:currency"`
	PaymentId          string     `gorm:"column:payment_id"`
	PaymentGatewayRef  string     `gorm:"column:payment_gateway_ref"`
	PolicyId           string     `gorm:"column:policy_id"`
	CancellationReason string     `gorm:"column:cancellation_reason"`
	FailureReason      string     `gorm:"column:failure_reason"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
	PaidAt             *time.Time `gorm:"column:paid_at"`
	// Phase-2 extended fields (fields 19–32)
	InvoiceId          string     `gorm:"column:invoice_id"`
	OrganisationId     string     `gorm:"column:organisation_id"`
	IdempotencyKey     string     `gorm:"column:idempotency_key"`
	CorrelationId      string     `gorm:"column:correlation_id"`
	PaymentStatus      string     `gorm:"column:payment_status"`
	BillingStatus      string     `gorm:"column:billing_status"`
	FulfillmentStatus  string     `gorm:"column:fulfillment_status"`
	ManualReviewReq    bool       `gorm:"column:manual_review_required"`
	PaymentDueAt       *time.Time `gorm:"column:payment_due_at"`
	CoverageStartAt    *time.Time `gorm:"column:coverage_start_at"`
	CoverageEndAt      *time.Time `gorm:"column:coverage_end_at"`
	ActorUserId        string     `gorm:"column:actor_user_id"`
	Portal             string     `gorm:"column:portal"`
	PurchaseOrderId    string     `gorm:"column:purchase_order_id"`
}

// scanRowToProto converts an orderScanRow to a proto Order by:
//   - Converting the status string to proto enum
//   - Constructing Money from amount + currency
//   - Converting time.Time to timestamppb.Timestamp
func scanRowToProto(row *orderScanRow) *ordersv1.Order {
	if row == nil {
		return nil
	}
	order := &ordersv1.Order{
		// Original fields
		OrderId:            row.OrderId,
		OrderNumber:        row.OrderNumber,
		TenantId:           row.TenantId,
		QuotationId:        row.QuotationId,
		CustomerId:         row.CustomerId,
		ProductId:          row.ProductId,
		PlanId:             row.PlanId,
		Status:             stringToOrderStatus(row.Status),
		Currency:           row.Currency,
		PaymentId:          row.PaymentId,
		PaymentGatewayRef:  row.PaymentGatewayRef,
		PolicyId:           row.PolicyId,
		CancellationReason: row.CancellationReason,
		FailureReason:      row.FailureReason,
		CreatedAt:          timestamppb.New(row.CreatedAt),
		UpdatedAt:          timestamppb.New(row.UpdatedAt),
		// Phase-2 extended fields
		InvoiceId:            row.InvoiceId,
		OrganisationId:       row.OrganisationId,
		IdempotencyKey:       row.IdempotencyKey,
		CorrelationId:        row.CorrelationId,
		PaymentStatus:        stringToOrderPaymentStatus(row.PaymentStatus),
		BillingStatus:        stringToOrderBillingStatus(row.BillingStatus),
		FulfillmentStatus:    stringToOrderFulfillmentStatus(row.FulfillmentStatus),
		ManualReviewRequired: row.ManualReviewReq,
		ActorUserId:          row.ActorUserId,
		Portal:               row.Portal,
		PurchaseOrderId:      row.PurchaseOrderId,
	}

	// Money proto from amount + currency
	order.TotalPayable = &commonv1.Money{
		Amount:        row.TotalPayable,
		Currency:      row.Currency,
		DecimalAmount: float64(row.TotalPayable) / 100.0,
	}

	// Nullable timestamps
	if row.PaidAt != nil {
		order.PaidAt = timestamppb.New(*row.PaidAt)
	}
	if row.PaymentDueAt != nil {
		order.PaymentDueAt = timestamppb.New(*row.PaymentDueAt)
	}
	if row.CoverageStartAt != nil {
		order.CoverageStartAt = timestamppb.New(*row.CoverageStartAt)
	}
	if row.CoverageEndAt != nil {
		order.CoverageEndAt = timestamppb.New(*row.CoverageEndAt)
	}

	return order
}

// ─── enum converters ──────────────────────────────────────────────────────────

func stringToOrderStatus(s string) ordersv1.OrderStatus {
	if v, ok := ordersv1.OrderStatus_value[s]; ok {
		return ordersv1.OrderStatus(v)
	}
	return ordersv1.OrderStatus_ORDER_STATUS_UNSPECIFIED
}

func stringToOrderPaymentStatus(s string) ordersv1.OrderPaymentStatus {
	if v, ok := ordersv1.OrderPaymentStatus_value[s]; ok {
		return ordersv1.OrderPaymentStatus(v)
	}
	return ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_UNSPECIFIED
}

func stringToOrderBillingStatus(s string) ordersv1.OrderBillingStatus {
	if v, ok := ordersv1.OrderBillingStatus_value[s]; ok {
		return ordersv1.OrderBillingStatus(v)
	}
	return ordersv1.OrderBillingStatus_ORDER_BILLING_STATUS_UNSPECIFIED
}

func stringToOrderFulfillmentStatus(s string) ordersv1.OrderFulfillmentStatus {
	if v, ok := ordersv1.OrderFulfillmentStatus_value[s]; ok {
		return ordersv1.OrderFulfillmentStatus(v)
	}
	return ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_UNSPECIFIED
}

// (stringToOrderStatus is defined above in the enum converters block)

func newOrderNumber() string {
	short := strings.ToUpper(uuid.NewString()[:6])
	return fmt.Sprintf("ORD-%s-%s", time.Now().UTC().Format("20060102"), short)
}

func statusStr(s ordersv1.OrderStatus) string { return s.String() }

// ─── CreateOrder ─────────────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) CreateOrder(ctx context.Context, input domain.OrderCreateInput) (*ordersv1.Order, error) {
	id := input.OrderID
	if id == "" {
		id = uuid.NewString()
	}
	number := input.OrderNumber
	if number == "" {
		number = newOrderNumber()
	}
	currency := input.Currency
	if currency == "" {
		currency = "BDT"
	}
	var amount int64
	if input.TotalPayable != nil {
		amount = input.TotalPayable.Amount
	}

	now := time.Now().UTC()
	row := map[string]any{
		// Original fields
		"order_id":      id,
		"order_number":  number,
		"tenant_id":     input.TenantID,
		"quotation_id":  input.QuotationID,
		"customer_id":   input.CustomerID,
		"product_id":    coalesceUUID(input.ProductID),
		"plan_id":       coalesceUUID(input.PlanID),
		"status":        statusStr(ordersv1.OrderStatus_ORDER_STATUS_PENDING),
		"total_payable": amount,
		"currency":      currency,
		"created_at":    now,
		"updated_at":    now,
		// Phase-2 extended fields
		"payment_status":          input.PaymentStatus.String(),
		"billing_status":          input.BillingStatus.String(),
		"fulfillment_status":      input.FulfillmentStatus.String(),
		"manual_review_required":  input.ManualReviewRequired,
		"actor_user_id":           nullableUUID(input.ActorUserID),
		"portal":                  input.Portal,
		"correlation_id":          input.CorrelationID,
	}
	// Only set nullable UUID columns when non-empty to avoid invalid UUID errors
	if input.IdempotencyKey != "" {
		row["idempotency_key"] = input.IdempotencyKey
	}
	if input.OrganisationID != "" {
		row["organisation_id"] = input.OrganisationID
	}
	if input.PurchaseOrderID != "" {
		row["purchase_order_id"] = input.PurchaseOrderID
	}
	if input.PaymentDueAt != nil {
		row["payment_due_at"] = *input.PaymentDueAt
	}
	if input.CoverageStartAt != nil {
		row["coverage_start_at"] = *input.CoverageStartAt
	}
	if input.CoverageEndAt != nil {
		row["coverage_end_at"] = *input.CoverageEndAt
	}

	if err := r.db.WithContext(ctx).Table(ordersTable).Create(row).Error; err != nil {
		return nil, fmt.Errorf("repository.CreateOrder: %w", err)
	}
	return r.GetOrder(ctx, id)
}

// ─── GetOrder ────────────────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) GetOrder(ctx context.Context, orderID string) (*ordersv1.Order, error) {
	var row orderScanRow
	query := `
		SELECT
			order_id, order_number, tenant_id, quotation_id, customer_id,
			product_id, plan_id, status, total_payable, currency,
			COALESCE(payment_id::text,'')          AS payment_id,
			COALESCE(payment_gateway_ref,'')        AS payment_gateway_ref,
			COALESCE(policy_id::text,'')            AS policy_id,
			COALESCE(cancellation_reason,'')        AS cancellation_reason,
			COALESCE(failure_reason,'')             AS failure_reason,
			created_at, updated_at, paid_at,
			COALESCE(invoice_id::text,'')           AS invoice_id,
			COALESCE(organisation_id::text,'')      AS organisation_id,
			COALESCE(idempotency_key,'')            AS idempotency_key,
			COALESCE(correlation_id,'')             AS correlation_id,
			COALESCE(payment_status,'')             AS payment_status,
			COALESCE(billing_status,'')             AS billing_status,
			COALESCE(fulfillment_status,'')         AS fulfillment_status,
			manual_review_required,
			payment_due_at, coverage_start_at, coverage_end_at,
			COALESCE(actor_user_id::text,'')        AS actor_user_id,
			COALESCE(portal,'')                     AS portal,
			COALESCE(purchase_order_id::text,'')    AS purchase_order_id
		FROM insurance_schema.orders
		WHERE order_id = ?
		LIMIT 1
	`
	err := r.db.WithContext(ctx).Raw(query, orderID).Scan(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("repository.GetOrder: %w", err)
	}
	if row.OrderId == "" {
		return nil, domain.ErrNotFound
	}
	return scanRowToProto(&row), nil
}

// ─── ListOrders ──────────────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) ListOrders(
	ctx context.Context,
	pageSize, offset int,
	customerID string,
	status ordersv1.OrderStatus,
) ([]*ordersv1.Order, int64, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Build WHERE clause
	var whereConditions []string
	var args []any
	if customerID != "" {
		whereConditions = append(whereConditions, "customer_id = ?")
		args = append(args, customerID)
	}
	if status != ordersv1.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		whereConditions = append(whereConditions, "status = ?")
		args = append(args, statusStr(status))
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM insurance_schema.orders " + whereClause
	var total int64
	if err := r.db.WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repository.ListOrders count: %w", err)
	}

	// Fetch paginated results
	query := `
		SELECT
			order_id, order_number, tenant_id, quotation_id, customer_id,
			product_id, plan_id, status, total_payable, currency,
			COALESCE(payment_id::text,'')          AS payment_id,
			COALESCE(payment_gateway_ref,'')        AS payment_gateway_ref,
			COALESCE(policy_id::text,'')            AS policy_id,
			COALESCE(cancellation_reason,'')        AS cancellation_reason,
			COALESCE(failure_reason,'')             AS failure_reason,
			created_at, updated_at, paid_at,
			COALESCE(invoice_id::text,'')           AS invoice_id,
			COALESCE(organisation_id::text,'')      AS organisation_id,
			COALESCE(idempotency_key,'')            AS idempotency_key,
			COALESCE(correlation_id,'')             AS correlation_id,
			COALESCE(payment_status,'')             AS payment_status,
			COALESCE(billing_status,'')             AS billing_status,
			COALESCE(fulfillment_status,'')         AS fulfillment_status,
			manual_review_required,
			payment_due_at, coverage_start_at, coverage_end_at,
			COALESCE(actor_user_id::text,'')        AS actor_user_id,
			COALESCE(portal,'')                     AS portal,
			COALESCE(purchase_order_id::text,'')    AS purchase_order_id
		FROM insurance_schema.orders
		` + whereClause + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, pageSize, offset)

	var rows []orderScanRow
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("repository.ListOrders: %w", err)
	}

	orders := make([]*ordersv1.Order, len(rows))
	for i := range rows {
		orders[i] = scanRowToProto(&rows[i])
	}
	return orders, total, nil
}

// ─── UpdateOrderStatus ───────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, orderID string, status ordersv1.OrderStatus) error {
	updates := map[string]any{
		"status":     statusStr(status),
		"updated_at": time.Now().UTC(),
	}
	if status == ordersv1.OrderStatus_ORDER_STATUS_PAID {
		updates["paid_at"] = time.Now().UTC()
	}
	res := r.db.WithContext(ctx).Table(ordersTable).
		Where("order_id = ?", orderID).
		Updates(updates)
	if res.Error != nil {
		return fmt.Errorf("repository.UpdateOrderStatus: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── SetPaymentInfo ──────────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) SetPaymentInfo(ctx context.Context, orderID, paymentID, gatewayRef string) error {
	res := r.db.WithContext(ctx).Table(ordersTable).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"payment_id":          paymentID,
			"payment_gateway_ref": gatewayRef,
			"status":              statusStr(ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED),
			"updated_at":          time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("repository.SetPaymentInfo: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── SetPolicyID ─────────────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) SetPolicyID(ctx context.Context, orderID, policyID string) error {
	res := r.db.WithContext(ctx).Table(ordersTable).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"policy_id":  policyID,
			"status":     statusStr(ordersv1.OrderStatus_ORDER_STATUS_POLICY_ISSUED),
			"updated_at": time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("repository.SetPolicyID: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── SetCancellationReason ───────────────────────────────────────────────────

func (r *OrderRepositoryImpl) SetCancellationReason(ctx context.Context, orderID, reason string) error {
	res := r.db.WithContext(ctx).Table(ordersTable).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"cancellation_reason": reason,
			"status":              statusStr(ordersv1.OrderStatus_ORDER_STATUS_CANCELLED),
			"updated_at":          time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("repository.SetCancellationReason: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── SetFailureReason ────────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) SetFailureReason(ctx context.Context, orderID, reason string) error {
	res := r.db.WithContext(ctx).Table(ordersTable).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"failure_reason": reason,
			"status":         statusStr(ordersv1.OrderStatus_ORDER_STATUS_FAILED),
			"updated_at":     time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("repository.SetFailureReason: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── GetOrderByIdempotencyKey ────────────────────────────────────────────────

func (r *OrderRepositoryImpl) GetOrderByIdempotencyKey(ctx context.Context, key string) (*ordersv1.Order, error) {
	var row orderScanRow
	query := `
		SELECT
			order_id, order_number, tenant_id, quotation_id, customer_id,
			product_id, plan_id, status, total_payable, currency,
			COALESCE(payment_id::text,'')          AS payment_id,
			COALESCE(payment_gateway_ref,'')        AS payment_gateway_ref,
			COALESCE(policy_id::text,'')            AS policy_id,
			COALESCE(cancellation_reason,'')        AS cancellation_reason,
			COALESCE(failure_reason,'')             AS failure_reason,
			created_at, updated_at, paid_at,
			COALESCE(invoice_id::text,'')           AS invoice_id,
			COALESCE(organisation_id::text,'')      AS organisation_id,
			COALESCE(idempotency_key,'')            AS idempotency_key,
			COALESCE(correlation_id,'')             AS correlation_id,
			COALESCE(payment_status,'')             AS payment_status,
			COALESCE(billing_status,'')             AS billing_status,
			COALESCE(fulfillment_status,'')         AS fulfillment_status,
			manual_review_required,
			payment_due_at, coverage_start_at, coverage_end_at,
			COALESCE(actor_user_id::text,'')        AS actor_user_id,
			COALESCE(portal,'')                     AS portal,
			COALESCE(purchase_order_id::text,'')    AS purchase_order_id
		FROM insurance_schema.orders
		WHERE idempotency_key = ?
		LIMIT 1
	`
	err := r.db.WithContext(ctx).Raw(query, key).Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("repository.GetOrderByIdempotencyKey: %w", err)
	}
	if row.OrderId == "" {
		return nil, domain.ErrNotFound
	}
	return scanRowToProto(&row), nil
}

// ─── SetInvoiceID ────────────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) SetInvoiceID(ctx context.Context, orderID, invoiceID string) error {
	res := r.db.WithContext(ctx).Table(ordersTable).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"invoice_id":     invoiceID,
			"billing_status": "ORDER_BILLING_STATUS_INVOICED",
			"updated_at":     time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("repository.SetInvoiceID: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── SetFulfillmentStatus ─────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) SetFulfillmentStatus(ctx context.Context, orderID string, status ordersv1.OrderFulfillmentStatus) error {
	res := r.db.WithContext(ctx).Table(ordersTable).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"fulfillment_status": status.String(),
			"updated_at":         time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("repository.SetFulfillmentStatus: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── SetPaymentStatus ─────────────────────────────────────────────────────────

func (r *OrderRepositoryImpl) SetPaymentStatus(ctx context.Context, orderID string, status ordersv1.OrderPaymentStatus) error {
	res := r.db.WithContext(ctx).Table(ordersTable).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"payment_status": status.String(),
			"updated_at":     time.Now().UTC(),
		})
	if res.Error != nil {
		return fmt.Errorf("repository.SetPaymentStatus: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── misc ─────────────────────────────────────────────────────────────────────

// coalesceUUID returns the nil UUID string if the input is empty.
// Used only for NOT NULL UUID columns (product_id, plan_id) that must have a value.
func coalesceUUID(s string) string {
	if strings.TrimSpace(s) == "" {
		return "00000000-0000-0000-0000-000000000000"
	}
	return s
}

// nullableUUID returns nil if empty (for nullable UUID columns), or the value.
func nullableUUID(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
