package models


// ForeignKey represents a foreign_key
type ForeignKey struct {
	ReferencesTable string `json:"references_table,omitempty"`
	ReferencesColumn string `json:"references_column,omitempty"`
	ReferencesSchema string `json:"references_schema,omitempty"`
	OnDelete *ReferentialAction `json:"on_delete,omitempty"`
	OnUpdate *ReferentialAction `json:"on_update,omitempty"`
	ConstraintName string `json:"constraint_name,omitempty"`
}
