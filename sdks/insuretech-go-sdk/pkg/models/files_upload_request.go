package models


// FilesUploadRequest represents a files_upload_request
type FilesUploadRequest struct {
	TenantId string `json:"tenant_id"`
	Files []*FileUpload `json:"files,omitempty"`
}
