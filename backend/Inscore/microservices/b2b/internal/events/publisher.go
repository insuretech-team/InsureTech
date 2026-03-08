package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	b2beventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/events/v1"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Kafka topic names for B2B events
const (
	TopicOrganisationCreated   = "b2b.organisation.created"
	TopicOrganisationUpdated   = "b2b.organisation.updated"
	TopicOrganisationApproved  = "b2b.organisation.approved"
	TopicOrganisationSuspended = "b2b.organisation.suspended"
	TopicOrgMemberAdded        = "b2b.org_member.added"
	TopicOrgMemberRemoved      = "b2b.org_member.removed"
	TopicB2BAdminAssigned      = "b2b.admin.assigned"
)

// EventProducer interface for decoupling Kafka/MsgQueue
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

func (p *Publisher) PublishOrganisationCreated(ctx context.Context, organisationID, tenantID, name, code, createdBy string) error {
	evt := &b2beventsv1.OrganisationCreatedEvent{
		EventId:        uuid.New().String(),
		OrganisationId: organisationID,
		TenantId:       tenantID,
		Name:           name,
		Code:           code,
		CreatedBy:      createdBy,
		Timestamp:      timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicOrganisationCreated, organisationID, evt); err != nil {
		appLogger.Warnf("Failed to publish OrganisationCreatedEvent for org %s: %v", organisationID, err)
	}
	return nil
}

func (p *Publisher) PublishOrganisationUpdated(ctx context.Context, organisationID, name string, status b2bv1.OrganisationStatus, updatedBy string) error {
	evt := &b2beventsv1.OrganisationUpdatedEvent{
		EventId:        uuid.New().String(),
		OrganisationId: organisationID,
		Name:           name,
		Status:         status,
		UpdatedBy:      updatedBy,
		Timestamp:      timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicOrganisationUpdated, organisationID, evt); err != nil {
		appLogger.Warnf("Failed to publish OrganisationUpdatedEvent for org %s: %v", organisationID, err)
	}
	return nil
}

func (p *Publisher) PublishOrganisationApproved(ctx context.Context, organisationID, approvedBy string) error {
	evt := &b2beventsv1.OrganisationApprovedEvent{
		EventId:        uuid.New().String(),
		OrganisationId: organisationID,
		ApprovedBy:     approvedBy,
		Timestamp:      timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicOrganisationApproved, organisationID, evt); err != nil {
		appLogger.Warnf("Failed to publish OrganisationApprovedEvent for org %s: %v", organisationID, err)
	}
	return nil
}

func (p *Publisher) PublishOrganisationSuspended(ctx context.Context, organisationID, reason, suspendedBy string) error {
	evt := &b2beventsv1.OrganisationSuspendedEvent{
		EventId:        uuid.New().String(),
		OrganisationId: organisationID,
		Reason:         reason,
		SuspendedBy:    suspendedBy,
		Timestamp:      timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicOrganisationSuspended, organisationID, evt); err != nil {
		appLogger.Warnf("Failed to publish OrganisationSuspendedEvent for org %s: %v", organisationID, err)
	}
	return nil
}

func (p *Publisher) PublishOrgMemberAdded(ctx context.Context, memberID, organisationID, userID string, role b2bv1.OrgMemberRole, addedBy string) error {
	evt := &b2beventsv1.OrgMemberAddedEvent{
		EventId:        uuid.New().String(),
		MemberId:       memberID,
		OrganisationId: organisationID,
		UserId:         userID,
		Role:           role,
		AddedBy:        addedBy,
		Timestamp:      timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicOrgMemberAdded, organisationID, evt); err != nil {
		appLogger.Warnf("Failed to publish OrgMemberAddedEvent for org %s: %v", organisationID, err)
	}
	return nil
}

func (p *Publisher) PublishOrgMemberRemoved(ctx context.Context, memberID, organisationID, userID, removedBy string) error {
	evt := &b2beventsv1.OrgMemberRemovedEvent{
		EventId:        uuid.New().String(),
		MemberId:       memberID,
		OrganisationId: organisationID,
		UserId:         userID,
		RemovedBy:      removedBy,
		Timestamp:      timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicOrgMemberRemoved, organisationID, evt); err != nil {
		appLogger.Warnf("Failed to publish OrgMemberRemovedEvent for org %s: %v", organisationID, err)
	}
	return nil
}

func (p *Publisher) PublishB2BAdminAssigned(ctx context.Context, organisationID, userID, assignedBy string) error {
	evt := &b2beventsv1.B2BAdminAssignedEvent{
		EventId:        uuid.New().String(),
		OrganisationId: organisationID,
		UserId:         userID,
		AssignedBy:     assignedBy,
		Timestamp:      timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicB2BAdminAssigned, organisationID, evt); err != nil {
		appLogger.Warnf("Failed to publish B2BAdminAssignedEvent for org %s: %v", organisationID, err)
	}
	return nil
}

// publish is the internal helper that sends to Kafka if a producer is wired,
// or logs a warning and no-ops if running without Kafka (e.g. dev/test).
func (p *Publisher) publish(ctx context.Context, topic, key string, msg interface{}) error {
	if p.producer == nil {
		appLogger.Infof("Kafka producer not configured - event dropped (topic=%s, key=%s)", topic, key)
		return nil
	}
	return p.producer.Produce(ctx, topic, key, msg)
}
