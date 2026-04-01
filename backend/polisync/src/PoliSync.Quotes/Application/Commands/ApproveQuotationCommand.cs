using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Command to approve a quotation after underwriting
/// </summary>
public sealed record ApproveQuotationCommand(Guid QuotationId) : IRequest<Result>;
