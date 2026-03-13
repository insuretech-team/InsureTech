package models


// PeerAnalyticsRetrievalResponse represents a peer_analytics_retrieval_response
type PeerAnalyticsRetrievalResponse struct {
	Analytics *ParticipantAnalytics `json:"analytics,omitempty"`
}
