package service

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bsvcv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	svcTestBusinessID = "22222222-2222-2222-2222-222222222001"
	svcTestTenantID   = "11111111-1111-1111-1111-111111111111"
	svcPlanHealth1    = "55555555-5555-5555-5555-555555555001"
)

var (
	svcLiveOnce sync.Once
	svcLiveDB   *gorm.DB
	svcLiveErr  error
	svcLiveSvc  *B2BService
)

func testSvc(t *testing.T) (*B2BService, *gorm.DB) {
	t.Helper()
	if os.Getenv("INSURETECH_LIVE_DB_TESTS") != "1" {
		t.Skip("skipping live service test: set INSURETECH_LIVE_DB_TESTS=1 to run")
	}
	svcLiveOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		cfgPath, err := config.ResolveConfigPath("database.yaml")
		if err != nil { svcLiveErr = err; return }
		svcLiveErr = db.InitializeManagerForService(cfgPath)
		if svcLiveErr != nil { return }
		svcLiveDB = db.GetDB()
		repo := repository.NewPortalRepository(svcLiveDB)
		svcLiveSvc = NewB2BService(repo, nil) // nil publisher for tests
	})
	if svcLiveErr != nil { t.Fatalf("live service init: %v", svcLiveErr) }
	if svcLiveSvc == nil { t.Fatal("live service is nil") }
	return svcLiveSvc, svcLiveDB
}

func svcMoney(amount float64) *commonv1.Money {
	return &commonv1.Money{Amount: int64(amount * 100), Currency: "BDT", DecimalAmount: amount}
}

func cleanupSvcOrg(t *testing.T, gdb *gorm.DB, orgID string) {
	t.Helper()
	gdb.Exec("UPDATE b2b_schema.organisations SET deleted_at = NOW() WHERE organisation_id = $1", orgID)
}

func cleanupSvcPO(t *testing.T, gdb *gorm.DB, poID string) {
	t.Helper()
	gdb.Exec("UPDATE b2b_schema.purchase_orders SET deleted_at = NOW() WHERE purchase_order_id = $1", poID)
}

// ─── ORGANISATION ─────────────────────────────────────────────────────────────

func TestB2BService_Organisation_CRUD(t *testing.T) {
	ctx := context.Background()
	svc, gdb := testSvc(t)

	orgName := fmt.Sprintf("Svc Org %s", uuid.NewString()[:8])

	// CREATE
	createResp, err := svc.CreateOrganisation(ctx, &b2bsvcv1.CreateOrganisationRequest{
		TenantId:     svcTestTenantID,
		Name:         orgName,
		Code:         fmt.Sprintf("SVC%s", uuid.NewString()[:4]),
		Industry:     "Technology",
		ContactEmail: "svctest@example.com",
		ContactPhone: "+8801700000000",
		Address:      "Dhaka, Bangladesh",
	})
	require.NoError(t, err)
	org := createResp.GetOrganisation()
	require.NotNil(t, org)
	require.NotEmpty(t, org.OrganisationId)
	assert.Equal(t, orgName, org.Name)
	assert.Equal(t, svcTestTenantID, org.TenantId)
	t.Cleanup(func() { cleanupSvcOrg(t, gdb, org.OrganisationId) })
	t.Logf("created org: id=%s name=%s", org.OrganisationId, org.Name)

	// GET
	getResp, err := svc.GetOrganisation(ctx, &b2bsvcv1.GetOrganisationRequest{OrganisationId: org.OrganisationId})
	require.NoError(t, err)
	assert.Equal(t, orgName, getResp.GetOrganisation().Name)

	// LIST
	listResp, err := svc.ListOrganisations(ctx, &b2bsvcv1.ListOrganisationsRequest{
		TenantId: svcTestTenantID,
		PageSize: 50,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, listResp.GetTotalCount(), int32(1))
	found := false
	for _, o := range listResp.GetOrganisations() {
		if o.OrganisationId == org.OrganisationId { found = true }
	}
	assert.True(t, found)

	// UPDATE
	updateResp, err := svc.UpdateOrganisation(ctx, &b2bsvcv1.UpdateOrganisationRequest{
		OrganisationId: org.OrganisationId,
		Name:           orgName + " Updated",
		ContactEmail:   "updated@example.com",
		Industry:       "Finance",
	})
	require.NoError(t, err)
	assert.Equal(t, orgName+" Updated", updateResp.GetOrganisation().Name)
	assert.Equal(t, "Finance", updateResp.GetOrganisation().Industry)
}

// ─── ORG MEMBER ───────────────────────────────────────────────────────────────

func TestB2BService_OrgMember_AddResolveRemove(t *testing.T) {
	ctx := context.Background()
	svc, gdb := testSvc(t)

	createResp, err := svc.CreateOrganisation(ctx, &b2bsvcv1.CreateOrganisationRequest{
		TenantId: svcTestTenantID,
		Name:     fmt.Sprintf("Member Org %s", uuid.NewString()[:8]),
		Code:     fmt.Sprintf("MBR%s", uuid.NewString()[:4]),
	})
	require.NoError(t, err)
	org := createResp.GetOrganisation()
	t.Cleanup(func() { cleanupSvcOrg(t, gdb, org.OrganisationId) })

	userID := uuid.NewString()

	// ADD
	addResp, err := svc.AddOrgMember(ctx, &b2bsvcv1.AddOrgMemberRequest{
		OrganisationId: org.OrganisationId,
		UserId:         userID,
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_HR_MANAGER,
	})
	require.NoError(t, err)
	member := addResp.GetMember()
	require.NotNil(t, member)
	assert.Equal(t, userID, member.UserId)
	assert.Equal(t, b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_HR_MANAGER, member.Role)
	t.Logf("added member: id=%s user=%s role=%s", member.MemberId, member.UserId, member.Role)

	// RESOLVE
	resolveResp, err := svc.ResolveMyOrganisation(ctx, &b2bsvcv1.ResolveMyOrganisationRequest{UserId: userID})
	require.NoError(t, err)
	assert.Equal(t, org.OrganisationId, resolveResp.GetOrganisationId())
	assert.Equal(t, b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_HR_MANAGER, resolveResp.GetRole())
	assert.NotEmpty(t, resolveResp.GetOrganisationName())

	// REMOVE
	_, err = svc.RemoveOrgMember(ctx, &b2bsvcv1.RemoveOrgMemberRequest{
		OrganisationId: org.OrganisationId,
		MemberId:       member.MemberId,
	})
	require.NoError(t, err)
}

// ─── DEPARTMENT ───────────────────────────────────────────────────────────────

func TestB2BService_Department_CRUD(t *testing.T) {
	ctx := context.Background()
	svc, _ := testSvc(t)

	deptName := fmt.Sprintf("Svc Dept %s", uuid.NewString()[:8])

	// CREATE
	createResp, err := svc.CreateDepartment(ctx, &b2bsvcv1.CreateDepartmentRequest{
		Name:       deptName,
		BusinessId: svcTestBusinessID,
	})
	require.NoError(t, err)
	dept := createResp.GetDepartment()
	require.NotNil(t, dept)
	assert.Equal(t, deptName, dept.Name)
	assert.Equal(t, svcTestBusinessID, dept.BusinessId)
	t.Logf("created dept: id=%s name=%s", dept.DepartmentId, dept.Name)

	// GET
	getResp, err := svc.GetDepartment(ctx, &b2bsvcv1.GetDepartmentRequest{DepartmentId: dept.DepartmentId})
	require.NoError(t, err)
	assert.Equal(t, dept.DepartmentId, getResp.GetDepartment().DepartmentId)

	// LIST
	listResp, err := svc.ListDepartments(ctx, &b2bsvcv1.ListDepartmentsRequest{
		BusinessId: svcTestBusinessID,
		PageSize:   50,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, listResp.GetTotalCount(), int32(1))
	found := false
	for _, d := range listResp.GetDepartments() {
		if d.DepartmentId == dept.DepartmentId { found = true }
	}
	assert.True(t, found)

	// UPDATE
	updateResp, err := svc.UpdateDepartment(ctx, &b2bsvcv1.UpdateDepartmentRequest{
		DepartmentId: dept.DepartmentId,
		Name:         deptName + " Updated",
	})
	require.NoError(t, err)
	assert.Equal(t, deptName+" Updated", updateResp.GetDepartment().Name)

	// DELETE
	_, err = svc.DeleteDepartment(ctx, &b2bsvcv1.DeleteDepartmentRequest{DepartmentId: dept.DepartmentId})
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetDepartment(ctx, &b2bsvcv1.GetDepartmentRequest{DepartmentId: dept.DepartmentId})
	require.Error(t, err)
}

// ─── EMPLOYEE ─────────────────────────────────────────────────────────────────

func TestB2BService_Employee_CRUD(t *testing.T) {
	ctx := context.Background()
	svc, _ := testSvc(t)

	// Create a department first
	deptResp, err := svc.CreateDepartment(ctx, &b2bsvcv1.CreateDepartmentRequest{
		Name:       fmt.Sprintf("EmpSvc Dept %s", uuid.NewString()[:8]),
		BusinessId: svcTestBusinessID,
	})
	require.NoError(t, err)
	dept := deptResp.GetDepartment()
	t.Cleanup(func() {
		_, _ = svc.DeleteDepartment(ctx, &b2bsvcv1.DeleteDepartmentRequest{DepartmentId: dept.DepartmentId})
	})

	empName := fmt.Sprintf("Svc Employee %s", uuid.NewString()[:8])
	empID := fmt.Sprintf("EMP-%s", uuid.NewString()[:6])

	// CREATE
	createResp, err := svc.CreateEmployee(ctx, &b2bsvcv1.CreateEmployeeRequest{
		Name:              empName,
		EmployeeId:        empID,
		DepartmentId:      dept.DepartmentId,
		BusinessId:        svcTestBusinessID,
		InsuranceCategory: commonv1.InsuranceType_INSURANCE_TYPE_HEALTH,
		AssignedPlanId:    svcPlanHealth1,
		CoverageAmount:    svcMoney(50000),
		NumberOfDependent: 2,
		Email:             "svctest@example.com",
		MobileNumber:      "+8801711223344",
		DateOfBirth:       "1990-01-15",
		DateOfJoining:     "2023-06-01",
		Gender:            b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE,
	})
	require.NoError(t, err)
	empView := createResp.GetEmployee()
	require.NotNil(t, empView)
	emp := empView.GetEmployee()
	require.NotNil(t, emp)
	assert.Equal(t, empName, emp.Name)
	assert.Equal(t, empID, emp.EmployeeId)
	assert.Equal(t, dept.DepartmentId, emp.DepartmentId)
	assert.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_HEALTH, emp.InsuranceCategory)
	assert.Equal(t, b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE, emp.Gender)
	assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE, emp.Status)
	assert.NotNil(t, emp.CoverageAmount)
	assert.InDelta(t, 50000.0, emp.CoverageAmount.DecimalAmount, 0.01)
	// EmployeeView should have department name enriched
	assert.NotEmpty(t, empView.DepartmentName)
	t.Logf("created employee: uuid=%s name=%s dept_name=%s", emp.EmployeeUuid, emp.Name, empView.DepartmentName)

	// GET
	getResp, err := svc.GetEmployee(ctx, &b2bsvcv1.GetEmployeeRequest{EmployeeUuid: emp.EmployeeUuid})
	require.NoError(t, err)
	assert.Equal(t, emp.EmployeeUuid, getResp.GetEmployee().GetEmployee().EmployeeUuid)
	assert.Equal(t, empName, getResp.GetEmployee().GetEmployee().Name)

	// LIST by businessID
	listResp, err := svc.ListEmployees(ctx, &b2bsvcv1.ListEmployeesRequest{
		BusinessId: svcTestBusinessID,
		PageSize:   50,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, listResp.GetTotalCount(), int32(1))
	found := false
	for _, ev := range listResp.GetEmployees() {
		if ev.GetEmployee().EmployeeUuid == emp.EmployeeUuid { found = true }
	}
	assert.True(t, found)

	// LIST by departmentID
	deptListResp, err := svc.ListEmployees(ctx, &b2bsvcv1.ListEmployeesRequest{
		DepartmentId: dept.DepartmentId,
		PageSize:     50,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, deptListResp.GetTotalCount(), int32(1))

	// UPDATE name + email
	updateResp, err := svc.UpdateEmployee(ctx, &b2bsvcv1.UpdateEmployeeRequest{
		EmployeeUuid: emp.EmployeeUuid,
		Name:         empName + " Updated",
		Email:        "updated@example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, empName+" Updated", updateResp.GetEmployee().GetEmployee().Name)
	assert.Equal(t, "updated@example.com", updateResp.GetEmployee().GetEmployee().Email)

	// UPDATE status INACTIVE
	statusResp, err := svc.UpdateEmployee(ctx, &b2bsvcv1.UpdateEmployeeRequest{
		EmployeeUuid: emp.EmployeeUuid,
		Status:       b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE,
	})
	require.NoError(t, err)
	assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE, statusResp.GetEmployee().GetEmployee().Status)

	// DELETE
	_, err = svc.DeleteEmployee(ctx, &b2bsvcv1.DeleteEmployeeRequest{EmployeeUuid: emp.EmployeeUuid})
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetEmployee(ctx, &b2bsvcv1.GetEmployeeRequest{EmployeeUuid: emp.EmployeeUuid})
	require.Error(t, err)
}

func TestB2BService_Employee_ListFilters(t *testing.T) {
	ctx := context.Background()
	svc, _ := testSvc(t)

	// List all by businessID
	byBiz, err := svc.ListEmployees(ctx, &b2bsvcv1.ListEmployeesRequest{
		BusinessId: svcTestBusinessID,
		PageSize:   5,
	})
	require.NoError(t, err)
	t.Logf("employees by businessID: total=%d", byBiz.GetTotalCount())

	// List ACTIVE
	byActive, err := svc.ListEmployees(ctx, &b2bsvcv1.ListEmployeesRequest{
		BusinessId: svcTestBusinessID,
		Status:     b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE,
		PageSize:   5,
	})
	require.NoError(t, err)
	for _, ev := range byActive.GetEmployees() {
		assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE, ev.GetEmployee().Status)
	}
	t.Logf("employees ACTIVE: total=%d", byActive.GetTotalCount())

	// List INACTIVE
	byInactive, err := svc.ListEmployees(ctx, &b2bsvcv1.ListEmployeesRequest{
		BusinessId: svcTestBusinessID,
		Status:     b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE,
		PageSize:   5,
	})
	require.NoError(t, err)
	for _, ev := range byInactive.GetEmployees() {
		assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE, ev.GetEmployee().Status)
	}
	t.Logf("employees INACTIVE: total=%d", byInactive.GetTotalCount())
}

// ─── PURCHASE ORDER ───────────────────────────────────────────────────────────

func TestB2BService_PurchaseOrder_CRUD(t *testing.T) {
	ctx := context.Background()
	svc, gdb := testSvc(t)

	// Create dept for PO
	deptResp, err := svc.CreateDepartment(ctx, &b2bsvcv1.CreateDepartmentRequest{
		Name:       fmt.Sprintf("PO Svc Dept %s", uuid.NewString()[:8]),
		BusinessId: svcTestBusinessID,
	})
	require.NoError(t, err)
	dept := deptResp.GetDepartment()
	t.Cleanup(func() {
		_, _ = svc.DeleteDepartment(ctx, &b2bsvcv1.DeleteDepartmentRequest{DepartmentId: dept.DepartmentId})
	})

	// CREATE
	createResp, err := svc.CreatePurchaseOrder(ctx, &b2bsvcv1.CreatePurchaseOrderRequest{
		DepartmentId:       dept.DepartmentId,
		PlanId:             svcPlanHealth1,
		EmployeeCount:      20,
		NumberOfDependents: 5,
		CoverageAmount:     svcMoney(400000),
		RequestedBy:        uuid.NewString(),
		Notes:              "service layer live test PO",
	})
	require.NoError(t, err)
	poView := createResp.GetPurchaseOrder()
	require.NotNil(t, poView)
	po := poView.GetPurchaseOrder()
	require.NotNil(t, po)
	require.NotEmpty(t, po.PurchaseOrderId)
	assert.Equal(t, dept.DepartmentId, po.DepartmentId)
	assert.Equal(t, svcPlanHealth1, po.PlanId)
	assert.EqualValues(t, 20, po.EmployeeCount)
	assert.EqualValues(t, 5, po.NumberOfDependents)
	assert.NotNil(t, po.CoverageAmount)
	assert.InDelta(t, 400000.0, po.CoverageAmount.DecimalAmount, 0.01)
	assert.Equal(t, "service layer live test PO", po.Notes)
	t.Cleanup(func() { cleanupSvcPO(t, gdb, po.PurchaseOrderId) })
	t.Logf("created PO: id=%s number=%s status=%s dept_name=%s", po.PurchaseOrderId, po.PurchaseOrderNumber, po.Status, poView.DepartmentName)

	// GET
	getResp, err := svc.GetPurchaseOrder(ctx, &b2bsvcv1.GetPurchaseOrderRequest{PurchaseOrderId: po.PurchaseOrderId})
	require.NoError(t, err)
	assert.Equal(t, po.PurchaseOrderId, getResp.GetPurchaseOrder().GetPurchaseOrder().PurchaseOrderId)
	assert.NotNil(t, getResp.GetPurchaseOrder().GetPurchaseOrder().CoverageAmount)

	// LIST by businessID
	listResp, err := svc.ListPurchaseOrders(ctx, &b2bsvcv1.ListPurchaseOrdersRequest{
		BusinessId: svcTestBusinessID,
		PageSize:   50,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, listResp.GetTotalCount(), int32(1))
	found := false
	for _, pov := range listResp.GetPurchaseOrders() {
		if pov.GetPurchaseOrder().PurchaseOrderId == po.PurchaseOrderId { found = true }
	}
	assert.True(t, found)

	// LIST by status SUBMITTED
	submittedResp, err := svc.ListPurchaseOrders(ctx, &b2bsvcv1.ListPurchaseOrdersRequest{
		BusinessId: svcTestBusinessID,
		Status:     b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED,
		PageSize:   50,
	})
	require.NoError(t, err)
	for _, pov := range submittedResp.GetPurchaseOrders() {
		assert.Equal(t, b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED, pov.GetPurchaseOrder().Status)
	}
}

// ─── CATALOG ─────────────────────────────────────────────────────────────────

func TestB2BService_Catalog_ListAndResolve(t *testing.T) {
	ctx := context.Background()
	svc, _ := testSvc(t)

	listResp, err := svc.ListPurchaseOrderCatalog(ctx, &b2bsvcv1.ListPurchaseOrderCatalogRequest{})
	require.NoError(t, err)
	t.Logf("catalog plans: %d", len(listResp.GetItems()))
	for i, item := range listResp.GetItems() {
		if i >= 5 { break }
		assert.NotEmpty(t, item.PlanId)
		assert.NotEmpty(t, item.ProductName)
		t.Logf("  plan[%d]: id=%s name=%q product=%q", i, item.PlanId, item.PlanName, item.ProductName)
	}
}

// ─── SMOKE TEST ───────────────────────────────────────────────────────────────

func TestB2BService_LiveDB_Smoke(t *testing.T) {
	ctx := context.Background()
	svc, _ := testSvc(t)

	depts, err := svc.ListDepartments(ctx, &b2bsvcv1.ListDepartmentsRequest{PageSize: 5})
	require.NoError(t, err)
	t.Logf("smoke: departments total=%d fetched=%d", depts.GetTotalCount(), len(depts.GetDepartments()))
	if len(depts.GetDepartments()) > 0 {
		d, err := svc.GetDepartment(ctx, &b2bsvcv1.GetDepartmentRequest{DepartmentId: depts.GetDepartments()[0].DepartmentId})
		require.NoError(t, err)
		require.NotNil(t, d.GetDepartment())
	}

	emps, err := svc.ListEmployees(ctx, &b2bsvcv1.ListEmployeesRequest{PageSize: 5})
	require.NoError(t, err)
	t.Logf("smoke: employees total=%d fetched=%d", emps.GetTotalCount(), len(emps.GetEmployees()))
	if len(emps.GetEmployees()) > 0 {
		e, err := svc.GetEmployee(ctx, &b2bsvcv1.GetEmployeeRequest{EmployeeUuid: emps.GetEmployees()[0].GetEmployee().EmployeeUuid})
		require.NoError(t, err)
		require.NotNil(t, e.GetEmployee())
	}

	catalog, err := svc.ListPurchaseOrderCatalog(ctx, &b2bsvcv1.ListPurchaseOrderCatalogRequest{})
	require.NoError(t, err)
	t.Logf("smoke: catalog plans=%d", len(catalog.GetItems()))

	pos, err := svc.ListPurchaseOrders(ctx, &b2bsvcv1.ListPurchaseOrdersRequest{PageSize: 5})
	require.NoError(t, err)
	t.Logf("smoke: purchase orders total=%d fetched=%d", pos.GetTotalCount(), len(pos.GetPurchaseOrders()))
	if len(pos.GetPurchaseOrders()) > 0 {
		po, err := svc.GetPurchaseOrder(ctx, &b2bsvcv1.GetPurchaseOrderRequest{
			PurchaseOrderId: pos.GetPurchaseOrders()[0].GetPurchaseOrder().PurchaseOrderId,
		})
		require.NoError(t, err)
		require.NotNil(t, po.GetPurchaseOrder())
	}

	orgs, err := svc.ListOrganisations(ctx, &b2bsvcv1.ListOrganisationsRequest{PageSize: 5})
	require.NoError(t, err)
	t.Logf("smoke: organisations total=%d fetched=%d", orgs.GetTotalCount(), len(orgs.GetOrganisations()))
	if len(orgs.GetOrganisations()) > 0 {
		o, err := svc.GetOrganisation(ctx, &b2bsvcv1.GetOrganisationRequest{OrganisationId: orgs.GetOrganisations()[0].OrganisationId})
		require.NoError(t, err)
		require.NotNil(t, o.GetOrganisation())
	}
}

// ─── original test (preserved) ────────────────────────────────────────────────

func TestB2BService_Live_PurchaseOrdersResolveProductAndPlanNames(t *testing.T) {
	ctx := context.Background()
	svc, _ := testSvc(t)
	_ = time.Now()

	resp, err := svc.ListPurchaseOrders(ctx, &b2bsvcv1.ListPurchaseOrdersRequest{
		BusinessId: svcTestBusinessID,
		PageSize:   10,
	})
	require.NoError(t, err)
	t.Logf("purchase orders: total=%d", resp.GetTotalCount())
	for _, pov := range resp.GetPurchaseOrders() {
		po := pov.GetPurchaseOrder()
		t.Logf("  po=%s plan=%s product=%s status=%s dept=%s plan_name=%s product_name=%s",
			po.PurchaseOrderId, po.PlanId, po.ProductId, po.Status,
			pov.DepartmentName, pov.PlanName, pov.ProductName)
	}
}

