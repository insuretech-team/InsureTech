package models


// SDPOffer represents a sdpoffer
type SDPOffer struct {
	FromPeerId string `json:"from_peer_id,omitempty"`
	ToPeerId string `json:"to_peer_id,omitempty"`
	Sdp string `json:"sdp,omitempty"`
	Type *SDPType `json:"type,omitempty"`
}
