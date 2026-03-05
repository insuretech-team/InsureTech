package domain

import (
	"context"

	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type CatalogPlan struct {
	ProductID         string
	ProductName       string
	PlanID            string
	PlanName          string
	InsuranceCategory commonv1.InsuranceType
	PremiumAmount     *commonv1.Money
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

type B2BRepository interface {
	ListDepartments(ctx context.Context, pageSize, offset int, businessID string) ([]*b2bv1.Department, int64, error)
	GetDepartment(ctx context.Context, departmentID string) (*b2bv1.Department, error)
	ListEmployees(ctx context.Context, pageSize, offset int, departmentID, businessID string) ([]*b2bv1.Employee, int64, error)
	GetEmployee(ctx context.Context, employeeUUID string) (*b2bv1.Employee, error)
	GetDepartmentNames(ctx context.Context, departmentIDs []string) (map[string]string, error)
	ListCatalogPlans(ctx context.Context) ([]*CatalogPlan, error)
	GetCatalogPlansByPlanIDs(ctx context.Context, planIDs []string) (map[string]*CatalogPlan, error)
	ListPurchaseOrders(ctx context.Context, pageSize, offset int, businessID string, status b2bv1.PurchaseOrderStatus) ([]*b2bv1.PurchaseOrder, int64, error)
	GetPurchaseOrder(ctx context.Context, purchaseOrderID string) (*b2bv1.PurchaseOrder, error)
	CreatePurchaseOrder(ctx context.Context, input PurchaseOrderCreateInput) (*b2bv1.PurchaseOrder, error)
}

type B2BService interface {
	ListPurchaseOrderCatalog(ctx context.Context, req *b2bservicev1.ListPurchaseOrderCatalogRequest) (*b2bservicev1.ListPurchaseOrderCatalogResponse, error)
	ListPurchaseOrders(ctx context.Context, req *b2bservicev1.ListPurchaseOrdersRequest) (*b2bservicev1.ListPurchaseOrdersResponse, error)
	GetPurchaseOrder(ctx context.Context, req *b2bservicev1.GetPurchaseOrderRequest) (*b2bservicev1.GetPurchaseOrderResponse, error)
	CreatePurchaseOrder(ctx context.Context, req *b2bservicev1.CreatePurchaseOrderRequest) (*b2bservicev1.CreatePurchaseOrderResponse, error)
	ListDepartments(ctx context.Context, req *b2bservicev1.ListDepartmentsRequest) (*b2bservicev1.ListDepartmentsResponse, error)
	ListEmployees(ctx context.Context, req *b2bservicev1.ListEmployeesRequest) (*b2bservicev1.ListEmployeesResponse, error)
	GetEmployee(ctx context.Context, req *b2bservicev1.GetEmployeeRequest) (*b2bservicev1.GetEmployeeResponse, error)
}
