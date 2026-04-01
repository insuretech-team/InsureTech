package models

// UnderwritingQuoteStatus represents a underwriting_quote_status
type UnderwritingQuoteStatus string

// UnderwritingQuoteStatus values
const (
	UnderwritingQuoteStatusQUOTESTATUSUNSPECIFIED UnderwritingQuoteStatus = "QUOTE_STATUS_UNSPECIFIED"
	UnderwritingQuoteStatusQUOTESTATUSDRAFT  = "QUOTE_STATUS_DRAFT"
	UnderwritingQuoteStatusQUOTESTATUSPENDINGUNDERWRITING  = "QUOTE_STATUS_PENDING_UNDERWRITING"
	UnderwritingQuoteStatusQUOTESTATUSAPPROVED  = "QUOTE_STATUS_APPROVED"
	UnderwritingQuoteStatusQUOTESTATUSREJECTED  = "QUOTE_STATUS_REJECTED"
	UnderwritingQuoteStatusQUOTESTATUSEXPIRED  = "QUOTE_STATUS_EXPIRED"
	UnderwritingQuoteStatusQUOTESTATUSCONVERTED  = "QUOTE_STATUS_CONVERTED"
)
