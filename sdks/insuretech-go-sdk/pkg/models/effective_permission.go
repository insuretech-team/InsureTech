package models


// EffectivePermission represents a effective_permission
type EffectivePermission struct {
	Action string `json:"action,omitempty"`
	Object string `json:"object,omitempty"`
	ViaRole string `json:"via_role,omitempty"`
}
