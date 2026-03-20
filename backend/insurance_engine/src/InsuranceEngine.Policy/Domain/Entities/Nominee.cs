using System;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Nominee/Beneficiary. Maps to 'policy_nominees' table.
/// Proto: insuretech.policy.entity.v1.Nominee
/// </summary>
public class Nominee
{
    public Guid Id { get; set; }
    public Guid PolicyId { get; set; }

    /// <summary>
    /// Optional link to Beneficiary entity. Proto stores nominee data inline,
    /// but we keep this FK for when beneficiary record exists.
    /// </summary>
    public Guid? BeneficiaryId { get; set; }
    public Beneficiary? Beneficiary { get; set; }

    // --- Proto-aligned inline fields ---

    /// <summary>
    /// Nominee's full name. Proto: full_name VARCHAR(200)
    /// </summary>
    public string FullName { get; set; } = string.Empty;

    /// <summary>
    /// Relationship to policyholder. Proto: relationship VARCHAR(50)
    /// </summary>
    public string Relationship { get; set; } = string.Empty;

    /// <summary>
    /// Share percentage (0-100). Proto: nominee_share_percent DOUBLE PRECISION
    /// Example: 50.00 means 50%
    /// </summary>
    public double SharePercentage { get; set; }

    /// <summary>
    /// Nominee's date of birth. Proto: date_of_birth DATE
    /// </summary>
    public DateTime? DateOfBirth { get; set; }

    /// <summary>
    /// Nominee's date of birth as text (for display when exact date unknown).
    /// Proto: nominee_dob_text VARCHAR(50)
    /// </summary>
    public string? NomineeDobText { get; set; }

    /// <summary>
    /// National ID number. Proto: nid_number VARCHAR(20)
    /// </summary>
    public string? NidNumber { get; set; }

    /// <summary>
    /// Phone number. Proto: phone_number VARCHAR(20)
    /// </summary>
    public string? PhoneNumber { get; set; }

    // Audit
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }
    public bool IsDeleted { get; set; }
}

