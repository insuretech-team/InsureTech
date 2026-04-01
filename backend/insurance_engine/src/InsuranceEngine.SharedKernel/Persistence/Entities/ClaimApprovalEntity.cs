using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'claim_approvals' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// Multi-level approval workflow (L1, L2, L3, Board).
/// </summary>
[Table("claim_approvals", Schema = "insurance_schema")]
public class ClaimApprovalEntity
{
    [Key]
    [Column("approval_id")]
    public Guid ApprovalId { get; set; }

    [Column("claim_id")]
    public Guid ClaimId { get; set; }

    [Column("approver_id")]
    public Guid ApproverId { get; set; }

    [Column("approver_role")]
    public string ApproverRole { get; set; } = string.Empty;

    [Column("approval_level")]
    public int ApprovalLevel { get; set; } // 1-4 (L1, L2, L3, Board)

    [Column("decision")]
    public string Decision { get; set; } = "PENDING";

    [Column("approved_amount")]
    public long? ApprovedAmount { get; set; }

    [Column("approved_currency")]
    public string ApprovedCurrency { get; set; } = "BDT";

    [Column("notes")]
    public string? Notes { get; set; }

    [Column("approved_at")]
    public DateTime? ApprovedAt { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    // Navigation
    public ClaimEntity Claim { get; set; } = null!;
}
