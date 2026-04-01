using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Command to cancel an order
/// </summary>
public sealed record CancelOrderCommand(
    Guid OrderId,
    string Reason) : IRequest<Result>;
