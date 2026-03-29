namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'policies' table in insurance_schema.
/// Aligned with insuretech.policy.entity.v1.Policy proto definition.
/// Partitioned by RANGE on created_at (yearly).
/// </summary>
public class PolicyEntity
{
    public Guid PolicyId { get; set; }
    public string PolicyNumber { get; set; } = string.Empty;
    public Guid ProductId { get; set; }
    public Guid CustomerId { get; set; }
    public Guid? PartnerId { get; set; }
    public Guid? AgentId { get; set; }
    public Guid? QuoteId { get; set; }
    public Guid? UnderwritingDecisionId { get; set; }
    public string Status { get; set; } = "PENDING_PAYMENT";
    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public long SumInsured { get; set; }
    public string SumInsuredCurrency { get; set; } = "BDT";
    public int TenureMonths { get; set; }
    public DateTime StartDate { get; set; }
    public DateTime EndDate { get; set; }
    public DateTime? IssuedAt { get; set; }
    public string? PolicyDocumentUrl { get; set; }
    public string? PaymentFrequency { get; set; }
    public long? VatTax { get; set; }
    public long? ServiceFee { get; set; }
    public long? TotalPayable { get; set; }
    public string? PaymentGatewayReference { get; set; }
    public string? ReceiptNumber { get; set; }
    public string? OccupationRiskClass { get; set; }
    public bool HasExistingPolicies { get; set; }
    public string? ClaimsHistorySummary { get; set; }
    public string? ProviderName { get; set; }
    public DateTime? EnrollmentStartDate { get; set; }
    public DateTime? EnrollmentEndDate { get; set; }
    public string? UnderwritingData { get; set; } // JSONB
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }

    // Navigation properties
    public ProductEntity Product { get; set; } = null!;
    public ICollection<PolicyNomineeEntity> Nominees { get; set; } = new List<PolicyNomineeEntity>();
    public ICollection<PolicyRiderEntity> Riders { get; set; } = new List<PolicyRiderEntity>();
    public ICollection<ClaimEntity> Claims { get; set; } = new List<ClaimEntity>();
}
