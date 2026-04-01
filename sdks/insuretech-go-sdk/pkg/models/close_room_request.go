package models


// CloseRoomRequest represents a close_room_request
type CloseRoomRequest struct {
	Reason string `json:"reason,omitempty"`
	RoomId string `json:"room_id"`
}
