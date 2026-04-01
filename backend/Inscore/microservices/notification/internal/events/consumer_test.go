package events

import (
	"context"
	"testing"

	authzeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/events/v1"
	notificationv1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/entity/v1"
	partnereventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/events/v1"
	storageeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/events/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type capturingNotificationService struct {
	inputs []*DomainNotificationInput
}

func (s *capturingNotificationService) QueueDomainNotification(_ context.Context, input *DomainNotificationInput) error {
	s.inputs = append(s.inputs, input)
	return nil
}

func TestConsumer_HandleAuthzRoleAssignedAggregate(t *testing.T) {
	service := &capturingNotificationService{}
	consumer := NewConsumer(service)

	payload, err := protojson.Marshal(&authzeventsv1.RoleAssignedEvent{
		UserId:     "user-1",
		RoleId:     "role-1",
		RoleName:   "Claims Admin",
		AssignedBy: "admin-1",
		Domain:     "b2b:tenant-1",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if err := consumer.HandleMessage(context.Background(), "authz.events", payload); err != nil {
		t.Fatalf("handle message: %v", err)
	}

	if len(service.inputs) != 1 {
		t.Fatalf("expected 1 notification input, got %d", len(service.inputs))
	}
	got := service.inputs[0]
	if got.RecipientID != "user-1" {
		t.Fatalf("recipient mismatch: got %q", got.RecipientID)
	}
	if got.Subject != "Role assigned" {
		t.Fatalf("subject mismatch: got %q", got.Subject)
	}
}

func TestConsumer_HandleStorageUploadedAggregate(t *testing.T) {
	service := &capturingNotificationService{}
	consumer := NewConsumer(service)

	payload, err := protojson.Marshal(&storageeventsv1.FileUploadedEvent{
		FileId:        "file-1",
		UploadedBy:    "user-2",
		Filename:      "policy.pdf",
		ReferenceId:   "policy-1",
		ReferenceType: "policy",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if err := consumer.HandleMessage(context.Background(), "storage.events", payload); err != nil {
		t.Fatalf("handle message: %v", err)
	}

	if len(service.inputs) != 1 {
		t.Fatalf("expected 1 notification input, got %d", len(service.inputs))
	}
	got := service.inputs[0]
	if got.RecipientID != "user-2" {
		t.Fatalf("recipient mismatch: got %q", got.RecipientID)
	}
	if got.RecipientRefKind != "file" {
		t.Fatalf("recipient ref kind mismatch: got %q", got.RecipientRefKind)
	}
}

func TestConsumer_HandlePartnerVerifiedAggregate(t *testing.T) {
	service := &capturingNotificationService{}
	consumer := NewConsumer(service)

	payload, err := protojson.Marshal(&partnereventsv1.PartnerVerifiedEvent{
		PartnerId:  "partner-1",
		VerifiedBy: "ops-1",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if err := consumer.HandleMessage(context.Background(), "partner-events", payload); err != nil {
		t.Fatalf("handle message: %v", err)
	}

	if len(service.inputs) != 1 {
		t.Fatalf("expected 1 notification input, got %d", len(service.inputs))
	}
	got := service.inputs[0]
	if got.RecipientRefKind != "partner" || got.RecipientRefID != "partner-1" {
		t.Fatalf("recipient ref mismatch: kind=%q id=%q", got.RecipientRefKind, got.RecipientRefID)
	}
}

func TestConsumer_HandleMediaAggregateJSON(t *testing.T) {
	service := &capturingNotificationService{}
	consumer := NewConsumer(service)

	payload := []byte(`{"event_type":"media.processing.failed","tenant_id":"tenant-1","media_id":"media-1","job_id":"job-1","data":{"reason":"ocr timeout"}}`)
	if err := consumer.HandleMessage(context.Background(), "media-events", payload); err != nil {
		t.Fatalf("handle message: %v", err)
	}

	if len(service.inputs) != 1 {
		t.Fatalf("expected 1 notification input, got %d", len(service.inputs))
	}
	got := service.inputs[0]
	if got.RecipientRefKind != "media" || got.RecipientRefID != "media-1" {
		t.Fatalf("recipient ref mismatch: kind=%q id=%q", got.RecipientRefKind, got.RecipientRefID)
	}
	if got.Priority != notificationv1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH {
		t.Fatalf("priority mismatch: got %v", got.Priority)
	}
}

func TestConsumer_HandleDocgenAggregateJSON(t *testing.T) {
	service := &capturingNotificationService{}
	consumer := NewConsumer(service)

	payload := []byte(`{"event_type":"document.generated","document_generation_id":"gen-1","entity_type":"policy","entity_id":"policy-22","file_url":"https://example/doc.pdf","correlation_id":"corr-1"}`)
	if err := consumer.HandleMessage(context.Background(), "docgen-events", payload); err != nil {
		t.Fatalf("handle message: %v", err)
	}

	if len(service.inputs) != 1 {
		t.Fatalf("expected 1 notification input, got %d", len(service.inputs))
	}
	got := service.inputs[0]
	if got.RecipientRefKind != "policy" || got.RecipientRefID != "policy-22" {
		t.Fatalf("recipient ref mismatch: kind=%q id=%q", got.RecipientRefKind, got.RecipientRefID)
	}
	if got.Subject != "Document ready" {
		t.Fatalf("subject mismatch: got %q", got.Subject)
	}
}
