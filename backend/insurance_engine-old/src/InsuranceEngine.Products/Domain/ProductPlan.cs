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
    public string? PlanDescription { get; set; }

    // Money fields — stored as bigint (paisa)
    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public long MinSumInsured { get; set; }
    public string MinSumInsuredCurrency { get; set; } = "BDT";
    public long MaxSumInsured { get; set; }
    public string MaxSumInsuredCurrency { get; set; } = "BDT";

    public string? Attributes { get; set; } // JSONB

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Money convenience accessors
    public Money PremiumMoney
    {
        get => new(PremiumAmount, PremiumCurrency);
        set { PremiumAmount = value.Amount; PremiumCurrency = value.CurrencyCode; }
    }

    public Money MinSumInsuredMoney
    {
        get => new(MinSumInsured, MinSumInsuredCurrency);
        set { MinSumInsured = value.Amount; MinSumInsuredCurrency = value.CurrencyCode; }
    }

    public Money MaxSumInsuredMoney
    {
        get => new(MaxSumInsured, MaxSumInsuredCurrency);
        set { MaxSumInsured = value.Amount; MaxSumInsuredCurrency = value.CurrencyCode; }
    }
}
