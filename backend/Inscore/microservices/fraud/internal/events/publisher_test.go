package events

import (
	"context"
	"sync"
	"testing"

	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
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

func TestPublisher_PublishFraudAlertTriggered(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "fraud.test.events")

	alert := &fraudv1.FraudAlert{
		Id:          "alert-001",
		AlertNumber: "FAL-20240101",
		EntityType:  "POLICY",
		EntityId:    "pol-001",
		RiskLevel:   "RISK_LEVEL_HIGH",
		FraudScore:  70,
	}
	err := pub.PublishFraudAlertTriggered(context.Background(), alert, "corr-123")
	require.NoError(t, err)

	evts := mp.Events()
	require.Len(t, evts, 1)
	require.Equal(t, "fraud.test.events", evts[0].Topic)
	require.Equal(t, "alert-001", evts[0].Key)
}

func TestPublisher_PublishFraudCaseCreated(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "fraud.test.events")

	fc := &fraudv1.FraudCase{
		Id:           "case-001",
		CaseNumber:   "FRC-20240101",
		FraudAlertId: "alert-001",
		Priority:     fraudv1.CasePriority_CASE_PRIORITY_HIGH,
	}
	err := pub.PublishFraudCaseCreated(context.Background(), fc, "corr-456")
	require.NoError(t, err)

	evts := mp.Events()
	require.Len(t, evts, 1)
	require.Equal(t, "case-001", evts[0].Key)
}

func TestPublisher_PublishFraudConfirmed(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "fraud.test.events")

	fc := &fraudv1.FraudCase{Id: "case-002"}
	err := pub.PublishFraudConfirmed(context.Background(), fc, "CLAIM", "clm-001", "corr-789")
	require.NoError(t, err)

	evts := mp.Events()
	require.Len(t, evts, 1)
	require.Equal(t, "case-002", evts[0].Key)
}

func TestPublisher_NilProducer_DoesNotPanic(t *testing.T) {
	pub := NewPublisher(nil, "fraud.events")
	err := pub.PublishFraudAlertTriggered(context.Background(), &fraudv1.FraudAlert{Id: "x"}, "")
	require.NoError(t, err)
}

func TestPublisher_NilAlert_NoOp(t *testing.T) {
	mp := &mockProducer{}
	pub := NewPublisher(mp, "fraud.events")
	err := pub.PublishFraudAlertTriggered(context.Background(), nil, "")
	require.NoError(t, err)
	require.Empty(t, mp.Events())
}
