package grpc

import (
	"context"
	"errors"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/service"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type B2BHandler struct {
	b2bservicev1.UnimplementedB2BServiceServer
	svc domain.B2BService
}

func NewB2BHandler(svc domain.B2BService) *B2BHandler {
	return &B2BHandler{svc: svc}
}

func (h *B2BHandler) ListPurchaseOrderCatalog(ctx context.Context, req *b2bservicev1.ListPurchaseOrderCatalogRequest) (*b2bservicev1.ListPurchaseOrderCatalogResponse, error) {
	resp, err := h.svc.ListPurchaseOrderCatalog(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) ListPurchaseOrders(ctx context.Context, req *b2bservicev1.ListPurchaseOrdersRequest) (*b2bservicev1.ListPurchaseOrdersResponse, error) {
	resp, err := h.svc.ListPurchaseOrders(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) GetPurchaseOrder(ctx context.Context, req *b2bservicev1.GetPurchaseOrderRequest) (*b2bservicev1.GetPurchaseOrderResponse, error) {
	resp, err := h.svc.GetPurchaseOrder(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) CreatePurchaseOrder(ctx context.Context, req *b2bservicev1.CreatePurchaseOrderRequest) (*b2bservicev1.CreatePurchaseOrderResponse, error) {
	resp, err := h.svc.CreatePurchaseOrder(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) ListDepartments(ctx context.Context, req *b2bservicev1.ListDepartmentsRequest) (*b2bservicev1.ListDepartmentsResponse, error) {
	resp, err := h.svc.ListDepartments(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) ListEmployees(ctx context.Context, req *b2bservicev1.ListEmployeesRequest) (*b2bservicev1.ListEmployeesResponse, error) {
	resp, err := h.svc.ListEmployees(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) GetEmployee(ctx context.Context, req *b2bservicev1.GetEmployeeRequest) (*b2bservicev1.GetEmployeeResponse, error) {
	resp, err := h.svc.GetEmployee(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
