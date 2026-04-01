package models


// RequestMoreDocumentsRequest represents a request_more_documents_request
type RequestMoreDocumentsRequest struct {
	ClaimId string `json:"claim_id"`
	Message string `json:"message,omitempty"`
	RequiredDocumentTypes []string `json:"required_document_types,omitempty"`
}
