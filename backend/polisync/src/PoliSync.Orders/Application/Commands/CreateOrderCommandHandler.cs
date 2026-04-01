using MediatR;
using PoliSync.Orders.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Handler for creating an order from an approved quotation
/// </summary>
public sealed class CreateOrderCommandHandler : IRequestHandler<CreateOrderCommand, Result<Guid>>
{
    private readonly IOrderDataGateway _gateway;

    public CreateOrderCommandHandler(IOrderDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result<Guid>> Handle(CreateOrderCommand request, CancellationToken cancellationToken)
    {
        var response = await _gateway.CreateOrderAsync(
            request.QuotationId.ToString(),
            request.CustomerId.ToString(),
            "SSLCommerz",
            cancellationToken,
            productId: request.ProductId == Guid.Empty ? null : request.ProductId.ToString(),
            planId: request.PlanId == Guid.Empty ? null : request.PlanId.ToString(),
            totalPayable: request.TotalPayable,
            currency: request.Currency);

        if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
        {
            return Result.Fail<Guid>(response.Error.Code, response.Error.Message);
        }

        var orderId = response.Order?.Order?.OrderId ?? string.Empty;
        if (!Guid.TryParse(orderId, out var parsedOrderId))
        {
            return Result.Fail<Guid>("INVALID_RESPONSE", "Order service returned an invalid order identifier");
        }

        return Result.Ok(parsedOrderId);
    }
}
