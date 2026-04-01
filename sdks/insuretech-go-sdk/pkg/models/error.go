package models


// Error represents a error
type Error struct {
	Code string `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
	DocumentationUrl string `json:"documentation_url,omitempty"`
	ErrorId string `json:"error_id,omitempty"`
	FieldViolations []*FieldViolation `json:"field_violations,omitempty"`
	HttpStatusCode int `json:"http_status_code,omitempty"`
	Message string `json:"message,omitempty"`
	RetryAfterSeconds int `json:"retry_after_seconds,omitempty"`
	Retryable bool `json:"retryable,omitempty"`
}
