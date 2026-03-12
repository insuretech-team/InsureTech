using System;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Health declaration for underwriting.
/// Maps to 'health_declarations' table in insurance_schema.
/// </summary>
public class UnderwritingHealthDeclaration
{
    public Guid Id { get; set; }
    public Guid QuoteId { get; set; }

    // Physical measurements
    public int HeightCm { get; set; }
    public decimal WeightKg { get; set; }
    public decimal Bmi { get; set; }

    // Medical history
    public bool HasPreExistingConditions { get; set; }
    public string? PreExistingConditionsJson { get; set; } // Encrypted JSONB
    public bool IsCurrentlyHospitalized { get; set; }
    public bool HasFamilyHistory { get; set; }
    public string? FamilyHistoryJson { get; set; } // Encrypted JSONB

    // Lifestyle
    public bool IsSmoker { get; set; }
    public bool IsAlcoholConsumer { get; set; }
    public string? OccupationRiskLevel { get; set; } // LOW, MEDIUM, HIGH

    // Medical examination
    public bool IsMedicalExamRequired { get; set; }
    public bool IsMedicalExamCompleted { get; set; }
    public string? MedicalExamResultsJson { get; set; } // Encrypted JSONB
    public DateTime? MedicalExamDate { get; set; }

    // Documents
    public string? MedicalDocumentsJson { get; set; } // JSONB list of document references

    // Audit Info JSONB
    public string? AuditInfoJson { get; set; }
    
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
