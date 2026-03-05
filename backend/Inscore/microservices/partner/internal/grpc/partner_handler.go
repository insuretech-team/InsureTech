package grpc

import (
	"context"
	"errors"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/service"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PartnerHandler exposes partner business logic through gRPC transport.
type PartnerHandler struct {
	partnerservicev1.UnimplementedPartnerServiceServer
	svc domain.PartnerService
}

func NewPartnerHandler(svc domain.PartnerService) *PartnerHandler {
	return &PartnerHandler{svc: svc}
}

func (h *PartnerHandler) CreatePartner(ctx context.Context, req *partnerservicev1.CreatePartnerRequest) (*partnerservicev1.CreatePartnerResponse, error) {
	resp, err := h.svc.CreatePartner(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) GetPartner(ctx context.Context, req *partnerservicev1.GetPartnerRequest) (*partnerservicev1.GetPartnerResponse, error) {
	resp, err := h.svc.GetPartner(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) UpdatePartner(ctx context.Context, req *partnerservicev1.UpdatePartnerRequest) (*partnerservicev1.UpdatePartnerResponse, error) {
	resp, err := h.svc.UpdatePartner(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) ListPartners(ctx context.Context, req *partnerservicev1.ListPartnersRequest) (*partnerservicev1.ListPartnersResponse, error) {
	resp, err := h.svc.ListPartners(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) DeletePartner(ctx context.Context, req *partnerservicev1.DeletePartnerRequest) (*partnerservicev1.DeletePartnerResponse, error) {
	resp, err := h.svc.DeletePartner(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) VerifyPartner(ctx context.Context, req *partnerservicev1.VerifyPartnerRequest) (*partnerservicev1.VerifyPartnerResponse, error) {
	resp, err := h.svc.VerifyPartner(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) UpdatePartnerStatus(ctx context.Context, req *partnerservicev1.UpdatePartnerStatusRequest) (*partnerservicev1.UpdatePartnerStatusResponse, error) {
	resp, err := h.svc.UpdatePartnerStatus(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) GetPartnerCommission(ctx context.Context, req *partnerservicev1.GetPartnerCommissionRequest) (*partnerservicev1.GetPartnerCommissionResponse, error) {
	resp, err := h.svc.GetPartnerCommission(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) UpdateCommissionStructure(ctx context.Context, req *partnerservicev1.UpdateCommissionStructureRequest) (*partnerservicev1.UpdateCommissionStructureResponse, error) {
	resp, err := h.svc.UpdateCommissionStructure(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) GetPartnerAPICredentials(ctx context.Context, req *partnerservicev1.GetPartnerAPICredentialsRequest) (*partnerservicev1.GetPartnerAPICredentialsResponse, error) {
	resp, err := h.svc.GetPartnerAPICredentials(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *PartnerHandler) RotatePartnerAPIKey(ctx context.Context, req *partnerservicev1.RotatePartnerAPIKeyRequest) (*partnerservicev1.RotatePartnerAPIKeyResponse, error) {
	resp, err := h.svc.RotatePartnerAPIKey(ctx, req)
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
	case errors.Is(err, service.ErrConflict):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrUnavailable):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
