package models


// ICECandidate represents a icecandidate
type ICECandidate struct {
	Candidate string `json:"candidate,omitempty"`
	SdpMid string `json:"sdp_mid,omitempty"`
	SdpMLineIndex int `json:"sdp_m_line_index,omitempty"`
	UsernameFragment string `json:"username_fragment,omitempty"`
	FromPeerId string `json:"from_peer_id,omitempty"`
	ToPeerId string `json:"to_peer_id,omitempty"`
}
