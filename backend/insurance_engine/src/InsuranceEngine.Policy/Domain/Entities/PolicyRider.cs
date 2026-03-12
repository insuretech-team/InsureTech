using System;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Policy rider/add-on attached to a policy. Maps to 'policy_riders' table.
/// </summary>
public class PolicyRider
{
    public Guid Id { get; set; }
    public Guid PolicyId { get; set; }

    public string RiderName { get; set; } = string.Empty;

    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public long CoverageAmount { get; set; }
    public string CoverageCurrency { get; set; } = "BDT";

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    public Money Premium
    {
        get => new(PremiumAmount, PremiumCurrency);
        set { PremiumAmount = value.Amount; PremiumCurrency = value.CurrencyCode; }
    }

    public Money Coverage
    {
        get => new(CoverageAmount, CoverageCurrency);
        set { CoverageAmount = value.Amount; CoverageCurrency = value.CurrencyCode; }
    }
}
