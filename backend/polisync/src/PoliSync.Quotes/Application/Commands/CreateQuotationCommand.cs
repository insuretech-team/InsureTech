using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Command to create a new quotation
/// </summary>
public sealed record CreateQuotationCommand(
    Guid TenantId,
    Guid ProductId,
    Guid PlanId,
    Guid CustomerId,
    long BasePremium,
    long RiderPremium,
    int ExpiryDays = 30) : IRequest<Result<Guid>>;
