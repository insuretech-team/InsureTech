package models

import (
	"time"
)

// CrmStatisticsRetrievalRequest represents a crm_statistics_retrieval_request
type CrmStatisticsRetrievalRequest struct {
	AgentId string `json:"agent_id"`
	EndDate time.Time `json:"end_date,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
}
