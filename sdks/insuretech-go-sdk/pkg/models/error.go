package models


// Error represents a error
type Error struct {
	Retryable bool `json:"retryable,omitempty"`
	RetryAfterSeconds int `json:"retry_after_seconds,omitempty"`
	Code string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	HttpStatusCode int `json:"http_status_code,omitempty"`
	ErrorId string `json:"error_id,omitempty"`
	DocumentationUrl string `json:"documentation_url,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
	FieldViolations []*FieldViolation `json:"field_violations,omitempty"`
}
