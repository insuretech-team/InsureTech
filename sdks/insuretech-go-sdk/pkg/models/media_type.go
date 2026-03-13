package models

// MediaType represents a media_type
type MediaType string

// MediaType values
const (
	MediaTypeMEDIATYPEUNSPECIFIED MediaType = "MEDIA_TYPE_UNSPECIFIED"
	MediaTypeMEDIATYPEIMAGE  = "MEDIA_TYPE_IMAGE"
	MediaTypeMEDIATYPEDOCUMENT  = "MEDIA_TYPE_DOCUMENT"
	MediaTypeMEDIATYPEVIDEO  = "MEDIA_TYPE_VIDEO"
)
