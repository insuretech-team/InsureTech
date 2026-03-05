package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetPartnerCommission aggregates commission earned over a period
func (s *PartnerService) GetPartnerCommission(ctx context.Context, req *partnerservicev1.GetPartnerCommissionRequest) (*partnerservicev1.GetPartnerCommissionResponse, error) {
	if req.PartnerId == "" {
		return nil, fmt.Errorf("%w: partner_id is required", ErrInvalidArgument)
	}
	start, end := toTimeRange(req.StartDate, req.EndDate)
	comms, _, err := s.commissionRepo.ListByPartnerAndDateRange(ctx, req.PartnerId, start, end, 200, 0)
	if err != nil {
		logger.Errorf("Failed to list partner commissions (partner_id=%s): %v", req.PartnerId, err)
		return nil, fmt.Errorf("load commission details: %w", err)
	}
	totalAmount, err := s.commissionRepo.SumByPartnerAndDateRange(ctx, req.PartnerId, start, end)
	if err != nil {
		logger.Errorf("Failed to sum partner commissions (partner_id=%s): %v", req.PartnerId, err)
		return nil, fmt.Errorf("aggregate commission total: %w", err)
	}

	logger.Infof("Fetching commission summary for partner: %s", req.PartnerId)
	s.metrics.IncCommissionCalculated()

	details := make([]*partnerservicev1.CommissionDetail, 0, len(comms))
	for _, c := range comms {
		amount := int64(0)
		currency := "BDT"
		if c.CommissionAmount != nil {
			amount = c.CommissionAmount.Amount
			if c.CommissionAmount.Currency != "" {
				currency = c.CommissionAmount.Currency
			}
		}
		details = append(details, &partnerservicev1.CommissionDetail{
			PolicyId: c.PolicyId,
			CommissionAmount: &commonv1.Money{
				Amount:        amount,
				Currency:      currency,
				DecimalAmount: float64(amount) / 100.0,
			},
			CommissionType: strings.TrimPrefix(c.Type.String(), "COMMISSION_TYPE_"),
			EarnedAt:       c.CreatedAt,
		})
	}

	return &partnerservicev1.GetPartnerCommissionResponse{
		PartnerId: req.PartnerId,
		Currency:  "BDT",
		TotalCommission: &commonv1.Money{
			Amount:        totalAmount,
			Currency:      "BDT",
			DecimalAmount: float64(totalAmount) / 100.0,
		},
		Details: details,
	}, nil
}

// UpdateCommissionStructure adjusts existing rates (requires focal person or admin authz)
func (s *PartnerService) UpdateCommissionStructure(ctx context.Context, req *partnerservicev1.UpdateCommissionStructureRequest) (*partnerservicev1.UpdateCommissionStructureResponse, error) {
	if strings.TrimSpace(req.PartnerId) == "" {
		return nil, fmt.Errorf("%w: partner_id is required", ErrInvalidArgument)
	}
	partner, err := s.partnerRepo.GetByID(ctx, req.PartnerId)
	if err != nil {
		return nil, fmt.Errorf("%w: partner not found", ErrNotFound)
	}
	applyCommissionRates(partner, req.CommissionRates)
	if err := s.partnerRepo.Update(ctx, req.PartnerId, partner, []string{
		"acquisition_commission_rate",
		"renewal_commission_rate",
		"claims_assistance_rate",
	}); err != nil {
		logger.Errorf("Failed to update commission structure (partner_id=%s): %v", req.PartnerId, err)
		return nil, fmt.Errorf("update commission structure: %w", err)
	}
	return &partnerservicev1.UpdateCommissionStructureResponse{
		Success: true,
		Message: "commission structure updated",
	}, nil
}

func toTimeRange(startTS, endTS *timestamppb.Timestamp) (*time.Time, *time.Time) {
	var start *time.Time
	var end *time.Time
	if startTS != nil && startTS.IsValid() {
		v := startTS.AsTime()
		start = &v
	}
	if endTS != nil && endTS.IsValid() {
		v := endTS.AsTime()
		end = &v
	}
	return start, end
}

func applyCommissionRates(partner *partnerv1.Partner, rates map[string]float64) {
	// Keep existing rates unless explicitly provided.
	acq := partner.AcquisitionCommissionRate
	ren := partner.RenewalCommissionRate
	claim := partner.ClaimsAssistanceRate

	for k, v := range rates {
		switch strings.ToLower(strings.TrimSpace(k)) {
		case "acquisition", "acquisition_rate", "acquisition_commission_rate":
			acq = v
		case "renewal", "renewal_rate", "renewal_commission_rate":
			ren = v
		case "claims_assistance", "claims_assistance_rate":
			claim = v
		}
	}
	partner.AcquisitionCommissionRate = acq
	partner.RenewalCommissionRate = ren
	partner.ClaimsAssistanceRate = claim
}
