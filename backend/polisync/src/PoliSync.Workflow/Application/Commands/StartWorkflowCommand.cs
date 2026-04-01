using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Workflow.Application.Commands;

/// <summary>
/// Starts a workflow instance for a given entity using a named template.
/// The template is resolved by name from the Go workflow-engine's definition registry.
/// </summary>
public sealed record StartWorkflowCommand(
    string TemplateName,    // matches WorkflowDefinition.Name in Go service
    string EntityType,      // e.g. "CLAIM", "ENDORSEMENT", "REFUND", "QUOTATION"
    string EntityId,        // UUID of the entity (claim_id, endorsement_id, etc.)
    string InitiatedBy,     // user UUID who triggered this
    Dictionary<string, string>? Context = null  // extra metadata passed to Go engine
) : IRequest<Result<StartWorkflowResult>>;

public sealed record StartWorkflowResult(
    string WorkflowInstanceId,
    string EntityType,
    string EntityId,
    string Status = "IN_PROGRESS"
);
