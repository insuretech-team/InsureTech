package models


// MarkAsReadResponse represents a mark_as_read_response
type MarkAsReadResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
