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

// AssignOrgAdmin Assign a platform user as an OrgAdmin
func (s *B2bService) AssignOrgAdmin(ctx context.Context, organisationId string, req *models.OrgAdminAssignmentRequest) (*models.OrgAdminAssignmentResponse, error) {
	path := "/v1/b2b/organisations/{organisation_id}/admins"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	var result models.OrgAdminAssignmentResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListPurchaseOrders List purchase orders for the authenticated organisation
func (s *B2bService) ListPurchaseOrders(ctx context.Context) (*models.PurchaseOrdersListingResponse, error) {
	path := "/v1/b2b/purchase-orders"
	var result models.PurchaseOrdersListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreatePurchaseOrder Create a purchase order for a product plan
func (s *B2bService) CreatePurchaseOrder(ctx context.Context, req *models.PurchaseOrderCreationRequest) (*models.PurchaseOrderCreationResponse, error) {
	path := "/v1/b2b/purchase-orders"
	var result models.PurchaseOrderCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPurchaseOrder Get a single purchase order
func (s *B2bService) GetPurchaseOrder(ctx context.Context, purchaseOrderId string) (*models.PurchaseOrderRetrievalResponse, error) {
	path := "/v1/b2b/purchase-orders/{purchase_order_id}"
	path = strings.ReplaceAll(path, "{purchase_order_id}", purchaseOrderId)
	var result models.PurchaseOrderRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateOrganisation Create a new organisation (SuperAdmin only)
func (s *B2bService) CreateOrganisation(ctx context.Context, req *models.OrganisationCreationRequest) (*models.OrganisationCreationResponse, error) {
	path := "/v1/b2b/organisations"
	var result models.OrganisationCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListOrganisations List all organisations (SuperAdmin: all; BizAdmin: own only)
func (s *B2bService) ListOrganisations(ctx context.Context) (*models.OrganisationsListingResponse, error) {
	path := "/v1/b2b/organisations"
	var result models.OrganisationsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListEmployees List employees for the authenticated organisation
func (s *B2bService) ListEmployees(ctx context.Context) (*models.EmployeesListingResponse, error) {
	path := "/v1/b2b/employees"
	var result models.EmployeesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateEmployee Create a new employee
func (s *B2bService) CreateEmployee(ctx context.Context, req *models.EmployeeCreationRequest) (*models.EmployeeCreationResponse, error) {
	path := "/v1/b2b/employees"
	var result models.EmployeeCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDepartment Get a single department
func (s *B2bService) GetDepartment(ctx context.Context, departmentId string) (*models.DepartmentRetrievalResponse, error) {
	path := "/v1/b2b/departments/{department_id}"
	path = strings.ReplaceAll(path, "{department_id}", departmentId)
	var result models.DepartmentRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateDepartment Update a department's name
func (s *B2bService) UpdateDepartment(ctx context.Context, departmentId string, req *models.DepartmentUpdateRequest) (*models.DepartmentUpdateResponse, error) {
	path := "/v1/b2b/departments/{department_id}"
	path = strings.ReplaceAll(path, "{department_id}", departmentId)
	var result models.DepartmentUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDepartment Soft-delete a department (only if no active employees)
func (s *B2bService) DeleteDepartment(ctx context.Context, departmentId string) error {
	path := "/v1/b2b/departments/{department_id}"
	path = strings.ReplaceAll(path, "{department_id}", departmentId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// GetEmployee Get a single employee by employee_uuid
func (s *B2bService) GetEmployee(ctx context.Context, employeeUuid string) (*models.EmployeeRetrievalResponse, error) {
	path := "/v1/b2b/employees/{employee_uuid}"
	path = strings.ReplaceAll(path, "{employee_uuid}", employeeUuid)
	var result models.EmployeeRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateEmployee Update an existing employee's details
func (s *B2bService) UpdateEmployee(ctx context.Context, employeeUuid string, req *models.EmployeeUpdateRequest) (*models.EmployeeUpdateResponse, error) {
	path := "/v1/b2b/employees/{employee_uuid}"
	path = strings.ReplaceAll(path, "{employee_uuid}", employeeUuid)
	var result models.EmployeeUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteEmployee Soft-delete an employee record
func (s *B2bService) DeleteEmployee(ctx context.Context, employeeUuid string) error {
	path := "/v1/b2b/employees/{employee_uuid}"
	path = strings.ReplaceAll(path, "{employee_uuid}", employeeUuid)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// CreateDepartment Create a new department
func (s *B2bService) CreateDepartment(ctx context.Context, req *models.DepartmentCreationRequest) (*models.DepartmentCreationResponse, error) {
	path := "/v1/b2b/departments"
	var result models.DepartmentCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDepartments List departments for the authenticated organisation
func (s *B2bService) ListDepartments(ctx context.Context) (*models.DepartmentsListingResponse, error) {
	path := "/v1/b2b/departments"
	var result models.DepartmentsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOrganisation Get a single organisation by ID
func (s *B2bService) GetOrganisation(ctx context.Context, organisationId string) (*models.OrganisationRetrievalResponse, error) {
	path := "/v1/b2b/organisations/{organisation_id}"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	var result models.OrganisationRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateOrganisation Update an organisation's profile
func (s *B2bService) UpdateOrganisation(ctx context.Context, organisationId string, req *models.OrganisationUpdateRequest) (*models.OrganisationUpdateResponse, error) {
	path := "/v1/b2b/organisations/{organisation_id}"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	var result models.OrganisationUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteOrganisation Soft-delete an organisation and revoke its memberships
func (s *B2bService) DeleteOrganisation(ctx context.Context, organisationId string) error {
	path := "/v1/b2b/organisations/{organisation_id}"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// ListPurchaseOrderCatalog List purchasable product plans for purchase orders
func (s *B2bService) ListPurchaseOrderCatalog(ctx context.Context) (*models.PurchaseOrderCatalogListingResponse, error) {
	path := "/v1/b2b/purchase-orders/catalog"
	var result models.PurchaseOrderCatalogListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveOrgMember Remove an OrgMember from the organisation
func (s *B2bService) RemoveOrgMember(ctx context.Context, organisationId string, memberId string) error {
	path := "/v1/b2b/organisations/{organisation_id}/members/{member_id}"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	path = strings.ReplaceAll(path, "{member_id}", memberId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// ListOrgMembers List members for an organisation
func (s *B2bService) ListOrgMembers(ctx context.Context, organisationId string) (*models.OrgMembersListingResponse, error) {
	path := "/v1/b2b/organisations/{organisation_id}/members"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	var result models.OrgMembersListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddOrgMember Add a platform user as an OrgMember
func (s *B2bService) AddOrgMember(ctx context.Context, organisationId string, req *models.AddOrgMemberRequest) (*models.AddOrgMemberResponse, error) {
	path := "/v1/b2b/organisations/{organisation_id}/members"
	path = strings.ReplaceAll(path, "{organisation_id}", organisationId)
	var result models.AddOrgMemberResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

