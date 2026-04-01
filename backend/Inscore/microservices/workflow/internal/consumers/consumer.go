// Package consumers provides Kafka event consumers for the workflow engine.
// Inbound events (claim filed, endorsement requested, etc.) trigger workflow
// instance creation using the appropriate registered definition.
package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/events"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	workflowservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/services/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

// WorkflowTriggerConsumer listens for entity-level events and auto-starts
// the appropriate workflow definition for each entity type.
type WorkflowTriggerConsumer struct {
	svc    domain.WorkflowService
	repo   domain.WorkflowRepository
	topics []string
}

func NewWorkflowTriggerConsumer(svc domain.WorkflowService, repo domain.WorkflowRepository) *WorkflowTriggerConsumer {
	return &WorkflowTriggerConsumer{
		svc:  svc,
		repo: repo,
		topics: []string{
			events.TopicClaimFiled,
			events.TopicEndorsementRequested,
			events.TopicRefundRequested,
		},
	}
}

func (c *WorkflowTriggerConsumer) Topics() []string { return c.topics }

// Setup is run at the beginning of a new session.
func (c *WorkflowTriggerConsumer) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session.
func (c *WorkflowTriggerConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from a Kafka partition.
func (c *WorkflowTriggerConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		c.handleMessage(session.Context(), msg)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *WorkflowTriggerConsumer) handleMessage(ctx context.Context, msg *sarama.ConsumerMessage) {
	var payload map[string]any
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		appLogger.Errorf("workflow consumer: failed to parse message from %s: %v", msg.Topic, err)
		return
	}

	entityID := firstString(payload, "claim_id", "endorsement_id", "refund_id", "entity_id")
	entityType := topicToEntityType(msg.Topic)
	initiatedBy := firstString(payload, "initiated_by", "user_id", "created_by")
	if initiatedBy == "" {
		appLogger.Warnf("workflow consumer: missing initiated_by/user_id in message from %s, skipping auto-start", msg.Topic)
		return
	}

	if entityID == "" || entityType == "" {
		appLogger.Warnf("workflow consumer: missing entity_id or entity_type in message from %s", msg.Topic)
		return
	}

	// Find active definition for this entity type
	defs, _, err := c.repo.ListDefinitions(ctx, entityType, 0, 1, 0)
	if err != nil || len(defs) == 0 {
		appLogger.Warnf("workflow consumer: no active definition for entity_type=%s, skipping auto-start", entityType)
		return
	}

	def := defs[0]

	// Build context struct from payload
	contextStruct, _ := structpb.NewStruct(payload)

	startCtx := metadata.NewIncomingContext(ctx, metadata.Pairs("x-user-id", initiatedBy))
	_, err = c.svc.StartWorkflow(startCtx, &workflowservicev1.StartWorkflowRequest{
		WorkflowDefinitionId: def.WorkflowDefinitionId,
		EntityType:           entityType,
		EntityId:             entityID,
		Context:              contextStruct,
	})
	if err != nil {
		appLogger.Errorf("workflow consumer: failed to start workflow for %s/%s: %v", entityType, entityID, err)
		return
	}

	appLogger.Infof("workflow consumer: auto-started workflow for %s/%s from topic %s", entityType, entityID, msg.Topic)
}

// ─── HELPERS ─────────────────────────────────────────────────────────────────

func topicToEntityType(topic string) string {
	switch topic {
	case events.TopicClaimFiled:
		return "CLAIM"
	case events.TopicEndorsementRequested:
		return "ENDORSEMENT"
	case events.TopicRefundRequested:
		return "REFUND"
	default:
		return ""
	}
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// correlationID generates a fallback correlation id.
func correlationID() string {
	return fmt.Sprintf("wf-%s", uuid.NewString()[:8])
}
