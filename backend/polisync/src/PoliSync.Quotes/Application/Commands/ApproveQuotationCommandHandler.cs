using MediatR;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Handler for approving a quotation
/// </summary>
public sealed class ApproveQuotationCommandHandler : IRequestHandler<ApproveQuotationCommand, Result>
{
    private readonly IQuotationDataGateway _gateway;
    private readonly IEventBus _eventBus;

    public ApproveQuotationCommandHandler(IQuotationDataGateway gateway, IEventBus eventBus)
    {
        _gateway = gateway;
        _eventBus = eventBus;
    }

    public async Task<Result> Handle(ApproveQuotationCommand request, CancellationToken cancellationToken)
    {
        // Retrieve quotation
        var quotation = await _gateway.GetByIdAsync(request.QuotationId, cancellationToken);
        if (quotation == null)
            return Result.Fail("NOT_FOUND", "Quotation not found");

        // Approve quotation
        var approveResult = quotation.Approve();
        if (!approveResult.IsSuccess)
            return approveResult;

        // Update via data gateway
        await _gateway.UpdateAsync(quotation, cancellationToken);

        // Publish domain events to Kafka
        foreach (var domainEvent in quotation.DomainEvents)
        {
            await _eventBus.PublishAsync(domainEvent, cancellationToken);
        }

        return Result.Ok();
    }
}
