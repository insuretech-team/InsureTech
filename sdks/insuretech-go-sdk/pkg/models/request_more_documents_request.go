package models


// RequestMoreDocumentsRequest represents a request_more_documents_request
type RequestMoreDocumentsRequest struct {
	RequiredDocumentTypes []string `json:"required_document_types,omitempty"`
	Message string `json:"message,omitempty"`
	ClaimId string `json:"claim_id"`
}
