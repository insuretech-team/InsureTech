using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Endorsement.Application.Commands;

public sealed record RequestEndorsementCommand(
    string PolicyId,
    string Type,
    string Reason,
    string Changes,
    string EffectiveDate
) : ICommand<RequestEndorsementResult>;

public sealed record RequestEndorsementResult(string EndorsementId, string EndorsementNumber);
