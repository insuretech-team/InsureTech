using MediatR;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Quotes.Application.Commands;

/// <summary>
/// Handler for rejecting a quotation
/// </summary>
public sealed class RejectQuotationCommandHandler : IRequestHandler<RejectQuotationCommand, Result>
{
    private readonly IQuotationDataGateway _gateway;
    private readonly IEventBus _eventBus;

    public RejectQuotationCommandHandler(IQuotationDataGateway gateway, IEventBus eventBus)
    {
        _gateway = gateway;
        _eventBus = eventBus;
    }

    public async Task<Result> Handle(RejectQuotationCommand request, CancellationToken cancellationToken)
    {
        // Retrieve quotation
        var quotation = await _gateway.GetByIdAsync(request.QuotationId, cancellationToken);
        if (quotation == null)
            return Result.Fail("NOT_FOUND", "Quotation not found");

        // Reject quotation
        var rejectResult = quotation.Reject(request.Reason);
        if (!rejectResult.IsSuccess)
            return rejectResult;

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
