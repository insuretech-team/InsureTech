package models


// ConnectionStats represents a connection_stats
type ConnectionStats struct {
	BytesSent string `json:"bytes_sent,omitempty"`
	BytesReceived string `json:"bytes_received,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	BitrateKbps float64 `json:"bitrate_kbps,omitempty"`
	PacketLossPercent float64 `json:"packet_loss_percent,omitempty"`
	RttMs float64 `json:"rtt_ms,omitempty"`
}
