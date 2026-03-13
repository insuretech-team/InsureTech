package models


// ConvertQuoteToPolicyResponse represents a convert_quote_to_policy_response
type ConvertQuoteToPolicyResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PolicyNumber string `json:"policy_number,omitempty"`
}
