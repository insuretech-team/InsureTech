package models

// ApiKeyOwnerType represents a api_key_owner_type
type ApiKeyOwnerType string

// ApiKeyOwnerType values
const (
	ApiKeyOwnerTypeAPIKEYOWNERTYPEUNSPECIFIED ApiKeyOwnerType = "API_KEY_OWNER_TYPE_UNSPECIFIED"
	ApiKeyOwnerTypeAPIKEYOWNERTYPEINSURER  = "API_KEY_OWNER_TYPE_INSURER"
	ApiKeyOwnerTypeAPIKEYOWNERTYPEPARTNER  = "API_KEY_OWNER_TYPE_PARTNER"
	ApiKeyOwnerTypeAPIKEYOWNERTYPEINTERNAL  = "API_KEY_OWNER_TYPE_INTERNAL"
)
