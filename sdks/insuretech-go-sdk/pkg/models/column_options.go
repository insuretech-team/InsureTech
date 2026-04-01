package models


// ColumnOptions represents a column_options
type ColumnOptions struct {
	AutoIncrement bool `json:"auto_increment,omitempty"`
	CheckConstraint string `json:"check_constraint,omitempty"`
	ColumnName string `json:"column_name,omitempty"`
	Comment string `json:"comment,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	Encrypted bool `json:"encrypted,omitempty"`
	ExcludeFromInsert bool `json:"exclude_from_insert,omitempty"`
	ExcludeFromUpdate bool `json:"exclude_from_update,omitempty"`
	ForeignKey *ForeignKey `json:"foreign_key,omitempty"`
	Index *IndexOptions `json:"index,omitempty"`
	IsJson bool `json:"is_json,omitempty"`
	NotNull bool `json:"not_null,omitempty"`
	PrimaryKey bool `json:"primary_key,omitempty"`
	SqlType string `json:"sql_type,omitempty"`
	Unique bool `json:"unique,omitempty"`
}
