package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	claimsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/claims/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type ClaimApprovalRepository struct {
	db *gorm.DB
}

func NewClaimApprovalRepository(db *gorm.DB) *ClaimApprovalRepository {
	return &ClaimApprovalRepository{db: db}
}

func (r *ClaimApprovalRepository) Create(ctx context.Context, approval *claimsv1.ClaimApproval) (*claimsv1.ClaimApproval, error) {
	if approval.ApprovalId == "" {
		return nil, fmt.Errorf("approval_id is required")
	}
	
	approvedAmount := int64(0)
	approvedCurrency := "BDT"
	if approval.ApprovedAmount != nil {
		approvedAmount = approval.ApprovedAmount.Amount
		approvedCurrency = approval.ApprovedAmount.Currency
	}
	
	var approvedAt interface{}
	if approval.ApprovedAt != nil {
		approvedAt = approval.ApprovedAt.AsTime()
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.claim_approvals
			(approval_id, claim_id, approver_id, approver_role, approval_level,
			 decision, approved_amount, approved_currency, notes, approved_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())`,
		approval.ApprovalId, approval.ClaimId, approval.ApproverId, approval.ApproverRole,
		approval.ApprovalLevel, strings.ToUpper(approval.Decision.String()),
		approvedAmount, approvedCurrency, approval.Notes, approvedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert claim approval: %w", err)
	}

	return r.GetByID(ctx, approval.ApprovalId)
}

func (r *ClaimApprovalRepository) GetByID(ctx context.Context, approvalID string) (*claimsv1.ClaimApproval, error) {
	var (
		a                claimsv1.ClaimApproval
		decisionStr      sql.NullString
		approvedAmount   sql.NullInt64
		approvedCurrency string
		notes            sql.NullString
		approvedAt       sql.NullTime
		createdAt        time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT approval_id, claim_id, approver_id, approver_role, approval_level,
		       decision, approved_amount, approved_currency, notes, approved_at, created_at
		FROM insurance_schema.claim_approvals
		WHERE approval_id = $1`,
		approvalID,
	).Row().Scan(
		&a.ApprovalId, &a.ClaimId, &a.ApproverId, &a.ApproverRole, &a.ApprovalLevel,
		&decisionStr, &approvedAmount, &approvedCurrency, &notes, &approvedAt, &createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get claim approval: %w", err)
	}

	if decisionStr.Valid {
		k := strings.ToUpper(decisionStr.String)
		if v, ok := claimsv1.ApprovalDecision_value[k]; ok {
			a.Decision = claimsv1.ApprovalDecision(v)
		}
	}
	
	if approvedAmount.Valid {
		a.ApprovedAmount = &commonv1.Money{Amount: approvedAmount.Int64, Currency: approvedCurrency}
	}
	a.ApprovedCurrency = approvedCurrency
	
	if notes.Valid {
		a.Notes = notes.String
	}
	if approvedAt.Valid {
		a.ApprovedAt = timestamppb.New(approvedAt.Time)
	}
	if !createdAt.IsZero() {
		a.CreatedAt = timestamppb.New(createdAt)
	}

	return &a, nil
}

func (r *ClaimApprovalRepository) Update(ctx context.Context, approval *claimsv1.ClaimApproval) (*claimsv1.ClaimApproval, error) {
	approvedAmount := int64(0)
	approvedCurrency := "BDT"
	if approval.ApprovedAmount != nil {
		approvedAmount = approval.ApprovedAmount.Amount
		approvedCurrency = approval.ApprovedAmount.Currency
	}
	
	var approvedAt interface{}
	if approval.ApprovedAt != nil {
		approvedAt = approval.ApprovedAt.AsTime()
	}
	
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.claim_approvals
		SET claim_id = $2, approver_id = $3, approver_role = $4, approval_level = $5,
		    decision = $6, approved_amount = $7, approved_currency = $8, notes = $9, approved_at = $10
		WHERE approval_id = $1`,
		approval.ApprovalId, approval.ClaimId, approval.ApproverId, approval.ApproverRole,
		approval.ApprovalLevel, strings.ToUpper(approval.Decision.String()),
		approvedAmount, approvedCurrency, approval.Notes, approvedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update claim approval: %w", err)
	}

	return r.GetByID(ctx, approval.ApprovalId)
}

func (r *ClaimApprovalRepository) Delete(ctx context.Context, approvalID string) error {
	return r.db.WithContext(ctx).Exec(`DELETE FROM insurance_schema.claim_approvals WHERE approval_id = $1`, approvalID).Error
}

func (r *ClaimApprovalRepository) ListByClaimID(ctx context.Context, claimID string) ([]*claimsv1.ClaimApproval, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT approval_id, claim_id, approver_id, approver_role, approval_level,
		       decision, approved_amount, approved_currency, notes, approved_at, created_at
		FROM insurance_schema.claim_approvals
		WHERE claim_id = $1
		ORDER BY approval_level ASC, created_at DESC`,
		claimID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list claim approvals: %w", err)
	}
	defer rows.Close()

	approvals := make([]*claimsv1.ClaimApproval, 0)
	for rows.Next() {
		var (
			a                claimsv1.ClaimApproval
			decisionStr      sql.NullString
			approvedAmount   sql.NullInt64
			approvedCurrency string
			notes            sql.NullString
			approvedAt       sql.NullTime
			createdAt        time.Time
		)

		err := rows.Scan(
			&a.ApprovalId, &a.ClaimId, &a.ApproverId, &a.ApproverRole, &a.ApprovalLevel,
			&decisionStr, &approvedAmount, &approvedCurrency, &notes, &approvedAt, &createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan claim approval: %w", err)
		}

		if decisionStr.Valid {
			k := strings.ToUpper(decisionStr.String)
			if v, ok := claimsv1.ApprovalDecision_value[k]; ok {
				a.Decision = claimsv1.ApprovalDecision(v)
			}
		}
		
		if approvedAmount.Valid {
			a.ApprovedAmount = &commonv1.Money{Amount: approvedAmount.Int64, Currency: approvedCurrency}
		}
		a.ApprovedCurrency = approvedCurrency
		
		if notes.Valid {
			a.Notes = notes.String
		}
		if approvedAt.Valid {
			a.ApprovedAt = timestamppb.New(approvedAt.Time)
		}
		if !createdAt.IsZero() {
			a.CreatedAt = timestamppb.New(createdAt)
		}

		approvals = append(approvals, &a)
	}

	return approvals, nil
}
