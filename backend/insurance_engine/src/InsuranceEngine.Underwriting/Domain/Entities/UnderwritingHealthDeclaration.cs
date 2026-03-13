using System;
using InsuranceEngine.Underwriting.Domain.Enums;

namespace InsuranceEngine.Underwriting.Domain.Entities;

public class UnderwritingHealthDeclaration
{
    public Guid Id { get; set; }
    public Guid QuoteId { get; set; }

    public int HeightCm { get; set; }
    public decimal WeightKg { get; set; }
    public decimal Bmi { get; set; }

    public bool HasPreExistingConditions { get; set; }
    public string? PreExistingConditionsJson { get; set; }
    public bool IsCurrentlyHospitalized { get; set; }
    public bool HasFamilyHistory { get; set; }
    public string? FamilyHistoryJson { get; set; }

    public bool IsSmoker { get; set; }
    public bool IsAlcoholConsumer { get; set; }
    public string? OccupationRiskLevel { get; set; }

    public bool IsMedicalExamRequired { get; set; }
    public bool IsMedicalExamCompleted { get; set; }
    public string? MedicalExamResultsJson { get; set; }
    public DateTime? MedicalExamDate { get; set; }

    public string? MedicalDocumentsJson { get; set; }

    public string? AuditInfoJson { get; set; }
    
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
