package repository_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"gorm.io/gorm"
)

// testPolicyNumber generates a valid policy number from a UUID.
// Format: LBT-YYYY-XXXX-NNNNNN where XXXX is 4 alphanumeric chars and NNNNNN is 6 digits.
func testPolicyNumber(prefix, id string) string {
	// Get digits from the UUID for the numeric part
	digits := ""
	for _, c := range id {
		if c >= '0' && c <= '9' {
			digits += string(c)
		}
	}
	// Pad with zeros to ensure at least 6 digits
	for len(digits) < 6 {
		digits = "0" + digits
	}
	numPart := digits[len(digits)-6:]
	return fmt.Sprintf("LBT-2026-%s-%s", prefix, numPart)
}

// testMobileNumber generates a valid +880 mobile number from a UUID.
// Format: +8801X-XXXXXXXX where X are digits (13 chars total).
func testMobileNumber(id string) string {
	digits := ""
	for _, c := range id {
		if c >= '0' && c <= '9' {
			digits += string(c)
			if len(digits) == 9 {
				break
			}
		}
	}
	// pad with zeros if UUID had fewer than 9 numeric chars
	for len(digits) < 9 {
		digits += "0"
	}
	return "+8801" + digits
}

// testProductCode generates a unique product code from a UUID.
// Format: ABC-123 where ABC are letters (a-f from UUID) and 123 are digits.
func testProductCode(id string) string {
	letters := ""
	for _, c := range id {
		if (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
			letters += strings.ToUpper(string(c))
			if len(letters) == 3 {
				break
			}
		}
	}
	for len(letters) < 3 {
		letters += "A"
	}

	digits := ""
	for _, c := range id {
		if c >= '0' && c <= '9' {
			digits += string(c)
			if len(digits) == 3 {
				break
			}
		}
	}
	for len(digits) < 3 {
		digits += "0"
	}

	return letters + "-" + digits
}

// ─── DB singleton ─────────────────────────────────────────────────────────────

var (
	ordersTestDBOnce sync.Once
	ordersTestDB     *gorm.DB
	ordersTestDBErr  error
)

func testOrdersDB(t *testing.T) *gorm.DB {
	t.Helper()
	ordersTestDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}
		ordersTestDBErr = db.InitializeManagerForService(configPath)
		if ordersTestDBErr != nil {
			return
		}
		ordersTestDB = db.GetDB()
	})
	if ordersTestDBErr != nil {
		t.Skipf("skipping live DB test: %v", ordersTestDBErr)
	}
	if ordersTestDB == nil {
		t.Skip("skipping live DB test: db is nil")
	}
	return ordersTestDB
}

// ─── Env ─────────────────────────────────────────────────────────────────────

type liveOrdersEnv struct {
	db      *gorm.DB
	repo    *repository.OrderRepositoryImpl
	fixture *testFixtures
}

func setupLiveTest(t *testing.T) *liveOrdersEnv {
	t.Helper()
	gormDB := testOrdersDB(t)
	fixtures := createFixtures(t, gormDB)
	return &liveOrdersEnv{
		db:      gormDB,
		repo:    repository.NewOrderRepository(gormDB),
		fixture: fixtures,
	}
}

// ─── Fixtures ─────────────────────────────────────────────────────────────────

// testFixtures holds IDs for records required by orders FK constraints.
type testFixtures struct {
	UserID      string
	ProductID   string
	PlanID      string
	QuotationID string
	TenantID    string
	PaymentID   string
	PolicyID    string
}

// createFixtures inserts minimal prerequisite records for orders tests and
// registers t.Cleanup to remove them in reverse FK order.
func createFixtures(t *testing.T, gormDB *gorm.DB) *testFixtures {
	t.Helper()
	ctx := context.Background()

	fx := &testFixtures{
		UserID:    uuid.NewString(),
		ProductID: uuid.NewString(),
		PlanID:    uuid.NewString(),
		TenantID:  uuid.NewString(),
	}

	// 1. Create user in authn_schema.users
	err := gormDB.WithContext(ctx).Exec(`
		INSERT INTO authn_schema.users
			(user_id, mobile_number, email, password_hash, status, user_type, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'ACTIVE', 'B2C_CUSTOMER', false, NOW(), NOW())`,
		fx.UserID,
		testMobileNumber(fx.UserID),
		fmt.Sprintf("orders-test-%s@test.local", fx.UserID[:8]),
		"$2a$10$testhashXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	).Error
	if err != nil {
		t.Skipf("cannot create test user fixture: %v", err)
	}
	t.Cleanup(func() {
		gormDB.Exec("DELETE FROM authn_schema.users WHERE user_id = ?", fx.UserID)
	})

	// 2. Create product
	err = gormDB.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.products
			(product_id, product_code, product_name, category, base_premium,
			 min_sum_insured, max_sum_insured, min_tenure_months, max_tenure_months,
			 status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, 'LIFE', 50000, 100000, 10000000, 12, 120, 'ACTIVE', $4, NOW(), NOW())`,
		fx.ProductID,
		testProductCode(fx.ProductID),
		fmt.Sprintf("Test Product %s", fx.ProductID[:8]),
		fx.UserID,
	).Error
	if err != nil {
		t.Skipf("cannot create test product fixture: %v", err)
	}
	t.Cleanup(func() {
		gormDB.Exec("DELETE FROM insurance_schema.products WHERE product_id = ?", fx.ProductID)
	})

	// 3. Create plan linked to product
	err = gormDB.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.product_plans
			(plan_id, product_id, plan_name, premium_amount, min_sum_insured, max_sum_insured, created_at, updated_at)
		VALUES ($1, $2, $3, 50000, 100000, 10000000, NOW(), NOW())`,
		fx.PlanID,
		fx.ProductID,
		fmt.Sprintf("Test Plan %s", fx.PlanID[:8]),
	).Error
	if err != nil {
		t.Skipf("cannot create test plan fixture: %v", err)
	}
	t.Cleanup(func() {
		gormDB.Exec("DELETE FROM insurance_schema.product_plans WHERE plan_id = ?", fx.PlanID)
	})

	// 4. Create payment fixture (used by SetPaymentInfo tests)
	fx.PaymentID = uuid.NewString()
	err = gormDB.WithContext(ctx).Exec(`
		INSERT INTO payment_schema.payments
			(payment_id, type, method, status, amount, currency, payer_id, created_at, updated_at)
		VALUES ($1, 'PREMIUM', 'bKash', 'INITIATED', 50000, 'BDT', $2, NOW(), NOW())`,
		fx.PaymentID, fx.UserID,
	).Error
	if err != nil {
		t.Skipf("cannot create payment fixture: %v", err)
	}
	t.Cleanup(func() {
		gormDB.Exec("DELETE FROM payment_schema.payments WHERE payment_id = ?", fx.PaymentID)
	})

	// 5. Create policy fixture (used by SetPolicyID tests)
	fx.PolicyID = uuid.NewString()
	err = gormDB.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.policies
			(policy_id, policy_number, product_id, customer_id, status, premium_amount, sum_insured, tenure_months,
			 start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'ACTIVE', 50000, 1000000, 12, CURRENT_DATE, CURRENT_DATE + INTERVAL '12 months', NOW(), NOW())`,
		fx.PolicyID,
		testPolicyNumber("TST1", fx.PolicyID),
		fx.ProductID,
		fx.UserID,
	).Error
	if err != nil {
		t.Skipf("cannot create policy fixture: %v", err)
	}
	t.Cleanup(func() {
		gormDB.Exec("DELETE FROM insurance_schema.policies WHERE policy_id = ?", fx.PolicyID)
	})

	// 6. Create quotation
	fx.QuotationID = uuid.NewString()
	err = gormDB.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.quotations
			(quotation_id, business_id, plan_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, 'APPROVED', NOW(), NOW())`,
		fx.QuotationID,
		fx.UserID,
		fx.PlanID,
	).Error
	if err != nil {
		t.Skipf("cannot create test quotation fixture: %v", err)
	}
	t.Cleanup(func() {
		gormDB.Exec("DELETE FROM insurance_schema.quotations WHERE quotation_id = ?", fx.QuotationID)
	})

	return fx
}

// ─── Input builder ───────────────────────────────────────────────────────────

// newTestOrderInput returns a valid OrderCreateInput using real fixture IDs.
// Includes all phase-2 fields with proper defaults to match DB constraints.
func newTestOrderInput(fx *testFixtures) domain.OrderCreateInput {
	return domain.OrderCreateInput{
		TenantID:    fx.TenantID,
		QuotationID: fx.QuotationID,
		CustomerID:  fx.UserID,
		ProductID:   fx.ProductID,
		PlanID:      fx.PlanID,
		Currency:    "BDT",
		TotalPayable: &commonv1.Money{
			Amount:   50000, // 500.00 BDT in paisa
			Currency: "BDT",
		},
		// Phase-2 defaults
		PaymentStatus:     ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_UNPAID,
		BillingStatus:     ordersv1.OrderBillingStatus_ORDER_BILLING_STATUS_NOT_INVOICED,
		FulfillmentStatus: ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_NOT_STARTED,
	}
}

// createTestOrder inserts an order and registers cleanup.
func createTestOrder(t *testing.T, env *liveOrdersEnv) *ordersv1.Order {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order, err := env.repo.CreateOrder(ctx, newTestOrderInput(env.fixture))
	if err != nil {
		t.Fatalf("createTestOrder: %v", err)
	}
	t.Cleanup(func() {
		env.db.Exec("DELETE FROM insurance_schema.orders WHERE order_id = ?", order.OrderId)
	})
	return order
}
