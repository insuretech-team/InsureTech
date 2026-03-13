# 8. Security, Privacy, and Compliance Requirements (Detailed)

Security and compliance are business requirements: they protect customers, protect funds, enable partner trust, and satisfy IDRA/BFIU expectations.
This section translates SRS Section 7 controls into business-operational requirements.

## 8.1 Security Control Catalog (SEC)

### SEC-001

- **Business control requirement:** Use separate secret vault - AWS KMS/Azure Key Vault/HashiCorp, 90-day key rotation
- **Priority:** M1
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-002

- **Business control requirement:** Use Data Masking: NID (last 3 digits), phone (mask middle), email (mask username)
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-003

- **Business control requirement:** Follow PCI-DSS compliance for card flows - Approach: Hosted payment page (redirect model) - DO NOT store card data, Level: SAQ-A (simplest, for redirecting merchants), Requirements: Annual SAQ, quarterly ASV scans, TLS 1.3, Tokenization: Store only gateway tokens for recurring payments
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-004

- **Business control requirement:** Have AML/CFT detection hooks - Transaction Monitoring: 20+ automated rules for AML detection including Rapid purchases (>3 policies in 7 days), High-value premiums (>BDT 5 lakh), Frequent cancellations, Mismatched nominees, Geographic/payment anomalies
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-005

- **Business control requirement:** Have IDRA reporting capabilities following IDRA data format - Monthly Reports: Premium Collection (Form IC-1), Claims Intimation (Form IC-2), Quarterly Reports: Claims Settlement (IC-3), Financial Performance (IC-4), Annual Reports: FCR (Financial Condition Report), CARAMELS Framework Returns, Event-Based: Significant incidents (48hrs), fraud cases (7 days), Platform: Report generator with IDRA Excel templates, audit trail, 20-year archive
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-006

- **Business control requirement:** Have regular penetration testing - Penetration Testing: Pre-launch + annually (SISA InfoSec or international firm)
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-007

- **Business control requirement:** Have regular security audits from various security auditors and regulatory bodies and maintain compliance
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-008

- **Business control requirement:** DAST: OWASP ZAP/Burp Suite (weekly on staging)
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-009

- **Business control requirement:** SAST: SonarQube/Checkmarx (every commit, block critical vulnerabilities)
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-010

- **Business control requirement:** Virus scanning: ClamAV on uploaded files
- **Priority:** M
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-021

- **Business control requirement:** Implement API rate limiting per user/IP: 1000 requests/hour for authenticated users, 100 requests/hour for anonymous
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-022

- **Business control requirement:** Maintain separate encryption keys for different data types with hierarchical key management
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-023

- **Business control requirement:** Implement real-time security incident response with automated threat isolation
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-024

- **Business control requirement:** Perform continuous vulnerability assessment with automated patching for critical vulnerabilities
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-025

- **Business control requirement:** Implement zero-trust network architecture with microsegmentation
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-011

- **Business control requirement:** IDRA Monthly Reports: Generate Form IC-1 (Premium Collection) by 10th of each month with breakdown by product line, geographic region, partner channel in Excel format per IDRA template v2024
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-012

- **Business control requirement:** IDRA Monthly Reports: Generate Form IC-2 (Claims Intimation) by 10th of each month listing all new claims with policy number, claim amount, claim type, date of intimation
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-013

- **Business control requirement:** IDRA Quarterly Reports: Generate Form IC-3 (Claims Settlement) within 15 days of quarter-end showing settlement ratio, average TAT, pending >30 days breakdown
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-014

- **Business control requirement:** IDRA Quarterly Reports: Generate Form IC-4 (Financial Performance) within 20 days of quarter-end with premium earned, claims paid, commission paid, net profit/loss
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-015

- **Business control requirement:** IDRA Annual FCR: Generate Financial Condition Report (FCR) within 90 days of year-end including full CARAMELS framework assessment with external auditor sign-off
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-016

- **Business control requirement:** IDRA Event-Based Reporting: Report significant incidents (fraud >BDT 1L, data breach, system outage >4hrs) within 48 hours via IDRA portal
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-017

- **Business control requirement:** AML/CFT Concrete Triggers: Flag transactions matching: (1) >3 policies in 7 days, (2) Premium >BDT 5L without income proof, (3) Nominee mismatch with no relationship doc, (4) Payment from third-party account, (5) Frequent cancellations >2 in 30 days, (6) Geographic anomaly (policy in Dhaka, payment from remote district), (7) Multiple failed KYC attempts >3, (8) PEP match in screening
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-018

- **Business control requirement:** SAR Workflow: (1) System auto-flags suspicious transaction → (2) Compliance Officer reviews within 24hrs → (3) If confirmed suspicious, escalate to Business Admin+Focal Person → (4) Prepare SAR with evidence → (5) Submit to BFIU within 3 business days → (6) Mark account for enhanced monitoring → (7) Do NOT notify customer (tipping off prohibited)
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-019

- **Business control requirement:** Data Deletion Exceptions: Customer data deletion requests processed within 30 days EXCEPT: (a) Active policy holders (deletion after policy expiry+7yrs), (b) Ongoing claims (deletion after settlement+7yrs), (c) Under SAR investigation (deletion prohibited until case closed), (d) Regulatory hold (deletion requires IDRA/BFIU approval)
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-020

- **Business control requirement:** Right to Erasure Workflow: Customer submits deletion request → System validates exceptions → If eligible, anonymize PII while retaining transaction records → Generate deletion certificate → Notify customer within 30 days
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

## 8.2 AML/CFT Operating Model (Business View)

The platform must support configurable AML monitoring rules, alerting, investigation workflow, and STR/SAR filing support with strict auditability.
(See SRS Section 7.7.x for rule tables and workflows.)

## 8.3 IDRA Reporting and Record-Keeping (Business View)

The platform must retain and produce long-term records (policies, payments, claims, cancellations, approvals, customer communications) with retrieval capability within required SLAs.

[[[PAGEBREAK]]]
