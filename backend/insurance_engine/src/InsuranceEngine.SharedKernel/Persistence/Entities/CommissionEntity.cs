using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'commissions' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("commissions", Schema = "insurance_schema")]
public class CommissionEntity
{
    [Key]
    [Column("commission_id")]
    public Guid CommissionId { get; set; }

    [Column("commission_number")]
    public string CommissionNumber { get; set; } = $"COM-{Guid.NewGuid().ToString()[..8].ToUpper()}";

    [Column("policy_id")]
    public Guid PolicyId { get; set; }

    [Column("partner_id")]
    public Guid? PartnerId { get; set; }

    [Column("agent_id")]
    public Guid? AgentId { get; set; }

    [Column("recipient_id")]
    public Guid? RecipientId { get; set; }

    [Column("type")]
    public string CommissionType { get; set; } = "ACQUISITION"; // ACQUISITION, RENEWAL, CLAIMS_ASSISTANCE

    [Column("commission_amount")]
    public long CommissionAmount { get; set; }

    [Column("commission_currency")]
    public string CommissionCurrency { get; set; } = "BDT";

    [Column("commission_rate")]
    public decimal CommissionRate { get; set; }

    [Column("status")]
    public string Status { get; set; } = "PENDING"; // PENDING, PROCESSING, PAID, CANCELLED

    [Column("payment_id")]
    public Guid? PayoutId { get; set; }

    [Column("calculation_breakdown")]
    public string? CalculationBreakdown { get; set; } // JSONB

    [Column("paid_at")]
    public DateTime? PaidAt { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [Column("deleted_at")]
    public DateTime? DeletedAt { get; set; }

    // Navigation
    public PolicyEntity Policy { get; set; } = null!;
}
