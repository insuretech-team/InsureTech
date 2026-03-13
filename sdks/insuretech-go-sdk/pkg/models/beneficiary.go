package models

import (
	"time"
)

// Beneficiary represents a beneficiary
type Beneficiary struct {
	BeneficiaryId string `json:"beneficiary_id"`
	Type *BeneficiaryType `json:"type"`
	Code string `json:"code"`
	Status interface{} `json:"status"`
	KycStatus interface{} `json:"kyc_status"`
	KycCompletedAt time.Time `json:"kyc_completed_at,omitempty"`
	UserId string `json:"user_id"`
	RiskScore string `json:"risk_score,omitempty"`
	ReferralCode string `json:"referral_code,omitempty"`
	ReferredBy string `json:"referred_by,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
}
