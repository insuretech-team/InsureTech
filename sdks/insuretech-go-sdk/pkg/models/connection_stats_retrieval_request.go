package models


// ConnectionStatsRetrievalRequest represents a connection_stats_retrieval_request
type ConnectionStatsRetrievalRequest struct {
	RoomId string `json:"room_id"`
	PeerId string `json:"peer_id"`
}
