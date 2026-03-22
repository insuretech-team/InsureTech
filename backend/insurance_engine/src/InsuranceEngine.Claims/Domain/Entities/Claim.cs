using System;
using System.Collections.Generic;
using System.Linq;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Claims.Domain.Entities;

public class Claim
{
    public Guid Id { get; set; }
    public string ClaimNumber { get; set; } = string.Empty;
    public Guid PolicyId { get; set; }
    public Guid CustomerId { get; set; }
    public ClaimStatus Status { get; set; }
    public ClaimType Type { get; set; }

    // Money fields (paisa)
    public long ClaimedAmount { get; set; }
    public string ClaimedCurrency { get; set; } = "BDT";
    public long ApprovedAmount { get; set; }
    public string ApprovedCurrency { get; set; } = "BDT";
    public long SettledAmount { get; set; }
    public string SettledCurrency { get; set; } = "BDT";

    public DateTime IncidentDate { get; set; }
    public string IncidentDescription { get; set; } = string.Empty;
    public string? PlaceOfIncident { get; set; }

    public DateTime SubmittedAt { get; set; }
    public DateTime? ApprovedAt { get; set; }
    public DateTime? SettledAt { get; set; }
    public string? RejectionReason { get; set; }

    public ClaimProcessingType ProcessingType { get; set; }

    // --- Proto-aligned fields (FR-100, FR-101) ---

    /// <summary>
    /// Deductible amount in paisa (BDT minor units). Proto: deductible_amount
    /// </summary>
    public long DeductibleAmount { get; set; }
    public string DeductibleCurrency { get; set; } = "BDT";

    /// <summary>
    /// Co-payment amount in paisa (BDT minor units). Proto: co_pay_amount
    /// </summary>
    public long CoPayAmount { get; set; }
    public string CoPayCurrency { get; set; } = "BDT";

    /// <summary>
    /// Encrypted bank details or ref to linked bank for payout. Proto: bank_details_for_payout
    /// PII — must be encrypted at rest (AES-256).
    /// </summary>
    public string? BankDetailsForPayout { get; set; }

    /// <summary>
    /// Whether the claimant can appeal a rejected claim. Proto: appeal_option_available
    /// </summary>
    public bool AppealOptionAvailable { get; set; }

    /// <summary>
    /// JSON array of in-app messages related to this claim. Proto: in_app_messages (JSONB)
    /// </summary>
    public string? InAppMessages { get; set; }

    /// <summary>
    /// Internal notes from claim processor. Proto: processor_notes
    /// </summary>
    public string? ProcessorNotes { get; set; }

    // --- Fraud check reference ---
    public Guid? FraudCheckId { get; set; }
    public FraudCheckResult? FraudCheck { get; set; }

    // --- Navigation properties ---
    public List<ClaimApproval> Approvals { get; set; } = new();
    public List<ClaimDocument> Documents { get; set; } = new();

    // Audit
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }
    public bool IsDeleted { get; set; }

    // Money convenience accessors
    public Money ClaimedMoney => new(ClaimedAmount, ClaimedCurrency);
    public Money ApprovedMoney => new(ApprovedAmount, ApprovedCurrency);
    public Money SettledMoney => new(SettledAmount, SettledCurrency);
    public Money DeductibleMoney => new(DeductibleAmount, DeductibleCurrency);
    public Money CoPayMoney => new(CoPayAmount, CoPayCurrency);

    // Constants for approval matrix based on BDT amounts (converted to paisa)
    private const long L1_THRESHOLD = 20_000_000;      // 200,000 BDT
    private const long L2_THRESHOLD = 50_000_000;      // 500,000 BDT
    private const long L3_THRESHOLD = 100_000_000;     // 1,000,000 BDT

    public int GetRequiredApprovalLevel()
    {
        if (ClaimedAmount <= L1_THRESHOLD) return 1;  // L1 (Officer/Auto)
        if (ClaimedAmount <= L2_THRESHOLD) return 2;  // L2 (Manager)
        if (ClaimedAmount <= L3_THRESHOLD) return 3;  // L3 (Head/Joint)
        return 4;                                     // L4 (Board/Insurer)
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
