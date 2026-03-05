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

var ErrAlertNotFound = errors.New("fraud alert not found")

// FraudAlertRepository handles CRUD for fraud alerts.
type FraudAlertRepository struct {
	db *gorm.DB
}

func NewFraudAlertRepository(db *gorm.DB) *FraudAlertRepository {
	return &FraudAlertRepository{db: db}
}

func (r *FraudAlertRepository) Create(ctx context.Context, alert *fraudv1.FraudAlert) error {
	if alert.Id == "" {
		alert.Id = uuid.NewString()
	}
	if alert.AlertNumber == "" {
		alert.AlertNumber = "FAL-" + time.Now().UTC().Format("20060102-150405")
	}
	if alert.Status == fraudv1.AlertStatus_ALERT_STATUS_UNSPECIFIED {
		alert.Status = fraudv1.AlertStatus_ALERT_STATUS_OPEN
	}
	if strings.TrimSpace(alert.RiskLevel) == "" {
		alert.RiskLevel = "RISK_LEVEL_MEDIUM"
	}

	now := time.Now().UTC()
	values := map[string]any{
		"alert_id":      alert.Id,
		"alert_number":  alert.AlertNumber,
		"entity_type":   alert.EntityType,
		"entity_id":     alert.EntityId,
		"fraud_rule_id": alert.FraudRuleId,
		"risk_level":    alert.RiskLevel,
		"fraud_score":   alert.FraudScore,
		"details":       alert.Details,
		"status":        alert.Status.String(),
		"assigned_to":   nil,
		"resolved_at":   nil,
		"audit_info":    alert.AuditInfo,
		"created_at":    now,
		"updated_at":    now,
	}
	if alert.AssignedTo != "" {
		values["assigned_to"] = alert.AssignedTo
	}
	if alert.ResolvedAt != nil {
		values["resolved_at"] = alert.ResolvedAt.AsTime()
	}

	return r.db.WithContext(ctx).Table("insurance_schema.fraud_alerts").Create(values).Error
}

func (r *FraudAlertRepository) GetByID(ctx context.Context, alertID string) (*fraudv1.FraudAlert, error) {
	var alert fraudv1.FraudAlert
	err := r.db.WithContext(ctx).Table("insurance_schema.fraud_alerts").
		Where("alert_id = ?", alertID).
		First(&alert).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAlertNotFound
		}
		return nil, err
	}
	return &alert, nil
}

func (r *FraudAlertRepository) List(ctx context.Context, status string, riskLevel string, start, end *time.Time, limit, offset int) ([]*fraudv1.FraudAlert, int32, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	q := r.db.WithContext(ctx).Table("insurance_schema.fraud_alerts")
	if s := normalizeAlertStatus(status); s != "" {
		q = q.Where("status = ?", s)
	}
	if rl := strings.TrimSpace(riskLevel); rl != "" {
		q = q.Where("risk_level = ?", rl)
	}
	if start != nil {
		q = q.Where("created_at >= ?", *start)
	}
	if end != nil {
		q = q.Where("created_at <= ?", *end)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var alerts []*fraudv1.FraudAlert
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&alerts).Error; err != nil {
		return nil, 0, err
	}
	return alerts, int32(total), nil
}

func (r *FraudAlertRepository) UpdateStatus(ctx context.Context, alertID string, status fraudv1.AlertStatus, assignedTo string) error {
	updates := map[string]any{
		"status":     status.String(),
		"updated_at": time.Now().UTC(),
	}
	if strings.TrimSpace(assignedTo) != "" {
		updates["assigned_to"] = strings.TrimSpace(assignedTo)
	}
	if status == fraudv1.AlertStatus_ALERT_STATUS_CLOSED ||
		status == fraudv1.AlertStatus_ALERT_STATUS_FALSE_POSITIVE ||
		status == fraudv1.AlertStatus_ALERT_STATUS_CONFIRMED {
		updates["resolved_at"] = time.Now().UTC()
	}

	res := r.db.WithContext(ctx).Table("insurance_schema.fraud_alerts").
		Where("alert_id = ?", alertID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAlertNotFound
	}
	return nil
}

func normalizeAlertStatus(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if _, ok := fraudv1.AlertStatus_value[s]; ok {
		return s
	}
	upper := strings.ToUpper(s)
	if _, ok := fraudv1.AlertStatus_value["ALERT_STATUS_"+upper]; ok {
		return "ALERT_STATUS_" + upper
	}
	return s
}
