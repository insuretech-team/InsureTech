package models

import (
	"time"
)

// TINInfo represents a tin_info
type TINInfo struct {
	TaxpayerName string `json:"taxpayer_name,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
}
