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
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ─── DB Setup ─────────────────────────────────────────────────────────────────
//
// Uses the shared db.Manager (database.yaml → DigitalOcean primary + Neon backup),
// the exact same setup as b2b_service_live_test.go.
// Set INSURETECH_LIVE_DB_TESTS=1 to enable; tests are skipped otherwise.

var (
	repoLiveOnce sync.Once
	repoLiveDB   *gorm.DB
	repoLiveErr  error
)

func testB2BDB(t *testing.T) *gorm.DB {
	t.Helper()
	if os.Getenv("INSURETECH_LIVE_DB_TESTS") != "1" {
		t.Skip("skipping live DB test: set INSURETECH_LIVE_DB_TESTS=1 to run")
	}
	repoLiveOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		cfgPath, err := config.ResolveConfigPath("database.yaml")
		if err != nil {
			repoLiveErr = err
			return
		}
		repoLiveErr = db.InitializeManagerForService(cfgPath)
		if repoLiveErr != nil {
			return
		}
		repoLiveDB = db.GetDB()
	})
	require.NoError(t, repoLiveErr, "live DB init")
	require.NotNil(t, repoLiveDB, "live DB handle")
	return repoLiveDB
}

func testRepo(t *testing.T) *PortalRepository {
	t.Helper()
	return NewPortalRepository(testB2BDB(t))
}

// ─── Cleanup helpers ──────────────────────────────────────────────────────────

func cleanupEmployee(ctx context.Context, t *testing.T, db *gorm.DB, employeeUUID string) {
	t.Helper()
	db.Exec("DELETE FROM b2b_schema.employees WHERE employee_uuid = $1", employeeUUID)
}

func cleanupDepartment(ctx context.Context, t *testing.T, db *gorm.DB, departmentID string) {
	t.Helper()
	db.Exec("DELETE FROM b2b_schema.departments WHERE department_id = $1", departmentID)
}

func cleanupOrganisation(ctx context.Context, t *testing.T, db *gorm.DB, orgID string) {
	t.Helper()
	db.Exec("DELETE FROM b2b_schema.org_members WHERE organisation_id = $1", orgID)
	db.Exec("DELETE FROM b2b_schema.organisations WHERE organisation_id = $1", orgID)
}

// ─── Organisation Tests ───────────────────────────────────────────────────────

func TestPortalRepository_Organisation_CRUD_LiveDB(t *testing.T) {
	ctx := context.Background()
	repo := testRepo(t)
	db := testB2BDB(t)

	orgID := uuid.NewString()
	code := fmt.Sprintf("TEST%d", time.Now().UnixMilli()%10000)

	t.Cleanup(func() { cleanupOrganisation(ctx, t, db, orgID) })

	// Create
	org, err := repo.CreateOrganisation(ctx, domain.OrganisationCreateInput{
		OrganisationID: orgID,
		TenantID:       uuid.NewString(),
		Name:           "Test Corp",
		Code:           code,
		Industry:       "Technology",
		ContactEmail:   "admin@testcorp.com",
		ContactPhone:   "+8801700000000",
		Address:        "Dhaka, Bangladesh",
	})
	require.NoError(t, err, "CreateOrganisation")
	require.Equal(t, orgID, org.GetOrganisationId())
	require.Equal(t, "Test Corp", org.GetName())
	assert.Equal(t, b2bv1.OrganisationStatus_ORGANISATION_STATUS_ACTIVE, org.GetStatus())

	// Get
	got, err := repo.GetOrganisation(ctx, orgID)
	require.NoError(t, err, "GetOrganisation")
	require.Equal(t, orgID, got.GetOrganisationId())
	require.Equal(t, code, got.GetCode())

	// List
	orgs, total, err := repo.ListOrganisations(ctx, 50, 0, "", b2bv1.OrganisationStatus_ORGANISATION_STATUS_UNSPECIFIED)
	require.NoError(t, err, "ListOrganisations")
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(orgs), 1)

	// Update
	updated, err := repo.UpdateOrganisation(ctx, domain.OrganisationUpdateInput{
		OrganisationID: orgID,
		Name:           "Test Corp Updated",
		Industry:       "FinTech",
	})
	require.NoError(t, err, "UpdateOrganisation")
	assert.Equal(t, "Test Corp Updated", updated.GetName())
	assert.Equal(t, "FinTech", updated.GetIndustry())
}

func TestPortalRepository_OrgMember_LiveDB(t *testing.T) {
	ctx := context.Background()
	repo := testRepo(t)
	db := testB2BDB(t)

	orgID := uuid.NewString()
	userID := uuid.NewString()
	code := fmt.Sprintf("ORG%d", time.Now().UnixMilli()%10000)

	t.Cleanup(func() { cleanupOrganisation(ctx, t, db, orgID) })

	// Seed org
	_, err := repo.CreateOrganisation(ctx, domain.OrganisationCreateInput{
		OrganisationID: orgID,
		TenantID:       uuid.NewString(),
		Name:           "Member Test Corp",
		Code:           code,
	})
	require.NoError(t, err)

	// AddOrgMember
	member, err := repo.AddOrgMember(ctx, domain.OrgMemberCreateInput{
		OrganisationID: orgID,
		UserID:         userID,
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN,
	})
	require.NoError(t, err, "AddOrgMember")
	require.Equal(t, orgID, member.GetOrganisationId())
	require.Equal(t, userID, member.GetUserId())
	assert.Equal(t, b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN, member.GetRole())

	// ResolveOrganisationByUserID — the core fix
	resolvedOrgID, role, orgName, err := repo.ResolveOrganisationByUserID(ctx, userID)
	require.NoError(t, err, "ResolveOrganisationByUserID")
	assert.Equal(t, orgID, resolvedOrgID)
	assert.Equal(t, b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN, role)
	assert.Equal(t, "Member Test Corp", orgName)

	// RemoveOrgMember
	err = repo.RemoveOrgMember(ctx, orgID, member.GetMemberId())
	require.NoError(t, err, "RemoveOrgMember")

	// After remove — resolve should fail
	_, _, _, err = repo.ResolveOrganisationByUserID(ctx, userID)
	assert.Error(t, err, "should not resolve after member removed")
}

// ─── Department Tests ─────────────────────────────────────────────────────────

func TestPortalRepository_Department_CRUD_LiveDB(t *testing.T) {
	ctx := context.Background()
	repo := testRepo(t)
	db := testB2BDB(t)

	deptID := uuid.NewString()
	bizID := uuid.NewString()

	t.Cleanup(func() { cleanupDepartment(ctx, t, db, deptID) })

	// Create
	dept, err := repo.CreateDepartment(ctx, domain.DepartmentCreateInput{
		DepartmentID: deptID,
		Name:         "Engineering",
		BusinessID:   bizID,
	})
	require.NoError(t, err, "CreateDepartment")
	require.Equal(t, deptID, dept.GetDepartmentId())
	require.Equal(t, "Engineering", dept.GetName())
	require.Equal(t, bizID, dept.GetBusinessId())
	assert.Equal(t, int32(0), dept.GetEmployeeNo())

	// Get
	got, err := repo.GetDepartment(ctx, deptID)
	require.NoError(t, err, "GetDepartment")
	assert.Equal(t, "Engineering", got.GetName())

	// List
	depts, total, err := repo.ListDepartments(ctx, 50, 0, bizID)
	require.NoError(t, err, "ListDepartments")
	assert.Equal(t, int64(1), total)
	require.Equal(t, 1, len(depts))
	assert.Equal(t, deptID, depts[0].GetDepartmentId())

	// Update
	updated, err := repo.UpdateDepartment(ctx, domain.DepartmentUpdateInput{
		DepartmentID: deptID,
		Name:         "Engineering & DevOps",
	})
	require.NoError(t, err, "UpdateDepartment")
	assert.Equal(t, "Engineering & DevOps", updated.GetName())

	// Delete (no employees, should succeed)
	err = repo.DeleteDepartment(ctx, deptID)
	require.NoError(t, err, "DeleteDepartment")

	// Get after delete should fail
	_, err = repo.GetDepartment(ctx, deptID)
	assert.Error(t, err, "should not find deleted department")
}

// ─── Employee Tests ───────────────────────────────────────────────────────────

func TestPortalRepository_Employee_CRUD_LiveDB(t *testing.T) {
	ctx := context.Background()
	repo := testRepo(t)
	db := testB2BDB(t)

	// Seed department
	deptID := uuid.NewString()
	bizID := uuid.NewString()
	empID := uuid.NewString()

	t.Cleanup(func() {
		cleanupEmployee(ctx, t, db, empID)
		cleanupDepartment(ctx, t, db, deptID)
	})

	_, err := repo.CreateDepartment(ctx, domain.DepartmentCreateInput{
		DepartmentID: deptID,
		Name:         "HR Department",
		BusinessID:   bizID,
	})
	require.NoError(t, err)

	// Create Employee
	emp, err := repo.CreateEmployee(ctx, domain.EmployeeCreateInput{
		EmployeeUUID:  empID,
		Name:          "Rahim Uddin",
		EmployeeID:    fmt.Sprintf("EMP-%d", time.Now().UnixMilli()),
		DepartmentID:  deptID,
		BusinessID:    bizID,
		Email:         "rahim@testcorp.com",
		MobileNumber:  "+8801711000001",
		DateOfBirth:   "1990-05-15",
		DateOfJoining: "2022-01-01",
		Gender:        b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE,
	})
	require.NoError(t, err, "CreateEmployee")
	require.Equal(t, empID, emp.GetEmployeeUuid())
	require.Equal(t, "Rahim Uddin", emp.GetName())
	assert.Equal(t, "rahim@testcorp.com", emp.GetEmail())
	assert.Equal(t, "+8801711000001", emp.GetMobileNumber())
	assert.Equal(t, "1990-05-15", emp.GetDateOfBirth())
	assert.Equal(t, b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE, emp.GetGender())
	assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE, emp.GetStatus())

	// Get
	got, err := repo.GetEmployee(ctx, empID)
	require.NoError(t, err, "GetEmployee")
	assert.Equal(t, empID, got.GetEmployeeUuid())
	assert.Equal(t, bizID, got.GetBusinessId())

	// List
	employees, total, err := repo.ListEmployees(ctx, 50, 0, "", bizID, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED)
	require.NoError(t, err, "ListEmployees")
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(employees), 1)

	// List with status filter
	activeEmployees, activeTotal, err := repo.ListEmployees(ctx, 50, 0, "", bizID, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE)
	require.NoError(t, err)
	assert.Equal(t, total, activeTotal) // all employees are active

	// List with department filter
	deptEmployees, deptTotal, err := repo.ListEmployees(ctx, 50, 0, deptID, bizID, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED)
	require.NoError(t, err)
	assert.Equal(t, int64(1), deptTotal)
	assert.Equal(t, 1, len(deptEmployees))
	_ = activeEmployees

	// GetDepartmentNames enrichment
	names, err := repo.GetDepartmentNames(ctx, []string{deptID})
	require.NoError(t, err)
	assert.Equal(t, "HR Department", names[deptID])

	// Update
	updated, err := repo.UpdateEmployee(ctx, domain.EmployeeUpdateInput{
		EmployeeUUID: empID,
		Name:         "Rahim Uddin Chowdhury",
		Email:        "rahim.updated@testcorp.com",
		Gender:       b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE,
	})
	require.NoError(t, err, "UpdateEmployee")
	assert.Equal(t, "Rahim Uddin Chowdhury", updated.GetName())
	assert.Equal(t, "rahim.updated@testcorp.com", updated.GetEmail())

	// Deactivate via UpdateEmployee status
	deactivated, err := repo.UpdateEmployee(ctx, domain.EmployeeUpdateInput{
		EmployeeUUID: empID,
		Status:       b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE,
	})
	require.NoError(t, err)
	assert.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE, deactivated.GetStatus())

	// Delete
	err = repo.DeleteEmployee(ctx, empID)
	require.NoError(t, err, "DeleteEmployee")

	// Get after delete should fail
	_, err = repo.GetEmployee(ctx, empID)
	assert.Error(t, err, "should not find deleted employee")
}

func TestPortalRepository_DeleteDepartment_WithActiveEmployees_LiveDB(t *testing.T) {
	ctx := context.Background()
	repo := testRepo(t)
	db := testB2BDB(t)

	deptID := uuid.NewString()
	bizID := uuid.NewString()
	empID := uuid.NewString()

	t.Cleanup(func() {
		cleanupEmployee(ctx, t, db, empID)
		cleanupDepartment(ctx, t, db, deptID)
	})

	_, err := repo.CreateDepartment(ctx, domain.DepartmentCreateInput{
		DepartmentID: deptID,
		Name:         "Finance",
		BusinessID:   bizID,
	})
	require.NoError(t, err)

	_, err = repo.CreateEmployee(ctx, domain.EmployeeCreateInput{
		EmployeeUUID:  empID,
		Name:          "Karim Ahmed",
		EmployeeID:    fmt.Sprintf("EMP-%d", time.Now().UnixMilli()),
		DepartmentID:  deptID,
		BusinessID:    bizID,
		DateOfJoining: "2023-06-01",
	})
	require.NoError(t, err)

	// Should fail — department has active employees
	err = repo.DeleteDepartment(ctx, deptID)
	require.Error(t, err, "should refuse to delete dept with active employees")
	assert.Contains(t, err.Error(), "active employee")
}
