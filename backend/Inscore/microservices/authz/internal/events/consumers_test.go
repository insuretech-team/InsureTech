package events

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/domain"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/stretchr/testify/require"
)

type fakeEnforcer struct {
	domain.EnforcerIface
	addRoleErr    error
	invalidateErr error
	addRoleCalls  int
	invCalls      int
	lastSub       string
	lastRole      string
	lastDomain    string
}

func (f *fakeEnforcer) AddRoleForUserInDomain(sub, role, dom string) error {
	f.addRoleCalls++
	f.lastSub, f.lastRole, f.lastDomain = sub, role, dom
	return f.addRoleErr
}
func (f *fakeEnforcer) InvalidateCache() error {
	f.invCalls++
	return f.invalidateErr
}

func TestNewUserRegisteredHandler(t *testing.T) {
	f := &fakeEnforcer{}
	h := NewUserRegisteredHandler(f)

	valid, _ := json.Marshal(UserRegisteredPayload{UserID: "u1", Portal: "b2c"})
	require.NoError(t, h(context.Background(), &kafkaconsumer.Message{Offset: 1, Value: valid}))
	require.Equal(t, 1, f.addRoleCalls)
	require.Equal(t, "user:u1", f.lastSub)
	require.Equal(t, "role:customer", f.lastRole)
	require.Equal(t, "b2c:root", f.lastDomain)

	badJSON := []byte("{")
	require.Error(t, h(context.Background(), &kafkaconsumer.Message{Offset: 2, Value: badJSON}))

	missing, _ := json.Marshal(UserRegisteredPayload{})
	require.NoError(t, h(context.Background(), &kafkaconsumer.Message{Offset: 3, Value: missing}))
	require.Equal(t, 1, f.addRoleCalls)

	unknownPortal, _ := json.Marshal(UserRegisteredPayload{UserID: "u2", Portal: "x"})
	require.NoError(t, h(context.Background(), &kafkaconsumer.Message{Offset: 4, Value: unknownPortal}))
	require.Equal(t, 1, f.addRoleCalls)

	f.addRoleErr = errors.New("boom")
	agent, _ := json.Marshal(UserRegisteredPayload{UserID: "u3", Portal: "agent", TenantID: "t1"})
	require.NoError(t, h(context.Background(), &kafkaconsumer.Message{Offset: 5, Value: agent}))
	require.Equal(t, 2, f.addRoleCalls)
}

func TestNewPolicyCacheInvalidatedHandlerAndFanout(t *testing.T) {
	f := &fakeEnforcer{}
	h := NewPolicyCacheInvalidatedHandler(f)

	msgBody, _ := json.Marshal(PolicyCacheInvalidatedPayload{Domain: "system:root"})
	require.NoError(t, h(context.Background(), &kafkaconsumer.Message{Offset: 1, Value: msgBody}))
	require.Equal(t, 1, f.invCalls)

	require.NoError(t, h(context.Background(), &kafkaconsumer.Message{Offset: 2, Value: []byte("{")}))

	f.invalidateErr = errors.New("reload failed")
	err := h(context.Background(), &kafkaconsumer.Message{Offset: 3, Value: msgBody})
	require.Error(t, err)

	fanout := FanOutHandler(TopicHandlers{"a": func(_ context.Context, _ *kafkaconsumer.Message) error { return nil }})
	require.NoError(t, fanout(context.Background(), &kafkaconsumer.Message{Topic: "a"}))
	require.NoError(t, fanout(context.Background(), &kafkaconsumer.Message{Topic: "b"}))
}

