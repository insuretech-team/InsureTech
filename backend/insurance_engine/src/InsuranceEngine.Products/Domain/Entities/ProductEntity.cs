using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.Products.Domain.Entities;

/// <summary>
/// EF Core entity for 'products' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations and used by PoliSync.
/// </summary>
[Table("products", Schema = "insurance_schema")]
public class ProductEntity
{
    [Key]
    [Column("product_id")]
    public Guid ProductId { get; set; }

    [Column("tenant_id")]
    public string TenantId { get; set; } = "default";

    [Column("product_code")]
    public string ProductCode { get; set; } = string.Empty;

    [Column("product_name")]
    public string ProductName { get; set; } = string.Empty;

    [Column("name_bn")]
    public string? NameBn { get; set; }

    [Column("product_type")]
    public string ProductType { get; set; } = "GENERAL";

    [Column("category")]
    public string Category { get; set; } = string.Empty;

    [Column("description")]
    public string? Description { get; set; }

    [Column("description_bn")]
    public string? DescriptionBn { get; set; }

    [Column("status")]
    public string Status { get; set; } = "DRAFT";

    [Column("is_active")]
    public bool IsActive { get; set; }

    [Column("base_premium")]
    public long BasePremium { get; set; }

    [Column("base_premium_currency")]
    public string BasePremiumCurrency { get; set; } = "BDT";

    [Column("min_sum_insured")]
    public long MinSumInsured { get; set; }

    [Column("min_sum_insured_currency")]
    public string MinSumInsuredCurrency { get; set; } = "BDT";

    [Column("max_sum_insured")]
    public long MaxSumInsured { get; set; }

    [Column("max_sum_insured_currency")]
    public string MaxSumInsuredCurrency { get; set; } = "BDT";

    [Column("unit_amount")]
    public long UnitAmount { get; set; } = 100000; // Default base unit of 1,000 BDT

    [Column("min_age")]
    public int MinAge { get; set; }

    [Column("max_age")]
    public int MaxAge { get; set; }

    [Column("min_term_months")]
    public int MinTenureMonths { get; set; }

    [Column("max_term_months")]
    public int MaxTenureMonths { get; set; }

    [Column("terms_url")]
    public string? TermsUrl { get; set; }

    [Column("questions")]
    public string? Questions { get; set; } // JSONB

    [Column("exclusions")]
    public string[] Exclusions { get; set; } = [];

    [Column("product_attributes")]
    public string? ProductAttributes { get; set; } // JSONB

    [Column("created_by")]
    public string CreatedBy { get; set; } = string.Empty;

    [Column("updated_by")]
    public string? UpdatedBy { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [Column("deleted_at")]
    public DateTime? DeletedAt { get; set; }

    [Column("version")]
    public int Version { get; set; } = 1;

    [Column("is_mandatory")]
    public bool IsMandatory { get; set; } = false;

    // Navigation properties
    public ICollection<ProductRiderEntity> Riders { get; set; } = new List<ProductRiderEntity>();
}
