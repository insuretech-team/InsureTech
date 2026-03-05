using PoliSync.SharedKernel.Domain;

namespace PoliSync.Underwriting.Domain;

public class UnderwritingDecision : Entity
{
    public Guid DecisionId { get; private set; }
    public Guid QuoteId { get; private set; } // FK to Quote
    
    public DecisionType Decision { get; private set; }
    public DecisionMethod Method { get; private set; }
    public float RiskScore { get; private set; }
    public RiskLevel RiskLevel { get; private set; }
    
    public string? Reason { get; private set; }
    public string? Conditions { get; private set; } // JSON array of strings
    public string? RiskFactors { get; private set; } // JSON map
    
    public long? AdjustedPremiumAmount { get; private set; }
    public Guid? UnderwriterId { get; private set; }
    
    public DateTime DecidedAt { get; private set; }

    private UnderwritingDecision() { }

    public static UnderwritingDecision Create(
        Guid quoteId,
        DecisionType decision,
        DecisionMethod method,
        float riskScore,
        RiskLevel riskLevel,
        string? reason = null,
        string? conditions = null,
        string? riskFactors = null,
        long? adjustedPremiumAmount = null,
        Guid? underwriterId = null)
    {
        var entity = new UnderwritingDecision
        {
            DecisionId = Guid.NewGuid(),
            QuoteId = quoteId,
            Decision = decision,
            Method = method,
            RiskScore = riskScore,
            RiskLevel = riskLevel,
            Reason = reason,
            Conditions = conditions ?? "[]",
            RiskFactors = riskFactors ?? "{}",
            AdjustedPremiumAmount = adjustedPremiumAmount,
            UnderwriterId = underwriterId,
            DecidedAt = DateTime.UtcNow
        };

        entity.RaiseDomainEvent(new UnderwritingDecisionMadeEvent(entity.DecisionId, quoteId, decision));
        return entity;
    }
}
