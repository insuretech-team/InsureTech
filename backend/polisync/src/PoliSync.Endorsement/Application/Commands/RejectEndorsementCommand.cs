using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Endorsement.Application.Commands;

public sealed record RejectEndorsementCommand(
    string EndorsementId,
    string RejectedBy,
    string Reason
) : ICommand;
