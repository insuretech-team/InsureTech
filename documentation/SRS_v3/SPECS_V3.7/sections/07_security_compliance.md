# 7. Security & Compliance Requirements

& Compliance Requirements



The LabAid InsureTech Platform implements a **Zero Trust Security Model** with defense-in-depth strategies:

**Core Security Principles:**
1. **Never Trust, Always Verify:** All users and devices authenticated and authorized
2. **Least Privilege Access:** Minimum required permissions for each role
3. **Assume Breach:** Monitor and respond as if compromise has occurred
4. **Encrypt Everything:** Data protection at all layers and states
5. **Continuous Monitoring:** Real-time threat detection and response

### 7.1 Security Infrastructure & Key Management

| ID      | Requirement Description                                                                                                                                                                                                                                                                                                                                                                                                                                                      | Priority |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| SEC-001 | The system shall use separate secret vault - AWS KMS/Azure Key Vault/HashiCorp, 90-day key rotation                                                                                                                                                                                                                                                                                                                                                                          | M1        |
| SEC-002 | The system shall use Data Masking: NID (last 3 digits), phone (mask middle), email (mask username)                                                                                                                                                                                                                                                                                                                                                                           | M2        |
| SEC-003 | The system shall follow PCI-DSS compliance for card flows - Approach: Hosted payment page (redirect model) - DO NOT store card data, Level: SAQ-A (simplest, for redirecting merchants), Requirements: Annual SAQ, quarterly ASV scans, TLS 1.3, Tokenization: Store only gateway tokens for recurring payments                                                                                                                                                              | M2        |
| SEC-004 | The system shall have AML/CFT detection hooks - Transaction Monitoring: 20+ automated rules for AML detection including Rapid purchases (>3 policies in 7 days), High-value premiums (>BDT 5 lakh), Frequent cancellations, Mismatched nominees, Geographic/payment anomalies                                                                                                                                                                                                | D        |
| SEC-005 | The system shall have IDRA reporting capabilities following IDRA data format - Monthly Reports: Premium Collection (Form IC-1), Claims Intimation (Form IC-2), Quarterly Reports: Claims Settlement (IC-3), Financial Performance (IC-4), Annual Reports: FCR (Financial Condition Report), CARAMELS Framework Returns, Event-Based: Significant incidents (48hrs), fraud cases (7 days), Platform: Report generator with IDRA Excel templates, audit trail, 20-year archive | D        |
| SEC-006 | The system shall have regular penetration testing - Penetration Testing: Pre-launch + annually (SISA InfoSec or international firm)                                                                                                                                                                                                                                                                                                                                          | D        |
| SEC-007 | The system shall have regular security audits from various security auditors and regulatory bodies and maintain compliance                                                                                                                                                                                                                                                                                                                                                   | D        |
| SEC-008 | DAST: OWASP ZAP/Burp Suite (weekly on staging)                                                                                                                                                                                                                                                                                                                                                                                                                               | D        |
| SEC-009 | SAST: SonarQube/Checkmarx (every commit, block critical vulnerabilities)                                                                                                                                                                                                                                                                                                                                                                                                     | D        |
| SEC-010 | Virus scanning: ClamAV on uploaded files                                                                                                                                                                                                                                                                                                                                                                                                                                     | M        |
| SEC-021 | The system shall implement API rate limiting per user/IP: 1000 requests/hour for authenticated users, 100 requests/hour for anonymous | M2 |
| SEC-022 | The system shall maintain separate encryption keys for different data types with hierarchical key management | M2 |
| SEC-023 | The system shall implement real-time security incident response with automated threat isolation | M2 |
| SEC-024 | The system shall perform continuous vulnerability assessment with automated patching for critical vulnerabilities | D |
| SEC-025 | The system shall implement zero-trust network architecture with microsegmentation | D |

### 7.2 Enhanced IDRA Compliance (MD FEEDBACK)

| ID      | Requirement Description                                                                                                                                                                        | Priority |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| SEC-011 | IDRA Monthly Reports: Generate Form IC-1 (Premium Collection) by 10th of each month with breakdown by product line, geographic region, partner channel in Excel format per IDRA template v2024 | M2        |
| SEC-012 | IDRA Monthly Reports: Generate Form IC-2 (Claims Intimation) by 10th of each month listing all new claims with policy number, claim amount, claim type, date of intimation                     | M2        |
| SEC-013 | IDRA Quarterly Reports: Generate Form IC-3 (Claims Settlement) within 15 days of quarter-end showing settlement ratio, average TAT, pending >30 days breakdown                                 | M2        |
| SEC-014 | IDRA Quarterly Reports: Generate Form IC-4 (Financial Performance) within 20 days of quarter-end with premium earned, claims paid, commission paid, net profit/loss                            | M3        |
| SEC-015 | IDRA Annual FCR: Generate Financial Condition Report (FCR) within 90 days of year-end including full CARAMELS framework assessment with external auditor sign-off                              | M3        |
| SEC-016 | IDRA Event-Based Reporting: Report significant incidents (fraud >BDT 1L, data breach, system outage >4hrs) within 48 hours via IDRA portal                                                     | M3        |

### 7.3 Enhanced AML/CFT Compliance (MD FEEDBACK)

| ID      | Requirement Description                                                                                                                                                                                                                                                                                                                                                                            | Priority |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| SEC-017 | AML/CFT Concrete Triggers: Flag transactions matching: (1) >3 policies in 7 days, (2) Premium >BDT 5L without income proof, (3) Nominee mismatch with no relationship doc, (4) Payment from third-party account, (5) Frequent cancellations >2 in 30 days, (6) Geographic anomaly (policy in Dhaka, payment from remote district), (7) Multiple failed KYC attempts >3, (8) PEP match in screening | M3        |
| SEC-018 | SAR Workflow: (1) System auto-flags suspicious transaction → (2) Compliance Officer reviews within 24hrs → (3) If confirmed suspicious, escalate to Business Admin+Focal Person → (4) Prepare SAR with evidence → (5) Submit to BFIU within 3 business days → (6) Mark account for enhanced monitoring → (7) Do NOT notify customer (tipping off prohibited)                                 | M3        |
| SEC-019 | Data Deletion Exceptions: Customer data deletion requests processed within 30 days EXCEPT: (a) Active policy holders (deletion after policy expiry+7yrs), (b) Ongoing claims (deletion after settlement+7yrs), (c) Under SAR investigation (deletion prohibited until case closed), (d) Regulatory hold (deletion requires IDRA/BFIU approval)                                                     | M3        |
| SEC-020 | Right to Erasure Workflow: Customer submits deletion request → System validates exceptions → If eligible, anonymize PII while retaining transaction records → Generate deletion certificate → Notify customer within 30 days                                                                                                                                                                   | D        |

### 7.4 Data Protection & Encryption Standards

| Data Classification                       | Encryption Standard                   | Key Management                     | Access Control                 |
| ----------------------------------------- | ------------------------------------- | ---------------------------------- | ------------------------------ |
| Personally Identifiable Information (PII) | AES-256                               | AWS KMS with 90-day rotation       | Role-based with audit logging  |
| Financial Transaction Data                | AES-256 + Additional Hashing          | TigerBeetle built-in encryption    | Restricted access with MFA     |
| KYC Documents                             | AES-256 with client-side encryption   | End-to-end encryption              | Compliance officer access only |
| Medical Records                           | AES-256 with additional anonymization | Healthcare-specific key management | Medical staff + consent-based  |
| Audit Logs                                | AES-256 with immutable storage        | Centralized key management         | Read-only access for auditors  |

---


### 7.5 Authentication & Authorization

**Multi-Factor Authentication (MFA):**
- SMS OTP for mobile number verification
- Email verification for account recovery
- Biometric authentication on supported mobile devices
- Hardware tokens for admin users

**Role-Based Access Control (RBAC):**
```
Roles Hierarchy:
├── System Admin
│   ├── Full system access
│   ├── User management
│   └── Security configuration
├── Business Admin
│   ├── Business operations
│   ├── Reporting access
│   └── Policy management
├── Partner Admin
│   ├── Agent management
│   ├── Commission tracking
│   └── Customer support
├── Agent
│   ├── Customer onboarding
│   ├── Policy sales
│   └── Basic support
└── Customer
    ├── Policy management
    ├── Claims submission
    └── Profile updates
```

**Session Management:**
- JWT tokens with 15-minute expiry
- Refresh token rotation
- Session invalidation on suspicious activity
- Device fingerprinting for fraud detection




### 7.6 IDRA Compliance Requirements

| IDRA ID | Requirement Description | Reporting Frequency | Priority | Owner |
|---------|------------------------|-------------------|----------|-------|
| IDRA-001 | Digital insurance product approval and registration | One-time + updates | M3 | Compliance Team |
| IDRA-002 | Customer data protection and privacy compliance | Quarterly review | M3 | Security Team |
| IDRA-003 | Policy issuance and documentation standards | Real-time compliance | M3 | Insurance Engine |
| IDRA-004 | Claims processing and settlement reporting | Monthly | M3 | Claims Team |
| IDRA-005 | Financial solvency and capital adequacy reporting | Quarterly | M3 | Finance Team |
| IDRA-006 | Agent licensing and training compliance | Ongoing | M3 | Partner Management |
| IDRA-007 | Marketing and sales practice compliance | Quarterly | M3 | Marketing Team |
| IDRA-008 | Actuarial and risk management reporting | Annual | D | Risk Management |
| IDRA-009 | Audit trail and record keeping requirements | Ongoing | M3 | Audit System |
| IDRA-010 | Regulatory change management and updates | As required | M3 | Compliance Team |

### 7.7 BFIU Anti-Money Laundering (AML) Compliance

| BFIU ID | Requirement Description | Threshold | Priority | Implementation |
|---------|------------------------|-----------|----------|----------------|
| BFIU-001 | Customer due diligence (CDD) for all policyholders | All customers | M3 | KYC verification system |
| BFIU-002 | Enhanced due diligence (EDD) for high-value policies | >50,000 BDT sum assured | M3 | Risk scoring system |
| BFIU-003 | Suspicious transaction monitoring and reporting | Real-time analysis | M3 | AI fraud detection |
| BFIU-004 | Cash transaction reporting | >10,000 BDT | M3 | Payment monitoring |
| BFIU-005 | Wire transfer monitoring | >100,000 BDT | M3 | Transaction screening |
| BFIU-006 | Politically exposed person (PEP) screening | All customers | M3 | PEP database integration |
| BFIU-007 | Sanctions list screening | All parties | M3| Sanctions database |
| BFIU-008 | Record retention for AML purposes | 5 years minimum | M3 | Data retention policies |
| BFIU-009 | AML training for employees and agents | Annual certification | M2 | Training management |
| BFIU-010 | AML audit and compliance reporting | Quarterly | M3| Compliance dashboard |

### 7.7.1 Customer Risk Scoring Matrix

**Risk Factors:**
1. **Transaction Frequency:** >3 claims in 6 months = +20 points
2. **Transaction Amount:** Single transaction >50K BDT = +15 points  
3. **Geographic Anomaly:** Claim location far from registered address = +10 points
4. **KYC Completeness:** Missing NID verification = +25 points
5. **Device Fingerprinting:** Multiple accounts from same device = +15 points
6. **Behavioral Anomaly:** Unusual activity patterns = +10 points

**Risk Categories:**
- **Low Risk:** 0-30 points → Annual review
- **Medium Risk:** 31-60 points → Semi-annual review  
- **High Risk:** >60 points → Quarterly review + Enhanced monitoring

### 7.7.2 Automated AML Monitoring Rules

**Transaction Monitoring Rules (20+ Rules):**

| Rule ID | Rule Description | Threshold | Action |
|---------|-----------------|-----------|--------|
| TM-001 | Structuring: Multiple transactions just below reporting threshold | 3+ transactions of 9K-10K BDT in 7 days | Flag for review |
| TM-002 | Rapid Movement: Quick policy purchase and claim | Claim within 7 days of purchase | Flag + manual review |
| TM-003 | Geographic Anomaly: Claim far from registered address | >100 km distance | Flag + location verification |
| TM-004 | Frequency Anomaly: Frequent claims | >3 claims in 6 months | Flag + pattern analysis |
| TM-005 | Amount Anomaly: Claim amount near coverage limit | >90% of coverage | Flag + document verification |
| TM-006 | Device Anomaly: Multiple accounts from same device | >3 accounts | Flag + fraud investigation |
| TM-007 | Payment Method Switch: Frequent payment method changes | >2 changes in 30 days | Flag + verification |
| TM-008 | Rapid Purchases: Multiple policies in short timeframe | >3 policies in 7 days | Flag + EDD |
| TM-009 | High-Value Premiums: Single premium exceeds threshold | >BDT 5 lakh | Enhanced due diligence |
| TM-010 | Frequent Cancellations: Policy cancellation patterns | >2 cancellations in 3 months | Flag + investigation |
| TM-011 | Mismatched Nominees: Nominee not family member | Non-relative nominee | Flag + verification |
| TM-012 | Payment Source Anomaly: Different payers for same policy | >2 different payers | Flag + source verification |
| TM-013 | Geographic Risk: High-risk geographic location | Blacklisted areas | Enhanced monitoring |
| TM-014 | Age Anomaly: Unusual age for product type | Outside typical range | Flag + verification |
| TM-015 | Occupation Risk: High-risk occupation categories | PEP, cash-intensive business | Enhanced due diligence |
| TM-016 | Document Inconsistency: Mismatched KYC documents | OCR verification failure | Flag + manual review |
| TM-017 | Refund Requests: Frequent refund requests | >2 refunds in 6 months | Flag + pattern analysis |
| TM-018 | Beneficiary Changes: Multiple beneficiary modifications | >2 changes in 12 months | Flag + verification |
| TM-019 | Third-Party Payments: Non-policyholder making payments | Different payer than insured | Flag + source verification |
| TM-020 | Dormant Activation: Long-dormant account suddenly active | No activity >6 months then sudden purchase | Flag + identity verification |

**Customer Risk Scoring:**

| Risk Factor | Points | Description |
|------------|--------|-------------|
| Transaction Frequency | +20 | >3 claims in 6 months |
| Transaction Amount | +15 | Single transaction >50K BDT |
| Geographic Anomaly | +10 | Claim location far from registered address |
| KYC Completeness | +25 | Missing NID verification |
| Device Fingerprinting | +15 | Multiple accounts from same device |
| Behavioral Anomaly | +10 | Unusual activity patterns |

**Risk Categories:**
- **Low Risk:** 0-30 points → Annual review
- **Medium Risk:** 31-60 points → Semi-annual review
- **High Risk:** >60 points → Quarterly review + Enhanced monitoring

**STR/SAR Filing Workflow:**
1. **Detection:** Automated rule triggers or manual reporting by staff
2. **Investigation:** Compliance Officer reviews flagged activity within 24 hours
3. **Decision:** Determine if suspicious (consult with Business Admin if needed)
4. **Filing:** Submit STR/SAR to BFIU portal within 7 days
5. **Action:** Freeze account if necessary, notify authorities
6. **Documentation:** Maintain records for 7 years
7. **No Tipping Off:** Customer must not be notified per law

| Rule ID | Rule Description | Threshold | Action |
|---------|-----------------|-----------|--------|
| TM-001 | Structuring: Multiple transactions just below reporting threshold | 3+ transactions of 9K-10K BDT in 7 days | Flag for review |
| TM-002 | Rapid Movement: Quick policy purchase and claim | Claim within 7 days of purchase | Flag + manual review |
| TM-003 | Geographic Anomaly: Claim far from registered address | >100 km distance | Flag + location verification |
| TM-004 | Frequency Anomaly: Frequent claims | >3 claims in 6 months | Flag + pattern analysis |
| TM-005 | Amount Anomaly: Claim amount near coverage limit | >90% of coverage | Flag + document verification |
| TM-006 | Device Anomaly: Multiple accounts from same device | >3 accounts | Flag + fraud investigation |
| TM-007 | Payment Method Switch: Frequent payment method changes | >2 changes in 30 days | Flag + verification |
| TM-008 | Time Anomaly: Transactions outside normal hours | 11 PM - 6 AM transactions | Flag + review |
| TM-009 | Velocity Check: High transaction volume | >10 transactions per day | Flag + velocity analysis |
| TM-010 | Round Amount: Suspicious round number patterns | Multiple round amounts (10K, 20K, 50K) | Flag + pattern review |

### 7.7.3 STR/SAR Filing Workflow

1. **Detection:** Automated rule triggers or manual reporting by staff
2. **Investigation:** Compliance Officer reviews flagged activity within 24 hours  
3. **Decision:** Determine if suspicious (consult with Business Admin if needed)
4. **Filing:** Submit STR/SAR to BFIU portal within 7 days
5. **Action:** Freeze account if necessary, notify authorities
6. **Documentation:** Maintain records for 7 years
7. **No Tipping Off:** Customer must not be notified per law

### 7.8. AML/CFT Compliance Requirements

#### 7.8.1 Customer Due Diligence (CDD) Framework

**Mandatory CDD Requirements for Bangladesh:**

| Requirement                     | Implementation                              | Compliance Standard  |
| ------------------------------- | ------------------------------------------- | -------------------- |
| **Identity Verification** | NID/Passport verification via approved eKYC | BFIU Guidelines      |
| **Address Verification**  | Utility bill or bank statement              | MLPA Requirements    |
| **Photo Identification**  | Selfie with liveness detection              | Enhanced CDD         |
| **Source of Funds**       | Income declaration for high-value policies  | Risk-based approach  |
| **PEP Screening**         | Automated screening against watchlists      | FATF Recommendations |

#### 7.8.2 Risk-Based Customer Categorization

| Risk Level            | Criteria                                  | CDD Requirements             | Monitoring Frequency |
| --------------------- | ----------------------------------------- | ---------------------------- | -------------------- |
| **Low Risk**    | Standard customers, low premium policies  | Standard CDD                 | Annual review        |
| **Medium Risk** | Higher premiums, multiple policies        | Enhanced documentation       | Quarterly review     |
| **High Risk**   | PEPs, large premiums, suspicious patterns | Enhanced Due Diligence (EDD) | Monthly monitoring   |
| **Prohibited**  | Sanctioned individuals, blocked entities  | Transaction rejection        | Real-time blocking   |

#### 7.8.3 Automated AML Monitoring Rules

| Monitoring Rule                        | Threshold                                         | Alert Level | Action Required          |
| -------------------------------------- | ------------------------------------------------- | ----------- | ------------------------ |
| **Rapid Policy Purchases**       | >3 policies in 7 days                             | High        | Enhanced verification    |
| **High-Value Premiums**          | >BDT 5 lakh                                       | High        | Management approval      |
| **Frequent Cancellations**       | >2 cancellations in 30 days                       | Medium      | Pattern analysis         |
| **Mismatched Nominees**          | Different family names without relationship proof | Medium      | Additional documentation |
| **Geographic Anomalies**         | Transaction from unusual location                 | Low         | Location verification    |
| **Payment Method Inconsistency** | Different mobile numbers vs NID                   | Medium      | Customer verification    |

#### 7.8.4 Record Keeping & Audit Trail

| Document Type                 | Retention Period                | Storage Requirements         | Access Controls               |
| ----------------------------- | ------------------------------- | ---------------------------- | ----------------------------- |
| **CDD Documentation**   | 5+ years after relationship end | Encrypted PostgreSQL + S3    | Compliance team only          |
| **Transaction Records** | 7+ years                        | TigerBeetle + Archive        | Audit and compliance          |
| **STR Documentation**   | 10+ years                       | Secured offline storage      | Business Admin + Focal Person |
| **Training Records**    | 5+ years                        | HR system integration        | HR and compliance             |
| **System Audit Logs**   | 20+ years                       | Immutable PostgreSQL logging | System administrators         |

### 7.9 Data Protection & Privacy

**Privacy by Design Implementation:**
- Data minimization in collection and processing
- Purpose limitation for data usage
- Storage limitation with automated purging
- Accuracy maintenance with user control
- Security safeguards at all layers
- Transparency through privacy notices
- User control over personal data

**Data Subject Rights (GDPR-Style):**
- Right to access personal data
- Right to rectification of inaccurate data
- Right to erasure ("right to be forgotten")
- Right to restrict processing
- Right to data portability
- Right to object to processing
- Rights related to automated decision making

### 7.10 Cybersecurity Measures

**Threat Protection:**
```
Defense Layers:
├── Network Security
│   ├── WAF (Web Application Firewall)
│   ├── DDoS protection
│   └── Network segmentation
├── Application Security
│   ├── OWASP Top 10 protection
│   ├── Input validation
│   └── SQL injection prevention
├── Data Security
│   ├── Encryption at rest/transit
│   ├── Key management
│   └── Data masking
└── Monitoring & Response
    ├── SIEM integration
    ├── SOC operations
    └── Incident response
```

**Security Monitoring:**
- 24/7 Security Operations Center (SOC)
- Real-time threat intelligence feeds
- Automated incident response workflows
- Vulnerability management program
- Regular penetration testing

### 7.11 Audit & Logging Requirements

**Audit Trail Requirements:**
- All user actions logged with timestamps
- Immutable audit records using blockchain/cryptographic hashing
- Real-time audit log streaming to SIEM
- Audit log retention for 7 years (regulatory requirement)
- Automated anomaly detection on audit patterns

**Compliance Monitoring:**
- Automated compliance rule checking
- Real-time policy violation alerts
- Regulatory reporting automation
- Compliance dashboard for management oversight
- Third-party audit support and evidence collection

[[[PAGEBREAK]]]
