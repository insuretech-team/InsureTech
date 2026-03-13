package models

// RoomState represents a room_state
type RoomState string

// RoomState values
const (
	RoomStateROOMSTATEUNSPECIFIED RoomState = "ROOM_STATE_UNSPECIFIED"
	RoomStateROOMSTATEACTIVE  = "ROOM_STATE_ACTIVE"
	RoomStateROOMSTATEIDLE  = "ROOM_STATE_IDLE"
	RoomStateROOMSTATECLOSED  = "ROOM_STATE_CLOSED"
)
