package models


// RoomCreationResponse represents a room_creation_response
type RoomCreationResponse struct {
	JoinToken string `json:"join_token,omitempty"`
	Room *Room `json:"room,omitempty"`
}
