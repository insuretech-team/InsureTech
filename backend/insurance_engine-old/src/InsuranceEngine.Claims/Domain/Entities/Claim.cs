using System;
using System.Collections.Generic;
using System.Linq;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.Claims.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Claims.Domain.Entities;

public class Claim : AggregateRoot<Guid>
{
    public string ClaimNumber { get; private set; } = string.Empty;
    public Guid PolicyId { get; private set; }
    public Guid CustomerId { get; private set; }
    public ClaimStatus Status { get; private set; }
    public ClaimType Type { get; private set; }

    // Money fields (paisa)
    public long ClaimedAmount { get; private set; }
    public string ClaimedCurrency { get; private set; } = "BDT";
    public long ApprovedAmount { get; private set; }
    public string ApprovedCurrency { get; private set; } = "BDT";
    public long SettledAmount { get; private set; }
    public string SettledCurrency { get; private set; } = "BDT";

    public DateTime IncidentDate { get; private set; }
    public string IncidentDescription { get; private set; } = string.Empty;
    public string? PlaceOfIncident { get; private set; }

    public DateTime SubmittedAt { get; private set; }
    public DateTime? ApprovedAt { get; private set; }
    public DateTime? SettledAt { get; private set; }
    public string? RejectionReason { get; private set; }

    public ClaimProcessingType ProcessingType { get; private set; }

    public long DeductibleAmount { get; private set; }
    public string DeductibleCurrency { get; private set; } = "BDT";
    public double CoPayPercentage { get; private set; }
    public long CoPayAmount { get; private set; }
    public string CoPayCurrency { get; private set; } = "BDT";

    public string? BankDetailsForPayout { get; set; }
    public bool AppealOptionAvailable { get; set; }
    public string? InAppMessages { get; set; }
    public string? ProcessorNotes { get; set; }

    // --- Fraud check reference ---
    public Guid? FraudCheckId { get; private set; }
    public FraudCheckResult? FraudCheck { get; private set; }

    // --- Navigation properties ---
    public List<ClaimApproval> Approvals { get; private set; } = new();
    public List<ClaimDocument> Documents { get; private set; } = new();

    // Audit
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }
    public string? AuditInfoJson { get; set; }
    public DateTime? DeletedAt { get; private set; }

    // Constants for approval matrix based on SRS v3.11 (Appendix C) & PoliSync
    private const long ZHTC_THRESHOLD = 1_000_000;    // 10,000 BDT
    private const long L1_THRESHOLD = 5_000_000;      // 50,000 BDT
    private const long L2_THRESHOLD = 20_000_000;     // 200,000 BDT
    private const long L3_THRESHOLD = 50_000_000;     // 500,000 BDT
    
    private const double ZHTC_FRAUD_THRESHOLD = 0.30;

    // --- Money convenience accessors ---
    public Money ClaimedMoney => new(ClaimedAmount, ClaimedCurrency);
    public Money ApprovedMoney => new(ApprovedAmount, ApprovedCurrency);
    public Money SettledMoney => new(SettledAmount, SettledCurrency);
    public Money DeductibleMoney => new(DeductibleAmount, DeductibleCurrency);

    // EF Core constructor
    public Claim() { }

    public static Claim File(
        string claimNumber,
        Guid policyId,
        Guid customerId,
        ClaimType type,
        long amount,
        DateTime incidentDate,
        string incidentDescription,
        string? placeOfIncident)
    {
        var claim = new Claim
        {
            Id = Guid.NewGuid(),
            ClaimNumber = claimNumber,
            PolicyId = policyId,
            CustomerId = customerId,
            Type = type,
            ClaimedAmount = amount,
            IncidentDate = incidentDate,
            IncidentDescription = incidentDescription,
            PlaceOfIncident = placeOfIncident,
            Status = ClaimStatus.Submitted,
            ProcessingType = ClaimProcessingType.Manual,
            SubmittedAt = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        claim.AddDomainEvent(new ClaimSubmittedEvent(
            claim.Id, claim.ClaimNumber, claim.PolicyId, claim.CustomerId, 
            claim.ClaimedAmount, claim.ClaimedCurrency, claim.IncidentDate));

        return claim;
    }

    public void ApplyFraudCheck(Guid fraudCheckId, double fraudScore, List<string> flags)
    {
        FraudCheckId = fraudCheckId;
        FraudCheck = new FraudCheckResult 
        { 
            Id = fraudCheckId, 
            ClaimId = Id, 
            FraudScore = fraudScore, 
            CreatedAt = DateTime.UtcNow 
        };

        // ZHTC (Zero Hassle Trust Claim) Logic - FR-093
        if (ClaimedAmount <= ZHTC_THRESHOLD && fraudScore < ZHTC_FRAUD_THRESHOLD)
        {
            ProcessingType = ClaimProcessingType.AutoAdjudicated;
            Status = ClaimStatus.Approved;
            ApprovedAmount = ClaimedAmount; // Simplified for ZHTC; financials applied later if needed
            ApprovedAt = DateTime.UtcNow;
        }

        UpdatedAt = DateTime.UtcNow;
    }

    public void CalculateFinancials(long deductible, double coPayPercentage)
    {
        DeductibleAmount = deductible;
        CoPayPercentage = coPayPercentage;

        // Formula: Approved = (Claimed - Deductible) * (1 - CoPay%)
        // or as per SRS FR-100 literal interpretation: (Claim - Deductible) * CoPay%
        // Given co-pay is usually what the insurer pays in this context (Trust claims), 
        // we'll follow the logic where ApprovedAmount is the final payout.
        
        var netAmount = Math.Max(0, ClaimedAmount - DeductibleAmount);
        ApprovedAmount = (long)(netAmount * (1.0 - CoPayPercentage));
        
        UpdatedAt = DateTime.UtcNow;
    }

    public void AddDocument(string type, string url, string hash)
    {
        var doc = new ClaimDocument
        {
            Id = Guid.NewGuid(),
            ClaimId = Id,
            DocumentType = type,
            FileUrl = url,
            FileHash = hash,
            UploadedAt = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
        Documents.Add(doc);
        UpdatedAt = DateTime.UtcNow;
    }

    public void UpdateStatus(ClaimStatus status)
    {
        Status = status;
        UpdatedAt = DateTime.UtcNow;
    }

    public void Settle(long amount, string currency)
    {
        SettledAmount = amount;
        SettledCurrency = currency;
        SettledAt = DateTime.UtcNow;
        Status = ClaimStatus.Settled;
        UpdatedAt = DateTime.UtcNow;
    }

    public int GetRequiredApprovalLevel()
    {
        if (ClaimedAmount <= ZHTC_THRESHOLD) return 0; // Level 0 (Auto/Officer ZHTC)
        if (ClaimedAmount <= L1_THRESHOLD) return 1;   // Level 1 (Officer)
        if (ClaimedAmount <= L2_THRESHOLD) return 2;   // Level 2 (Manager)
        if (ClaimedAmount <= L3_THRESHOLD) return 3;   // Level 3 (Director/Head)
        return 4;                                      // Level 4 (Board/CEO)
    }

    public Result AddApproval(Guid approverId, string role, int level, ApprovalDecision decision, long approvedAmount, string notes)
    {
        var approval = new ClaimApproval
        {
            Id = Guid.NewGuid(),
            ClaimId = Id,
            ApproverId = approverId,
            ApproverRole = role,
            ApprovalLevel = level,
            Decision = decision,
            ApprovedAmount = approvedAmount,
            Notes = notes,
            ApprovedAt = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow
        };

        Approvals.Add(approval);

        if (decision == ApprovalDecision.Approved)
        {
            ApprovedAmount = approvedAmount;
            var requiredLevel = GetRequiredApprovalLevel();
            if (level >= requiredLevel)
            {
                Status = ClaimStatus.Approved;
                ApprovedAt = DateTime.UtcNow;
            }
            else
            {
                Status = ClaimStatus.UnderReview;
            }
        }
        else if (decision == ApprovalDecision.Rejected)
        {
            Status = ClaimStatus.Rejected;
            RejectionReason = notes;
        }
        else if (decision == ApprovalDecision.Escalated || decision == ApprovalDecision.NeedsMoreInfo)
        {
            Status = decision == ApprovalDecision.NeedsMoreInfo
                ? ClaimStatus.PendingDocuments
                : ClaimStatus.UnderReview;
        }

        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }
}
