# API Endpoint Map

Auto-generated endpoint documentation

## Table of Contents

- [Admins](#admins)
- [Api-Keys](#api-keys)
- [Audit-Events](#audit-events)
- [Audit-Logs](#audit-logs)
- [Audit-Trail](#audit-trail)
- [Audits](#audits)
- [Auth](#auth)
- [Beneficiaries](#beneficiaries)
- [Biometric](#biometric)
- [Business](#business)
- [Catalog](#catalog)
- [Chat](#chat)
- [Check](#check)
- [Claims](#claims)
- [Commands](#commands)
- [Commission](#commission)
- [Commission-Payouts](#commission-payouts)
- [Commission-Statement](#commission-statement)
- [Commissions](#commissions)
- [Compliance-Logs](#compliance-logs)
- [Compliance-Reports](#compliance-reports)
- [Config](#config)
- [Credentials](#credentials)
- [Csrf](#csrf)
- [Current](#current)
- [Dashboards](#dashboards)
- [Decision](#decision)
- [Departments](#departments)
- [Devices](#devices)
- [Document-Templates](#document-templates)
- [Document-Types](#document-types)
- [Documents](#documents)
- [Download](#download)
- [Email](#email)
- [Employees](#employees)
- [Endorsements](#endorsements)
- [Faqs](#faqs)
- [Fraud](#fraud)
- [Fraud-Alerts](#fraud-alerts)
- [Fraud-Cases](#fraud-cases)
- [Fraud-Checks](#fraud-checks)
- [Fraud-Rules](#fraud-rules)
- [Grace-Period](#grace-period)
- [Health-Declaration](#health-declaration)
- [Individual](#individual)
- [Insurer-Products](#insurer-products)
- [Insurers](#insurers)
- [Invoice](#invoice)
- [Invoices](#invoices)
- [Jwks.Json](#jwks.json)
- [Knowledge-Base](#knowledge-base)
- [Kyc](#kyc)
- [Kyc-Verifications](#kyc-verifications)
- [Login](#login)
- [Logout](#logout)
- [Media](#media)
- [Members](#members)
- [Messages](#messages)
- [Metrics](#metrics)
- [My-Tasks](#my-tasks)
- [Notification-Preferences](#notification-preferences)
- [Notification-Templates](#notification-templates)
- [Notifications](#notifications)
- [Optimized](#optimized)
- [Orders](#orders)
- [Organisations](#organisations)
- [Otp](#otp)
- [Partners](#partners)
- [Password](#password)
- [Payment-Methods](#payment-methods)
- [Payments](#payments)
- [Pdf](#pdf)
- [Permissions](#permissions)
- [Photo](#photo)
- [Policies](#policies)
- [Process](#process)
- [Processing-Jobs](#processing-jobs)
- [Products](#products)
- [Profile](#profile)
- [Purchase-Orders](#purchase-orders)
- [Queries](#queries)
- [Quotes](#quotes)
- [Receipt](#receipt)
- [References](#references)
- [Refund](#refund)
- [Refunds](#refunds)
- [Register](#register)
- [Reminders](#reminders)
- [Renewal-Schedule](#renewal-schedule)
- [Report-Definitions](#report-definitions)
- [Report-Executions](#report-executions)
- [Report-Schedules](#report-schedules)
- [Reports](#reports)
- [Revenue-Share](#revenue-share)
- [Risk](#risk)
- [Risk-Score](#risk-score)
- [Roles](#roles)
- [Search](#search)
- [Sessions](#sessions)
- [Status](#status)
- [Tasks](#tasks)
- [Telemetry](#telemetry)
- [Tenants](#tenants)
- [Thumbnail](#thumbnail)
- [Tickets](#tickets)
- [Token](#token)
- [Totp](#totp)
- [Transactions](#transactions)
- [Transcript](#transcript)
- [Unknown](#unknown)
- [Upcoming](#upcoming)
- [Usage](#usage)
- [Webhook](#webhook)
- [Webhooks](#webhooks)
- [Workflow-Definitions](#workflow-definitions)
- [Workflow-History](#workflow-history)
- [Workflow-Instances](#workflow-instances)
- [Workflow-Tasks](#workflow-tasks)

## Admins

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/b2b/organisations/{organisation_id}/admins` | Resource | `POST` | POST: Assign a platform user as an OrgAdmin |

## Api-Keys

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/api-keys` | Collection | `GET`, `POST` | GET: List API keys for owner<br>POST: Validate API key and check scopes |
| `/v1/api-keys/{api_key_id}` | Resource | `DELETE`, `GET`, `POST` | DELETE: Revoke API key<br>GET: Get API key details<br>POST: Rotate API key |
| `/v1/auth/api-keys` | Collection | `GET`, `POST` | GET: List API keys for an owner<br>POST: Create a new API key for a user or service |
| `/v1/auth/api-keys/{key_id}:revoke` | Custom Action: `revoke` | `POST` | POST: Revoke an API key |
| `/v1/auth/api-keys/{key_id}:rotate` | Custom Action: `rotate` | `POST` | POST: Rotate an API key (generates new key, marks old one for graceful expiry) |

## Audit-Events

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/audit-events` | Collection | `GET`, `POST` | GET: Get audit events<br>POST: Create audit event |

## Audit-Logs

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/audit-logs` | Collection | `GET`, `POST` | GET: Get audit logs for entity<br>POST: Create audit log |

## Audit-Trail

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/entities/{entity_type}/{entity_id}/audit-trail` | Resource | `GET` | GET: Get audit trail for entity |

## Audits

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/authz/audits` | Collection | `GET` | GET: List access decision audits |

## Auth

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/voice-biometric:initiate` | Custom Action: `initiate` | `POST` | POST: 🎤 Voice Biometric Auth (Sprint 1 |
| `/v1/auth/voice-biometric:submit` | Custom Action: `submit` | `POST` | POST: Submit voice sample |
| `/v1/auth/voice-biometric:verify` | Custom Action: `verify` | `POST` | POST: Verify voice session |
| `/v1/auth/voice-sessions` | Collection | `POST` | POST: ── Voice Sessions ── |
| `/v1/auth/voice-sessions/{voice_session_id}` | Resource | `GET` | GET: Get voice session |
| `/v1/auth/voice-sessions/{voice_session_id}:end` | Custom Action: `end` | `POST` | POST: End voice session |

## Beneficiaries

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/beneficiaries` | Collection | `GET` | GET: List beneficiaries (admin) |
| `/v1/beneficiaries/{beneficiary_id}` | Resource | `GET`, `PATCH` | GET: Get beneficiary details<br>PATCH: Update beneficiary |

## Biometric

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/biometric:authenticate` | Custom Action: `authenticate` | `POST` | POST: Authenticate using a device-bound biometric token (mobile only) |

## Business

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/beneficiaries/business` | Collection | `POST` | POST: Create business beneficiary |

## Catalog

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/b2b/purchase-orders/catalog` | Collection | `GET` | GET: List purchasable product plans for purchase orders |

## Chat

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/ai/chat` | Collection | `POST` | POST: Chat with AI agent |

## Check

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/authz/check` | Collection | `POST` | POST: CheckAccess — single authorization check (gateway + per-service interceptor) |
| `/v1/authz/check:batch` | Custom Action: `batch` | `POST` | POST: BatchCheckAccess — check multiple (sub, dom, obj, act) tuples in one call |

## Claims

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/ai/claims:evaluate` | Custom Action: `evaluate` | `POST` | POST: Evaluate claim |
| `/v1/claims` | Collection | `POST` | POST: Submit claim |
| `/v1/claims/{claim_id}` | Resource | `GET` | GET: Get claim details |
| `/v1/claims/{claim_id}:approve` | Custom Action: `approve` | `POST` | POST: Approve claim |
| `/v1/claims/{claim_id}:dispute` | Custom Action: `dispute` | `POST` | POST: Dispute claim (by customer) |
| `/v1/claims/{claim_id}:reject` | Custom Action: `reject` | `POST` | POST: Reject claim |
| `/v1/claims/{claim_id}:request-documents` | Custom Action: `request-documents` | `POST` | POST: Request more documents from claimant |
| `/v1/claims/{claim_id}:settle` | Custom Action: `settle` | `POST` | POST: Settle claim |
| `/v1/users/{customer_id}/claims` | Resource | `GET` | GET: List user claims |

## Commands

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/voice-sessions/{voice_session_id}/commands` | Resource | `POST` | POST: Process voice command |

## Commission

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/partners/{partner_id}/commission` | Resource | `GET`, `PUT` | GET: Partner Commission & Financials<br>PUT: Update commission structure |

## Commission-Payouts

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/commission-payouts` | Collection | `POST` | POST: Create payout batch |
| `/v1/commission-payouts/{payout_id}` | Resource | `POST` | POST: Process payout |

## Commission-Statement

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/recipients/{recipient_id}/commission-statement` | Resource | `GET` | GET: Get commission statement |

## Commissions

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/commissions` | Collection | `GET`, `POST` | GET: List commissions for recipient<br>POST: Calculate and record commission for policy |
| `/v1/commissions/{commission_id}` | Resource | `GET` | GET: Get commission details |

## Compliance-Logs

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/compliance-logs` | Collection | `GET`, `POST` | GET: Get compliance logs<br>POST: Create compliance log |

## Compliance-Reports

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/compliance-reports` | Collection | `POST` | POST: Generate compliance report |

## Config

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/authz/portals/{portal}/config` | Resource | `GET`, `PATCH` | GET: Get portal config<br>PATCH: Update portal config |
| `/v1/insurers/{insurer_id}/config` | Resource | `PUT` | PUT: Update insurer config |
| `/v1/tenants/{tenant_id}/config` | Resource | `GET`, `PUT` | GET: Get tenant config<br>PUT: Update tenant config |

## Credentials

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/partners/{partner_id}/credentials` | Resource | `GET` | GET: Partner Integration |
| `/v1/partners/{partner_id}/credentials:rotate` | Custom Action: `rotate` | `POST` | POST: Rotate partner a p i key |

## Csrf

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/csrf:validate` | Custom Action: `validate` | `POST` | POST: Validate CSRF token (server-side sessions only) |

## Current

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/session/current` | Collection | `GET` | GET: Get current user's active session |

## Dashboards

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/analytics/dashboards` | Collection | `POST` | POST: Create dashboard |
| `/v1/analytics/dashboards/{dashboard_id}` | Resource | `GET` | GET: Get dashboard |

## Decision

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/quotes/{quote_id}/decision` | Resource | `GET` | GET: Get underwriting decision |

## Departments

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/b2b/departments` | Collection | `GET`, `POST` | GET: List departments for the authenticated organisation<br>POST: Create a new department |
| `/v1/b2b/departments/{department_id}` | Resource | `DELETE`, `GET`, `PATCH` | DELETE: Soft-delete a department (only if no active employees)<br>GET: Get a single department<br>PATCH: Update a department's name |

## Devices

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/iot/devices` | Collection | `POST` | POST: Register device |
| `/v1/iot/devices/{device_id}` | Resource | `GET` | GET: Get device status |
| `/v1/iot/devices/{device_id}:deactivate` | Custom Action: `deactivate` | `POST` | POST: Deactivate device |

## Document-Templates

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/document-templates` | Collection | `GET`, `POST` | GET: List templates<br>POST: Create template |
| `/v1/document-templates/{template_id}` | Resource | `DELETE`, `GET`, `PATCH` | DELETE: Delete template<br>GET: Get template<br>PATCH: Update template |
| `/v1/document-templates/{template_id}:deactivate` | Custom Action: `deactivate` | `POST` | POST: Deactivate template |

## Document-Types

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/document-types` | Collection | `GET` | GET: List document types |

## Documents

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/ai/documents:analyze` | Custom Action: `analyze` | `POST` | POST: Analyze document |
| `/v1/auth/documents/{user_document_id}` | Resource | `DELETE`, `GET`, `PATCH` | DELETE: Delete user document<br>GET: Get user document<br>PATCH: Update user document |
| `/v1/auth/documents/{user_document_id}:verify` | Custom Action: `verify` | `POST` | POST: ── Document Verification (Admin) ── |
| `/v1/auth/users/{user_id}/documents` | Resource | `GET`, `POST` | GET: List user documents<br>POST: Upload user document |
| `/v1/claims/{claim_id}/documents` | Resource | `POST` | POST: Upload claim document |
| `/v1/documents` | Collection | `POST` | POST: Generate document |
| `/v1/documents/{document_id}` | Resource | `DELETE`, `GET` | DELETE: Delete document<br>GET: Get document |
| `/v1/entities/{entity_type}/{entity_id}/documents` | Resource | `GET` | GET: List documents for entity |
| `/v1/kyc-verifications/{kyc_verification_id}/documents` | Resource | `POST` | POST: Upload document |

## Download

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/documents/{document_id}/download` | Resource | `GET` | GET: Download document |
| `/v1/media/{media_id}/download` | Resource | `GET` | GET: Download media file |
| `/v1/report-executions/{report_execution_id}/download` | Resource | `GET` | GET: Download report |

## Email

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/email/verify` | Collection | `POST` | POST: Verify email address using OTP (must call before email login is allowed) |

## Employees

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/b2b/employees` | Collection | `GET`, `POST` | GET: List employees for the authenticated organisation<br>POST: Create a new employee |
| `/v1/b2b/employees/{employee_uuid}` | Resource | `DELETE`, `GET`, `PATCH` | DELETE: Soft-delete an employee record<br>GET: Get a single employee by employee_uuid<br>PATCH: Update an existing employee's details |

## Endorsements

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/endorsements/{endorsement_id}` | Resource | `GET`, `POST` | GET: Get endorsement<br>POST: Reject endorsement |
| `/v1/policies/{policy_id}/endorsements` | Resource | `GET`, `POST` | GET: List endorsements for policy<br>POST: Request endorsement |

## Faqs

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/faqs` | Collection | `GET`, `POST` | GET: List FAQs<br>POST: Create FAQ |
| `/v1/faqs/{faq_id}` | Resource | `DELETE`, `PATCH` | DELETE: Delete FAQ<br>PATCH: Update FAQ |

## Fraud

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/ai/fraud:detect` | Custom Action: `detect` | `POST` | POST: Detect fraud |

## Fraud-Alerts

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/fraud-alerts` | Collection | `GET` | GET: List fraud alerts |
| `/v1/fraud-alerts/{fraud_alert_id}` | Resource | `GET` | GET: Get fraud alert |

## Fraud-Cases

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/fraud-cases` | Collection | `POST` | POST: Create fraud case |
| `/v1/fraud-cases/{fraud_case_id}` | Resource | `GET`, `PATCH` | GET: Get fraud case<br>PATCH: Update fraud case |

## Fraud-Checks

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/fraud-checks` | Collection | `POST` | POST: Check for fraud |

## Fraud-Rules

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/fraud-rules` | Collection | `GET`, `POST` | GET: List fraud rules<br>POST: Create fraud rule |
| `/v1/fraud-rules/{rule_id}` | Resource | `PATCH` | PATCH: Update fraud rule |
| `/v1/fraud-rules/{rule_id}:activate` | Custom Action: `activate` | `POST` | POST: Activate fraud rule |
| `/v1/fraud-rules/{rule_id}:deactivate` | Custom Action: `deactivate` | `POST` | POST: Deactivate fraud rule |

## Grace-Period

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/policies/{policy_id}/grace-period` | Resource | `GET` | GET: Get grace period status |

## Health-Declaration

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/quotes/{quote_id}/health-declaration` | Resource | `GET`, `POST` | GET: Get health declaration<br>POST: Submit health declaration |

## Individual

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/beneficiaries/individual` | Collection | `POST` | POST: Create individual beneficiary |

## Insurer-Products

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/insurer-products/{insurer_product_id}` | Resource | `GET`, `PATCH` | GET: Get insurer product<br>PATCH: Update insurer product |

## Insurers

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/insurers` | Collection | `GET`, `POST` | GET: List insurers<br>POST: Create insurer |
| `/v1/insurers/{insurer_id}` | Resource | `GET`, `PATCH` | GET: Get insurer details<br>PATCH: Update insurer |

## Invoice

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/orders/{order_id}/invoice` | Resource | `GET` | GET: Get invoice by order ID — used by orders-service to link invoice after creation |

## Invoices

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/invoices` | Collection | `GET`, `POST` | GET: List invoices with optional filters<br>POST: Create a new invoice for an order (B2C) or purchase order (B2B) |
| `/v1/invoices/{invoice_id}` | Resource | `GET` | GET: Get a single invoice by ID |
| `/v1/invoices/{invoice_id}:cancel` | Custom Action: `cancel` | `POST` | POST: Cancel an invoice (only allowed before PAID) |
| `/v1/invoices/{invoice_id}:generate-pdf` | Custom Action: `generate-pdf` | `POST` | POST: Trigger async invoice PDF generation |
| `/v1/invoices/{invoice_id}:issue` | Custom Action: `issue` | `POST` | POST: Issue an invoice — transitions from DRAFT → ISSUED and sends to customer/org |
| `/v1/invoices/{invoice_id}:mark-paid` | Custom Action: `mark-paid` | `POST` | POST: Mark invoice as paid — called by payment-service after payment confirmation |

## Jwks.Json

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/.well-known/jwks.json` | Collection | `GET` | GET: GetJWKS — serves the RS256 public key set for JWT verification |
| `/v1/auth/.well-known/jwks.json` | Collection | `GET` | GET: 🔑 JWKS 🔑 |

## Knowledge-Base

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/knowledge-base` | Collection | `POST` | POST: Create Knowledge Base Article |
| `/v1/knowledge-base/{article_id}` | Resource | `DELETE`, `PATCH` | DELETE: Delete Knowledge Base Article<br>PATCH: Update Knowledge Base Article |
| `/v1/knowledge-base/{slug}` | Resource | `GET` | GET: Get knowledge base article |

## Kyc

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/kyc/{kyc_id}:approve` | Custom Action: `approve` | `POST` | POST: Approve k y c |
| `/v1/auth/kyc/{kyc_id}:reject` | Custom Action: `reject` | `POST` | POST: Reject k y c |
| `/v1/auth/users/{user_id}/kyc` | Resource | `GET`, `POST` | GET: Get k y c status<br>POST: ── KYC Verification ── |
| `/v1/auth/users/{user_id}/kyc:complete` | Custom Action: `complete` | `POST` | POST: Complete k y c session |
| `/v1/auth/users/{user_id}/kyc:submit-frame` | Custom Action: `submit-frame` | `POST` | POST: Submit k y c frame |
| `/v1/beneficiaries/{beneficiary_id}/kyc` | Resource | `POST` | POST: Complete KYC |

## Kyc-Verifications

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/kyc-verifications` | Collection | `POST` | POST: Start KYC verification |
| `/v1/kyc-verifications/{kyc_verification_id}` | Resource | `GET`, `POST` | GET: Get KYC verification<br>POST: Reject KYC |
| `/v1/kyc-verifications:pending` | Custom Action: `pending` | `GET` | GET: List pending KYC verifications (admin review queue) |

## Login

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/email/login` | Collection | `POST` | POST: Login via email + OTP (Business Beneficiary / System User only → |
| `/v1/auth/login` | Collection | `POST` | POST: Login with credentials |

## Logout

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/logout` | Collection | `POST` | POST: Logout |

## Media

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/entities/{entity_type}/{entity_id}/media` | Resource | `GET` | GET: List media files for entity |
| `/v1/media` | Collection | `POST` | POST: Upload media file |
| `/v1/media/{media_id}` | Resource | `DELETE`, `GET` | DELETE: Delete media file<br>GET: Get media file |
| `/v1/media/{media_id}:validate` | Custom Action: `validate` | `POST` | POST: Validate media file |

## Members

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/b2b/organisations/{organisation_id}/members` | Resource | `GET`, `POST` | GET: List members for an organisation<br>POST: Add a platform user as an OrgMember |
| `/v1/b2b/organisations/{organisation_id}/members/{member_id}` | Resource | `DELETE` | DELETE: Remove an OrgMember from the organisation |

## Messages

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/tickets/{ticket_id}/messages` | Resource | `POST` | POST: Add ticket message |

## Metrics

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/analytics/metrics` | Collection | `POST` | POST: Get metrics |

## My-Tasks

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/tasks/my-tasks` | Collection | `GET` | GET: List my tasks |
| `/v1/workflow-tasks/my-tasks` | Collection | `GET` | GET: Get my tasks |

## Notification-Preferences

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/users/{user_id}/notification-preferences` | Resource | `PATCH` | PATCH: ── Notification Preferences ── |
| `/v1/users/{user_id}/notification-preferences` | Resource | `PUT` | PUT: Update notification preferences |

## Notification-Templates

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/notification-templates` | Collection | `POST` | POST: Create notification template |
| `/v1/notification-templates/{template_id}` | Resource | `PATCH` | PATCH: Update notification template |
| `/v1/notification-templates/{template_id}:deactivate` | Custom Action: `deactivate` | `POST` | POST: Deactivate notification template |

## Notifications

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/notifications` | Collection | `POST` | POST: Send notification |
| `/v1/notifications/{notification_id}` | Resource | `GET` | GET: Get notification status |
| `/v1/notifications:mark-as-read` | Custom Action: `mark-as-read` | `POST` | POST: Mark as read |
| `/v1/notifications:send-bulk` | Custom Action: `send-bulk` | `POST` | POST: Send bulk notifications |
| `/v1/users/{user_id}/notifications` | Resource | `GET` | GET: Get user notifications |

## Optimized

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/media/{media_id}/optimized` | Resource | `GET` | GET: Download optimized version |

## Orders

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/orders` | Collection | `GET`, `POST` | GET: List orders with optional filters<br>POST: Create a new order from an approved quotation |
| `/v1/orders/{order_id}` | Resource | `GET` | GET: Get a single order by ID |
| `/v1/orders/{order_id}:cancel` | Custom Action: `cancel` | `POST` | POST: Cancel an order (only allowed before PAID or POLICY_ISSUED) |
| `/v1/orders/{order_id}:confirm-payment` | Custom Action: `confirm-payment` | `POST` | POST: Confirm payment for an order (called by payment gateway callback) |
| `/v1/orders/{order_id}:pay` | Custom Action: `pay` | `POST` | POST: Initiate payment for a pending order |

## Organisations

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/b2b/organisations` | Collection | `GET`, `POST` | GET: List all organisations (SuperAdmin: all; BizAdmin: own only)<br>POST: Create a new organisation (SuperAdmin only) |
| `/v1/b2b/organisations/{organisation_id}` | Resource | `DELETE`, `GET`, `PATCH` | DELETE: Soft-delete an organisation and revoke its memberships<br>GET: Get a single organisation by ID<br>PATCH: Update an organisation's profile |

## Otp

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/email/otp:send` | Custom Action: `send` | `POST` | POST: Send email OTP (verification or login) |
| `/v1/auth/otp:resend` | Custom Action: `resend` | `POST` | POST: Resend OTP (invalidates previous OTP, generates fresh one) |
| `/v1/auth/otp:send` | Custom Action: `send` | `POST` | POST: Send OTP for verification |
| `/v1/auth/otp:verify` | Custom Action: `verify` | `POST` | POST: Verify OTP |

## Partners

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/partners` | Collection | `GET`, `POST` | GET: List partners<br>POST: Partner Management |
| `/v1/partners/{partner_id}` | Resource | `DELETE`, `GET`, `PATCH` | DELETE: Delete partner<br>GET: Get partner<br>PATCH: Update partner |
| `/v1/partners/{partner_id}:update-status` | Custom Action: `update-status` | `POST` | POST: Update partner status |
| `/v1/partners/{partner_id}:verify` | Custom Action: `verify` | `POST` | POST: Partner Verification & Onboarding |

## Password

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/email/password:reset` | Custom Action: `reset` | `POST` | POST: Complete password reset using email OTP |
| `/v1/auth/email/password:reset-request` | Custom Action: `reset-request` | `POST` | POST: Request password reset via email OTP |
| `/v1/auth/password:change` | Custom Action: `change` | `POST` | POST: Change password |
| `/v1/auth/password:reset` | Custom Action: `reset` | `POST` | POST: Reset password |

## Payment-Methods

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/users/{user_id}/payment-methods` | Resource | `GET`, `POST` | GET: Payment Methods<br>POST: Add payment method |

## Payments

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/mfs/payments` | Collection | `POST` | POST: Initiate payment |
| `/v1/payments` | Collection | `GET`, `POST` | GET: List payments<br>POST: Payment Processing |
| `/v1/payments/{payment_id}` | Resource | `GET` | GET: Get payment |
| `/v1/payments/{payment_id}:generate-receipt` | Custom Action: `generate-receipt` | `POST` | POST: Trigger async receipt PDF generation after payment is verified |
| `/v1/payments/{payment_id}:review` | Custom Action: `review` | `POST` | POST: Admin/agent reviews and approves or rejects a manual payment proof |
| `/v1/payments/{payment_id}:submit-proof` | Custom Action: `submit-proof` | `POST` | POST: Manual bank transfer: customer submits payment proof (scanned deposit slip / screenshot) |
| `/v1/payments/{payment_id}:verify` | Custom Action: `verify` | `POST` | POST: Verify payment |
| `/v1/payments:reconcile` | Custom Action: `reconcile` | `POST` | POST: Reconciliation |

## Pdf

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/invoices/{invoice_id}/pdf` | Resource | `GET` | GET: Get invoice PDF (pre-signed download URL or file ID) |

## Permissions

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/authz/users/{user_id}/permissions` | Resource | `GET` | GET: GetUserPermissions — resolves all effective permissions for a user in a domain |

## Photo

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/users/{user_id}/profile/photo:upload-url` | Custom Action: `upload-url` | `POST` | POST: ── Profile Photo Upload URL ── |

## Policies

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/authz/policies` | Collection | `GET`, `POST` | GET: List policy rules<br>POST: Create policy rule |
| `/v1/authz/policies/{policy_id}` | Resource | `DELETE`, `PATCH` | DELETE: Delete policy rule<br>PATCH: Update policy rule |
| `/v1/policies` | Collection | `POST` | POST: Create policy |
| `/v1/policies/{policy_id}` | Resource | `GET`, `PATCH`, `POST` | GET: Get policy details<br>PATCH: Update policy<br>POST: Revive lapsed policy |
| `/v1/policies/{policy_id}:cancel` | Custom Action: `cancel` | `POST` | POST: Cancel policy |
| `/v1/policies/{policy_id}:generate-document` | Custom Action: `generate-document` | `POST` | POST: Generate policy document |
| `/v1/policies/{policy_id}:issue` | Custom Action: `issue` | `POST` | POST: Issue policy (explicit transition from pending to active after payment) |
| `/v1/policies/{policy_id}:renew` | Custom Action: `renew` | `POST` | POST: Renew policy |
| `/v1/users/{customer_id}/policies` | Resource | `GET` | GET: List user policies |

## Process

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/media/{media_id}/process` | Resource | `POST` | POST: Request processing (OCR, optimization, etc |

## Processing-Jobs

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/processing-jobs` | Collection | `GET` | GET: List processing jobs |
| `/v1/processing-jobs/{job_id}` | Resource | `GET` | GET: Get processing job status |

## Products

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/insurers/{insurer_id}/products` | Resource | `GET`, `POST` | GET: List insurer products<br>POST: Add insurer product |
| `/v1/products` | Collection | `GET`, `POST` | GET: List all active products<br>POST: Create product (admin) |
| `/v1/products/{product_id}` | Resource | `GET`, `PATCH` | GET: Get product details<br>PATCH: Update product (admin) |
| `/v1/products/{product_id}:activate` | Custom Action: `activate` | `POST` | POST: Activate product |
| `/v1/products/{product_id}:calculate-premium` | Custom Action: `calculate-premium` | `POST` | POST: Calculate premium |
| `/v1/products/{product_id}:deactivate` | Custom Action: `deactivate` | `POST` | POST: Deactivate product |
| `/v1/products/{product_id}:discontinue` | Custom Action: `discontinue` | `POST` | POST: Discontinue product |
| `/v1/products:search` | Custom Action: `search` | `GET` | GET: Search products |

## Profile

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/users/{user_id}/profile` | Resource | `GET`, `PATCH`, `POST` | GET: Get user profile<br>PATCH: Update user profile<br>POST: Create user profile |

## Purchase-Orders

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/b2b/purchase-orders` | Collection | `GET`, `POST` | GET: List purchase orders for the authenticated organisation<br>POST: Create a purchase order for a product plan |
| `/v1/b2b/purchase-orders/{purchase_order_id}` | Resource | `GET` | GET: Get a single purchase order |

## Queries

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/analytics/queries:run` | Custom Action: `run` | `POST` | POST: Run custom query |

## Quotes

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/beneficiaries/{beneficiary_id}/quotes` | Resource | `GET` | GET: List quotes for beneficiary |
| `/v1/quotes` | Collection | `POST` | POST: Request premium quote |
| `/v1/quotes/{quote_id}` | Resource | `GET`, `POST` | GET: Get quote details<br>POST: Convert quote to policy |

## Receipt

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/payments/{payment_id}/receipt` | Resource | `GET` | GET: Retrieve generated receipt (pre-signed URL or file ID) |

## References

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/payments/provider/{provider}/references/{provider_reference}` | Resource | `GET` | GET: Lookup by provider-specific reference (e |

## Refund

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/policies/{policy_id}/refund` | Resource | `POST` | POST: Request refund |

## Refunds

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/mfs/refunds` | Collection | `POST` | POST: Execute refund |
| `/v1/payments/{payment_id}/refunds` | Resource | `POST` | POST: Refund Management |
| `/v1/policies/{policy_id}/refunds` | Resource | `POST` | POST: Calculate refund amount |
| `/v1/refunds` | Collection | `GET` | GET: List refunds |
| `/v1/refunds/{refund_id}` | Resource | `GET`, `POST` | GET: Get refund<br>POST: Process refund payment |

## Register

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/email/register` | Collection | `POST` | POST: Register a portal user with email (requires email, triggers email |
| `/v1/auth/register` | Collection | `POST` | POST: Register new user |

## Reminders

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/renewal-schedules/{renewal_schedule_id}/reminders` | Resource | `POST` | POST: Send renewal reminder |

## Renewal-Schedule

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/policies/{policy_id}/renewal-schedule` | Resource | `GET` | GET: Get renewal schedule |

## Report-Definitions

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/report-definitions` | Collection | `GET` | GET: List report definitions |

## Report-Executions

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/report-executions` | Collection | `GET` | GET: List report executions |
| `/v1/report-executions/{report_execution_id}` | Resource | `GET` | GET: Get report execution |

## Report-Schedules

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/report-schedules` | Collection | `GET`, `POST` | GET: List report schedules<br>POST: Create report schedule |

## Reports

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/analytics/reports/{report_id}:generate` | Custom Action: `generate` | `POST` | POST: Generate report |
| `/v1/analytics/reports:schedule` | Custom Action: `schedule` | `POST` | POST: Schedule report |
| `/v1/reports/{report_definition_id}` | Resource | `POST` | POST: Execute report |

## Revenue-Share

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/insurers/{insurer_id}/revenue-share` | Resource | `GET` | GET: Get revenue share report |

## Risk

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/ai/risk:assess` | Custom Action: `assess` | `POST` | POST: Assess risk |
| `/v1/iot/devices/{device_id}/risk` | Resource | `GET` | GET: Get risk assessment |

## Risk-Score

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/beneficiaries/{beneficiary_id}/risk-score` | Resource | `POST` | POST: Update risk score |

## Roles

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/authz/roles` | Collection | `GET`, `POST` | GET: List roles<br>POST: Create role |
| `/v1/authz/roles/{role_id}` | Resource | `DELETE`, `GET`, `PATCH` | DELETE: Delete role<br>GET: Get role<br>PATCH: Update role |
| `/v1/authz/users/{user_id}/roles` | Resource | `GET`, `POST` | GET: List user roles<br>POST: AssignRole — assign a role to a user within domain (portal:tenant_id) |
| `/v1/authz/users/{user_id}/roles/{role_id}` | Resource | `DELETE` | DELETE: Remove role |

## Search

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/knowledge-base/search` | Collection | `GET` | GET: Search knowledge base |

## Sessions

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/sessions/{session_id}` | Resource | `DELETE`, `GET` | DELETE: Revoke a specific session<br>GET: Get session details |
| `/v1/auth/users/{user_id}/sessions` | Resource | `GET` | GET: List all sessions for a user |
| `/v1/auth/users/{user_id}/sessions:revoke-all` | Custom Action: `revoke-all` | `POST` | POST: Revoke all sessions for a user (logout from all devices) |

## Status

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/orders/{order_id}/status` | Resource | `GET` | GET: Get lightweight order status |
| `/v1/tickets/{ticket_id}/status` | Resource | `PATCH` | PATCH: Update ticket status |

## Tasks

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/tasks` | Collection | `POST` | POST: Create task |
| `/v1/tasks/{task_id}` | Resource | `GET`, `PATCH`, `POST` | GET: Get task<br>PATCH: Update task<br>POST: Complete task |

## Telemetry

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/iot/telemetry` | Collection | `POST` | POST: Send telemetry data |

## Tenants

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/tenants` | Collection | `GET`, `POST` | GET: List tenants<br>POST: Create tenant |
| `/v1/tenants/{tenant_id}` | Resource | `GET`, `PATCH` | GET: Get tenant<br>PATCH: Update tenant |

## Thumbnail

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/media/{media_id}/thumbnail` | Resource | `GET` | GET: Download thumbnail |

## Tickets

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/tickets` | Collection | `GET`, `POST` | GET: List tickets<br>POST: Create ticket |
| `/v1/tickets/{ticket_id}` | Resource | `GET`, `POST` | GET: Get ticket<br>POST: Assign ticket |

## Token

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/token:refresh` | Custom Action: `refresh` | `POST` | POST: Refresh access token |
| `/v1/auth/token:validate` | Custom Action: `validate` | `POST` | POST: Validate token |

## Totp

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/auth/users/{user_id}/totp:disable` | Custom Action: `disable` | `POST` | POST: Disable t o t p |
| `/v1/auth/users/{user_id}/totp:enable` | Custom Action: `enable` | `POST` | POST: 🔐 TOTP / 2FA 🔐 |
| `/v1/auth/users/{user_id}/totp:verify` | Custom Action: `verify` | `POST` | POST: Verify t o t p |

## Transactions

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/mfs/transactions` | Collection | `GET` | GET: List transactions |
| `/v1/mfs/transactions/{mfs_transaction_id}` | Resource | `GET` | GET: Get transaction |

## Transcript

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/voice-sessions/{voice_session_id}/transcript` | Resource | `GET` | GET: Get transcript |

## Unknown

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/voice-sessions` | Collection | `POST` | POST: Start voice session |
| `/v1/voice-sessions/{voice_session_id}` | Resource | `GET`, `POST` | GET: Get voice session<br>POST: End voice session |

## Upcoming

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/renewals/upcoming` | Collection | `GET` | GET: List upcoming renewals |

## Usage

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/api-keys/{api_key_id}/usage` | Resource | `GET` | GET: Get usage statistics |

## Webhook

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/payments/webhook/{provider}` | Resource | `POST` | POST: Gateway webhook — called by API gateway when SSLCommerz/bKash/Nagad posts callback |

## Webhooks

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/mfs/webhooks/{provider}` | Resource | `POST` | POST: Process webhook |

## Workflow-Definitions

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/workflow-definitions` | Collection | `POST` | POST: Create workflow definition |
| `/v1/workflow-definitions/{workflow_definition_id}` | Resource | `GET` | GET: Get workflow definition |

## Workflow-History

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/entities/{entity_type}/{entity_id}/workflow-history` | Resource | `GET` | GET: Get workflow history for entity |

## Workflow-Instances

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/workflow-instances` | Collection | `POST` | POST: Start workflow instance |
| `/v1/workflow-instances/{workflow_instance_id}` | Resource | `GET` | GET: Get workflow instance |

## Workflow-Tasks

| Path | Type | Methods | Operations |
|------|------|---------|------------|
| `/v1/workflow-tasks/{task_id}` | Resource | `POST` | POST: Complete task |

## Statistics

- **Total Endpoints:** 260
- **Total Operations:** 335
- **Collections:** 75
- **Resources:** 111
- **Custom Actions:** 74
