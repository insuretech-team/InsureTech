package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
	storageeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/events/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	TopicStorageEvents = "storage.events"
)

// EventProducer decouples Kafka implementation.
type EventProducer interface {
	Produce(ctx context.Context, topic string, key string, msg interface{}) error
	Close() error
}

type Publisher struct {
	producer EventProducer
}

func NewPublisher(producer EventProducer) *Publisher {
	return &Publisher{producer: producer}
}

func (p *Publisher) publish(ctx context.Context, topic, key string, msg interface{}) error {
	if p.producer == nil {
		appLogger.Infof("Storage event dropped: producer not configured (topic=%s, key=%s)", topic, key)
		return nil
	}
	return p.producer.Produce(ctx, topic, key, msg)
}

func (p *Publisher) PublishFileUploaded(ctx context.Context, file *storageentityv1.StoredFile, source, uploadedBy string) error {
	if file == nil {
		return nil
	}
	evt := &storageeventsv1.FileUploadedEvent{
		EventId:       uuid.New().String(),
		FileId:        file.FileId,
		TenantId:      file.TenantId,
		Filename:      file.Filename,
		ContentType:   file.ContentType,
		SizeBytes:     file.SizeBytes,
		StorageKey:    file.StorageKey,
		Bucket:        file.Bucket,
		Url:           file.Url,
		CdnUrl:        file.CdnUrl,
		ReferenceId:   file.ReferenceId,
		ReferenceType: file.ReferenceType,
		IsPublic:      file.IsPublic,
		UploadedBy:    uploadedBy,
		Source:        source,
		Timestamp:     timestamppb.New(time.Now().UTC()),
	}
	if err := p.publish(ctx, TopicStorageEvents, file.FileId, evt); err != nil {
		appLogger.Warnf("Failed to publish FileUploadedEvent (file_id=%s): %v", file.FileId, err)
	}
	return nil
}

func (p *Publisher) PublishUploadURLIssued(
	ctx context.Context,
	tenantID, fileID, filename, storageKey, referenceID, referenceType string,
	isPublic bool,
	expiresAt time.Time,
	requestedBy string,
) error {
	evt := &storageeventsv1.FileUploadURLIssuedEvent{
		EventId:       uuid.New().String(),
		TenantId:      tenantID,
		FileId:        fileID,
		Filename:      filename,
		StorageKey:    storageKey,
		ReferenceId:   referenceID,
		ReferenceType: referenceType,
		IsPublic:      isPublic,
		ExpiresAt:     timestamppb.New(expiresAt),
		RequestedBy:   requestedBy,
		Timestamp:     timestamppb.New(time.Now().UTC()),
	}
	key := fileID
	if key == "" {
		key = tenantID
	}
	if err := p.publish(ctx, TopicStorageEvents, key, evt); err != nil {
		appLogger.Warnf("Failed to publish FileUploadURLIssuedEvent (file_id=%s): %v", fileID, err)
	}
	return nil
}

func (p *Publisher) PublishFileUploadFinalized(ctx context.Context, file *storageentityv1.StoredFile, finalizedBy string) error {
	if file == nil {
		return nil
	}
	evt := &storageeventsv1.FileUploadFinalizedEvent{
		EventId:     uuid.New().String(),
		FileId:      file.FileId,
		TenantId:    file.TenantId,
		SizeBytes:   file.SizeBytes,
		ContentType: file.ContentType,
		FinalizedBy: finalizedBy,
		Timestamp:   timestamppb.New(time.Now().UTC()),
	}
	if err := p.publish(ctx, TopicStorageEvents, file.FileId, evt); err != nil {
		appLogger.Warnf("Failed to publish FileUploadFinalizedEvent (file_id=%s): %v", file.FileId, err)
	}
	return nil
}

func (p *Publisher) PublishFileMetadataUpdated(ctx context.Context, tenantID, fileID string, updatedFields []string, updatedBy string) error {
	evt := &storageeventsv1.FileMetadataUpdatedEvent{
		EventId:       uuid.New().String(),
		FileId:        fileID,
		TenantId:      tenantID,
		UpdatedFields: updatedFields,
		UpdatedBy:     updatedBy,
		Timestamp:     timestamppb.New(time.Now().UTC()),
	}
	key := fileID
	if key == "" {
		key = tenantID
	}
	if err := p.publish(ctx, TopicStorageEvents, key, evt); err != nil {
		appLogger.Warnf("Failed to publish FileMetadataUpdatedEvent (file_id=%s): %v", fileID, err)
	}
	return nil
}

func (p *Publisher) PublishFileDeleted(ctx context.Context, tenantID, fileID, storageKey, deletedBy string) error {
	evt := &storageeventsv1.FileDeletedEvent{
		EventId:    uuid.New().String(),
		FileId:     fileID,
		TenantId:   tenantID,
		StorageKey: storageKey,
		DeletedBy:  deletedBy,
		Timestamp:  timestamppb.New(time.Now().UTC()),
	}
	key := fileID
	if key == "" {
		key = tenantID
	}
	if err := p.publish(ctx, TopicStorageEvents, key, evt); err != nil {
		appLogger.Warnf("Failed to publish FileDeletedEvent (file_id=%s): %v", fileID, err)
	}
	return nil
}
