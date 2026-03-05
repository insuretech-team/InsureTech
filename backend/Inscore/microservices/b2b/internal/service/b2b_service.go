package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"gorm.io/gorm"
)

type B2BService struct {
	repo domain.B2BRepository
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

func NewB2BService(repo domain.B2BRepository) *B2BService {
	return &B2BService{repo: repo}
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
