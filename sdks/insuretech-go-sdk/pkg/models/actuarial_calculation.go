package models

import (
	"time"
)

// ActuarialCalculation represents a actuarial_calculation
type ActuarialCalculation struct {
	CalculatedAt time.Time `json:"calculated_at"`
	CalculatedBy string `json:"calculated_by,omitempty"`
	CalculatedPremium float64 `json:"calculated_premium,omitempty"`
	CalculationId string `json:"calculation_id"`
	CalculationReference string `json:"calculation_reference"`
	CalculationType *ActuarialCalculationType `json:"calculation_type"`
	CombinedRatio float64 `json:"combined_ratio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	EffectiveDate time.Time `json:"effective_date"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	LossRatio float64 `json:"loss_ratio,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	ParametersJson string `json:"parameters_json"`
	ReserveAmount float64 `json:"reserve_amount,omitempty"`
	ResultsJson string `json:"results_json"`
	Status interface{} `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	ValidationErrors string `json:"validation_errors,omitempty"`
}
