package models


// RoomUpdateRequest represents a room_update_request
type RoomUpdateRequest struct {
	RoomId string `json:"room_id"`
	Config *RoomConfig `json:"config,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
