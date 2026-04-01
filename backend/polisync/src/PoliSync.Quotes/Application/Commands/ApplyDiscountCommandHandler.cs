using MediatR;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Handler for applying discount to a quotation
/// </summary>
public sealed class ApplyDiscountCommandHandler : IRequestHandler<ApplyDiscountCommand, Result>
{
    private readonly IQuotationDataGateway _gateway;

    public ApplyDiscountCommandHandler(IQuotationDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result> Handle(ApplyDiscountCommand request, CancellationToken cancellationToken)
    {
        var quotation = await _gateway.GetByIdAsync(request.QuotationId, cancellationToken);
        if (quotation == null)
            return Result.Fail("NOT_FOUND", "Quotation not found");

        var result = quotation.ApplyDiscount(request.DiscountAmount);
        if (!result.IsSuccess)
            return result;

        await _gateway.UpdateAsync(quotation, cancellationToken);

        return Result.Ok();
    }
}
