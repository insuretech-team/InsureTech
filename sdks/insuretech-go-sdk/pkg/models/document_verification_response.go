package models


// DocumentVerificationResponse represents a document_verification_response
type DocumentVerificationResponse struct {
	Document *UserDocument `json:"document,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
