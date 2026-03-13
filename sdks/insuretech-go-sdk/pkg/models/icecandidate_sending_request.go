package models


// ICECandidateSendingRequest represents a icecandidate_sending_request
type ICECandidateSendingRequest struct {
	ToPeerId string `json:"to_peer_id"`
	Candidate string `json:"candidate,omitempty"`
	SdpMid string `json:"sdp_mid,omitempty"`
	SdpMLineIndex int `json:"sdp_m_line_index,omitempty"`
	FromPeerId string `json:"from_peer_id"`
}
