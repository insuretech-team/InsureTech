package models

import (
	"time"
)

// OfferReceivedEvent represents a offer_received_event
type OfferReceivedEvent struct {
	Offer *SDPOffer `json:"offer,omitempty"`
	ReceivedAt time.Time `json:"received_at,omitempty"`
	RoomId string `json:"room_id,omitempty"`
}
