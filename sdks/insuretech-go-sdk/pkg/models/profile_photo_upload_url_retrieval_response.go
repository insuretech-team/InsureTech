package models


// ProfilePhotoUploadURLRetrievalResponse represents a profile_photo_upload_url_retrieval_response
type ProfilePhotoUploadURLRetrievalResponse struct {
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	UploadUrl string `json:"upload_url,omitempty"`
}
