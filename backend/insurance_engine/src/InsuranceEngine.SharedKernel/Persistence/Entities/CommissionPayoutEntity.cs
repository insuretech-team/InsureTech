namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'commission_payouts' table in insurance_schema.
/// Batch payout records for commissions.
/// </summary>
public class CommissionPayoutEntity
{
    public Guid PayoutId { get; set; }
    public string PayoutNumber { get; set; } = string.Empty;
    public string RecipientType { get; set; } = string.Empty;
    public Guid RecipientId { get; set; }
    public long TotalAmount { get; set; }
    public string TotalCurrency { get; set; } = "BDT";
    public int CommissionCount { get; set; }
    public string PeriodStart { get; set; } = string.Empty;
    public string PeriodEnd { get; set; } = string.Empty;
    public string Status { get; set; } = "PENDING"; // PENDING, PROCESSED, FAILED
    public string? PaymentMethod { get; set; }
    public string? PaymentReference { get; set; }
    public DateTime? PaidAt { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
