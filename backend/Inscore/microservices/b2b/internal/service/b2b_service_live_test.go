package service

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	b2bServiceLiveDBOnce sync.Once
	b2bServiceLiveDB     *gorm.DB
	b2bServiceLiveDBErr  error
)

func testB2BServiceLiveDB(t *testing.T) *gorm.DB {
	t.Helper()

	if os.Getenv("INSURETECH_LIVE_DB_TESTS") != "1" {
		t.Skip("skipping live DB test: INSURETECH_LIVE_DB_TESTS != 1")
	}

	b2bServiceLiveDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()

		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			resolvedPath, err := config.ResolveConfigPath("database.yaml")
			if err != nil {
				b2bServiceLiveDBErr = err
				return
			}
			configPath = resolvedPath
		}

		b2bServiceLiveDBErr = db.InitializeManagerForService(configPath)
		if b2bServiceLiveDBErr != nil {
			return
		}
		b2bServiceLiveDB = db.GetDB()
	})

	if b2bServiceLiveDBErr != nil {
		t.Fatalf("live DB init failed: %v", b2bServiceLiveDBErr)
	}
	if b2bServiceLiveDB == nil {
		t.Fatal("live DB init failed: db is nil")
	}
	return b2bServiceLiveDB
}

func TestB2BService_Live_PurchaseOrdersResolveProductAndPlanNames(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	svc := NewB2BService(repository.NewPortalRepository(testB2BServiceLiveDB(t)))
	dbConn := testB2BServiceLiveDB(t)

	resp, err := svc.ListPurchaseOrders(context.Background(), &b2bservicev1.ListPurchaseOrdersRequest{
		PageSize: 25,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.PurchaseOrders)

	catalog, err := repository.NewPortalRepository(testB2BServiceLiveDB(t)).ListCatalogPlans(context.Background())
	require.NoError(t, err)
	for i, item := range catalog {
		if i >= 10 {
			break
		}
		t.Logf("catalog[%d]: product_id=%s product_name=%s plan_id=%s plan_name=%s", i, item.ProductID, item.ProductName, item.PlanID, item.PlanName)
	}

	type rawPlan struct {
		ProductID   string `gorm:"column:product_id"`
		ProductName string `gorm:"column:product_name"`
		PlanID      string `gorm:"column:plan_id"`
		PlanName    string `gorm:"column:plan_name"`
	}
	var rawPlans []rawPlan
	err = dbConn.Table("insurance_schema.product_plans AS pp").
		Select("pp.product_id, p.product_name, pp.plan_id, pp.plan_name").
		Joins("LEFT JOIN insurance_schema.products AS p ON p.product_id = pp.product_id").
		Where("pp.plan_id IN ?", []string{
			"55555555-5555-5555-5555-555555555001",
			"55555555-5555-5555-5555-555555555002",
			"55555555-5555-5555-5555-555555555004",
		}).
		Find(&rawPlans).Error
	require.NoError(t, err)
	for _, item := range rawPlans {
		t.Logf("raw plan: product_id=%s product_name=%q plan_id=%s plan_name=%q", item.ProductID, item.ProductName, item.PlanID, item.PlanName)
	}

	for _, item := range resp.PurchaseOrders {
		require.NotNil(t, item)
		require.NotNil(t, item.PurchaseOrder)
		require.NotEmpty(t, item.PurchaseOrder.PurchaseOrderId)

		poNumber := item.PurchaseOrder.GetPurchaseOrderNumber()
		if strings.HasPrefix(poNumber, "PO-20260227-") {
			t.Logf("seeded purchase order %s -> product_id=%s plan_id=%s resolved_product=%q resolved_plan=%q",
				poNumber,
				item.PurchaseOrder.GetProductId(),
				item.PurchaseOrder.GetPlanId(),
				item.GetProductName(),
				item.GetPlanName(),
			)
			require.NotEqual(t, "Unknown Product", item.GetProductName(), "seeded purchase order should resolve product name")
			require.NotEqual(t, "Unknown Plan", item.GetPlanName(), "seeded purchase order should resolve plan name")
			require.NotEmpty(t, strings.TrimSpace(item.GetProductName()))
			require.NotEmpty(t, strings.TrimSpace(item.GetPlanName()))
		}
	}
}
