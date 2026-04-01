using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Commission.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Commission.Application.Commands;

public sealed class CreateCommissionPayoutCommandHandler
    : IRequestHandler<CreateCommissionPayoutCommand, Result<CreateCommissionPayoutResult>>
{
    private readonly ICommissionDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<CreateCommissionPayoutCommandHandler> _logger;

    public CreateCommissionPayoutCommandHandler(
        ICommissionDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<CreateCommissionPayoutCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<CreateCommissionPayoutResult>> Handle(
        CreateCommissionPayoutCommand request, CancellationToken cancellationToken)
    {
        try
        {
            if (string.IsNullOrWhiteSpace(request.RecipientId))
                return Result.Fail<CreateCommissionPayoutResult>("VALIDATION_ERROR", "RecipientId is required");

            var response = await _dataGateway.CreatePayoutAsync(
                request.RecipientType,
                request.RecipientId,
                request.PeriodStart,
                request.PeriodEnd,
                request.CommissionIds,
                cancellationToken);

            if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
                return Result.Fail<CreateCommissionPayoutResult>(response.Error.Code, response.Error.Message);

            await _eventBus.PublishAsync(
                new CommissionPayoutCreatedEvent(response.PayoutId, request.RecipientId, response.CommissionCount),
                cancellationToken);

            _logger.LogInformation("Commission payout created: {PayoutId} for recipient {RecipientId}",
                response.PayoutId, request.RecipientId);

            return Result.Ok(new CreateCommissionPayoutResult(
                response.PayoutId,
                response.PayoutNumber,
                response.TotalAmount!,
                response.CommissionCount));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create commission payout for recipient {RecipientId}", request.RecipientId);
            return Result.Fail<CreateCommissionPayoutResult>("CREATE_PAYOUT_FAILED", ex.Message);
        }
    }
}

public sealed record CommissionPayoutCreatedEvent(string PayoutId, string RecipientId, int CommissionCount)
    : DomainEvent;
