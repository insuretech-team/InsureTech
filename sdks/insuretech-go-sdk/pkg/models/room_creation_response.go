package models


// RoomCreationResponse represents a room_creation_response
type RoomCreationResponse struct {
	Room *Room `json:"room,omitempty"`
	JoinToken string `json:"join_token,omitempty"`
}
