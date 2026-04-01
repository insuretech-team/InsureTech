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

	"github.com/newage-saint/insuretech/backend/inscore/pkg/mobile"
)

type SMSConfig struct {
	Provider          string
	APIBase           string
	SID               string
	APIKey            string
	MaskingEnabled    bool
	MaskingSenderID   string
	NonMaskingEnabled bool
	NonMaskingSender  string
}

type SMSClient struct {
	config     SMSConfig
	httpClient *http.Client
}

type SMSRequest struct {
	MSISDN     string
	Message    string
	UseMasking bool
	CSMSID     string
}

type SMSResponse struct {
	MessageID string
	Status    string
	ErrorCode string
	ErrorMsg  string
}

func NewSMSClient(cfg SMSConfig) *SMSClient {
	return &SMSClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *SMSClient) Send(ctx context.Context, req *SMSRequest) (*SMSResponse, error) {
	if strings.TrimSpace(req.MSISDN) == "" {
		return nil, errors.New("sms recipient is required")
	}
	if strings.TrimSpace(req.Message) == "" {
		return nil, errors.New("sms message is required")
	}
	if strings.TrimSpace(c.config.APIBase) == "" {
		return nil, errors.New("sms delivery is not configured")
	}

	msisdn, err := mobile.NormalizeBangladeshMobileDigits(req.MSISDN)
	if err != nil {
		return nil, fmt.Errorf("invalid sms recipient: %w", err)
	}

	sender, err := c.resolveSender(req.UseMasking)
	if err != nil {
		return nil, err
	}

	payload := map[string]any{
		"api_token": c.config.APIKey,
		"sid":       c.config.SID,
		"msisdn":    msisdn,
		"sms":       req.Message,
		"sender_id": sender,
		"csms_id":   req.CSMSID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal SMS request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.config.APIBase, "/")+"/api/v3/send-sms", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create SMS request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send SMS request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read SMS response: %w", err)
	}

	var apiResp struct {
		Status       string `json:"status"`
		StatusCode   int    `json:"status_code"`
		ErrorMessage string `json:"error_message"`
		SMSInfo      []struct {
			SMSStatus     string `json:"sms_status"`
			StatusMessage string `json:"status_message"`
			ReferenceID   string `json:"reference_id"`
		} `json:"smsinfo"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse SMS response: %w", err)
	}
	if apiResp.Status != "SUCCESS" {
		return nil, fmt.Errorf("sms failed: %s", strings.TrimSpace(apiResp.ErrorMessage))
	}
	if len(apiResp.SMSInfo) == 0 || apiResp.SMSInfo[0].SMSStatus != "SUCCESS" {
		errMsg := "provider did not accept SMS"
		if len(apiResp.SMSInfo) > 0 && strings.TrimSpace(apiResp.SMSInfo[0].StatusMessage) != "" {
			errMsg = apiResp.SMSInfo[0].StatusMessage
		}
		return nil, fmt.Errorf("sms failed: %s", errMsg)
	}

	return &SMSResponse{
		MessageID: apiResp.SMSInfo[0].ReferenceID,
		Status:    "PENDING",
	}, nil
}

func (c *SMSClient) resolveSender(useMasking bool) (string, error) {
	if useMasking && c.config.MaskingEnabled && strings.TrimSpace(c.config.MaskingSenderID) != "" {
		return c.config.MaskingSenderID, nil
	}
	if c.config.NonMaskingEnabled && strings.TrimSpace(c.config.NonMaskingSender) != "" {
		return c.config.NonMaskingSender, nil
	}
	return "", errors.New("no sms sender configured")
}
