package models


// SDPAnswer represents a sdpanswer
type SDPAnswer struct {
	ToPeerId string `json:"to_peer_id,omitempty"`
	Sdp string `json:"sdp,omitempty"`
	Type *SDPType `json:"type,omitempty"`
	FromPeerId string `json:"from_peer_id,omitempty"`
}
