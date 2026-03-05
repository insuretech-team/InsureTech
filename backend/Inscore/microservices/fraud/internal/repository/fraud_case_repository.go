package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	"gorm.io/gorm"
)

var ErrCaseNotFound = errors.New("fraud case not found")

// FraudCaseRepository handles CRUD for fraud cases.
type FraudCaseRepository struct {
	db *gorm.DB
}

func NewFraudCaseRepository(db *gorm.DB) *FraudCaseRepository {
	return &FraudCaseRepository{db: db}
}

func (r *FraudCaseRepository) Create(ctx context.Context, fraudCase *fraudv1.FraudCase) error {
	if fraudCase.Id == "" {
		fraudCase.Id = uuid.NewString()
	}
	if fraudCase.CaseNumber == "" {
		fraudCase.CaseNumber = "FRC-" + time.Now().UTC().Format("20060102-150405")
	}
	if fraudCase.Priority == fraudv1.CasePriority_CASE_PRIORITY_UNSPECIFIED {
		fraudCase.Priority = fraudv1.CasePriority_CASE_PRIORITY_MEDIUM
	}
	if fraudCase.Status == fraudv1.CaseStatus_CASE_STATUS_UNSPECIFIED {
		fraudCase.Status = fraudv1.CaseStatus_CASE_STATUS_OPEN
	}

	now := time.Now().UTC()
	values := map[string]any{
		"case_id":             fraudCase.Id,
		"case_number":         fraudCase.CaseNumber,
		"fraud_alert_id":      fraudCase.FraudAlertId,
		"priority":            fraudCase.Priority.String(),
		"investigation_notes": fraudCase.InvestigationNotes,
		"evidence":            fraudCase.Evidence,
		"status":              fraudCase.Status.String(),
		"outcome":             nil,
		"investigator_id":     nil,
		"closed_at":           nil,
		"audit_info":          fraudCase.AuditInfo,
		"created_at":          now,
		"updated_at":          now,
	}
	if fraudCase.Outcome != fraudv1.CaseOutcome_CASE_OUTCOME_UNSPECIFIED {
		values["outcome"] = fraudCase.Outcome.String()
	}
	if fraudCase.InvestigatorId != "" {
		values["investigator_id"] = fraudCase.InvestigatorId
	}
	if fraudCase.ClosedAt != nil {
		values["closed_at"] = fraudCase.ClosedAt.AsTime()
	}

	return r.db.WithContext(ctx).Table("insurance_schema.fraud_cases").Create(values).Error
}

func (r *FraudCaseRepository) GetByID(ctx context.Context, caseID string) (*fraudv1.FraudCase, error) {
	var fraudCase fraudv1.FraudCase
	err := r.db.WithContext(ctx).Table("insurance_schema.fraud_cases").
		Where("case_id = ?", caseID).
		First(&fraudCase).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCaseNotFound
		}
		return nil, err
	}
	return &fraudCase, nil
}

func (r *FraudCaseRepository) Update(ctx context.Context, caseID string, status fraudv1.CaseStatus, outcome fraudv1.CaseOutcome, notes string, evidence string) error {
	updates := map[string]any{
		"updated_at": time.Now().UTC(),
	}
	if status != fraudv1.CaseStatus_CASE_STATUS_UNSPECIFIED {
		updates["status"] = status.String()
	}
	if outcome != fraudv1.CaseOutcome_CASE_OUTCOME_UNSPECIFIED {
		updates["outcome"] = outcome.String()
	}
	if strings.TrimSpace(notes) != "" {
		updates["investigation_notes"] = notes
	}
	if strings.TrimSpace(evidence) != "" {
		updates["evidence"] = evidence
	}
	if status == fraudv1.CaseStatus_CASE_STATUS_CLOSED {
		updates["closed_at"] = time.Now().UTC()
	}

	res := r.db.WithContext(ctx).Table("insurance_schema.fraud_cases").
		Where("case_id = ?", caseID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCaseNotFound
	}
	return nil
}
