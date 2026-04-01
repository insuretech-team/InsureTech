package models

import (
	"time"
)

// PortalConfigUpdatedEvent represents a portal_config_updated_event
type PortalConfigUpdatedEvent struct {
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	Portal *Portal `json:"portal,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
