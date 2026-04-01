using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Refund.Application.Commands;

public sealed record ApproveRefundCommand(
    string RefundId,
    string ApprovedBy,
    string Comments
) : ICommand;
