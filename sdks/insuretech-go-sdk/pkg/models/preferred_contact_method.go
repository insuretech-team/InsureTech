package models

// PreferredContactMethod represents a preferred_contact_method
type PreferredContactMethod string

// PreferredContactMethod values
const (
	PreferredContactMethodPREFERREDCONTACTMETHODUNSPECIFIED PreferredContactMethod = "PREFERRED_CONTACT_METHOD_UNSPECIFIED"
	PreferredContactMethodPREFERREDCONTACTMETHODEMAIL  = "PREFERRED_CONTACT_METHOD_EMAIL"
	PreferredContactMethodPREFERREDCONTACTMETHODPHONE  = "PREFERRED_CONTACT_METHOD_PHONE"
	PreferredContactMethodPREFERREDCONTACTMETHODSMS  = "PREFERRED_CONTACT_METHOD_SMS"
	PreferredContactMethodPREFERREDCONTACTMETHODWHATSAPP  = "PREFERRED_CONTACT_METHOD_WHATSAPP"
)
