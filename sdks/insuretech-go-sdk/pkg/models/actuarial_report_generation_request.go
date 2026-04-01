package models

import (
	"time"
)

// ActuarialReportGenerationRequest represents a actuarial_report_generation_request
type ActuarialReportGenerationRequest struct {
	Format string `json:"format,omitempty"`
	GeneratedBy string `json:"generated_by,omitempty"`
	LineOfBusiness string `json:"line_of_business,omitempty"`
	PeriodEnd time.Time `json:"period_end,omitempty"`
	PeriodStart time.Time `json:"period_start,omitempty"`
	ProductId string `json:"product_id"`
	ReportType string `json:"report_type,omitempty"`
	SegmentationFields []string `json:"segmentation_fields,omitempty"`
}
