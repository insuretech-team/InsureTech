package models


// HealthDeclarationSubmissionResponse represents a health_declaration_submission_response
type HealthDeclarationSubmissionResponse struct {
	Message string `json:"message,omitempty"`
	MedicalExamRequired bool `json:"medical_exam_required,omitempty"`
	AutoApprovalPossible bool `json:"auto_approval_possible,omitempty"`
	Error *Error `json:"error,omitempty"`
}
