using MediatR;
using PoliSync.Orders.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Handler for confirming payment completion
/// </summary>
public sealed class ConfirmPaymentCommandHandler : IRequestHandler<ConfirmPaymentCommand, Result>
{
    private readonly IOrderDataGateway _gateway;

    public ConfirmPaymentCommandHandler(IOrderDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result> Handle(ConfirmPaymentCommand request, CancellationToken cancellationToken)
    {
        var order = await _gateway.GetOrderAsync(request.OrderId.ToString(), cancellationToken);
        if (order?.Order is null)
            return Result.Fail("NOT_FOUND", "Order not found");

        if (!string.IsNullOrWhiteSpace(order.Order.PaymentId) && order.Order.PaymentId != request.PaymentId)
            return Result.Fail("PAYMENT_MISMATCH", "Payment ID does not match order");

        var response = await _gateway.ConfirmPaymentAsync(
            request.OrderId.ToString(),
            request.PaymentId,
            request.TransactionId,
            cancellationToken);

        if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
        {
            return Result.Fail(response.Error.Code, response.Error.Message);
        }

        return Result.Ok();
    }
}
