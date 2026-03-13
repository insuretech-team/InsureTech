package models

import (
	"time"
)

// MCPServer represents a mcp_server
type MCPServer struct {
	ServerId string `json:"server_id,omitempty"`
	ServerName string `json:"server_name,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
	Status *ServerStatus `json:"status,omitempty"`
	LastPingAt time.Time `json:"last_ping_at,omitempty"`
}
