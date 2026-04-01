using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Commands;

/// <summary>
/// Command to create an order from an approved quotation
/// </summary>
public sealed record CreateOrderCommand(
    Guid QuotationId,
    Guid CustomerId,
    Guid ProductId,
    Guid PlanId,
    long TotalPayable,
    string Currency = "BDT",
    DateTime? PaymentDueAt = null,
    DateTime? CoverageStartAt = null,
    DateTime? CoverageEndAt = null
) : IRequest<Result<Guid>>;
