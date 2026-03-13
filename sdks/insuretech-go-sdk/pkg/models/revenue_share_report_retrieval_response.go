package models


// RevenueShareReportRetrievalResponse represents a revenue_share_report_retrieval_response
type RevenueShareReportRetrievalResponse struct {
	TotalGrossPremium *Money `json:"total_gross_premium,omitempty"`
	TotalPlatformShare *Money `json:"total_platform_share,omitempty"`
	TotalInsurerShare *Money `json:"total_insurer_share,omitempty"`
	PolicyCount int `json:"policy_count,omitempty"`
	ByRevenueModel map[string]interface{} `json:"by_revenue_model,omitempty"`
	Error *Error `json:"error,omitempty"`
	InsurerId string `json:"insurer_id,omitempty"`
}
