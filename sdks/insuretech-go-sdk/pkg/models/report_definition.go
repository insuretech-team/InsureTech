package models


// ReportDefinition represents a report_definition
type ReportDefinition struct {
	AuditInfo interface{} `json:"audit_info"`
	Category *ReportCategory `json:"category"`
	Description string `json:"description,omitempty"`
	FormatConfig string `json:"format_config,omitempty"`
	Id string `json:"id"`
	IsActive bool `json:"is_active,omitempty"`
	Name string `json:"name"`
	Parameters string `json:"parameters,omitempty"`
	Query string `json:"query"`
}
