package models


// JoinRoomResponse represents a join_room_response
type JoinRoomResponse struct {
	ExistingPeers []*Peer `json:"existing_peers,omitempty"`
	Peer *Peer `json:"peer,omitempty"`
	Room *Room `json:"room,omitempty"`
}
