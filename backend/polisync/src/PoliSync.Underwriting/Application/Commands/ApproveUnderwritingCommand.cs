using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Underwriting.Application.Commands;

public sealed record ApproveUnderwritingCommand(
    string QuoteId,
    string UnderwriterId,
    string Comments,
    string? ConditionsJson,
    bool PremiumAdjusted,
    long AdjustedPremiumPaisa,
    string RiskLevel
) : ICommand<string>; // Returns decision_id
