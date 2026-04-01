using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Command to apply a discount to a quotation
/// </summary>
public sealed record ApplyDiscountCommand(
    Guid QuotationId,
    long DiscountAmount
) : IRequest<Result>;
