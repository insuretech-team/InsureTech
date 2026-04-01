using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Claims.Application.Commands;

public sealed record ApproveClaimCommand(
    string ClaimId,
    string ApproverId,
    long ApprovedAmountPaisa,
    string Notes
) : ICommand;
