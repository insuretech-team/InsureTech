package repository

import (
	"context"
	"strings"
	"time"

	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// accessDecisionAuditRow is a plain Go struct for GORM scanning.
// The proto-generated AccessDecisionAudit has a *timestamppb.Timestamp DecidedAt
// field that GORM cannot reflect — so we use this intermediate struct instead.
type accessDecisionAuditRow struct {
	AuditID     string     `gorm:"column:audit_id"`
	UserID      string     `gorm:"column:user_id"`
	SessionID   string     `gorm:"column:session_id"`
	Domain      string     `gorm:"column:domain"`
	Subject     string     `gorm:"column:subject"`
	Object      string     `gorm:"column:object"`
	Action      string     `gorm:"column:action"`
	Decision    string     `gorm:"column:decision"`
	MatchedRule string     `gorm:"column:matched_rule"`
	IPAddress   string     `gorm:"column:ip_address"`
	UserAgent   string     `gorm:"column:user_agent"`
	DecidedAt   *time.Time `gorm:"column:decided_at"`
}

func (accessDecisionAuditRow) TableName() string {
	return "authz_schema.access_decision_audits"
}

func auditRowToProto(r *accessDecisionAuditRow) *authzentityv1.AccessDecisionAudit {
	a := &authzentityv1.AccessDecisionAudit{
		AuditId:     r.AuditID,
		UserId:      r.UserID,
		SessionId:   r.SessionID,
		Domain:      r.Domain,
		Subject:     r.Subject,
		Object:      r.Object,
		Action:      r.Action,
		MatchedRule: r.MatchedRule,
		IpAddress:   r.IPAddress,
		UserAgent:   r.UserAgent,
	}
	if r.DecidedAt != nil {
		a.DecidedAt = timestamppb.New(*r.DecidedAt)
	}
	switch strings.ToLower(r.Decision) {
	case "allow", "policy_effect_allow":
		a.Decision = authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW
	case "deny", "policy_effect_deny":
		a.Decision = authzentityv1.PolicyEffect_POLICY_EFFECT_DENY
	}
	return a
}

// AuditRepo writes access decision audit rows.
type AuditRepo struct{ db *gorm.DB }

func NewAuditRepo(db *gorm.DB) *AuditRepo { return &AuditRepo{db: db} }

func (r *AuditRepo) Create(ctx context.Context, a *authzentityv1.AccessDecisionAudit) error {
	decidedAt := time.Now()
	if a.DecidedAt != nil {
		decidedAt = a.DecidedAt.AsTime()
	}
	// Map PolicyEffect enum to DB string
	decision := "deny"
	if a.Decision == authzentityv1.PolicyEffect_POLICY_EFFECT_ALLOW {
		decision = "allow"
	}
	row := &accessDecisionAuditRow{
		AuditID:     a.AuditId,
		UserID:      a.UserId,
		SessionID:   a.SessionId,
		Domain:      a.Domain,
		Subject:     a.Subject,
		Object:      a.Object,
		Action:      a.Action,
		Decision:    decision,
		MatchedRule: a.MatchedRule,
		IPAddress:   a.IpAddress,
		UserAgent:   a.UserAgent,
		DecidedAt:   &decidedAt,
	}
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *AuditRepo) List(ctx context.Context, req *authzservicev1.ListAccessDecisionAuditsRequest) ([]*authzentityv1.AccessDecisionAudit, int64, error) {
	q := r.db.WithContext(ctx).Model(&accessDecisionAuditRow{})

	if req.UserId != "" {
		q = q.Where("user_id = ?", req.UserId)
	}
	if req.Domain != "" {
		q = q.Where("domain = ?", req.Domain)
	}
	if req.Decision != authzentityv1.PolicyEffect_POLICY_EFFECT_UNSPECIFIED {
		decisionName := strings.ToLower(req.Decision.String())
		decisionShort := strings.TrimPrefix(decisionName, "policy_effect_")
		q = q.Where(
			"(LOWER(decision::text) = ? OR LOWER(decision::text) = ? OR LOWER(decision::text) LIKE ?)",
			decisionName, decisionShort, "%"+decisionShort+"%",
		)
	}
	if req.From != nil {
		q = q.Where("decided_at >= ?", req.From.AsTime())
	}
	if req.To != nil {
		q = q.Where("decided_at <= ?", req.To.AsTime())
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 || pageSize > 500 {
		pageSize = 50
	}

	var rows []*accessDecisionAuditRow
	if err := q.Order("decided_at DESC").Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	audits := make([]*authzentityv1.AccessDecisionAudit, len(rows))
	for i, row := range rows {
		audits[i] = auditRowToProto(row)
	}
	return audits, total, nil
}
