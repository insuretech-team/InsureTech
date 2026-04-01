package models


// AgentDashboardRetrievalResponse represents a agent_dashboard_retrieval_response
type AgentDashboardRetrievalResponse struct {
	LeadsConvertedThisMonth int `json:"leads_converted_this_month,omitempty"`
	MyContactsCount int `json:"my_contacts_count,omitempty"`
	MyLeadsCount int `json:"my_leads_count,omitempty"`
	NewLeadsToday int `json:"new_leads_today,omitempty"`
	RecentContacts []*Contact `json:"recent_contacts,omitempty"`
	RecentLeads []*Lead `json:"recent_leads,omitempty"`
}
