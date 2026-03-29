using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'products' table in insurance_schema.
/// Aligned with insuretech.products.entity.v1.Product proto definition.
/// </summary>
public class ProductEntity
{
    public Guid ProductId { get; set; }
    public string ProductCode { get; set; } = string.Empty;
    public string ProductName { get; set; } = string.Empty;
    public string Category { get; set; } = string.Empty;
    public string? Description { get; set; }
    public long BasePremium { get; set; }
    public string BasePremiumCurrency { get; set; } = "BDT";
    public long MinSumInsured { get; set; }
    public string MinSumInsuredCurrency { get; set; } = "BDT";
    public long MaxSumInsured { get; set; }
    public string MaxSumInsuredCurrency { get; set; } = "BDT";
    public int MinTenureMonths { get; set; }
    public int MaxTenureMonths { get; set; }
    public int MinAge { get; set; }
    public int MaxAge { get; set; }
    public string? TermsUrl { get; set; }
    public string? Questions { get; set; } // JSONB (RiskAssessmentQuestion[])
    public string[] Exclusions { get; set; } = [];
    public string Status { get; set; } = "DRAFT";
    public string? ProductAttributes { get; set; } // JSONB
    public string CreatedBy { get; set; } = string.Empty;
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }
    public int Version { get; set; } = 1;

    // Navigation properties
    public ICollection<PolicyEntity> Policies { get; set; } = new List<PolicyEntity>();
    public ICollection<ProductRiderEntity> Riders { get; set; } = new List<ProductRiderEntity>();
}
