using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Workflow.Application.Commands;

/// <summary>
/// Completes a workflow task with a decision (APPROVED, REJECTED, RETURNED).
/// Triggers instance advancement in the Go workflow engine.
/// </summary>
public sealed record CompleteWorkflowTaskCommand(
    string TaskId,
    string Decision,   // APPROVED | REJECTED | RETURNED
    string Comments,
    string CompletedBy // user UUID
) : IRequest<Result<CompleteWorkflowTaskResult>>;

public sealed record CompleteWorkflowTaskResult(
    string TaskId,
    string Decision,
    string WorkflowInstanceId,
    bool WorkflowCompleted
);
