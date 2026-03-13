package models

// ConfigType represents a config_type
type ConfigType string

// ConfigType values
const (
	ConfigTypeCONFIGTYPEUNSPECIFIED ConfigType = "CONFIG_TYPE_UNSPECIFIED"
	ConfigTypeCONFIGTYPESTRING  = "CONFIG_TYPE_STRING"
	ConfigTypeCONFIGTYPENUMBER  = "CONFIG_TYPE_NUMBER"
	ConfigTypeCONFIGTYPEBOOLEAN  = "CONFIG_TYPE_BOOLEAN"
	ConfigTypeCONFIGTYPEJSON  = "CONFIG_TYPE_JSON"
)
