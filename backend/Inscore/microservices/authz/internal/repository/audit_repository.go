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

// AuditRepo writes access decision audit rows using the proto-generated struct directly.
type AuditRepo struct{ db *gorm.DB }

func NewAuditRepo(db *gorm.DB) *AuditRepo { return &AuditRepo{db: db} }

func (r *AuditRepo) Create(ctx context.Context, a *authzentityv1.AccessDecisionAudit) error {
	// Set DecidedAt to now if not provided
	if a.DecidedAt == nil {
		a.DecidedAt = timestamppb.New(time.Now())
	}
	return r.db.WithContext(ctx).Table("authz_schema.access_decision_audits").Create(a).Error
}

func (r *AuditRepo) List(ctx context.Context, req *authzservicev1.ListAccessDecisionAuditsRequest) ([]*authzentityv1.AccessDecisionAudit, int64, error) {
	q := r.db.WithContext(ctx).Table("authz_schema.access_decision_audits")

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

	var rows []*authzentityv1.AccessDecisionAudit
	if err := q.Order("decided_at DESC").Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
