using System;
using System.Collections.Generic;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Premium quote generated during underwriting.
/// Maps to 'quotes' table in insurance_schema.
/// </summary>
public class Quote
{
    public Guid Id { get; set; }
    public string QuoteNumber { get; set; } = string.Empty;
    public Guid BeneficiaryId { get; set; }
    public Guid InsurerProductId { get; set; }
    public QuoteStatus Status { get; set; }

    // Money fields stored as bigint (paisa)
    public long SumAssuredAmount { get; set; }
    public string SumAssuredCurrency { get; set; } = "BDT";
    
    public int TermYears { get; set; }
    public string PremiumPaymentMode { get; set; } = "YEARLY";

    public long BasePremiumAmount { get; set; }
    public long RiderPremiumAmount { get; set; }
    public long TaxAmount { get; set; }
    public long TotalPremiumAmount { get; set; }
    public string Currency { get; set; } = "BDT";

    // JSONB fields
    public string? PremiumCalculationJson { get; set; }
    public string? SelectedRidersJson { get; set; }

    // Applicant data snapshot
    public int ApplicantAge { get; set; }
    public string? ApplicantOccupation { get; set; }
    public bool IsSmoker { get; set; }

    public DateTime ValidUntil { get; set; }
    public Guid? ConvertedPolicyId { get; set; }
    public DateTime? ConvertedAt { get; set; }

    // Audit Info JSONB
    public string? AuditInfoJson { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }
    public bool IsDeleted { get; set; }

    // Convenience accessors
    public Money SumAssured => new(SumAssuredAmount, SumAssuredCurrency);
    public Money BasePremium => new(BasePremiumAmount, Currency);
    public Money TotalPremium => new(TotalPremiumAmount, Currency);
}
