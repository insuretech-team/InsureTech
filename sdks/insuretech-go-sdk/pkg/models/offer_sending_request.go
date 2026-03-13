package models


// OfferSendingRequest represents a offer_sending_request
type OfferSendingRequest struct {
	FromPeerId string `json:"from_peer_id"`
	ToPeerId string `json:"to_peer_id"`
	Sdp string `json:"sdp,omitempty"`
}
