using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Command to reject a quotation with a reason
/// </summary>
public sealed record RejectQuotationCommand(Guid QuotationId, string Reason) : IRequest<Result>;
