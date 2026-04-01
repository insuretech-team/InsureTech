package models


// ICECandidate represents a icecandidate
type ICECandidate struct {
	Candidate string `json:"candidate,omitempty"`
	FromPeerId string `json:"from_peer_id,omitempty"`
	SdpMLineIndex int `json:"sdp_m_line_index,omitempty"`
	SdpMid string `json:"sdp_mid,omitempty"`
	ToPeerId string `json:"to_peer_id,omitempty"`
	UsernameFragment string `json:"username_fragment,omitempty"`
}
