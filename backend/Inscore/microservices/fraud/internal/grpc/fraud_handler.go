package grpc

import (
	"context"
	"errors"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/service"
	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FraudHandler exposes FraudService through gRPC transport.
type FraudHandler struct {
	fraudservicev1.UnimplementedFraudServiceServer
	svc *service.FraudService
}

func NewFraudHandler(svc *service.FraudService) *FraudHandler {
	return &FraudHandler{svc: svc}
}

func (h *FraudHandler) CheckFraud(ctx context.Context, req *fraudservicev1.CheckFraudRequest) (*fraudservicev1.CheckFraudResponse, error) {
	resp, err := h.svc.CheckFraud(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) GetFraudAlert(ctx context.Context, req *fraudservicev1.GetFraudAlertRequest) (*fraudservicev1.GetFraudAlertResponse, error) {
	resp, err := h.svc.GetFraudAlert(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) ListFraudAlerts(ctx context.Context, req *fraudservicev1.ListFraudAlertsRequest) (*fraudservicev1.ListFraudAlertsResponse, error) {
	resp, err := h.svc.ListFraudAlerts(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) CreateFraudCase(ctx context.Context, req *fraudservicev1.CreateFraudCaseRequest) (*fraudservicev1.CreateFraudCaseResponse, error) {
	resp, err := h.svc.CreateFraudCase(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) GetFraudCase(ctx context.Context, req *fraudservicev1.GetFraudCaseRequest) (*fraudservicev1.GetFraudCaseResponse, error) {
	resp, err := h.svc.GetFraudCase(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) UpdateFraudCase(ctx context.Context, req *fraudservicev1.UpdateFraudCaseRequest) (*fraudservicev1.UpdateFraudCaseResponse, error) {
	resp, err := h.svc.UpdateFraudCase(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) ListFraudRules(ctx context.Context, req *fraudservicev1.ListFraudRulesRequest) (*fraudservicev1.ListFraudRulesResponse, error) {
	resp, err := h.svc.ListFraudRules(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) CreateFraudRule(ctx context.Context, req *fraudservicev1.CreateFraudRuleRequest) (*fraudservicev1.CreateFraudRuleResponse, error) {
	resp, err := h.svc.CreateFraudRule(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) UpdateFraudRule(ctx context.Context, req *fraudservicev1.UpdateFraudRuleRequest) (*fraudservicev1.UpdateFraudRuleResponse, error) {
	resp, err := h.svc.UpdateFraudRule(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) ActivateFraudRule(ctx context.Context, req *fraudservicev1.ActivateFraudRuleRequest) (*fraudservicev1.ActivateFraudRuleResponse, error) {
	resp, err := h.svc.ActivateFraudRule(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *FraudHandler) DeactivateFraudRule(ctx context.Context, req *fraudservicev1.DeactivateFraudRuleRequest) (*fraudservicev1.DeactivateFraudRuleResponse, error) {
	resp, err := h.svc.DeactivateFraudRule(ctx, req)
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
