namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'claim_approvals' table in insurance_schema.
/// Aligned with insuretech.claims.entity.v1.ClaimApproval proto definition.
/// Multi-level approval workflow (L1, L2, L3, Board).
/// </summary>
public class ClaimApprovalEntity
{
    public Guid ApprovalId { get; set; }
    public Guid ClaimId { get; set; }
    public Guid ApproverId { get; set; }
    public string ApproverRole { get; set; } = string.Empty;
    public int ApprovalLevel { get; set; } // 1-4 (L1, L2, L3, Board)
    public string Decision { get; set; } = "PENDING";
    public long? ApprovedAmount { get; set; }
    public string ApprovedCurrency { get; set; } = "BDT";
    public string? Notes { get; set; }
    public DateTime? ApprovedAt { get; set; }
    public DateTime CreatedAt { get; set; }

    // Navigation
    public ClaimEntity Claim { get; set; } = null!;
}
