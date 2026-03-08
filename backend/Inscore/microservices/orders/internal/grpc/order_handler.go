package grpc

import (
	"context"
	"errors"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/service"
	orderservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	orderservicev1.UnimplementedOrderServiceServer
	svc domain.OrderService
}

func NewOrderHandler(svc domain.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// ─── ORDER LIFECYCLE ──────────────────────────────────────────────────────────

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderservicev1.CreateOrderRequest) (*orderservicev1.CreateOrderResponse, error) {
	resp, err := h.svc.CreateOrder(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *orderservicev1.GetOrderRequest) (*orderservicev1.GetOrderResponse, error) {
	resp, err := h.svc.GetOrder(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *orderservicev1.ListOrdersRequest) (*orderservicev1.ListOrdersResponse, error) {
	resp, err := h.svc.ListOrders(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *OrderHandler) InitiatePayment(ctx context.Context, req *orderservicev1.InitiatePaymentRequest) (*orderservicev1.InitiatePaymentResponse, error) {
	resp, err := h.svc.InitiatePayment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *OrderHandler) ConfirmPayment(ctx context.Context, req *orderservicev1.ConfirmPaymentRequest) (*orderservicev1.ConfirmPaymentResponse, error) {
	resp, err := h.svc.ConfirmPayment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *orderservicev1.CancelOrderRequest) (*orderservicev1.CancelOrderResponse, error) {
	resp, err := h.svc.CancelOrder(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *OrderHandler) GetOrderStatus(ctx context.Context, req *orderservicev1.GetOrderStatusRequest) (*orderservicev1.GetOrderStatusResponse, error) {
	resp, err := h.svc.GetOrderStatus(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

// ─── ERROR MAPPER ─────────────────────────────────────────────────────────────

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrInvalidTransition):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, service.ErrPaymentFailed):
		return status.Error(codes.Aborted, err.Error())
	case errors.Is(err, service.ErrNotImplemented):
		return status.Error(codes.Unimplemented, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
