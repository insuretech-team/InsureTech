package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	"gorm.io/gorm"
)

// ProcessPolicyCommissionEvent computes and persists commission for a policy lifecycle event.
func (s *PartnerService) ProcessPolicyCommissionEvent(ctx context.Context, policyID string, cType partnerv1.CommissionType) error {
	policyID = strings.TrimSpace(policyID)
	if policyID == "" {
		return fmt.Errorf("%w: policy_id is required", ErrInvalidArgument)
	}
	if cType == partnerv1.CommissionType_COMMISSION_TYPE_UNSPECIFIED {
		cType = partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION
	}

	exists, err := s.commissionRepo.ExistsByPolicyAndType(ctx, policyID, cType)
	if err != nil {
		return fmt.Errorf("check existing commission: %w", err)
	}
	if exists {
		logger.Infof("Commission already exists (policy_id=%s, type=%s); skipping duplicate", policyID, cType.String())
		return nil
	}

	input, err := s.commissionRepo.ResolvePolicyCommissionInput(ctx, policyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: policy not found", ErrNotFound)
		}
		return fmt.Errorf("resolve policy commission input: %w", err)
	}

	if strings.TrimSpace(input.PartnerID) == "" {
		// Policy has no partner context; no partner commission to calculate.
		return nil
	}

	partner, err := s.partnerRepo.GetByID(ctx, input.PartnerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: partner not found", ErrNotFound)
		}
		return fmt.Errorf("load partner for commission: %w", err)
	}

	rate := commissionRateByType(partner, cType)
	if rate <= 0 {
		// Zero-rate partner contract => no payable commission.
		return nil
	}

	amount := int64(math.Round(float64(input.PremiumAmount) * rate / 100.0))
	if amount <= 0 {
		return nil
	}

	currency := strings.TrimSpace(input.Currency)
	if currency == "" {
		currency = "BDT"
	}

	commission := &partnerv1.Commission{
		PolicyId:  input.PolicyID,
		PartnerId: input.PartnerID,
		AgentId:   input.AgentID,
		Type:      cType,
		CommissionAmount: &commonv1.Money{
			Amount:        amount,
			Currency:      currency,
			DecimalAmount: float64(amount) / 100.0,
		},
		CommissionRate: rate,
		Status:         partnerv1.CommissionStatus_COMMISSION_STATUS_PENDING,
	}

	if err := s.commissionRepo.Create(ctx, commission); err != nil {
		return fmt.Errorf("create commission: %w", err)
	}

	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishCommissionCalculated(ctx, commission)
	}
	s.metrics.IncCommissionCalculated()
	return nil
}

func commissionRateByType(partner *partnerv1.Partner, cType partnerv1.CommissionType) float64 {
	if partner == nil {
		return 0
	}
	switch cType {
	case partnerv1.CommissionType_COMMISSION_TYPE_RENEWAL:
		return partner.RenewalCommissionRate
	case partnerv1.CommissionType_COMMISSION_TYPE_CLAIMS_ASSISTANCE:
		return partner.ClaimsAssistanceRate
	default:
		return partner.AcquisitionCommissionRate
	}
}
