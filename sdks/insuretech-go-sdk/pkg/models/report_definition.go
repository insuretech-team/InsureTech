package models


// ReportDefinition represents a report_definition
type ReportDefinition struct {
	IsActive bool `json:"is_active,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Category *ReportCategory `json:"category"`
	Description string `json:"description,omitempty"`
	FormatConfig string `json:"format_config,omitempty"`
	Id string `json:"id"`
	Name string `json:"name"`
	Query string `json:"query"`
	Parameters string `json:"parameters,omitempty"`
}
