package models


// ConnectionStatsRetrievalRequest represents a connection_stats_retrieval_request
type ConnectionStatsRetrievalRequest struct {
	PeerId string `json:"peer_id"`
	RoomId string `json:"room_id"`
}
