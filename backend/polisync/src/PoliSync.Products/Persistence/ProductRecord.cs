using System.Text.Json;

namespace PoliSync.Products.Persistence;

/// <summary>
/// EF Core POCO for the insurance_schema.products table.
/// Separate from the proto-generated Insuretech.Products.Entity.V1.Product class.
/// Proto types are used at the gRPC boundary; this type is used only for DB persistence.
/// </summary>
public sealed class ProductRecord
{
    public Guid ProductId { get; set; }
    public string ProductCode { get; set; } = string.Empty;
    public string ProductName { get; set; } = string.Empty;
    public string Category { get; set; } = string.Empty;        // stored as VARCHAR enum string
    public string? Description { get; set; }
    public long BasePremium { get; set; }                       // paisa
    public string BasePremiumCurrency { get; set; } = "BDT";
    public long MinSumInsured { get; set; }                     // paisa
    public string MinSumInsuredCurrency { get; set; } = "BDT";
    public long MaxSumInsured { get; set; }                     // paisa
    public string MaxSumInsuredCurrency { get; set; } = "BDT";
    public int MinTenureMonths { get; set; }
    public int MaxTenureMonths { get; set; }
    public string[] Exclusions { get; set; } = [];
    public string Status { get; set; } = "PRODUCT_STATUS_DRAFT"; // proto enum string
    public string? ProductAttributes { get; set; }               // JSONB
    public Guid CreatedBy { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }                     // soft delete

    // Navigation
    public List<ProductPlanRecord> Plans { get; set; } = [];
    public List<RiderRecord> Riders { get; set; } = [];
    public PricingConfigRecord? PricingConfig { get; set; }
}

/// <summary>
/// EF Core POCO for insurance_schema.product_plans.
/// </summary>
public sealed class ProductPlanRecord
{
    public Guid PlanId { get; set; }
    public Guid ProductId { get; set; }
    public string PlanName { get; set; } = string.Empty;
    public string? PlanDescription { get; set; }
    public long PremiumAmount { get; set; }                    // paisa
    public string PremiumCurrency { get; set; } = "BDT";
    public long MinSumInsured { get; set; }
    public string MinSumInsuredCurrency { get; set; } = "BDT";
    public long MaxSumInsured { get; set; }
    public string MaxSumInsuredCurrency { get; set; } = "BDT";
    public string? Attributes { get; set; }                    // JSONB
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public ProductRecord Product { get; set; } = null!;
}

/// <summary>
/// EF Core POCO for insurance_schema.product_riders.
/// </summary>
public sealed class RiderRecord
{
    public Guid RiderId { get; set; }
    public Guid ProductId { get; set; }
    public string RiderName { get; set; } = string.Empty;
    public string? Description { get; set; }
    public long PremiumAmount { get; set; }                    // paisa
    public string PremiumCurrency { get; set; } = "BDT";
    public long CoverageAmount { get; set; }                   // paisa
    public string CoverageCurrency { get; set; } = "BDT";
    public bool IsMandatory { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public ProductRecord Product { get; set; } = null!;
}

/// <summary>
/// EF Core POCO for insurance_schema.pricing_configs.
/// Rules stored as JSONB — deserialized to PricingRuleJson on read.
/// </summary>
public sealed class PricingConfigRecord
{
    public Guid PricingConfigId { get; set; }
    public Guid ProductId { get; set; }
    public string Rules { get; set; } = "[]";                  // JSONB stored as string
    public DateTime EffectiveFrom { get; set; }
    public DateTime? EffectiveTo { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public ProductRecord Product { get; set; } = null!;
}
