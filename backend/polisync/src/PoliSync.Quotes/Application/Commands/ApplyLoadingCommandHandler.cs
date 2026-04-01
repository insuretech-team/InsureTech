using MediatR;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Handler for applying underwriting loading to a quotation
/// </summary>
public sealed class ApplyLoadingCommandHandler : IRequestHandler<ApplyLoadingCommand, Result>
{
    private readonly IQuotationDataGateway _gateway;

    public ApplyLoadingCommandHandler(IQuotationDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result> Handle(ApplyLoadingCommand request, CancellationToken cancellationToken)
    {
        // Retrieve quotation
        var quotation = await _gateway.GetByIdAsync(request.QuotationId, cancellationToken);
        if (quotation == null)
            return Result.Fail("NOT_FOUND", "Quotation not found");

        // Apply loading
        var loadingResult = quotation.ApplyLoading(request.LoadingAmount);
        if (!loadingResult.IsSuccess)
            return loadingResult;

        // Update via data gateway
        await _gateway.UpdateAsync(quotation, cancellationToken);

        return Result.Ok();
    }
}
