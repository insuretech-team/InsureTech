package models


// FilesListingRequest represents a files_listing_request
type FilesListingRequest struct {
	TenantId string `json:"tenant_id"`
	Page *PaginationRequest `json:"page,omitempty"`
	FileType *FileType `json:"file_type,omitempty"`
	ReferenceId string `json:"reference_id"`
	ReferenceType string `json:"reference_type,omitempty"`
	UploadedBy string `json:"uploaded_by,omitempty"`
}
