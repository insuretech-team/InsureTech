package models

import (
	"time"
)

// RoleMFAConfig represents a role_mfaconfig
type RoleMFAConfig struct {
	RoleId string `json:"role_id"`
	MfaRequired bool `json:"mfa_required"`
	MfaMethods []string `json:"mfa_methods,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
