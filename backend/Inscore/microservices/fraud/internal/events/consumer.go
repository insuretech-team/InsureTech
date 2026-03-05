package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// FraudChecker defines the async check API required by the consumer.
type FraudChecker interface {
	CheckFraud(ctx context.Context, req *fraudservicev1.CheckFraudRequest) (*fraudservicev1.CheckFraudResponse, error)
}

// Consumer performs async fraud checks from upstream lifecycle events.
type Consumer struct {
	svc FraudChecker
}

func NewConsumer(svc FraudChecker) *Consumer {
	return &Consumer{svc: svc}
}

// HandleMessage consumes a JSON event payload and executes fraud checks when possible.
func (c *Consumer) HandleMessage(ctx context.Context, topic string, key string, payload []byte) error {
	if c == nil || c.svc == nil {
		return nil
	}
	if len(payload) == 0 {
		return nil
	}

	var event map[string]any
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("decode event payload: %w", err)
	}

	entityType := firstString(event,
		"entity_type", "entityType", "subject_type", "subjectType", "type",
	)
	entityID := firstString(event,
		"entity_id", "entityId", "claim_id", "claimId", "policy_id", "policyId", "customer_id", "customerId", "id",
	)
	if entityType == "" || entityID == "" {
		appLogger.Debug("fraud consumer: event skipped (missing entity fields)")
		return nil
	}

	data, err := structpb.NewStruct(event)
	if err != nil {
		return fmt.Errorf("build struct payload: %w", err)
	}

	_, err = c.svc.CheckFraud(ctx, &fraudservicev1.CheckFraudRequest{
		EntityType: strings.ToUpper(entityType),
		EntityId:   entityID,
		Data:       data,
	})
	if err != nil {
		return err
	}

	appLogger.Infof("fraud consumer: processed topic=%s key=%s entity=%s:%s", topic, key, entityType, entityID)
	return nil
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		v, ok := m[k]
		if !ok {
			continue
		}
		s := strings.TrimSpace(fmt.Sprint(v))
		if s != "" && s != "<nil>" {
			return s
		}
	}
	return ""
}
