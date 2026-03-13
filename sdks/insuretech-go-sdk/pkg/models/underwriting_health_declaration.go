package models

import (
	"time"
)

// UnderwritingHealthDeclaration represents a underwriting_health_declaration
type UnderwritingHealthDeclaration struct {
	HasPreExistingConditions bool `json:"has_pre_existing_conditions,omitempty"`
	AlcoholConsumer bool `json:"alcohol_consumer,omitempty"`
	MedicalExamRequired bool `json:"medical_exam_required,omitempty"`
	MedicalExamResults string `json:"medical_exam_results,omitempty"`
	MedicalExamDate time.Time `json:"medical_exam_date,omitempty"`
	QuoteId string `json:"quote_id"`
	Smoker bool `json:"smoker,omitempty"`
	Id string `json:"id"`
	HeightCm int `json:"height_cm,omitempty"`
	IsCurrentlyHospitalized bool `json:"is_currently_hospitalized,omitempty"`
	HasFamilyHistory bool `json:"has_family_history,omitempty"`
	OccupationRiskLevel string `json:"occupation_risk_level,omitempty"`
	MedicalExamCompleted bool `json:"medical_exam_completed,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	WeightKg string `json:"weight_kg,omitempty"`
	Bmi string `json:"bmi,omitempty"`
	PreExistingConditions string `json:"pre_existing_conditions,omitempty"`
	FamilyHistory string `json:"family_history,omitempty"`
	MedicalDocuments string `json:"medical_documents,omitempty"`
}
