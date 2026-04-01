using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'business_beneficiaries' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("business_beneficiaries", Schema = "insurance_schema")]
public class BusinessBeneficiaryEntity
{
    [Key]
    [Column("beneficiary_id")]
    public Guid BeneficiaryId { get; set; } // Primary key

    [Column("parent_beneficiary_id")]
    public Guid ParentBeneficiaryId { get; set; } // Link to beneficiaries table

    [Column("business_name")]
    public string BusinessName { get; set; } = string.Empty;

    [Column("business_name_bn")]
    public string? BusinessNameBn { get; set; }

    [Column("trade_license_number")]
    public string TradeLicenseNumber { get; set; } = string.Empty;

    [Column("trade_license_issue_date")]
    public DateTime? TradeLicenseIssueDate { get; set; }

    [Column("trade_license_expiry_date")]
    public DateTime? TradeLicenseExpiryDate { get; set; }

    [Column("tin_number")]
    public string TinNumber { get; set; } = string.Empty;

    [Column("bin_number")]
    public string? BinNumber { get; set; }

    [Column("business_type")]
    public string BusinessType { get; set; } = string.Empty;

    [Column("industry_sector")]
    public string? IndustrySector { get; set; }

    [Column("employee_count")]
    public int? EmployeeCount { get; set; }

    [Column("incorporation_date")]
    public DateTime? IncorporationDate { get; set; }

    [Column("contact_info")]
    public string? ContactInfo { get; set; } // JSONB

    [Column("registered_address")]
    public string? RegisteredAddress { get; set; } // JSONB

    [Column("business_address")]
    public string? BusinessAddress { get; set; } // JSONB

    [Column("focal_person_name")]
    public string FocalPersonName { get; set; } = string.Empty;

    [Column("focal_person_designation")]
    public string? FocalPersonDesignation { get; set; }

    [Column("focal_person_nid")]
    public string? FocalPersonNid { get; set; }

    [Column("focal_person_contact")]
    public string? FocalPersonContact { get; set; } // JSONB

    [Column("registration_number")]
    public string? RegistrationNumber { get; set; }

    [Column("tax_id")]
    public string? TaxId { get; set; }

    [Column("primary_contact")]
    public string? PrimaryContact { get; set; } // JSONB

    [Column("total_employees_covered")]
    public int TotalEmployeesCovered { get; set; }

    [Column("active_policies_count")]
    public int ActivePoliciesCount { get; set; }

    [Column("total_premium_amount")]
    public long TotalPremiumAmount { get; set; } // Paisa

    [Column("pending_actions_count")]
    public int PendingActionsCount { get; set; }

    [Column("audit_info")]
    public string? AuditInfo { get; set; } // JSONB

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public BeneficiaryEntity Beneficiary { get; set; } = null!;
}
