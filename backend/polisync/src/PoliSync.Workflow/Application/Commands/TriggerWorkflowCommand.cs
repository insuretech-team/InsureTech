using MediatR;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Workflow.Domain;

namespace PoliSync.Workflow.Application.Commands;

/// <summary>
/// Triggers a workflow for an entity by resolving the appropriate template
/// from IWorkflowTemplateProvider based on the entity context.
///
/// This is the primary entry point for domain handlers (Claims, Endorsements, Refunds)
/// to start workflows without knowing template names.
///
/// Usage in a command handler:
/// <code>
///   await _mediator.Send(new TriggerWorkflowCommand(new WorkflowTriggerContext
///   {
///       EntityType = "CLAIM",
///       EntityId = claimId,
///       InitiatedBy = userId,
///       AmountPaisa = claimedAmountPaisa
///   }), cancellationToken);
/// </code>
/// </summary>
public sealed record TriggerWorkflowCommand(WorkflowTriggerContext Context)
    : IRequest<Result<TriggerWorkflowResult>>;

public sealed record TriggerWorkflowResult(
    string WorkflowInstanceId,
    string TemplateName,
    string EntityType,
    string EntityId,
    bool WasTriggered // false = no template found, workflow skipped
);
