package models

// TrackType represents a track_type
type TrackType string

// TrackType values
const (
	TrackTypeTRACKTYPEUNSPECIFIED TrackType = "TRACK_TYPE_UNSPECIFIED"
	TrackTypeTRACKTYPEAUDIO  = "TRACK_TYPE_AUDIO"
	TrackTypeTRACKTYPEVIDEO  = "TRACK_TYPE_VIDEO"
	TrackTypeTRACKTYPESCREEN  = "TRACK_TYPE_SCREEN"
	TrackTypeTRACKTYPEDATA  = "TRACK_TYPE_DATA"
)
