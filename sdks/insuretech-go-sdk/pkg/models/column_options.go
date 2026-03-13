package models


// ColumnOptions represents a column_options
type ColumnOptions struct {
	Unique bool `json:"unique,omitempty"`
	PrimaryKey bool `json:"primary_key,omitempty"`
	AutoIncrement bool `json:"auto_increment,omitempty"`
	Comment string `json:"comment,omitempty"`
	Encrypted bool `json:"encrypted,omitempty"`
	SqlType string `json:"sql_type,omitempty"`
	ForeignKey *ForeignKey `json:"foreign_key,omitempty"`
	ExcludeFromInsert bool `json:"exclude_from_insert,omitempty"`
	ExcludeFromUpdate bool `json:"exclude_from_update,omitempty"`
	ColumnName string `json:"column_name,omitempty"`
	NotNull bool `json:"not_null,omitempty"`
	CheckConstraint string `json:"check_constraint,omitempty"`
	Index *IndexOptions `json:"index,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	IsJson bool `json:"is_json,omitempty"`
}
