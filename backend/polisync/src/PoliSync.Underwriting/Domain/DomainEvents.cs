using PoliSync.SharedKernel.Domain;

namespace PoliSync.Underwriting.Domain;

public record QuoteRequestedEvent(Guid QuoteId, string QuoteNumber) : DomainEvent
{
    public override string EventType => "Quote.Requested";
}

public record QuoteApprovedEvent(Guid QuoteId) : DomainEvent
{
    public override string EventType => "Quote.Approved";
}

public record QuoteRejectedEvent(Guid QuoteId) : DomainEvent
{
    public override string EventType => "Quote.Rejected";
}

public record UnderwritingDecisionMadeEvent(Guid DecisionId, Guid QuoteId, DecisionType Decision) : DomainEvent
{
    public override string EventType => "UnderwritingDecision.Made";
}
