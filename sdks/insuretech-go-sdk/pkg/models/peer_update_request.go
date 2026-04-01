package models


// PeerUpdateRequest represents a peer_update_request
type PeerUpdateRequest struct {
	DisplayName string `json:"display_name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PeerId string `json:"peer_id"`
}
