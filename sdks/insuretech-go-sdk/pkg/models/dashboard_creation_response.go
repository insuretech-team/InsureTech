package models


// DashboardCreationResponse represents a dashboard_creation_response
type DashboardCreationResponse struct {
	DashboardId string `json:"dashboard_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
