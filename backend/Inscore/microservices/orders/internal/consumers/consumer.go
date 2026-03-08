// Package consumers handles incoming Kafka events for the orders microservice.
// Consumes payment.completed / payment.failed / payment.verified / policy.issued /
// b2b.purchase_order.approved events using typed protojson unmarshalling.
package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/service"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	b2beventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/events/v1"
	documentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/document/services/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	paymenteventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/events/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// EventConsumer handles incoming Kafka events for the orders domain.
type EventConsumer struct {
	repo          domain.OrderRepository
	orderSvc      *service.OrderServiceImpl // for CreateOrder in B2B flow
	docgenClient  documentservicev1.DocumentServiceClient // receipt PDF generation
	storageClient storageservicev1.StorageServiceClient   // file reference validation
}

// NewEventConsumer creates an EventConsumer with all service dependencies.
// docgenClient and storageClient are optional (nil = feature disabled, logs warning).
func NewEventConsumer(
	repo domain.OrderRepository,
	orderSvc *service.OrderServiceImpl,
	docgenClient documentservicev1.DocumentServiceClient,
	storageClient storageservicev1.StorageServiceClient,
) *EventConsumer {
	return &EventConsumer{
		repo:          repo,
		orderSvc:      orderSvc,
		docgenClient:  docgenClient,
		storageClient: storageClient,
	}
}


// HandlePaymentCompleted processes insuretech.payment.v1.payment.completed events.
// Uses typed protojson unmarshalling to read order_id from PaymentCompletedEvent.
// Transitions the order to PAID status and updates payment_status dimension.
func (c *EventConsumer) HandlePaymentCompleted(ctx context.Context, data []byte) error {
	if len(data) == 0 {
		appLogger.Warnf("HandlePaymentCompleted: empty payload — skipping")
		return nil
	}

	var evt paymenteventsv1.PaymentCompletedEvent
	if err := protojson.Unmarshal(data, &evt); err != nil {
		appLogger.Errorf("HandlePaymentCompleted: protojson unmarshal failed: %v", err)
		return fmt.Errorf("unmarshal PaymentCompletedEvent: %w", err)
	}

	// Use typed order_id field (populated after payment_events.proto was extended).
	// Fall back to correlation_id for events published by older service versions.
	orderID := evt.GetOrderId()
	if orderID == "" {
		orderID = evt.GetCorrelationId()
	}
	if orderID == "" {
		appLogger.Warnf("HandlePaymentCompleted: missing order_id in PaymentCompletedEvent — skipping")
		return nil // not retryable — event is missing required field
	}

	appLogger.Infof("HandlePaymentCompleted: order=%s payment=%s provider=%s — transitioning to PAID",
		orderID, evt.GetPaymentId(), evt.GetProvider())

	if err := c.repo.UpdateOrderStatus(ctx, orderID, ordersv1.OrderStatus_ORDER_STATUS_PAID); err != nil {
		return fmt.Errorf("HandlePaymentCompleted UpdateOrderStatus: %w", err)
	}
	// Update payment dimension — non-fatal if it fails (primary status already set)
	if err := c.repo.SetPaymentStatus(ctx, orderID, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID); err != nil {
		appLogger.Warnf("HandlePaymentCompleted: SetPaymentStatus failed for order %s: %v", orderID, err)
	}
	return nil
}

// HandlePaymentFailed processes insuretech.payment.v1.payment.failed events.
// Transitions the order to FAILED status with the failure reason.
func (c *EventConsumer) HandlePaymentFailed(ctx context.Context, data []byte) error {
	if len(data) == 0 {
		appLogger.Warnf("HandlePaymentFailed: empty payload — skipping")
		return nil
	}

	var evt paymenteventsv1.PaymentFailedEvent
	if err := protojson.Unmarshal(data, &evt); err != nil {
		appLogger.Errorf("HandlePaymentFailed: protojson unmarshal failed: %v", err)
		return fmt.Errorf("unmarshal PaymentFailedEvent: %w", err)
	}

	// Use typed order_id field (populated after payment_events.proto was extended).
	// Fall back to correlation_id for events published by older service versions.
	orderID := evt.GetOrderId()
	if orderID == "" {
		orderID = evt.GetCorrelationId()
	}
	if orderID == "" {
		appLogger.Warnf("HandlePaymentFailed: missing order_id in PaymentFailedEvent — skipping")
		return nil
	}

	reason := evt.GetErrorMessage()
	if reason == "" {
		reason = evt.GetErrorCode()
	}
	if reason == "" {
		reason = "Payment failed"
	}

	appLogger.Infof("HandlePaymentFailed: order=%s — marking FAILED: %s", orderID, reason)

	if err := c.repo.SetFailureReason(ctx, orderID, reason); err != nil {
		return fmt.Errorf("HandlePaymentFailed SetFailureReason: %w", err)
	}
	if err := c.repo.SetPaymentStatus(ctx, orderID, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAYMENT_FAILED); err != nil {
		appLogger.Warnf("HandlePaymentFailed: SetPaymentStatus failed for order %s: %v", orderID, err)
	}
	return nil
}

// HandlePolicyIssued processes insuretech.insurance.v1.policy.issued events.
// Updates the order with the issued policy_id and transitions to POLICY_ISSUED.
// Uses flat JSON parsing since insurance-service events don't yet carry typed order_id.
func (c *EventConsumer) HandlePolicyIssued(ctx context.Context, data []byte) error {
	if len(data) == 0 {
		appLogger.Warnf("HandlePolicyIssued: empty payload — skipping")
		return nil
	}

	// TODO: switch to typed protojson once insurance_events.proto has typed order_id field.
	payload, err := flattenJSONMap(data)
	if err != nil {
		appLogger.Errorf("HandlePolicyIssued: failed to parse payload: %v", err)
		return fmt.Errorf("parse policy.issued: %w", err)
	}

	orderID := payload["order_id"]
	policyID := payload["policy_id"]

	if orderID == "" || policyID == "" {
		appLogger.Warnf("HandlePolicyIssued: missing order_id or policy_id in payload — skipping")
		return nil
	}

	appLogger.Infof("HandlePolicyIssued: policy=%s → order=%s", policyID, orderID)

	if err := c.repo.SetPolicyID(ctx, orderID, policyID); err != nil {
		return fmt.Errorf("HandlePolicyIssued SetPolicyID: %w", err)
	}
	if err := c.repo.SetFulfillmentStatus(ctx, orderID, ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLED); err != nil {
		appLogger.Warnf("HandlePolicyIssued: SetFulfillmentStatus failed for order %s: %v", orderID, err)
	}
	return nil
}

// HandlePaymentVerified processes insuretech.payment.v1.payment.verified events.
// This event is published by the payment-service when a system operator has manually
// verified a payment that was held for review. On receipt the order is re-transitioned
// to PAID and the payment-status dimension updated so it can continue to fulfillment.
func (c *EventConsumer) HandlePaymentVerified(ctx context.Context, data []byte) error {
	if len(data) == 0 {
		appLogger.Warnf("HandlePaymentVerified: empty payload — skipping")
		return nil
	}

	var evt paymenteventsv1.PaymentVerifiedEvent
	if err := protojson.Unmarshal(data, &evt); err != nil {
		appLogger.Errorf("HandlePaymentVerified: protojson unmarshal failed: %v", err)
		return fmt.Errorf("unmarshal PaymentVerifiedEvent: %w", err)
	}

	orderID := evt.GetOrderId()
	if orderID == "" {
		orderID = evt.GetCorrelationId()
	}
	if orderID == "" {
		appLogger.Warnf("HandlePaymentVerified: missing order_id — skipping")
		return nil
	}

	appLogger.Infof("HandlePaymentVerified: payment=%s order=%s verifiedBy=%s → PAID",
		evt.GetPaymentId(), orderID, evt.GetVerifiedBy())

	if err := c.repo.UpdateOrderStatus(ctx, orderID, ordersv1.OrderStatus_ORDER_STATUS_PAID); err != nil {
		return fmt.Errorf("HandlePaymentVerified UpdateOrderStatus: %w", err)
	}
	if err := c.repo.SetPaymentStatus(ctx, orderID, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID); err != nil {
		appLogger.Warnf("HandlePaymentVerified: SetPaymentStatus failed for order %s: %v", orderID, err)
	}

	// Trigger receipt generation now that payment is confirmed.
	c.triggerReceiptDocgen(ctx, orderID, evt.GetPaymentId(), evt.GetTenantId(), evt.GetOrganisationId())
	return nil
}

// HandleManualReviewRequested processes insuretech.payment.v1.payment.manual_review_requested.
// Holds the order in PENDING_REVIEW status while staff verify the manual proof of payment.
func (c *EventConsumer) HandleManualReviewRequested(ctx context.Context, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var evt paymenteventsv1.ManualPaymentProofSubmittedEvent
	if err := protojson.Unmarshal(data, &evt); err != nil {
		appLogger.Errorf("HandleManualReviewRequested: unmarshal failed: %v", err)
		return fmt.Errorf("unmarshal ManualPaymentProofSubmittedEvent: %w", err)
	}

	orderID := evt.GetOrderId()
	if orderID == "" {
		appLogger.Warnf("HandleManualReviewRequested: missing order_id — skipping")
		return nil
	}

	appLogger.Infof("HandleManualReviewRequested: payment=%s order=%s → PENDING_REVIEW", evt.GetPaymentId(), orderID)

	// Transition order to PENDING_MANUAL_REVIEW — holds fulfillment until verified.
	// ORDER_STATUS_PENDING_REVIEW does not exist in the proto enum — hold the order
	// in PAYMENT_INITIATED state while awaiting manual verification.
	if err := c.repo.UpdateOrderStatus(ctx, orderID, ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED); err != nil {
		return fmt.Errorf("HandleManualReviewRequested UpdateOrderStatus: %w", err)
	}
	if err := c.repo.SetPaymentStatus(ctx, orderID, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAYMENT_IN_PROGRESS); err != nil {
		appLogger.Warnf("HandleManualReviewRequested: SetPaymentStatus failed for order %s: %v", orderID, err)
	}

	// Validate the proof file exists in storage (non-fatal — auditing only).
	if evt.GetManualProofFileId() != "" {
		c.validateStorageFile(ctx, evt.GetManualProofFileId())
	}
	return nil
}

// HandleManualPaymentReviewed processes insuretech.payment.v1.payment.manual_review_completed.
// When approved, transitions the order to PAID; when rejected, to FAILED.
func (c *EventConsumer) HandleManualPaymentReviewed(ctx context.Context, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var evt paymenteventsv1.ManualPaymentReviewedEvent
	if err := protojson.Unmarshal(data, &evt); err != nil {
		appLogger.Errorf("HandleManualPaymentReviewed: unmarshal failed: %v", err)
		return fmt.Errorf("unmarshal ManualPaymentReviewedEvent: %w", err)
	}

	orderID := evt.GetOrderId()
	if orderID == "" {
		appLogger.Warnf("HandleManualPaymentReviewed: missing order_id — skipping")
		return nil
	}

	if evt.GetApproved() {
		appLogger.Infof("HandleManualPaymentReviewed: order=%s APPROVED by %s → PAID", orderID, evt.GetReviewedBy())
		if err := c.repo.UpdateOrderStatus(ctx, orderID, ordersv1.OrderStatus_ORDER_STATUS_PAID); err != nil {
			return fmt.Errorf("HandleManualPaymentReviewed (approve) UpdateOrderStatus: %w", err)
		}
		if err := c.repo.SetPaymentStatus(ctx, orderID, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID); err != nil {
			appLogger.Warnf("HandleManualPaymentReviewed: SetPaymentStatus failed: %v", err)
		}
		// Trigger receipt generation after approval.
		c.triggerReceiptDocgen(ctx, orderID, evt.GetPaymentId(), evt.GetTenantId(), evt.GetOrganisationId())
	} else {
		reason := evt.GetRejectionReason()
		if reason == "" {
			reason = evt.GetReviewNotes()
		}
		if reason == "" {
			reason = "Manual payment review rejected"
		}
		appLogger.Infof("HandleManualPaymentReviewed: order=%s REJECTED by %s: %s", orderID, evt.GetReviewedBy(), reason)
		if err := c.repo.SetFailureReason(ctx, orderID, reason); err != nil {
			return fmt.Errorf("HandleManualPaymentReviewed (reject) SetFailureReason: %w", err)
		}
		if err := c.repo.SetPaymentStatus(ctx, orderID, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAYMENT_FAILED); err != nil {
			appLogger.Warnf("HandleManualPaymentReviewed: SetPaymentStatus failed: %v", err)
		}
	}
	return nil
}

// HandleB2BPurchaseOrderApproved processes insuretech.b2b.v1.purchase_order.status_changed
// events where the new status is APPROVED (purchase_order_status == "APPROVED").
// When the B2B back-office approves a PO, the orders-service creates a corresponding
// order so the payment flow can begin for that organisation.
func (c *EventConsumer) HandleB2BPurchaseOrderApproved(ctx context.Context, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var evt b2beventsv1.PurchaseOrderStatusChangedEvent
	if err := protojson.Unmarshal(data, &evt); err != nil {
		appLogger.Errorf("HandleB2BPurchaseOrderApproved: unmarshal failed: %v", err)
		return fmt.Errorf("unmarshal PurchaseOrderStatusChangedEvent: %w", err)
	}

	// Only act on APPROVED transitions — GetNewStatus() returns a proto enum.
	newStatus := evt.GetNewStatus()
	newStatusStr := strings.ToUpper(newStatus.String())
	if !strings.Contains(newStatusStr, "APPROVED") {
		appLogger.Infof("HandleB2BPurchaseOrderApproved: ignoring status=%s (not APPROVED)", newStatusStr)
		return nil
	}

	poID := evt.GetPurchaseOrderId()
	orgID := evt.GetOrganisationId()
	if poID == "" {
		appLogger.Warnf("HandleB2BPurchaseOrderApproved: missing purchase_order_id — skipping")
		return nil
	}

	appLogger.Infof("HandleB2BPurchaseOrderApproved: po=%s org=%s status→APPROVED → creating order", poID, orgID)

	if c.orderSvc == nil {
		appLogger.Warnf("HandleB2BPurchaseOrderApproved: orderSvc not wired — skipping order creation for po=%s", poID)
		return nil
	}

	// Delegate to order service — it will create the order with B2B metadata
	// (purchase_order_id, organisation_id, tenant_id) from the event.
	// PurchaseOrderStatusChangedEvent has no TenantId or TotalAmount fields.
	// Pass empty/nil for now — the order service can enrich from the PO record.
	if err := c.orderSvc.CreateOrderForB2BPurchaseOrder(ctx, poID, orgID, "", nil); err != nil {
		return fmt.Errorf("HandleB2BPurchaseOrderApproved CreateOrderForB2BPurchaseOrder po=%s: %w", poID, err)
	}
	return nil
}

// ─── internal helpers ─────────────────────────────────────────────────────────

// triggerReceiptDocgen calls the docgen-service to generate a payment receipt PDF.
// Non-fatal: a warning is logged if docgen is unavailable (receipt can be re-generated on demand).
func (c *EventConsumer) triggerReceiptDocgen(ctx context.Context, orderID, paymentID, tenantID, orgID string) {
	if c.docgenClient == nil {
		appLogger.Warnf("triggerReceiptDocgen: docgen client not wired — skipping receipt for order=%s", orderID)
		return
	}
	_ = uuid.NewString() // correlationID reserved for future tracing
	req := &documentservicev1.GenerateDocumentRequest{
		// Template slug "payment-receipt" is seeded in the document-service DB.
		TemplateId: "payment-receipt",
		EntityType: "order",
		EntityId:   orderID,
	}
	resp, err := c.docgenClient.GenerateDocument(ctx, req)
	if err != nil {
		appLogger.Warnf("triggerReceiptDocgen: docgen call failed for order=%s: %v", orderID, err)
		return
	}
	appLogger.Infof("triggerReceiptDocgen: receipt queued — order=%s doc=%s", orderID, resp.GetDocumentId())
}

// validateStorageFile checks that a file_id exists in the storage-service.
// Non-fatal: failures are logged for auditing but do not block the consumer.
func (c *EventConsumer) validateStorageFile(ctx context.Context, fileID string) {
	if c.storageClient == nil {
		return
	}
	_, err := c.storageClient.GetFile(ctx, &storageservicev1.GetFileRequest{FileId: fileID})
	if err != nil {
		appLogger.Warnf("validateStorageFile: storage check failed for file=%s: %v", fileID, err)
	}
}

// ─── json helpers ─────────────────────────────────────────────────────────────

// flattenJSONMap decodes raw Kafka bytes into a flat string map.
// Used as a fallback for events from services that don't yet emit typed proto events.
func flattenJSONMap(data []byte) (map[string]string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		switch sv := v.(type) {
		case string:
			out[k] = strings.TrimSpace(sv)
		case float64:
			out[k] = fmt.Sprintf("%v", sv)
		case bool:
			if sv {
				out[k] = "true"
			} else {
				out[k] = "false"
			}
		}
	}
	return out, nil
}
