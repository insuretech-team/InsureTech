using Google.Protobuf.WellKnownTypes;
using Insuretech.Endorsement.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Endorsement.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Endorsement.Application.Commands;

public sealed class ApproveEndorsementCommandHandler : IRequestHandler<ApproveEndorsementCommand, Result>
{
    private readonly IEndorsementDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<ApproveEndorsementCommandHandler> _logger;

    public ApproveEndorsementCommandHandler(
        IEndorsementDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<ApproveEndorsementCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result> Handle(ApproveEndorsementCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var endorsement = await _dataGateway.GetEndorsementAsync(request.EndorsementId, cancellationToken);
            if (endorsement is null)
                return Result.Fail("ENDORSEMENT_NOT_FOUND", $"Endorsement {request.EndorsementId} not found");

            if (endorsement.Status != EndorsementStatus.Pending)
                return Result.Fail("INVALID_STATUS", $"Endorsement in status {endorsement.Status} cannot be approved");

            endorsement.Status = EndorsementStatus.Applied;
            endorsement.ApprovedBy = request.ApprovedBy;
            endorsement.ApprovedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            if (!string.IsNullOrWhiteSpace(request.Comments))
                endorsement.Changes = $"{endorsement.Changes};comments={request.Comments}";

            await _dataGateway.UpdateEndorsementAsync(endorsement, cancellationToken);
            await _eventBus.PublishAsync(
                new EndorsementApprovedEvent(endorsement.Id, endorsement.PolicyId),
                cancellationToken);

            _logger.LogInformation("Endorsement {EndorsementId} approved by {ApprovedBy}",
                request.EndorsementId, request.ApprovedBy);
            return Result.Ok();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve endorsement {EndorsementId}", request.EndorsementId);
            return Result.Fail("APPROVE_ENDORSEMENT_FAILED", ex.Message);
        }
    }
}

public sealed record EndorsementApprovedEvent(string EndorsementId, string PolicyId)
    : PoliSync.SharedKernel.Domain.DomainEvent;

