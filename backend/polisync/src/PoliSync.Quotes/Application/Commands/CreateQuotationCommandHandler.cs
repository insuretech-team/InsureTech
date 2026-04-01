using MediatR;
using PoliSync.Quotes.Domain;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Handler for creating a new quotation
/// </summary>
public sealed class CreateQuotationCommandHandler : IRequestHandler<CreateQuotationCommand, Result<Guid>>
{
    private readonly IQuotationDataGateway _gateway;

    public CreateQuotationCommandHandler(IQuotationDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result<Guid>> Handle(CreateQuotationCommand request, CancellationToken cancellationToken)
    {
        // Create quotation aggregate
        var quotationResult = Quotation.Create(
            request.TenantId,
            request.ProductId,
            request.PlanId,
            request.CustomerId,
            request.BasePremium,
            request.RiderPremium,
            request.ExpiryDays);

        if (!quotationResult.IsSuccess)
            return Result.Fail<Guid>(quotationResult.Error!.Code, quotationResult.Error.Message);

        var quotation = quotationResult.Value!;

        // Persist via data gateway
        await _gateway.CreateAsync(quotation, cancellationToken);

        return Result.Ok(quotation.Id);
    }
}
