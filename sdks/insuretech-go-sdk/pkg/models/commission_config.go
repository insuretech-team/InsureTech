package models

import (
	"time"
)

// CommissionConfig represents a commission_config
type CommissionConfig struct {
	RenewalFlatFee *Money `json:"renewal_flat_fee,omitempty"`
	AgentSplitConfig string `json:"agent_split_config,omitempty"`
	PerformanceTiers string `json:"performance_tiers,omitempty"`
	EffectiveFrom time.Time `json:"effective_from"`
	Id string `json:"id"`
	InsurerProductId string `json:"insurer_product_id"`
	AcquisitionFlatFee *Money `json:"acquisition_flat_fee,omitempty"`
	RenewalRate string `json:"renewal_rate,omitempty"`
	ClaimsAssistanceRate string `json:"claims_assistance_rate,omitempty"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	RevenueModel *RevenueModel `json:"revenue_model"`
	AcquisitionRate string `json:"acquisition_rate,omitempty"`
}
