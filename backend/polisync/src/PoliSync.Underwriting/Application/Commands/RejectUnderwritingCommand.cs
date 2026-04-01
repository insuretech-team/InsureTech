using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Underwriting.Application.Commands;

public sealed record RejectUnderwritingCommand(
    string QuoteId,
    string UnderwriterId,
    string Reason,
    string Comments,
    string RiskLevel
) : ICommand<string>; // Returns decision_id
