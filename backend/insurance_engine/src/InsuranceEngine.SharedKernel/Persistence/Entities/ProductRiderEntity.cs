namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'product_riders' table in insurance_schema.
/// Aligned with insuretech.products.entity.v1.Rider proto definition.
/// </summary>
public class ProductRiderEntity
{
    public Guid RiderId { get; set; }
    public Guid ProductId { get; set; }
    public string NameEn { get; set; } = string.Empty;
    public string NameBn { get; set; } = string.Empty;
    public long AdditionalPremium { get; set; }
    public string AdditionalPremiumCurrency { get; set; } = "BDT";
    public long AdditionalCoverage { get; set; }
    public string AdditionalCoverageCurrency { get; set; } = "BDT";
    public DateTime CreatedAt { get; set; }

    // Navigation
    public ProductEntity Product { get; set; } = null!;
}
