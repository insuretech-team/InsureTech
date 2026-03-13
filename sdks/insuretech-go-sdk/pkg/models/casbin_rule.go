package models


// CasbinRule represents a casbin_rule
type CasbinRule struct {
	V2 string `json:"v2,omitempty"`
	V3 string `json:"v3,omitempty"`
	V4 string `json:"v4,omitempty"`
	V5 string `json:"v5,omitempty"`
	Id string `json:"id,omitempty"`
	Ptype string `json:"ptype,omitempty"`
	V0 string `json:"v0,omitempty"`
	V1 string `json:"v1,omitempty"`
}
