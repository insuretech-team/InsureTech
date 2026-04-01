using Insuretech.Endorsement.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Endorsement.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Endorsement.Application.Commands;

public sealed class RejectEndorsementCommandHandler : IRequestHandler<RejectEndorsementCommand, Result>
{
    private readonly IEndorsementDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<RejectEndorsementCommandHandler> _logger;

    public RejectEndorsementCommandHandler(
        IEndorsementDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<RejectEndorsementCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result> Handle(RejectEndorsementCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var endorsement = await _dataGateway.GetEndorsementAsync(request.EndorsementId, cancellationToken);
            if (endorsement is null)
                return Result.Fail("ENDORSEMENT_NOT_FOUND", $"Endorsement {request.EndorsementId} not found");

            if (endorsement.Status != EndorsementStatus.Pending)
                return Result.Fail("INVALID_STATUS", $"Endorsement in status {endorsement.Status} cannot be rejected");

            endorsement.Status = EndorsementStatus.Rejected;
            endorsement.Reason = request.Reason;

            await _dataGateway.UpdateEndorsementAsync(endorsement, cancellationToken);
            await _eventBus.PublishAsync(
                new EndorsementRejectedEvent(endorsement.Id, endorsement.PolicyId, request.Reason),
                cancellationToken);

            _logger.LogInformation("Endorsement {EndorsementId} rejected. Reason: {Reason}",
                request.EndorsementId, request.Reason);
            return Result.Ok();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to reject endorsement {EndorsementId}", request.EndorsementId);
            return Result.Fail("REJECT_ENDORSEMENT_FAILED", ex.Message);
        }
    }
}

public sealed record EndorsementRejectedEvent(string EndorsementId, string PolicyId, string Reason)
    : PoliSync.SharedKernel.Domain.DomainEvent;

