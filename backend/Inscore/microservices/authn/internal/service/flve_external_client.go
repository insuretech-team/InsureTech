package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	"github.com/google/uuid"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"google.golang.org/grpc"
)

// flveExternalKYCClient adapts the existing ExternalKYCClient interface to
// FLVE's HTTP eKYC endpoints.
type flveExternalKYCClient struct {
	baseURL string
	token   string
	client  *http.Client
}

type flveStartResponse struct {
	SessionID string `json:"session_id"`
	State     string `json:"state"`
	Error     string `json:"error"`
}

type flveFrameResponse struct {
	SessionID      string  `json:"session_id"`
	SessionState   string  `json:"session_state"`
	StepCompleted  bool    `json:"step_completed"`
	CurrentStep    string  `json:"current_step"`
	CompletedSteps int32   `json:"completed_steps"`
	TotalSteps     int32   `json:"total_steps"`
	LivenessScore  float64 `json:"liveness_score"`
	LivenessConf   float64 `json:"liveness_confidence"`
	Error          string  `json:"error"`
}

type flveCompleteResponse struct {
	SessionID          string  `json:"session_id"`
	Success            bool    `json:"success"`
	State              string  `json:"state"`
	LivenessConfidence float64 `json:"liveness_confidence"`
	ProfileImageURL    string  `json:"profile_image_url"`
	Error              string  `json:"error"`
}

// NewFLVEExternalKYCClient creates an ExternalKYCClient backed by FLVE HTTP endpoints.
func NewFLVEExternalKYCClient(baseURL, token string, timeout time.Duration) ExternalKYCClient {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &flveExternalKYCClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *flveExternalKYCClient) StartKYCVerification(ctx context.Context, in *kycservicev1.StartKYCVerificationRequest, _ ...grpc.CallOption) (*kycservicev1.StartKYCVerificationResponse, error) {
	reqBody := map[string]interface{}{
		"user_id": in.GetEntityId(),
		"metadata": map[string]string{
			"entity_type":       in.GetEntityType(),
			"method":            in.GetMethod(),
			"verification_type": in.GetType(),
		},
	}

	var resp flveStartResponse
	if err := c.doJSON(ctx, http.MethodPost, "/ekyc/start", reqBody, &resp); err != nil {
		logger.Errorf("flve start eKYC: %v", err)
		return nil, errors.New("flve start eKYC")
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("flve start eKYC error: %s", resp.Error)
	}

	sessionID, err := normalizeFLVESessionID(resp.SessionID)
	if err != nil {
		return nil, err
	}

	return &kycservicev1.StartKYCVerificationResponse{
		KycVerificationId: sessionID,
		Message:           "FLVE eKYC session started",
	}, nil
}

func (c *flveExternalKYCClient) UploadDocument(ctx context.Context, in *kycservicev1.UploadDocumentRequest, _ ...grpc.CallOption) (*kycservicev1.UploadDocumentResponse, error) {
	frameBytes, err := decodeFramePayload(in.GetDocumentUrl())
	if err != nil {
		logger.Errorf("flve upload frame decode: %v", err)
		return nil, errors.New("flve upload frame decode")
	}

	query := url.Values{}
	query.Set("session_id", in.GetKycVerificationId())
	endpoint := "/ekyc/frame?" + query.Encode()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "frame.jpg")
	if err != nil {
		logger.Errorf("create multipart form: %v", err)
		return nil, errors.New("create multipart form")
	}
	if _, err := part.Write(frameBytes); err != nil {
		logger.Errorf("write multipart frame: %v", err)
		return nil, errors.New("write multipart frame")
	}
	if err := writer.Close(); err != nil {
		logger.Errorf("close multipart writer: %v", err)
		return nil, errors.New("close multipart writer")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, &body)
	if err != nil {
		logger.Errorf("build flve frame request: %v", err)
		return nil, errors.New("build flve frame request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.applyAuth(req)

	httpResp, err := c.client.Do(req)
	if err != nil {
		logger.Errorf("flve frame request failed: %v", err)
		return nil, errors.New("flve frame request failed")
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		logger.Errorf("read flve frame response: %v", err)
		return nil, errors.New("read flve frame response")
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("flve frame status %d: %s", httpResp.StatusCode, string(respBody))
	}

	var resp flveFrameResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		logger.Errorf("decode flve frame response: %v", err)
		return nil, errors.New("decode flve frame response")
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("flve frame error: %s", resp.Error)
	}

	docID := in.GetDocumentNumber()
	if docID == "" {
		docID = uuid.NewString()
	}

	return &kycservicev1.UploadDocumentResponse{
		DocumentVerificationId: docID,
		Message:                "FLVE frame accepted",
	}, nil
}

func (c *flveExternalKYCClient) VerifyKYC(ctx context.Context, in *kycservicev1.VerifyKYCRequest, _ ...grpc.CallOption) (*kycservicev1.VerifyKYCResponse, error) {
	query := url.Values{}
	query.Set("session_id", in.GetKycVerificationId())
	endpoint := "/ekyc/complete?" + query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, http.NoBody)
	if err != nil {
		logger.Errorf("build flve complete request: %v", err)
		return nil, errors.New("build flve complete request")
	}
	c.applyAuth(req)

	httpResp, err := c.client.Do(req)
	if err != nil {
		logger.Errorf("flve complete request failed: %v", err)
		return nil, errors.New("flve complete request failed")
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		logger.Errorf("read flve complete response: %v", err)
		return nil, errors.New("read flve complete response")
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("flve complete status %d: %s", httpResp.StatusCode, string(respBody))
	}

	var resp flveCompleteResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		logger.Errorf("decode flve complete response: %v", err)
		return nil, errors.New("decode flve complete response")
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("flve complete error: %s", resp.Error)
	}
	if !resp.Success {
		return nil, fmt.Errorf("flve complete unsuccessful (state=%s)", resp.State)
	}

	msg := "FLVE eKYC verified"
	if resp.ProfileImageURL != "" {
		msg = fmt.Sprintf("FLVE eKYC verified (profile_image_url=%s)", resp.ProfileImageURL)
	}
	return &kycservicev1.VerifyKYCResponse{
		Message: msg,
	}, nil
}

func (c *flveExternalKYCClient) doJSON(ctx context.Context, method, endpoint string, reqBody interface{}, out interface{}) error {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorf("marshal request: %v", err)
		return errors.New("marshal request")
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		logger.Errorf("create request: %v", err)
		return errors.New("create request")
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("read response: %v", err)
		return errors.New("read response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		logger.Errorf("decode response: %v", err)
		return errors.New("decode response")
	}
	return nil
}

func (c *flveExternalKYCClient) applyAuth(req *http.Request) {
	if strings.TrimSpace(c.token) != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}

func decodeFramePayload(data string) ([]byte, error) {
	s := strings.TrimSpace(data)
	if s == "" {
		logger.Errorf("empty document_url")
		return nil, errors.New("empty document_url")
	}
	if strings.HasPrefix(s, "data:") {
		idx := strings.Index(s, ",")
		if idx < 0 || idx == len(s)-1 {
			logger.Errorf("invalid data URL payload")
			return nil, errors.New("invalid data URL payload")
		}
		return base64.StdEncoding.DecodeString(s[idx+1:])
	}
	if decoded, err := base64.StdEncoding.DecodeString(s); err == nil && len(decoded) > 0 {
		return decoded, nil
	}
	logger.Errorf("unsupported frame payload format")
	return nil, errors.New("unsupported frame payload format")
}

func normalizeFLVESessionID(raw string) (string, error) {
	id := strings.TrimSpace(raw)
	if id == "" {
		logger.Errorf("flve returned empty session_id")
		return "", errors.New("flve returned empty session_id")
	}
	if u, err := uuid.Parse(id); err == nil {
		return u.String(), nil
	}

	legacy := strings.TrimPrefix(id, "ekyc_")
	if len(legacy) == 32 {
		legacy = fmt.Sprintf("%s-%s-%s-%s-%s", legacy[0:8], legacy[8:12], legacy[12:16], legacy[16:20], legacy[20:32])
	}
	if u, err := uuid.Parse(legacy); err == nil {
		return u.String(), nil
	}
	return "", fmt.Errorf("invalid flve session_id format: %q", raw)
}
