# Claims Processing Workflow Example

## Claim Submission Flow

```
Step 1: Claim Initiation
- User selects active policy
- Clicks "File Claim"
- System validates policy status (must be ACTIVE)

Step 2: Incident Details
- Incident Date: 2025-01-10
- Incident Type: Hospitalization
- Hospital Name: LabAid Hospital
- Claimed Amount: 25,000 BDT

Step 3: Document Upload
- Hospital Bill (PDF)
- Prescription (Image)
- Discharge Summary (PDF)
- System validates:
  * File size < 5MB
  * Image quality check
  * OCR extraction of bill details

Step 4: Claim Submission
- Claim Number Generated: CLM-2025-0001-000045
- Digital hash created (SHA-256)
- Notification sent to partner/insurer
- Status: SUBMITTED

Step 5: Review Process
- Auto-assigned to Claims Officer (Amount < 10K)
- OR Assigned to Claims Manager (Amount 10K-50K)
- Fraud check runs automatically

Step 6: Approval
- Approver reviews documents
- Approver adds notes
- Decision: APPROVED for 23,000 BDT (excluded non-covered items)

Step 7: Settlement
- Payment initiated via bKash
- Settlement time: 24 hours
- Customer notified
- Status: SETTLED
```

## Claims Approval Matrix

| Claimed Amount | Approval Level | Approver | TAT |
|----------------|----------------|----------|-----|
| 0-10K | L1 Auto/Officer | System OR Claims Officer | 24 Hours |
| 10K-50K | L2 Manager | Claims Manager | 3 Days |
| 50K-2L | L3 Head | Business Admin + Focal Person | 7 Days |
| 2L+ | Board | Board + Insurer Approval | 15 Days |
