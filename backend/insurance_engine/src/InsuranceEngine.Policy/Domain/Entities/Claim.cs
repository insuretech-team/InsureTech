using System;
using System.Collections.Generic;
using System.Linq;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;

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
    public double FraudScore { get; set; }
    public string? FraudCheckData { get; set; } // JSONB
    
    public List<ClaimApproval> Approvals { get; set; } = new();
    public List<ClaimDocument> Documents { get; set; } = new();
    
    // Audit
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public bool IsDeleted { get; set; }

    // Constants for approval matrix
    private const long ZHTC_THRESHOLD = 1_000_000;      // 10,000 BDT
    private const long L1_THRESHOLD = 5_000_000;        // 50,000 BDT
    private const long L2_THRESHOLD = 20_000_000;       // 200,000 BDT
    private const long L3_THRESHOLD = 50_000_000;       // 500,000 BDT

    public int GetRequiredApprovalLevel()
    {
        if (ClaimedAmount <= ZHTC_THRESHOLD) return 0;
        if (ClaimedAmount <= L1_THRESHOLD) return 1;
        if (ClaimedAmount <= L2_THRESHOLD) return 2;
        if (ClaimedAmount <= L3_THRESHOLD) return 3;
        return 4;
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
        else if (decision == ApprovalDecision.Escalated)
        {
            Status = ClaimStatus.UnderReview;
        }

        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }
}
