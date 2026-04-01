package models


// LeaveRoomRequest represents a leave_room_request
type LeaveRoomRequest struct {
	PeerId string `json:"peer_id"`
	Reason string `json:"reason,omitempty"`
	RoomId string `json:"room_id"`
}
