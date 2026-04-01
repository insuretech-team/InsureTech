using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'claims' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("claims", Schema = "insurance_schema")]
public class ClaimEntity
{
    [Key]
    [Column("claim_id")]
    public Guid ClaimId { get; set; }

    [Column("claim_number")]
    public string ClaimNumber { get; set; } = string.Empty;

    [Column("policy_id")]
    public Guid PolicyId { get; set; }

    [Column("customer_id")]
    public Guid CustomerId { get; set; }

    [Column("status")]
    public string Status { get; set; } = "SUBMITTED";

    [Column("type")]
    public string Type { get; set; } = string.Empty;

    [Column("claimed_amount")]
    public long ClaimedAmount { get; set; }

    [Column("claimed_currency")]
    public string ClaimedCurrency { get; set; } = "BDT";

    [Column("approved_amount")]
    public long? ApprovedAmount { get; set; }

    [Column("approved_currency")]
    public string ApprovedCurrency { get; set; } = "BDT";

    [Column("settled_amount")]
    public long? SettledAmount { get; set; }

    [Column("settled_currency")]
    public string SettledCurrency { get; set; } = "BDT";

    [Column("incident_date")]
    public DateTime IncidentDate { get; set; }

    [Column("incident_description")]
    public string IncidentDescription { get; set; } = string.Empty;

    [Column("submitted_at")]
    public DateTime SubmittedAt { get; set; }

    [Column("approved_at")]
    public DateTime? ApprovedAt { get; set; }

    [Column("settled_at")]
    public DateTime? SettledAt { get; set; }

    [Column("rejection_reason")]
    public string? RejectionReason { get; set; }

    [Column("place_of_incident")]
    public string? PlaceOfIncident { get; set; }

    [Column("bank_details_for_payout")]
    public string? BankDetailsForPayout { get; set; } // PII, encrypted

    [Column("appeal_option_available")]
    public bool AppealOptionAvailable { get; set; }

    [Column("in_app_messages")]
    public string? InAppMessages { get; set; } // JSONB

    [Column("processing_type")]
    public string ProcessingType { get; set; } = "MANUAL";

    [Column("claim_source")]
    public string? ClaimSource { get; set; } // JSONB

    [Column("is_priority")]
    public bool IsPriority { get; set; }

    [Column("deductible_amount")]
    public long? DeductibleAmount { get; set; }

    [Column("co_pay_amount")]
    public long? CoPayAmount { get; set; }

    [Column("processor_notes")]
    public string? ProcessorNotes { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [Column("deleted_at")]
    public DateTime? DeletedAt { get; set; }

    // Navigation properties
    public PolicyEntity Policy { get; set; } = null!;
    public ICollection<ClaimDocumentEntity> Documents { get; set; } = new List<ClaimDocumentEntity>();
    public ICollection<ClaimApprovalEntity> Approvals { get; set; } = new List<ClaimApprovalEntity>();
    public FraudCheckEntity? FraudCheck { get; set; }
}
