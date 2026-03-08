package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── CreateOrder ─────────────────────────────────────────────────────────────

func TestLiveRepository_CreateOrder(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := newTestOrderInput(env.fixture)
	order, err := env.repo.CreateOrder(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, order)

	t.Cleanup(func() {
		env.db.Exec("DELETE FROM insurance_schema.orders WHERE order_id = ?", order.OrderId)
	})

	assert.NotEmpty(t, order.OrderId)
	assert.NotEmpty(t, order.OrderNumber)
	assert.Equal(t, input.TenantID, order.TenantId)
	assert.Equal(t, input.QuotationID, order.QuotationId)
	assert.Equal(t, input.CustomerID, order.CustomerId)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PENDING, order.Status)

	// Money reconstruction
	require.NotNil(t, order.TotalPayable)
	assert.Equal(t, int64(50000), order.TotalPayable.Amount)
	assert.Equal(t, "BDT", order.TotalPayable.Currency)
	assert.InDelta(t, 500.0, order.TotalPayable.DecimalAmount, 0.001)
	assert.Equal(t, "BDT", order.Currency)

	// Timestamps populated by DB
	assert.NotNil(t, order.CreatedAt)
	assert.NotNil(t, order.UpdatedAt)
}

func TestLiveRepository_CreateOrder_DuplicateOrderNumber(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := newTestOrderInput(env.fixture)
	input.OrderNumber = "ORD-DUPLICATE-" + uuid.NewString()[:6]

	order1, err := env.repo.CreateOrder(ctx, input)
	require.NoError(t, err)
	t.Cleanup(func() {
		env.db.Exec("DELETE FROM insurance_schema.orders WHERE order_id = ?", order1.OrderId)
	})

	// Second create with same order_number must fail (UNIQUE constraint)
	input2 := newTestOrderInput(env.fixture)
	input2.OrderNumber = input.OrderNumber
	_, err = env.repo.CreateOrder(ctx, input2)
	require.Error(t, err, "duplicate order_number should return an error")
}

// ─── GetOrder ────────────────────────────────────────────────────────────────

func TestLiveRepository_GetOrder(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	created := createTestOrder(t, env)

	got, err := env.repo.GetOrder(ctx, created.OrderId)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, created.OrderId, got.OrderId)
	assert.Equal(t, created.OrderNumber, got.OrderNumber)
	assert.Equal(t, created.TenantId, got.TenantId)
	assert.Equal(t, created.QuotationId, got.QuotationId)
	assert.Equal(t, created.CustomerId, got.CustomerId)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PENDING, got.Status)
	require.NotNil(t, got.TotalPayable)
	assert.Equal(t, int64(50000), got.TotalPayable.Amount)
	assert.Equal(t, "BDT", got.TotalPayable.Currency)
}

func TestLiveRepository_GetOrder_NotFound(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := env.repo.GetOrder(ctx, uuid.NewString())
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// ─── ListOrders ──────────────────────────────────────────────────────────────

func TestLiveRepository_ListOrders(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create 3 orders for the same customer (same fixture)
	var created []*ordersv1.Order
	for i := 0; i < 3; i++ {
		input := newTestOrderInput(env.fixture)
		input.TotalPayable = &commonv1.Money{Amount: int64((i + 1) * 10000)}
		o, err := env.repo.CreateOrder(ctx, input)
		require.NoError(t, err)
		created = append(created, o)
		capturedID := o.OrderId
		t.Cleanup(func() {
			env.db.Exec("DELETE FROM insurance_schema.orders WHERE order_id = ?", capturedID)
		})
	}

	// List by customerID
	orders, total, err := env.repo.ListOrders(ctx, 10, 0, env.fixture.UserID, ordersv1.OrderStatus_ORDER_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(orders), 3)

	// Pagination page 1
	page1, total1, err := env.repo.ListOrders(ctx, 2, 0, env.fixture.UserID, ordersv1.OrderStatus_ORDER_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total1, int64(3))
	assert.Len(t, page1, 2)

	// Pagination page 2
	page2, _, err := env.repo.ListOrders(ctx, 2, 2, env.fixture.UserID, ordersv1.OrderStatus_ORDER_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(page2), 1)

	// Filter by status: PENDING → all 3 orders
	pending, _, err := env.repo.ListOrders(ctx, 10, 0, env.fixture.UserID, ordersv1.OrderStatus_ORDER_STATUS_PENDING)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(pending), 3)

	// Filter by status: PAID → none of our new orders
	paid, paidTotal, err := env.repo.ListOrders(ctx, 10, 0, env.fixture.UserID, ordersv1.OrderStatus_ORDER_STATUS_PAID)
	require.NoError(t, err)
	_ = paid
	_ = paidTotal // may be non-zero if prior test left paid orders; just ensure no error

	_ = created
}

// ─── UpdateOrderStatus ───────────────────────────────────────────────────────

func TestLiveRepository_UpdateOrderStatus(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	order := createTestOrder(t, env)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PENDING, order.Status)

	// → PAYMENT_INITIATED
	require.NoError(t, env.repo.UpdateOrderStatus(ctx, order.OrderId, ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED))
	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED, got.Status)
	assert.Nil(t, got.PaidAt)

	// → PAID (should also set paid_at)
	require.NoError(t, env.repo.UpdateOrderStatus(ctx, order.OrderId, ordersv1.OrderStatus_ORDER_STATUS_PAID))
	paid, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PAID, paid.Status)
	assert.NotNil(t, paid.PaidAt, "paid_at must be set when status transitions to PAID")
}

func TestLiveRepository_UpdateOrderStatus_NotFound(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := env.repo.UpdateOrderStatus(ctx, uuid.NewString(), ordersv1.OrderStatus_ORDER_STATUS_PAID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// ─── SetPaymentInfo ──────────────────────────────────────────────────────────

func TestLiveRepository_SetPaymentInfo(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order := createTestOrder(t, env)

	paymentID := env.fixture.PaymentID
	gatewayRef := "GW-TEST-" + uuid.NewString()[:8]

	require.NoError(t, env.repo.SetPaymentInfo(ctx, order.OrderId, paymentID, gatewayRef))

	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED, got.Status)
	assert.Equal(t, paymentID, got.PaymentId)
	assert.Equal(t, gatewayRef, got.PaymentGatewayRef)
}

// ─── SetPolicyID ─────────────────────────────────────────────────────────────

func TestLiveRepository_SetPolicyID(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order := createTestOrder(t, env)
	policyID := env.fixture.PolicyID

	require.NoError(t, env.repo.SetPolicyID(ctx, order.OrderId, policyID))

	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_POLICY_ISSUED, got.Status)
	assert.Equal(t, policyID, got.PolicyId)
}

// ─── SetCancellationReason ───────────────────────────────────────────────────

func TestLiveRepository_SetCancellationReason(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order := createTestOrder(t, env)
	reason := "Customer requested cancellation"

	require.NoError(t, env.repo.SetCancellationReason(ctx, order.OrderId, reason))

	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_CANCELLED, got.Status)
	assert.Equal(t, reason, got.CancellationReason)
}

// ─── SetFailureReason ────────────────────────────────────────────────────────

func TestLiveRepository_SetFailureReason(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order := createTestOrder(t, env)
	reason := "Insufficient funds"

	require.NoError(t, env.repo.SetFailureReason(ctx, order.OrderId, reason))

	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_FAILED, got.Status)
	assert.Equal(t, reason, got.FailureReason)
}

// ─── Phase-2: GetOrderByIdempotencyKey ───────────────────────────────────────

func TestLiveRepository_GetOrderByIdempotencyKey(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := "idem-test-" + uuid.NewString()
	input := newTestOrderInput(env.fixture)
	input.IdempotencyKey = key

	order, err := env.repo.CreateOrder(ctx, input)
	require.NoError(t, err)
	t.Cleanup(func() { env.db.Exec("DELETE FROM insurance_schema.orders WHERE order_id = ?", order.OrderId) })

	// Lookup by idempotency key
	found, err := env.repo.GetOrderByIdempotencyKey(ctx, key)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, order.OrderId, found.OrderId)
	assert.Equal(t, key, found.IdempotencyKey)

	// Missing key → ErrNotFound
	_, err = env.repo.GetOrderByIdempotencyKey(ctx, "non-existent-key-"+uuid.NewString())
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestLiveRepository_IdempotencyKey_UniqueConstraint(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := "idem-dup-" + uuid.NewString()
	input := newTestOrderInput(env.fixture)
	input.IdempotencyKey = key

	o1, err := env.repo.CreateOrder(ctx, input)
	require.NoError(t, err)
	t.Cleanup(func() { env.db.Exec("DELETE FROM insurance_schema.orders WHERE order_id = ?", o1.OrderId) })

	// Second create with same key must fail (unique constraint)
	input2 := newTestOrderInput(env.fixture)
	input2.IdempotencyKey = key
	_, err = env.repo.CreateOrder(ctx, input2)
	require.Error(t, err, "duplicate idempotency_key must return an error")
}

// ─── Phase-2: SetPaymentStatus ────────────────────────────────────────────────

func TestLiveRepository_SetPaymentStatus(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order := createTestOrder(t, env)
	assert.Equal(t, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_UNPAID, order.PaymentStatus)

	// Transition to IN_PROGRESS
	require.NoError(t, env.repo.SetPaymentStatus(ctx, order.OrderId, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAYMENT_IN_PROGRESS))
	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAYMENT_IN_PROGRESS, got.PaymentStatus)

	// Transition to PAID
	require.NoError(t, env.repo.SetPaymentStatus(ctx, order.OrderId, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID))
	got, err = env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID, got.PaymentStatus)

	// NotFound
	err = env.repo.SetPaymentStatus(ctx, uuid.NewString(), ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

// ─── Phase-2: SetFulfillmentStatus ────────────────────────────────────────────

func TestLiveRepository_SetFulfillmentStatus(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order := createTestOrder(t, env)
	assert.Equal(t, ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_NOT_STARTED, order.FulfillmentStatus)

	// → IN_PROGRESS
	require.NoError(t, env.repo.SetFulfillmentStatus(ctx, order.OrderId, ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLMENT_IN_PROGRESS))
	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLMENT_IN_PROGRESS, got.FulfillmentStatus)

	// → FULFILLED
	require.NoError(t, env.repo.SetFulfillmentStatus(ctx, order.OrderId, ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLED))
	got, err = env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLED, got.FulfillmentStatus)

	// NotFound
	err = env.repo.SetFulfillmentStatus(ctx, uuid.NewString(), ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLED)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

// ─── Phase-2: SetInvoiceID ────────────────────────────────────────────────────

func TestLiveRepository_SetInvoiceID(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order := createTestOrder(t, env)
	assert.Empty(t, order.InvoiceId)
	assert.Equal(t, ordersv1.OrderBillingStatus_ORDER_BILLING_STATUS_NOT_INVOICED, order.BillingStatus)

	invoiceID := uuid.NewString()
	require.NoError(t, env.repo.SetInvoiceID(ctx, order.OrderId, invoiceID))

	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)
	assert.Equal(t, invoiceID, got.InvoiceId)
	assert.Equal(t, ordersv1.OrderBillingStatus_ORDER_BILLING_STATUS_INVOICED, got.BillingStatus)

	// NotFound
	err = env.repo.SetInvoiceID(ctx, uuid.NewString(), invoiceID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

// ─── Phase-2: Phase2 fields round-trip ───────────────────────────────────────

func TestLiveRepository_Phase2Fields_RoundTrip(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now().UTC().Truncate(time.Second)
	coverageEnd := now.Add(365 * 24 * time.Hour)

	input := newTestOrderInput(env.fixture)
	input.IdempotencyKey  = "p2-roundtrip-" + uuid.NewString()
	input.CorrelationID   = uuid.NewString()
	input.ActorUserID     = env.fixture.UserID
	input.Portal          = "b2c"
	input.PaymentDueAt    = &now
	input.CoverageStartAt = &now
	input.CoverageEndAt   = &coverageEnd

	order, err := env.repo.CreateOrder(ctx, input)
	require.NoError(t, err)
	t.Cleanup(func() { env.db.Exec("DELETE FROM insurance_schema.orders WHERE order_id = ?", order.OrderId) })

	got, err := env.repo.GetOrder(ctx, order.OrderId)
	require.NoError(t, err)

	assert.Equal(t, input.IdempotencyKey, got.IdempotencyKey)
	assert.Equal(t, input.CorrelationID, got.CorrelationId)
	assert.Equal(t, input.ActorUserID, got.ActorUserId)
	assert.Equal(t, "b2c", got.Portal)
	assert.Equal(t, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_UNPAID, got.PaymentStatus)
	assert.Equal(t, ordersv1.OrderBillingStatus_ORDER_BILLING_STATUS_NOT_INVOICED, got.BillingStatus)
	assert.Equal(t, ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_NOT_STARTED, got.FulfillmentStatus)
	assert.False(t, got.ManualReviewRequired)
	require.NotNil(t, got.PaymentDueAt)
	require.NotNil(t, got.CoverageStartAt)
	require.NotNil(t, got.CoverageEndAt)
}

// ─── Full lifecycle ───────────────────────────────────────────────────────────

func TestLiveRepository_FullOrderLifecycle(t *testing.T) {
	env := setupLiveTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Create
	order := createTestOrder(t, env)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PENDING, order.Status)

	// 2. Initiate payment (use fixture paymentID to satisfy FK constraint)
	paymentID := env.fixture.PaymentID
	gatewayRef := "GW-LIFECYCLE-" + uuid.NewString()[:8]
	require.NoError(t, env.repo.SetPaymentInfo(ctx, order.OrderId, paymentID, gatewayRef))
	got, _ := env.repo.GetOrder(ctx, order.OrderId)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED, got.Status)
	assert.Equal(t, paymentID, got.PaymentId)

	// 3. Mark paid
	require.NoError(t, env.repo.UpdateOrderStatus(ctx, order.OrderId, ordersv1.OrderStatus_ORDER_STATUS_PAID))
	got, _ = env.repo.GetOrder(ctx, order.OrderId)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PAID, got.Status)
	assert.NotNil(t, got.PaidAt)

	// 4. Issue policy
	policyID := env.fixture.PolicyID
	require.NoError(t, env.repo.SetPolicyID(ctx, order.OrderId, policyID))
	got, _ = env.repo.GetOrder(ctx, order.OrderId)
	assert.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_POLICY_ISSUED, got.Status)
	assert.Equal(t, policyID, got.PolicyId)
}
