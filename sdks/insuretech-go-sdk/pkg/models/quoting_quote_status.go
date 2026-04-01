package models

// QuotingQuoteStatus represents a quoting_quote_status
type QuotingQuoteStatus string

// QuotingQuoteStatus values
const (
	QuotingQuoteStatusQUOTESTATUSUNSPECIFIED QuotingQuoteStatus = "QUOTE_STATUS_UNSPECIFIED"
	QuotingQuoteStatusQUOTESTATUSDRAFT  = "QUOTE_STATUS_DRAFT"
	QuotingQuoteStatusQUOTESTATUSGENERATED  = "QUOTE_STATUS_GENERATED"
	QuotingQuoteStatusQUOTESTATUSSENT  = "QUOTE_STATUS_SENT"
	QuotingQuoteStatusQUOTESTATUSVIEWED  = "QUOTE_STATUS_VIEWED"
	QuotingQuoteStatusQUOTESTATUSACCEPTED  = "QUOTE_STATUS_ACCEPTED"
	QuotingQuoteStatusQUOTESTATUSDECLINED  = "QUOTE_STATUS_DECLINED"
	QuotingQuoteStatusQUOTESTATUSEXPIRED  = "QUOTE_STATUS_EXPIRED"
	QuotingQuoteStatusQUOTESTATUSCONVERTED  = "QUOTE_STATUS_CONVERTED"
)
