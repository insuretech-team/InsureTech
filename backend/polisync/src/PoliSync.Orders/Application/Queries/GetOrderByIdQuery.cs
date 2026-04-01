using MediatR;

namespace PoliSync.Orders.Application.Queries;

/// <summary>
/// Query to get an order by ID
/// </summary>
public sealed record GetOrderByIdQuery(Guid OrderId) : IRequest<OrderDto?>;
