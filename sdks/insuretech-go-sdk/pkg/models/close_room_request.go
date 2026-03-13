package models


// CloseRoomRequest represents a close_room_request
type CloseRoomRequest struct {
	RoomId string `json:"room_id"`
	Reason string `json:"reason,omitempty"`
}
