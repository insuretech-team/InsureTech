package models


// EffectivePermission represents a effective_permission
type EffectivePermission struct {
	ViaRole string `json:"via_role,omitempty"`
	Object string `json:"object,omitempty"`
	Action string `json:"action,omitempty"`
}
