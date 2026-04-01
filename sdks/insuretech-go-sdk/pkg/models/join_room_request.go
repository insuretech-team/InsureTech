package models


// JoinRoomRequest represents a join_room_request
type JoinRoomRequest struct {
	Capabilities []*MediaCapability `json:"capabilities,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	JoinToken string `json:"join_token,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	RoomId string `json:"room_id"`
}
