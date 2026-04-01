using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Commission.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Commission.Application.Commands;

public sealed class CalculateCommissionCommandHandler
    : IRequestHandler<CalculateCommissionCommand, Result<CalculateCommissionResult>>
{
    private readonly ICommissionDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<CalculateCommissionCommandHandler> _logger;

    public CalculateCommissionCommandHandler(
        ICommissionDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<CalculateCommissionCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<CalculateCommissionResult>> Handle(
        CalculateCommissionCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var response = await _dataGateway.CalculateCommissionAsync(
                request.PolicyId,
                request.CommissionType,
                request.RecipientType,
                request.RecipientId,
                cancellationToken);

            if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
                return Result.Fail<CalculateCommissionResult>(response.Error.Code, response.Error.Message);

            await _eventBus.PublishAsync(
                new CommissionCalculatedEvent(response.CommissionId, request.PolicyId, response.Amount?.Amount ?? 0),
                cancellationToken);

            _logger.LogInformation("Commission calculated: {CommissionId} for policy {PolicyId}",
                response.CommissionId, request.PolicyId);

            return Result.Ok(new CalculateCommissionResult(
                response.CommissionId,
                response.CommissionNumber,
                response.Amount!,
                response.CalculationBreakdown));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to calculate commission for policy {PolicyId}", request.PolicyId);
            return Result.Fail<CalculateCommissionResult>("CALCULATE_COMMISSION_FAILED", ex.Message);
        }
    }
}

public sealed record CommissionCalculatedEvent(string CommissionId, string PolicyId, long AmountPaisa)
    : DomainEvent;
