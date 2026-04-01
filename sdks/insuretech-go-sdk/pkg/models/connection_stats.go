package models


// ConnectionStats represents a connection_stats
type ConnectionStats struct {
	BitrateKbps float64 `json:"bitrate_kbps,omitempty"`
	BytesReceived string `json:"bytes_received,omitempty"`
	BytesSent string `json:"bytes_sent,omitempty"`
	PacketLossPercent float64 `json:"packet_loss_percent,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	RttMs float64 `json:"rtt_ms,omitempty"`
}
