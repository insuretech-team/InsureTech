package models

import (
	"time"
)

// ResponseMeta represents a response_meta
type ResponseMeta struct {
	ApiVersion string `json:"api_version,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	RequestId string `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
