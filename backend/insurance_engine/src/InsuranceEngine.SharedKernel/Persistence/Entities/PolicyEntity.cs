using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'policies' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("policies", Schema = "insurance_schema")]
public class PolicyEntity
{
    [Key]
    [Column("policy_id")]
    public Guid PolicyId { get; set; }

    [Column("policy_number")]
    public string PolicyNumber { get; set; } = string.Empty;

    [Column("product_id")]
    public Guid ProductId { get; set; }

    [Column("customer_id")]
    public Guid CustomerId { get; set; }

    [Column("partner_id")]
    public Guid? PartnerId { get; set; }

    [Column("agent_id")]
    public Guid? AgentId { get; set; }

    [Column("quote_id")]
    public Guid? QuoteId { get; set; }

    [Column("underwriting_decision_id")]
    public Guid? UnderwritingDecisionId { get; set; }

    [Column("status")]
    public string Status { get; set; } = "PENDING_PAYMENT";

    [Column("premium_amount")]
    public long PremiumAmount { get; set; }

    [NotMapped]
    public long Premium { get => PremiumAmount; set => PremiumAmount = value; }

    [Column("premium_currency")]
    public string PremiumCurrency { get; set; } = "BDT";

    [Column("sum_insured_amount")]
    public long SumInsuredAmount { get; set; }

    [NotMapped]
    public long SumInsured { get => SumInsuredAmount; set => SumInsuredAmount = value; }

    [Column("sum_insured_currency")]
    public string SumInsuredCurrency { get; set; } = "BDT";

    [Column("tenure_months")]
    public int TenureMonths { get; set; }

    [Column("start_date")]
    public DateTime StartDate { get; set; }

    [Column("end_date")]
    public DateTime EndDate { get; set; }

    [Column("issued_at")]
    public DateTime? IssuedAt { get; set; }

    [Column("policy_document_url")]
    public string? PolicyDocumentUrl { get; set; }

    [Column("payment_frequency")]
    public string? PaymentFrequency { get; set; }

    [Column("vat_tax_amount")]
    public long? VatTaxAmount { get; set; }

    [NotMapped]
    public long? VatTax { get => VatTaxAmount; set => VatTaxAmount = value; }

    [Column("vat_tax_currency")]
    public string VatTaxCurrency { get; set; } = "BDT";

    [Column("service_fee_amount")]
    public long? ServiceFeeAmount { get; set; }

    [NotMapped]
    public long? ServiceFee { get => ServiceFeeAmount; set => ServiceFeeAmount = value; }

    [Column("service_fee_currency")]
    public string ServiceFeeCurrency { get; set; } = "BDT";

    [Column("total_payable_amount")]
    public long? TotalPayableAmount { get; set; }

    [NotMapped]
    public long? TotalPayable { get => TotalPayableAmount; set => TotalPayableAmount = value; }

    [Column("total_payable_currency")]
    public string TotalPayableCurrency { get; set; } = "BDT";

    [Column("payment_gateway_reference")]
    public string? PaymentGatewayReference { get; set; }

    [Column("receipt_number")]
    public string? ReceiptNumber { get; set; }

    [Column("occupation_risk_class")]
    public string? OccupationRiskClass { get; set; }

    [Column("has_existing_policies")]
    public bool HasExistingPolicies { get; set; }

    [Column("claims_history_summary")]
    public string? ClaimsHistorySummary { get; set; }

    [Column("provider_name")]
    public string? ProviderName { get; set; }

    [Column("enrollment_start_date")]
    public DateTime? EnrollmentStartDate { get; set; }

    [Column("enrollment_end_date")]
    public DateTime? EnrollmentEndDate { get; set; }

    [Column("underwriting_data")]
    public string? UnderwritingData { get; set; } // JSONB

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [Column("deleted_at")]
    public DateTime? DeletedAt { get; set; }

    [Column("cancellation_approvals")]
    public string? CancellationApprovals { get; set; } // JSONB

    // Navigation properties
    public ICollection<PolicyNomineeEntity> Nominees { get; set; } = new List<PolicyNomineeEntity>();
    public ICollection<PolicyRiderEntity> Riders { get; set; } = new List<PolicyRiderEntity>();
    public ICollection<ClaimEntity> Claims { get; set; } = new List<ClaimEntity>();
}
