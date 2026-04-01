package models


// TableOptions represents a table_options
type TableOptions struct {
	AuditFields bool `json:"audit_fields,omitempty"`
	Comment string `json:"comment,omitempty"`
	EnableRls bool `json:"enable_rls,omitempty"`
	IsTable bool `json:"is_table,omitempty"`
	MigrationOrder int `json:"migration_order,omitempty"`
	PartitionColumn string `json:"partition_column,omitempty"`
	PartitionStrategy *PartitionStrategy `json:"partition_strategy,omitempty"`
	SchemaName string `json:"schema_name,omitempty"`
	SoftDelete bool `json:"soft_delete,omitempty"`
	TableName string `json:"table_name,omitempty"`
}
