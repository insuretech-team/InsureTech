package models


// RoomRetrievalResponse represents a room_retrieval_response
type RoomRetrievalResponse struct {
	Peers []*Peer `json:"peers,omitempty"`
	Room *Room `json:"room,omitempty"`
}
