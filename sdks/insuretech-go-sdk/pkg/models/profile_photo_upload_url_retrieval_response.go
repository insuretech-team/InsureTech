package models


// ProfilePhotoUploadURLRetrievalResponse represents a profile_photo_upload_url_retrieval_response
type ProfilePhotoUploadURLRetrievalResponse struct {
	UploadUrl string `json:"upload_url,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
}
