using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Command to apply underwriting loading to a quotation
/// </summary>
public sealed record ApplyLoadingCommand(Guid QuotationId, long LoadingAmount) : IRequest<Result>;
