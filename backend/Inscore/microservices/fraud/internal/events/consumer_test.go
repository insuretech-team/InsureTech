package events

import (
	"context"
	"encoding/json"
	"testing"

	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	"github.com/stretchr/testify/require"
)

// mockFraudChecker captures CheckFraud calls.
type mockFraudChecker struct {
	calls []checkCall
}

type checkCall struct {
	EntityType string
	EntityID   string
}

func (m *mockFraudChecker) CheckFraud(_ context.Context, req *fraudservicev1.CheckFraudRequest) (*fraudservicev1.CheckFraudResponse, error) {
	m.calls = append(m.calls, checkCall{EntityType: req.EntityType, EntityID: req.EntityId})
	return &fraudservicev1.CheckFraudResponse{IsFraudDetected: false}, nil
}

func TestConsumer_HandleMessage_ValidPayload(t *testing.T) {
	checker := &mockFraudChecker{}
	consumer := NewConsumer(checker)

	payload, _ := json.Marshal(map[string]any{
		"entity_type": "POLICY",
		"entity_id":   "pol-001",
		"premium":     50000,
	})

	err := consumer.HandleMessage(context.Background(), "policy.issued", "key-1", payload)
	require.NoError(t, err)
	require.Len(t, checker.calls, 1)
	require.Equal(t, "POLICY", checker.calls[0].EntityType)
	require.Equal(t, "pol-001", checker.calls[0].EntityID)
}

func TestConsumer_HandleMessage_MissingEntityFields(t *testing.T) {
	checker := &mockFraudChecker{}
	consumer := NewConsumer(checker)

	payload, _ := json.Marshal(map[string]any{
		"premium": 50000,
	})

	err := consumer.HandleMessage(context.Background(), "policy.issued", "key-1", payload)
	require.NoError(t, err) // should skip, not error
	require.Empty(t, checker.calls)
}

func TestConsumer_HandleMessage_InvalidJSON(t *testing.T) {
	checker := &mockFraudChecker{}
	consumer := NewConsumer(checker)

	err := consumer.HandleMessage(context.Background(), "policy.issued", "key-1", []byte("not json"))
	require.Error(t, err)
}

func TestConsumer_HandleMessage_EmptyPayload(t *testing.T) {
	checker := &mockFraudChecker{}
	consumer := NewConsumer(checker)

	err := consumer.HandleMessage(context.Background(), "policy.issued", "key-1", nil)
	require.NoError(t, err)
	require.Empty(t, checker.calls)
}

func TestConsumer_HandleMessage_NilSvc(t *testing.T) {
	consumer := NewConsumer(nil)
	err := consumer.HandleMessage(context.Background(), "t", "k", []byte(`{"a":"b"}`))
	require.NoError(t, err)
}
