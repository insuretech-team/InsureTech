package models

import (
	"time"
)

// ApiKeyUsage represents a api_key_usage
type ApiKeyUsage struct {
	ApiKeyId string `json:"api_key_id"`
	Endpoint string `json:"endpoint"`
	HttpMethod string `json:"http_method"`
	Id string `json:"id"`
	RequestIp string `json:"request_ip,omitempty"`
	RequestPayload string `json:"request_payload,omitempty"`
	ResponsePayload string `json:"response_payload,omitempty"`
	ResponseTimeMs int `json:"response_time_ms,omitempty"`
	StatusCode int `json:"status_code"`
	Timestamp time.Time `json:"timestamp"`
	TraceId string `json:"trace_id,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
