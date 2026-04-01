using Google.Protobuf.WellKnownTypes;
using Insuretech.Endorsement.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Endorsement.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Domain;
using EndorsementEntity = Insuretech.Endorsement.Entity.V1.Endorsement;

namespace PoliSync.Endorsement.Application.Commands;

public sealed class RequestEndorsementCommandHandler
    : IRequestHandler<RequestEndorsementCommand, Result<RequestEndorsementResult>>
{
    private readonly IEndorsementDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly IMediator _mediator;
    private readonly ILogger<RequestEndorsementCommandHandler> _logger;

    public RequestEndorsementCommandHandler(
        IEndorsementDataGateway dataGateway,
        IEventBus eventBus,
        IMediator mediator,
        ILogger<RequestEndorsementCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _mediator = mediator;
        _logger = logger;
    }

    public async Task<Result<RequestEndorsementResult>> Handle(
        RequestEndorsementCommand request,
        CancellationToken cancellationToken)
    {
        try
        {
            if (string.IsNullOrWhiteSpace(request.PolicyId))
                return Result.Fail<RequestEndorsementResult>("VALIDATION_ERROR", "PolicyId is required");

            var endorsementType = ParseEndorsementType(request.Type);
            if (endorsementType == EndorsementType.Unspecified)
                return Result.Fail<RequestEndorsementResult>("VALIDATION_ERROR", $"Invalid endorsement type: {request.Type}");

            var now = DateTime.UtcNow;
            var effectiveDate = DateTime.TryParse(request.EffectiveDate, out var parsed)
                ? Timestamp.FromDateTime(DateTime.SpecifyKind(parsed, DateTimeKind.Utc))
                : Timestamp.FromDateTime(now.Date.AddDays(1).ToUniversalTime());

            var premiumDelta = endorsementType switch
            {
                EndorsementType.SumAssuredChange  => 15_000L,
                EndorsementType.PremiumAdjustment => 10_000L,
                EndorsementType.RiderAddition     => 5_000L,
                EndorsementType.RiderRemoval      => -5_000L,
                _                                 => 0L
            };

            var endorsement = new EndorsementEntity
            {
                Id = Guid.NewGuid().ToString("N"),
                EndorsementNumber = $"END-{now:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}",
                PolicyId = request.PolicyId,
                Type = endorsementType,
                Reason = request.Reason,
                Changes = request.Changes,
                PremiumAdjustment = new Insuretech.Common.V1.Money { Amount = premiumDelta, Currency = "BDT" },
                PremiumRefundRequired = endorsementType == EndorsementType.RiderRemoval,
                Status = EndorsementStatus.Pending,
                RequestedBy = "SYSTEM",
                EffectiveDate = effectiveDate
            };

            var created = await _dataGateway.CreateEndorsementAsync(endorsement, cancellationToken);

            await _eventBus.PublishAsync(
                new EndorsementRequestedEvent(created.Id, request.PolicyId, endorsementType.ToString()),
                cancellationToken);

            _logger.LogInformation("Endorsement requested: {EndorsementId} for policy {PolicyId}",
                created.Id, request.PolicyId);

            // Trigger approval workflow — major vs standard resolved by sub-type
            var workflowResult = await _mediator.Send(new TriggerWorkflowCommand(
                new WorkflowTriggerContext
                {
                    EntityType  = "ENDORSEMENT",
                    EntityId    = created.Id,
                    InitiatedBy = "SYSTEM",
                    SubType     = endorsementType.ToString(),
                    Portal      = "B2C",
                    Metadata    = new Dictionary<string, string>
                    {
                        ["policy_id"]          = request.PolicyId,
                        ["endorsement_number"] = created.EndorsementNumber,
                        ["endorsement_type"]   = endorsementType.ToString()
                    }
                }), cancellationToken);

            if (workflowResult.IsSuccess && workflowResult.Value!.WasTriggered)
                _logger.LogInformation(
                    "Endorsement approval workflow started: instance={InstanceId} template='{Template}'",
                    workflowResult.Value.WorkflowInstanceId, workflowResult.Value.TemplateName);

            return Result.Ok(new RequestEndorsementResult(created.Id, created.EndorsementNumber));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to request endorsement for policy {PolicyId}", request.PolicyId);
            return Result.Fail<RequestEndorsementResult>("REQUEST_ENDORSEMENT_FAILED", ex.Message);
        }
    }

    private static EndorsementType ParseEndorsementType(string value)
    {
        if (string.IsNullOrWhiteSpace(value)) return EndorsementType.Unspecified;
        if (System.Enum.TryParse<EndorsementType>(value, true, out var direct)) return direct;
        var candidate = string.Concat(value.Split('_', StringSplitOptions.RemoveEmptyEntries)
            .Select(s => char.ToUpperInvariant(s[0]) + s[1..].ToLowerInvariant()));
        return System.Enum.TryParse<EndorsementType>(candidate, true, out var parsed) ? parsed : EndorsementType.Unspecified;
    }
}

public sealed record EndorsementRequestedEvent(string EndorsementId, string PolicyId, string Type)
    : PoliSync.SharedKernel.Domain.DomainEvent;

