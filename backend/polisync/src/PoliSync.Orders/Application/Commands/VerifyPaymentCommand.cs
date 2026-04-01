using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Command to verify payment status from webhook callback
/// </summary>
public sealed record VerifyPaymentCommand(
    Guid OrderId,
    string PaymentId,
    string TransactionId,
    string Status,
    decimal Amount,
    string Currency) : IRequest<Result<bool>>;
