package models


// CommissionStructureUpdateResponse represents a commission_structure_update_response
type CommissionStructureUpdateResponse struct {
	Error *Error `json:"error,omitempty"`
	Success bool `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}
