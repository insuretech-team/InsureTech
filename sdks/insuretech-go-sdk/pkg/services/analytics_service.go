package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// AnalyticsService handles analytics-related API calls
type AnalyticsService struct {
	Client Client
}

// GetDashboard Get dashboard
func (s *AnalyticsService) GetDashboard(ctx context.Context, dashboardId string) (*models.DashboardRetrievalResponse, error) {
	path := "/v1/analytics/dashboards/{dashboard_id}"
	path = strings.ReplaceAll(path, "{dashboard_id}", dashboardId)
	var result models.DashboardRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMetrics Get metrics
func (s *AnalyticsService) GetMetrics(ctx context.Context, req *models.MetricsRetrievalRequest) (*models.MetricsRetrievalResponse, error) {
	path := "/v1/analytics/metrics"
	var result models.MetricsRetrievalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ScheduleReport Schedule report
func (s *AnalyticsService) ScheduleReport(ctx context.Context, req *models.ScheduleReportRequest) (*models.ScheduleReportResponse, error) {
	path := "/v1/analytics/reports:schedule"
	var result models.ScheduleReportResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateDashboard Create dashboard
func (s *AnalyticsService) CreateDashboard(ctx context.Context, req *models.DashboardCreationRequest) (*models.DashboardCreationResponse, error) {
	path := "/v1/analytics/dashboards"
	var result models.DashboardCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GenerateReport Generate report
func (s *AnalyticsService) GenerateReport(ctx context.Context, reportId string, req *models.ReportGenerationRequest) (*models.ReportGenerationResponse, error) {
	path := "/v1/analytics/reports/{report_id}:generate"
	path = strings.ReplaceAll(path, "{report_id}", reportId)
	var result models.ReportGenerationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RunQuery Run custom query
func (s *AnalyticsService) RunQuery(ctx context.Context, req *models.RunQueryRequest) (*models.RunQueryResponse, error) {
	path := "/v1/analytics/queries:run"
	var result models.RunQueryResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

