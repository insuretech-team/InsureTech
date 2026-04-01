package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type fakeB2BRepo struct {
	listCatalogPlansFn            func(context.Context) ([]*domain.CatalogPlan, error)
	getCatalogPlansByPlanIDsFn    func(context.Context, []string) (map[string]*domain.CatalogPlan, error)
	listPurchaseOrdersFn          func(context.Context, int, int, string, b2bv1.PurchaseOrderStatus) ([]*b2bv1.PurchaseOrder, int64, error)
	getPurchaseOrderFn            func(context.Context, string) (*b2bv1.PurchaseOrder, error)
	createPurchaseOrderFn         func(context.Context, domain.PurchaseOrderCreateInput) (*b2bv1.PurchaseOrder, error)
	getDepartmentFn               func(context.Context, string) (*b2bv1.Department, error)
	getDepartmentNamesFn          func(context.Context, []string) (map[string]string, error)
	searchOrganisationsFn         func(context.Context, string, int) ([]*b2bv1.Organisation, error)
	getEmployeeFn                 func(context.Context, string) (*b2bv1.Employee, error)
	getEmployeeByUserIDFn         func(context.Context, string) (*b2bv1.Employee, error)
	getEmployeeByBusinessIDEmlFn  func(context.Context, string, string, string) (*b2bv1.Employee, error)
	createEmployeeFn              func(context.Context, domain.EmployeeCreateInput) (*b2bv1.Employee, error)
	updateEmployeeFn              func(context.Context, domain.EmployeeUpdateInput) (*b2bv1.Employee, error)
	createOrganisationFn          func(context.Context, domain.OrganisationCreateInput) (*b2bv1.Organisation, error)
	getOrganisationFn             func(context.Context, string) (*b2bv1.Organisation, error)
	getOrganisationByCodeFn       func(context.Context, string) (*b2bv1.Organisation, error)
	addOrgMemberFn                func(context.Context, domain.OrgMemberCreateInput) (*b2bv1.OrgMember, error)
	resolveOrganisationByUserIDFn func(context.Context, string) (string, b2bv1.OrgMemberRole, string, error)
}

func (f *fakeB2BRepo) CreateOrganisation(ctx context.Context, input domain.OrganisationCreateInput) (*b2bv1.Organisation, error) {
	return f.createOrganisationFn(ctx, input)
}
func (f *fakeB2BRepo) GetOrganisation(ctx context.Context, organisationID string) (*b2bv1.Organisation, error) {
	if f.getOrganisationFn != nil {
		return f.getOrganisationFn(ctx, organisationID)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) GetOrganisationByCode(ctx context.Context, code string) (*b2bv1.Organisation, error) {
	if f.getOrganisationByCodeFn != nil {
		return f.getOrganisationByCodeFn(ctx, code)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) SearchOrganisationsForEmployeeActivation(ctx context.Context, query string, limit int) ([]*b2bv1.Organisation, error) {
	if f.searchOrganisationsFn != nil {
		return f.searchOrganisationsFn(ctx, query, limit)
	}
	return []*b2bv1.Organisation{}, nil
}
func (f *fakeB2BRepo) ListOrganisations(context.Context, int, int, string, b2bv1.OrganisationStatus) ([]*b2bv1.Organisation, int64, error) {
	return []*b2bv1.Organisation{}, 0, nil
}
func (f *fakeB2BRepo) UpdateOrganisation(context.Context, domain.OrganisationUpdateInput) (*b2bv1.Organisation, error) {
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) DeleteOrganisation(context.Context, string) error {
	return nil
}
func (f *fakeB2BRepo) ListOrgMembers(context.Context, string) ([]*b2bv1.OrgMember, error) {
	return []*b2bv1.OrgMember{}, nil
}
func (f *fakeB2BRepo) AddOrgMember(ctx context.Context, input domain.OrgMemberCreateInput) (*b2bv1.OrgMember, error) {
	return f.addOrgMemberFn(ctx, input)
}
func (f *fakeB2BRepo) AssignOrgAdmin(context.Context, string, string) (*b2bv1.OrgMember, error) {
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) RemoveOrgMember(context.Context, string, string) error {
	return nil
}
func (f *fakeB2BRepo) ResolveOrganisationByUserID(ctx context.Context, userID string) (string, b2bv1.OrgMemberRole, string, error) {
	return f.resolveOrganisationByUserIDFn(ctx, userID)
}
func (f *fakeB2BRepo) ListDepartments(context.Context, int, int, string) ([]*b2bv1.Department, int64, error) {
	return []*b2bv1.Department{}, 0, nil
}
func (f *fakeB2BRepo) GetDepartment(ctx context.Context, departmentID string) (*b2bv1.Department, error) {
	return f.getDepartmentFn(ctx, departmentID)
}
func (f *fakeB2BRepo) CreateDepartment(context.Context, domain.DepartmentCreateInput) (*b2bv1.Department, error) {
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) UpdateDepartment(context.Context, domain.DepartmentUpdateInput) (*b2bv1.Department, error) {
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) UpdateDepartmentTotalPremium(context.Context, string) error { return nil }
func (f *fakeB2BRepo) DeleteDepartment(context.Context, string) error {
	return nil
}
func (f *fakeB2BRepo) ListEmployees(context.Context, int, int, string, string, b2bv1.EmployeeStatus) ([]*b2bv1.Employee, int64, error) {
	return []*b2bv1.Employee{}, 0, nil
}
func (f *fakeB2BRepo) GetEmployee(ctx context.Context, employeeUUID string) (*b2bv1.Employee, error) {
	if f.getEmployeeFn != nil {
		return f.getEmployeeFn(ctx, employeeUUID)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) GetEmployeeByUserID(ctx context.Context, userID string) (*b2bv1.Employee, error) {
	if f.getEmployeeByUserIDFn != nil {
		return f.getEmployeeByUserIDFn(ctx, userID)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) GetEmployeeByBusinessEmployeeIDEmail(ctx context.Context, businessID, employeeID, email string) (*b2bv1.Employee, error) {
	if f.getEmployeeByBusinessIDEmlFn != nil {
		return f.getEmployeeByBusinessIDEmlFn(ctx, businessID, employeeID, email)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) CreateEmployee(ctx context.Context, input domain.EmployeeCreateInput) (*b2bv1.Employee, error) {
	if f.createEmployeeFn != nil {
		return f.createEmployeeFn(ctx, input)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) UpdateEmployee(ctx context.Context, input domain.EmployeeUpdateInput) (*b2bv1.Employee, error) {
	if f.updateEmployeeFn != nil {
		return f.updateEmployeeFn(ctx, input)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeB2BRepo) DeleteEmployee(context.Context, string) error {
	return nil
}
func (f *fakeB2BRepo) GetDepartmentNames(ctx context.Context, departmentIDs []string) (map[string]string, error) {
	return f.getDepartmentNamesFn(ctx, departmentIDs)
}
func (f *fakeB2BRepo) ListCatalogPlans(ctx context.Context) ([]*domain.CatalogPlan, error) {
	return f.listCatalogPlansFn(ctx)
}
func (f *fakeB2BRepo) GetCatalogPlansByPlanIDs(ctx context.Context, planIDs []string) (map[string]*domain.CatalogPlan, error) {
	return f.getCatalogPlansByPlanIDsFn(ctx, planIDs)
}
func (f *fakeB2BRepo) ListPurchaseOrders(ctx context.Context, pageSize, offset int, businessID string, status b2bv1.PurchaseOrderStatus) ([]*b2bv1.PurchaseOrder, int64, error) {
	return f.listPurchaseOrdersFn(ctx, pageSize, offset, businessID, status)
}
func (f *fakeB2BRepo) GetPurchaseOrder(ctx context.Context, purchaseOrderID string) (*b2bv1.PurchaseOrder, error) {
	return f.getPurchaseOrderFn(ctx, purchaseOrderID)
}
func (f *fakeB2BRepo) CreatePurchaseOrder(ctx context.Context, input domain.PurchaseOrderCreateInput) (*b2bv1.PurchaseOrder, error) {
	return f.createPurchaseOrderFn(ctx, input)
}

type fakePublisher struct {
	created    []string
	memberAdds []string
	adminAdds  []string
}

type fakeEmployeeAuthNClient struct {
	provisionEmployeeUserFn       func(context.Context, *authnservicev1.ProvisionEmployeeUserRequest) (*authnservicev1.ProvisionEmployeeUserResponse, error)
	requestPasswordResetByEmailFn func(context.Context, *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error)
}

func (f *fakeEmployeeAuthNClient) ProvisionEmployeeUser(ctx context.Context, req *authnservicev1.ProvisionEmployeeUserRequest) (*authnservicev1.ProvisionEmployeeUserResponse, error) {
	if f.provisionEmployeeUserFn != nil {
		return f.provisionEmployeeUserFn(ctx, req)
	}
	return &authnservicev1.ProvisionEmployeeUserResponse{}, nil
}

func (f *fakeEmployeeAuthNClient) RequestPasswordResetByEmail(ctx context.Context, req *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error) {
	if f.requestPasswordResetByEmailFn != nil {
		return f.requestPasswordResetByEmailFn(ctx, req)
	}
	return &authnservicev1.RequestPasswordResetByEmailResponse{}, nil
}

type fakeEmployeeAuthZClient struct {
	assignRoleFn func(context.Context, *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error)
	listRolesFn  func(context.Context, *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error)
}

func (f *fakeEmployeeAuthZClient) AssignRole(ctx context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error) {
	if f.assignRoleFn != nil {
		return f.assignRoleFn(ctx, req)
	}
	return &authzservicev1.AssignRoleResponse{}, nil
}

func (f *fakeEmployeeAuthZClient) ListRoles(ctx context.Context, req *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error) {
	if f.listRolesFn != nil {
		return f.listRolesFn(ctx, req)
	}
	return &authzservicev1.ListRolesResponse{
		Roles: []*authzentityv1.Role{},
	}, nil
}

func (f *fakePublisher) PublishOrganisationCreated(context.Context, string, string, string, string, string) error {
	f.created = append(f.created, "created")
	return nil
}
func (f *fakePublisher) PublishOrganisationUpdated(context.Context, string, string, b2bv1.OrganisationStatus, string) error {
	return nil
}
func (f *fakePublisher) PublishOrganisationApproved(context.Context, string, string) error {
	return nil
}
func (f *fakePublisher) PublishOrgMemberAdded(context.Context, string, string, string, b2bv1.OrgMemberRole, string) error {
	f.memberAdds = append(f.memberAdds, "member")
	return nil
}
func (f *fakePublisher) PublishOrgMemberRemoved(context.Context, string, string, string, string) error {
	return nil
}
func (f *fakePublisher) PublishB2BAdminAssigned(context.Context, string, string, string) error {
	f.adminAdds = append(f.adminAdds, "admin")
	return nil
}

func TestB2BService_HelperFunctions(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-tenant-id", "tenant-meta"))
	t.Setenv("DEFAULT_TENANT_ID", "tenant-env")

	require.Equal(t, "tenant-explicit", resolveTenantID(ctx, "tenant-explicit"))
	require.Equal(t, "tenant-meta", resolveTenantID(ctx, ""))
	require.Equal(t, 0, parseOffset(""))
	require.Equal(t, 0, parseOffset("-1"))
	require.Equal(t, 42, parseOffset("42"))

	poNumber := makePurchaseOrderNumber(time.Date(2026, 3, 13, 0, 0, 0, 0, time.UTC))
	require.Contains(t, poNumber, "PO-20260313-")

	original := &commonv1.Money{Amount: 1000, Currency: "BDT", DecimalAmount: 10}
	cloned := cloneMoney(original)
	require.NotSame(t, original, cloned)
	require.Equal(t, original.Amount, cloned.Amount)

	multiplied := multiplyMoney(original, 3)
	require.Equal(t, int64(3000), multiplied.Amount)
	require.Equal(t, 30.0, multiplied.DecimalAmount)
	require.Equal(t, "Health Insurance", insuranceCategoryDisplayName(commonv1.InsuranceType_INSURANCE_TYPE_HEALTH))
	require.Equal(t, "Insurance Product", insuranceCategoryDisplayName(commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED))

	require.Nil(t, fallbackCatalogPlan("", "", commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED))
	require.NotNil(t, fallbackCatalogPlan(svcPlanHealth1, "", commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED))

	items := mergeCatalogWithSeedFallback([]*domain.CatalogPlan{{PlanID: "custom-plan", PlanName: "Custom"}})
	require.NotEmpty(t, items)
	require.GreaterOrEqual(t, len(items), len(seededCatalogPlans)+1)

	catalogMap := mergeCatalogMapWithSeedFallback(map[string]*domain.CatalogPlan{
		"custom-plan": {PlanID: "custom-plan", PlanName: "Custom"},
	})
	require.Contains(t, catalogMap, "custom-plan")
	require.Contains(t, catalogMap, svcPlanHealth1)

	view := purchaseOrderView(&b2bv1.PurchaseOrder{
		PurchaseOrderId: "po-1",
		DepartmentId:    "dept-1",
		ProductId:       "prod-1",
		PlanId:          svcPlanHealth1,
	}, map[string]string{}, map[string]*domain.CatalogPlan{})
	require.Equal(t, "Unassigned", view.DepartmentName)
	require.NotEmpty(t, view.ProductName)
	require.NotEmpty(t, view.PlanName)

	respItems := catalogItemsToResponse([]*domain.CatalogPlan{{PlanID: "plan-1", ProductID: "prod-1", ProductName: "Prod", PlanName: "Plan", PremiumAmount: original}})
	require.Len(t, respItems, 1)
	require.Equal(t, "plan-1", respItems[0].PlanId)
}

func TestB2BService_PurchaseOrderFlows(t *testing.T) {
	repo := &fakeB2BRepo{
		listCatalogPlansFn: func(context.Context) ([]*domain.CatalogPlan, error) {
			return []*domain.CatalogPlan{{PlanID: "custom-plan", ProductID: "prod-1", ProductName: "Custom Product", PlanName: "Custom Plan"}}, nil
		},
		listPurchaseOrdersFn: func(context.Context, int, int, string, b2bv1.PurchaseOrderStatus) ([]*b2bv1.PurchaseOrder, int64, error) {
			return []*b2bv1.PurchaseOrder{{
				PurchaseOrderId: "po-1",
				DepartmentId:    "dept-1",
				ProductId:       "prod-1",
				PlanId:          svcPlanHealth1,
			}}, 2, nil
		},
		getDepartmentNamesFn: func(context.Context, []string) (map[string]string, error) {
			return map[string]string{}, nil
		},
		getCatalogPlansByPlanIDsFn: func(context.Context, []string) (map[string]*domain.CatalogPlan, error) {
			return map[string]*domain.CatalogPlan{}, nil
		},
		getPurchaseOrderFn: func(context.Context, string) (*b2bv1.PurchaseOrder, error) {
			return nil, gorm.ErrRecordNotFound
		},
		getDepartmentFn: func(context.Context, string) (*b2bv1.Department, error) {
			return &b2bv1.Department{DepartmentId: "dept-1", BusinessId: "biz-1", Name: "Engineering"}, nil
		},
		createPurchaseOrderFn: func(_ context.Context, input domain.PurchaseOrderCreateInput) (*b2bv1.PurchaseOrder, error) {
			return &b2bv1.PurchaseOrder{
				PurchaseOrderId:     input.PurchaseOrderID,
				PurchaseOrderNumber: input.PurchaseOrderNumber,
				DepartmentId:        input.DepartmentID,
				ProductId:           input.ProductID,
				PlanId:              input.PlanID,
				InsuranceCategory:   input.InsuranceCategory,
				CoverageAmount:      input.CoverageAmount,
				EstimatedPremium:    input.EstimatedPremium,
				EmployeeCount:       input.EmployeeCount,
			}, nil
		},
	}
	svc := NewB2BService(repo, nil)

	catalogResp, err := svc.ListPurchaseOrderCatalog(context.Background(), nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(catalogResp.Items), len(seededCatalogPlans)+1)

	listResp, err := svc.ListPurchaseOrders(context.Background(), &b2bservicev1.ListPurchaseOrdersRequest{BusinessId: "biz-1", PageSize: 1})
	require.NoError(t, err)
	require.Len(t, listResp.PurchaseOrders, 1)
	require.Equal(t, "1", listResp.NextPageToken)
	require.Equal(t, "Unassigned", listResp.PurchaseOrders[0].DepartmentName)

	_, err = svc.GetPurchaseOrder(context.Background(), &b2bservicev1.GetPurchaseOrderRequest{PurchaseOrderId: "missing"})
	require.ErrorIs(t, err, ErrNotFound)

	_, err = svc.CreatePurchaseOrder(context.Background(), nil)
	require.ErrorIs(t, err, ErrInvalidArgument)

	createResp, err := svc.CreatePurchaseOrder(context.Background(), &b2bservicev1.CreatePurchaseOrderRequest{
		DepartmentId:   "dept-1",
		PlanId:         svcPlanHealth1,
		EmployeeCount:  2,
		CoverageAmount: &commonv1.Money{Amount: 10000, Currency: "BDT", DecimalAmount: 100},
		RequestedBy:    uuid.NewString(),
	})
	require.NoError(t, err)
	require.Equal(t, "Engineering", createResp.PurchaseOrder.DepartmentName)
	require.Equal(t, svcPlanHealth1, createResp.PurchaseOrder.PurchaseOrder.PlanId)
	require.NotNil(t, createResp.PurchaseOrder.PurchaseOrder.EstimatedPremium)
}

func TestB2BService_OrganisationMemberFlows(t *testing.T) {
	pub := &fakePublisher{}
	repo := &fakeB2BRepo{
		createOrganisationFn: func(_ context.Context, input domain.OrganisationCreateInput) (*b2bv1.Organisation, error) {
			return &b2bv1.Organisation{
				OrganisationId: input.OrganisationID,
				TenantId:       input.TenantID,
				Name:           input.Name,
				Code:           input.Code,
			}, nil
		},
		addOrgMemberFn: func(_ context.Context, input domain.OrgMemberCreateInput) (*b2bv1.OrgMember, error) {
			return &b2bv1.OrgMember{
				MemberId:       input.MemberID,
				OrganisationId: input.OrganisationID,
				UserId:         input.UserID,
				Role:           input.Role,
			}, nil
		},
		resolveOrganisationByUserIDFn: func(context.Context, string) (string, b2bv1.OrgMemberRole, string, error) {
			return "org-1", b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN, "Acme", nil
		},
	}
	svc := NewB2BService(repo, pub)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-tenant-id", "tenant-meta", "x-user-id", "caller-1"))

	orgResp, err := svc.CreateOrganisation(ctx, &b2bservicev1.CreateOrganisationRequest{
		Name: "Acme",
		Code: "acme",
	})
	require.NoError(t, err)
	require.Equal(t, "tenant-meta", orgResp.Organisation.TenantId)
	require.Equal(t, "ACME", orgResp.Organisation.Code)
	require.Len(t, pub.created, 1)

	_, err = svc.CreateOrganisation(context.Background(), &b2bservicev1.CreateOrganisationRequest{})
	require.ErrorIs(t, err, ErrInvalidArgument)

	memberResp, err := svc.AddOrgMember(ctx, &b2bservicev1.AddOrgMemberRequest{
		OrganisationId: "org-1",
		UserId:         "user-1",
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN,
	})
	require.NoError(t, err)
	require.Equal(t, "user-1", memberResp.Member.UserId)
	require.Len(t, pub.memberAdds, 1)
	require.Len(t, pub.adminAdds, 1)

	resolveResp, err := svc.ResolveMyOrganisation(context.Background(), &b2bservicev1.ResolveMyOrganisationRequest{UserId: "user-1"})
	require.NoError(t, err)
	require.Equal(t, "org-1", resolveResp.OrganisationId)

	repo.resolveOrganisationByUserIDFn = func(context.Context, string) (string, b2bv1.OrgMemberRole, string, error) {
		return "", 0, "", gorm.ErrRecordNotFound
	}
	_, err = svc.ResolveMyOrganisation(context.Background(), &b2bservicev1.ResolveMyOrganisationRequest{UserId: "missing"})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestB2BService_CreateEmployeeNormalizesMobile(t *testing.T) {
	repo := &fakeB2BRepo{
		getCatalogPlansByPlanIDsFn: func(context.Context, []string) (map[string]*domain.CatalogPlan, error) {
			return map[string]*domain.CatalogPlan{}, nil
		},
		createEmployeeFn: func(_ context.Context, input domain.EmployeeCreateInput) (*b2bv1.Employee, error) {
			require.Equal(t, "+8801712345678", input.MobileNumber)
			return &b2bv1.Employee{
				EmployeeUuid:   input.EmployeeUUID,
				DepartmentId:   input.DepartmentID,
				AssignedPlanId: input.AssignedPlanID,
				MobileNumber:   input.MobileNumber,
			}, nil
		},
		getDepartmentNamesFn: func(context.Context, []string) (map[string]string, error) {
			return map[string]string{"dept-1": "Admin"}, nil
		},
	}
	svc := NewB2BService(repo, nil)

	resp, err := svc.CreateEmployee(context.Background(), &b2bservicev1.CreateEmployeeRequest{
		Name:          "Admin",
		EmployeeId:    "EMP-1",
		DepartmentId:  "dept-1",
		BusinessId:    "org-1",
		Email:         "admin@example.com",
		MobileNumber:  "01712345678",
		DateOfJoining: "2026-03-21",
	})
	require.NoError(t, err)
	require.Equal(t, "+8801712345678", resp.Employee.Employee.MobileNumber)
}

func TestB2BService_CreateEmployeeRejectsInvalidMobile(t *testing.T) {
	svc := NewB2BService(&fakeB2BRepo{}, nil)

	_, err := svc.CreateEmployee(context.Background(), &b2bservicev1.CreateEmployeeRequest{
		Name:          "Admin",
		EmployeeId:    "EMP-1",
		DepartmentId:  "dept-1",
		BusinessId:    "org-1",
		Email:         "admin@example.com",
		MobileNumber:  "01112345678",
		DateOfJoining: "2026-03-21",
	})
	require.ErrorIs(t, err, ErrInvalidArgument)
	require.Contains(t, err.Error(), "mobile_number must be a valid Bangladesh number")
}

func TestB2BService_ActivateEmployeeUsesOrganisationCode(t *testing.T) {
	const (
		orgCode = "LPL"
		orgID   = "org-1"
		userID  = "user-1"
	)

	authn := &fakeEmployeeAuthNClient{
		requestPasswordResetByEmailFn: func(_ context.Context, req *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error) {
			require.Equal(t, "employee@example.com", req.GetEmail())
			return &authnservicev1.RequestPasswordResetByEmailResponse{
				OtpId:            "otp-1",
				Message:          "Verification code sent",
				ExpiresInSeconds: 600,
			}, nil
		},
	}
	authz := &fakeEmployeeAuthZClient{
		listRolesFn: func(context.Context, *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error) {
			return &authzservicev1.ListRolesResponse{
				Roles: []*authzentityv1.Role{{RoleId: "role-1", Name: "b2b_beneficiary"}},
			}, nil
		},
		assignRoleFn: func(_ context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error) {
			require.Equal(t, userID, req.GetUserId())
			require.Equal(t, "b2b:"+orgID, req.GetDomain())
			return &authzservicev1.AssignRoleResponse{}, nil
		},
	}

	repo := &fakeB2BRepo{
		getOrganisationByCodeFn: func(_ context.Context, code string) (*b2bv1.Organisation, error) {
			require.Equal(t, orgCode, code)
			return &b2bv1.Organisation{OrganisationId: orgID, Code: code}, nil
		},
		getEmployeeByBusinessIDEmlFn: func(_ context.Context, businessID, employeeID, email string) (*b2bv1.Employee, error) {
			require.Equal(t, orgID, businessID)
			require.Equal(t, "EMP-1", employeeID)
			require.Equal(t, "employee@example.com", email)
			return &b2bv1.Employee{
				EmployeeUuid: "emp-uuid",
				EmployeeId:   employeeID,
				BusinessId:   businessID,
				Email:        email,
				UserId:       userID,
			}, nil
		},
	}

	svc := NewB2BService(repo, nil).WithEmployeeIdentity(authn, authz)

	resp, err := svc.ActivateEmployee(context.Background(), &b2bservicev1.ActivateEmployeeRequest{
		OrganisationCode: strings.ToLower(orgCode),
		EmployeeId:       "EMP-1",
		Email:            "employee@example.com",
	})
	require.NoError(t, err)
	require.Equal(t, userID, resp.GetUserId())
	require.Equal(t, "otp-1", resp.GetOtpId())
}

func TestB2BService_ActivateEmployeeFallsBackToOrganisationID(t *testing.T) {
	const (
		orgID  = "org-1"
		userID = "user-1"
	)

	authn := &fakeEmployeeAuthNClient{
		requestPasswordResetByEmailFn: func(_ context.Context, req *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error) {
			return &authnservicev1.RequestPasswordResetByEmailResponse{
				OtpId:            "otp-1",
				Message:          "Verification code sent",
				ExpiresInSeconds: 600,
			}, nil
		},
	}
	authz := &fakeEmployeeAuthZClient{
		listRolesFn: func(context.Context, *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error) {
			return &authzservicev1.ListRolesResponse{
				Roles: []*authzentityv1.Role{{RoleId: "role-1", Name: "b2b_beneficiary"}},
			}, nil
		},
		assignRoleFn: func(_ context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error) {
			require.Equal(t, "b2b:"+orgID, req.GetDomain())
			return &authzservicev1.AssignRoleResponse{}, nil
		},
	}

	repo := &fakeB2BRepo{
		getOrganisationByCodeFn: func(_ context.Context, code string) (*b2bv1.Organisation, error) {
			require.Empty(t, code)
			return nil, gorm.ErrRecordNotFound
		},
		getEmployeeByBusinessIDEmlFn: func(_ context.Context, businessID, employeeID, email string) (*b2bv1.Employee, error) {
			require.Equal(t, orgID, businessID)
			return &b2bv1.Employee{
				EmployeeUuid: "emp-uuid",
				EmployeeId:   employeeID,
				BusinessId:   businessID,
				Email:        email,
				UserId:       userID,
			}, nil
		},
	}
	repo.getOrganisationByCodeFn = func(_ context.Context, code string) (*b2bv1.Organisation, error) {
		require.Empty(t, code)
		return nil, gorm.ErrRecordNotFound
	}
	repo.getOrganisationFn = func(_ context.Context, requestedID string) (*b2bv1.Organisation, error) {
		require.Equal(t, orgID, requestedID)
		return &b2bv1.Organisation{OrganisationId: orgID}, nil
	}

	svc := NewB2BService(repo, nil).WithEmployeeIdentity(authn, authz)

	resp, err := svc.ActivateEmployee(context.Background(), &b2bservicev1.ActivateEmployeeRequest{
		OrganisationId: orgID,
		EmployeeId:     "EMP-1",
		Email:          "employee@example.com",
	})
	require.NoError(t, err)
	require.Equal(t, userID, resp.GetUserId())
}

func TestB2BService_ListEmployeeLoginOrganisationsFiltersByQuery(t *testing.T) {
	repo := &fakeB2BRepo{
		searchOrganisationsFn: func(_ context.Context, query string, limit int) ([]*b2bv1.Organisation, error) {
			require.Equal(t, "Alpha Security", query)
			require.Equal(t, 5, limit)
			return []*b2bv1.Organisation{
				{
					OrganisationId: "org-1",
					Name:           "Alpha Security Ltd",
					Code:           "ASL",
				},
			}, nil
		},
	}

	svc := NewB2BService(repo, nil)
	resp, err := svc.ListEmployeeLoginOrganisations(context.Background(), &b2bservicev1.ListEmployeeLoginOrganisationsRequest{
		Query:    "Alpha Security",
		PageSize: 5,
	})
	require.NoError(t, err)
	require.Len(t, resp.GetOrganisations(), 1)
	require.Equal(t, "Alpha Security Ltd", resp.GetOrganisations()[0].GetOrganisationName())
	require.Equal(t, "ASL", resp.GetOrganisations()[0].GetOrganisationCode())
}
