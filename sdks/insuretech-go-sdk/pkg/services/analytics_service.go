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

// CreateDashboard Create dashboard
func (s *AnalyticsService) CreateDashboard(ctx context.Context, req *models.DashboardCreationRequest) error {
	path := "/v1/analytics/dashboards"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetDashboard Get dashboard
func (s *AnalyticsService) GetDashboard(ctx context.Context, dashboardId string) error {
	path := "/v1/analytics/dashboards/{dashboard_id}"
	path = strings.ReplaceAll(path, "{dashboard_id}", dashboardId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetMetrics Get metrics
func (s *AnalyticsService) GetMetrics(ctx context.Context, req *models.MetricsRetrievalRequest) error {
	path := "/v1/analytics/metrics"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RunQuery Run custom query
func (s *AnalyticsService) RunQuery(ctx context.Context, req *models.RunQueryRequest) error {
	path := "/v1/analytics/queries:run"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GenerateReport Generate report
func (s *AnalyticsService) GenerateReport(ctx context.Context, reportId string, req *models.ReportGenerationRequest) error {
	path := "/v1/analytics/reports/{report_id}:generate"
	path = strings.ReplaceAll(path, "{report_id}", reportId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ScheduleReport Schedule report
func (s *AnalyticsService) ScheduleReport(ctx context.Context, req *models.ScheduleReportRequest) error {
	path := "/v1/analytics/reports:schedule"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

