using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Refund.Application.Commands;

public sealed record RequestRefundCommand(
    string PolicyId,
    string Reason,
    string ReasonDetails,
    string RequestedBy,
    long AmountPaisa = 0  // used for workflow template routing (high vs standard)
) : ICommand<RequestRefundResult>;

public sealed record RequestRefundResult(string RefundId, string RefundNumber);
