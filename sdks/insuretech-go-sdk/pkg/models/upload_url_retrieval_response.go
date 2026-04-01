package models


// UploadURLRetrievalResponse represents a upload_url_retrieval_response
type UploadURLRetrievalResponse struct {
	FileId string `json:"file_id,omitempty"`
	StorageKey string `json:"storage_key,omitempty"`
	UploadUrl string `json:"upload_url,omitempty"`
}
