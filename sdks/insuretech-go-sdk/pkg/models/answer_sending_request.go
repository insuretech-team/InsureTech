package models


// AnswerSendingRequest represents a answer_sending_request
type AnswerSendingRequest struct {
	Sdp string `json:"sdp,omitempty"`
	FromPeerId string `json:"from_peer_id"`
	ToPeerId string `json:"to_peer_id"`
}
