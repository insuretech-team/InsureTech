package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/domain"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type CommissionRepository struct {
	db *gorm.DB
}

func NewCommissionRepository(db *gorm.DB) *CommissionRepository {
	return &CommissionRepository{db: db}
}

// Create inserts a newly generated commission entry
func (r *CommissionRepository) Create(ctx context.Context, comm *partnerv1.Commission) error {
	if comm.CommissionId == "" {
		comm.CommissionId = uuid.New().String()
	}
	now := time.Now()
	comm.CreatedAt = timestamppb.New(now)
	comm.UpdatedAt = timestamppb.New(now)

	if comm.Status == partnerv1.CommissionStatus_COMMISSION_STATUS_UNSPECIFIED {
		comm.Status = partnerv1.CommissionStatus_COMMISSION_STATUS_PENDING
	}

	values := map[string]any{
		"commission_id":     comm.CommissionId,
		"policy_id":         comm.PolicyId,
		"partner_id":        nil,
		"agent_id":          nil,
		"payment_id":        nil,
		"type":              comm.Type.String(),
		"commission_amount": comm.CommissionAmount.Amount, // Stored as BIGINT in paisa
		"commission_rate":   comm.CommissionRate,
		"status":            comm.Status.String(),
		"created_at":        now,
		"updated_at":        now,
	}

	if comm.PartnerId != "" {
		values["partner_id"] = comm.PartnerId
	}
	if comm.AgentId != "" {
		values["agent_id"] = comm.AgentId
	}
	if comm.PaymentId != "" {
		values["payment_id"] = comm.PaymentId
	}

	return r.db.WithContext(ctx).Table("partner_schema.commissions").Create(values).Error
}

// GetByID looks up a single commission record
func (r *CommissionRepository) GetByID(ctx context.Context, id string) (*partnerv1.Commission, error) {
	var comm partnerv1.Commission
	err := r.db.WithContext(ctx).Table("partner_schema.commissions").
		Where("commission_id = ?", id).
		First(&comm).Error
	if err != nil {
		return nil, err
	}
	return &comm, nil
}

// ListPendingByPartner gets outstanding commissions ready for payout compilation
func (r *CommissionRepository) ListPendingByPartner(ctx context.Context, partnerID string, limit int) ([]*partnerv1.Commission, error) {
	var comms []*partnerv1.Commission
	err := r.db.WithContext(ctx).Table("partner_schema.commissions").
		Where("partner_id = ? AND status = ?", partnerID, partnerv1.CommissionStatus_COMMISSION_STATUS_PENDING.String()).
		Limit(limit).
		Find(&comms).Error
	return comms, err
}

// MarkAsPaid links a payment ID to the commission and sets it as PAID
func (r *CommissionRepository) MarkAsPaid(ctx context.Context, commissionID, paymentID string) error {
	return r.db.WithContext(ctx).Table("partner_schema.commissions").
		Where("commission_id = ?", commissionID).
		Updates(map[string]interface{}{
			"status":     partnerv1.CommissionStatus_COMMISSION_STATUS_PAID.String(),
			"payment_id": paymentID,
			"paid_at":    time.Now(),
			"updated_at": time.Now(),
		}).Error
}

// ListByPartnerAndDateRange returns commissions for a partner within optional date range.
func (r *CommissionRepository) ListByPartnerAndDateRange(
	ctx context.Context,
	partnerID string,
	start, end *time.Time,
	limit, offset int,
) ([]*partnerv1.Commission, int32, error) {
	var (
		comms []*partnerv1.Commission
		total int64
	)
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}
	q := r.db.WithContext(ctx).Table("partner_schema.commissions").Where("partner_id = ?", partnerID)
	if start != nil {
		q = q.Where("created_at >= ?", *start)
	}
	if end != nil {
		q = q.Where("created_at <= ?", *end)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&comms).Error; err != nil {
		return nil, 0, err
	}
	return comms, int32(total), nil
}

// SumByPartnerAndDateRange returns the total commission amount in paisa.
func (r *CommissionRepository) SumByPartnerAndDateRange(ctx context.Context, partnerID string, start, end *time.Time) (int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Table("partner_schema.commissions").Where("partner_id = ?", partnerID)
	if start != nil {
		q = q.Where("created_at >= ?", *start)
	}
	if end != nil {
		q = q.Where("created_at <= ?", *end)
	}
	if err := q.Select("COALESCE(SUM(commission_amount),0)").Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// ExistsByPolicyAndType checks whether a commission already exists for policy+type.
func (r *CommissionRepository) ExistsByPolicyAndType(ctx context.Context, policyID string, cType partnerv1.CommissionType) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("partner_schema.commissions").
		Where("policy_id = ? AND type = ?", policyID, cType.String()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type policyCommissionInputRow struct {
	PolicyID      string `gorm:"column:policy_id"`
	PartnerID     string `gorm:"column:partner_id"`
	AgentID       string `gorm:"column:agent_id"`
	PremiumAmount int64  `gorm:"column:premium_amount"`
	Currency      string `gorm:"column:premium_currency"`
}

// ResolvePolicyCommissionInput loads policy fields needed to compute commission.
func (r *CommissionRepository) ResolvePolicyCommissionInput(ctx context.Context, policyID string) (*domain.PolicyCommissionInput, error) {
	var row policyCommissionInputRow
	err := r.db.WithContext(ctx).Table("insurance_schema.policies").
		Select("policy_id, partner_id, agent_id, premium_amount, COALESCE(premium_currency, 'BDT') AS premium_currency").
		Where("policy_id = ?", policyID).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &domain.PolicyCommissionInput{
		PolicyID:      row.PolicyID,
		PartnerID:     row.PartnerID,
		AgentID:       row.AgentID,
		PremiumAmount: row.PremiumAmount,
		Currency:      row.Currency,
	}, nil
}
