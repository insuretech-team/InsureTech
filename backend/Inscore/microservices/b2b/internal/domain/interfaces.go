package domain

import (
	"context"

	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

// ─── CATALOG ─────────────────────────────────────────────────────────────────

type CatalogPlan struct {
	ProductID         string
	ProductName       string
	PlanID            string
	PlanName          string
	InsuranceCategory commonv1.InsuranceType
	PremiumAmount     *commonv1.Money
}

// ─── INPUT TYPES ─────────────────────────────────────────────────────────────

type OrganisationCreateInput struct {
	OrganisationID string
	TenantID       string
	Name           string
	Code           string
	Industry       string
	ContactEmail   string
	ContactPhone   string
	Address        string
}

type OrganisationUpdateInput struct {
	OrganisationID string
	Name           string
	Industry       string
	ContactEmail   string
	ContactPhone   string
	Address        string
	Status         b2bv1.OrganisationStatus
}

type OrgMemberCreateInput struct {
	MemberID       string
	OrganisationID string
	UserID         string
	Role           b2bv1.OrgMemberRole
}

type DepartmentCreateInput struct {
	DepartmentID string
	Name         string
	BusinessID   string
}

type DepartmentUpdateInput struct {
	DepartmentID string
	Name         string
}

type EmployeeCreateInput struct {
	EmployeeUUID      string
	Name              string
	EmployeeID        string
	DepartmentID      string
	BusinessID        string
	InsuranceCategory commonv1.InsuranceType
	AssignedPlanID    string
	CoverageAmount    *commonv1.Money
	// PremiumAmount is resolved from the assigned plan's catalog price and stored
	// on the employee row so it is readable without joining the plan catalog.
	PremiumAmount     *commonv1.Money
	NumberOfDependent int32
	// PII fields
	Email         string
	MobileNumber  string
	DateOfBirth   string
	DateOfJoining string
	Gender        b2bv1.EmployeeGender
	// B2C portal bridge — set when portal access is granted
	UserID string
}

type EmployeeUpdateInput struct {
	EmployeeUUID      string
	Name              string
	DepartmentID      string
	Email             string
	MobileNumber      string
	DateOfBirth       string
	DateOfJoining     string
	Gender            b2bv1.EmployeeGender
	InsuranceCategory commonv1.InsuranceType
	AssignedPlanID    string
	CoverageAmount    *commonv1.Money
	NumberOfDependent int32
	Status            b2bv1.EmployeeStatus
}

type PurchaseOrderCreateInput struct {
	PurchaseOrderID     string
	PurchaseOrderNumber string
	BusinessID          string
	DepartmentID        string
	ProductID           string
	PlanID              string
	InsuranceCategory   commonv1.InsuranceType
	EmployeeCount       int32
	NumberOfDependents  int32
	CoverageAmount      *commonv1.Money
	EstimatedPremium    *commonv1.Money
	Status              b2bv1.PurchaseOrderStatus
	RequestedBy         string
	Notes               string
}

// ─── REPOSITORY INTERFACE ────────────────────────────────────────────────────

type B2BRepository interface {
	// Organisation
	CreateOrganisation(ctx context.Context, input OrganisationCreateInput) (*b2bv1.Organisation, error)
	GetOrganisation(ctx context.Context, organisationID string) (*b2bv1.Organisation, error)
	ListOrganisations(ctx context.Context, pageSize, offset int, tenantID string, status b2bv1.OrganisationStatus) ([]*b2bv1.Organisation, int64, error)
	UpdateOrganisation(ctx context.Context, input OrganisationUpdateInput) (*b2bv1.Organisation, error)
	DeleteOrganisation(ctx context.Context, organisationID string) error

	// OrgMember
	ListOrgMembers(ctx context.Context, organisationID string) ([]*b2bv1.OrgMember, error)
	AddOrgMember(ctx context.Context, input OrgMemberCreateInput) (*b2bv1.OrgMember, error)
	AssignOrgAdmin(ctx context.Context, organisationID, memberID string) (*b2bv1.OrgMember, error)
	RemoveOrgMember(ctx context.Context, organisationID, memberID string) error
	// ResolveOrganisationByUserID returns the organisation_id + role for a given user.
	// Used by the REST gateway / service layer to inject business_id into every request.
	ResolveOrganisationByUserID(ctx context.Context, userID string) (organisationID string, role b2bv1.OrgMemberRole, organisationName string, err error)

	// Department
	ListDepartments(ctx context.Context, pageSize, offset int, businessID string) ([]*b2bv1.Department, int64, error)
	GetDepartment(ctx context.Context, departmentID string) (*b2bv1.Department, error)
	CreateDepartment(ctx context.Context, input DepartmentCreateInput) (*b2bv1.Department, error)
	UpdateDepartment(ctx context.Context, input DepartmentUpdateInput) (*b2bv1.Department, error)
	UpdateDepartmentTotalPremium(ctx context.Context, departmentID string) error
	DeleteDepartment(ctx context.Context, departmentID string) error

	// Employee
	ListEmployees(ctx context.Context, pageSize, offset int, departmentID, businessID string, status b2bv1.EmployeeStatus) ([]*b2bv1.Employee, int64, error)
	GetEmployee(ctx context.Context, employeeUUID string) (*b2bv1.Employee, error)
	CreateEmployee(ctx context.Context, input EmployeeCreateInput) (*b2bv1.Employee, error)
	UpdateEmployee(ctx context.Context, input EmployeeUpdateInput) (*b2bv1.Employee, error)
	DeleteEmployee(ctx context.Context, employeeUUID string) error

	// Enrichment helpers
	GetDepartmentNames(ctx context.Context, departmentIDs []string) (map[string]string, error)

	// Catalog
	ListCatalogPlans(ctx context.Context) ([]*CatalogPlan, error)
	GetCatalogPlansByPlanIDs(ctx context.Context, planIDs []string) (map[string]*CatalogPlan, error)

	// Purchase Orders
	ListPurchaseOrders(ctx context.Context, pageSize, offset int, businessID string, status b2bv1.PurchaseOrderStatus) ([]*b2bv1.PurchaseOrder, int64, error)
	GetPurchaseOrder(ctx context.Context, purchaseOrderID string) (*b2bv1.PurchaseOrder, error)
	CreatePurchaseOrder(ctx context.Context, input PurchaseOrderCreateInput) (*b2bv1.PurchaseOrder, error)
}

// ─── SERVICE INTERFACE ───────────────────────────────────────────────────────

type B2BService interface {
	// Organisation
	CreateOrganisation(ctx context.Context, req *b2bservicev1.CreateOrganisationRequest) (*b2bservicev1.CreateOrganisationResponse, error)
	GetOrganisation(ctx context.Context, req *b2bservicev1.GetOrganisationRequest) (*b2bservicev1.GetOrganisationResponse, error)
	ListOrganisations(ctx context.Context, req *b2bservicev1.ListOrganisationsRequest) (*b2bservicev1.ListOrganisationsResponse, error)
	UpdateOrganisation(ctx context.Context, req *b2bservicev1.UpdateOrganisationRequest) (*b2bservicev1.UpdateOrganisationResponse, error)
	DeleteOrganisation(ctx context.Context, req *b2bservicev1.DeleteOrganisationRequest) (*b2bservicev1.DeleteOrganisationResponse, error)
	ListOrgMembers(ctx context.Context, req *b2bservicev1.ListOrgMembersRequest) (*b2bservicev1.ListOrgMembersResponse, error)
	AddOrgMember(ctx context.Context, req *b2bservicev1.AddOrgMemberRequest) (*b2bservicev1.AddOrgMemberResponse, error)
	AssignOrgAdmin(ctx context.Context, req *b2bservicev1.AssignOrgAdminRequest) (*b2bservicev1.AssignOrgAdminResponse, error)
	RemoveOrgMember(ctx context.Context, req *b2bservicev1.RemoveOrgMemberRequest) (*b2bservicev1.RemoveOrgMemberResponse, error)
	ResolveMyOrganisation(ctx context.Context, req *b2bservicev1.ResolveMyOrganisationRequest) (*b2bservicev1.ResolveMyOrganisationResponse, error)

	// Department
	ListDepartments(ctx context.Context, req *b2bservicev1.ListDepartmentsRequest) (*b2bservicev1.ListDepartmentsResponse, error)
	GetDepartment(ctx context.Context, req *b2bservicev1.GetDepartmentRequest) (*b2bservicev1.GetDepartmentResponse, error)
	CreateDepartment(ctx context.Context, req *b2bservicev1.CreateDepartmentRequest) (*b2bservicev1.CreateDepartmentResponse, error)
	UpdateDepartment(ctx context.Context, req *b2bservicev1.UpdateDepartmentRequest) (*b2bservicev1.UpdateDepartmentResponse, error)
	DeleteDepartment(ctx context.Context, req *b2bservicev1.DeleteDepartmentRequest) (*b2bservicev1.DeleteDepartmentResponse, error)

	// Employee
	ListEmployees(ctx context.Context, req *b2bservicev1.ListEmployeesRequest) (*b2bservicev1.ListEmployeesResponse, error)
	GetEmployee(ctx context.Context, req *b2bservicev1.GetEmployeeRequest) (*b2bservicev1.GetEmployeeResponse, error)
	CreateEmployee(ctx context.Context, req *b2bservicev1.CreateEmployeeRequest) (*b2bservicev1.CreateEmployeeResponse, error)
	UpdateEmployee(ctx context.Context, req *b2bservicev1.UpdateEmployeeRequest) (*b2bservicev1.UpdateEmployeeResponse, error)
	DeleteEmployee(ctx context.Context, req *b2bservicev1.DeleteEmployeeRequest) (*b2bservicev1.DeleteEmployeeResponse, error)

	// Purchase Orders
	ListPurchaseOrderCatalog(ctx context.Context, req *b2bservicev1.ListPurchaseOrderCatalogRequest) (*b2bservicev1.ListPurchaseOrderCatalogResponse, error)
	ListPurchaseOrders(ctx context.Context, req *b2bservicev1.ListPurchaseOrdersRequest) (*b2bservicev1.ListPurchaseOrdersResponse, error)
	GetPurchaseOrder(ctx context.Context, req *b2bservicev1.GetPurchaseOrderRequest) (*b2bservicev1.GetPurchaseOrderResponse, error)
	CreatePurchaseOrder(ctx context.Context, req *b2bservicev1.CreatePurchaseOrderRequest) (*b2bservicev1.CreatePurchaseOrderResponse, error)
}
