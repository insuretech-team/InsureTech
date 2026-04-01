package models


// FilesUploadRequest represents a files_upload_request
type FilesUploadRequest struct {
	Files []*FileUpload `json:"files,omitempty"`
	TenantId string `json:"tenant_id"`
}
