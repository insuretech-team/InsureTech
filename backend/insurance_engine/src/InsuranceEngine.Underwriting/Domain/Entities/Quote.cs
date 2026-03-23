using System;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Underwriting.Domain.Entities;

public class Quote
{
    public Guid Id { get; set; }
    public string QuoteNumber { get; set; } = string.Empty;
    public Guid BeneficiaryId { get; set; }
    public Guid InsurerProductId { get; set; }
    public QuoteStatus Status { get; set; }

    public long SumAssuredAmount { get; set; }
    public string SumAssuredCurrency { get; set; } = "BDT";
    
    public int TermYears { get; set; }
    public string PremiumPaymentMode { get; set; } = "YEARLY";

    public long BasePremiumAmount { get; set; }
    public string BasePremiumCurrency { get; set; } = "BDT";
    public long RiderPremiumAmount { get; set; }
    public string RiderPremiumCurrency { get; set; } = "BDT";
    public long TaxAmount { get; set; }
    public string TaxCurrency { get; set; } = "BDT";
    public long TotalPremiumAmount { get; set; }
    public string TotalPremiumCurrency { get; set; } = "BDT";

    public string? PremiumCalculationJson { get; set; }
    public string? SelectedRidersJson { get; set; }

    public int ApplicantAge { get; set; }
    public string? ApplicantOccupation { get; set; }
    public bool IsSmoker { get; set; }

    public DateTime ValidUntil { get; set; }
    public Guid? ConvertedPolicyId { get; set; }
    public DateTime? ConvertedAt { get; set; }

    public string? AuditInfoJson { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }

    public Money SumAssured => new(SumAssuredAmount, SumAssuredCurrency);
    public Money BasePremium => new(BasePremiumAmount, BasePremiumCurrency);
    public Money RiderPremium => new(RiderPremiumAmount, RiderPremiumCurrency);
    public Money Tax => new(TaxAmount, TaxCurrency);
    public Money TotalPremium => new(TotalPremiumAmount, TotalPremiumCurrency);
}
