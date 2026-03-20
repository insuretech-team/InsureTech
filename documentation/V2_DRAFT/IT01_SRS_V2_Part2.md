## 3. System Features & Functional Requirements

Each requirement below has a unique ID (e.g., FR-1). Priority levels: M = Mandatory (Phase 1), D = Desirable (Phase 1), F = Future (Phase 2/3).

### 3.1 Authentication & User Management

**FR-1 (M)** — **Account Registration**: Support phone-based registration with OTP validation (SMS) and capture minimal profile fields (name, DOB, phone, email optional).

- **Bangladesh-specific**: NID validation with format checks (10/13/17 digit patterns)
- **Duplicate Prevention**: Cross-check against existing NID/phone combinations
- **Fallback**: Support for passport numbers for non-residents

**FR-2 (M)** — **Login**: OTP login and password-based login; session management and refresh tokens.

- **Security**: JWT tokens with 24-hour expiry
- **Session Management**: Single sign-on across web and mobile
- **Rate Limiting**: Max 5 OTP requests per hour per number

**FR-3 (M)** — **Profile Management**: Update personal info, nominee details, and document uploads.

- **Change Control**: Audit trail for all profile changes
- **Document Management**: Version control for uploaded documents
- **Nominee Management**: Support multiple nominees with percentage allocation

**FR-4 (M)** — **Duplicate Prevention**: Block duplicate accounts by national ID/phone; provide merge or support flow.

- **Detection Logic**: Check NID, phone, and biometric hash (if available)
- **Resolution Workflow**: Manual review queue for suspected duplicates
- **Customer Support**: Guided merge process with verification

### 3.2 Digital KYC & Document Verification

**FR-5 (M)** — **Document Upload**: Capture images/PDFs for NID/passport, photos, medical docs. UI flows per screens.

- **File Formats**: JPEG, PNG, PDF (max 10MB per file)
- **Quality Checks**: Blur detection, resolution validation
- **Security**: Immediate encryption and secure storage

**FR-6 (M)** — **OCR & Validation**: Extract NID/passport fields via OCR and pre-validate (format + checksum).

- **OCR Engine**: Support for Bengali and English text
- **Validation Rules**: NID checksum algorithm, date format validation
- **Confidence Scoring**: Auto-approval for high confidence (>90%), manual review for medium (70-90%)

**FR-7 (D)** — **eKYC Integration**: Integrate with third-party eKYC services for automated verification (Phase 1 if available).

- **Approved Providers**: Bangladesh Bank approved eKYC providers only
- **Fallback Process**: Manual KYC if eKYC unavailable
- **Compliance**: Full audit trail for regulatory inspection

 **FR-8 (M)** — **Liveness Detection**: Implement selfie with liveness detection to prevent document fraud.

- **Technology**: 3D face mapping or eye movement detection
- **Storage**: Biometric templates, not raw images
- **Privacy**: Option to delete biometric data after verification

### 3.3 Product Catalog & Policy Discovery

**FR-9 (M)** — **Product Catalog API**: Serve product list with metadata (name, coverage, premiums, insurer, T&Cs).

- **Dynamic Pricing**: Real-time premium calculation based on user profile
- **Localization**: Bengali product descriptions and terms
- **Insurer Integration**: Real-time product availability from multiple insurers

**FR-10 (M)** — **Filter & Compare**: Support filtering (category, premium, coverage) and side-by-side comparison (up to 3).

- **Filter Options**: Premium range, coverage amount, deductible, insurer rating
- **Comparison Matrix**: Key features highlighted in tabular format
- **Educational Content**: Tooltips explaining insurance terminology

**FR-11 (M)** — **Product Detail**: Full policy wording, exclusions, and illustrative premium breakdown.

- **Regulatory**: Complete IDRA-approved policy wordings
- **Transparency**: Clear exclusions and waiting periods
- **Calculator**: Interactive premium calculator with scenario modeling

### 3.4 Policy Purchase & Issuance

**FR-12 (M)** — **Multi-step Purchase Flow**: Personal details, nominee, document upload, review, payment, confirmation. Flow follows provided mockups.

- **Progress Tracking**: Visual progress indicator (5-step process)
- **Save & Resume**: Allow users to complete purchase later
- **Mobile Optimization**: Touch-friendly forms with smart defaults

**FR-13 (M)** — **Premium Calculation**: Request insurer API or local pricing engine and show breakdown.

- **Calculation Rules**:
  - Base premium + age factor + health loading + tax
  - Discount for bulk purchases or partner channels
  - Currency: BDT with proper formatting (commas, decimals)
- **Fallback Strategy**: Local pricing engine if insurer API unavailable
- **Audit Trail**: Log all premium calculations for regulatory review

**FR-14 (M)** — **Payment Integration**: Support bKash, Nagad, card (PCI-DSS compliant), bank transfer.

- **Payment Methods**:
  - bKash: Personal and merchant accounts
  - Nagad: Wallet and bank account funding
  - Cards: Visa, Mastercard, local debit cards
  - Bank Transfer: BEFTN integration for bulk payments
- **Security**: PCI-DSS Level 1 compliance
- **Reconciliation**: Daily settlement files and automated matching

**FR-15 (M)** — **Policy Issuance**: On successful payment, generate a digital policy certificate (PDF) and store in user account.

- **Document Generation**: PDF with QR code for verification
- **Digital Signature**: Insurer's digital signature on policy documents
- **Delivery**: SMS with download link + email attachment
- **Backup Storage**: Immutable storage for 7+ years

**FR-16 (D)** — **Promo/Discounts**: Support coupon codes and partner-discount flows.

- **Coupon System**: Time-limited, usage-limited promotional codes
- **Partner Discounts**: Automatic discounts for certain channels
- **Validation**: Real-time coupon validation and abuse prevention

### 3.5 Claims Management

**FR-17 (M)** — **Claim Initiation**: Policy prefill, claim reason selection, document upload (images, bills), and submission.

- **Smart Prefill**: Auto-populate policy and customer details
- **Document Requirements**: Dynamic document checklist based on claim type
- **Geo-location**: Optional location capture for claim verification

**FR-18 (M)** — **Status Tracking**: Provide claim status updates and admin notes to user.

- **Status Workflow**:
  - Submitted → Acknowledged (24h)
  - Under Review → Document Review (3-5 days)
  - Approved/Rejected → Payment/Appeal (1-2 days)
- **Communication**: SMS/email notifications at each status change
- **Transparency**: Reason codes for rejections with appeal process

**FR-19 (M)** — **Admin Workflow**: Claims dashboard for triage, verification, approval, rejection, and payment initiation.

- **Triage Rules**:
  - Auto-approve: Claims <BDT 10,000 with complete documentation
  - Manual Review: Claims >BDT 50,000 or flagged by fraud detection
  - Escalation: Claims >BDT 200,000 require manager approval
- **SLA Monitoring**: Automated escalation for overdue claims
- **Audit Trail**: Complete history of claim processing decisions

**FR-20 (D)** — **Automated Triage**: OCR, image verification and rule-based auto-accept for small claims.

- **OCR Integration**: Extract data from medical bills, receipts
- **Fraud Detection**: Cross-check against known fraud patterns
- **Auto-Settlement**: Immediate payout for pre-approved claim types

 **FR-21 (M)** — **Zero Human Touch Claims (ZHCT)**: Fully automated processing for claims under BDT 10,000.

- **Eligibility**: Active policy, complete documentation, no previous fraud flags
- **Processing Time**: 30 minutes maximum for eligible claims
- **Quality Control**: Random audit of 5% of auto-processed claims

### 3.6 Policy Management & Renewals

**FR-22 (M)** — **Policy Dashboard**: Active & past policies, download documents, renewal prompts.

- **Dashboard Features**: Policy status, premium due dates, claim history
- **Document Access**: Download policy certificates, receipts, claim forms
- **Renewal Alerts**: 60, 30, 15, 7 days before expiry

**FR-23 (M)** — **Renewals**: Auto-renew option and manual renew flows, with reminders.

- **Auto-renewal Logic**:
  - Default opt-in for health/life policies
  - Grace period: 30 days with continued coverage
  - Premium adjustment: Annual rate review with customer notification
- **Manual Renewal**: Option to modify coverage or switch insurers
- **Payment**: Use stored payment method or request new payment

**FR-24 (D)** — **Partial Adjustments**: Allow address, nominee change requests.

- **Change Types**: Address, nominee details, beneficiary percentages
- **Approval Process**: Instant for address, verification required for nominees
- **Documentation**: Updated policy certificates issued automatically

### 3.7 Notifications & Communication

**FR-25 (M)** — **Notification Engine**: Trigger SMS/Push/Email for OTP, purchase confirmation, claims updates, renewal reminders.

- **Multi-channel**: SMS (primary), push notifications, email (secondary)
- **Localization**: Bengali and English message templates
- **Delivery Confirmation**: Track delivery status and retry failed messages

**FR-26 (D)** — **Marketing Opt-in/Opt-out**: Manage user preferences.

- **Granular Control**: Separate preferences for transactional vs promotional messages
- **Compliance**: PDPA-ready consent management
- **Segmentation**: Targeted campaigns based on user behavior and preferences

### 3.8 Admin & Reporting

**FR-27 (M)** — **Admin Portal**: User management, product management, claim dashboards, and role-based access control.

- **Role-Based Access**:
  - Super Admin: Full system access
  - Claims Officer: Claims processing only
  - Product Manager: Product and pricing management
  - Customer Service: Read-only customer data with limited edit rights
- **Audit Logging**: All administrative actions logged with user identification

**FR-28 (M)** — **Reports**: Daily sales, claims ratio, partner performance, policy counts and KPIs (aligned to business plan targets).

- **Standard Reports**:
  - Daily Sales Report: Policies sold, premium collected, channel breakdown
  - Claims Analysis: Claims ratio, average settlement time, fraud detection stats
  - Regulatory Reports: IDRA quarterly submissions, AML suspicious activity
- **Custom Reports**: Ad-hoc query builder for business intelligence
- **Automated Distribution**: Scheduled email delivery to stakeholders

### 3.9 Partner / Agent Portal

**FR-29 (M)** — **Embedded Flow**: Partner can initiate policy purchase for end-customer via API or partner portal.

- **White-label Integration**: Partner-branded experience with configurable UI
- **Commission Tracking**: Real-time commission calculations and statements
- **Lead Management**: Track customer referrals and conversion rates

**FR-30 (M)** — **Partner Dashboard**: Commission statements, leads, onboarding analytics.

- **Performance Metrics**: Conversion rates, customer satisfaction scores
- **Training Resources**: Product guides, sales scripts, compliance training
- **Support Tools**: Direct chat with partner success team

### 3.10 Audit & Logging

**FR-31 (M)** — **Audit Trail**: Immutable logs for critical actions (policy issue, claim approval, payment).

- **Immutable Storage**: Blockchain or write-once storage for audit logs
- **Log Retention**: 7 years minimum for regulatory compliance
- **Access Control**: Audit logs accessible only to compliance team and auditors

**FR-32 (M)** — **Data Retention**: Maintain records to meet regulatory retention periods (configurable).

- **Retention Schedule**:
  - Customer data: 7 years after relationship ends
  - Policy data: 10 years after policy expiry
  - Claims data: 15 years for life insurance, 7 years for general insurance
  - Audit logs: 10 years minimum
- **Automated Purging**: Scheduled deletion of expired data with audit trail
