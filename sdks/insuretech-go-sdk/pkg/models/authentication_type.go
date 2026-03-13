package models

// AuthenticationType represents a authentication_type
type AuthenticationType string

// AuthenticationType values
const (
	AuthenticationTypeAUTHENTICATIONTYPEUNSPECIFIED AuthenticationType = "AUTHENTICATION_TYPE_UNSPECIFIED"
	AuthenticationTypeAUTHENTICATIONTYPEAPIKEY  = "AUTHENTICATION_TYPE_API_KEY"
	AuthenticationTypeAUTHENTICATIONTYPEOAUTH2  = "AUTHENTICATION_TYPE_OAUTH2"
	AuthenticationTypeAUTHENTICATIONTYPEBASICAUTH  = "AUTHENTICATION_TYPE_BASIC_AUTH"
	AuthenticationTypeAUTHENTICATIONTYPEMUTUALTLS  = "AUTHENTICATION_TYPE_MUTUAL_TLS"
)
