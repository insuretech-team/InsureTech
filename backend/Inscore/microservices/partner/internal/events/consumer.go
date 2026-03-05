package events

import (
	"context"
	"fmt"
	"strings"

	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
)

// CommissionEventProcessor defines business callback for policy commission events.
type CommissionEventProcessor interface {
	ProcessPolicyCommissionEvent(ctx context.Context, policyID string, cType partnerv1.CommissionType) error
}

// NewPolicyLifecycleHandler returns a Kafka handler for policy issued/renewed events.
func NewPolicyLifecycleHandler(processor CommissionEventProcessor) kafkaconsumer.HandlerFunc {
	return func(ctx context.Context, msg *kafkaconsumer.Message) error {
		if processor == nil || msg == nil {
			return nil
		}

		var payload map[string]any
		if err := msg.Unmarshal(&payload); err != nil {
			return fmt.Errorf("partner commission consumer: invalid payload: %w", err)
		}

		policyID, cType := resolvePolicyEvent(msg.Topic, payload)
		if policyID == "" {
			appLogger.Debug("partner commission consumer: skipped (missing policy id)")
			return nil
		}

		if err := processor.ProcessPolicyCommissionEvent(ctx, policyID, cType); err != nil {
			return err
		}

		appLogger.Infof("partner commission consumer: processed topic=%s policy_id=%s type=%s", msg.Topic, policyID, cType.String())
		return nil
	}
}

func resolvePolicyEvent(topic string, payload map[string]any) (string, partnerv1.CommissionType) {
	topicLower := strings.ToLower(strings.TrimSpace(topic))

	eventType := strings.ToLower(firstString(payload,
		"event_type", "eventType", "type", "name",
	))
	isRenewal := strings.Contains(topicLower, "renew") || strings.Contains(eventType, "renew")
	if !isRenewal {
		// new_policy_id strongly indicates renewal workflow payload.
		if strings.TrimSpace(firstString(payload, "new_policy_id", "newPolicyId")) != "" {
			isRenewal = true
		}
	}

	if isRenewal {
		policyID := firstString(payload, "new_policy_id", "newPolicyId", "policy_id", "policyId")
		return policyID, partnerv1.CommissionType_COMMISSION_TYPE_RENEWAL
	}

	policyID := firstString(payload, "policy_id", "policyId")
	return policyID, partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION
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
