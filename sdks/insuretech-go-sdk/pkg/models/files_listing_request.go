package models


// FilesListingRequest represents a files_listing_request
type FilesListingRequest struct {
	FileType *FileType `json:"file_type,omitempty"`
	Page *PaginationRequest `json:"page,omitempty"`
	ReferenceId string `json:"reference_id"`
	ReferenceType string `json:"reference_type,omitempty"`
	TenantId string `json:"tenant_id"`
	UploadedBy string `json:"uploaded_by,omitempty"`
}
