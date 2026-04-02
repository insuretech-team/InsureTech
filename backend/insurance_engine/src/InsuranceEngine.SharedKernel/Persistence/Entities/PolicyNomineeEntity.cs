using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'policy_nominees' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("policy_nominees", Schema = "insurance_schema")]
public class PolicyNomineeEntity
{
    [Key]
    [Column("nominee_id")]
    public Guid NomineeId { get; set; }

    [Column("policy_id")]
    public Guid PolicyId { get; set; }

    [Column("full_name")]
    public string FullName { get; set; } = string.Empty;

    [Column("relationship")]
    public string Relationship { get; set; } = string.Empty;

    [Column("share_percentage")]
    public double SharePercentage { get; set; }

    [Column("date_of_birth")]
    public DateTime DateOfBirth { get; set; } = DateTime.UnixEpoch;

    [Column("nid_number")]
    public string? NidNumber { get; set; } // PII, encrypted

    [Column("identity_type")]
    public string? IdentityType { get; set; }

    [Column("identity_number")]
    public string? IdentityNumber { get; set; }

    [Column("phone_number")]
    public string? PhoneNumber { get; set; } // PII, encrypted

    [Column("nominee_dob_text")]
    public string? NomineeDobText { get; set; }

    [Column("nominee_share_percent")]
    public double? NomineeSharePercent { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public PolicyEntity Policy { get; set; } = null!;
}
