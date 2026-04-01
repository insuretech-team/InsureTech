package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	policyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/policy/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type InsuranceProposalRepository struct {
	db *gorm.DB
}

func NewInsuranceProposalRepository(db *gorm.DB) *InsuranceProposalRepository {
	return &InsuranceProposalRepository{db: db}
}

func (r *InsuranceProposalRepository) Create(ctx context.Context, proposal *policyv1.InsuranceProposal) (*policyv1.InsuranceProposal, error) {
	if proposal.ProposalId == "" {
		return nil, fmt.Errorf("proposal_id is required")
	}

	proposedPremium, proposedPremiumCurrency := moneyParts(proposal.ProposedPremium)
	proposedSumInsured, proposedSumInsuredCurrency := moneyParts(proposal.ProposedSumInsured)

	submittedAt := time.Now().UTC()
	if proposal.SubmittedAt != nil {
		submittedAt = proposal.SubmittedAt.AsTime()
	}

	reviewedAt := nullableTime(proposal.ReviewedAt)
	reviewedByUserID := nullableString(proposal.ReviewedByUserId)
	approvedPolicyID := nullableString(proposal.ApprovedPolicyId)
	refundID := nullableString(proposal.RefundId)
	submissionPayload := nullableString(proposal.SubmissionPayload)
	insurerResponsePayload := nullableString(proposal.InsurerResponsePayload)
	decisionReason := nullableString(proposal.DecisionReason)
	correlationID := nullableString(proposal.CorrelationId)

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.insurance_proposals
			(proposal_id, proposal_number, tenant_id, order_id, quotation_id, customer_id,
			 insurer_id, product_id, plan_id, proposed_premium, proposed_premium_currency,
			 proposed_sum_insured, proposed_sum_insured_currency, status, submission_payload,
			 insurer_response_payload, decision_reason, submitted_at, reviewed_at,
			 reviewed_by_user_id, approved_policy_id, refund_id, correlation_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
		        $16, $17, $18, $19, $20, $21, $22, $23)`,
		proposal.ProposalId,
		proposal.ProposalNumber,
		proposal.TenantId,
		proposal.OrderId,
		proposal.QuotationId,
		proposal.CustomerId,
		proposal.InsurerId,
		proposal.ProductId,
		proposal.PlanId,
		proposedPremium,
		proposedPremiumCurrency,
		proposedSumInsured,
		proposedSumInsuredCurrency,
		strings.ToUpper(proposal.Status.String()),
		submissionPayload,
		insurerResponsePayload,
		decisionReason,
		submittedAt,
		reviewedAt,
		reviewedByUserID,
		approvedPolicyID,
		refundID,
		correlationID,
	).Error
	if err != nil {
		return nil, fmt.Errorf("failed to insert insurance proposal: %w", err)
	}

	return r.GetByID(ctx, proposal.ProposalId)
}

func (r *InsuranceProposalRepository) GetByID(ctx context.Context, proposalID string) (*policyv1.InsuranceProposal, error) {
	var (
		proposal                 policyv1.InsuranceProposal
		proposedPremium          int64
		proposedPremiumCurrency  string
		proposedSumInsured       int64
		proposedSumCurrency      string
		statusStr                sql.NullString
		submissionPayload        sql.NullString
		insurerResponsePayload   sql.NullString
		decisionReason           sql.NullString
		reviewedByUserID         sql.NullString
		approvedPolicyID         sql.NullString
		refundID                 sql.NullString
		correlationID            sql.NullString
		submittedAt              time.Time
		reviewedAt               sql.NullTime
		createdAt                time.Time
		updatedAt                time.Time
		deletedAt                sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT proposal_id, proposal_number, tenant_id, order_id, quotation_id, customer_id,
		       insurer_id, product_id, plan_id, proposed_premium, proposed_premium_currency,
		       proposed_sum_insured, proposed_sum_insured_currency, status, submission_payload,
		       insurer_response_payload, decision_reason, submitted_at, reviewed_at,
		       reviewed_by_user_id, approved_policy_id, refund_id, correlation_id,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.insurance_proposals
		WHERE proposal_id = $1 AND deleted_at IS NULL`,
		proposalID,
	).Row().Scan(
		&proposal.ProposalId,
		&proposal.ProposalNumber,
		&proposal.TenantId,
		&proposal.OrderId,
		&proposal.QuotationId,
		&proposal.CustomerId,
		&proposal.InsurerId,
		&proposal.ProductId,
		&proposal.PlanId,
		&proposedPremium,
		&proposedPremiumCurrency,
		&proposedSumInsured,
		&proposedSumCurrency,
		&statusStr,
		&submissionPayload,
		&insurerResponsePayload,
		&decisionReason,
		&submittedAt,
		&reviewedAt,
		&reviewedByUserID,
		&approvedPolicyID,
		&refundID,
		&correlationID,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get insurance proposal: %w", err)
	}

	proposal.ProposedPremium = &commonv1.Money{
		Amount:   proposedPremium,
		Currency: proposedPremiumCurrency,
	}
	proposal.ProposedSumInsured = &commonv1.Money{
		Amount:   proposedSumInsured,
		Currency: proposedSumCurrency,
	}
	if statusStr.Valid {
		if v, ok := policyv1.ProposalStatus_value[strings.ToUpper(statusStr.String)]; ok {
			proposal.Status = policyv1.ProposalStatus(v)
		}
	}
	if submissionPayload.Valid {
		proposal.SubmissionPayload = submissionPayload.String
	}
	if insurerResponsePayload.Valid {
		proposal.InsurerResponsePayload = insurerResponsePayload.String
	}
	if decisionReason.Valid {
		proposal.DecisionReason = decisionReason.String
	}
	if reviewedByUserID.Valid {
		proposal.ReviewedByUserId = reviewedByUserID.String
	}
	if approvedPolicyID.Valid {
		proposal.ApprovedPolicyId = approvedPolicyID.String
	}
	if refundID.Valid {
		proposal.RefundId = refundID.String
	}
	if correlationID.Valid {
		proposal.CorrelationId = correlationID.String
	}

	proposal.SubmittedAt = timestamppb.New(submittedAt)
	if reviewedAt.Valid {
		proposal.ReviewedAt = timestamppb.New(reviewedAt.Time)
	}
	proposal.CreatedAt = timestamppb.New(createdAt)
	proposal.UpdatedAt = timestamppb.New(updatedAt)
	if deletedAt.Valid {
		proposal.DeletedAt = timestamppb.New(deletedAt.Time)
	}

	return &proposal, nil
}

func (r *InsuranceProposalRepository) Update(ctx context.Context, proposal *policyv1.InsuranceProposal) (*policyv1.InsuranceProposal, error) {
	proposedPremium, proposedPremiumCurrency := moneyParts(proposal.ProposedPremium)
	proposedSumInsured, proposedSumInsuredCurrency := moneyParts(proposal.ProposedSumInsured)

	submittedAt := time.Now().UTC()
	if proposal.SubmittedAt != nil {
		submittedAt = proposal.SubmittedAt.AsTime()
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.insurance_proposals
		SET proposal_number = $2,
		    tenant_id = $3,
		    order_id = $4,
		    quotation_id = $5,
		    customer_id = $6,
		    insurer_id = $7,
		    product_id = $8,
		    plan_id = $9,
		    proposed_premium = $10,
		    proposed_premium_currency = $11,
		    proposed_sum_insured = $12,
		    proposed_sum_insured_currency = $13,
		    status = $14,
		    submission_payload = $15,
		    insurer_response_payload = $16,
		    decision_reason = $17,
		    submitted_at = $18,
		    reviewed_at = $19,
		    reviewed_by_user_id = $20,
		    approved_policy_id = $21,
		    refund_id = $22,
		    correlation_id = $23
		WHERE proposal_id = $1 AND deleted_at IS NULL`,
		proposal.ProposalId,
		proposal.ProposalNumber,
		proposal.TenantId,
		proposal.OrderId,
		proposal.QuotationId,
		proposal.CustomerId,
		proposal.InsurerId,
		proposal.ProductId,
		proposal.PlanId,
		proposedPremium,
		proposedPremiumCurrency,
		proposedSumInsured,
		proposedSumInsuredCurrency,
		strings.ToUpper(proposal.Status.String()),
		nullableString(proposal.SubmissionPayload),
		nullableString(proposal.InsurerResponsePayload),
		nullableString(proposal.DecisionReason),
		submittedAt,
		nullableTime(proposal.ReviewedAt),
		nullableString(proposal.ReviewedByUserId),
		nullableString(proposal.ApprovedPolicyId),
		nullableString(proposal.RefundId),
		nullableString(proposal.CorrelationId),
	).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update insurance proposal: %w", err)
	}

	return r.GetByID(ctx, proposal.ProposalId)
}

func (r *InsuranceProposalRepository) Delete(ctx context.Context, proposalID string) error {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.insurance_proposals
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE proposal_id = $1 AND deleted_at IS NULL`,
		proposalID,
	).Error
	if err != nil {
		return fmt.Errorf("failed to delete insurance proposal: %w", err)
	}

	return nil
}

func (r *InsuranceProposalRepository) List(ctx context.Context, orderID, insurerID, customerID string, status policyv1.ProposalStatus, page, pageSize int) ([]*policyv1.InsuranceProposal, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	where, args := buildProposalFilters(orderID, insurerID, customerID, status)

	countQuery := `SELECT COUNT(*) FROM insurance_schema.insurance_proposals WHERE deleted_at IS NULL` + where
	var total int64
	if err := r.db.WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count insurance proposals: %w", err)
	}

	query := `
		SELECT proposal_id, proposal_number, tenant_id, order_id, quotation_id, customer_id,
		       insurer_id, product_id, plan_id, proposed_premium, proposed_premium_currency,
		       proposed_sum_insured, proposed_sum_insured_currency, status, submission_payload,
		       insurer_response_payload, decision_reason, submitted_at, reviewed_at,
		       reviewed_by_user_id, approved_policy_id, refund_id, correlation_id,
		       created_at, updated_at, deleted_at
		FROM insurance_schema.insurance_proposals
		WHERE deleted_at IS NULL` + where + `
		ORDER BY created_at DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	queryArgs := append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.WithContext(ctx).Raw(query, queryArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list insurance proposals: %w", err)
	}
	defer rows.Close()

	proposals := make([]*policyv1.InsuranceProposal, 0)
	for rows.Next() {
		proposal, err := scanInsuranceProposal(rows)
		if err != nil {
			return nil, 0, err
		}
		proposals = append(proposals, proposal)
	}

	return proposals, total, nil
}

func buildProposalFilters(orderID, insurerID, customerID string, status policyv1.ProposalStatus) (string, []interface{}) {
	args := make([]interface{}, 0, 4)
	conditions := make([]string, 0, 4)

	if orderID != "" {
		args = append(args, orderID)
		conditions = append(conditions, fmt.Sprintf("order_id = $%d", len(args)))
	}
	if insurerID != "" {
		args = append(args, insurerID)
		conditions = append(conditions, fmt.Sprintf("insurer_id = $%d", len(args)))
	}
	if customerID != "" {
		args = append(args, customerID)
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", len(args)))
	}
	if status != policyv1.ProposalStatus_PROPOSAL_STATUS_UNSPECIFIED {
		args = append(args, strings.ToUpper(status.String()))
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " AND " + strings.Join(conditions, " AND "), args
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanInsuranceProposal(scanner rowScanner) (*policyv1.InsuranceProposal, error) {
	var (
		proposal                policyv1.InsuranceProposal
		proposedPremium         int64
		proposedPremiumCurrency string
		proposedSumInsured      int64
		proposedSumCurrency     string
		statusStr               sql.NullString
		submissionPayload       sql.NullString
		responsePayload         sql.NullString
		decisionReason          sql.NullString
		submittedAt             time.Time
		reviewedAt              sql.NullTime
		reviewedByUserID        sql.NullString
		approvedPolicyID        sql.NullString
		refundID                sql.NullString
		correlationID           sql.NullString
		createdAt               time.Time
		updatedAt               time.Time
		deletedAt               sql.NullTime
	)

	if err := scanner.Scan(
		&proposal.ProposalId,
		&proposal.ProposalNumber,
		&proposal.TenantId,
		&proposal.OrderId,
		&proposal.QuotationId,
		&proposal.CustomerId,
		&proposal.InsurerId,
		&proposal.ProductId,
		&proposal.PlanId,
		&proposedPremium,
		&proposedPremiumCurrency,
		&proposedSumInsured,
		&proposedSumCurrency,
		&statusStr,
		&submissionPayload,
		&responsePayload,
		&decisionReason,
		&submittedAt,
		&reviewedAt,
		&reviewedByUserID,
		&approvedPolicyID,
		&refundID,
		&correlationID,
		&createdAt,
		&updatedAt,
		&deletedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to scan insurance proposal: %w", err)
	}

	proposal.ProposedPremium = &commonv1.Money{
		Amount:   proposedPremium,
		Currency: proposedPremiumCurrency,
	}
	proposal.ProposedSumInsured = &commonv1.Money{
		Amount:   proposedSumInsured,
		Currency: proposedSumCurrency,
	}
	if statusStr.Valid {
		if v, ok := policyv1.ProposalStatus_value[strings.ToUpper(statusStr.String)]; ok {
			proposal.Status = policyv1.ProposalStatus(v)
		}
	}
	if submissionPayload.Valid {
		proposal.SubmissionPayload = submissionPayload.String
	}
	if responsePayload.Valid {
		proposal.InsurerResponsePayload = responsePayload.String
	}
	if decisionReason.Valid {
		proposal.DecisionReason = decisionReason.String
	}
	if reviewedAt.Valid {
		proposal.ReviewedAt = timestamppb.New(reviewedAt.Time)
	}
	if reviewedByUserID.Valid {
		proposal.ReviewedByUserId = reviewedByUserID.String
	}
	if approvedPolicyID.Valid {
		proposal.ApprovedPolicyId = approvedPolicyID.String
	}
	if refundID.Valid {
		proposal.RefundId = refundID.String
	}
	if correlationID.Valid {
		proposal.CorrelationId = correlationID.String
	}

	proposal.SubmittedAt = timestamppb.New(submittedAt)
	proposal.CreatedAt = timestamppb.New(createdAt)
	proposal.UpdatedAt = timestamppb.New(updatedAt)
	if deletedAt.Valid {
		proposal.DeletedAt = timestamppb.New(deletedAt.Time)
	}

	return &proposal, nil
}

func moneyParts(value *commonv1.Money) (int64, string) {
	if value == nil {
		return 0, "BDT"
	}
	if value.Currency == "" {
		return value.Amount, "BDT"
	}
	return value.Amount, value.Currency
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullableTime(value *timestamppb.Timestamp) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: value.AsTime(), Valid: true}
}
