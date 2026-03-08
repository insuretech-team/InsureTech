package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ─── constants (seeded in dev DB) ─────────────────────────────────────────────

const (
	realBusinessID = "22222222-2222-2222-2222-222222222001"
	realTenantID   = "11111111-1111-1111-1111-111111111111"

	// seeded department IDs
	deptHR      = "44444444-4444-4444-4444-444444444001"
	deptFinance = "44444444-4444-4444-4444-444444444002"
	deptIT      = "44444444-4444-4444-4444-444444444003"

	// seeded employee UUIDs
	empJohnDoe    = "66666666-6666-6666-6666-666666666001"
	empJaneSmith  = "66666666-6666-6666-6666-666666666002"
	empBobJohnson = "66666666-6666-6666-6666-666666666003"

	// seeded purchase order IDs
	poApproved  = "77777777-7777-7777-7777-777777777001"
	poSubmitted = "77777777-7777-7777-7777-777777777002"
	poFulfilled = "77777777-7777-7777-7777-777777777003"

	// seeded plan IDs (for catalog)
	planHealth1 = "55555555-5555-5555-5555-555555555001"
	planHealth2 = "55555555-5555-5555-5555-555555555002"
	planLife1   = "55555555-5555-5555-5555-555555555004"
)

// ─── shared DB bootstrap ──────────────────────────────────────────────────────

var (
	b2bLiveOnce sync.Once
	b2bLiveDB   *gorm.DB
	b2bLiveErr  error
)

func testB2BLiveDB(t *testing.T) *gorm.DB {
	t.Helper()
	if os.Getenv("INSURETECH_LIVE_DB_TESTS") != "1" {
		t.Skip("skipping live DB test: set INSURETECH_LIVE_DB_TESTS=1 to run")
	}
	b2bLiveOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		cfgPath, err := config.ResolveConfigPath("database.yaml")
		if err != nil {
			b2bLiveErr = err
			return
		}
		b2bLiveErr = db.InitializeManagerForService(cfgPath)
		if b2bLiveErr == nil {
			b2bLiveDB = db.GetDB()
		}
	})
	if b2bLiveErr != nil {
		t.Fatalf("live DB init: %v", b2bLiveErr)
	}
	if b2bLiveDB == nil {
		t.Fatal("live DB is nil")
	}
	return b2bLiveDB
}

func newRepo(t *testing.T) *PortalRepository {
	t.Helper()
	return NewPortalRepository(testB2BLiveDB(t))
}

func testMoney(amount float64) *commonv1.Money {
	return &commonv1.Money{
		Amount:        int64(amount * 100),
		Currency:      "BDT",
		DecimalAmount: amount,
	}
}

// cleanupOrg hard-deletes a test organisation (no soft-delete in interface)
func cleanupOrg(t *testing.T, orgID string) {
	t.Helper()
	testB2BLiveDB(t).Exec(
		"UPDATE b2b_schema.organisations SET deleted_at = NOW() WHERE organisation_id = $1", orgID,
	)
}

// cleanupPO soft-deletes a test purchase order via deleted_at
func cleanupPO(t *testing.T, poID string) {
	t.Helper()
	testB2BLiveDB(t).Exec(
		"UPDATE b2b_schema.purchase_orders SET deleted_at = NOW() WHERE purchase_order_id = $1", poID,
	)
}

// ─── ORGANISATION CRUD ────────────────────────────────────────────────────────

func TestPortalRepository_Organisation_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	orgName := fmt.Sprintf("Test Org %s", uuid.NewString()[:8])
	orgCode := fmt.Sprintf("TST%s", uuid.NewString()[:4])
	orgID := uuid.NewString()

	t.Cleanup(func() { cleanupOrg(t, orgID) })

	// CREATE
	org, err := repo.CreateOrganisation(ctx, domain.OrganisationCreateInput{
		OrganisationID: orgID,
		TenantID:       realTenantID,
		Name:           orgName,
		Code:           orgCode,
		Industry:       "Technology",
		ContactEmail:   "test@example.com",
		ContactPhone:   "+8801700000000",
		Address:        "Dhaka, Bangladesh",
	})
	require.NoError(t, err)
	require.NotNil(t, org)
	assert.Equal(t, orgID, org.OrganisationId)
	assert.Equal(t, orgName, org.Name)
	assert.Equal(t, realTenantID, org.TenantId)
	assert.Equal(t, "Technology", org.Industry)
	assert.Equal(t, "test@example.com", org.ContactEmail)
	assert.NotNil(t, org.CreatedAt)
	t.Logf("created org: id=%s name=%s code=%s", org.OrganisationId, org.Name, org.Code)

	// GET
	fetched, err := repo.GetOrganisation(ctx, orgID)
	require.NoError(t, err)
	assert.Equal(t, orgID, fetched.OrganisationId)
	assert.Equal(t, orgName, fetched.Name)

	// LIST
	orgs, total, err := repo.ListOrganisations(ctx, 50, 0, realTenantID, b2bv1.OrganisationStatus_ORGANISATION_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	found := false
	for _, o := range orgs {
		if o.OrganisationId == orgID {
			found = true
		}
	}
	assert.True(t, found)

	// UPDATE
	updated, err := repo.UpdateOrganisation(ctx, domain.OrganisationUpdateInput{
		OrganisationID: orgID,
		Name:           orgName + " Updated",
		ContactEmail:   "updated@example.com",
		Industry:       "Finance",
	})
	require.NoError(t, err)
	assert.Equal(t, orgName+" Updated", updated.Name)
	assert.Equal(t, "Finance", updated.Industry)
	assert.Equal(t, "updated@example.com", updated.ContactEmail)

	// UPDATE STATUS
	statusUpdated, err := repo.UpdateOrganisation(ctx, domain.OrganisationUpdateInput{
		OrganisationID: orgID,
		Status:         b2bv1.OrganisationStatus_ORGANISATION_STATUS_INACTIVE,
	})
	require.NoError(t, err)
	assert.Equal(t, b2bv1.OrganisationStatus_ORGANISATION_STATUS_INACTIVE, statusUpdated.Status)
}

func TestPortalRepository_Organisation_GetNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)
	_, err := repo.GetOrganisation(ctx, uuid.NewString())
	require.Error(t, err)
}

// ─── ORG MEMBER ───────────────────────────────────────────────────────────────

func TestPortalRepository_OrgMember_AddResolveRemove(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	// create a fresh org
	orgID := uuid.NewString()
	t.Cleanup(func() { cleanupOrg(t, orgID) })
	_, err := repo.CreateOrganisation(ctx, domain.OrganisationCreateInput{
		OrganisationID: orgID,
		TenantID:       realTenantID,
		Name:           fmt.Sprintf("MemberTest Org %s", uuid.NewString()[:8]),
		Code:           fmt.Sprintf("MBR%s", uuid.NewString()[:4]),
	})
	require.NoError(t, err)

	memberID := uuid.NewString()
	userID := uuid.NewString()

	// ADD
	member, err := repo.AddOrgMember(ctx, domain.OrgMemberCreateInput{
		MemberID:       memberID,
		OrganisationID: orgID,
		UserID:         userID,
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_HR_MANAGER,
	})
	require.NoError(t, err)
	require.NotNil(t, member)
	assert.Equal(t, memberID, member.MemberId)
	assert.Equal(t, orgID, member.OrganisationId)
	assert.Equal(t, userID, member.UserId)
	assert.Equal(t, b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_HR_MANAGER, member.Role)
	assert.NotNil(t, member.JoinedAt)
	t.Logf("added member: member_id=%s user_id=%s role=%s", member.MemberId, member.UserId, member.Role)

	// RESOLVE
	resolvedOrgID, resolvedRole, resolvedName, err := repo.ResolveOrganisationByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, orgID, resolvedOrgID)
	assert.Equal(t, b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_HR_MANAGER, resolvedRole)
	assert.NotEmpty(t, resolvedName)
	t.Logf("resolved: org_id=%s name=%s role=%s", resolvedOrgID, resolvedName, resolvedRole)

	// REMOVE (hard delete in repo)
	err = repo.RemoveOrgMember(ctx, orgID, memberID)
	require.NoError(t, err)

	// After removal, resolve should return not found
	resolvedOrgID2, _, _, err2 := repo.ResolveOrganisationByUserID(ctx, userID)
	assert.True(t, resolvedOrgID2 == "" || err2 != nil, "removed member should not resolve to any org")
}

func TestPortalRepository_ResolveOrganisation_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)
	orgID, _, _, err := repo.ResolveOrganisationByUserID(ctx, uuid.NewString())
	assert.True(t, orgID == "" || err != nil)
}

// ─── DEPARTMENT CRUD ──────────────────────────────────────────────────────────

func TestPortalRepository_Department_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	deptID := uuid.NewString()
	deptName := fmt.Sprintf("Test Dept %s", uuid.NewString()[:8])

	// CREATE
	dept, err := repo.CreateDepartment(ctx, domain.DepartmentCreateInput{
		DepartmentID: deptID,
		Name:         deptName,
		BusinessID:   realBusinessID,
	})
	require.NoError(t, err)
	require.NotNil(t, dept)
	assert.Equal(t, deptID, dept.DepartmentId)
	assert.Equal(t, deptName, dept.Name)
	assert.Equal(t, realBusinessID, dept.BusinessId)
	assert.EqualValues(t, 0, dept.EmployeeNo)
	assert.NotNil(t, dept.CreatedAt)
	t.Logf("created dept: id=%s name=%s", dept.DepartmentId, dept.Name)

	// GET
	fetched, err := repo.GetDepartment(ctx, deptID)
	require.NoError(t, err)
	assert.Equal(t, deptID, fetched.DepartmentId)
	assert.Equal(t, deptName, fetched.Name)

	// LIST by businessID
	depts, total, err := repo.ListDepartments(ctx, 100, 0, realBusinessID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	found := false
	for _, d := range depts {
		if d.DepartmentId == deptID {
			found = true
		}
	}
	assert.True(t, found, "created department should appear in list")

	// LIST unfiltered
	_, allTotal, err := repo.ListDepartments(ctx, 100, 0, "")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, allTotal, total)

	// GET DEPARTMENT NAMES
	names, err := repo.GetDepartmentNames(ctx, []string{deptID})
	require.NoError(t, err)
	assert.Equal(t, deptName, names[deptID])

	// UPDATE
	updatedName := deptName + " Updated"
	updated, err := repo.UpdateDepartment(ctx, domain.DepartmentUpdateInput{
		DepartmentID: deptID,
		Name:         updatedName,
	})
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.Name)

	// verify update persisted
	refetched, err := repo.GetDepartment(ctx, deptID)
	require.NoError(t, err)
	assert.Equal(t, updatedName, refetched.Name)

	// DELETE (soft-delete via deleted_at)
	err = repo.DeleteDepartment(ctx, deptID)
	require.NoError(t, err)

	// verify soft-deleted
	_, err = repo.GetDepartment(ctx, deptID)
	require.Error(t, err, "soft-deleted department should not be found")
}

func TestPortalRepository_Department_GetNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)
	_, err := repo.GetDepartment(ctx, uuid.NewString())
	require.Error(t, err)
}

func TestPortalRepository_Department_DeleteBlocked_ActiveEmployees(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	// create department
	deptID := uuid.NewString()
	_, err := repo.CreateDepartment(ctx, domain.DepartmentCreateInput{
		DepartmentID: deptID,
		Name:         fmt.Sprintf("Block Dept %s", uuid.NewString()[:8]),
		BusinessID:   realBusinessID,
	})
	require.NoError(t, err)

	// create active employee in that department
	empUUID := uuid.NewString()
	emp, err := repo.CreateEmployee(ctx, domain.EmployeeCreateInput{
		EmployeeUUID:  empUUID,
		Name:          "Block Test Employee",
		EmployeeID:    fmt.Sprintf("EMP-%s", uuid.NewString()[:6]),
		DepartmentID:  deptID,
		BusinessID:    realBusinessID,
		DateOfJoining: time.Now().Format("2006-01-02"),
	})
	require.NoError(t, err)
	require.NotNil(t, emp)

	// delete should be blocked
	err = repo.DeleteDepartment(ctx, deptID)
	require.Error(t, err, "delete should be blocked when active employees exist")
	t.Logf("correctly blocked: %v", err)

	// cleanup
	require.NoError(t, repo.DeleteEmployee(ctx, empUUID))
	require.NoError(t, repo.DeleteDepartment(ctx, deptID))
}

func TestPortalRepository_Department_ListPagination(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	page1, total, err := repo.ListDepartments(ctx, 2, 0, realBusinessID)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(page1), 2)

	if total > 2 {
		page2, _, err := repo.ListDepartments(ctx, 2, 2, realBusinessID)
		require.NoError(t, err)
		assert.NotEmpty(t, page2)
		ids1 := make(map[string]struct{})
		for _, d := range page1 {
			ids1[d.DepartmentId] = struct{}{}
		}
		for _, d := range page2 {
			_, overlap := ids1[d.DepartmentId]
			assert.False(t, overlap, "pages should not overlap")
		}
	}
}

// ─── EMPLOYEE CRUD ────────────────────────────────────────────────────────────

func TestPortalRepository_Employee_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	// create a fresh department to attach employee
	deptID := uuid.NewString()
	_, err := repo.CreateDepartment(ctx, domain.DepartmentCreateInput{
		DepartmentID: deptID,
		Name:         fmt.Sprintf("EmpTest Dept %s", uuid.NewString()[:8]),
		BusinessID:   realBusinessID,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = repo.DeleteDepartment(ctx, deptID) })

	empUUID := uuid.NewString()
	empName := fmt.Sprintf("Test Employee %s", uuid.NewString()[:8])
	empID := fmt.Sprintf("EMP-%s", uuid.NewString()[:6])

	// CREATE
	emp, err := repo.CreateEmployee(ctx, domain.EmployeeCreateInput{
		EmployeeUUID:      empUUID,
		Name:              empName,
		EmployeeID:        empID,
		DepartmentID:      deptID,
		BusinessID:        realBusinessID,
		InsuranceCategory: commonv1.InsuranceType_INSURANCE_TYPE_HEALTH,
		AssignedPlanID:    planHealth1,
		CoverageAmount:    testMoney(50000),
		NumberOfDependent: 2,
		Email:             "employee@example.com",
		MobileNumber:      "+8801700000001",
		DateOfBirth:       "1990-01-15",
		DateOfJoining:     "2023-06-01",
		Gender:            b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE,
	})
	require.NoError(t, err)
	require.NotNil(t, emp)
	assert.Equal(t, empUUID, emp.EmployeeUuid)
	assert.Equal(t, empName, emp.Name)
	assert.Equal(t, empID, emp.EmployeeId)
	assert.Equal(t, deptID, emp.DepartmentId)
	assert.Equal(t, realBusinessID, emp.BusinessId)
	assert.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_HEALTH, emp.InsuranceCategory)
	assert.Equal(t, planHealth1, emp.AssignedPlanId)
	assert.NotNil(t, emp.CoverageAmount)
	assert.InDelta(t, 50000.0, emp.CoverageAmount.DecimalAmount, 0.01)
	assert.EqualValues(t, 2, emp.NumberOfDependent)
	assert.Equal(t, "employee@example.com", emp.Email)
	assert.Equal(t, "+8801700000001", emp.MobileNumber)
	assert.Equal(t, b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE, emp.Gender)
	assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE, emp.Status)
	assert.NotNil(t, emp.CreatedAt)
	t.Logf("created employee: uuid=%s name=%s empID=%s", emp.EmployeeUuid, emp.Name, emp.EmployeeId)

	// GET
	fetched, err := repo.GetEmployee(ctx, empUUID)
	require.NoError(t, err)
	assert.Equal(t, empUUID, fetched.EmployeeUuid)
	assert.Equal(t, empName, fetched.Name)
	assert.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_HEALTH, fetched.InsuranceCategory)

	// LIST by businessID
	emps, total, err := repo.ListEmployees(ctx, 100, 0, "", realBusinessID, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	found := false
	for _, e := range emps {
		if e.EmployeeUuid == empUUID {
			found = true
		}
	}
	assert.True(t, found)

	// LIST by departmentID
	deptEmps, deptTotal, err := repo.ListEmployees(ctx, 100, 0, deptID, "", b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, deptTotal, int64(1))
	found = false
	for _, e := range deptEmps {
		if e.EmployeeUuid == empUUID {
			found = true
		}
	}
	assert.True(t, found)

	// LIST by status ACTIVE
	activeEmps, activeTotal, err := repo.ListEmployees(ctx, 100, 0, "", realBusinessID, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, activeTotal, int64(1))
	for _, e := range activeEmps {
		assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE, e.Status)
	}

	// GET DEPARTMENT NAMES enrichment
	names, err := repo.GetDepartmentNames(ctx, []string{deptID})
	require.NoError(t, err)
	assert.NotEmpty(t, names[deptID])

	// UPDATE
	updated, err := repo.UpdateEmployee(ctx, domain.EmployeeUpdateInput{
		EmployeeUUID:      empUUID,
		Name:              empName + " Updated",
		Email:             "updated@example.com",
		MobileNumber:      "+8801700000002",
		InsuranceCategory: commonv1.InsuranceType_INSURANCE_TYPE_LIFE,
		Gender:            b2bv1.EmployeeGender_EMPLOYEE_GENDER_FEMALE,
		CoverageAmount:    testMoney(75000),
		NumberOfDependent: 3,
		AssignedPlanID:    planLife1,
		DateOfBirth:       "1992-05-20",
		DateOfJoining:     "2024-01-01",
	})
	require.NoError(t, err)
	assert.Equal(t, empName+" Updated", updated.Name)
	assert.Equal(t, "updated@example.com", updated.Email)
	assert.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_LIFE, updated.InsuranceCategory)
	assert.Equal(t, b2bv1.EmployeeGender_EMPLOYEE_GENDER_FEMALE, updated.Gender)
	assert.InDelta(t, 75000.0, updated.CoverageAmount.DecimalAmount, 0.01)
	assert.EqualValues(t, 3, updated.NumberOfDependent)

	// UPDATE STATUS - deactivate
	deactivated, err := repo.UpdateEmployee(ctx, domain.EmployeeUpdateInput{
		EmployeeUUID: empUUID,
		Status:       b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE,
	})
	require.NoError(t, err)
	assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE, deactivated.Status)

	// UPDATE STATUS - reactivate
	reactivated, err := repo.UpdateEmployee(ctx, domain.EmployeeUpdateInput{
		EmployeeUUID: empUUID,
		Status:       b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE,
	})
	require.NoError(t, err)
	assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE, reactivated.Status)

	// DELETE (soft-delete)
	err = repo.DeleteEmployee(ctx, empUUID)
	require.NoError(t, err)

	// verify soft-deleted
	_, err = repo.GetEmployee(ctx, empUUID)
	require.Error(t, err, "soft-deleted employee should not be found")
}

func TestPortalRepository_Employee_GetNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)
	_, err := repo.GetEmployee(ctx, uuid.NewString())
	require.Error(t, err)
}

func TestPortalRepository_Employee_ListPagination(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	page1, total, err := repo.ListEmployees(ctx, 2, 0, "", realBusinessID, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(page1), 2)

	if total > 2 {
		page2, _, err := repo.ListEmployees(ctx, 2, 2, "", realBusinessID, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED)
		require.NoError(t, err)
		assert.NotEmpty(t, page2)
		ids1 := make(map[string]struct{})
		for _, e := range page1 {
			ids1[e.EmployeeUuid] = struct{}{}
		}
		for _, e := range page2 {
			_, overlap := ids1[e.EmployeeUuid]
			assert.False(t, overlap, "pages should not overlap")
		}
	}
}

// ─── PURCHASE ORDER CRUD ──────────────────────────────────────────────────────

func TestPortalRepository_PurchaseOrder_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	// Create a fresh department for the PO
	deptID := uuid.NewString()
	_, err := repo.CreateDepartment(ctx, domain.DepartmentCreateInput{
		DepartmentID: deptID,
		Name:         fmt.Sprintf("PO Test Dept %s", uuid.NewString()[:8]),
		BusinessID:   realBusinessID,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = repo.DeleteDepartment(ctx, deptID) })

	poID := uuid.NewString()
	poNumber := fmt.Sprintf("PO-TEST-%s", uuid.NewString()[:8])
	t.Cleanup(func() { cleanupPO(t, poID) })

	// CREATE
	po, err := repo.CreatePurchaseOrder(ctx, domain.PurchaseOrderCreateInput{
		PurchaseOrderID:     poID,
		PurchaseOrderNumber: poNumber,
		BusinessID:          realBusinessID,
		DepartmentID:        deptID,
		ProductID:           "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1",
		PlanID:              planHealth1,
		InsuranceCategory:   commonv1.InsuranceType_INSURANCE_TYPE_HEALTH,
		EmployeeCount:       25,
		NumberOfDependents:  10,
		CoverageAmount:      testMoney(500000),
		EstimatedPremium:    testMoney(12500),
		Status:              b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED,
		RequestedBy:         uuid.NewString(),
		Notes:               "Live test purchase order",
	})
	require.NoError(t, err)
	require.NotNil(t, po)
	assert.Equal(t, poID, po.PurchaseOrderId)
	assert.Equal(t, poNumber, po.PurchaseOrderNumber)
	assert.Equal(t, realBusinessID, po.BusinessId)
	assert.Equal(t, deptID, po.DepartmentId)
	assert.Equal(t, planHealth1, po.PlanId)
	assert.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_HEALTH, po.InsuranceCategory)
	assert.EqualValues(t, 25, po.EmployeeCount)
	assert.EqualValues(t, 10, po.NumberOfDependents)
	assert.NotNil(t, po.CoverageAmount)
	assert.InDelta(t, 500000.0, po.CoverageAmount.DecimalAmount, 0.01)
	assert.NotNil(t, po.EstimatedPremium)
	assert.InDelta(t, 12500.0, po.EstimatedPremium.DecimalAmount, 0.01)
	assert.Equal(t, b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED, po.Status)
	assert.Equal(t, "Live test purchase order", po.Notes)
	assert.NotNil(t, po.CreatedAt)
	t.Logf("created PO: id=%s number=%s status=%s", po.PurchaseOrderId, po.PurchaseOrderNumber, po.Status)

	// GET
	fetched, err := repo.GetPurchaseOrder(ctx, poID)
	require.NoError(t, err)
	assert.Equal(t, poID, fetched.PurchaseOrderId)
	assert.Equal(t, poNumber, fetched.PurchaseOrderNumber)
	assert.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_HEALTH, fetched.InsuranceCategory)
	assert.NotNil(t, fetched.CoverageAmount)
	assert.NotNil(t, fetched.EstimatedPremium)
	assert.Equal(t, b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED, fetched.Status)

	// LIST by businessID
	orders, total, err := repo.ListPurchaseOrders(ctx, 100, 0, realBusinessID, b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	found := false
	for _, o := range orders {
		if o.PurchaseOrderId == poID {
			found = true
		}
	}
	assert.True(t, found, "created PO should appear in list")

	// LIST by status SUBMITTED
	submitted, submittedTotal, err := repo.ListPurchaseOrders(ctx, 100, 0, realBusinessID, b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, submittedTotal, int64(1))
	found = false
	for _, o := range submitted {
		if o.PurchaseOrderId == poID {
			found = true
		}
		assert.Equal(t, b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED, o.Status)
	}
	assert.True(t, found)

	// LIST unfiltered
	_, allTotal, err := repo.ListPurchaseOrders(ctx, 100, 0, "", b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, allTotal, total)
}

func TestPortalRepository_PurchaseOrder_GetNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)
	_, err := repo.GetPurchaseOrder(ctx, uuid.NewString())
	require.Error(t, err)
}

func TestPortalRepository_PurchaseOrder_ListPagination(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	page1, total, err := repo.ListPurchaseOrders(ctx, 1, 0, "", b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED)
	require.NoError(t, err)

	if total > 1 {
		require.Len(t, page1, 1)
		page2, _, err := repo.ListPurchaseOrders(ctx, 1, 1, "", b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED)
		require.NoError(t, err)
		require.Len(t, page2, 1)
		assert.NotEqual(t, page1[0].PurchaseOrderId, page2[0].PurchaseOrderId)
	}
}

// ─── CATALOG (read-only) ──────────────────────────────────────────────────────

func TestPortalRepository_Catalog_ListAll(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	items, err := repo.ListCatalogPlans(ctx)
	require.NoError(t, err)
	t.Logf("catalog plans found: %d", len(items))
	for i, item := range items {
		if i >= 5 {
			break
		}
		assert.NotEmpty(t, item.PlanID)
		assert.NotEmpty(t, item.ProductID)
		t.Logf("  plan[%d]: plan_id=%s plan_name=%q product_name=%q category=%s",
			i, item.PlanID, item.PlanName, item.ProductName, item.InsuranceCategory)
	}
}

func TestPortalRepository_Catalog_GetByPlanIDs_Seeded(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	seededIDs := []string{planHealth1, planHealth2, planLife1}
	plans, err := repo.GetCatalogPlansByPlanIDs(ctx, seededIDs)
	require.NoError(t, err)
	t.Logf("seeded plans resolved: %d/%d", len(plans), len(seededIDs))
	for planID, plan := range plans {
		assert.Equal(t, planID, plan.PlanID)
		assert.NotEmpty(t, plan.ProductName)
		assert.NotEmpty(t, plan.PlanName)
		t.Logf("  plan_id=%s name=%q product=%q premium=%v", plan.PlanID, plan.PlanName, plan.ProductName, plan.PremiumAmount)
	}
}

func TestPortalRepository_Catalog_GetByPlanIDs_Empty(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	empty, err := repo.GetCatalogPlansByPlanIDs(ctx, []string{})
	require.NoError(t, err)
	assert.Empty(t, empty)
}

func TestPortalRepository_Catalog_GetByPlanIDs_Missing(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	missing, err := repo.GetCatalogPlansByPlanIDs(ctx, []string{uuid.NewString()})
	require.NoError(t, err)
	assert.Empty(t, missing)
}

// ─── SMOKE TEST (read-only, all repos in one pass) ────────────────────────────

func TestPortalRepository_LiveDB_Smoke(t *testing.T) {
	ctx := context.Background()
	repo := newRepo(t)

	// departments
	depts, deptTotal, err := repo.ListDepartments(ctx, 5, 0, "")
	require.NoError(t, err)
	t.Logf("smoke: departments total=%d fetched=%d", deptTotal, len(depts))
	if len(depts) > 0 {
		d, err := repo.GetDepartment(ctx, depts[0].DepartmentId)
		require.NoError(t, err)
		require.NotNil(t, d)
		names, err := repo.GetDepartmentNames(ctx, []string{depts[0].DepartmentId})
		require.NoError(t, err)
		assert.Equal(t, depts[0].Name, names[depts[0].DepartmentId])
	}

	// employees
	emps, empTotal, err := repo.ListEmployees(ctx, 5, 0, "", "", b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	t.Logf("smoke: employees total=%d fetched=%d", empTotal, len(emps))
	if len(emps) > 0 {
		e, err := repo.GetEmployee(ctx, emps[0].EmployeeUuid)
		require.NoError(t, err)
		require.NotNil(t, e)
	}

	// catalog
	catalog, err := repo.ListCatalogPlans(ctx)
	require.NoError(t, err)
	t.Logf("smoke: catalog plans=%d", len(catalog))
	if len(catalog) > 0 {
		plans, err := repo.GetCatalogPlansByPlanIDs(ctx, []string{catalog[0].PlanID})
		require.NoError(t, err)
		assert.Contains(t, plans, catalog[0].PlanID)
	}

	// purchase orders
	pos, poTotal, err := repo.ListPurchaseOrders(ctx, 5, 0, "", b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	t.Logf("smoke: purchase orders total=%d fetched=%d", poTotal, len(pos))
	if len(pos) > 0 {
		po, err := repo.GetPurchaseOrder(ctx, pos[0].PurchaseOrderId)
		require.NoError(t, err)
		require.NotNil(t, po)
	}

	// organisations
	orgs, orgTotal, err := repo.ListOrganisations(ctx, 5, 0, "", b2bv1.OrganisationStatus_ORGANISATION_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	t.Logf("smoke: organisations total=%d fetched=%d", orgTotal, len(orgs))
	if len(orgs) > 0 {
		o, err := repo.GetOrganisation(ctx, orgs[0].OrganisationId)
		require.NoError(t, err)
		require.NotNil(t, o)
	}
}
