package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// IotService handles iot-related API calls
type IotService struct {
	Client Client
}

// RegisterDevice Register device
func (s *IotService) RegisterDevice(ctx context.Context, req *models.DeviceRegistrationRequest) error {
	path := "/v1/iot/devices"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetDeviceStatus Get device status
func (s *IotService) GetDeviceStatus(ctx context.Context, deviceId string) error {
	path := "/v1/iot/devices/{device_id}"
	path = strings.ReplaceAll(path, "{device_id}", deviceId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetRiskAssessment Get risk assessment
func (s *IotService) GetRiskAssessment(ctx context.Context, deviceId string) error {
	path := "/v1/iot/devices/{device_id}/risk"
	path = strings.ReplaceAll(path, "{device_id}", deviceId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// DeactivateDevice Deactivate device
func (s *IotService) DeactivateDevice(ctx context.Context, deviceId string, req *models.DeviceDeactivationRequest) error {
	path := "/v1/iot/devices/{device_id}:deactivate"
	path = strings.ReplaceAll(path, "{device_id}", deviceId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SendTelemetry Send telemetry data
func (s *IotService) SendTelemetry(ctx context.Context, req *models.TelemetrySendingRequest) error {
	path := "/v1/iot/telemetry"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

