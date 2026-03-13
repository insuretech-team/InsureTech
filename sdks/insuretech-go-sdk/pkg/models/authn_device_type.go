package models

// AuthnDeviceType represents a authn_device_type
type AuthnDeviceType string

// AuthnDeviceType values
const (
	AuthnDeviceTypeDEVICETYPEUNSPECIFIED AuthnDeviceType = "DEVICE_TYPE_UNSPECIFIED"
	AuthnDeviceTypeDEVICETYPEWEB  = "DEVICE_TYPE_WEB"
	AuthnDeviceTypeDEVICETYPEMOBILEANDROID  = "DEVICE_TYPE_MOBILE_ANDROID"
	AuthnDeviceTypeDEVICETYPEMOBILEIOS  = "DEVICE_TYPE_MOBILE_IOS"
	AuthnDeviceTypeDEVICETYPEAPI  = "DEVICE_TYPE_API"
	AuthnDeviceTypeDEVICETYPEDESKTOP  = "DEVICE_TYPE_DESKTOP"
)
