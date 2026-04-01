package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// B2bService handles b2b-related API calls
type B2bService struct {
	Client Client
}

// GetMyEmployeeCoverage Resolve the authenticated employee's assigned plan and coverage
func (s *B2bService) GetMyEmployeeCoverage(ctx context.Context) error {
	path := "/v1/b2b-self/coverage"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetMyEmployeeProfile Resolve the authenticated employee's own profile
func (s *B2bService) GetMyEmployeeProfile(ctx context.Context) error {
	path := "/v1/b2b-self/profile"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ListDepartments List departments for the authenticated organisation
func (s *B2bService) ListDepartments(ctx context.Context) error {
	path := "/v1/b2b/departments"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateDepartment Create a new department
func (s *B2bService) CreateDepartment(ctx context.Context, req *models.DepartmentCreationRequest) error {
	path := "/v1/b2b/departments"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetDepartment Get a single department
func (s *B2bService) GetDepartment(ctx context.Context, departmentId string) error {
	path := "/v1/b2b/departments/{department_id}"
	path = strings.ReplaceAll(path, "{department_id}", departmentId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateDepartment Update a department's name
func (s *B2bService) UpdateDepartment(ctx context.Context, departmentId string, req *models.DepartmentUpdateRequest) error {
	path := "/v1/b2b/departments/{department_id}"
	path = strings.ReplaceAll(path, "{department_id}", departmentId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeleteDepartment Soft-delete a department (only if no active employees)
func (s *B2bService) DeleteDepartment(ctx context.Context, departmentId string) error {
	path := "/v1/b2b/departments/{department_id}"
	path = strings.ReplaceAll(path, "{department_id}", departmentId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// ListEmployees List employees for the authenticated organisation
func (s *B2bService) ListEmployees(ctx context.Context) error {
	path := "/v1/b2b/employees"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateEmployee Create a new employee
func (s *B2bService) CreateEmployee(ctx context.Context, req *models.EmployeeCreationRequest) error {
	path := "/v1/b2b/employees"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetEmployee Get a single employee by employee_uuid
func (s *B2bService) GetEmployee(ctx context.Context, employeeUuid string) error {
	path := "/v1/b2b/employees/{employee_uuid}"
	path = strings.ReplaceAll(path, "{employee_uuid}", employeeUuid)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateEmployee Update an existing employee's details
func (s *B2bService) UpdateEmployee(ctx context.Context, employeeUuid string, req *models.EmployeeUpdateRequest) error {
	path := "/v1/b2b/employees/{employee_uuid}"
	path = strings.ReplaceAll(path, "{employee_uuid}", employeeUuid)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeleteEmployee Soft-delete an employee record
func (s *B2bService) DeleteEmployee(ctx context.Context, employeeUuid string) error {
	path := "/v1/b2b/employees/{employee_uuid}"
	path = strings.ReplaceAll(path, "{employee_uuid}", employeeUuid)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// ActivateEmployee Start employee self-service activation using organisation code +
func (s *B2bService) ActivateEmployee(ctx context.Context, req *models.EmployeeActivationRequest) error {
	path := "/v1/b2b/employees:activate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListOrganisations List all organisations (SuperAdmin: all; BizAdmin: own only)
func (s *B2bService) ListOrganisations(ctx context.Context) error {
	path := "/v1/b2b/organisations"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateOrganisation Create a new organisation (SuperAdmin only)
func (s *B2bService) CreateOrganisation(ctx context.Context, req *models.OrganisationCreationRequest) error {
	path := "/v1/b2b/organisations"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetOrganisation Get a single organisation by ID
func (s *B2bService) GetOrganisation(ctx context.Context, organisationId string) error {
	path := "/v1/b2b/organisations/{organisation_id}"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateOrganisation Update an organisation's profile
func (s *B2bService) UpdateOrganisation(ctx context.Context, organisationId string, req *models.OrganisationUpdateRequest) error {
	path := "/v1/b2b/organisations/{organisation_id}"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeleteOrganisation Soft-delete an organisation and revoke its memberships
func (s *B2bService) DeleteOrganisation(ctx context.Context, organisationId string) error {
	path := "/v1/b2b/organisations/{organisation_id}"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// AssignOrgAdmin Assign a platform user as an OrgAdmin
func (s *B2bService) AssignOrgAdmin(ctx context.Context, organisationId string, req *models.OrgAdminAssignmentRequest) error {
	path := "/v1/b2b/organisations/{organisation_id}/admins:assign"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListOrgMembers List members for an organisation
func (s *B2bService) ListOrgMembers(ctx context.Context, organisationId string) error {
	path := "/v1/b2b/organisations/{organisation_id}/members"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// AddOrgMember Add a platform user as an OrgMember
func (s *B2bService) AddOrgMember(ctx context.Context, organisationId string, req *models.AddOrgMemberRequest) error {
	path := "/v1/b2b/organisations/{organisation_id}/members"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RemoveOrgMember Remove an OrgMember from the organisation
func (s *B2bService) RemoveOrgMember(ctx context.Context, organisationId string, memberId string) error {
	path := "/v1/b2b/organisations/{organisation_id}/members/{member_id}"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	path = strings.ReplaceAll(path, "{member_id}", memberId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// ListEmployeeLoginOrganisations List organisations matching a partial name/code for employee activation
func (s *B2bService) ListEmployeeLoginOrganisations(ctx context.Context) error {
	path := "/v1/b2b/organisations:employee-login"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ListPurchaseOrders List purchase orders for the authenticated organisation
func (s *B2bService) ListPurchaseOrders(ctx context.Context) error {
	path := "/v1/b2b/purchase-orders"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreatePurchaseOrder Create a purchase order for a product plan
func (s *B2bService) CreatePurchaseOrder(ctx context.Context, req *models.PurchaseOrderCreationRequest) error {
	path := "/v1/b2b/purchase-orders"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListPurchaseOrderCatalog List purchasable product plans for purchase orders
func (s *B2bService) ListPurchaseOrderCatalog(ctx context.Context) error {
	path := "/v1/b2b/purchase-orders/catalog"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetPurchaseOrder Get a single purchase order
func (s *B2bService) GetPurchaseOrder(ctx context.Context, purchaseOrderId string) error {
	path := "/v1/b2b/purchase-orders/{purchase_order_id}"
	path = strings.ReplaceAll(path, "{purchase_order_id}", purchaseOrderId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

