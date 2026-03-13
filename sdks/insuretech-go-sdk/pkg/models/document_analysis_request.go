package models


// DocumentAnalysisRequest represents a document_analysis_request
type DocumentAnalysisRequest struct {
	DocumentUrl string `json:"document_url"`
	DocumentType string `json:"document_type,omitempty"`
	DocumentData string `json:"document_data,omitempty"`
}
