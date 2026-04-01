package models


// SDPAnswer represents a sdpanswer
type SDPAnswer struct {
	FromPeerId string `json:"from_peer_id,omitempty"`
	Sdp string `json:"sdp,omitempty"`
	ToPeerId string `json:"to_peer_id,omitempty"`
	Type *SDPType `json:"type,omitempty"`
}
