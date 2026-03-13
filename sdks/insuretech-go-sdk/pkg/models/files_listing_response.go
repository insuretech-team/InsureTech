package models


// FilesListingResponse represents a files_listing_response
type FilesListingResponse struct {
	Files []*StoredFile `json:"files,omitempty"`
	Page *PaginationResponse `json:"page,omitempty"`
}
