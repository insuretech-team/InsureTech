package models

// TrackState represents a track_state
type TrackState string

// TrackState values
const (
	TrackStateTRACKSTATEUNSPECIFIED TrackState = "TRACK_STATE_UNSPECIFIED"
	TrackStateTRACKSTATEACTIVE  = "TRACK_STATE_ACTIVE"
	TrackStateTRACKSTATEINACTIVE  = "TRACK_STATE_INACTIVE"
	TrackStateTRACKSTATEENDED  = "TRACK_STATE_ENDED"
	TrackStateTRACKSTATEPAUSED  = "TRACK_STATE_PAUSED"
)
