package models


// RevenueShareReportRetrievalResponse represents a revenue_share_report_retrieval_response
type RevenueShareReportRetrievalResponse struct {
	ByRevenueModel map[string]interface{} `json:"by_revenue_model,omitempty"`
	InsurerId string `json:"insurer_id,omitempty"`
	PolicyCount int `json:"policy_count,omitempty"`
	TotalGrossPremium *Money `json:"total_gross_premium,omitempty"`
	TotalInsurerShare *Money `json:"total_insurer_share,omitempty"`
	TotalPlatformShare *Money `json:"total_platform_share,omitempty"`
}
