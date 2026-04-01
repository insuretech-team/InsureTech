using MediatR;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Handler for setting service fee on a quotation
/// </summary>
public sealed class SetServiceFeeCommandHandler : IRequestHandler<SetServiceFeeCommand, Result>
{
    private readonly IQuotationDataGateway _gateway;

    public SetServiceFeeCommandHandler(IQuotationDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result> Handle(SetServiceFeeCommand request, CancellationToken cancellationToken)
    {
        var quotation = await _gateway.GetByIdAsync(request.QuotationId, cancellationToken);
        if (quotation == null)
            return Result.Fail("NOT_FOUND", "Quotation not found");

        var result = quotation.SetServiceFee(request.ServiceFee);
        if (!result.IsSuccess)
            return result;

        await _gateway.UpdateAsync(quotation, cancellationToken);

        return Result.Ok();
    }
}
