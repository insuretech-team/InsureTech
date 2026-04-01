using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'product_riders' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("product_riders", Schema = "insurance_schema")]
public class ProductRiderEntity
{
    [Key]
    [Column("rider_id")]
    public Guid RiderId { get; set; }

    [Column("product_id")]
    public Guid ProductId { get; set; }

    [Column("rider_name")]
    public string RiderName { get; set; } = string.Empty;

    [Column("name_en")]
    public string NameEn { get; set; } = string.Empty;

    [Column("name_bn")]
    public string NameBn { get; set; } = string.Empty;

    [Column("description")]
    public string? Description { get; set; }

    [Column("premium_amount")]
    public long PremiumAmount { get; set; }

    [Column("premium_currency")]
    public string PremiumCurrency { get; set; } = "BDT";

    [Column("additional_premium")]
    public long AdditionalPremium { get; set; }

    [Column("additional_premium_currency")]
    public string AdditionalPremiumCurrency { get; set; } = "BDT";

    [Column("coverage_amount")]
    public long CoverageAmount { get; set; }

    [Column("coverage_currency")]
    public string CoverageCurrency { get; set; } = "BDT";

    [Column("additional_coverage")]
    public long AdditionalCoverage { get; set; }

    [Column("additional_coverage_currency")]
    public string AdditionalCoverageCurrency { get; set; } = "BDT";

    [Column("is_mandatory")]
    public bool IsMandatory { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public ProductEntity Product { get; set; } = null!;
}
