namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'health_declarations' table in insurance_schema.
/// Aligned with insuretech.underwriting.entity.v1.HealthDeclaration proto.
/// </summary>
public class HealthDeclarationEntity
{
    public Guid DeclarationId { get; set; }
    public Guid QuoteId { get; set; }
    public int HeightCm { get; set; }
    public string WeightKg { get; set; } = string.Empty;
    public bool HasPreExistingConditions { get; set; }
    public string? PreExistingConditions { get; set; } // JSON
    public bool Smoker { get; set; }
    public bool AlcoholConsumer { get; set; }
    public string? OccupationRiskLevel { get; set; }
    public bool MedicalExamRequired { get; set; }
    public bool AutoApprovalPossible { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public QuoteEntity Quote { get; set; } = null!;
}
