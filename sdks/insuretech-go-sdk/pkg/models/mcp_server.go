package models

import (
	"time"
)

// MCPServer represents a mcp_server
type MCPServer struct {
	Capabilities []string `json:"capabilities,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	LastPingAt time.Time `json:"last_ping_at,omitempty"`
	ServerId string `json:"server_id,omitempty"`
	ServerName string `json:"server_name,omitempty"`
	Status *ServerStatus `json:"status,omitempty"`
}
