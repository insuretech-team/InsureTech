using MediatR;
using PoliSync.Orders.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Handler for canceling an order
/// </summary>
public sealed class CancelOrderCommandHandler : IRequestHandler<CancelOrderCommand, Result>
{
    private readonly IOrderDataGateway _gateway;

    public CancelOrderCommandHandler(IOrderDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result> Handle(CancelOrderCommand request, CancellationToken cancellationToken)
    {
        var response = await _gateway.CancelOrderAsync(
            request.OrderId.ToString(),
            request.Reason,
            cancellationToken);

        if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
        {
            return Result.Fail(response.Error.Code, response.Error.Message);
        }

        return Result.Ok();
    }
}
