package consumers

import (
	"context"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/stretchr/testify/require"
)

func TestPortalConfigCache_BasicOps(t *testing.T) {
	c := &PortalConfigCache{configs: map[string]*PortalConfigUpdatedPayload{}}
	cfg := &PortalConfigUpdatedPayload{Portal: "PORTAL_B2C", MfaRequired: true}

	c.Set("PORTAL_B2C", cfg)
	got := c.Get("PORTAL_B2C")
	require.NotNil(t, got)
	require.True(t, got.MfaRequired)

	c.Invalidate("PORTAL_B2C")
	require.Nil(t, c.Get("PORTAL_B2C"))

	c.Set("PORTAL_B2C", cfg)
	c.Set("PORTAL_AGENT", &PortalConfigUpdatedPayload{Portal: "PORTAL_AGENT"})
	c.InvalidateAll()
	require.Nil(t, c.Get("PORTAL_B2C"))
	require.Nil(t, c.Get("PORTAL_AGENT"))
}

func TestFanOutHandler_DispatchAndUnknown(t *testing.T) {
	called := false
	h := FanOutHandler(TopicHandlers{
		"topic.a": func(ctx context.Context, msg *consumer.Message) error {
			called = true
			return nil
		},
	})

	err := h(context.Background(), &consumer.Message{Topic: "topic.a"})
	require.NoError(t, err)
	require.True(t, called)

	err = h(context.Background(), &consumer.Message{Topic: "topic.unknown", Value: []byte(`{"x":1}`)})
	require.NoError(t, err)
}

func TestPortalConfigUpdatedHandler(t *testing.T) {
	h := NewPortalConfigUpdatedHandler()

	err := h(context.Background(), &consumer.Message{
		Topic:  "authz.events",
		Offset: 1,
		Value:  []byte(`{"portal":"PORTAL_B2C","mfa_required":true}`),
	})
	require.NoError(t, err)
	got := GlobalPortalConfigCache.Get("PORTAL_B2C")
	require.NotNil(t, got)
	require.True(t, got.MfaRequired)

	err = h(context.Background(), &consumer.Message{Offset: 2, Value: []byte(`{bad-json`)})
	require.Error(t, err)

	err = h(context.Background(), &consumer.Message{Offset: 3, Value: []byte(`{"portal":""}`)})
	require.NoError(t, err)
}

func TestSimpleHandlers_NoRepoPath(t *testing.T) {
	ctx := context.Background()

	require.NoError(t, NewSMSDLRHandler(nil)(ctx, &consumer.Message{
		Offset: 1,
		Value:  []byte(`{"status":"DELIVERED"}`), // missing message_id -> skip path
	}))

	require.NoError(t, NewAccountLockedHandler(nil, nil)(ctx, &consumer.Message{
		Offset: 1,
		Value:  []byte(`{"reason":"test"}`), // missing user_id
	}))

	require.NoError(t, NewUserRegisteredHandler(nil, nil)(ctx, &consumer.Message{
		Offset: 1,
		Value:  []byte(`{"mobile_number":"+8801712345678"}`), // missing user_id
	}))

	require.NoError(t, NewPasswordChangedHandler(nil, nil)(ctx, &consumer.Message{
		Offset: 1,
		Value:  []byte(`{"ip_address":"127.0.0.1"}`), // missing user_id
	}))

	require.NoError(t, NewPasswordResetRequestedHandler(nil, nil)(ctx, &consumer.Message{
		Offset: 1,
		Value:  []byte(`{"ip_address":"127.0.0.1"}`), // missing user_id + mobile
	}))

	require.NoError(t, NewSessionRevokedAllHandler(nil, nil)(ctx, &consumer.Message{
		Offset: 1,
		Value:  []byte(`{"user_id":"u1","session_id":"s1"}`), // single-session revoke -> no SMS path
	}))
}
