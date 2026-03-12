using System;
using System.Collections.Generic;

namespace InsuranceEngine.Policy.Domain.ValueObjects;

/// <summary>
/// Applicant/Proposer information. Stored as JSONB in policies table.
/// </summary>
public record Applicant
{
    public string FullName { get; init; } = string.Empty;
    public DateTime? DateOfBirth { get; init; }
    public string? NidNumber { get; init; }      // Encrypted at rest
    public string? Occupation { get; init; }
    public long AnnualIncome { get; init; }       // paisa
    public string? Address { get; init; }
    public string? PhoneNumber { get; init; }     // Encrypted at rest
    public HealthDeclaration? HealthDeclaration { get; init; }
}

/// <summary>
/// Health declaration information. Stored as JSONB.
/// </summary>
public record HealthDeclaration
{
    public bool HasPreExistingConditions { get; init; }
    public List<string> Conditions { get; init; } = new();
    public bool IsSmoker { get; init; }
    public string? BloodGroup { get; init; }
}
