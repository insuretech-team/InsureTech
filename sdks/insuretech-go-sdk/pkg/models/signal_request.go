package models


// SignalRequest represents a signal_request
type SignalRequest struct {
	Answer *AnswerSendingRequest `json:"answer,omitempty"`
	IceCandidate *ICECandidateSendingRequest `json:"ice_candidate,omitempty"`
	Offer *OfferSendingRequest `json:"offer,omitempty"`
	PeerId string `json:"peer_id"`
	Ping *PingRequest `json:"ping,omitempty"`
	RoomId string `json:"room_id"`
}
