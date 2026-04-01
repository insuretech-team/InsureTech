package models

import (
	"time"
)

// PolicyCreatedEvent represents a policy_created_event
type PolicyCreatedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PolicyNumber string `json:"policy_number,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	SumInsured *Money `json:"sum_insured,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
