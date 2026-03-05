package core

// Common types used across the WebRTC package

// ConnectionState represents the state of a peer connection
type ConnectionState string

const (
	ConnectionStateNew          ConnectionState = "new"
	ConnectionStateConnecting   ConnectionState = "connecting"
	ConnectionStateConnected    ConnectionState = "connected"
	ConnectionStateDisconnected ConnectionState = "disconnected"
	ConnectionStateFailed       ConnectionState = "failed"
	ConnectionStateClosed       ConnectionState = "closed"
)

// RoomState represents the state of a conference room
type RoomState string

const (
	RoomStateActive RoomState = "active"
	RoomStateIdle   RoomState = "idle"
	RoomStateClosed RoomState = "closed"
)

// TrackKind represents the type of media track
type TrackKind string

const (
	TrackKindAudio TrackKind = "audio"
	TrackKindVideo TrackKind = "video"
	TrackKindScreen TrackKind = "screen"
)

// TrackState represents the state of a media track
type TrackState string

const (
	TrackStateActive TrackState = "active"
	TrackStateMuted  TrackState = "muted"
	TrackStateEnded  TrackState = "ended"
)
