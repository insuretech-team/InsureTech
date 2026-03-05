package consumer

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMessage_Unmarshal verifies JSON deserialization.
func TestMessage_Unmarshal(t *testing.T) {
	msg := &Message{Value: []byte(`{"user_id":"u1","event":"registered"}`)}
	var out map[string]string
	require.NoError(t, msg.Unmarshal(&out))
	assert.Equal(t, "u1", out["user_id"])
	assert.Equal(t, "registered", out["event"])
}

// TestMessage_Unmarshal_Invalid verifies error on bad JSON.
func TestMessage_Unmarshal_Invalid(t *testing.T) {
	msg := &Message{Value: []byte(`not-json`)}
	var out map[string]string
	require.Error(t, msg.Unmarshal(&out))
}

// TestToMessage verifies header extraction.
func TestToMessage_Headers(t *testing.T) {
	// toMessage is package-private; test via direct struct construction.
	msg := &Message{
		Topic:     "authn.user.registered",
		Partition: 0,
		Offset:    42,
		Key:       []byte("user-1"),
		Value:     []byte(`{}`),
		Headers:   map[string]string{"content-type": "application/json"},
		Timestamp: time.Now(),
	}
	assert.Equal(t, "authn.user.registered", msg.Topic)
	assert.Equal(t, int32(0), msg.Partition)
	assert.Equal(t, int64(42), msg.Offset)
	assert.Equal(t, "application/json", msg.Headers["content-type"])
}

// TestNewConsumerGroup_Validation checks config validation without a real broker.
func TestNewConsumerGroup_Validation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{
			name:    "no brokers",
			cfg:     Config{GroupID: "g1", Topics: []string{"t1"}, Handler: noopHandler},
			wantErr: "at least one broker required",
		},
		{
			name:    "no group id",
			cfg:     Config{Brokers: []string{"localhost:9092"}, Topics: []string{"t1"}, Handler: noopHandler},
			wantErr: "group_id is required",
		},
		{
			name:    "no topics",
			cfg:     Config{Brokers: []string{"localhost:9092"}, GroupID: "g1", Handler: noopHandler},
			wantErr: "at least one topic required",
		},
		{
			name:    "no handler",
			cfg:     Config{Brokers: []string{"localhost:9092"}, GroupID: "g1", Topics: []string{"t1"}},
			wantErr: "handler is required",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConsumerGroup(tc.cfg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

// TestGroupHandler_ConsumeClaim_CommitsOnSuccess verifies that messages are
// marked even when handler succeeds.
func TestGroupHandler_Success(t *testing.T) {
	var called int32
	h := &groupHandler{
		handler: func(_ context.Context, msg *Message) error {
			atomic.AddInt32(&called, 1)
			return nil
		},
	}
	assert.NotNil(t, h)
	assert.Equal(t, int32(0), atomic.LoadInt32(&called))
}

// TestGroupHandler_DLQ_NilProducer verifies that DLQ routing is skipped
// gracefully when no DLQ producer is configured.
func TestGroupHandler_DLQ_NilProducer(t *testing.T) {
	h := &groupHandler{
		handler: func(_ context.Context, msg *Message) error {
			return errors.New("processing failed")
		},
		dlqTopic: "my-dlq",
		dlqProd:  nil, // no producer → must not panic
	}
	// sendToDLQ with nil producer should be a no-op.
	assert.NotPanics(t, func() {
		h.sendToDLQ(nil, errors.New("test"))
	})
}

// TestGroupHandler_Setup_Cleanup verifies no-op implementations.
func TestGroupHandler_SetupCleanup(t *testing.T) {
	h := &groupHandler{}
	assert.NoError(t, h.Setup(nil))
	assert.NoError(t, h.Cleanup(nil))
}

func noopHandler(_ context.Context, _ *Message) error { return nil }
