using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Command to confirm payment completion for an order
/// </summary>
public sealed record ConfirmPaymentCommand(
    Guid OrderId,
    string PaymentId,
    string TransactionId
) : IRequest<Result>;
