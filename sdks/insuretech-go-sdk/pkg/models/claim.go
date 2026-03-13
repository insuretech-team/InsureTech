package models

import (
	"time"
)

// Claim represents a claim
type Claim struct {
	DeductibleAmount *Money `json:"deductible_amount,omitempty"`
	ClaimedCurrency string `json:"claimed_currency,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	SettledAmount *Money `json:"settled_amount,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	PlaceOfIncident string `json:"place_of_incident,omitempty"`
	SettledCurrency string `json:"settled_currency,omitempty"`
	Type *ClaimType `json:"type,omitempty"`
	IncidentDate time.Time `json:"incident_date,omitempty"`
	SubmittedAt time.Time `json:"submitted_at,omitempty"`
	SettledAt time.Time `json:"settled_at,omitempty"`
	InAppMessages string `json:"in_app_messages,omitempty"`
	ClaimNumber string `json:"claim_number,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Documents []*ClaimDocument `json:"documents,omitempty"`
	BankDetailsForPayout string `json:"bank_details_for_payout,omitempty"`
	CoPayAmount *Money `json:"co_pay_amount,omitempty"`
	ProcessorNotes string `json:"processor_notes,omitempty"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Approvals []*ClaimApproval `json:"approvals,omitempty"`
	ApprovedAmount *Money `json:"approved_amount,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	ProcessingType *ClaimProcessingType `json:"processing_type,omitempty"`
	ApprovedCurrency string `json:"approved_currency,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	Status *ClaimStatus `json:"status,omitempty"`
	ClaimedAmount *Money `json:"claimed_amount,omitempty"`
	IncidentDescription string `json:"incident_description,omitempty"`
	FraudCheck *FraudCheckResult `json:"fraud_check,omitempty"`
	AppealOptionAvailable bool `json:"appeal_option_available,omitempty"`
}
