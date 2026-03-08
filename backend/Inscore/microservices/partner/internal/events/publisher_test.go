package events

import (
	"context"
	"sync"
	"testing"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	partnerentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	"github.com/stretchr/testify/require"
)

// mockProducer captures published events for testing.
type mockProducer struct {
	mu       sync.Mutex
	messages []mockMessage
}

type mockMessage struct {
	Topic string
	Key   string
	Msg   interface{}
}

func (m *mockProducer) Produce(_ context.Context, topic string, key string, msg interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, mockMessage{Topic: topic, Key: key, Msg: msg})
	return nil
}

func (m *mockProducer) Close() error { return nil }

func (m *mockProducer) Events() []mockMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]mockMessage, len(m.messages))
	copy(out, m.messages)
	return out
}

func TestPublisher_PublishPartnerOnboarded(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "partner.test.events")

	partner := &partnerentityv1.Partner{
		PartnerId:        "p-001",
		OrganizationName: "Test Org",
		Type:             partnerentityv1.PartnerType_PARTNER_TYPE_CORPORATE,
		FocalPersonId:    "fp-001",
	}
	err := pub.PublishPartnerOnboarded(context.Background(), partner)
	require.NoError(t, err)

	evts := mp.Events()
	require.Len(t, evts, 1)
	require.Equal(t, "partner.test.events", evts[0].Topic)
	require.Equal(t, "p-001", evts[0].Key)
}

func TestPublisher_PublishPartnerVerified(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "partner.test.events")

	err := pub.PublishPartnerVerified(context.Background(), "p-002", "admin")
	require.NoError(t, err)

	evts := mp.Events()
	require.Len(t, evts, 1)
	require.Equal(t, "p-002", evts[0].Key)
}

func TestPublisher_PublishAgentRegistered(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "partner.test.events")

	agent := &partnerentityv1.Agent{
		AgentId:   "a-001",
		PartnerId: "p-001",
		FullName:  "Test Agent",
	}
	err := pub.PublishAgentRegistered(context.Background(), agent)
	require.NoError(t, err)

	evts := mp.Events()
	require.Len(t, evts, 1)
	require.Equal(t, "a-001", evts[0].Key)
}

func TestPublisher_PublishCommissionCalculated(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "partner.test.events")

	commission := &partnerentityv1.Commission{
		CommissionId: "c-001",
		PartnerId:    "p-001",
		AgentId:      "a-001",
		PolicyId:     "pol-001",
		Type:         partnerentityv1.CommissionType_COMMISSION_TYPE_ACQUISITION,
		CommissionAmount: &commonv1.Money{
			Amount:   5000,
			Currency: "BDT",
		},
	}
	err := pub.PublishCommissionCalculated(context.Background(), commission)
	require.NoError(t, err)

	evts := mp.Events()
	require.Len(t, evts, 1)
	require.Equal(t, "c-001", evts[0].Key)
}

func TestPublisher_NilProducer_DoesNotPanic(t *testing.T) {
	pub := NewPublisher(nil, "partner-events")
	err := pub.PublishPartnerOnboarded(context.Background(), &partnerentityv1.Partner{PartnerId: "x"})
	require.NoError(t, err)
}

func TestPublisher_NilEntity_NoOp(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "partner-events")
	err := pub.PublishPartnerOnboarded(context.Background(), nil)
	require.NoError(t, err)
	require.Empty(t, mp.Events())
}
