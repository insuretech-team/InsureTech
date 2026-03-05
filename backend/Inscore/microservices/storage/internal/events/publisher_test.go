package events

import (
	"context"
	"testing"

	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
	storageeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/events/v1"
)

type mockProducer struct {
	calls []produceCall
}

type produceCall struct {
	topic string
	key   string
	msg   interface{}
}

func (m *mockProducer) Produce(_ context.Context, topic string, key string, msg interface{}) error {
	m.calls = append(m.calls, produceCall{
		topic: topic,
		key:   key,
		msg:   msg,
	})
	return nil
}

func (m *mockProducer) Close() error { return nil }

func TestPublishFileUploaded(t *testing.T) {
	t.Parallel()

	mp := &mockProducer{}
	pub := NewPublisher(mp)
	file := &storageentityv1.StoredFile{
		FileId:      "file-1",
		TenantId:    "tenant-1",
		Filename:    "id-front.jpg",
		ContentType: "image/jpeg",
		SizeBytes:   12,
		StorageKey:  "insuretech/tenant-1/file-1.jpg",
	}

	if err := pub.PublishFileUploaded(context.Background(), file, "DIRECT", "user-1"); err != nil {
		t.Fatalf("PublishFileUploaded returned error: %v", err)
	}
	if len(mp.calls) != 1 {
		t.Fatalf("expected 1 produce call, got %d", len(mp.calls))
	}
	call := mp.calls[0]
	if call.topic != TopicStorageEvents {
		t.Fatalf("topic mismatch: got=%q want=%q", call.topic, TopicStorageEvents)
	}
	if call.key != "file-1" {
		t.Fatalf("key mismatch: got=%q want=%q", call.key, "file-1")
	}
	if _, ok := call.msg.(*storageeventsv1.FileUploadedEvent); !ok {
		t.Fatalf("unexpected message type: %T", call.msg)
	}
}

func TestPublishNoProducerDoesNotError(t *testing.T) {
	t.Parallel()

	pub := NewPublisher(nil)
	err := pub.PublishFileDeleted(context.Background(), "tenant-1", "file-1", "k", "user-1")
	if err != nil {
		t.Fatalf("expected nil error with nil producer, got %v", err)
	}
}
