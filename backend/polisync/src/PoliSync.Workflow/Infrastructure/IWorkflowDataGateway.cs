using Google.Protobuf.WellKnownTypes;
using Insuretech.Workflow.Entity.V1;
using PoliSync.Workflow.Application.Queries;

namespace PoliSync.Workflow.Infrastructure;

/// <summary>
/// Gateway interface for the Go workflow-engine gRPC service.
/// Abstracts all communication with the workflow microservice.
/// </summary>
public interface IWorkflowDataGateway
{
    // ── Definitions ──────────────────────────────────────────────────────────

    /// <summary>Creates a new workflow definition and returns its ID.</summary>
    Task<string?> CreateDefinitionAsync(
        string name,
        string description,
        string workflowType,
        string entityType,
        string stepsJson,
        string conditionsJson,
        CancellationToken cancellationToken = default);

    /// <summary>Resolves a definition ID by template name. Returns null if not found.</summary>
    Task<string?> ResolveDefinitionIdByNameAsync(
        string name,
        CancellationToken cancellationToken = default);

    // ── Instances ─────────────────────────────────────────────────────────────

    /// <summary>Starts a workflow instance and returns the instance ID.</summary>
    Task<string?> StartWorkflowAsync(
        string definitionId,
        string entityType,
        string entityId,
        Struct? context,
        CancellationToken cancellationToken = default);

    /// <summary>Gets a workflow instance with its tasks.</summary>
    Task<GetWorkflowInstanceResult?> GetWorkflowInstanceAsync(
        string instanceId,
        CancellationToken cancellationToken = default);

    /// <summary>Gets all workflow instances for an entity.</summary>
    Task<IReadOnlyList<WorkflowInstance>> GetWorkflowHistoryAsync(
        string entityType,
        string entityId,
        CancellationToken cancellationToken = default);

    // ── Tasks ─────────────────────────────────────────────────────────────────

    /// <summary>Gets tasks assigned to a user.</summary>
    Task<GetMyTasksResult> GetMyTasksAsync(
        string userId,
        string? status,
        int page,
        int pageSize,
        CancellationToken cancellationToken = default);

    /// <summary>Completes a workflow task.</summary>
    Task<CompleteTaskResult> CompleteTaskAsync(
        string taskId,
        string decision,
        string comments,
        string completedBy,
        CancellationToken cancellationToken = default);
}

/// <summary>Result of completing a workflow task.</summary>
public sealed record CompleteTaskResult(
    bool Success,
    string? WorkflowInstanceId = null,
    bool WorkflowCompleted = false,
    string? ErrorCode = null,
    string? ErrorMessage = null
);
