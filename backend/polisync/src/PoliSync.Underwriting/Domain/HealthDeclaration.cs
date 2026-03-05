using PoliSync.SharedKernel.Domain;

namespace PoliSync.Underwriting.Domain;

public class HealthDeclaration : Entity
{
    public Guid DeclarationId { get; private set; }
    public Guid QuoteId { get; private set; } // FK to Quote
    
    // Physical attributes
    public float HeightCm { get; private set; }
    public float WeightKg { get; private set; }
    public float Bmi { get; private set; }
    
    // Habits
    public bool IsSmoker { get; private set; }
    public bool ConsumesAlcohol { get; private set; }
    
    // Medical history
    public bool HasPreExistingConditions { get; private set; }
    public string? ConditionDetails { get; private set; } // JSONB
    
    // Risk factors
    public bool HasFamilyHistoryOfCriticalIllness { get; private set; }
    public string? OccupationRiskLevel { get; private set; } 
    
    public DateTime SubmittedAt { get; private set; }

    private HealthDeclaration() { }

    public static HealthDeclaration Create(
        Guid quoteId,
        float heightCm,
        float weightKg,
        bool isSmoker,
        bool consumesAlcohol,
        bool hasPreExistingConditions,
        string? conditionDetails,
        bool hasFamilyHistoryOfCriticalIllness,
        string occupationRiskLevel)
    {
        float heightM = heightCm / 100f;
        float bmi = weightKg / (heightM * heightM);

        return new HealthDeclaration
        {
            DeclarationId = Guid.NewGuid(),
            QuoteId = quoteId,
            HeightCm = heightCm,
            WeightKg = weightKg,
            Bmi = bmi,
            IsSmoker = isSmoker,
            ConsumesAlcohol = consumesAlcohol,
            HasPreExistingConditions = hasPreExistingConditions,
            ConditionDetails = conditionDetails ?? "[]",
            HasFamilyHistoryOfCriticalIllness = hasFamilyHistoryOfCriticalIllness,
            OccupationRiskLevel = occupationRiskLevel,
            SubmittedAt = DateTime.UtcNow
        };
    }
}
