package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type B2BServiceHandler struct {
	client B2BServiceClient
}

type B2BServiceClient interface {
	ListPurchaseOrderCatalog(ctx context.Context, in *b2bservicev1.ListPurchaseOrderCatalogRequest, opts ...grpc.CallOption) (*b2bservicev1.ListPurchaseOrderCatalogResponse, error)
	ListPurchaseOrders(ctx context.Context, in *b2bservicev1.ListPurchaseOrdersRequest, opts ...grpc.CallOption) (*b2bservicev1.ListPurchaseOrdersResponse, error)
	GetPurchaseOrder(ctx context.Context, in *b2bservicev1.GetPurchaseOrderRequest, opts ...grpc.CallOption) (*b2bservicev1.GetPurchaseOrderResponse, error)
	CreatePurchaseOrder(ctx context.Context, in *b2bservicev1.CreatePurchaseOrderRequest, opts ...grpc.CallOption) (*b2bservicev1.CreatePurchaseOrderResponse, error)
	ListDepartments(ctx context.Context, in *b2bservicev1.ListDepartmentsRequest, opts ...grpc.CallOption) (*b2bservicev1.ListDepartmentsResponse, error)
	ListEmployees(ctx context.Context, in *b2bservicev1.ListEmployeesRequest, opts ...grpc.CallOption) (*b2bservicev1.ListEmployeesResponse, error)
	GetEmployee(ctx context.Context, in *b2bservicev1.GetEmployeeRequest, opts ...grpc.CallOption) (*b2bservicev1.GetEmployeeResponse, error)
}

func NewB2BServiceHandler(conn *grpc.ClientConn) *B2BServiceHandler {
	return &B2BServiceHandler{client: b2bservicev1.NewB2BServiceClient(conn)}
}

func parsePurchaseOrderStatusQuery(value string) b2bv1.PurchaseOrderStatus {
	value = strings.TrimSpace(value)
	if value == "" {
		return b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED
	}
	if n, err := strconv.Atoi(value); err == nil {
		return b2bv1.PurchaseOrderStatus(n)
	}
	if enumValue, ok := b2bv1.PurchaseOrderStatus_value[strings.ToUpper(value)]; ok {
		return b2bv1.PurchaseOrderStatus(enumValue)
	}
	normalized := "PURCHASE_ORDER_STATUS_" + strings.ToUpper(value)
	if enumValue, ok := b2bv1.PurchaseOrderStatus_value[normalized]; ok {
		return b2bv1.PurchaseOrderStatus(enumValue)
	}
	return b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED
}

func (h *B2BServiceHandler) ListPurchaseOrderCatalog(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListPurchaseOrderCatalog(ctx, &b2bservicev1.ListPurchaseOrderCatalogRequest{})
	})
}

func (h *B2BServiceHandler) ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &b2bservicev1.ListPurchaseOrdersRequest{
			PageToken:  r.URL.Query().Get("page_token"),
			BusinessId: r.URL.Query().Get("business_id"),
			Status:     parsePurchaseOrderStatusQuery(r.URL.Query().Get("status")),
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListPurchaseOrders(ctx, req)
	})
}

func (h *B2BServiceHandler) GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	purchaseOrderID := r.PathValue("purchase_order_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetPurchaseOrder(ctx, &b2bservicev1.GetPurchaseOrderRequest{
			PurchaseOrderId: purchaseOrderID,
		})
	})
}

func (h *B2BServiceHandler) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.CreatePurchaseOrderRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreatePurchaseOrder(ctx, &req)
	})
}

func (h *B2BServiceHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &b2bservicev1.ListDepartmentsRequest{
			PageToken:  r.URL.Query().Get("page_token"),
			BusinessId: r.URL.Query().Get("business_id"),
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListDepartments(ctx, req)
	})
}

func (h *B2BServiceHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &b2bservicev1.ListEmployeesRequest{
			PageToken:    r.URL.Query().Get("page_token"),
			DepartmentId: r.URL.Query().Get("department_id"),
			BusinessId:   r.URL.Query().Get("business_id"),
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListEmployees(ctx, req)
	})
}

func (h *B2BServiceHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	employeeUUID := r.PathValue("employee_uuid")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetEmployee(ctx, &b2bservicev1.GetEmployeeRequest{
			EmployeeUuid: employeeUUID,
		})
	})
}
