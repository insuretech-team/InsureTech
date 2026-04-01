package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/mobile"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type B2BService struct {
	repo      domain.B2BRepository
	publisher EventPublisher
	authn     EmployeeAuthNClient
	authz     EmployeeAuthZClient
}

type EmployeeAuthNClient interface {
	ProvisionEmployeeUser(ctx context.Context, req *authnservicev1.ProvisionEmployeeUserRequest) (*authnservicev1.ProvisionEmployeeUserResponse, error)
	RequestPasswordResetByEmail(ctx context.Context, req *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error)
}

type EmployeeAuthZClient interface {
	AssignRole(ctx context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error)
	ListRoles(ctx context.Context, req *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error)
}

// EventPublisher interface for publishing B2B events
type EventPublisher interface {
	PublishOrganisationCreated(ctx context.Context, organisationID, tenantID, name, code, createdBy string) error
	PublishOrganisationUpdated(ctx context.Context, organisationID, name string, status b2bv1.OrganisationStatus, updatedBy string) error
	PublishOrganisationApproved(ctx context.Context, organisationID, approvedBy string) error
	PublishOrgMemberAdded(ctx context.Context, memberID, organisationID, userID string, role b2bv1.OrgMemberRole, addedBy string) error
	PublishOrgMemberRemoved(ctx context.Context, memberID, organisationID, userID, removedBy string) error
	PublishB2BAdminAssigned(ctx context.Context, organisationID, userID, assignedBy string) error
}

var seededCatalogPlans = map[string]*domain.CatalogPlan{
	"55555555-5555-5555-5555-555555555001": {
		ProductID:         "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1",
		ProductName:       "Health Insurance",
		PlanID:            "55555555-5555-5555-5555-555555555001",
		PlanName:          "Seba",
		InsuranceCategory: commonv1.InsuranceType_INSURANCE_TYPE_HEALTH,
		PremiumAmount: &commonv1.Money{
			Amount:        50000,
			Currency:      "BDT",
			DecimalAmount: 500,
		},
	},
	"55555555-5555-5555-5555-555555555002": {
		ProductID:         "cccccccc-cccc-cccc-cccc-ccccccccccc3",
		ProductName:       "Health Insurance",
		PlanID:            "55555555-5555-5555-5555-555555555002",
		PlanName:          "Surokkha",
		InsuranceCategory: commonv1.InsuranceType_INSURANCE_TYPE_HEALTH,
		PremiumAmount: &commonv1.Money{
			Amount:        43000,
			Currency:      "BDT",
			DecimalAmount: 430,
		},
	},
	"55555555-5555-5555-5555-555555555004": {
		ProductID:         "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2",
		ProductName:       "Life Insurance",
		PlanID:            "55555555-5555-5555-5555-555555555004",
		PlanName:          "Verosa",
		InsuranceCategory: commonv1.InsuranceType_INSURANCE_TYPE_LIFE,
		PremiumAmount: &commonv1.Money{
			Amount:        85000,
			Currency:      "BDT",
			DecimalAmount: 850,
		},
	},
}

func NewB2BService(repo domain.B2BRepository, publisher EventPublisher) *B2BService {
	return &B2BService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *B2BService) WithEmployeeIdentity(authnClient EmployeeAuthNClient, authzClient EmployeeAuthZClient) *B2BService {
	s.authn = authnClient
	s.authz = authzClient
	return s
}

func resolveTenantID(ctx context.Context, requestedTenantID string) string {
	tenantID := strings.TrimSpace(requestedTenantID)
	if tenantID != "" {
		return tenantID
	}

	if tenantID = grpcmeta.FirstFromContext(ctx, "x-tenant-id"); tenantID != "" {
		return tenantID
	}

	if envTenantID := strings.TrimSpace(os.Getenv("DEFAULT_TENANT_ID")); envTenantID != "" {
		return envTenantID
	}

	return "00000000-0000-0000-0000-000000000001"
}

// resolveCallerID extracts the acting user's ID from gRPC metadata (x-user-id header).
// Falls back to the provided default if not present.
func resolveCallerID(ctx context.Context, fallback string) string {
	if callerID := grpcmeta.FirstFromContext(ctx, "x-user-id"); callerID != "" {
		return callerID
	}
	if strings.TrimSpace(fallback) != "" {
		return strings.TrimSpace(fallback)
	}
	return "system"
}

func parseOffset(token string) int {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0
	}
	n, err := strconv.Atoi(token)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

func makePurchaseOrderNumber(now time.Time) string {
	suffix := strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))[:8]
	return fmt.Sprintf("PO-%s-%s", now.UTC().Format("20060102"), suffix)
}

func cloneMoney(value *commonv1.Money) *commonv1.Money {
	if value == nil {
		return nil
	}
	return &commonv1.Money{
		Amount:        value.Amount,
		Currency:      value.Currency,
		DecimalAmount: value.DecimalAmount,
	}
}

func resolveMetadataValue(ctx context.Context, key string) string {
	return grpcmeta.FirstFromContext(ctx, key)
}

func chooseEmployeeFullName(employee *b2bv1.Employee) string {
	if employee == nil {
		return "Employee"
	}
	if name := strings.TrimSpace(employee.GetName()); name != "" {
		return name
	}
	if email := strings.TrimSpace(employee.GetEmail()); email != "" {
		return email
	}
	return "Employee"
}

func (s *B2BService) ensureB2BBeneficiaryRole(ctx context.Context, userID, organisationID string) error {
	if s.authz == nil || strings.TrimSpace(userID) == "" || strings.TrimSpace(organisationID) == "" {
		return nil
	}

	resp, err := s.authz.ListRoles(ctx, &authzservicev1.ListRolesRequest{
		Portal:     authzentityv1.Portal_PORTAL_B2B,
		ActiveOnly: true,
		PageSize:   200,
	})
	if err != nil {
		return err
	}

	roleID := ""
	for _, role := range resp.GetRoles() {
		if role.GetName() == "b2b_beneficiary" {
			roleID = role.GetRoleId()
			break
		}
	}
	if roleID == "" {
		return fmt.Errorf("%w: authz role b2b_beneficiary not found", ErrNotFound)
	}

	_, err = s.authz.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		Domain:     "b2b:" + organisationID,
		AssignedBy: resolveCallerID(ctx, userID),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			msg := strings.ToLower(st.Message())
			if strings.Contains(msg, "duplicate") || strings.Contains(msg, "already exists") {
				return nil
			}
		}
		return err
	}
	return nil
}

func (s *B2BService) ensureEmployeeUser(ctx context.Context, employee *b2bv1.Employee) (string, error) {
	if employee == nil {
		return "", fmt.Errorf("%w: employee is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(employee.GetEmail()) == "" {
		return "", fmt.Errorf("%w: employee email is required for self-service access", ErrInvalidArgument)
	}
	if s.authn == nil {
		return strings.TrimSpace(employee.GetUserId()), nil
	}

	resp, err := s.authn.ProvisionEmployeeUser(ctx, &authnservicev1.ProvisionEmployeeUserRequest{
		Email:        strings.TrimSpace(employee.GetEmail()),
		FullName:     chooseEmployeeFullName(employee),
		MobileNumber: strings.TrimSpace(employee.GetMobileNumber()),
		EmployeeId:   strings.TrimSpace(employee.GetEmployeeId()),
		BusinessId:   strings.TrimSpace(employee.GetBusinessId()),
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.GetUserId()), nil
}

func (s *B2BService) assignedPlanNameForEmployee(ctx context.Context, employee *b2bv1.Employee) (string, error) {
	if employee == nil || strings.TrimSpace(employee.GetAssignedPlanId()) == "" {
		return "N/A", nil
	}
	planCatalog, err := s.repo.GetCatalogPlansByPlanIDs(ctx, []string{employee.GetAssignedPlanId()})
	if err != nil {
		return "", err
	}
	planCatalog = mergeCatalogMapWithSeedFallback(planCatalog)
	if plan := planCatalog[employee.GetAssignedPlanId()]; plan != nil && strings.TrimSpace(plan.PlanName) != "" {
		return plan.PlanName, nil
	}
	if fallback := fallbackCatalogPlan(employee.GetAssignedPlanId(), "", employee.GetInsuranceCategory()); fallback != nil && strings.TrimSpace(fallback.PlanName) != "" {
		return fallback.PlanName, nil
	}
	return employee.GetAssignedPlanId(), nil
}

func multiplyMoney(value *commonv1.Money, factor int32) *commonv1.Money {
	if value == nil {
		return nil
	}
	if factor < 0 {
		factor = 0
	}
	return &commonv1.Money{
		Amount:        value.Amount * int64(factor),
		Currency:      value.Currency,
		DecimalAmount: value.DecimalAmount * float64(factor),
	}
}

func insuranceCategoryDisplayName(value commonv1.InsuranceType) string {
	switch value {
	case commonv1.InsuranceType_INSURANCE_TYPE_HEALTH:
		return "Health Insurance"
	case commonv1.InsuranceType_INSURANCE_TYPE_LIFE:
		return "Life Insurance"
	case commonv1.InsuranceType_INSURANCE_TYPE_AUTO:
		return "Motor Insurance"
	case commonv1.InsuranceType_INSURANCE_TYPE_TRAVEL:
		return "Travel Insurance"
	default:
		return "Insurance Product"
	}
}

func fallbackCatalogPlan(planID, productID string, insuranceCategory commonv1.InsuranceType) *domain.CatalogPlan {
	if plan := seededCatalogPlans[planID]; plan != nil {
		return &domain.CatalogPlan{
			ProductID:         plan.ProductID,
			ProductName:       plan.ProductName,
			PlanID:            plan.PlanID,
			PlanName:          plan.PlanName,
			InsuranceCategory: plan.InsuranceCategory,
			PremiumAmount:     cloneMoney(plan.PremiumAmount),
		}
	}

	if strings.TrimSpace(planID) == "" && strings.TrimSpace(productID) == "" {
		return nil
	}

	return &domain.CatalogPlan{
		ProductID:         productID,
		ProductName:       insuranceCategoryDisplayName(insuranceCategory),
		PlanID:            planID,
		PlanName:          planID,
		InsuranceCategory: insuranceCategory,
	}
}

func mergeCatalogWithSeedFallback(items []*domain.CatalogPlan) []*domain.CatalogPlan {
	result := make([]*domain.CatalogPlan, 0, len(items)+len(seededCatalogPlans))
	seen := make(map[string]struct{}, len(items)+len(seededCatalogPlans))

	for _, item := range items {
		if item == nil || strings.TrimSpace(item.PlanID) == "" {
			continue
		}
		result = append(result, item)
		seen[item.PlanID] = struct{}{}
	}

	for planID, item := range seededCatalogPlans {
		if _, ok := seen[planID]; ok {
			continue
		}
		result = append(result, fallbackCatalogPlan(item.PlanID, item.ProductID, item.InsuranceCategory))
	}

	return result
}

func mergeCatalogMapWithSeedFallback(items map[string]*domain.CatalogPlan) map[string]*domain.CatalogPlan {
	result := make(map[string]*domain.CatalogPlan, len(items)+len(seededCatalogPlans))
	for planID, item := range items {
		if item == nil {
			continue
		}
		result[planID] = item
	}
	for planID, item := range seededCatalogPlans {
		if _, ok := result[planID]; ok {
			continue
		}
		result[planID] = fallbackCatalogPlan(item.PlanID, item.ProductID, item.InsuranceCategory)
	}
	return result
}

func normalizeOptionalEmployeeMobile(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}

	normalized, err := mobile.NormalizeBangladeshMobileE164(trimmed)
	if err != nil {
		return "", fmt.Errorf("%w: mobile_number must be a valid Bangladesh number in formats like 01712345678, +8801712345678, or 008801712345678", ErrInvalidArgument)
	}
	return normalized, nil
}

func purchaseOrderView(
	order *b2bv1.PurchaseOrder,
	departmentNames map[string]string,
	planCatalog map[string]*domain.CatalogPlan,
) *b2bservicev1.PurchaseOrderView {
	departmentName := departmentNames[order.GetDepartmentId()]
	if strings.TrimSpace(departmentName) == "" {
		departmentName = "Unassigned"
	}

	productName := "Unknown Product"
	planName := "Unknown Plan"
	plan := planCatalog[order.GetPlanId()]
	if plan == nil {
		plan = fallbackCatalogPlan(order.GetPlanId(), order.GetProductId(), order.GetInsuranceCategory())
	}
	if plan != nil {
		if strings.TrimSpace(plan.ProductName) != "" {
			productName = plan.ProductName
		}
		if strings.TrimSpace(plan.PlanName) != "" {
			planName = plan.PlanName
		}
	}

	return &b2bservicev1.PurchaseOrderView{
		PurchaseOrder:  order,
		DepartmentName: departmentName,
		ProductName:    productName,
		PlanName:       planName,
	}
}

func catalogItemsToResponse(items []*domain.CatalogPlan) []*b2bservicev1.PurchaseOrderCatalogItem {
	result := make([]*b2bservicev1.PurchaseOrderCatalogItem, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		result = append(result, &b2bservicev1.PurchaseOrderCatalogItem{
			ProductId:         item.ProductID,
			ProductName:       item.ProductName,
			PlanId:            item.PlanID,
			PlanName:          item.PlanName,
			InsuranceCategory: item.InsuranceCategory,
			PremiumAmount:     cloneMoney(item.PremiumAmount),
		})
	}
	return result
}

func (s *B2BService) ListPurchaseOrderCatalog(
	ctx context.Context,
	req *b2bservicev1.ListPurchaseOrderCatalogRequest,
) (*b2bservicev1.ListPurchaseOrderCatalogResponse, error) {
	if req == nil {
		req = &b2bservicev1.ListPurchaseOrderCatalogRequest{}
	}

	items, err := s.repo.ListCatalogPlans(ctx)
	if err != nil {
		return nil, err
	}
	items = mergeCatalogWithSeedFallback(items)

	return &b2bservicev1.ListPurchaseOrderCatalogResponse{
		Items: catalogItemsToResponse(items),
	}, nil
}

func (s *B2BService) ListPurchaseOrders(
	ctx context.Context,
	req *b2bservicev1.ListPurchaseOrdersRequest,
) (*b2bservicev1.ListPurchaseOrdersResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}

	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 100
	}
	if pageSize > 500 {
		pageSize = 500
	}
	offset := parseOffset(req.GetPageToken())

	orders, total, err := s.repo.ListPurchaseOrders(ctx, pageSize, offset, req.GetBusinessId(), req.GetStatus())
	if err != nil {
		return nil, err
	}

	departmentIDs := make([]string, 0, len(orders))
	planIDs := make([]string, 0, len(orders))
	for _, order := range orders {
		if order.GetDepartmentId() != "" {
			departmentIDs = append(departmentIDs, order.GetDepartmentId())
		}
		if order.GetPlanId() != "" {
			planIDs = append(planIDs, order.GetPlanId())
		}
	}

	departmentNames, err := s.repo.GetDepartmentNames(ctx, departmentIDs)
	if err != nil {
		return nil, err
	}
	planCatalog, err := s.repo.GetCatalogPlansByPlanIDs(ctx, planIDs)
	if err != nil {
		return nil, err
	}
	planCatalog = mergeCatalogMapWithSeedFallback(planCatalog)

	items := make([]*b2bservicev1.PurchaseOrderView, 0, len(orders))
	for _, order := range orders {
		items = append(items, purchaseOrderView(order, departmentNames, planCatalog))
	}

	next := ""
	if int64(offset+len(items)) < total {
		next = strconv.Itoa(offset + len(items))
	}

	return &b2bservicev1.ListPurchaseOrdersResponse{
		PurchaseOrders: items,
		NextPageToken:  next,
		TotalCount:     int32(total),
	}, nil
}

func (s *B2BService) GetPurchaseOrder(
	ctx context.Context,
	req *b2bservicev1.GetPurchaseOrderRequest,
) (*b2bservicev1.GetPurchaseOrderResponse, error) {
	if req == nil || strings.TrimSpace(req.GetPurchaseOrderId()) == "" {
		return nil, fmt.Errorf("%w: purchase_order_id is required", ErrInvalidArgument)
	}

	order, err := s.repo.GetPurchaseOrder(ctx, req.GetPurchaseOrderId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: purchase order not found", ErrNotFound)
		}
		return nil, err
	}

	departmentNames, err := s.repo.GetDepartmentNames(ctx, []string{order.GetDepartmentId()})
	if err != nil {
		return nil, err
	}
	planCatalog, err := s.repo.GetCatalogPlansByPlanIDs(ctx, []string{order.GetPlanId()})
	if err != nil {
		return nil, err
	}
	planCatalog = mergeCatalogMapWithSeedFallback(planCatalog)

	return &b2bservicev1.GetPurchaseOrderResponse{
		PurchaseOrder: purchaseOrderView(order, departmentNames, planCatalog),
	}, nil
}

func (s *B2BService) CreatePurchaseOrder(
	ctx context.Context,
	req *b2bservicev1.CreatePurchaseOrderRequest,
) (*b2bservicev1.CreatePurchaseOrderResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetDepartmentId()) == "" {
		return nil, fmt.Errorf("%w: department_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetPlanId()) == "" {
		return nil, fmt.Errorf("%w: plan_id is required", ErrInvalidArgument)
	}
	if req.GetEmployeeCount() <= 0 {
		return nil, fmt.Errorf("%w: employee_count must be greater than zero", ErrInvalidArgument)
	}
	if req.GetCoverageAmount() == nil || req.GetCoverageAmount().GetAmount() <= 0 {
		return nil, fmt.Errorf("%w: coverage_amount must be greater than zero", ErrInvalidArgument)
	}

	department, err := s.repo.GetDepartment(ctx, req.GetDepartmentId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: department not found", ErrInvalidArgument)
		}
		return nil, err
	}

	planCatalog, err := s.repo.GetCatalogPlansByPlanIDs(ctx, []string{req.GetPlanId()})
	if err != nil {
		return nil, err
	}
	selectedPlan := planCatalog[req.GetPlanId()]
	if selectedPlan == nil {
		selectedPlan = fallbackCatalogPlan(req.GetPlanId(), "", 0)
	}
	if selectedPlan == nil {
		return nil, fmt.Errorf("%w: plan not found", ErrInvalidArgument)
	}

	now := time.Now()
	order, err := s.repo.CreatePurchaseOrder(ctx, domain.PurchaseOrderCreateInput{
		PurchaseOrderID:     uuid.NewString(),
		PurchaseOrderNumber: makePurchaseOrderNumber(now),
		BusinessID:          department.GetBusinessId(),
		DepartmentID:        req.GetDepartmentId(),
		ProductID:           selectedPlan.ProductID,
		PlanID:              req.GetPlanId(),
		InsuranceCategory:   selectedPlan.InsuranceCategory,
		EmployeeCount:       req.GetEmployeeCount(),
		NumberOfDependents:  req.GetNumberOfDependents(),
		CoverageAmount:      cloneMoney(req.GetCoverageAmount()),
		EstimatedPremium:    multiplyMoney(selectedPlan.PremiumAmount, req.GetEmployeeCount()),
		Status:              b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED,
		RequestedBy:         strings.TrimSpace(req.GetRequestedBy()),
		Notes:               strings.TrimSpace(req.GetNotes()),
	})
	if err != nil {
		return nil, err
	}

	departmentNames := map[string]string{department.GetDepartmentId(): department.GetName()}
	return &b2bservicev1.CreatePurchaseOrderResponse{
		PurchaseOrder: purchaseOrderView(order, departmentNames, planCatalog),
		Message:       "Purchase order submitted successfully",
	}, nil
}

func (s *B2BService) ListDepartments(
	ctx context.Context,
	req *b2bservicev1.ListDepartmentsRequest,
) (*b2bservicev1.ListDepartmentsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}

	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 100
	}
	if pageSize > 500 {
		pageSize = 500
	}
	offset := parseOffset(req.GetPageToken())

	departments, total, err := s.repo.ListDepartments(ctx, pageSize, offset, req.GetBusinessId())
	if err != nil {
		return nil, err
	}

	next := ""
	if int64(offset+len(departments)) < total {
		next = strconv.Itoa(offset + len(departments))
	}

	return &b2bservicev1.ListDepartmentsResponse{
		Departments:   departments,
		NextPageToken: next,
		TotalCount:    int32(total),
	}, nil
}

func (s *B2BService) ListEmployees(
	ctx context.Context,
	req *b2bservicev1.ListEmployeesRequest,
) (*b2bservicev1.ListEmployeesResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}

	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 100
	}
	if pageSize > 500 {
		pageSize = 500
	}
	offset := parseOffset(req.GetPageToken())

	employees, total, err := s.repo.ListEmployees(
		ctx,
		pageSize,
		offset,
		req.GetDepartmentId(),
		req.GetBusinessId(),
		req.GetStatus(),
	)
	if err != nil {
		return nil, err
	}

	departmentIDs := make([]string, 0, len(employees))
	planIDs := make([]string, 0, len(employees))
	for _, employee := range employees {
		if employee.GetDepartmentId() != "" {
			departmentIDs = append(departmentIDs, employee.GetDepartmentId())
		}
		if employee.GetAssignedPlanId() != "" {
			planIDs = append(planIDs, employee.GetAssignedPlanId())
		}
	}

	departmentNames, err := s.repo.GetDepartmentNames(ctx, departmentIDs)
	if err != nil {
		return nil, err
	}
	planCatalog, err := s.repo.GetCatalogPlansByPlanIDs(ctx, planIDs)
	if err != nil {
		return nil, err
	}
	planCatalog = mergeCatalogMapWithSeedFallback(planCatalog)

	items := make([]*b2bservicev1.EmployeeView, 0, len(employees))
	for _, employee := range employees {
		departmentName := departmentNames[employee.GetDepartmentId()]
		if strings.TrimSpace(departmentName) == "" {
			departmentName = "Unassigned"
		}

		assignedPlanName := "N/A"
		if strings.TrimSpace(employee.GetAssignedPlanId()) != "" {
			if plan := planCatalog[employee.GetAssignedPlanId()]; plan != nil && strings.TrimSpace(plan.PlanName) != "" {
				assignedPlanName = plan.PlanName
			} else {
				if fallback := fallbackCatalogPlan(employee.GetAssignedPlanId(), "", employee.GetInsuranceCategory()); fallback != nil && strings.TrimSpace(fallback.PlanName) != "" {
					assignedPlanName = fallback.PlanName
				} else {
					assignedPlanName = employee.GetAssignedPlanId()
				}
			}
		}

		items = append(items, &b2bservicev1.EmployeeView{
			Employee:         employee,
			DepartmentName:   departmentName,
			AssignedPlanName: assignedPlanName,
		})
	}

	next := ""
	if int64(offset+len(items)) < total {
		next = strconv.Itoa(offset + len(items))
	}

	return &b2bservicev1.ListEmployeesResponse{
		Employees:     items,
		NextPageToken: next,
		TotalCount:    int32(total),
	}, nil
}

func (s *B2BService) GetEmployee(
	ctx context.Context,
	req *b2bservicev1.GetEmployeeRequest,
) (*b2bservicev1.GetEmployeeResponse, error) {
	if req == nil || strings.TrimSpace(req.GetEmployeeUuid()) == "" {
		return nil, fmt.Errorf("%w: employee_uuid is required", ErrInvalidArgument)
	}

	employee, err := s.repo.GetEmployee(ctx, req.GetEmployeeUuid())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee not found", ErrNotFound)
		}
		return nil, err
	}

	departmentNames, err := s.repo.GetDepartmentNames(ctx, []string{employee.GetDepartmentId()})
	if err != nil {
		return nil, err
	}
	departmentName := departmentNames[employee.GetDepartmentId()]
	if strings.TrimSpace(departmentName) == "" {
		departmentName = "Unassigned"
	}

	assignedPlanName := "N/A"
	if strings.TrimSpace(employee.GetAssignedPlanId()) != "" {
		planCatalog, err := s.repo.GetCatalogPlansByPlanIDs(ctx, []string{employee.GetAssignedPlanId()})
		if err != nil {
			return nil, err
		}
		planCatalog = mergeCatalogMapWithSeedFallback(planCatalog)
		if plan := planCatalog[employee.GetAssignedPlanId()]; plan != nil && strings.TrimSpace(plan.PlanName) != "" {
			assignedPlanName = plan.PlanName
		} else {
			if fallback := fallbackCatalogPlan(employee.GetAssignedPlanId(), "", employee.GetInsuranceCategory()); fallback != nil && strings.TrimSpace(fallback.PlanName) != "" {
				assignedPlanName = fallback.PlanName
			} else {
				assignedPlanName = employee.GetAssignedPlanId()
			}
		}
	}

	return &b2bservicev1.GetEmployeeResponse{
		Employee: &b2bservicev1.EmployeeView{
			Employee:         employee,
			DepartmentName:   departmentName,
			AssignedPlanName: assignedPlanName,
		},
	}, nil
}

func (s *B2BService) ListEmployeeLoginOrganisations(
	ctx context.Context,
	req *b2bservicev1.ListEmployeeLoginOrganisationsRequest,
) (*b2bservicev1.ListEmployeeLoginOrganisationsResponse, error) {
	if req == nil {
		req = &b2bservicev1.ListEmployeeLoginOrganisationsRequest{}
	}
	query := strings.TrimSpace(req.GetQuery())
	if len(query) < 2 {
		return &b2bservicev1.ListEmployeeLoginOrganisationsResponse{
			Organisations: []*b2bservicev1.EmployeeLoginOrganisation{},
		}, nil
	}

	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 8
	}
	if pageSize > 10 {
		pageSize = 10
	}

	organisations, err := s.repo.SearchOrganisationsForEmployeeActivation(ctx, query, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*b2bservicev1.EmployeeLoginOrganisation, 0, len(organisations))
	for _, organisation := range organisations {
		items = append(items, &b2bservicev1.EmployeeLoginOrganisation{
			OrganisationId:   organisation.GetOrganisationId(),
			OrganisationName: organisation.GetName(),
			OrganisationCode: organisation.GetCode(),
		})
	}

	return &b2bservicev1.ListEmployeeLoginOrganisationsResponse{
		Organisations: items,
	}, nil
}

// ─── CreateEmployee ───────────────────────────────────────────────────────────

func (s *B2BService) CreateEmployee(
	ctx context.Context,
	req *b2bservicev1.CreateEmployeeRequest,
) (*b2bservicev1.CreateEmployeeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetName()) == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetEmployeeId()) == "" {
		return nil, fmt.Errorf("%w: employee_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetDepartmentId()) == "" {
		return nil, fmt.Errorf("%w: department_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetBusinessId()) == "" {
		return nil, fmt.Errorf("%w: business_id is required — must be injected from session", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetEmail()) == "" {
		return nil, fmt.Errorf("%w: email is required for employee self-service login", ErrInvalidArgument)
	}
	normalizedMobileNumber, err := normalizeOptionalEmployeeMobile(req.GetMobileNumber())
	if err != nil {
		return nil, err
	}

	// Resolve plan premium before insert so it is persisted on the employee row.
	var planPremiumAmount *commonv1.Money
	assignedPlanName := "N/A"
	if strings.TrimSpace(req.GetAssignedPlanId()) != "" {
		planCatalog, err := s.repo.GetCatalogPlansByPlanIDs(ctx, []string{req.GetAssignedPlanId()})
		if err != nil {
			return nil, err
		}
		planCatalog = mergeCatalogMapWithSeedFallback(planCatalog)
		if plan := planCatalog[req.GetAssignedPlanId()]; plan != nil {
			assignedPlanName = plan.PlanName
			planPremiumAmount = cloneMoney(plan.PremiumAmount)
		}
	}

	input := domain.EmployeeCreateInput{
		EmployeeUUID:      uuid.NewString(),
		Name:              strings.TrimSpace(req.GetName()),
		EmployeeID:        strings.TrimSpace(req.GetEmployeeId()),
		DepartmentID:      req.GetDepartmentId(),
		BusinessID:        req.GetBusinessId(),
		InsuranceCategory: req.GetInsuranceCategory(),
		AssignedPlanID:    req.GetAssignedPlanId(),
		CoverageAmount:    req.GetCoverageAmount(),
		PremiumAmount:     planPremiumAmount,
		NumberOfDependent: req.GetNumberOfDependent(),
		Email:             strings.TrimSpace(req.GetEmail()),
		MobileNumber:      normalizedMobileNumber,
		DateOfBirth:       req.GetDateOfBirth(),
		DateOfJoining:     req.GetDateOfJoining(),
		Gender:            req.GetGender(),
		UserID:            strings.TrimSpace(req.GetUserId()),
	}

	if strings.TrimSpace(input.UserID) == "" {
		input.UserID, err = s.ensureEmployeeUser(ctx, &b2bv1.Employee{
			Name:         input.Name,
			EmployeeId:   input.EmployeeID,
			BusinessId:   input.BusinessID,
			Email:        input.Email,
			MobileNumber: input.MobileNumber,
		})
		if err != nil {
			return nil, err
		}
	}

	employee, err := s.repo.CreateEmployee(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := s.ensureB2BBeneficiaryRole(ctx, employee.GetUserId(), employee.GetBusinessId()); err != nil {
		return nil, err
	}

	_ = s.repo.UpdateDepartmentTotalPremium(ctx, employee.GetDepartmentId())

	// Enrich with department name (plan name already resolved above)
	departmentNames, err := s.repo.GetDepartmentNames(ctx, []string{employee.GetDepartmentId()})
	if err != nil {
		return nil, err
	}
	departmentName := departmentNames[employee.GetDepartmentId()]
	if strings.TrimSpace(departmentName) == "" {
		departmentName = "Unassigned"
	}

	return &b2bservicev1.CreateEmployeeResponse{
		Employee: &b2bservicev1.EmployeeView{
			Employee:         employee,
			DepartmentName:   departmentName,
			AssignedPlanName: assignedPlanName,
		},
		Message: "Employee created successfully",
	}, nil
}

// ─── UpdateEmployee ───────────────────────────────────────────────────────────

func (s *B2BService) UpdateEmployee(
	ctx context.Context,
	req *b2bservicev1.UpdateEmployeeRequest,
) (*b2bservicev1.UpdateEmployeeResponse, error) {
	if req == nil || strings.TrimSpace(req.GetEmployeeUuid()) == "" {
		return nil, fmt.Errorf("%w: employee_uuid is required", ErrInvalidArgument)
	}

	oldEmp, _ := s.repo.GetEmployee(ctx, req.GetEmployeeUuid())
	var oldDeptID string
	if oldEmp != nil {
		oldDeptID = oldEmp.GetDepartmentId()
	}
	normalizedMobileNumber, err := normalizeOptionalEmployeeMobile(req.GetMobileNumber())
	if err != nil {
		return nil, err
	}

	input := domain.EmployeeUpdateInput{
		EmployeeUUID:      req.GetEmployeeUuid(),
		Name:              req.GetName(),
		DepartmentID:      req.GetDepartmentId(),
		Email:             req.GetEmail(),
		MobileNumber:      normalizedMobileNumber,
		DateOfBirth:       req.GetDateOfBirth(),
		DateOfJoining:     req.GetDateOfJoining(),
		Gender:            req.GetGender(),
		InsuranceCategory: req.GetInsuranceCategory(),
		AssignedPlanID:    req.GetAssignedPlanId(),
		CoverageAmount:    req.GetCoverageAmount(),
		NumberOfDependent: req.GetNumberOfDependent(),
		Status:            req.GetStatus(),
		UserID:            strings.TrimSpace(req.GetUserId()),
	}

	employee, err := s.repo.UpdateEmployee(ctx, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee not found", ErrNotFound)
		}
		return nil, err
	}

	_ = s.repo.UpdateDepartmentTotalPremium(ctx, employee.GetDepartmentId())
	if oldDeptID != "" && oldDeptID != employee.GetDepartmentId() {
		_ = s.repo.UpdateDepartmentTotalPremium(ctx, oldDeptID)
	}

	departmentNames, err := s.repo.GetDepartmentNames(ctx, []string{employee.GetDepartmentId()})
	if err != nil {
		return nil, err
	}
	departmentName := departmentNames[employee.GetDepartmentId()]
	if strings.TrimSpace(departmentName) == "" {
		departmentName = "Unassigned"
	}

	assignedPlanName := "N/A"
	if strings.TrimSpace(employee.GetAssignedPlanId()) != "" {
		planCatalog, err := s.repo.GetCatalogPlansByPlanIDs(ctx, []string{employee.GetAssignedPlanId()})
		if err != nil {
			return nil, err
		}
		planCatalog = mergeCatalogMapWithSeedFallback(planCatalog)
		if plan := planCatalog[employee.GetAssignedPlanId()]; plan != nil {
			assignedPlanName = plan.PlanName
		}
	}

	return &b2bservicev1.UpdateEmployeeResponse{
		Employee: &b2bservicev1.EmployeeView{
			Employee:         employee,
			DepartmentName:   departmentName,
			AssignedPlanName: assignedPlanName,
		},
		Message: "Employee updated successfully",
	}, nil
}

func (s *B2BService) ActivateEmployee(
	ctx context.Context,
	req *b2bservicev1.ActivateEmployeeRequest,
) (*b2bservicev1.ActivateEmployeeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	organisationCode := strings.ToUpper(strings.TrimSpace(req.GetOrganisationCode()))
	organisationID := strings.TrimSpace(req.GetOrganisationId())
	if organisationCode == "" && organisationID == "" {
		return nil, fmt.Errorf("%w: organisation_code or organisation_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetEmployeeId()) == "" {
		return nil, fmt.Errorf("%w: employee_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetEmail()) == "" {
		return nil, fmt.Errorf("%w: email is required", ErrInvalidArgument)
	}
	if s.authn == nil {
		return nil, fmt.Errorf("%w: employee activation is unavailable", ErrInvalidArgument)
	}

	var organisation *b2bv1.Organisation
	var err error
	if organisationCode != "" {
		organisation, err = s.repo.GetOrganisationByCode(ctx, organisationCode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if organisation == nil && organisationID != "" {
		organisation, err = s.repo.GetOrganisation(ctx, organisationID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if organisation == nil {
		return nil, fmt.Errorf("%w: organisation not found for this activation request", ErrNotFound)
	}

	employee, err := s.repo.GetEmployeeByBusinessEmployeeIDEmail(
		ctx,
		organisation.GetOrganisationId(),
		strings.TrimSpace(req.GetEmployeeId()),
		strings.TrimSpace(req.GetEmail()),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee not found for this organisation and email", ErrNotFound)
		}
		return nil, err
	}

	userID := strings.TrimSpace(employee.GetUserId())
	if userID == "" {
		userID, err = s.ensureEmployeeUser(ctx, employee)
		if err != nil {
			return nil, err
		}
		employee, err = s.repo.UpdateEmployee(ctx, domain.EmployeeUpdateInput{
			EmployeeUUID: employee.GetEmployeeUuid(),
			UserID:       userID,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := s.ensureB2BBeneficiaryRole(ctx, userID, employee.GetBusinessId()); err != nil {
		return nil, err
	}

	resetResp, err := s.authn.RequestPasswordResetByEmail(ctx, &authnservicev1.RequestPasswordResetByEmailRequest{
		Email: strings.TrimSpace(employee.GetEmail()),
	})
	if err != nil {
		return nil, err
	}

	return &b2bservicev1.ActivateEmployeeResponse{
		UserId:           userID,
		OtpId:            resetResp.GetOtpId(),
		Message:          resetResp.GetMessage(),
		ExpiresInSeconds: resetResp.GetExpiresInSeconds(),
	}, nil
}

func (s *B2BService) GetMyEmployeeProfile(
	ctx context.Context,
	req *b2bservicev1.GetMyEmployeeProfileRequest,
) (*b2bservicev1.GetMyEmployeeProfileResponse, error) {
	if req == nil {
		req = &b2bservicev1.GetMyEmployeeProfileRequest{}
	}
	userID := resolveMetadataValue(ctx, "x-user-id")
	if userID == "" {
		return nil, fmt.Errorf("%w: missing user_id", ErrInvalidArgument)
	}

	employee, err := s.repo.GetEmployeeByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee not found", ErrNotFound)
		}
		return nil, err
	}

	departmentNames, err := s.repo.GetDepartmentNames(ctx, []string{employee.GetDepartmentId()})
	if err != nil {
		return nil, err
	}
	departmentName := departmentNames[employee.GetDepartmentId()]
	if strings.TrimSpace(departmentName) == "" {
		departmentName = "Unassigned"
	}
	assignedPlanName, err := s.assignedPlanNameForEmployee(ctx, employee)
	if err != nil {
		return nil, err
	}

	return &b2bservicev1.GetMyEmployeeProfileResponse{
		Employee: &b2bservicev1.EmployeeView{
			Employee:         employee,
			DepartmentName:   departmentName,
			AssignedPlanName: assignedPlanName,
		},
	}, nil
}

func (s *B2BService) GetMyEmployeeCoverage(
	ctx context.Context,
	req *b2bservicev1.GetMyEmployeeCoverageRequest,
) (*b2bservicev1.GetMyEmployeeCoverageResponse, error) {
	if req == nil {
		req = &b2bservicev1.GetMyEmployeeCoverageRequest{}
	}
	userID := resolveMetadataValue(ctx, "x-user-id")
	if userID == "" {
		return nil, fmt.Errorf("%w: missing user_id", ErrInvalidArgument)
	}

	employee, err := s.repo.GetEmployeeByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee not found", ErrNotFound)
		}
		return nil, err
	}

	organisation, err := s.repo.GetOrganisation(ctx, employee.GetBusinessId())
	if err != nil {
		return nil, err
	}
	assignedPlanName, err := s.assignedPlanNameForEmployee(ctx, employee)
	if err != nil {
		return nil, err
	}

	return &b2bservicev1.GetMyEmployeeCoverageResponse{
		Coverage: &b2bservicev1.EmployeeCoverageView{
			OrganisationId:    organisation.GetOrganisationId(),
			OrganisationName:  organisation.GetName(),
			EmployeeUuid:      employee.GetEmployeeUuid(),
			EmployeeId:        employee.GetEmployeeId(),
			AssignedPlanId:    employee.GetAssignedPlanId(),
			AssignedPlanName:  assignedPlanName,
			InsuranceCategory: employee.GetInsuranceCategory(),
			CoverageAmount:    cloneMoney(employee.GetCoverageAmount()),
			PremiumAmount:     cloneMoney(employee.GetPremiumAmount()),
			NumberOfDependent: employee.GetNumberOfDependent(),
		},
	}, nil
}

// ─── DeleteEmployee ───────────────────────────────────────────────────────────

func (s *B2BService) DeleteEmployee(
	ctx context.Context,
	req *b2bservicev1.DeleteEmployeeRequest,
) (*b2bservicev1.DeleteEmployeeResponse, error) {
	if req == nil || strings.TrimSpace(req.GetEmployeeUuid()) == "" {
		return nil, fmt.Errorf("%w: employee_uuid is required", ErrInvalidArgument)
	}

	emp, _ := s.repo.GetEmployee(ctx, req.GetEmployeeUuid())
	var deptID string
	if emp != nil {
		deptID = emp.GetDepartmentId()
	}

	if err := s.repo.DeleteEmployee(ctx, req.GetEmployeeUuid()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee not found", ErrNotFound)
		}
		return nil, err
	}

	if deptID != "" {
		_ = s.repo.UpdateDepartmentTotalPremium(ctx, deptID)
	}

	return &b2bservicev1.DeleteEmployeeResponse{Message: "Employee deleted successfully"}, nil
}

// ─── DEPARTMENT CRUD ──────────────────────────────────────────────────────────

func (s *B2BService) GetDepartment(
	ctx context.Context,
	req *b2bservicev1.GetDepartmentRequest,
) (*b2bservicev1.GetDepartmentResponse, error) {
	if req == nil || strings.TrimSpace(req.GetDepartmentId()) == "" {
		return nil, fmt.Errorf("%w: department_id is required", ErrInvalidArgument)
	}
	dept, err := s.repo.GetDepartment(ctx, req.GetDepartmentId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: department not found", ErrNotFound)
		}
		return nil, err
	}
	return &b2bservicev1.GetDepartmentResponse{Department: dept}, nil
}

func (s *B2BService) CreateDepartment(
	ctx context.Context,
	req *b2bservicev1.CreateDepartmentRequest,
) (*b2bservicev1.CreateDepartmentResponse, error) {
	if req == nil || strings.TrimSpace(req.GetName()) == "" {
		return nil, fmt.Errorf("%w: department name is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetBusinessId()) == "" {
		return nil, fmt.Errorf("%w: business_id is required (must be injected from session)", ErrInvalidArgument)
	}

	input := domain.DepartmentCreateInput{
		DepartmentID: uuid.NewString(),
		Name:         strings.TrimSpace(req.GetName()),
		BusinessID:   req.GetBusinessId(),
	}

	dept, err := s.repo.CreateDepartment(ctx, input)
	if err != nil {
		return nil, err
	}
	return &b2bservicev1.CreateDepartmentResponse{
		Department: dept,
		Message:    "Department created successfully",
	}, nil
}

func (s *B2BService) UpdateDepartment(
	ctx context.Context,
	req *b2bservicev1.UpdateDepartmentRequest,
) (*b2bservicev1.UpdateDepartmentResponse, error) {
	if req == nil || strings.TrimSpace(req.GetDepartmentId()) == "" {
		return nil, fmt.Errorf("%w: department_id is required", ErrInvalidArgument)
	}
	input := domain.DepartmentUpdateInput{
		DepartmentID: req.GetDepartmentId(),
		Name:         strings.TrimSpace(req.GetName()),
	}
	dept, err := s.repo.UpdateDepartment(ctx, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: department not found", ErrNotFound)
		}
		return nil, err
	}
	return &b2bservicev1.UpdateDepartmentResponse{
		Department: dept,
		Message:    "Department updated successfully",
	}, nil
}

func (s *B2BService) DeleteDepartment(
	ctx context.Context,
	req *b2bservicev1.DeleteDepartmentRequest,
) (*b2bservicev1.DeleteDepartmentResponse, error) {
	if req == nil || strings.TrimSpace(req.GetDepartmentId()) == "" {
		return nil, fmt.Errorf("%w: department_id is required", ErrInvalidArgument)
	}
	if err := s.repo.DeleteDepartment(ctx, req.GetDepartmentId()); err != nil {
		return nil, err
	}
	return &b2bservicev1.DeleteDepartmentResponse{Message: "Department deleted successfully"}, nil
}

// ─── ORGANISATION CRUD ────────────────────────────────────────────────────────

func (s *B2BService) CreateOrganisation(
	ctx context.Context,
	req *b2bservicev1.CreateOrganisationRequest,
) (*b2bservicev1.CreateOrganisationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetName()) == "" {
		return nil, fmt.Errorf("%w: organisation name is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetCode()) == "" {
		return nil, fmt.Errorf("%w: organisation code is required", ErrInvalidArgument)
	}
	tenantID := resolveTenantID(ctx, req.GetTenantId())
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidArgument)
	}

	input := domain.OrganisationCreateInput{
		OrganisationID: uuid.NewString(),
		TenantID:       tenantID,
		Name:           strings.TrimSpace(req.GetName()),
		Code:           strings.ToUpper(strings.TrimSpace(req.GetCode())),
		Industry:       req.GetIndustry(),
		ContactEmail:   req.GetContactEmail(),
		ContactPhone:   req.GetContactPhone(),
		Address:        req.GetAddress(),
	}

	org, err := s.repo.CreateOrganisation(ctx, input)
	if err != nil {
		return nil, err
	}

	// Publish event
	if s.publisher != nil {
		callerID := resolveCallerID(ctx, "system")
		_ = s.publisher.PublishOrganisationCreated(ctx, org.OrganisationId, org.TenantId, org.Name, org.Code, callerID)
	}

	return &b2bservicev1.CreateOrganisationResponse{
		Organisation: org,
		Message:      "Organisation created successfully",
	}, nil
}

func (s *B2BService) GetOrganisation(
	ctx context.Context,
	req *b2bservicev1.GetOrganisationRequest,
) (*b2bservicev1.GetOrganisationResponse, error) {
	if req == nil || strings.TrimSpace(req.GetOrganisationId()) == "" {
		return nil, fmt.Errorf("%w: organisation_id is required", ErrInvalidArgument)
	}
	org, err := s.repo.GetOrganisation(ctx, req.GetOrganisationId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: organisation not found", ErrNotFound)
		}
		return nil, err
	}
	return &b2bservicev1.GetOrganisationResponse{Organisation: org}, nil
}

func (s *B2BService) ListOrganisations(
	ctx context.Context,
	req *b2bservicev1.ListOrganisationsRequest,
) (*b2bservicev1.ListOrganisationsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}
	offset := parseOffset(req.GetPageToken())
	tenantID := resolveTenantID(ctx, req.GetTenantId())

	orgs, total, err := s.repo.ListOrganisations(ctx, pageSize, offset, tenantID, req.GetStatus())
	if err != nil {
		return nil, err
	}

	next := ""
	if int64(offset+len(orgs)) < total {
		next = strconv.Itoa(offset + len(orgs))
	}
	return &b2bservicev1.ListOrganisationsResponse{
		Organisations: orgs,
		NextPageToken: next,
		TotalCount:    int32(total),
	}, nil
}

func (s *B2BService) UpdateOrganisation(
	ctx context.Context,
	req *b2bservicev1.UpdateOrganisationRequest,
) (*b2bservicev1.UpdateOrganisationResponse, error) {
	if req == nil || strings.TrimSpace(req.GetOrganisationId()) == "" {
		return nil, fmt.Errorf("%w: organisation_id is required", ErrInvalidArgument)
	}
	input := domain.OrganisationUpdateInput{
		OrganisationID: req.GetOrganisationId(),
		Name:           req.GetName(),
		Industry:       req.GetIndustry(),
		ContactEmail:   req.GetContactEmail(),
		ContactPhone:   req.GetContactPhone(),
		Address:        req.GetAddress(),
		Status:         req.GetStatus(),
	}
	org, err := s.repo.UpdateOrganisation(ctx, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: organisation not found", ErrNotFound)
		}
		return nil, err
	}
	return &b2bservicev1.UpdateOrganisationResponse{
		Organisation: org,
		Message:      "Organisation updated successfully",
	}, nil
}

func (s *B2BService) DeleteOrganisation(
	ctx context.Context,
	req *b2bservicev1.DeleteOrganisationRequest,
) (*b2bservicev1.DeleteOrganisationResponse, error) {
	if req == nil || strings.TrimSpace(req.GetOrganisationId()) == "" {
		return nil, fmt.Errorf("%w: organisation_id is required", ErrInvalidArgument)
	}
	if err := s.repo.DeleteOrganisation(ctx, req.GetOrganisationId()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: organisation not found", ErrNotFound)
		}
		return nil, err
	}
	return &b2bservicev1.DeleteOrganisationResponse{
		Message: "Organisation deleted successfully",
	}, nil
}

func (s *B2BService) ListOrgMembers(
	ctx context.Context,
	req *b2bservicev1.ListOrgMembersRequest,
) (*b2bservicev1.ListOrgMembersResponse, error) {
	if req == nil || strings.TrimSpace(req.GetOrganisationId()) == "" {
		return nil, fmt.Errorf("%w: organisation_id is required", ErrInvalidArgument)
	}
	members, err := s.repo.ListOrgMembers(ctx, req.GetOrganisationId())
	if err != nil {
		return nil, err
	}
	return &b2bservicev1.ListOrgMembersResponse{Members: members}, nil
}

func (s *B2BService) AddOrgMember(
	ctx context.Context,
	req *b2bservicev1.AddOrgMemberRequest,
) (*b2bservicev1.AddOrgMemberResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetOrganisationId()) == "" {
		return nil, fmt.Errorf("%w: organisation_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetUserId()) == "" {
		return nil, fmt.Errorf("%w: user_id is required", ErrInvalidArgument)
	}

	input := domain.OrgMemberCreateInput{
		MemberID:       uuid.NewString(),
		OrganisationID: req.GetOrganisationId(),
		UserID:         req.GetUserId(),
		Role:           req.GetRole(),
	}
	member, err := s.repo.AddOrgMember(ctx, input)
	if err != nil {
		return nil, err
	}

	// Publish event
	if s.publisher != nil {
		callerID := resolveCallerID(ctx, "system")
		_ = s.publisher.PublishOrgMemberAdded(ctx, member.MemberId, member.OrganisationId, member.UserId, member.Role, callerID)

		// If role is BUSINESS_ADMIN, also publish admin assignment event so
		// the authz consumer assigns the b2b_org_admin Casbin role immediately.
		if member.Role == b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN {
			_ = s.publisher.PublishB2BAdminAssigned(ctx, member.OrganisationId, member.UserId, callerID)
		}
	}

	return &b2bservicev1.AddOrgMemberResponse{
		Member:  member,
		Message: "Member added successfully",
	}, nil
}

func (s *B2BService) AssignOrgAdmin(
	ctx context.Context,
	req *b2bservicev1.AssignOrgAdminRequest,
) (*b2bservicev1.AssignOrgAdminResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetOrganisationId()) == "" || strings.TrimSpace(req.GetMemberId()) == "" {
		return nil, fmt.Errorf("%w: organisation_id and member_id are required", ErrInvalidArgument)
	}

	member, err := s.repo.AssignOrgAdmin(ctx, req.GetOrganisationId(), req.GetMemberId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: member not found", ErrNotFound)
		}
		return nil, err
	}

	if s.publisher != nil {
		callerID := resolveCallerID(ctx, "system")
		_ = s.publisher.PublishB2BAdminAssigned(ctx, member.OrganisationId, member.UserId, callerID)
	}

	return &b2bservicev1.AssignOrgAdminResponse{
		Member:  member,
		Message: "B2B admin assigned successfully",
	}, nil
}

func (s *B2BService) RemoveOrgMember(
	ctx context.Context,
	req *b2bservicev1.RemoveOrgMemberRequest,
) (*b2bservicev1.RemoveOrgMemberResponse, error) {
	if req == nil || strings.TrimSpace(req.GetOrganisationId()) == "" || strings.TrimSpace(req.GetMemberId()) == "" {
		return nil, fmt.Errorf("%w: organisation_id and member_id are required", ErrInvalidArgument)
	}
	if err := s.repo.RemoveOrgMember(ctx, req.GetOrganisationId(), req.GetMemberId()); err != nil {
		return nil, err
	}

	// Publish event
	if s.publisher != nil {
		callerID := resolveCallerID(ctx, "system")
		_ = s.publisher.PublishOrgMemberRemoved(ctx, req.GetMemberId(), req.GetOrganisationId(), "", callerID)
	}

	return &b2bservicev1.RemoveOrgMemberResponse{Message: "Member removed successfully"}, nil
}

// ResolveMyOrganisation resolves the organisation_id for the authenticated user.
// This is the core fix for the hardcoded business_id problem:
// The REST gateway calls this RPC with the user_id from the validated session,
// then injects the returned organisation_id as x-business-id gRPC metadata
// into every subsequent B2B service call.
func (s *B2BService) ResolveMyOrganisation(
	ctx context.Context,
	req *b2bservicev1.ResolveMyOrganisationRequest,
) (*b2bservicev1.ResolveMyOrganisationResponse, error) {
	if req == nil || strings.TrimSpace(req.GetUserId()) == "" {
		return nil, fmt.Errorf("%w: user_id is required", ErrInvalidArgument)
	}

	orgID, role, orgName, err := s.repo.ResolveOrganisationByUserID(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: no organisation found for user", ErrNotFound)
		}
		return nil, err
	}

	return &b2bservicev1.ResolveMyOrganisationResponse{
		OrganisationId:   orgID,
		OrganisationName: orgName,
		Role:             role,
	}, nil
}

// keepTimeImport prevents the time import from being flagged unused during transition.
var keepTimeImport = time.RFC3339
