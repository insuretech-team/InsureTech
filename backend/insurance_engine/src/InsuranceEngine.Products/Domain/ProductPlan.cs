using System;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Products.Domain;

/// <summary>
/// Product plan variant (e.g., Basic, Silver, Gold). Maps to 'product_plans' table.
/// </summary>
public class ProductPlan
{
    public Guid Id { get; set; }
    public Guid ProductId { get; set; }
    public Product? Product { get; set; }

    public string PlanName { get; set; } = string.Empty;
    public string? PlanNameBn { get; set; }
    public string? PlanDescription { get; set; }
    public string? DescriptionBn { get; set; }

    // Money fields — stored as bigint (paisa)
    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public long MinSumInsuredAmount { get; set; }
    public string MinSumInsuredCurrency { get; set; } = "BDT";
    public long MaxSumInsuredAmount { get; set; }
    public string MaxSumInsuredCurrency { get; set; } = "BDT";

    public string? Attributes { get; set; } // JSONB

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Money convenience accessors
    public Money Premium
    {
        get => new(PremiumAmount, PremiumCurrency);
        set { PremiumAmount = value.Amount; PremiumCurrency = value.CurrencyCode; }
    }

    public Money MinSumInsured
    {
        get => new(MinSumInsuredAmount, MinSumInsuredCurrency);
        set { MinSumInsuredAmount = value.Amount; MinSumInsuredCurrency = value.CurrencyCode; }
    }

    public Money MaxSumInsured
    {
        get => new(MaxSumInsuredAmount, MaxSumInsuredCurrency);
        set { MaxSumInsuredAmount = value.Amount; MaxSumInsuredCurrency = value.CurrencyCode; }
    }
}
