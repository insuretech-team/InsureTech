package models

import (
	"time"
)

// Claim represents a claim
type Claim struct {
	AppealOptionAvailable bool `json:"appeal_option_available,omitempty"`
	Approvals []*ClaimApproval `json:"approvals,omitempty"`
	ApprovedAmount *Money `json:"approved_amount,omitempty"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	ApprovedCurrency string `json:"approved_currency,omitempty"`
	BankDetailsForPayout string `json:"bank_details_for_payout,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	ClaimNumber string `json:"claim_number,omitempty"`
	ClaimedAmount *Money `json:"claimed_amount,omitempty"`
	ClaimedCurrency string `json:"claimed_currency,omitempty"`
	CoPayAmount *Money `json:"co_pay_amount,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	DeductibleAmount *Money `json:"deductible_amount,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Documents []*ClaimDocument `json:"documents,omitempty"`
	FraudCheck *FraudCheckResult `json:"fraud_check,omitempty"`
	InAppMessages string `json:"in_app_messages,omitempty"`
	IncidentDate time.Time `json:"incident_date,omitempty"`
	IncidentDescription string `json:"incident_description,omitempty"`
	PlaceOfIncident string `json:"place_of_incident,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	ProcessingType *ClaimProcessingType `json:"processing_type,omitempty"`
	ProcessorNotes string `json:"processor_notes,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	SettledAmount *Money `json:"settled_amount,omitempty"`
	SettledAt time.Time `json:"settled_at,omitempty"`
	SettledCurrency string `json:"settled_currency,omitempty"`
	Status *ClaimStatus `json:"status,omitempty"`
	SubmittedAt time.Time `json:"submitted_at,omitempty"`
	Type *ClaimType `json:"type,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
