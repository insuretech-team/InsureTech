using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Command to mark a quotation as received by underwriting
/// </summary>
public sealed record MarkQuotationReceivedCommand(
    Guid QuotationId
) : IRequest<Result>;
