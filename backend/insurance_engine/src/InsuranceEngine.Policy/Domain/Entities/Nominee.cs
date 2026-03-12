using System;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Nominee/Beneficiary. Maps to 'policy_nominees' table.
/// </summary>
public class Nominee
{
    public Guid Id { get; set; }
    public Guid PolicyId { get; set; }
    public Guid BeneficiaryId { get; set; }
    public Beneficiary Beneficiary { get; set; } = null!;

    public string Relationship { get; set; } = string.Empty;
    public double SharePercentage { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public bool IsDeleted { get; set; }
}
