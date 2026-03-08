package repository

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	underwritingv1 "github.com/newage-saint/insuretech/gen/go/insuretech/underwriting/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type HealthDeclarationRepository struct {
	db *gorm.DB
}

func NewHealthDeclarationRepository(db *gorm.DB) *HealthDeclarationRepository {
	return &HealthDeclarationRepository{db: db}
}

func (r *HealthDeclarationRepository) Create(ctx context.Context, declaration *underwritingv1.HealthDeclaration) (*underwritingv1.HealthDeclaration, error) {
	if declaration.Id == "" {
		return nil, fmt.Errorf("declaration_id is required")
	}

	// Handle JSONB fields
	var preExistingConditions interface{}
	if declaration.PreExistingConditions != "" {
		preExistingConditions = declaration.PreExistingConditions
	}

	var familyHistory interface{}
	if declaration.FamilyHistory != "" {
		familyHistory = declaration.FamilyHistory
	}

	var medicalExamResults interface{}
	if declaration.MedicalExamResults != "" {
		medicalExamResults = declaration.MedicalExamResults
	}

	var medicalDocuments interface{}
	if declaration.MedicalDocuments != "" {
		medicalDocuments = declaration.MedicalDocuments
	}

	// Handle timestamps
	var medicalExamDate sql.NullTime
	if declaration.MedicalExamDate != nil {
		medicalExamDate = sql.NullTime{Time: declaration.MedicalExamDate.AsTime(), Valid: true}
	}

	var auditInfo interface{}
	if declaration.AuditInfo != nil {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.health_declarations
			(declaration_id, quote_id, height_cm, weight_kg, bmi,
			 has_pre_existing_conditions, pre_existing_conditions,
			 is_currently_hospitalized, has_family_history, family_history,
			 smoker, alcohol_consumer, occupation_risk_level,
			 medical_exam_required, medical_exam_completed, medical_exam_results,
			 medical_exam_date, medical_documents, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
		declaration.Id,
		declaration.QuoteId,
		declaration.HeightCm,
		declaration.WeightKg,
		declaration.Bmi,
		declaration.HasPreExistingConditions,
		preExistingConditions,
		declaration.IsCurrentlyHospitalized,
		declaration.HasFamilyHistory,
		familyHistory,
		declaration.Smoker,
		declaration.AlcoholConsumer,
		declaration.OccupationRiskLevel,
		declaration.MedicalExamRequired,
		declaration.MedicalExamCompleted,
		medicalExamResults,
		medicalExamDate,
		medicalDocuments,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert health declaration: %w", err)
	}

	return r.GetByID(ctx, declaration.Id)
}

func (r *HealthDeclarationRepository) GetByID(ctx context.Context, declarationID string) (*underwritingv1.HealthDeclaration, error) {
	var (
		h                       underwritingv1.HealthDeclaration
		heightCm                sql.NullInt64
		weightKg                sql.NullString
		bmi                     sql.NullString
		preExistingConditions   sql.NullString
		familyHistory           sql.NullString
		occupationRiskLevel     sql.NullString
		medicalExamResults      sql.NullString
		medicalExamDate         sql.NullTime
		medicalDocuments        sql.NullString
		auditInfo               sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT declaration_id, quote_id, height_cm, weight_kg, bmi,
		       has_pre_existing_conditions, pre_existing_conditions,
		       is_currently_hospitalized, has_family_history, family_history,
		       smoker, alcohol_consumer, occupation_risk_level,
		       medical_exam_required, medical_exam_completed, medical_exam_results,
		       medical_exam_date, medical_documents, audit_info
		FROM insurance_schema.health_declarations
		WHERE declaration_id = $1`,
		declarationID,
	).Row().Scan(
		&h.Id,
		&h.QuoteId,
		&heightCm,
		&weightKg,
		&bmi,
		&h.HasPreExistingConditions,
		&preExistingConditions,
		&h.IsCurrentlyHospitalized,
		&h.HasFamilyHistory,
		&familyHistory,
		&h.Smoker,
		&h.AlcoholConsumer,
		&occupationRiskLevel,
		&h.MedicalExamRequired,
		&h.MedicalExamCompleted,
		&medicalExamResults,
		&medicalExamDate,
		&medicalDocuments,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get health declaration: %w", err)
	}

	// Set optional fields
	if heightCm.Valid {
		h.HeightCm = int32(heightCm.Int64)
	}
	if weightKg.Valid {
		h.WeightKg = weightKg.String
	}
	if bmi.Valid {
		h.Bmi = bmi.String
	}
	if preExistingConditions.Valid {
		h.PreExistingConditions = preExistingConditions.String
	}
	if familyHistory.Valid {
		h.FamilyHistory = familyHistory.String
	}
	if occupationRiskLevel.Valid {
		h.OccupationRiskLevel = occupationRiskLevel.String
	}
	if medicalExamResults.Valid {
		h.MedicalExamResults = medicalExamResults.String
	}
	if medicalDocuments.Valid {
		h.MedicalDocuments = medicalDocuments.String
	}

	// Set timestamps
	if medicalExamDate.Valid {
		h.MedicalExamDate = timestamppb.New(medicalExamDate.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		h.AuditInfo = &commonv1.AuditInfo{}
	}

	return &h, nil
}

func (r *HealthDeclarationRepository) Update(ctx context.Context, declaration *underwritingv1.HealthDeclaration) (*underwritingv1.HealthDeclaration, error) {
	// Handle JSONB fields
	var preExistingConditions interface{}
	if declaration.PreExistingConditions != "" {
		preExistingConditions = declaration.PreExistingConditions
	}

	var familyHistory interface{}
	if declaration.FamilyHistory != "" {
		familyHistory = declaration.FamilyHistory
	}

	var medicalExamResults interface{}
	if declaration.MedicalExamResults != "" {
		medicalExamResults = declaration.MedicalExamResults
	}

	var medicalDocuments interface{}
	if declaration.MedicalDocuments != "" {
		medicalDocuments = declaration.MedicalDocuments
	}

	// Handle timestamps
	var medicalExamDate sql.NullTime
	if declaration.MedicalExamDate != nil {
		medicalExamDate = sql.NullTime{Time: declaration.MedicalExamDate.AsTime(), Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.health_declarations
		SET quote_id = $2,
		    height_cm = $3,
		    weight_kg = $4,
		    bmi = $5,
		    has_pre_existing_conditions = $6,
		    pre_existing_conditions = $7,
		    is_currently_hospitalized = $8,
		    has_family_history = $9,
		    family_history = $10,
		    smoker = $11,
		    alcohol_consumer = $12,
		    occupation_risk_level = $13,
		    medical_exam_required = $14,
		    medical_exam_completed = $15,
		    medical_exam_results = $16,
		    medical_exam_date = $17,
		    medical_documents = $18
		WHERE declaration_id = $1`,
		declaration.Id,
		declaration.QuoteId,
		declaration.HeightCm,
		declaration.WeightKg,
		declaration.Bmi,
		declaration.HasPreExistingConditions,
		preExistingConditions,
		declaration.IsCurrentlyHospitalized,
		declaration.HasFamilyHistory,
		familyHistory,
		declaration.Smoker,
		declaration.AlcoholConsumer,
		declaration.OccupationRiskLevel,
		declaration.MedicalExamRequired,
		declaration.MedicalExamCompleted,
		medicalExamResults,
		medicalExamDate,
		medicalDocuments,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update health declaration: %w", err)
	}

	return r.GetByID(ctx, declaration.Id)
}

func (r *HealthDeclarationRepository) Delete(ctx context.Context, declarationID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.health_declarations
		WHERE declaration_id = $1`,
		declarationID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete health declaration: %w", err)
	}

	return nil
}

func (r *HealthDeclarationRepository) GetByQuoteID(ctx context.Context, quoteID string) (*underwritingv1.HealthDeclaration, error) {
	var (
		h                       underwritingv1.HealthDeclaration
		heightCm                sql.NullInt64
		weightKg                sql.NullString
		bmi                     sql.NullString
		preExistingConditions   sql.NullString
		familyHistory           sql.NullString
		occupationRiskLevel     sql.NullString
		medicalExamResults      sql.NullString
		medicalExamDate         sql.NullTime
		medicalDocuments        sql.NullString
		auditInfo               sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT declaration_id, quote_id, height_cm, weight_kg, bmi,
		       has_pre_existing_conditions, pre_existing_conditions,
		       is_currently_hospitalized, has_family_history, family_history,
		       smoker, alcohol_consumer, occupation_risk_level,
		       medical_exam_required, medical_exam_completed, medical_exam_results,
		       medical_exam_date, medical_documents, audit_info
		FROM insurance_schema.health_declarations
		WHERE quote_id = $1`,
		quoteID,
	).Row().Scan(
		&h.Id,
		&h.QuoteId,
		&heightCm,
		&weightKg,
		&bmi,
		&h.HasPreExistingConditions,
		&preExistingConditions,
		&h.IsCurrentlyHospitalized,
		&h.HasFamilyHistory,
		&familyHistory,
		&h.Smoker,
		&h.AlcoholConsumer,
		&occupationRiskLevel,
		&h.MedicalExamRequired,
		&h.MedicalExamCompleted,
		&medicalExamResults,
		&medicalExamDate,
		&medicalDocuments,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get health declaration by quote_id: %w", err)
	}

	// Set optional fields
	if heightCm.Valid {
		h.HeightCm = int32(heightCm.Int64)
	}
	if weightKg.Valid {
		h.WeightKg = weightKg.String
	}
	if bmi.Valid {
		h.Bmi = bmi.String
	}
	if preExistingConditions.Valid {
		h.PreExistingConditions = preExistingConditions.String
	}
	if familyHistory.Valid {
		h.FamilyHistory = familyHistory.String
	}
	if occupationRiskLevel.Valid {
		h.OccupationRiskLevel = occupationRiskLevel.String
	}
	if medicalExamResults.Valid {
		h.MedicalExamResults = medicalExamResults.String
	}
	if medicalDocuments.Valid {
		h.MedicalDocuments = medicalDocuments.String
	}

	// Set timestamps
	if medicalExamDate.Valid {
		h.MedicalExamDate = timestamppb.New(medicalExamDate.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		h.AuditInfo = &commonv1.AuditInfo{}
	}

	return &h, nil
}
