package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/metrics"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"gorm.io/gorm"
)

// PartnerService contains partner business logic.
type PartnerService struct {
	partnerRepo    domain.PartnerRepository
	agentRepo      domain.AgentRepository
	commissionRepo domain.CommissionRepository
	eventPublisher domain.EventPublisher
	authnClient    domain.AuthNClient
	metrics        *metrics.RuntimeMetrics
}

// NewPartnerService creates a new PartnerService
func NewPartnerService(
	pRepo domain.PartnerRepository,
	aRepo domain.AgentRepository,
	cRepo domain.CommissionRepository,
	evtPub domain.EventPublisher,
	authnClient domain.AuthNClient,
) *PartnerService {
	return &PartnerService{
		partnerRepo:    pRepo,
		agentRepo:      aRepo,
		commissionRepo: cRepo,
		eventPublisher: evtPub,
		authnClient:    authnClient,
		metrics:        metrics.NewRuntimeMetrics(),
	}
}

// CreatePartner creates a new partner record
func (s *PartnerService) CreatePartner(ctx context.Context, req *partnerservicev1.CreatePartnerRequest) (*partnerservicev1.CreatePartnerResponse, error) {
	if req.Partner == nil {
		return nil, fmt.Errorf("%w: partner details are required", ErrInvalidArgument)
	}

	err := s.partnerRepo.Create(ctx, req.Partner)
	if err != nil {
		logger.Errorf("Failed to create partner: %v", err)
		return nil, fmt.Errorf("create partner: %w", err)
	}
	s.metrics.IncPartnerCreated()

	logger.Infof("Successfully onboarded new partner: %s", req.Partner.PartnerId)

	// Emit PartnerOnboardedEvent
	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishPartnerOnboarded(ctx, req.Partner)
	}

	return &partnerservicev1.CreatePartnerResponse{
		PartnerId: req.Partner.PartnerId,
		Partner:   req.Partner,
	}, nil
}

// GetPartner retrieves a specific partner
func (s *PartnerService) GetPartner(ctx context.Context, req *partnerservicev1.GetPartnerRequest) (*partnerservicev1.GetPartnerResponse, error) {
	if req.PartnerId == "" {
		return nil, fmt.Errorf("%w: partner_id is required", ErrInvalidArgument)
	}

	partner, err := s.partnerRepo.GetByID(ctx, req.PartnerId)
	if err != nil {
		logger.Errorf("Partner not found: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: partner not found", ErrNotFound)
		}
		return nil, fmt.Errorf("get partner: %w", err)
	}
	s.metrics.IncPartnerFetched()

	return &partnerservicev1.GetPartnerResponse{
		Partner: partner,
	}, nil
}

// UpdatePartner updates a specific partner's core fields
func (s *PartnerService) UpdatePartner(ctx context.Context, req *partnerservicev1.UpdatePartnerRequest) (*partnerservicev1.UpdatePartnerResponse, error) {
	if req.PartnerId == "" || req.Partner == nil {
		return nil, fmt.Errorf("%w: partner_id and partner details are required", ErrInvalidArgument)
	}
	if err := s.partnerRepo.Update(ctx, req.PartnerId, req.Partner, req.UpdateMask); err != nil {
		logger.Errorf("Partner update failed (partner_id=%s): %v", req.PartnerId, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: partner not found", ErrNotFound)
		}
		return nil, fmt.Errorf("update partner: %w", err)
	}
	s.metrics.IncPartnerUpdated()
	partner, err := s.partnerRepo.GetByID(ctx, req.PartnerId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: updated partner not found", ErrNotFound)
		}
		return nil, fmt.Errorf("load updated partner: %w", err)
	}
	logger.Infof("Partner profile updated (partner_id=%s)", req.PartnerId)

	return &partnerservicev1.UpdatePartnerResponse{
		Partner: partner,
	}, nil
}

// ListPartners retrieves a paginated and filtered list of partners
func (s *PartnerService) ListPartners(ctx context.Context, req *partnerservicev1.ListPartnersRequest) (*partnerservicev1.ListPartnersResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 50
	}
	offset := decodePartnerPageToken(req.PageToken)
	partners, total, err := s.partnerRepo.ListWithFilters(ctx, limit, offset, req.Filter, req.OrderBy)
	if err != nil {
		return nil, fmt.Errorf("list partners: %w", err)
	}
	s.metrics.IncPartnerListed()

	nextToken := ""
	if int32(offset+len(partners)) < total {
		nextToken = strconv.Itoa(offset + len(partners))
	}
	return &partnerservicev1.ListPartnersResponse{
		Partners:      partners,
		TotalCount:    total,
		NextPageToken: nextToken,
	}, nil
}

// DeletePartner soft-deletes a partner
func (s *PartnerService) DeletePartner(ctx context.Context, req *partnerservicev1.DeletePartnerRequest) (*partnerservicev1.DeletePartnerResponse, error) {
	if strings.TrimSpace(req.PartnerId) == "" {
		return nil, fmt.Errorf("%w: partner_id is required", ErrInvalidArgument)
	}
	if err := s.partnerRepo.SoftDelete(ctx, req.PartnerId); err != nil {
		logger.Errorf("Partner delete failed (partner_id=%s): %v", req.PartnerId, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: partner not found", ErrNotFound)
		}
		return nil, fmt.Errorf("delete partner: %w", err)
	}
	s.metrics.IncPartnerDeleted()
	return &partnerservicev1.DeletePartnerResponse{
		Message: "partner deleted successfully",
	}, nil
}

func decodePartnerPageToken(token string) int {
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

// MetricsSnapshot returns current partner runtime counters.
func (s *PartnerService) MetricsSnapshot() map[string]int64 {
	if s.metrics == nil {
		return map[string]int64{}
	}
	return s.metrics.Snapshot()
}
