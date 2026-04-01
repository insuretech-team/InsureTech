package models


// ReportScheduleCreationRequest represents a report_schedule_creation_request
type ReportScheduleCreationRequest struct {
	CronExpression string `json:"cron_expression,omitempty"`
	Frequency string `json:"frequency,omitempty"`
	Name string `json:"name"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	ReportDefinitionId string `json:"report_definition_id"`
}
