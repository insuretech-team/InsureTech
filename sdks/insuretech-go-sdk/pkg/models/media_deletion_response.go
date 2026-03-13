package models


// MediaDeletionResponse represents a media_deletion_response
type MediaDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
