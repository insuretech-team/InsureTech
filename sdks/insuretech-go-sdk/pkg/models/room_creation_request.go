package models


// RoomCreationRequest represents a room_creation_request
type RoomCreationRequest struct {
	Config *RoomConfig `json:"config,omitempty"`
	CreatorId string `json:"creator_id"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Name string `json:"name"`
}
