using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Refund.Application.Commands;

public sealed record ProcessRefundCommand(
    string RefundId,
    string PaymentReference,
    string PaymentMethod,
    long RefundAmountPaisa,
    string Reason,
    string InitiatedBy
) : ICommand<ProcessRefundResult>;

public sealed record ProcessRefundResult(string PaymentRefundId, string ProcessedAt);
