package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// RenewalService handles renewal-related API calls
type RenewalService struct {
	Client Client
}

// ListUpcomingRenewals List upcoming renewals
func (s *RenewalService) ListUpcomingRenewals(ctx context.Context) (*models.UpcomingRenewalsListingResponse, error) {
	path := "/v1/renewals/upcoming"
	var result models.UpcomingRenewalsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGracePeriod Get grace period status
func (s *RenewalService) GetGracePeriod(ctx context.Context, policyId string) (*models.RenewalGracePeriodRetrievalResponse, error) {
	path := "/v1/policies/{policy_id}/grace-period"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.RenewalGracePeriodRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RenewPolicy Renew policy manually
func (s *RenewalService) RenewPolicy(ctx context.Context, policyId string, req *models.RenewalPolicyRenewalRequest) (*models.RenewalPolicyRenewalResponse, error) {
	path := "/v1/policies/{policy_id}"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.RenewalPolicyRenewalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRenewalSchedule Get renewal schedule
func (s *RenewalService) GetRenewalSchedule(ctx context.Context, policyId string) (*models.RenewalRenewalScheduleRetrievalResponse, error) {
	path := "/v1/policies/{policy_id}/renewal-schedule"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.RenewalRenewalScheduleRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendRenewalReminder Send renewal reminder
func (s *RenewalService) SendRenewalReminder(ctx context.Context, renewalScheduleId string, req *models.RenewalReminderSendingRequest) (*models.RenewalReminderSendingResponse, error) {
	path := "/v1/renewal-schedules/{renewal_schedule_id}/reminders"
	path = strings.ReplaceAll(path, "{renewal_schedule_id}", renewalScheduleId)
	var result models.RenewalReminderSendingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

