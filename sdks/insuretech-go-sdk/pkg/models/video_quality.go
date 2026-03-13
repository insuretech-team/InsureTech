package models

// VideoQuality represents a video_quality
type VideoQuality string

// VideoQuality values
const (
	VideoQualityVIDEOQUALITYUNSPECIFIED VideoQuality = "VIDEO_QUALITY_UNSPECIFIED"
	VideoQualityVIDEOQUALITYLOW  = "VIDEO_QUALITY_LOW"
	VideoQualityVIDEOQUALITYMEDIUM  = "VIDEO_QUALITY_MEDIUM"
	VideoQualityVIDEOQUALITYHIGH  = "VIDEO_QUALITY_HIGH"
	VideoQualityVIDEOQUALITYHD  = "VIDEO_QUALITY_HD"
	VideoQualityVIDEOQUALITY4K  = "VIDEO_QUALITY_4K"
)
