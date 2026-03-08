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

type FraudCaseRepository struct {
	db *gorm.DB
}

func NewFraudCaseRepository(db *gorm.DB) *FraudCaseRepository {
	return &FraudCaseRepository{db: db}
}

func (r *FraudCaseRepository) Create(ctx context.Context, fraudCase *fraudv1.FraudCase) (*fraudv1.FraudCase, error) {
	if fraudCase.Id == "" {
		return nil, fmt.Errorf("case_id is required")
	}

	// Get valid user UUID for audit_info
	var createdBy string
	err := r.db.WithContext(ctx).Raw(`SELECT user_id FROM authn_schema.users LIMIT 1`).Scan(&createdBy).Error
	if err != nil || createdBy == "" {
		return nil, fmt.Errorf("failed to get valid user for created_by: %w", err)
	}

	auditInfoJSON := fmt.Sprintf(`{"created_by":"%s","created_at":"%s"}`, createdBy, time.Now().UTC().Format(time.RFC3339))

	// Handle nullable fields
	var investigatorID sql.NullString
	if fraudCase.InvestigatorId != "" {
		investigatorID = sql.NullString{String: fraudCase.InvestigatorId, Valid: true}
	}

	var closedAt sql.NullTime
	if fraudCase.ClosedAt != nil {
		closedAt = sql.NullTime{Time: fraudCase.ClosedAt.AsTime(), Valid: true}
	}

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.fraud_cases
			(case_id, case_number, fraud_alert_id, priority, investigation_notes, 
			 evidence, status, outcome, investigator_id, closed_at, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		fraudCase.Id,
		fraudCase.CaseNumber,
		fraudCase.FraudAlertId,
		strings.ToUpper(fraudCase.Priority.String()),
		fraudCase.InvestigationNotes,
		fraudCase.Evidence,
		strings.ToUpper(fraudCase.Status.String()),
		strings.ToUpper(fraudCase.Outcome.String()),
		investigatorID,
		closedAt,
		auditInfoJSON,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert fraud case: %w", err)
	}

	return r.GetByID(ctx, fraudCase.Id)
}

func (r *FraudCaseRepository) GetByID(ctx context.Context, caseID string) (*fraudv1.FraudCase, error) {
	var (
		fraudCase      fraudv1.FraudCase
		priorityStr    sql.NullString
		statusStr      sql.NullString
		outcomeStr     sql.NullString
		evidence       sql.NullString
		investigatorID sql.NullString
		closedAt       sql.NullTime
		auditInfo      sql.NullString
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT case_id, case_number, fraud_alert_id, priority, 
		       COALESCE(investigation_notes, '') as investigation_notes,
		       evidence, status, outcome, investigator_id, closed_at,
		       audit_info, created_at, updated_at
		FROM insurance_schema.fraud_cases
		WHERE case_id = $1`,
		caseID,
	).Row().Scan(
		&fraudCase.Id,
		&fraudCase.CaseNumber,
		&fraudCase.FraudAlertId,
		&priorityStr,
		&fraudCase.InvestigationNotes,
		&evidence,
		&statusStr,
		&outcomeStr,
		&investigatorID,
		&closedAt,
		&auditInfo,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get fraud case: %w", err)
	}

	// Parse priority enum
	if priorityStr.Valid {
		k := strings.ToUpper(priorityStr.String)
		if v, ok := fraudv1.CasePriority_value[k]; ok {
			fraudCase.Priority = fraudv1.CasePriority(v)
		}
	}

	// Parse status enum
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := fraudv1.CaseStatus_value[k]; ok {
			fraudCase.Status = fraudv1.CaseStatus(v)
		}
	}

	// Parse outcome enum
	if outcomeStr.Valid {
		k := strings.ToUpper(outcomeStr.String)
		if v, ok := fraudv1.CaseOutcome_value[k]; ok {
			fraudCase.Outcome = fraudv1.CaseOutcome(v)
		}
	}

	// Set evidence
	if evidence.Valid {
		fraudCase.Evidence = evidence.String
	}

	// Set investigator_id
	if investigatorID.Valid {
		fraudCase.InvestigatorId = investigatorID.String
	}

	// Set closed_at
	if closedAt.Valid {
		fraudCase.ClosedAt = timestamppb.New(closedAt.Time)
	}

	return &fraudCase, nil
}

func (r *FraudCaseRepository) Update(ctx context.Context, fraudCase *fraudv1.FraudCase) (*fraudv1.FraudCase, error) {
	// Handle nullable fields
	var investigatorID sql.NullString
	if fraudCase.InvestigatorId != "" {
		investigatorID = sql.NullString{String: fraudCase.InvestigatorId, Valid: true}
	}

	var closedAt sql.NullTime
	if fraudCase.ClosedAt != nil {
		closedAt = sql.NullTime{Time: fraudCase.ClosedAt.AsTime(), Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.fraud_cases
		SET case_number = $2,
		    fraud_alert_id = $3,
		    priority = $4,
		    investigation_notes = $5,
		    evidence = $6,
		    status = $7,
		    outcome = $8,
		    investigator_id = $9,
		    closed_at = $10,
		    updated_at = NOW()
		WHERE case_id = $1`,
		fraudCase.Id,
		fraudCase.CaseNumber,
		fraudCase.FraudAlertId,
		strings.ToUpper(fraudCase.Priority.String()),
		fraudCase.InvestigationNotes,
		fraudCase.Evidence,
		strings.ToUpper(fraudCase.Status.String()),
		strings.ToUpper(fraudCase.Outcome.String()),
		investigatorID,
		closedAt,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update fraud case: %w", err)
	}

	return r.GetByID(ctx, fraudCase.Id)
}

func (r *FraudCaseRepository) Delete(ctx context.Context, caseID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.fraud_cases
		WHERE case_id = $1`,
		caseID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete fraud case: %w", err)
	}

	return nil
}

func (r *FraudCaseRepository) ListByAlertID(ctx context.Context, alertID string) ([]*fraudv1.FraudCase, error) {
	query := `
		SELECT case_id, case_number, fraud_alert_id, priority, 
		       COALESCE(investigation_notes, '') as investigation_notes,
		       evidence, status, outcome, investigator_id, closed_at,
		       audit_info, created_at, updated_at
		FROM insurance_schema.fraud_cases
		WHERE fraud_alert_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.WithContext(ctx).Raw(query, alertID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list fraud cases: %w", err)
	}
	defer rows.Close()

	cases := make([]*fraudv1.FraudCase, 0)
	for rows.Next() {
		var (
			fraudCase      fraudv1.FraudCase
			priorityStr    sql.NullString
			statusStr      sql.NullString
			outcomeStr     sql.NullString
			evidence       sql.NullString
			investigatorID sql.NullString
			closedAt       sql.NullTime
			auditInfo      sql.NullString
			createdAt      time.Time
			updatedAt      time.Time
		)

		err := rows.Scan(
			&fraudCase.Id,
			&fraudCase.CaseNumber,
			&fraudCase.FraudAlertId,
			&priorityStr,
			&fraudCase.InvestigationNotes,
			&evidence,
			&statusStr,
			&outcomeStr,
			&investigatorID,
			&closedAt,
			&auditInfo,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan fraud case: %w", err)
		}

		// Parse enums
		if priorityStr.Valid {
			k := strings.ToUpper(priorityStr.String)
			if v, ok := fraudv1.CasePriority_value[k]; ok {
				fraudCase.Priority = fraudv1.CasePriority(v)
			}
		}

		if statusStr.Valid {
			k := strings.ToUpper(statusStr.String)
			if v, ok := fraudv1.CaseStatus_value[k]; ok {
				fraudCase.Status = fraudv1.CaseStatus(v)
			}
		}

		if outcomeStr.Valid {
			k := strings.ToUpper(outcomeStr.String)
			if v, ok := fraudv1.CaseOutcome_value[k]; ok {
				fraudCase.Outcome = fraudv1.CaseOutcome(v)
			}
		}

		if evidence.Valid {
			fraudCase.Evidence = evidence.String
		}

		if investigatorID.Valid {
			fraudCase.InvestigatorId = investigatorID.String
		}

		if closedAt.Valid {
			fraudCase.ClosedAt = timestamppb.New(closedAt.Time)
		}

		cases = append(cases, &fraudCase)
	}

	return cases, nil
}
