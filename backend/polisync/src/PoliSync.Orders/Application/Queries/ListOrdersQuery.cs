using MediatR;
using PoliSync.Orders.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Queries;

/// <summary>
/// Query to list orders with filtering
/// </summary>
public sealed record ListOrdersQuery(
    Guid? CustomerId = null,
    OrderStatus? Status = null,
    int PageNumber = 1,
    int PageSize = 20) : IRequest<Result<ListOrdersResponse>>;

/// <summary>
/// Response containing list of orders
/// </summary>
public sealed record ListOrdersResponse(
    IReadOnlyList<OrderDto> Orders,
    int TotalCount,
    int PageNumber,
    int PageSize);

/// <summary>
/// Order data transfer object
/// </summary>
public sealed record OrderDto(
    Guid Id,
    string OrderNumber,
    Guid QuotationId,
    Guid CustomerId,
    Guid ProductId,
    Guid PlanId,
    OrderStatus Status,
    long TotalPayable,
    string Currency,
    string? PaymentId,
    string? PaymentGatewayRef,
    OrderPaymentStatus PaymentStatus,
    string? PolicyId,
    string? CancellationReason,
    string? FailureReason,
    DateTime? PaymentDueAt,
    DateTime? CoverageStartAt,
    DateTime? CoverageEndAt,
    DateTime CreatedAt,
    DateTime? UpdatedAt);
