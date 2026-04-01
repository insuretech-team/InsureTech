using MediatR;
using PoliSync.Quotes.Domain;
using PoliSync.Quotes.Infrastructure;

namespace PoliSync.Quotes.Application.Queries;

/// <summary>
/// Handler for retrieving a quotation by ID
/// </summary>
public sealed class GetQuotationByIdQueryHandler : IRequestHandler<GetQuotationByIdQuery, Quotation?>
{
    private readonly IQuotationDataGateway _gateway;

    public GetQuotationByIdQueryHandler(IQuotationDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Quotation?> Handle(GetQuotationByIdQuery request, CancellationToken cancellationToken)
    {
        return await _gateway.GetByIdAsync(request.QuotationId, cancellationToken);
    }
}
