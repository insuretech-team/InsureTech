package models


// DocumentAnalysisRequest represents a document_analysis_request
type DocumentAnalysisRequest struct {
	DocumentData string `json:"document_data,omitempty"`
	DocumentType string `json:"document_type,omitempty"`
	DocumentUrl string `json:"document_url"`
}
