package models

import (
	"time"
)

// TINInfo represents a tin_info
type TINInfo struct {
	TinNumber string `json:"tin_number,omitempty"`
	TaxpayerName string `json:"taxpayer_name,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
}
