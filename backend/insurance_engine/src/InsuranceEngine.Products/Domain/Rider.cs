using System;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Products.Domain;

/// <summary>
/// Product rider/add-on. Maps to 'product_riders' table.
/// </summary>
public class Rider
{
    public Guid Id { get; set; }
    public Guid ProductId { get; set; }
    public Product? Product { get; set; }

    public string RiderName { get; set; } = string.Empty;
    public string? Description { get; set; }

    // Money fields — stored as bigint (paisa)
    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public long CoverageAmount { get; set; }
    public string CoverageCurrency { get; set; } = "BDT";

    public bool IsMandatory { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Money convenience accessors
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
