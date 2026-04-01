using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Commission.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Commission.Application.Commands;

public sealed class ProcessCommissionPayoutCommandHandler
    : IRequestHandler<ProcessCommissionPayoutCommand, Result<string>>
{
    private readonly ICommissionDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<ProcessCommissionPayoutCommandHandler> _logger;

    public ProcessCommissionPayoutCommandHandler(
        ICommissionDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<ProcessCommissionPayoutCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(
        ProcessCommissionPayoutCommand request, CancellationToken cancellationToken)
    {
        try
        {
            if (string.IsNullOrWhiteSpace(request.PayoutId))
                return Result.Fail<string>("VALIDATION_ERROR", "PayoutId is required");

            var response = await _dataGateway.ProcessPayoutAsync(
                request.PayoutId,
                request.PaymentMethod,
                request.PaymentReference,
                cancellationToken);

            if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
                return Result.Fail<string>(response.Error.Code, response.Error.Message);

            var paidAt = string.IsNullOrWhiteSpace(response.PaidAt)
                ? DateTime.UtcNow.ToString("O")
                : response.PaidAt;

            await _eventBus.PublishAsync(
                new CommissionPayoutProcessedEvent(request.PayoutId, request.PaymentReference),
                cancellationToken);

            _logger.LogInformation("Commission payout processed: {PayoutId}", request.PayoutId);
            return Result.Ok(paidAt);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to process commission payout {PayoutId}", request.PayoutId);
            return Result.Fail<string>("PROCESS_PAYOUT_FAILED", ex.Message);
        }
    }
}

public sealed record CommissionPayoutProcessedEvent(string PayoutId, string PaymentReference)
    : DomainEvent;
