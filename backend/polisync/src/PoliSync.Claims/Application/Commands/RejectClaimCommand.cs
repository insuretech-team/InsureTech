using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Claims.Application.Commands;

public sealed record RejectClaimCommand(
    string ClaimId,
    string ApproverId,
    string Reason
) : ICommand;
