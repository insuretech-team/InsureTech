using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Endorsement.Application.Commands;

public sealed record ApproveEndorsementCommand(
    string EndorsementId,
    string ApprovedBy,
    string Comments
) : ICommand;
