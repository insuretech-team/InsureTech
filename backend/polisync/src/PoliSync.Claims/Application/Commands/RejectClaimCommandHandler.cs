using Google.Protobuf.WellKnownTypes;
using Insuretech.Claims.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Claims.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Claims.Application.Commands;

public sealed class RejectClaimCommandHandler : IRequestHandler<RejectClaimCommand, Result>
{
    private readonly IClaimDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<RejectClaimCommandHandler> _logger;

    public RejectClaimCommandHandler(
        IClaimDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<RejectClaimCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result> Handle(RejectClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim is null)
                return Result.Fail("CLAIM_NOT_FOUND", $"Claim {request.ClaimId} not found");

            if (claim.Status == ClaimStatus.Settled || claim.Status == ClaimStatus.Rejected)
                return Result.Fail("INVALID_STATUS", $"Claim in status {claim.Status} cannot be rejected");

            var now = Timestamp.FromDateTime(DateTime.UtcNow);
            claim.Status = ClaimStatus.Rejected;
            claim.RejectionReason = request.Reason;
            claim.UpdatedAt = now;
            var approval = new ClaimApproval
            {
                ApprovalId = Guid.NewGuid().ToString("N"),
                ClaimId = claim.ClaimId,
                ApproverId = request.ApproverId,
                ApproverRole = "L1",
                ApprovalLevel = 1,
                Decision = ApprovalDecision.Rejected,
                ApprovedAmount = new Insuretech.Common.V1.Money { Amount = 0, Currency = "BDT" },
                Notes = request.Reason,
                ApprovedAt = now,
                CreatedAt = now,
                ApprovedCurrency = "BDT"
            };

            await _dataGateway.UpdateClaimAsync(claim, cancellationToken);
            await _dataGateway.CreateClaimApprovalAsync(approval, cancellationToken);
            await _eventBus.PublishAsync(new ClaimRejectedEvent(claim.ClaimId, claim.PolicyId, request.Reason), cancellationToken);

            _logger.LogInformation("Claim {ClaimId} rejected by {ApproverId}. Reason: {Reason}", request.ClaimId, request.ApproverId, request.Reason);
            return Result.Ok();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to reject claim {ClaimId}", request.ClaimId);
            return Result.Fail("REJECT_CLAIM_FAILED", ex.Message);
        }
    }
}

public sealed record ClaimRejectedEvent(string ClaimId, string PolicyId, string Reason) : PoliSync.SharedKernel.Domain.DomainEvent;
