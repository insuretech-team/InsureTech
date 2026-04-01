using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Command to set the service fee for a quotation
/// </summary>
public sealed record SetServiceFeeCommand(
    Guid QuotationId,
    long ServiceFee
) : IRequest<Result>;
