using MediatR;
using PoliSync.Orders.Infrastructure;

namespace PoliSync.Orders.Application.Queries;

/// <summary>
/// Handler for retrieving an order by ID
/// </summary>
public sealed class GetOrderByIdQueryHandler : IRequestHandler<GetOrderByIdQuery, OrderDto?>
{
    private readonly IOrderDataGateway _gateway;

    public GetOrderByIdQueryHandler(IOrderDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<OrderDto?> Handle(GetOrderByIdQuery request, CancellationToken cancellationToken)
    {
        var order = await _gateway.GetOrderAsync(request.OrderId.ToString(), cancellationToken);
        return order is null ? null : OrderViewMapper.ToDto(order);
    }
}
