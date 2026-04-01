using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'policy_riders' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("policy_riders", Schema = "insurance_schema")]
public class PolicyRiderEntity
{
    [Key]
    [Column("rider_id")]
    public Guid RiderId { get; set; }

    [Column("policy_id")]
    public Guid PolicyId { get; set; }

    [Column("rider_name")]
    public string RiderName { get; set; } = string.Empty;

    [Column("premium_amount")]
    public long PremiumAmount { get; set; }

    [Column("premium_currency")]
    public string PremiumCurrency { get; set; } = "BDT";

    [Column("coverage_amount")]
    public long CoverageAmount { get; set; }

    [Column("coverage_currency")]
    public string CoverageCurrency { get; set; } = "BDT";

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public PolicyEntity Policy { get; set; } = null!;
}
