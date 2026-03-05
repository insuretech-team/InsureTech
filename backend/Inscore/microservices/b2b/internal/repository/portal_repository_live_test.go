package repository

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	b2bLiveDBOnce sync.Once
	b2bLiveDB     *gorm.DB
	b2bLiveDBErr  error
)

func testB2BLiveDB(t *testing.T) *gorm.DB {
	t.Helper()

	if os.Getenv("INSURETECH_LIVE_DB_TESTS") != "1" {
		t.Skip("skipping live DB test: INSURETECH_LIVE_DB_TESTS != 1")
	}

	b2bLiveDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()

		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			resolvedPath, err := config.ResolveConfigPath("database.yaml")
			if err != nil {
				b2bLiveDBErr = err
				return
			}
			configPath = resolvedPath
		}

		b2bLiveDBErr = db.InitializeManagerForService(configPath)
		if b2bLiveDBErr != nil {
			return
		}
		b2bLiveDB = db.GetDB()
	})

	if b2bLiveDBErr != nil {
		t.Fatalf("live DB init failed: %v", b2bLiveDBErr)
	}
	if b2bLiveDB == nil {
		t.Fatal("live DB init failed: db is nil")
	}
	return b2bLiveDB
}

func TestPortalRepository_LiveDB_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	repo := NewPortalRepository(testB2BLiveDB(t))

	departments, deptTotal, err := repo.ListDepartments(ctx, 5, 0, "")
	require.NoError(t, err)
	require.GreaterOrEqual(t, deptTotal, int64(len(departments)))
	if len(departments) > 0 {
		department, err := repo.GetDepartment(ctx, departments[0].DepartmentId)
		require.NoError(t, err)
		require.NotNil(t, department)

		names, err := repo.GetDepartmentNames(ctx, []string{departments[0].DepartmentId})
		require.NoError(t, err)
		require.Equal(t, departments[0].Name, names[departments[0].DepartmentId])
	}

	employees, employeeTotal, err := repo.ListEmployees(ctx, 5, 0, "", "")
	require.NoError(t, err)
	require.GreaterOrEqual(t, employeeTotal, int64(len(employees)))
	if len(employees) > 0 {
		employee, err := repo.GetEmployee(ctx, employees[0].EmployeeUuid)
		require.NoError(t, err)
		require.NotNil(t, employee)
	}

	catalog, err := repo.ListCatalogPlans(ctx)
	require.NoError(t, err)
	if len(catalog) > 0 {
		plans, err := repo.GetCatalogPlansByPlanIDs(ctx, []string{catalog[0].PlanID})
		require.NoError(t, err)
		require.Contains(t, plans, catalog[0].PlanID)
	}

	purchaseOrders, poTotal, err := repo.ListPurchaseOrders(ctx, 5, 0, "", 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, poTotal, int64(len(purchaseOrders)))
	if len(purchaseOrders) > 0 {
		purchaseOrder, err := repo.GetPurchaseOrder(ctx, purchaseOrders[0].PurchaseOrderId)
		require.NoError(t, err)
		require.NotNil(t, purchaseOrder)
	}
}
