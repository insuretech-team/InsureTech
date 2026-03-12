using System;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Underwriting decision for quote approval/rejection.
/// Maps to 'underwriting_decisions' table in insurance_schema.
/// </summary>
public class UnderwritingDecision
{
    public Guid Id { get; set; }
    public Guid QuoteId { get; set; }
    
    public DecisionType Decision { get; set; }
    public DecisionMethod Method { get; set; }

    // Risk assessment
    public decimal RiskScore { get; set; }
    public RiskLevel RiskLevel { get; set; }
    public string? RiskFactorsJson { get; set; } // JSONB

    // Decision details
    public string? Reason { get; set; }
    public string? ConditionsJson { get; set; } // JSONB

    // Premium adjustment
    public bool IsPremiumAdjusted { get; set; }
    public long AdjustedPremiumAmount { get; set; }
    public string AdjustedPremiumCurrency { get; set; } = "BDT";
    public string? AdjustmentReason { get; set; }

    // Underwriter info
    public Guid? UnderwriterId { get; set; }
    public string? UnderwriterComments { get; set; }

    public DateTime DecidedAt { get; set; }
    public DateTime? ValidUntil { get; set; }

    // Audit Info JSONB
    public string? AuditInfoJson { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    public Money AdjustedPremium => new(AdjustedPremiumAmount, AdjustedPremiumCurrency);
}
