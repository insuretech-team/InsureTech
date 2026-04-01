# Notification Topic Model

Notification is treated as a platform orchestrator, not a single-domain consumer.

## Group Model

- `customer_identity`
  - registration, KYC, account security
- `customer_commerce`
  - orders, invoices, payments, refunds, policy issuance/cancellation/renewal
- `claims_lifecycle`
  - claim submitted, docs requested, approved, rejected, settled, fraud escalation
- `renewal_retention`
  - renewal reminders, grace period, lapse, reinstatement
- `document_artifacts`
  - receipts, policy packs, generated documents, finalized storage artifacts
- `partner_ops`
  - partner/agent/employer/insurer coordination events in the multi-insurer network
- `support_escalation`
  - ticketing, tasks, SLA breaches, operational escalations
- `marketing_campaign`
  - audience segmentation and approved campaign fanout
- `iot_risk_alerts`
  - telematics and threshold-based customer alerts
- `compliance_ops`
  - AML, fraud, regulator, and internal-control alerts
- `notification_state`
  - outbound lifecycle emitted by notification service itself
- `webhook_fanout`
  - downstream external delivery requests

## Subscription Profiles

- `transactional_core`
  - minimal customer-critical event set
- `customer_core`
  - default profile for customer-facing notification flows
- `operations`
  - internal and partner/operator-heavy profile
- `platform_all`
  - all known groups for broad integration or staging validation

## Dynamic Controls

The service supports env-driven topic resolution:

- `NOTIFICATION_SUBSCRIPTION_PROFILE`
- `NOTIFICATION_TOPIC_GROUPS`
- `NOTIFICATION_DISABLED_TOPIC_GROUPS`
- `NOTIFICATION_TOPIC_ALLOWLIST`
- `NOTIFICATION_TOPIC_DENYLIST`
- `NOTIFICATION_EXTRA_TOPICS`
- `NOTIFICATION_INCLUDE_RESERVED_TOPICS`

## Current vs Reserved

Each topic descriptor is marked as:

- `CurrentContract`
  - a topic with an existing producer/proto contract today
- `Reserved`
  - provisioned in the model because the SRS/platform requires it, even if the producer is not fully implemented yet

Reserved topics are excluded from consumer subscriptions unless explicitly enabled.
