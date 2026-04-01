using MediatR;
using PoliSync.Quotes.Domain;

namespace PoliSync.Quotes.Application.Queries;

/// <summary>
/// Query to retrieve a quotation by ID
/// </summary>
public sealed record GetQuotationByIdQuery(Guid QuotationId) : IRequest<Quotation?>;
