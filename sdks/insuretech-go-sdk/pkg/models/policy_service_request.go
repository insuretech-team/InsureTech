package models

import (
	"time"
)

// PolicyServiceRequest represents a policy_service_request
type PolicyServiceRequest struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	CustomerId string `json:"customer_id"`
	PolicyId string `json:"policy_id"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
	ProcessedBy string `json:"processed_by,omitempty"`
	RequestData string `json:"request_data,omitempty"`
	RequestId string `json:"request_id"`
	RequestType *ServiceRequestType `json:"request_type,omitempty"`
	Status *ServiceRequestStatus `json:"status,omitempty"`
}
