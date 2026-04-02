using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'commission_payouts' table in payment_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("commission_payouts", Schema = "payment_schema")]
public class CommissionPayoutEntity
{
    [Key]
    [Column("payout_id")]
    public Guid PayoutId { get; set; }

    [Column("payout_number")]
    public string PayoutNumber { get; set; } = $"PAY-{Guid.NewGuid().ToString()[..8].ToUpper()}";

    [Column("recipient_type")]
    public string RecipientType { get; set; } = string.Empty; // PARTNER, AGENT

    [Column("recipient_id")]
    public Guid RecipientId { get; set; }

    [Column("total_amount")]
    public long TotalAmount { get; set; }

    [Column("total_currency")]
    public string TotalCurrency { get; set; } = "BDT";

    [Column("commission_count")]
    public int CommissionCount { get; set; }

    [Column("period_start")]
    public DateTime PeriodStart { get; set; }

    [Column("period_end")]
    public DateTime PeriodEnd { get; set; }

    [Column("status")]
    public string Status { get; set; } = "PENDING"; // PENDING, APPROVED, PROCESSING, PAID, FAILED, CANCELLED

    [Column("payment_method")]
    public string? PaymentMethod { get; set; }

    [Column("payment_reference")]
    public string? PaymentReference { get; set; }

    [Column("paid_at")]
    public DateTime? PaidAt { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }
}
