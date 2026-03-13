package models


// SignalRequest represents a signal_request
type SignalRequest struct {
	Ping *PingRequest `json:"ping,omitempty"`
	PeerId string `json:"peer_id"`
	RoomId string `json:"room_id"`
	Offer *OfferSendingRequest `json:"offer,omitempty"`
	Answer *AnswerSendingRequest `json:"answer,omitempty"`
	IceCandidate *ICECandidateSendingRequest `json:"ice_candidate,omitempty"`
}
