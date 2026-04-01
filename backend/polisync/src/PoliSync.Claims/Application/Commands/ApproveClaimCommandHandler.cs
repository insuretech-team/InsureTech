using Google.Protobuf.WellKnownTypes;
using Insuretech.Claims.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Claims.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Claims.Application.Commands;

public sealed class ApproveClaimCommandHandler : IRequestHandler<ApproveClaimCommand, Result>
{
    private readonly IClaimDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<ApproveClaimCommandHandler> _logger;

    public ApproveClaimCommandHandler(
        IClaimDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<ApproveClaimCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result> Handle(ApproveClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim is null)
                return Result.Fail("CLAIM_NOT_FOUND", $"Claim {request.ClaimId} not found");

            if (claim.Status != ClaimStatus.Submitted && claim.Status != ClaimStatus.UnderReview)
                return Result.Fail("INVALID_STATUS", $"Claim in status {claim.Status} cannot be approved");

            var now = Timestamp.FromDateTime(DateTime.UtcNow);
            var approvedAmount = request.ApprovedAmountPaisa > 0
                ? new Insuretech.Common.V1.Money { Amount = request.ApprovedAmountPaisa, Currency = "BDT" }
                : claim.ClaimedAmount;

            claim.Status = ClaimStatus.Approved;
            claim.ApprovedAmount = approvedAmount;
            claim.ApprovedAt = now;
            claim.UpdatedAt = now;
            var approval = new ClaimApproval
            {
                ApprovalId = Guid.NewGuid().ToString("N"),
                ClaimId = claim.ClaimId,
                ApproverId = request.ApproverId,
                ApproverRole = "L1",
                ApprovalLevel = 1,
                Decision = ApprovalDecision.Approved,
                ApprovedAmount = approvedAmount,
                Notes = request.Notes,
                ApprovedAt = now,
                CreatedAt = now,
                ApprovedCurrency = approvedAmount.Currency
            };

            await _dataGateway.UpdateClaimAsync(claim, cancellationToken);
            await _dataGateway.CreateClaimApprovalAsync(approval, cancellationToken);
            await _eventBus.PublishAsync(new ClaimApprovedEvent(claim.ClaimId, claim.PolicyId, approvedAmount.Amount), cancellationToken);

            _logger.LogInformation("Claim {ClaimId} approved by {ApproverId}", request.ClaimId, request.ApproverId);
            return Result.Ok();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve claim {ClaimId}", request.ClaimId);
            return Result.Fail("APPROVE_CLAIM_FAILED", ex.Message);
        }
    }
}

public sealed record ClaimApprovedEvent(string ClaimId, string PolicyId, long ApprovedAmountPaisa) : PoliSync.SharedKernel.Domain.DomainEvent;
