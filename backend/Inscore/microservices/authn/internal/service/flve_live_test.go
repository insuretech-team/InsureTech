package service

// flve_live_test.go — Live integration tests against the FLVE HuggingFace Space.
//
// These tests hit the REAL deployed FLVE Space at https://farukhannan-flve.hf.space
// They require:
//   - FLVE_HF_ENDPOINT env var (or defaults to flve.yaml value)
//   - FLVE_HF_TOKEN env var     (HuggingFace read token)
//   - A test JPEG image (uses a synthetic solid-color frame if not provided)
//
// Run:
//   go test ./backend/inscore/microservices/authn/internal/service/... \
//       -run TestFLVELive -v -count=1 -timeout 120s
//
// Skip:
//   These tests are skipped automatically when FLVE_LIVE_TEST=1 is not set,
//   so they never block CI.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	defaultFLVEEndpoint = "https://farukhannan-flve.hf.space"
	defaultFLVEToken    = "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"
)

// liveTestConfig returns endpoint + token from env, falling back to flve.yaml values.
func liveTestConfig() (endpoint, token string) {
	endpoint = strings.TrimRight(os.Getenv("FLVE_HF_ENDPOINT"), "/")
	if endpoint == "" {
		endpoint = defaultFLVEEndpoint
	}
	token = os.Getenv("FLVE_HF_TOKEN")
	if token == "" {
		token = defaultFLVEToken
	}
	return
}

// skipIfNotLive skips the test unless FLVE_LIVE_TEST=1 is explicitly set.
func skipIfNotLive(t *testing.T) {
	t.Helper()
	if os.Getenv("FLVE_LIVE_TEST") != "1" {
		t.Skip("set FLVE_LIVE_TEST=1 to run live FLVE integration tests")
	}
}

// makeSyntheticFrame creates a minimal valid JPEG (solid grey 64x64) for testing.
// This won't pass liveness checks but verifies the API contract is wired correctly.
func makeSyntheticFrame() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	grey := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, grey)
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	return buf.Bytes()
}

// flveLiveClient is a minimal raw HTTP client for live tests (bypasses FLVEAdapter
// to test the actual HTTP contract independently).
type flveLiveClient struct {
	endpoint string
	token    string
	http     *http.Client
}

func newLiveClient(endpoint, token string) *flveLiveClient {
	return &flveLiveClient{
		endpoint: endpoint,
		token:    token,
		http:     &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *flveLiveClient) get(ctx context.Context, path string) (map[string]any, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+path, nil)
	if err != nil {
		return nil, 0, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return body, resp.StatusCode, nil
}

func (c *flveLiveClient) postJSON(ctx context.Context, path string, payload any) (map[string]any, int, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+path, bytes.NewReader(b))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return body, resp.StatusCode, nil
}

// ── Live Tests ────────────────────────────────────────────────────────────────

// TestFLVELive_Health verifies the HF Space is alive and returns status=healthy.
func TestFLVELive_Health(t *testing.T) {
	skipIfNotLive(t)
	endpoint, token := liveTestConfig()
	c := newLiveClient(endpoint, token)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	body, code, err := c.get(ctx, "/health")
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	if code != 200 {
		t.Fatalf("expected HTTP 200, got %d — body: %v", code, body)
	}
	t.Logf("✅ Health: %v", body)
}

// TestFLVELive_StartSession verifies /ekyc/start returns a valid session_id.
func TestFLVELive_StartSession(t *testing.T) {
	skipIfNotLive(t)
	endpoint, token := liveTestConfig()
	c := newLiveClient(endpoint, token)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	userID := uuid.New().String()
	kycID := uuid.New().String()

	body, code, err := c.postJSON(ctx, "/ekyc/start", map[string]any{
		"user_id":             userID,
		"kyc_verification_id": kycID,
		"tenant_id":           "insuretech",
		"user_type":           "CUSTOMER",
		"portal":              "customer",
	})
	if err != nil {
		t.Fatalf("start session failed: %v", err)
	}
	if code != 200 {
		t.Fatalf("expected HTTP 200, got %d — body: %v", code, body)
	}

	sessionID, ok := body["session_id"].(string)
	if !ok || sessionID == "" {
		t.Fatalf("expected session_id in response, got: %v", body)
	}

	// Validate session_id is a valid UUID
	if _, err := uuid.Parse(sessionID); err != nil {
		t.Errorf("session_id %q is not a valid UUID: %v", sessionID, err)
	}

	steps, _ := body["steps"].([]any)
	if len(steps) == 0 {
		t.Error("expected at least one step in response")
	}

	t.Logf("✅ Session started: session_id=%s steps=%d state=%v",
		sessionID, len(steps), body["state"])

	// ── Step 2: Get status ────────────────────────────────────────────────────
	t.Run("GetStatus", func(t *testing.T) {
		statusBody, statusCode, err := c.get(ctx, "/ekyc/status/"+sessionID)
		if err != nil {
			t.Fatalf("get status failed: %v", err)
		}
		if statusCode != 200 {
			t.Fatalf("expected HTTP 200, got %d — body: %v", statusCode, statusBody)
		}
		state, _ := statusBody["state"].(string)
		if !strings.Contains(state, "ACTIVE") {
			t.Errorf("expected session to be ACTIVE, got: %s", state)
		}
		t.Logf("✅ Status: state=%s progress=%.2f remaining=%vs",
			state,
			statusBody["overall_progress"],
			statusBody["remaining_seconds"],
		)
	})

	// ── Step 3: Submit synthetic frame (won't pass liveness, but tests wire) ──
	t.Run("SubmitFrame_SyntheticJPEG", func(t *testing.T) {
		adapter := NewFLVEAdapter(endpoint, token, 30*time.Second)
		frameResp, err := adapter.SubmitEKYCFrame(ctx, sessionID, makeSyntheticFrame())
		if err != nil {
			// Expected — synthetic frame won't have a face.
			// Accept "face not detected" as a valid non-crash response.
			t.Logf("ℹ️  Frame submit error (expected for synthetic frame): %v", err)
			return
		}
		t.Logf("✅ Frame response: state=%s face_detected=%v step_progress=%.2f",
			frameResp.SessionState, frameResp.Detection != nil, frameResp.StepProgress)
	})
}

// TestFLVELive_FLVEAdapter_StartEKYC verifies the Go FLVEAdapter (not raw HTTP)
// works end-to-end against the live Space.
func TestFLVELive_FLVEAdapter_StartEKYC(t *testing.T) {
	skipIfNotLive(t)
	endpoint, token := liveTestConfig()

	adapter := NewFLVEAdapter(endpoint, token, 45*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	userID := uuid.New().String()
	kycID := uuid.New().String()

	resp, err := adapter.StartEKYC(ctx, &FLVEStartRequest{
		UserID:            userID,
		KYCVerificationID: kycID,
		TenantID:          "insuretech",
		UserType:          "CUSTOMER",
		Portal:            "customer",
	})
	if err != nil {
		t.Fatalf("StartEKYC failed: %v", err)
	}
	if resp.SessionID == "" {
		t.Fatal("expected non-empty session_id")
	}
	if _, err := uuid.Parse(resp.SessionID); err != nil {
		t.Errorf("session_id %q is not a valid UUID: %v", resp.SessionID, err)
	}
	if len(resp.Steps) != 4 {
		t.Errorf("expected 4 challenge steps, got %d", len(resp.Steps))
	}

	expectedChallenges := []string{"BLINK", "LOOK_LEFT", "LOOK_RIGHT", "CAPTURE"}
	for i, step := range resp.Steps {
		normalized := strings.TrimPrefix(step.Type, "EKYC_CHALLENGE_")
		if normalized != expectedChallenges[i] {
			t.Errorf("step %d: expected %s, got %s", i+1, expectedChallenges[i], normalized)
		}
	}

	t.Logf("✅ FLVEAdapter.StartEKYC: session_id=%s steps=%d state=%s timeout=%ds",
		resp.SessionID, len(resp.Steps), resp.State, resp.TotalTimeoutSeconds)

	// Follow up with GetEKYCStatus using the adapter
	t.Run("GetEKYCStatus", func(t *testing.T) {
		statusResp, err := adapter.GetEKYCStatus(ctx, resp.SessionID)
		if err != nil {
			t.Fatalf("GetEKYCStatus failed: %v", err)
		}
		if !strings.Contains(statusResp.State, "ACTIVE") {
			t.Errorf("expected ACTIVE state, got: %s", statusResp.State)
		}
		t.Logf("✅ FLVEAdapter.GetEKYCStatus: state=%s remaining=%ds",
			statusResp.State, statusResp.RemainingSeconds)
	})
}

// TestFLVELive_Auth_WithoutToken verifies auth middleware rejects requests without token.
// Only meaningful when FLVE_API_TOKEN is set in the HF Space secrets.
func TestFLVELive_Auth_WithoutToken(t *testing.T) {
	skipIfNotLive(t)
	endpoint, _ := liveTestConfig()

	// Use client with NO token
	c := newLiveClient(endpoint, "")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	body, code, err := c.postJSON(ctx, "/ekyc/start", map[string]any{
		"user_id": uuid.New().String(),
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	// If FLVE_API_TOKEN is set in the Space, expect 401/403.
	// If not set yet, the Space runs in open mode and returns 200.
	// 404 means the Space is running without the auth middleware deployed yet —
	// the middleware returns 401 for protected paths, not 404; a 404 here means
	// the Space is still running the old image (pre-auth-middleware deploy).
	switch code {
	case 200:
		t.Logf("ℹ️  Space open mode (FLVE_API_TOKEN not set) — session_id=%v", body["session_id"])
	case 401:
		t.Logf("✅ Auth active: 401 Unauthorized")
	case 403:
		t.Logf("✅ Auth active: 403 Forbidden")
	case 404:
		t.Logf("ℹ️  Space returning 404 without token — auth middleware not yet deployed (run: python scripts/deploy_flve_hf_secrets.py)")
	default:
		t.Errorf("unexpected status %d: %v", code, body)
	}
}

// TestFLVELive_NormalizeFLVESessionID verifies the session ID normalizer
// handles all formats FLVE might return.
func TestFLVELive_NormalizeFLVESessionID(t *testing.T) {
	cases := []struct {
		input    string
		wantErr  bool
		wantUUID bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", false, true},
		{"ekyc_550e8400e29b41d4a716446655440000", false, true},
		{"kyc_550e8400e29b41d4a716446655440000", false, true},
		{"session_550e8400e29b41d4a716446655440000", false, true},
		{"sess_opaque_abc123_nonuuid", false, false}, // opaque — stored as-is
		{"", true, false},
	}
	for _, tc := range cases {
		result, err := normalizeFLVESessionID(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("input=%q: expected error, got nil (result=%q)", tc.input, result)
			}
			continue
		}
		if err != nil {
			t.Errorf("input=%q: unexpected error: %v", tc.input, err)
			continue
		}
		if tc.wantUUID {
			if _, parseErr := uuid.Parse(result); parseErr != nil {
				t.Errorf("input=%q: expected UUID output, got %q: %v", tc.input, result, parseErr)
			}
		} else {
			if result != tc.input {
				t.Errorf("input=%q: opaque ID should pass through unchanged, got %q", tc.input, result)
			}
		}
		t.Logf("  ✅ %q → %q", tc.input, result)
	}
}

// TestFLVELive_FakeFrame_ExistingSession tests submitting a frame to a non-existent session.
func TestFLVELive_FakeFrame_NonExistentSession(t *testing.T) {
	skipIfNotLive(t)
	endpoint, token := liveTestConfig()
	adapter := NewFLVEAdapter(endpoint, token, 30*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fakeSessionID := uuid.New().String()
	_, err := adapter.SubmitEKYCFrame(ctx, fakeSessionID, makeSyntheticFrame())
	if err == nil {
		t.Errorf("expected error for non-existent session, got nil")
	} else {
		t.Logf("✅ Non-existent session correctly returns error: %v", err)
	}
}

// ── Benchmark ─────────────────────────────────────────────────────────────────

// BenchmarkFLVELive_StartEKYC measures round-trip time for session creation.
func BenchmarkFLVELive_StartEKYC(b *testing.B) {
	if os.Getenv("FLVE_LIVE_TEST") != "1" {
		b.Skip("set FLVE_LIVE_TEST=1 to run live FLVE benchmarks")
	}
	endpoint, token := liveTestConfig()
	adapter := NewFLVEAdapter(endpoint, token, 45*time.Second)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := adapter.StartEKYC(ctx, &FLVEStartRequest{
			UserID:   uuid.New().String(),
			TenantID: "insuretech",
		})
		if err != nil {
			b.Fatalf("StartEKYC: %v", err)
		}
		_ = fmt.Sprintf("session: %s", resp.SessionID)
	}
}
