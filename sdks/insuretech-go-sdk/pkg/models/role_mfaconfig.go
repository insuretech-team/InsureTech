package models

import (
	"time"
)

// RoleMFAConfig represents a role_mfaconfig
type RoleMFAConfig struct {
	MfaMethods []string `json:"mfa_methods,omitempty"`
	MfaRequired bool `json:"mfa_required"`
	RoleId string `json:"role_id"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
