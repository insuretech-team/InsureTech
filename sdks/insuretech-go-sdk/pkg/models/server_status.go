package models

// ServerStatus represents a server_status
type ServerStatus string

// ServerStatus values
const (
	ServerStatusSERVERSTATUSUNSPECIFIED ServerStatus = "SERVER_STATUS_UNSPECIFIED"
	ServerStatusSERVERSTATUSONLINE  = "SERVER_STATUS_ONLINE"
	ServerStatusSERVERSTATUSOFFLINE  = "SERVER_STATUS_OFFLINE"
	ServerStatusSERVERSTATUSDEGRADED  = "SERVER_STATUS_DEGRADED"
)
