package models


// CrmStatisticsRetrievalResponse represents a crm_statistics_retrieval_response
type CrmStatisticsRetrievalResponse struct {
	ActiveContacts int `json:"active_contacts,omitempty"`
	ContactedLeads int `json:"contacted_leads,omitempty"`
	ConversionRate float64 `json:"conversion_rate,omitempty"`
	ConvertedLeads int `json:"converted_leads,omitempty"`
	LeadsBySource map[string]interface{} `json:"leads_by_source,omitempty"`
	LostLeads int `json:"lost_leads,omitempty"`
	NewLeads int `json:"new_leads,omitempty"`
	QualifiedLeads int `json:"qualified_leads,omitempty"`
	TotalContacts int `json:"total_contacts,omitempty"`
	TotalLeads int `json:"total_leads,omitempty"`
}
