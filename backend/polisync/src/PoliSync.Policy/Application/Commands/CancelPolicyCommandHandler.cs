using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Policy.Domain;
using PoliSync.Policy.Infrastructure;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Domain;

namespace PoliSync.Policy.Application.Commands;

public sealed class CancelPolicyCommandHandler : IRequestHandler<CancelPolicyCommand, Result>
{
    private readonly IPolicyDataGateway _policyDataGateway;
    private readonly IEventBus _eventBus;
    private readonly IMediator _mediator;
    private readonly ILogger<CancelPolicyCommandHandler> _logger;

    public CancelPolicyCommandHandler(
        IPolicyDataGateway policyDataGateway,
        IEventBus eventBus,
        IMediator mediator,
        ILogger<CancelPolicyCommandHandler> logger)
    {
        _policyDataGateway = policyDataGateway;
        _eventBus = eventBus;
        _mediator = mediator;
        _logger = logger;
    }

    public async Task<Result> Handle(CancelPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policy = await _policyDataGateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policy is null)
                return Result.Fail("POLICY_NOT_FOUND", $"Policy {request.PolicyId} not found");

            if (policy.Status == Insuretech.Policy.Entity.V1.PolicyStatus.Cancelled)
                return Result.Fail("ALREADY_CANCELLED", "Policy is already cancelled");

            if (policy.Status == Insuretech.Policy.Entity.V1.PolicyStatus.Expired)
                return Result.Fail("POLICY_EXPIRED", "Expired policies cannot be cancelled");

            var aggregate = new PolicyAggregate(policy);
            aggregate.CancelPolicy(request.Reason);

            await _policyDataGateway.UpdatePolicyAsync(aggregate.Policy, cancellationToken);

            foreach (var domainEvent in aggregate.DomainEvents)
                await _eventBus.PublishAsync(domainEvent, cancellationToken);

            _logger.LogInformation("Policy {PolicyId} cancelled. Reason: {Reason}", request.PolicyId, request.Reason);

            // Trigger policy cancellation approval workflow
            var workflowResult = await _mediator.Send(new TriggerWorkflowCommand(
                new WorkflowTriggerContext
                {
                    EntityType  = "POLICY",
                    EntityId    = request.PolicyId,
                    InitiatedBy = request.RequestedBy ?? "SYSTEM",
                    SubType     = "CANCELLATION",
                    Portal      = request.Portal ?? "B2C",
                    Metadata    = new Dictionary<string, string>
                    {
                        ["policy_id"] = request.PolicyId,
                        ["reason"]    = request.Reason ?? string.Empty
                    }
                }), cancellationToken);

            if (workflowResult.IsSuccess && workflowResult.Value!.WasTriggered)
                _logger.LogInformation(
                    "Policy cancellation approval workflow started: instance={InstanceId} template='{Template}'",
                    workflowResult.Value.WorkflowInstanceId, workflowResult.Value.TemplateName);

            return Result.Ok();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to cancel policy {PolicyId}", request.PolicyId);
            return Result.Fail("CANCEL_POLICY_FAILED", ex.Message);
        }
    }
}
