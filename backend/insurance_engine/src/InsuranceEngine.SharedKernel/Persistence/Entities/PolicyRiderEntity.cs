namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'policy_riders' table in insurance_schema.
/// Aligned with insuretech.policy.entity.v1.Rider proto definition.
/// </summary>
public class PolicyRiderEntity
{
    public Guid RiderId { get; set; }
    public Guid PolicyId { get; set; }
    public string RiderName { get; set; } = string.Empty;
    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public long CoverageAmount { get; set; }
    public string CoverageCurrency { get; set; } = "BDT";
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public PolicyEntity Policy { get; set; } = null!;
}
