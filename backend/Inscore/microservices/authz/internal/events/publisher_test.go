package events

import (
	"context"
	"errors"
	"testing"

	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"github.com/stretchr/testify/require"
)

type fakeProducer struct {
	topic string
	key   string
	msg   interface{}
	err   error
}

func (f *fakeProducer) Produce(_ context.Context, topic, key string, msg interface{}) error {
	f.topic = topic
	f.key = key
	f.msg = msg
	return f.err
}
func (f *fakeProducer) Close() error { return nil }

func TestPublisher_NoProducerAndNilInputs(t *testing.T) {
	p := NewPublisher(nil)
	require.NoError(t, p.PublishRoleCreated(context.Background(), nil))
	require.NoError(t, p.PublishPolicyRuleCreated(context.Background(), nil))
	require.NoError(t, p.PublishRoleDeleted(context.Background(), "r1", "admin", authzentityv1.Portal_PORTAL_SYSTEM, "u1"))
}

func TestPublisher_EmitsEventsAndHandlesProducerError(t *testing.T) {
	fp := &fakeProducer{}
	p := NewPublisher(fp)
	ctx := context.Background()

	role := &authzentityv1.Role{RoleId: "role-1", Name: "admin", Portal: authzentityv1.Portal_PORTAL_SYSTEM}
	require.NoError(t, p.PublishRoleCreated(ctx, role))
	require.Equal(t, TopicAuthZEvents, fp.topic)
	require.Equal(t, "role-1", fp.key)
	require.NotNil(t, fp.msg)

	pr := &authzentityv1.PolicyRule{PolicyId: "policy-1", Subject: "role:admin", Domain: "system:root", Object: "svc:user/*", Action: "GET"}
	require.NoError(t, p.PublishPolicyRuleCreated(ctx, pr))
	require.Equal(t, "policy-1", fp.key)
	require.NoError(t, p.PublishPolicyRuleUpdated(ctx, pr))
	require.NoError(t, p.PublishPolicyRuleDeleted(ctx, "policy-1", "role:admin", "system:root", "u2"))

	require.NoError(t, p.PublishRoleAssigned(ctx, "u1", "r1", "admin", "system:root", "u2"))
	require.NoError(t, p.PublishRoleRemoved(ctx, "u1", "r1", "admin", "system:root", "u2"))

	require.NoError(t, p.PublishAccessDenied(ctx, "u1", "system:root", "svc:user/get", "GET", "s1", "127.0.0.1"))
	require.NoError(t, p.PublishAccessGranted(ctx, "u1", "system:root", "svc:user/get", "GET"))
	require.NoError(t, p.PublishPortalConfigUpdated(ctx, authzentityv1.Portal_PORTAL_SYSTEM, true, 10, 20, 30, 40, "u1"))
	require.NoError(t, p.PublishPolicyCacheInvalidated(ctx, "", "u1"))
	require.Equal(t, "global", fp.key)

	fp.err = errors.New("produce failed")
	require.NoError(t, p.PublishRoleUpdated(ctx, role))
}
