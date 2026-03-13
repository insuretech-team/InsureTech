package models


// DocumentAnalysisResponse represents a document_analysis_response
type DocumentAnalysisResponse struct {
	AnalysisId string `json:"analysis_id,omitempty"`
	ExtractedData map[string]interface{} `json:"extracted_data,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	VerificationPassed bool `json:"verification_passed,omitempty"`
	Error *Error `json:"error,omitempty"`
}
