package models

import (
	"time"
)

// CommissionConfig represents a commission_config
type CommissionConfig struct {
	AcquisitionFlatFee *Money `json:"acquisition_flat_fee,omitempty"`
	AcquisitionRate string `json:"acquisition_rate,omitempty"`
	AgentSplitConfig string `json:"agent_split_config,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	ClaimsAssistanceRate string `json:"claims_assistance_rate,omitempty"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	Id string `json:"id"`
	InsurerProductId string `json:"insurer_product_id"`
	PerformanceTiers string `json:"performance_tiers,omitempty"`
	RenewalFlatFee *Money `json:"renewal_flat_fee,omitempty"`
	RenewalRate string `json:"renewal_rate,omitempty"`
	RevenueModel *RevenueModel `json:"revenue_model"`
}
