package models

import (
	"time"
)

// UnderwritingHealthDeclaration represents a underwriting_health_declaration
type UnderwritingHealthDeclaration struct {
	AlcoholConsumer bool `json:"alcohol_consumer,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Bmi string `json:"bmi,omitempty"`
	FamilyHistory string `json:"family_history,omitempty"`
	HasFamilyHistory bool `json:"has_family_history,omitempty"`
	HasPreExistingConditions bool `json:"has_pre_existing_conditions,omitempty"`
	HeightCm int `json:"height_cm,omitempty"`
	Id string `json:"id"`
	IsCurrentlyHospitalized bool `json:"is_currently_hospitalized,omitempty"`
	MedicalDocuments string `json:"medical_documents,omitempty"`
	MedicalExamCompleted bool `json:"medical_exam_completed,omitempty"`
	MedicalExamDate time.Time `json:"medical_exam_date,omitempty"`
	MedicalExamRequired bool `json:"medical_exam_required,omitempty"`
	MedicalExamResults string `json:"medical_exam_results,omitempty"`
	OccupationRiskLevel string `json:"occupation_risk_level,omitempty"`
	PreExistingConditions string `json:"pre_existing_conditions,omitempty"`
	QuoteId string `json:"quote_id"`
	Smoker bool `json:"smoker,omitempty"`
	WeightKg string `json:"weight_kg,omitempty"`
}
