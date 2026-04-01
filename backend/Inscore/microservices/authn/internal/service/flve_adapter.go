package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// ── FLVEAdapter Interface ────────────────────────────────────────────────────

// FLVEAdapter is the purpose-built interface for the FLVE HuggingFace eKYC
// endpoints. Unlike the old ExternalKYCClient (which forced FLVE into a
// generic KYC shape), this interface preserves the full FLVE response data.
type FLVEAdapter interface {
	StartEKYC(ctx context.Context, req *FLVEStartRequest) (*FLVEStartResponse, error)
	SubmitEKYCFrame(ctx context.Context, sessionID string, imageData []byte) (*FLVEFrameResponse, error)
	CompleteEKYC(ctx context.Context, sessionID string) (*FLVECompleteResponse, error)
	GetEKYCStatus(ctx context.Context, sessionID string) (*FLVEStatusResponse, error)
}

// ── Request / Response Types ─────────────────────────────────────────────────
// Hand-mapped from authoritative Pydantic models in
// LabaidAi-Retina/deployments/huggingface/src/models/ekyc_schemas_proto.py

type FLVEStartRequest struct {
	UserID            string            `json:"user_id"`
	TenantID          string            `json:"tenant_id,omitempty"`
	UserType          string            `json:"user_type,omitempty"`
	Portal            string            `json:"portal,omitempty"`
	KYCVerificationID string            `json:"kyc_verification_id,omitempty"`
	ReferenceImageURL string            `json:"reference_image_url,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

type FLVEStartResponse struct {
	SessionID           string     `json:"session_id"`
	Steps               []FLVEStep `json:"steps"`
	TotalTimeoutSeconds int        `json:"total_timeout_seconds"`
	State               string     `json:"state"`
	Error               string     `json:"error,omitempty"`
}

type FLVEStep struct {
	StepNumber     int     `json:"step_number"`
	Type           string  `json:"type"`
	State          string  `json:"state"`
	Instruction    string  `json:"instruction"`
	InstructionKey string  `json:"instruction_key"`
	TimeoutSeconds int     `json:"timeout_seconds"`
	Confidence     float64 `json:"confidence"`
}

type FLVEEyeContours struct {
	Left  *FLVEEyeContour `json:"left,omitempty"`
	Right *FLVEEyeContour `json:"right,omitempty"`
}

type FLVEEyeContour struct {
	Edges  [][2]int                    `json:"edges"`
	Points map[string]FLVEContourPoint `json:"points"`
}

type FLVEContourPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type FLVEFrameResponse struct {
	SessionID       string                 `json:"session_id"`
	SessionState    string                 `json:"session_state"`
	CurrentStep     *FLVEStep              `json:"current_step,omitempty"`
	NextStep        *FLVEStep              `json:"next_step,omitempty"`
	StepCompleted   bool                   `json:"step_completed"`
	StepProgress    float64                `json:"step_progress"`
	OverallProgress float64                `json:"overall_progress"`
	Detection       map[string]interface{} `json:"detection,omitempty"`
	HeadPose        *FLVEHeadPose          `json:"head_pose,omitempty"`
	EyeState        *FLVEEyeState          `json:"eye_state,omitempty"`
	EyeContours     *FLVEEyeContours       `json:"eye_contours,omitempty"`
	LivenessScore   float64                `json:"liveness_score"`
	Guidance        []string               `json:"guidance"`
	Error           string                 `json:"error,omitempty"`
}

type FLVECompleteResponse struct {
	SessionID          string              `json:"session_id"`
	Success            bool                `json:"success"`
	State              string              `json:"state"`
	ProfileImageURL    string              `json:"profile_image_url,omitempty"`
	ProfileImageID     string              `json:"profile_image_id,omitempty"`
	CapturedImageB64   string              `json:"captured_image_base64,omitempty"`
	Embedding          []float64           `json:"embedding"`
	LivenessConfidence float64             `json:"liveness_confidence"`
	IdentityMatch      bool                `json:"identity_match"`
	MatchScore         float64             `json:"match_score"`
	Summary            *FLVESessionSummary `json:"summary,omitempty"`
	CompletedAt        string              `json:"completed_at,omitempty"`
	Error              string              `json:"error,omitempty"`
}

type FLVEStatusResponse struct {
	SessionID        string     `json:"session_id"`
	State            string     `json:"state"`
	Steps            []FLVEStep `json:"steps"`
	OverallProgress  float64    `json:"overall_progress"`
	ElapsedSeconds   int        `json:"elapsed_seconds"`
	RemainingSeconds int        `json:"remaining_seconds"`
	Error            string     `json:"error,omitempty"`
}

type FLVEHeadPose struct {
	Yaw   float64 `json:"yaw"`
	Pitch float64 `json:"pitch"`
	Roll  float64 `json:"roll"`
}

type FLVEEyeState struct {
	LeftOpenness  float64 `json:"left_openness"`
	RightOpenness float64 `json:"right_openness"`
	IsBlinking    bool    `json:"is_blinking"`
}

type FLVESessionSummary struct {
	TotalSteps           int              `json:"total_steps"`
	CompletedSteps       int              `json:"completed_steps"`
	FailedSteps          int              `json:"failed_steps"`
	TotalFramesProcessed int              `json:"total_frames_processed"`
	ElapsedMs            int              `json:"elapsed_ms"`
	StepResults          []FLVEStepResult `json:"step_results"`
}

type FLVEStepResult struct {
	Type            string  `json:"type"`
	State           string  `json:"state"`
	Confidence      float64 `json:"confidence"`
	FramesProcessed int     `json:"frames_processed"`
	ElapsedMs       int     `json:"elapsed_ms"`
}

// ── HTTP Adapter Implementation ──────────────────────────────────────────────

type flveHTTPAdapter struct {
	baseURL    string
	token      string
	client     *http.Client
	maxRetries int
	retryBase  time.Duration
}

// NewFLVEAdapter creates an FLVEAdapter backed by FLVE HTTP endpoints.
func NewFLVEAdapter(baseURL, token string, timeout time.Duration) FLVEAdapter {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &flveHTTPAdapter{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		client: &http.Client{
			Timeout: timeout,
		},
		maxRetries: 3,
		retryBase:  2 * time.Second,
	}
}

func (a *flveHTTPAdapter) StartEKYC(ctx context.Context, req *FLVEStartRequest) (*FLVEStartResponse, error) {
	var resp FLVEStartResponse
	if err := a.doJSONWithRetry(ctx, http.MethodPost, "/ekyc/start", req, &resp); err != nil {
		return nil, fmt.Errorf("flve start ekyc: %w", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("flve start ekyc error: %s", resp.Error)
	}
	return &resp, nil
}

func (a *flveHTTPAdapter) SubmitEKYCFrame(ctx context.Context, sessionID string, imageData []byte) (*FLVEFrameResponse, error) {
	query := url.Values{}
	query.Set("session_id", sessionID)
	endpoint := "/ekyc/frame?" + query.Encode()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "frame.jpg")
	if err != nil {
		return nil, fmt.Errorf("create multipart form: %w", err)
	}
	if _, err := part.Write(imageData); err != nil {
		return nil, fmt.Errorf("write multipart frame: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+endpoint, &body)
	if err != nil {
		return nil, fmt.Errorf("build flve frame request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	a.applyAuth(req)

	httpResp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("flve frame request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read flve frame response: %w", err)
	}
	if isHTMLResponse(httpResp, respBody) {
		switch httpResp.StatusCode {
		case 401:
			return nil, fmt.Errorf("FLVE authentication failed (HTTP 401): check FLVE_API_TOKEN matches the HuggingFace Space token")
		case 403:
			return nil, fmt.Errorf("FLVE access forbidden (HTTP 403): check FLVE_API_TOKEN matches the HuggingFace Space token")
		case 404:
			return nil, fmt.Errorf("FLVE endpoint not found (HTTP 404): Space may be sleeping or FLVE_API_TOKEN is wrong")
		default:
			return nil, fmt.Errorf("FLVE returned non-JSON response (HTTP %d): Space may be unavailable", httpResp.StatusCode)
		}
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("flve frame status %d: %s", httpResp.StatusCode, string(respBody))
	}

	var resp FLVEFrameResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decode flve frame response: %w", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("flve frame error: %s", resp.Error)
	}
	return &resp, nil
}

func (a *flveHTTPAdapter) CompleteEKYC(ctx context.Context, sessionID string) (*FLVECompleteResponse, error) {
	query := url.Values{}
	query.Set("session_id", sessionID)
	endpoint := "/ekyc/complete?" + query.Encode()

	var resp FLVECompleteResponse
	if err := a.doRawWithRetry(ctx, http.MethodPost, endpoint, nil, &resp); err != nil {
		return nil, fmt.Errorf("flve complete ekyc: %w", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("flve complete ekyc error: %s", resp.Error)
	}
	return &resp, nil
}

func (a *flveHTTPAdapter) GetEKYCStatus(ctx context.Context, sessionID string) (*FLVEStatusResponse, error) {
	endpoint := "/ekyc/status/" + url.PathEscape(sessionID)
	var resp FLVEStatusResponse
	if err := a.doJSONWithRetry(ctx, http.MethodGet, endpoint, nil, &resp); err != nil {
		return nil, fmt.Errorf("flve get ekyc status: %w", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("flve get ekyc status error: %s", resp.Error)
	}
	return &resp, nil
}

// ── HTTP helpers ─────────────────────────────────────────────────────────────

func (a *flveHTTPAdapter) applyAuth(req *http.Request) {
	if strings.TrimSpace(a.token) != "" {
		req.Header.Set("Authorization", "Bearer "+a.token)
	}
}

func (a *flveHTTPAdapter) doJSONWithRetry(ctx context.Context, method, endpoint string, reqBody interface{}, out interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= a.maxRetries; attempt++ {
		if attempt > 0 {
			delay := a.retryBase * time.Duration(math.Pow(2, float64(attempt-1)))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
		err := a.doJSON(ctx, method, endpoint, reqBody, out)
		if err == nil {
			return nil
		}
		lastErr = err
		if !isRetryable(err) {
			return err
		}
		logger.Warnf("flve request retry %d/%d: %v", attempt+1, a.maxRetries, err)
	}
	return lastErr
}

func (a *flveHTTPAdapter) doRawWithRetry(ctx context.Context, method, endpoint string, reqBody io.Reader, out interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= a.maxRetries; attempt++ {
		if attempt > 0 {
			delay := a.retryBase * time.Duration(math.Pow(2, float64(attempt-1)))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, a.baseURL+endpoint, reqBody)
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}
		a.applyAuth(req)

		resp, err := a.client.Do(req)
		if err != nil {
			lastErr = err
			if !isRetryable(err) {
				return err
			}
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response: %w", err)
		}
		if isHTMLResponse(resp, body) {
			switch resp.StatusCode {
			case 401:
				return fmt.Errorf("FLVE authentication failed (HTTP 401): check FLVE_API_TOKEN matches the HuggingFace Space token")
			case 403:
				return fmt.Errorf("FLVE access forbidden (HTTP 403): check FLVE_API_TOKEN matches the HuggingFace Space token")
			case 404:
				return fmt.Errorf("FLVE endpoint not found (HTTP 404): Space may be sleeping or FLVE_API_TOKEN is wrong")
			case 503:
				lastErr = fmt.Errorf("FLVE Space unavailable (HTTP 503): Space is waking up, retry in a moment")
				continue
			default:
				return fmt.Errorf("FLVE returned non-JSON response (HTTP %d): Space may be unavailable", resp.StatusCode)
			}
		}
		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
		}
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		return nil
	}
	return lastErr
}

func (a *flveHTTPAdapter) doJSON(ctx context.Context, method, endpoint string, reqBody interface{}, out interface{}) error {
	var bodyReader io.Reader
	if reqBody != nil {
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, a.baseURL+endpoint, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	a.applyAuth(req)

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// Detect HTML responses — HuggingFace private Space proxy returns HTML
	// pages (404/503) instead of JSON when the Bearer token is wrong or the
	// Space is waking up. Surface a clean error instead of a JSON parse panic.
	if isHTMLResponse(resp, respBody) {
		switch resp.StatusCode {
		case 401:
			return fmt.Errorf("FLVE authentication failed (HTTP 401): check FLVE_API_TOKEN matches the HuggingFace Space token")
		case 403:
			return fmt.Errorf("FLVE access forbidden (HTTP 403): check FLVE_API_TOKEN matches the HuggingFace Space token")
		case 404:
			return fmt.Errorf("FLVE endpoint not found (HTTP 404): Space may be sleeping or FLVE_API_TOKEN is wrong — check HF Space is running and token is correct")
		case 503:
			return fmt.Errorf("FLVE Space unavailable (HTTP 503): Space is waking up, retry in a moment")
		default:
			return fmt.Errorf("FLVE returned non-JSON response (HTTP %d): Space may be unavailable", resp.StatusCode)
		}
	}

	if resp.StatusCode >= 500 {
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// isHTMLResponse returns true when the response body looks like an HTML page
// rather than a JSON API response. Used to detect HuggingFace proxy error pages.
func isHTMLResponse(resp *http.Response, body []byte) bool {
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "text/html") {
		return true
	}
	// Fallback: check body prefix (handles cases where Content-Type is missing)
	trimmed := bytes.TrimSpace(body)
	return bytes.HasPrefix(trimmed, []byte("<!DOCTYPE")) || bytes.HasPrefix(trimmed, []byte("<html"))
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if strings.Contains(msg, "status 5") {
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}
	// Connection errors are retryable
	if strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "i/o timeout") {
		return true
	}
	return false
}
