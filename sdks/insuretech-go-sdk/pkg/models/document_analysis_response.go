package models


// DocumentAnalysisResponse represents a document_analysis_response
type DocumentAnalysisResponse struct {
	AnalysisId string `json:"analysis_id,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	ExtractedData map[string]interface{} `json:"extracted_data,omitempty"`
	VerificationPassed bool `json:"verification_passed,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
