namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'policy_nominees' table in insurance_schema.
/// Aligned with insuretech.policy.entity.v1.Nominee proto definition.
/// </summary>
public class PolicyNomineeEntity
{
    public Guid NomineeId { get; set; }
    public Guid PolicyId { get; set; }
    public string FullName { get; set; } = string.Empty;
    public string Relationship { get; set; } = string.Empty;
    public double SharePercentage { get; set; }
    public DateTime DateOfBirth { get; set; }
    public string? NidNumber { get; set; } // PII, encrypted
    public string? PhoneNumber { get; set; } // PII, encrypted
    public string? NomineeDobText { get; set; }
    public double? NomineeSharePercent { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public PolicyEntity Policy { get; set; } = null!;
}
