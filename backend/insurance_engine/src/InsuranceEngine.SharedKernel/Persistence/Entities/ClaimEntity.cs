namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'claims' table in insurance_schema.
/// Aligned with insuretech.claims.entity.v1.Claim proto definition.
/// Partitioned by RANGE on created_at (monthly).
/// </summary>
public class ClaimEntity
{
    public Guid ClaimId { get; set; }
    public string ClaimNumber { get; set; } = string.Empty;
    public Guid PolicyId { get; set; }
    public Guid CustomerId { get; set; }
    public string Status { get; set; } = "SUBMITTED";
    public string Type { get; set; } = string.Empty;
    public long ClaimedAmount { get; set; }
    public string ClaimedCurrency { get; set; } = "BDT";
    public long? ApprovedAmount { get; set; }
    public string ApprovedCurrency { get; set; } = "BDT";
    public long? SettledAmount { get; set; }
    public string SettledCurrency { get; set; } = "BDT";
    public DateTime IncidentDate { get; set; }
    public string IncidentDescription { get; set; } = string.Empty;
    public DateTime SubmittedAt { get; set; }
    public DateTime? ApprovedAt { get; set; }
    public DateTime? SettledAt { get; set; }
    public string? RejectionReason { get; set; }
    public string? PlaceOfIncident { get; set; }
    public string? BankDetailsForPayout { get; set; } // PII, encrypted
    public bool AppealOptionAvailable { get; set; }
    public string? InAppMessages { get; set; } // JSONB
    public string ProcessingType { get; set; } = "MANUAL";
    public long? DeductibleAmount { get; set; }
    public long? CoPayAmount { get; set; }
    public string? ProcessorNotes { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }

    // Navigation properties
    public PolicyEntity Policy { get; set; } = null!;
    public ICollection<ClaimDocumentEntity> Documents { get; set; } = new List<ClaimDocumentEntity>();
    public ICollection<ClaimApprovalEntity> Approvals { get; set; } = new List<ClaimApprovalEntity>();
    public FraudCheckEntity? FraudCheck { get; set; }
}
