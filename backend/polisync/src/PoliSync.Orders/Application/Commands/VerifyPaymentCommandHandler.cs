using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Orders.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Handler for verifying payment from webhook callback
/// </summary>
public sealed class VerifyPaymentCommandHandler : IRequestHandler<VerifyPaymentCommand, Result<bool>>
{
    private readonly IOrderDataGateway _orderGateway;
    private readonly ILogger<VerifyPaymentCommandHandler> _logger;

    public VerifyPaymentCommandHandler(
        IOrderDataGateway orderGateway,
        ILogger<VerifyPaymentCommandHandler> logger)
    {
        _orderGateway = orderGateway;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(VerifyPaymentCommand request, CancellationToken cancellationToken)
    {
        var order = await _orderGateway.GetOrderAsync(request.OrderId.ToString(), cancellationToken);
        if (order?.Order is null)
        {
            _logger.LogWarning("Order {OrderId} not found for payment verification", request.OrderId);
            return Result.Fail<bool>("ORDER_NOT_FOUND", "Order not found");
        }

        // Verify payment status
        if (request.Status.Equals("VALID", StringComparison.OrdinalIgnoreCase) ||
            request.Status.Equals("VALIDATED", StringComparison.OrdinalIgnoreCase))
        {
            var confirmResult = await _orderGateway.ConfirmPaymentAsync(
                request.OrderId.ToString(),
                request.PaymentId,
                request.TransactionId,
                cancellationToken);
            if (confirmResult.Error is not null && !string.IsNullOrWhiteSpace(confirmResult.Error.Code))
            {
                _logger.LogError("Failed to confirm payment for order {OrderId}: {Error}", 
                    request.OrderId, confirmResult.Error.Message);
                return Result.Fail<bool>(confirmResult.Error.Code, confirmResult.Error.Message);
            }

            _logger.LogInformation("Payment verified and confirmed for order {OrderId}", request.OrderId);
            return Result.Ok(true);
        }
        else
        {
            var cancelResult = await _orderGateway.CancelOrderAsync(
                request.OrderId.ToString(),
                $"Payment verification failed with status: {request.Status}",
                cancellationToken);
            if (cancelResult.Error is not null && !string.IsNullOrWhiteSpace(cancelResult.Error.Code))
            {
                _logger.LogError("Failed to mark order as failed for order {OrderId}: {Error}", 
                    request.OrderId, cancelResult.Error.Message);
                return Result.Fail<bool>(cancelResult.Error.Code, cancelResult.Error.Message);
            }

            _logger.LogWarning("Payment failed for order {OrderId} with status {Status}", 
                request.OrderId, request.Status);
            return Result.Ok(false);
        }
    }
}
