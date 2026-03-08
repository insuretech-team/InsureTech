using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace PoliSync.Infrastructure.Entities;

/// <summary>
/// EF Core entity for Product table
/// Maps to insurance_schema.products table managed by Go migrations
/// </summary>
[Table("products", Schema = "insurance_schema")]
public class ProductEntity
{
    [Key]
    [Column("product_id")]
    public string ProductId { get; set; } = string.Empty;

    [Column("tenant_id")]
    public string TenantId { get; set; } = string.Empty;

    [Column("product_code")]
    public string ProductCode { get; set; } = string.Empty;

    [Column("product_name")]
    public string ProductName { get; set; } = string.Empty;

    [Column("product_type")]
    public string ProductType { get; set; } = string.Empty;

    [Column("category")]
    public string Category { get; set; } = string.Empty;

    [Column("description")]
    public string? Description { get; set; }

    [Column("status")]
    public string Status { get; set; } = string.Empty;

    [Column("is_active")]
    public bool IsActive { get; set; }

    [Column("min_coverage_amount")]
    public decimal MinCoverageAmount { get; set; }

    [Column("max_coverage_amount")]
    public decimal MaxCoverageAmount { get; set; }

    [Column("currency")]
    public string Currency { get; set; } = "BDT";

    [Column("min_age")]
    public int MinAge { get; set; }

    [Column("max_age")]
    public int MaxAge { get; set; }

    [Column("min_term_months")]
    public int MinTermMonths { get; set; }

    [Column("max_term_months")]
    public int MaxTermMonths { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [Column("created_by")]
    public string? CreatedBy { get; set; }

    [Column("updated_by")]
    public string? UpdatedBy { get; set; }
}
