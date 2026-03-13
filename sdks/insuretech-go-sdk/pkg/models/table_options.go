package models


// TableOptions represents a table_options
type TableOptions struct {
	SoftDelete bool `json:"soft_delete,omitempty"`
	EnableRls bool `json:"enable_rls,omitempty"`
	IsTable bool `json:"is_table,omitempty"`
	Comment string `json:"comment,omitempty"`
	AuditFields bool `json:"audit_fields,omitempty"`
	PartitionStrategy *PartitionStrategy `json:"partition_strategy,omitempty"`
	PartitionColumn string `json:"partition_column,omitempty"`
	TableName string `json:"table_name,omitempty"`
	SchemaName string `json:"schema_name,omitempty"`
	MigrationOrder int `json:"migration_order,omitempty"`
}
