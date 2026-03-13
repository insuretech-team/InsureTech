package models

// DocumentStatus represents a document_status
type DocumentStatus string

// DocumentStatus values
const (
	DocumentStatusDOCUMENTSTATUSUNSPECIFIED DocumentStatus = "DOCUMENT_STATUS_UNSPECIFIED"
	DocumentStatusDOCUMENTSTATUSPENDING  = "DOCUMENT_STATUS_PENDING"
	DocumentStatusDOCUMENTSTATUSVERIFIED  = "DOCUMENT_STATUS_VERIFIED"
	DocumentStatusDOCUMENTSTATUSREJECTED  = "DOCUMENT_STATUS_REJECTED"
)
