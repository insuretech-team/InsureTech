package models


// ProfilePhotoUploadURLRetrievalRequest represents a profile_photo_upload_url_retrieval_request
type ProfilePhotoUploadURLRetrievalRequest struct {
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	UserId string `json:"user_id"`
	ContentType string `json:"content_type,omitempty"`
}
