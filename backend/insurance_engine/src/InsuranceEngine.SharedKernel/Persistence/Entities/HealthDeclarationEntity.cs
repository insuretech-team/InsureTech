using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'health_declarations' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("health_declarations", Schema = "insurance_schema")]
public class HealthDeclarationEntity
{
    [Key]
    [Column("declaration_id")]
    public Guid DeclarationId { get; set; }

    [Column("quote_id")]
    public Guid QuoteId { get; set; }

    [Column("height_cm")]
    public int HeightCm { get; set; }

    [Column("weight_kg")]
    public string WeightKg { get; set; } = string.Empty;

    [Column("bmi")]
    public decimal? Bmi { get; set; }

    [Column("has_pre_existing_conditions")]
    public bool HasPreExistingConditions { get; set; }

    [Column("pre_existing_conditions")]
    public string? PreExistingConditions { get; set; } // JSONB

    [Column("is_currently_hospitalized")]
    public bool IsCurrentlyHospitalized { get; set; }

    [Column("has_family_history")]
    public bool HasFamilyHistory { get; set; }

    [Column("family_history")]
    public string? FamilyHistory { get; set; } // JSONB

    [Column("smoker")]
    public bool Smoker { get; set; }

    [Column("alcohol_consumer")]
    public bool AlcoholConsumer { get; set; }

    [Column("occupation_risk_level")]
    public string? OccupationRiskLevel { get; set; }

    [Column("medical_exam_required")]
    public bool MedicalExamRequired { get; set; }

    [Column("medical_exam_completed")]
    public bool MedicalExamCompleted { get; set; }

    [Column("medical_exam_results")]
    public string? MedicalExamResults { get; set; } // JSONB

    [Column("medical_exam_status")]
    public string? MedicalExamStatus { get; set; }

    [Column("medical_exam_date")]
    public DateTime? MedicalExamDate { get; set; }

    [Column("medical_record_numbers")]
    public string[]? MedicalRecordNumbers { get; set; }

    [Column("medical_comments")]
    public string? MedicalComments { get; set; }

    [Column("medical_review_status")]
    public string? MedicalReviewStatus { get; set; }

    [Column("medical_documents")]
    public string? MedicalDocuments { get; set; } // JSONB

    [Column("auto_approval_possible")]
    public bool AutoApprovalPossible { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public QuoteEntity Quote { get; set; } = null!;
}
