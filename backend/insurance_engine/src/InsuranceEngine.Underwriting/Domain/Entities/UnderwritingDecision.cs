using System;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Underwriting.Domain.Entities;

public class UnderwritingDecision
{
    public Guid Id { get; set; }
    public Guid QuoteId { get; set; }
    
    public DecisionType Decision { get; set; }
    public DecisionMethod Method { get; set; }

    public decimal RiskScore { get; set; }
    public RiskLevel RiskLevel { get; set; }
    public string? RiskFactorsJson { get; set; }

    public string? Reason { get; set; }
    public string? ConditionsJson { get; set; }

    public bool IsPremiumAdjusted { get; set; }
    public long AdjustedPremiumAmount { get; set; }
    public string AdjustedPremiumCurrency { get; set; } = "BDT";
    public string? AdjustmentReason { get; set; }

    public Guid? UnderwriterId { get; set; }
    public string? UnderwriterComments { get; set; }

    public DateTime DecidedAt { get; set; }
    public DateTime? ValidUntil { get; set; }

    public string? AuditInfoJson { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    public Money AdjustedPremium => new(AdjustedPremiumAmount, AdjustedPremiumCurrency);
}
