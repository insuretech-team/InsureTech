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

// ListReportExecutions List report executions
func (s *ReportService) ListReportExecutions(ctx context.Context) (*models.ReportExecutionsListingResponse, error) {
	path := "/v1/report-executions"
	var result models.ReportExecutionsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListReportDefinitions List report definitions
func (s *ReportService) ListReportDefinitions(ctx context.Context) (*models.ReportDefinitionsListingResponse, error) {
	path := "/v1/report-definitions"
	var result models.ReportDefinitionsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetReportExecution Get report execution
func (s *ReportService) GetReportExecution(ctx context.Context, reportExecutionId string) (*models.ReportExecutionRetrievalResponse, error) {
	path := "/v1/report-executions/{report_execution_id}"
	path = strings.ReplaceAll(path, "{report_execution_id}", reportExecutionId)
	var result models.ReportExecutionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecuteReport Execute report
func (s *ReportService) ExecuteReport(ctx context.Context, reportDefinitionId string, req *models.ReportExecutionRequest) (*models.ReportExecutionResponse, error) {
	path := "/v1/reports/{report_definition_id}"
	path = strings.ReplaceAll(path, "{report_definition_id}", reportDefinitionId)
	var result models.ReportExecutionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateReportSchedule Create report schedule
func (s *ReportService) CreateReportSchedule(ctx context.Context, req *models.ReportScheduleCreationRequest) (*models.ReportScheduleCreationResponse, error) {
	path := "/v1/report-schedules"
	var result models.ReportScheduleCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListReportSchedules List report schedules
func (s *ReportService) ListReportSchedules(ctx context.Context) (*models.ReportSchedulesListingResponse, error) {
	path := "/v1/report-schedules"
	var result models.ReportSchedulesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadReport Download report
func (s *ReportService) DownloadReport(ctx context.Context, reportExecutionId string) (*models.ReportDownloadResponse, error) {
	path := "/v1/report-executions/{report_execution_id}/download"
	path = strings.ReplaceAll(path, "{report_execution_id}", reportExecutionId)
	var result models.ReportDownloadResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

