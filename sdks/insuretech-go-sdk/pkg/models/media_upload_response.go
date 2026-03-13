package models


// MediaUploadResponse represents a media_upload_response
type MediaUploadResponse struct {
	MediaId string `json:"media_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
