using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Command to submit a quotation for underwriting
/// </summary>
public sealed record SubmitQuotationCommand(Guid QuotationId) : IRequest<Result>;
