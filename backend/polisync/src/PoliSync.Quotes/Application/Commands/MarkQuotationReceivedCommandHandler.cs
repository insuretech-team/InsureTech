using MediatR;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Handler for marking a quotation as received by underwriting
/// </summary>
public sealed class MarkQuotationReceivedCommandHandler : IRequestHandler<MarkQuotationReceivedCommand, Result>
{
    private readonly IQuotationDataGateway _gateway;

    public MarkQuotationReceivedCommandHandler(IQuotationDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result> Handle(MarkQuotationReceivedCommand request, CancellationToken cancellationToken)
    {
        var quotation = await _gateway.GetByIdAsync(request.QuotationId, cancellationToken);
        if (quotation == null)
            return Result.Fail("NOT_FOUND", "Quotation not found");

        var result = quotation.MarkAsReceived();
        if (!result.IsSuccess)
            return result;

        await _gateway.UpdateAsync(quotation, cancellationToken);

        return Result.Ok();
    }
}
