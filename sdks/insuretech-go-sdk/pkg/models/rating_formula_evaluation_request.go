package models


// RatingFormulaEvaluationRequest represents a rating_formula_evaluation_request
type RatingFormulaEvaluationRequest struct {
	CalculatedBy string `json:"calculated_by,omitempty"`
	CalculationReference string `json:"calculation_reference,omitempty"`
	FormulaCode string `json:"formula_code,omitempty"`
	FormulaId string `json:"formula_id"`
	SaveCalculation bool `json:"save_calculation,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}
