using MediatR;
using PoliSync.Orders.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Handler for initiating payment for an order
/// Integrates with Go payment-service via gRPC
/// </summary>
public sealed class InitiatePaymentCommandHandler : IRequestHandler<InitiatePaymentCommand, Result<InitiatePaymentResponse>>
{
    private readonly IOrderDataGateway _orderGateway;

    public InitiatePaymentCommandHandler(IOrderDataGateway orderGateway)
    {
        _orderGateway = orderGateway;
    }

    public async Task<Result<InitiatePaymentResponse>> Handle(InitiatePaymentCommand request, CancellationToken cancellationToken)
    {
        var response = await _orderGateway.InitiatePaymentAsync(
            request.OrderId.ToString(),
            request.PaymentMethod,
            request.CallbackUrl ?? string.Empty,
            request.IdempotencyKey ?? Guid.NewGuid().ToString(),
            cancellationToken);

        if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
        {
            return Result.Fail<InitiatePaymentResponse>(response.Error.Code, response.Error.Message);
        }

        return Result.Ok(new InitiatePaymentResponse(
            response.PaymentId,
            response.PaymentUrl,
            response.PaymentUrl,
            response.ExpiresAt.ToDateTime()));
    }
}
