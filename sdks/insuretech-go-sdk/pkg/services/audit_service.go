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

// GenerateComplianceReport Generate compliance report
func (s *AuditService) GenerateComplianceReport(ctx context.Context, req *models.ComplianceReportGenerationRequest) (*models.ComplianceReportGenerationResponse, error) {
	path := "/v1/compliance-reports"
	var result models.ComplianceReportGenerationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateAuditEvent Create audit event
func (s *AuditService) CreateAuditEvent(ctx context.Context, req *models.AuditEventCreationRequest) (*models.AuditEventCreationResponse, error) {
	path := "/v1/audit-events"
	var result models.AuditEventCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAuditEvents Get audit events
func (s *AuditService) GetAuditEvents(ctx context.Context) (*models.AuditEventsRetrievalResponse, error) {
	path := "/v1/audit-events"
	var result models.AuditEventsRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateAuditLog Create audit log
func (s *AuditService) CreateAuditLog(ctx context.Context, req *models.AuditLogCreationRequest) (*models.AuditLogCreationResponse, error) {
	path := "/v1/audit-logs"
	var result models.AuditLogCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAuditLogs Get audit logs for entity
func (s *AuditService) GetAuditLogs(ctx context.Context) (*models.AuditLogsRetrievalResponse, error) {
	path := "/v1/audit-logs"
	var result models.AuditLogsRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateComplianceLog Create compliance log
func (s *AuditService) CreateComplianceLog(ctx context.Context, req *models.ComplianceLogCreationRequest) (*models.ComplianceLogCreationResponse, error) {
	path := "/v1/compliance-logs"
	var result models.ComplianceLogCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetComplianceLogs Get compliance logs
func (s *AuditService) GetComplianceLogs(ctx context.Context) (*models.ComplianceLogsRetrievalResponse, error) {
	path := "/v1/compliance-logs"
	var result models.ComplianceLogsRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAuditTrail Get audit trail for entity
func (s *AuditService) GetAuditTrail(ctx context.Context, entityType string, entityId string) (*models.AuditTrailRetrievalResponse, error) {
	path := "/v1/entities/{entity_type}/{entity_id}/audit-trail"
	path = strings.ReplaceAll(path, "{entity_type}", entityType)
	path = strings.ReplaceAll(path, "{entity_id}", entityId)
	var result models.AuditTrailRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

