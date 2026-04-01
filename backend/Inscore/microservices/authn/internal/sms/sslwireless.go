package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/mobile"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
)

// SSLWirelessClient implements SMS sending via SSL Wireless API
type SSLWirelessClient struct {
	config     *config.Config
	httpClient *http.Client
}

// SendSMSRequest represents an SMS send request
type SendSMSRequest struct {
	MSISDN     string // 8801XXXXXXXXX format
	Message    string
	UseMasking bool
	CSMSId     string // Optional client message ID
}

// SendSMSResponse represents SSL Wireless API response
type SendSMSResponse struct {
	MessageID string
	Status    string
	ErrorCode string
	ErrorMsg  string
}

// DLRStatus represents delivery report status
type DLRStatus string

const (
	DLRStatusPending   DLRStatus = "PENDING"
	DLRStatusDelivered DLRStatus = "DELIVERED"
	DLRStatusFailed    DLRStatus = "FAILED"
	DLRStatusExpired   DLRStatus = "EXPIRED"
)

// DLRWebhookPayload represents incoming DLR webhook data
type DLRWebhookPayload struct {
	MessageID   string    `json:"message_id"`
	MSISDN      string    `json:"msisdn"`
	Status      string    `json:"status"`
	ErrorCode   string    `json:"error_code"`
	DeliveredAt time.Time `json:"delivered_at"`
	Carrier     string    `json:"carrier"`
}

// NewSSLWirelessClient creates a new SSL Wireless SMS client
func NewSSLWirelessClient(cfg *config.Config) *SSLWirelessClient {
	return &SSLWirelessClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendSMS sends an SMS via SSL Wireless
// Supports both masking (alphanumeric sender) and non-masking (numeric sender)
func (c *SSLWirelessClient) SendSMS(ctx context.Context, req *SendSMSRequest) (*SendSMSResponse, error) {
	// Normalize MSISDN to 8801XXXXXXXXX format
	msisdn := NormalizeMSISDN(req.MSISDN)
	if msisdn == "" {
		return nil, errors.New("invalid MSISDN format: " + req.MSISDN)
	}

	// Determine sender based on masking preference
	var sender string

	if req.UseMasking && c.config.SMS.MaskingEnabled {
		sender = c.config.SMS.MaskingSenderID // e.g., "LABAIDINS"
	} else if c.config.SMS.NonMaskingEnabled {
		sender = c.config.SMS.NonMaskingSender // Numeric sender
	} else {
		return nil, errors.New("no SMS sender configured")
	}

	// Build API request
	apiURL := c.config.SMS.APIBase + "/api/v3/send-sms"

	payload := map[string]interface{}{
		"api_token": c.config.SMS.APIKey,
		"sid":       c.config.SMS.SID,
		"msisdn":    msisdn,
		"sms":       req.Message,
		"sender_id": sender,
		"csms_id":   req.CSMSId,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("failed to marshal request: %v", err)
		return nil, errors.New("failed to marshal request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Errorf("failed to create HTTP request: %v", err)
		return nil, errors.New("failed to create HTTP request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		logger.Errorf("failed to send HTTP request: %v", err)
		return nil, errors.New("failed to send HTTP request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("failed to read response body: %v", err)
		return nil, errors.New("failed to read response body")
	}

	// Parse response — SSL Wireless v3 API format:
	// {"status":"SUCCESS","status_code":200,"error_message":"","smsinfo":[{"sms_status":"SUCCESS","reference_id":"..."}]}
	var apiResp struct {
		Status       string `json:"status"`
		StatusCode   int    `json:"status_code"`
		ErrorMessage string `json:"error_message"`
		SMSInfo      []struct {
			SMSStatus     string `json:"sms_status"`
			StatusMessage string `json:"status_message"`
			MSISDN        string `json:"msisdn"`
			CSMSId        string `json:"csms_id"`
			ReferenceID   string `json:"reference_id"`
		} `json:"smsinfo"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		logger.Errorf("failed to parse SSL Wireless response: %v | body: %s", err, string(body))
		return nil, errors.New("failed to parse response")
	}

	// Check for errors — API returns STATUS in uppercase
	if apiResp.Status != "SUCCESS" {
		logger.Errorf("SSL Wireless send failed: status=%s error=%s body=%s", apiResp.Status, apiResp.ErrorMessage, string(body))
		return &SendSMSResponse{
			Status:    "FAILED",
			ErrorCode: apiResp.ErrorMessage,
			ErrorMsg:  apiResp.ErrorMessage,
		}, errors.New("SMS send failed: " + apiResp.ErrorMessage)
	}

	// Also check individual SMS status
	if len(apiResp.SMSInfo) == 0 || apiResp.SMSInfo[0].SMSStatus != "SUCCESS" {
		errMsg := "unknown error"
		if len(apiResp.SMSInfo) > 0 {
			errMsg = apiResp.SMSInfo[0].StatusMessage
		}
		logger.Errorf("SSL Wireless SMS status failed: %s", errMsg)
		return &SendSMSResponse{
			Status:   "FAILED",
			ErrorMsg: errMsg,
		}, errors.New("SMS send failed: " + errMsg)
	}

	refID := apiResp.SMSInfo[0].ReferenceID
	logger.Infof("SSL Wireless SMS sent: reference_id=%s msisdn=%s", refID, msisdn)

	return &SendSMSResponse{
		MessageID: refID,
		Status:    "PENDING",
	}, nil
}

// ParseDLRWebhook parses incoming DLR webhook payload
func (c *SSLWirelessClient) ParseDLRWebhook(payload []byte) (*DLRWebhookPayload, error) {
	var dlr DLRWebhookPayload
	if err := json.Unmarshal(payload, &dlr); err != nil {
		logger.Errorf("failed to parse DLR payload: %v", err)
		return nil, errors.New("failed to parse DLR payload")
	}
	return &dlr, nil
}

// normalizeMSISDN normalizes phone numbers to 8801XXXXXXXXX format
// Accepts: +8801XXXXXXXXX, 8801XXXXXXXXX, 01XXXXXXXXX
// NormalizeMSISDN normalizes phone numbers to 8801XXXXXXXXX format
// Accepts: +8801XXXXXXXXX, 8801XXXXXXXXX, 01XXXXXXXXX
func NormalizeMSISDN(phone string) string {
	normalized, err := mobile.NormalizeBangladeshMobileDigits(phone)
	if err != nil {
		return ""
	}
	return normalized
}

// detectCarrier detects mobile carrier from MSISDN prefix
// GP: 8801[3,7,8]
// Robi: 8801[6,8]
// Banglalink: 8801[9]
// Teletalk: 8801[5]
// DetectCarrier detects mobile carrier from MSISDN prefix
// GP: 8801[3,7,8], Robi: 8801[6,8], Banglalink: 8801[9], Teletalk: 8801[5]
func DetectCarrier(msisdn string) string {
	if len(msisdn) < 6 {
		return "UNKNOWN"
	}

	// Get the 5th digit (X in 8801X)
	digit := string(msisdn[4])

	switch digit {
	case "3", "7":
		return "GP"
	case "8":
		// 8 is used by both GP and Robi, check 6th digit
		if len(msisdn) > 5 && msisdn[5] >= '0' && msisdn[5] <= '9' {
			return "GP" // Assume GP for simplicity
		}
		return "GP"
	case "6":
		return "ROBI"
	case "9":
		return "BANGLALINK"
	case "5":
		return "TELETALK"
	default:
		return "UNKNOWN"
	}
}

// ValidateMSISDN validates if a phone number is in valid Bangladesh format
func ValidateMSISDN(phone string) bool {
	normalized := NormalizeMSISDN(phone)
	return normalized != ""
}

// MaskMSISDN masks a phone number for logging/display
// Example: 8801712345678 -> 8801XXX***678
func MaskMSISDN(msisdn string) string {
	if len(msisdn) < 13 {
		return "XXXXXXXXXXXXX"
	}
	return msisdn[:4] + "XXX***" + msisdn[len(msisdn)-3:]
}
