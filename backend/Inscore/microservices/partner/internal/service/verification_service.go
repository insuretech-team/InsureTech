package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// VerifyPartner handles KYB or manual verification checks
func (s *PartnerService) VerifyPartner(ctx context.Context, req *partnerservicev1.VerifyPartnerRequest) (*partnerservicev1.VerifyPartnerResponse, error) {
	if req.PartnerId == "" || req.VerificationType == "" {
		return nil, fmt.Errorf("%w: invalid verify request", ErrInvalidArgument)
	}
	verifiedBy := "system"
	if req.VerificationData != nil {
		if v := strings.TrimSpace(req.VerificationData["verified_by"]); v != "" {
			verifiedBy = v
		}
	}

	err := s.partnerRepo.UpdateStatus(ctx, req.PartnerId, partnerv1.PartnerStatus_PARTNER_STATUS_ACTIVE)
	if err != nil {
		logger.Errorf("Failed to update status during verification: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: partner not found", ErrNotFound)
		}
		return nil, fmt.Errorf("verify partner: %w", err)
	}

	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishPartnerVerified(ctx, req.PartnerId, verifiedBy)
	}
	s.metrics.IncPartnerVerified()

	return &partnerservicev1.VerifyPartnerResponse{
		Verified:           true,
		VerificationStatus: "APPROVED",
		VerifiedAt:         timestamppb.New(time.Now()),
		VerifiedBy:         verifiedBy,
	}, nil
}

// UpdatePartnerStatus allows admin or focal person to pause operations
func (s *PartnerService) UpdatePartnerStatus(ctx context.Context, req *partnerservicev1.UpdatePartnerStatusRequest) (*partnerservicev1.UpdatePartnerStatusResponse, error) {
	if req.PartnerId == "" || req.Status == "" {
		return nil, fmt.Errorf("%w: partner_id and status are required", ErrInvalidArgument)
	}
	status, ok := parsePartnerStatus(req.Status)
	if !ok {
		return nil, fmt.Errorf("%w: invalid status enum string", ErrInvalidArgument)
	}

	err := s.partnerRepo.UpdateStatus(ctx, req.PartnerId, status)
	if err != nil {
		logger.Errorf("Failed to update status manually: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: partner not found", ErrNotFound)
		}
		return nil, fmt.Errorf("update partner status: %w", err)
	}

	updated, err := s.partnerRepo.GetByID(ctx, req.PartnerId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: partner not found", ErrNotFound)
		}
		return nil, fmt.Errorf("load partner after status update: %w", err)
	}

	if status == partnerv1.PartnerStatus_PARTNER_STATUS_ACTIVE && s.eventPublisher != nil {
		_ = s.eventPublisher.PublishPartnerVerified(ctx, req.PartnerId, "admin")
	}
	s.metrics.IncPartnerStatusUpdate()

	return &partnerservicev1.UpdatePartnerStatusResponse{
		Partner: updated,
	}, nil
}

func parsePartnerStatus(raw string) (partnerv1.PartnerStatus, bool) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return partnerv1.PartnerStatus_PARTNER_STATUS_UNSPECIFIED, false
	}
	if iv, ok := partnerv1.PartnerStatus_value[v]; ok {
		return partnerv1.PartnerStatus(iv), true
	}
	if !strings.HasPrefix(v, "PARTNER_STATUS_") {
		key := "PARTNER_STATUS_" + strings.ToUpper(v)
		if iv, ok := partnerv1.PartnerStatus_value[key]; ok {
			return partnerv1.PartnerStatus(iv), true
		}
	}
	return partnerv1.PartnerStatus_PARTNER_STATUS_UNSPECIFIED, false
}
