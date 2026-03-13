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

// DeactivateDevice Deactivate device
func (s *IotService) DeactivateDevice(ctx context.Context, deviceId string, req *models.DeviceDeactivationRequest) (*models.DeviceDeactivationResponse, error) {
	path := "/v1/iot/devices/{device_id}:deactivate"
	path = strings.ReplaceAll(path, "{device_id}", deviceId)
	var result models.DeviceDeactivationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RegisterDevice Register device
func (s *IotService) RegisterDevice(ctx context.Context, req *models.DeviceRegistrationRequest) (*models.DeviceRegistrationResponse, error) {
	path := "/v1/iot/devices"
	var result models.DeviceRegistrationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDeviceStatus Get device status
func (s *IotService) GetDeviceStatus(ctx context.Context, deviceId string) (*models.DeviceStatusRetrievalResponse, error) {
	path := "/v1/iot/devices/{device_id}"
	path = strings.ReplaceAll(path, "{device_id}", deviceId)
	var result models.DeviceStatusRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendTelemetry Send telemetry data
func (s *IotService) SendTelemetry(ctx context.Context, req *models.TelemetrySendingRequest) (*models.TelemetrySendingResponse, error) {
	path := "/v1/iot/telemetry"
	var result models.TelemetrySendingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRiskAssessment Get risk assessment
func (s *IotService) GetRiskAssessment(ctx context.Context, deviceId string) (*models.RiskAssessmentRetrievalResponse, error) {
	path := "/v1/iot/devices/{device_id}/risk"
	path = strings.ReplaceAll(path, "{device_id}", deviceId)
	var result models.RiskAssessmentRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

