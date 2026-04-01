package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// AuditService handles audit-related API calls
type AuditService struct {
	Client Client
}

// GetAuditEvents Get audit events
func (s *AuditService) GetAuditEvents(ctx context.Context) error {
	path := "/v1/audit-events"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateAuditEvent Create audit event
func (s *AuditService) CreateAuditEvent(ctx context.Context, req *models.AuditEventCreationRequest) error {
	path := "/v1/audit-events"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetAuditLogs Get audit logs for entity
func (s *AuditService) GetAuditLogs(ctx context.Context) error {
	path := "/v1/audit-logs"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateAuditLog Create audit log
func (s *AuditService) CreateAuditLog(ctx context.Context, req *models.AuditLogCreationRequest) error {
	path := "/v1/audit-logs"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetComplianceLogs Get compliance logs
func (s *AuditService) GetComplianceLogs(ctx context.Context) error {
	path := "/v1/compliance-logs"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateComplianceLog Create compliance log
func (s *AuditService) CreateComplianceLog(ctx context.Context, req *models.ComplianceLogCreationRequest) error {
	path := "/v1/compliance-logs"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GenerateComplianceReport Generate compliance report
func (s *AuditService) GenerateComplianceReport(ctx context.Context, req *models.ComplianceReportGenerationRequest) error {
	path := "/v1/compliance-reports:generate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetAuditTrail Get audit trail for entity
func (s *AuditService) GetAuditTrail(ctx context.Context, entityType string, entityId string) error {
	path := "/v1/entities/{entity_type}/{entity_id}/audit-trail"
	path = strings.ReplaceAll(path, "{entity_type}", entityType)
	path = strings.ReplaceAll(path, "{entity_id}", entityId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

