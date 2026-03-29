namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'commissions' table in insurance_schema.
/// Aligned with insuretech.partner.entity.v1.Commission proto.
/// </summary>
public class CommissionEntity
{
    public Guid CommissionId { get; set; }
    public string CommissionNumber { get; set; } = string.Empty;
    public Guid PolicyId { get; set; }
    public string CommissionType { get; set; } = "ACQUISITION"; // ACQUISITION, RENEWAL
    public string RecipientType { get; set; } = string.Empty; // AGENT, PARTNER, PLATFORM
    public Guid RecipientId { get; set; }
    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public decimal CommissionRate { get; set; }
    public long CommissionAmount { get; set; }
    public string CommissionCurrency { get; set; } = "BDT";
    public string? CalculationBreakdown { get; set; } // JSONB
    public string Status { get; set; } = "PENDING"; // PENDING, APPROVED, PAID
    public Guid? PayoutId { get; set; }
    public DateTime? PaidAt { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }

    // Navigation
    public PolicyEntity Policy { get; set; } = null!;
}
