package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrPushNotConfigured = errors.New("push delivery requires device token registry and FCM/APNs integration")
var ErrPushNoActiveDeviceTokens = errors.New("push recipient has no active device tokens")
var ErrPushAllTargetsRejected = errors.New("push provider rejected all device tokens")

type PushConfig struct {
	Provider    string
	Endpoint    string
	ServerKey   string
	Timeout     time.Duration
	MockSuccess bool
}

type PushTarget struct {
	Provider    string
	Platform    string
	DeviceToken string
	DeviceID    string
	AppID       string
}

type PushRequest struct {
	RecipientID string
	Title       string
	Body        string
	Data        map[string]string
	Targets     []PushTarget
}

type PushResponse struct {
	MessageID     string
	Status        string
	SentAt        time.Time
	SentCount     int
	FailureCount  int
	InvalidTokens []string
}

type PushClient struct {
	config     PushConfig
	httpClient *http.Client
}

func NewPushClient(cfg PushConfig) *PushClient {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &PushClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *PushClient) Send(ctx context.Context, req *PushRequest) (*PushResponse, error) {
	if req == nil || strings.TrimSpace(req.RecipientID) == "" {
		return nil, Permanent(errors.New("push recipient is required"))
	}
	targets := c.supportedTargets(req.Targets)
	if len(targets) == 0 {
		return nil, Permanent(ErrPushNoActiveDeviceTokens)
	}

	switch strings.ToLower(strings.TrimSpace(c.config.Provider)) {
	case "mock":
		return &PushResponse{
			MessageID: uuid.NewString(),
			Status:    "DELIVERED",
			SentAt:    time.Now().UTC(),
			SentCount: len(targets),
		}, nil
	case "fcm":
		return c.sendFCM(ctx, req, targets)
	default:
		return nil, Permanent(ErrPushNotConfigured)
	}
}

func (c *PushClient) supportedTargets(targets []PushTarget) []PushTarget {
	if len(targets) == 0 {
		return nil
	}
	provider := strings.ToUpper(strings.TrimSpace(c.config.Provider))
	if provider == "" || provider == "MOCK" {
		return targets
	}
	filtered := make([]PushTarget, 0, len(targets))
	for _, target := range targets {
		targetProvider := strings.ToUpper(strings.TrimSpace(target.Provider))
		if targetProvider == "" || targetProvider == provider {
			filtered = append(filtered, target)
		}
	}
	return filtered
}

func (c *PushClient) sendFCM(ctx context.Context, req *PushRequest, targets []PushTarget) (*PushResponse, error) {
	if strings.TrimSpace(c.config.Endpoint) == "" || strings.TrimSpace(c.config.ServerKey) == "" {
		return nil, Permanent(ErrPushNotConfigured)
	}

	registrationIDs := make([]string, 0, len(targets))
	for _, target := range targets {
		if trimmed := strings.TrimSpace(target.DeviceToken); trimmed != "" {
			registrationIDs = append(registrationIDs, trimmed)
		}
	}
	if len(registrationIDs) == 0 {
		return nil, Permanent(ErrPushNoActiveDeviceTokens)
	}

	payload := map[string]any{
		"registration_ids": registrationIDs,
		"priority":         "high",
		"notification": map[string]string{
			"title": req.Title,
			"body":  req.Body,
		},
		"data": cloneStringMap(req.Data),
	}
	if c.config.MockSuccess {
		return &PushResponse{
			MessageID: uuid.NewString(),
			Status:    "SENT",
			SentAt:    time.Now().UTC(),
			SentCount: len(registrationIDs),
		}, nil
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, Permanent(fmt.Errorf("marshal push payload: %w", err))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, Permanent(fmt.Errorf("create push request: %w", err))
	}
	httpReq.Header.Set("Authorization", "key="+c.config.ServerKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send push request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return nil, fmt.Errorf("read push response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusRequestTimeout || resp.StatusCode >= 500 {
			return nil, fmt.Errorf("push provider returned status %d", resp.StatusCode)
		}
		return nil, Permanent(fmt.Errorf("push provider returned status %d", resp.StatusCode))
	}

	var apiResp struct {
		Success int `json:"success"`
		Failure int `json:"failure"`
		Results []struct {
			MessageID string `json:"message_id"`
			Error     string `json:"error"`
		} `json:"results"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parse push response: %w", err)
	}

	response := &PushResponse{
		Status:       "SENT",
		SentAt:       time.Now().UTC(),
		SentCount:    apiResp.Success,
		FailureCount: apiResp.Failure,
	}
	for idx, result := range apiResp.Results {
		if response.MessageID == "" && strings.TrimSpace(result.MessageID) != "" {
			response.MessageID = result.MessageID
		}
		if !isInvalidFCMTokenError(result.Error) || idx >= len(registrationIDs) {
			continue
		}
		response.InvalidTokens = append(response.InvalidTokens, registrationIDs[idx])
	}

	if apiResp.Success > 0 {
		return response, nil
	}
	if len(response.InvalidTokens) == len(registrationIDs) && len(registrationIDs) > 0 {
		return response, Permanent(ErrPushAllTargetsRejected)
	}
	return response, fmt.Errorf("push provider did not accept any device tokens")
}

func isInvalidFCMTokenError(value string) bool {
	switch strings.TrimSpace(value) {
	case "InvalidRegistration", "NotRegistered", "MismatchSenderId":
		return true
	default:
		return false
	}
}

func cloneStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return map[string]string{}
	}
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}
