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

// GetGracePeriod Get grace period status
func (s *RenewalService) GetGracePeriod(ctx context.Context, policyId string) error {
	path := "/v1/policies/{policy_id}/grace-period"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetRenewalSchedule Get renewal schedule
func (s *RenewalService) GetRenewalSchedule(ctx context.Context, policyId string) error {
	path := "/v1/policies/{policy_id}/renewal-schedule"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RenewPolicy Renew policy manually
func (s *RenewalService) RenewPolicy(ctx context.Context, policyId string, req *models.PolicyRenewalRequest) error {
	path := "/v1/policies/{policy_id}:renew"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RevivePolicy Revive lapsed policy
func (s *RenewalService) RevivePolicy(ctx context.Context, policyId string, req *models.RevivePolicyRequest) error {
	path := "/v1/policies/{policy_id}:revive"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SendRenewalReminder Send renewal reminder
func (s *RenewalService) SendRenewalReminder(ctx context.Context, renewalScheduleId string, req *models.RenewalReminderSendingRequest) error {
	path := "/v1/renewal-schedules/{renewal_schedule_id}/reminders"
	path = strings.ReplaceAll(path, "{renewal_schedule_id}", renewalScheduleId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListUpcomingRenewals List upcoming renewals
func (s *RenewalService) ListUpcomingRenewals(ctx context.Context) error {
	path := "/v1/renewals/upcoming"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

