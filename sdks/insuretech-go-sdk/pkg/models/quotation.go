package models

import (
	"time"
)

// Quotation represents a quotation
type Quotation struct {
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	ApprovedByUserId string `json:"approved_by_user_id,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedByUserId string `json:"created_by_user_id,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	EmployeeNo int `json:"employee_no,omitempty"`
	EstimatedPremium *Money `json:"estimated_premium,omitempty"`
	InsuranceCategory *InsuranceType `json:"insurance_category,omitempty"`
	InsurerName string `json:"insurer_name,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
	QuotationId string `json:"quotation_id,omitempty"`
	QuotationNumber string `json:"quotation_number,omitempty"`
	QuotedAmount *Money `json:"quoted_amount,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	Status *QuotationStatus `json:"status,omitempty"`
	SubmissionDate time.Time `json:"submission_date,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}
