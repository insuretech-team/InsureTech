using MediatR;
using PoliSync.SharedKernel.CQRS;
using Insuretech.Workflow.Entity.V1;

namespace PoliSync.Workflow.Application.Queries;

/// <summary>Gets a workflow instance with its tasks.</summary>
public sealed record GetWorkflowInstanceQuery(string WorkflowInstanceId)
    : IRequest<Result<GetWorkflowInstanceResult>>;

public sealed record GetWorkflowInstanceResult(
    WorkflowInstance Instance,
    IReadOnlyList<WorkflowTask> Tasks
);

/// <summary>Gets all workflow instances for a given entity (history).</summary>
public sealed record GetWorkflowHistoryQuery(string EntityType, string EntityId)
    : IRequest<Result<IReadOnlyList<WorkflowInstance>>>;

/// <summary>Gets all tasks assigned to the current user.</summary>
public sealed record GetMyWorkflowTasksQuery(
    string UserId,
    string? Status = null,
    int Page = 1,
    int PageSize = 20
) : IRequest<Result<GetMyTasksResult>>;

public sealed record GetMyTasksResult(
    IReadOnlyList<WorkflowTask> Tasks,
    int TotalCount
);
