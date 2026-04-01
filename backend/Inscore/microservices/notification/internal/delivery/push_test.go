package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPushClientSendMock(t *testing.T) {
	client := NewPushClient(PushConfig{Provider: "mock"})
	resp, err := client.Send(t.Context(), &PushRequest{
		RecipientID: "user-1",
		Title:       "Hello",
		Body:        "World",
		Targets: []PushTarget{
			{Provider: "FCM", DeviceToken: "token-1"},
			{Provider: "FCM", DeviceToken: "token-2"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "DELIVERED", resp.Status)
	require.Equal(t, 2, resp.SentCount)
}

func TestPushClientSendFCMInvalidTokensPermanent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":0,"failure":2,"results":[{"error":"InvalidRegistration"},{"error":"NotRegistered"}]}`))
	}))
	defer srv.Close()

	client := NewPushClient(PushConfig{
		Provider:  "fcm",
		Endpoint:  srv.URL,
		ServerKey: "server-key",
		Timeout:   2 * time.Second,
	})
	resp, err := client.Send(t.Context(), &PushRequest{
		RecipientID: "user-1",
		Title:       "Hello",
		Body:        "World",
		Targets: []PushTarget{
			{Provider: "FCM", DeviceToken: "token-1"},
			{Provider: "FCM", DeviceToken: "token-2"},
		},
	})
	require.Error(t, err)
	require.True(t, IsPermanent(err))
	require.NotNil(t, resp)
	require.ElementsMatch(t, []string{"token-1", "token-2"}, resp.InvalidTokens)
}
