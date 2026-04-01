package models


// ICECandidateSendingRequest represents a icecandidate_sending_request
type ICECandidateSendingRequest struct {
	Candidate string `json:"candidate,omitempty"`
	FromPeerId string `json:"from_peer_id"`
	SdpMLineIndex int `json:"sdp_m_line_index,omitempty"`
	SdpMid string `json:"sdp_mid,omitempty"`
	ToPeerId string `json:"to_peer_id"`
}
