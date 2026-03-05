package events

import (
	"context"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	partnerentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	partnereventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/events/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EventProducer interface for decoupling Kafka/MsgQueue
type EventProducer interface {
	Produce(ctx context.Context, topic string, key string, msg interface{}) error
	Close() error
}

// Publisher publishes Partner domain events to Kafka.
type Publisher struct {
	producer EventProducer
	topic    string
}

// NewPublisher creates a new Publisher. producer may be nil (events will be dropped with a log).
func NewPublisher(producer EventProducer, topic string) *Publisher {
	return &Publisher{producer: producer, topic: topic}
}

func (p *Publisher) publish(ctx context.Context, key string, msg interface{}) error {
	if p == nil || p.producer == nil {
		topic := "partner-events"
		if p != nil && p.topic != "" {
			topic = p.topic
		}
		appLogger.Infof("partner: Kafka producer not configured — event dropped (topic=%s)", topic)
		return nil
	}
	// Fallback to default topic if not set
	topic := p.topic
	if topic == "" {
		topic = "partner-events"
	}
	return p.producer.Produce(ctx, topic, key, msg)
}

// PublishPartnerOnboarded emits a PartnerOnboardedEvent.
func (p *Publisher) PublishPartnerOnboarded(ctx context.Context, partner *partnerentityv1.Partner) error {
	if partner == nil {
		return nil
	}
	evt := &partnereventsv1.PartnerOnboardedEvent{
		EventId:          uuid.New().String(),
		PartnerId:        partner.PartnerId,
		OrganizationName: partner.OrganizationName,
		PartnerType:      partner.Type.String(),
		FocalPersonId:    partner.FocalPersonId,
		Timestamp:        timestamppb.Now(),
	}
	if err := p.publish(ctx, partner.PartnerId, evt); err != nil {
		appLogger.Warnf("partner: failed to publish PartnerOnboardedEvent (partner_id=%s): %v", partner.PartnerId, err)
	}
	return nil
}

// PublishPartnerVerified emits a PartnerVerifiedEvent.
func (p *Publisher) PublishPartnerVerified(ctx context.Context, partnerId string, verifiedBy string) error {
	evt := &partnereventsv1.PartnerVerifiedEvent{
		EventId:    uuid.New().String(),
		PartnerId:  partnerId,
		VerifiedBy: verifiedBy,
		Timestamp:  timestamppb.Now(),
	}
	if err := p.publish(ctx, partnerId, evt); err != nil {
		appLogger.Warnf("partner: failed to publish PartnerVerifiedEvent (partner_id=%s): %v", partnerId, err)
	}
	return nil
}

// PublishAgentRegistered emits an AgentRegisteredEvent.
func (p *Publisher) PublishAgentRegistered(ctx context.Context, agent *partnerentityv1.Agent) error {
	if agent == nil {
		return nil
	}
	evt := &partnereventsv1.AgentRegisteredEvent{
		EventId:   uuid.New().String(),
		AgentId:   agent.AgentId,
		PartnerId: agent.PartnerId,
		AgentName: agent.FullName,
		Timestamp: timestamppb.Now(),
	}
	if err := p.publish(ctx, agent.AgentId, evt); err != nil {
		appLogger.Warnf("partner: failed to publish AgentRegisteredEvent (agent_id=%s): %v", agent.AgentId, err)
	}
	return nil
}

// PublishCommissionCalculated emits a CommissionCalculatedEvent.
func (p *Publisher) PublishCommissionCalculated(ctx context.Context, commission *partnerentityv1.Commission) error {
	if commission == nil {
		return nil
	}
	evt := &partnereventsv1.CommissionCalculatedEvent{
		EventId:      uuid.New().String(),
		CommissionId: commission.CommissionId,
		PartnerId:    commission.PartnerId,
		AgentId:      commission.AgentId,
		PolicyId:     commission.PolicyId,
		CommissionAmount: &commonv1.Money{
			Amount:        commission.GetCommissionAmount().GetAmount(),
			Currency:      commission.GetCommissionAmount().GetCurrency(),
			DecimalAmount: commission.GetCommissionAmount().GetDecimalAmount(),
		},
		CommissionType: commission.Type.String(),
		Timestamp:      timestamppb.Now(),
	}
	if err := p.publish(ctx, commission.CommissionId, evt); err != nil {
		appLogger.Warnf("partner: failed to publish CommissionCalculatedEvent (commission_id=%s): %v", commission.CommissionId, err)
	}
	return nil
}
