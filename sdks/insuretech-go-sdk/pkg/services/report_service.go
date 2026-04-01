package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// ReportService handles report-related API calls
type ReportService struct {
	Client Client
}

// ListReportDefinitions List report definitions
func (s *ReportService) ListReportDefinitions(ctx context.Context) error {
	path := "/v1/report-definitions"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ListReportExecutions List report executions
func (s *ReportService) ListReportExecutions(ctx context.Context) error {
	path := "/v1/report-executions"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetReportExecution Get report execution
func (s *ReportService) GetReportExecution(ctx context.Context, reportExecutionId string) error {
	path := "/v1/report-executions/{report_execution_id}"
	path = strings.ReplaceAll(path, "{report_execution_id}", reportExecutionId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// DownloadReport Download report
func (s *ReportService) DownloadReport(ctx context.Context, reportExecutionId string) error {
	path := "/v1/report-executions/{report_execution_id}/download"
	path = strings.ReplaceAll(path, "{report_execution_id}", reportExecutionId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ListReportSchedules List report schedules
func (s *ReportService) ListReportSchedules(ctx context.Context) error {
	path := "/v1/report-schedules"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateReportSchedule Create report schedule
func (s *ReportService) CreateReportSchedule(ctx context.Context, req *models.ReportScheduleCreationRequest) error {
	path := "/v1/report-schedules"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ExecuteReport Execute report
func (s *ReportService) ExecuteReport(ctx context.Context, reportDefinitionId string, req *models.ReportExecutionRequest) error {
	path := "/v1/reports/{report_definition_id}:execute"
	path = strings.ReplaceAll(path, "{report_definition_id}", reportDefinitionId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

