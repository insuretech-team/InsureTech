package grpc

import (
	"context"
	"errors"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/service"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	paymentservicev1.UnimplementedPaymentServiceServer
	svc domain.PaymentService
}

func NewPaymentHandler(svc domain.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) InitiatePayment(ctx context.Context, req *paymentservicev1.InitiatePaymentRequest) (*paymentservicev1.InitiatePaymentResponse, error) {
	resp, err := h.svc.InitiatePayment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) VerifyPayment(ctx context.Context, req *paymentservicev1.VerifyPaymentRequest) (*paymentservicev1.VerifyPaymentResponse, error) {
	resp, err := h.svc.VerifyPayment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *paymentservicev1.GetPaymentRequest) (*paymentservicev1.GetPaymentResponse, error) {
	resp, err := h.svc.GetPayment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) ListPayments(ctx context.Context, req *paymentservicev1.ListPaymentsRequest) (*paymentservicev1.ListPaymentsResponse, error) {
	resp, err := h.svc.ListPayments(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) InitiateRefund(ctx context.Context, req *paymentservicev1.InitiateRefundRequest) (*paymentservicev1.InitiateRefundResponse, error) {
	resp, err := h.svc.InitiateRefund(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) GetRefundStatus(ctx context.Context, req *paymentservicev1.GetRefundStatusRequest) (*paymentservicev1.GetRefundStatusResponse, error) {
	resp, err := h.svc.GetRefundStatus(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) ListPaymentMethods(ctx context.Context, req *paymentservicev1.ListPaymentMethodsRequest) (*paymentservicev1.ListPaymentMethodsResponse, error) {
	resp, err := h.svc.ListPaymentMethods(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) AddPaymentMethod(ctx context.Context, req *paymentservicev1.AddPaymentMethodRequest) (*paymentservicev1.AddPaymentMethodResponse, error) {
	resp, err := h.svc.AddPaymentMethod(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) ReconcilePayments(ctx context.Context, req *paymentservicev1.ReconcilePaymentsRequest) (*paymentservicev1.ReconcilePaymentsResponse, error) {
	resp, err := h.svc.ReconcilePayments(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) HandleGatewayWebhook(ctx context.Context, req *paymentservicev1.HandleGatewayWebhookRequest) (*paymentservicev1.HandleGatewayWebhookResponse, error) {
	resp, err := h.svc.HandleGatewayWebhook(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) GetPaymentByProviderReference(ctx context.Context, req *paymentservicev1.GetPaymentByProviderReferenceRequest) (*paymentservicev1.GetPaymentByProviderReferenceResponse, error) {
	resp, err := h.svc.GetPaymentByProviderReference(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) SubmitManualPaymentProof(ctx context.Context, req *paymentservicev1.SubmitManualPaymentProofRequest) (*paymentservicev1.SubmitManualPaymentProofResponse, error) {
	resp, err := h.svc.SubmitManualPaymentProof(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) ReviewManualPayment(ctx context.Context, req *paymentservicev1.ReviewManualPaymentRequest) (*paymentservicev1.ReviewManualPaymentResponse, error) {
	resp, err := h.svc.ReviewManualPayment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) GenerateReceipt(ctx context.Context, req *paymentservicev1.GenerateReceiptRequest) (*paymentservicev1.GenerateReceiptResponse, error) {
	resp, err := h.svc.GenerateReceipt(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PaymentHandler) GetPaymentReceipt(ctx context.Context, req *paymentservicev1.GetPaymentReceiptRequest) (*paymentservicev1.GetPaymentReceiptResponse, error) {
	resp, err := h.svc.GetPaymentReceipt(ctx, req)
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
