package models


// HealthDeclarationSubmissionResponse represents a health_declaration_submission_response
type HealthDeclarationSubmissionResponse struct {
	AutoApprovalPossible bool `json:"auto_approval_possible,omitempty"`
	MedicalExamRequired bool `json:"medical_exam_required,omitempty"`
}
