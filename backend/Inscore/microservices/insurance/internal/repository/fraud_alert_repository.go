package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
)

type FraudAlertRepository struct {
	db *gorm.DB
}

func NewFraudAlertRepository(db *gorm.DB) *FraudAlertRepository {
	return &FraudAlertRepository{db: db}
}

func (r *FraudAlertRepository) Create(ctx context.Context, alert *fraudv1.FraudAlert) (*fraudv1.FraudAlert, error) {
	if alert.Id == "" {
		return nil, fmt.Errorf("alert_id is required")
	}

	var createdBy string
	err := r.db.WithContext(ctx).Raw(`SELECT user_id FROM authn_schema.users LIMIT 1`).Scan(&createdBy).Error
	if err != nil || createdBy == "" {
		return nil, fmt.Errorf("failed to get valid user for created_by: %w", err)
	}

	auditInfoJSON := fmt.Sprintf(`{"created_by":"%s","created_at":"%s"}`, createdBy, time.Now().UTC().Format(time.RFC3339))

	var assignedTo sql.NullString
	if alert.AssignedTo != "" {
		assignedTo = sql.NullString{String: alert.AssignedTo, Valid: true}
	}

	var resolvedAt sql.NullTime
	if alert.ResolvedAt != nil {
		resolvedAt = sql.NullTime{Time: alert.ResolvedAt.AsTime(), Valid: true}
	}

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.fraud_alerts
			(alert_id, alert_number, entity_type, entity_id, fraud_rule_id, 
			 risk_level, fraud_score, details, status, assigned_to, resolved_at, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		alert.Id,
		alert.AlertNumber,
		alert.EntityType,
		alert.EntityId,
		alert.FraudRuleId,
		alert.RiskLevel,
		alert.FraudScore,
		alert.Details,
		strings.ToUpper(alert.Status.String()),
		assignedTo,
		resolvedAt,
		auditInfoJSON,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert fraud alert: %w", err)
	}

	return r.GetByID(ctx, alert.Id)
}

func (r *FraudAlertRepository) GetByID(ctx context.Context, alertID string) (*fraudv1.FraudAlert, error) {
	var (
		alert      fraudv1.FraudAlert
		statusStr  sql.NullString
		details    sql.NullString
		assignedTo sql.NullString
		resolvedAt sql.NullTime
		auditInfo  sql.NullString
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT alert_id, alert_number, entity_type, entity_id, fraud_rule_id, 
		       risk_level, fraud_score, details, status, assigned_to, resolved_at,
		       audit_info, created_at, updated_at
		FROM insurance_schema.fraud_alerts
		WHERE alert_id = $1`,
		alertID,
	).Row().Scan(
		&alert.Id,
		&alert.AlertNumber,
		&alert.EntityType,
		&alert.EntityId,
		&alert.FraudRuleId,
		&alert.RiskLevel,
		&alert.FraudScore,
		&details,
		&statusStr,
		&assignedTo,
		&resolvedAt,
		&auditInfo,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get fraud alert: %w", err)
	}

	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := fraudv1.AlertStatus_value[k]; ok {
			alert.Status = fraudv1.AlertStatus(v)
		}
	}

	if details.Valid {
		alert.Details = details.String
	}

	if assignedTo.Valid {
		alert.AssignedTo = assignedTo.String
	}

	if resolvedAt.Valid {
		alert.ResolvedAt = timestamppb.New(resolvedAt.Time)
	}

	return &alert, nil
}

func (r *FraudAlertRepository) Update(ctx context.Context, alert *fraudv1.FraudAlert) (*fraudv1.FraudAlert, error) {
	var assignedTo sql.NullString
	if alert.AssignedTo != "" {
		assignedTo = sql.NullString{String: alert.AssignedTo, Valid: true}
	}

	var resolvedAt sql.NullTime
	if alert.ResolvedAt != nil {
		resolvedAt = sql.NullTime{Time: alert.ResolvedAt.AsTime(), Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.fraud_alerts
		SET alert_number = $2,
		    entity_type = $3,
		    entity_id = $4,
		    fraud_rule_id = $5,
		    risk_level = $6,
		    fraud_score = $7,
		    details = $8,
		    status = $9,
		    assigned_to = $10,
		    resolved_at = $11,
		    updated_at = NOW()
		WHERE alert_id = $1`,
		alert.Id,
		alert.AlertNumber,
		alert.EntityType,
		alert.EntityId,
		alert.FraudRuleId,
		alert.RiskLevel,
		alert.FraudScore,
		alert.Details,
		strings.ToUpper(alert.Status.String()),
		assignedTo,
		resolvedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update fraud alert: %w", err)
	}

	return r.GetByID(ctx, alert.Id)
}

func (r *FraudAlertRepository) Delete(ctx context.Context, alertID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.fraud_alerts
		WHERE alert_id = $1`,
		alertID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete fraud alert: %w", err)
	}

	return nil
}

func (r *FraudAlertRepository) ListByEntityID(ctx context.Context, entityID string) ([]*fraudv1.FraudAlert, error) {
	query := `
		SELECT alert_id, alert_number, entity_type, entity_id, fraud_rule_id, 
		       risk_level, fraud_score, details, status, assigned_to, resolved_at,
		       audit_info, created_at, updated_at
		FROM insurance_schema.fraud_alerts
		WHERE entity_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.WithContext(ctx).Raw(query, entityID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list fraud alerts: %w", err)
	}
	defer rows.Close()

	alerts := make([]*fraudv1.FraudAlert, 0)
	for rows.Next() {
		var (
			alert      fraudv1.FraudAlert
			statusStr  sql.NullString
			details    sql.NullString
			assignedTo sql.NullString
			resolvedAt sql.NullTime
			auditInfo  sql.NullString
			createdAt  time.Time
			updatedAt  time.Time
		)

		err := rows.Scan(
			&alert.Id,
			&alert.AlertNumber,
			&alert.EntityType,
			&alert.EntityId,
			&alert.FraudRuleId,
			&alert.RiskLevel,
			&alert.FraudScore,
			&details,
			&statusStr,
			&assignedTo,
			&resolvedAt,
			&auditInfo,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan fraud alert: %w", err)
		}

		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := fraudv1.AlertStatus_value[k]; ok {
				alert.Status = fraudv1.AlertStatus(v)
			}
		}

		if details.Valid {
			alert.Details = details.String
		}

		if assignedTo.Valid {
			alert.AssignedTo = assignedTo.String
		}

		if resolvedAt.Valid {
			alert.ResolvedAt = timestamppb.New(resolvedAt.Time)
		}

		alerts = append(alerts, &alert)
	}

	return alerts, nil
}
