package domain

import (
	"context"
	"time"

	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
)

// RuleRepository defines fraud rule persistence operations.
type RuleRepository interface {
	Create(ctx context.Context, rule *fraudv1.FraudRule) error
	GetByID(ctx context.Context, ruleID string) (*fraudv1.FraudRule, error)
	Update(ctx context.Context, ruleID string, rule *fraudv1.FraudRule) error
	List(ctx context.Context, category fraudv1.RuleCategory, activeOnly bool, limit, offset int) ([]*fraudv1.FraudRule, int32, error)
	SetActive(ctx context.Context, ruleID string, active bool) error
}

// AlertRepository defines fraud alert persistence operations.
type AlertRepository interface {
	Create(ctx context.Context, alert *fraudv1.FraudAlert) error
	GetByID(ctx context.Context, alertID string) (*fraudv1.FraudAlert, error)
	List(ctx context.Context, status string, riskLevel string, start, end *time.Time, limit, offset int) ([]*fraudv1.FraudAlert, int32, error)
	UpdateStatus(ctx context.Context, alertID string, status fraudv1.AlertStatus, assignedTo string) error
}

// CaseRepository defines fraud case persistence operations.
type CaseRepository interface {
	Create(ctx context.Context, fraudCase *fraudv1.FraudCase) error
	GetByID(ctx context.Context, caseID string) (*fraudv1.FraudCase, error)
	Update(ctx context.Context, caseID string, status fraudv1.CaseStatus, outcome fraudv1.CaseOutcome, notes string, evidence string) error
}
