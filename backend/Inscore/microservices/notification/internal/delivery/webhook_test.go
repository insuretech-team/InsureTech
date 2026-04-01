package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWebhookClientSendSignsPayload(t *testing.T) {
	var gotSignature string
	var gotTimestamp string
	var gotEvent string
	var gotBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSignature = r.Header.Get("X-InsureTech-Signature")
		gotTimestamp = r.Header.Get("X-InsureTech-Timestamp")
		gotEvent = r.Header.Get("X-InsureTech-Event")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	payload := json.RawMessage(`{"notification_id":"n-1"}`)
	client := NewWebhookClient(WebhookConfig{
		Enabled:   true,
		Timeout:   2 * time.Second,
		UserAgent: "test-agent",
	})

	resp, err := client.Send(t.Context(), &WebhookRequest{
		TargetURL:    srv.URL,
		Secret:       "secret-123",
		EventType:    "NOTIFICATION.SENT",
		Payload:      payload,
		Subscription: "sub-1",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
	require.Equal(t, "NOTIFICATION.SENT", gotEvent)
	require.NotEmpty(t, gotTimestamp)
	require.Equal(t, "sha256="+signWebhookPayload("secret-123", gotTimestamp, payload), gotSignature)
	require.JSONEq(t, string(payload), string(gotBody))
}

func TestWebhookClientSendMarks4xxPermanent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer srv.Close()

	client := NewWebhookClient(WebhookConfig{
		Enabled: true,
		Timeout: 2 * time.Second,
	})
	resp, err := client.Send(t.Context(), &WebhookRequest{
		TargetURL: srv.URL,
		Secret:    "secret-123",
		EventType: "NOTIFICATION.FAILED",
		Payload:   json.RawMessage(`{"x":1}`),
	})
	require.Error(t, err)
	require.True(t, IsPermanent(err))
	require.NotNil(t, resp)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
