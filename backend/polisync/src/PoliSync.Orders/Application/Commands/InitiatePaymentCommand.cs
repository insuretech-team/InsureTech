using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Command to initiate payment for an order
/// </summary>
public sealed record InitiatePaymentCommand(
    Guid OrderId,
    string PaymentMethod,
    string? CallbackUrl = null,
    string? IdempotencyKey = null
) : IRequest<Result<InitiatePaymentResponse>>;

public sealed record InitiatePaymentResponse(
    string PaymentId,
    string PaymentUrl,
    string PaymentGatewayRef,
    DateTime ExpiresAt
);
