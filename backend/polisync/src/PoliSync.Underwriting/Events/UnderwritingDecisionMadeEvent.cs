using PoliSync.SharedKernel.Domain;

namespace PoliSync.Underwriting.Events;

public sealed record UnderwritingDecisionMadeEvent : DomainEvent
{
    public string QuoteId { get; init; } = string.Empty;
    public string DecisionId { get; init; } = string.Empty;
    public string Decision { get; init; } = string.Empty;
    public string RiskLevel { get; init; } = string.Empty;
    public bool PremiumAdjusted { get; init; }
    public long QuotedAmount { get; init; }
    public string Currency { get; init; } = "BDT";
    public string Reason { get; init; } = string.Empty;
}
