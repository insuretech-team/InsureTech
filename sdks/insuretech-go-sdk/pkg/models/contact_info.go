package models


// ContactInfo represents a contact_info
type ContactInfo struct {
	AlternateMobile string `json:"alternate_mobile,omitempty"`
	Email string `json:"email,omitempty"`
	Landline string `json:"landline,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
}
