using Google.Protobuf.WellKnownTypes;
using Insuretech.Claims.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Claims.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Claims.Application.Commands;

public sealed class SettleClaimCommandHandler : IRequestHandler<SettleClaimCommand, Result<string>>
{
    private readonly IClaimDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<SettleClaimCommandHandler> _logger;

    public SettleClaimCommandHandler(
        IClaimDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<SettleClaimCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(SettleClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim is null)
                return Result.Fail<string>("CLAIM_NOT_FOUND", $"Claim {request.ClaimId} not found");

            if (claim.Status != ClaimStatus.Approved)
                return Result.Fail<string>("INVALID_STATUS", "Only approved claims can be settled");

            var now = Timestamp.FromDateTime(DateTime.UtcNow);
            var settledAmount = claim.ApprovedAmount?.Amount > 0 ? claim.ApprovedAmount : claim.ClaimedAmount;
            var paymentId = $"PAY-{Guid.NewGuid():N}"[..18];

            claim.Status = ClaimStatus.Settled;
            claim.SettledAmount = settledAmount;
            claim.SettledAt = now;
            claim.UpdatedAt = now;

            await _dataGateway.UpdateClaimAsync(claim, cancellationToken);
            await _eventBus.PublishAsync(new ClaimSettledEvent(claim.ClaimId, claim.PolicyId, settledAmount?.Amount ?? 0, paymentId), cancellationToken);

            _logger.LogInformation("Claim {ClaimId} settled. PaymentId: {PaymentId}", request.ClaimId, paymentId);
            return Result.Ok(paymentId);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to settle claim {ClaimId}", request.ClaimId);
            return Result.Fail<string>("SETTLE_CLAIM_FAILED", ex.Message);
        }
    }
}

public sealed record ClaimSettledEvent(string ClaimId, string PolicyId, long SettledAmountPaisa, string PaymentId) : PoliSync.SharedKernel.Domain.DomainEvent;

