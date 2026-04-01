using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'individual_beneficiaries' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("individual_beneficiaries", Schema = "insurance_schema")]
public class IndividualBeneficiaryEntity
{
    [Key]
    [Column("beneficiary_id")]
    public Guid BeneficiaryId { get; set; }

    [Column("full_name")]
    public string FullName { get; set; } = string.Empty;

    [Column("full_name_bn")]
    public string? FullNameBn { get; set; }

    [Column("date_of_birth")]
    public DateTime DateOfBirth { get; set; }

    [Column("gender")]
    public string Gender { get; set; } = string.Empty;

    [Column("nid_number")]
    public string? NidNumber { get; set; }

    [Column("passport_number")]
    public string? PassportNumber { get; set; }

    [Column("birth_certificate_number")]
    public string? BirthCertificateNumber { get; set; }

    [Column("tin_number")]
    public string? TinNumber { get; set; }

    [Column("marital_status")]
    public string? MaritalStatus { get; set; }

    [Column("occupation")]
    public string? Occupation { get; set; }

    [Column("contact_info")]
    public string? ContactInfo { get; set; } // JSONB

    [Column("permanent_address")]
    public string? PermanentAddress { get; set; } // JSONB

    [Column("present_address")]
    public string? PresentAddress { get; set; } // JSONB

    [Column("nominee_name")]
    public string? NomineeName { get; set; }

    [Column("nominee_relationship")]
    public string? NomineeRelationship { get; set; }

    [Column("audit_info")]
    public string? AuditInfo { get; set; } // JSONB

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public BeneficiaryEntity Beneficiary { get; set; } = null!;
}
