package models

import (
	"time"
)

// Beneficiary represents a beneficiary
type Beneficiary struct {
	AuditInfo interface{} `json:"audit_info"`
	BeneficiaryId string `json:"beneficiary_id"`
	Code string `json:"code"`
	KycCompletedAt time.Time `json:"kyc_completed_at,omitempty"`
	KycStatus interface{} `json:"kyc_status"`
	PartnerId string `json:"partner_id,omitempty"`
	ReferralCode string `json:"referral_code,omitempty"`
	ReferredBy string `json:"referred_by,omitempty"`
	RiskScore string `json:"risk_score,omitempty"`
	Status interface{} `json:"status"`
	Type *BeneficiaryType `json:"type"`
	UserId string `json:"user_id"`
}
