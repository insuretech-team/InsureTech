package events

import (
	"context"
	"testing"

	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
)

func TestPublisherNoOpPaths(t *testing.T) {
	var nilPublisher *Publisher
	nilPublisher.Publish(context.Background(), TopicOrderCreated, "key", &ordersv1.Order{})

	publisher := NewPublisher(nil)
	publisher.Publish(context.Background(), TopicOrderCreated, "key", &ordersv1.Order{})
}

func TestTopicConstantsPresent(t *testing.T) {
	topics := []string{
		TopicOrderCreated,
		TopicOrderPaymentInitiated,
		TopicOrderPaymentConfirmed,
		TopicOrderCancelled,
		TopicOrderFailed,
		TopicOrderFulfillmentCompleted,
		TopicPaymentCompleted,
		TopicPaymentFailed,
		TopicPaymentVerified,
		TopicManualReviewRequested,
		TopicManualPaymentReviewed,
		TopicPolicyIssued,
		TopicB2BPurchaseOrderApproved,
	}
	for _, topic := range topics {
		if topic == "" {
			t.Fatalf("expected non-empty topic")
		}
	}
}
