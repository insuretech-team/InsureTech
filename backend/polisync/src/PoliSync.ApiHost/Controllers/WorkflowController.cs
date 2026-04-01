using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Application.Queries;

namespace PoliSync.ApiHost.Controllers;

/// <summary>
/// REST API for the workflow engine.
/// Exposes task management and workflow history to portals.
/// </summary>
[ApiController]
[Route("api/workflows")]
public sealed class WorkflowController : ControllerBase
{
    private readonly IMediator _mediator;

    public WorkflowController(IMediator mediator) => _mediator = mediator;

    private string UserId => HttpContext.Request.Headers["x-user-id"].FirstOrDefault() ?? string.Empty;

    // ── Instance ──────────────────────────────────────────────────────────────

    /// <summary>Start a workflow for an entity by template name.</summary>
    [HttpPost("start")]
    public async Task<IActionResult> Start([FromBody] StartWorkflowRequest req, CancellationToken ct)
    {
        var result = await _mediator.Send(new StartWorkflowCommand(
            req.TemplateName, req.EntityType, req.EntityId,
            InitiatedBy: UserId,
            Context: req.Context), ct);

        return result.IsSuccess
            ? Ok(result.Value)
            : BadRequest(new { error = result.Error?.Code, message = result.Error?.Message });
    }

    /// <summary>Get a workflow instance with its tasks.</summary>
    [HttpGet("instances/{instanceId}")]
    public async Task<IActionResult> GetInstance(string instanceId, CancellationToken ct)
    {
        var result = await _mediator.Send(new GetWorkflowInstanceQuery(instanceId), ct);
        return result.IsSuccess ? Ok(result.Value) : NotFound(result.Error?.Message);
    }

    /// <summary>Get all workflow instances for an entity (audit trail).</summary>
    [HttpGet("{entityType}/{entityId}/history")]
    public async Task<IActionResult> GetHistory(string entityType, string entityId, CancellationToken ct)
    {
        var result = await _mediator.Send(new GetWorkflowHistoryQuery(entityType, entityId), ct);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error?.Message);
    }

    // ── Tasks ─────────────────────────────────────────────────────────────────

    /// <summary>Get tasks assigned to the current user.</summary>
    [HttpGet("my-tasks")]
    public async Task<IActionResult> GetMyTasks(
        [FromQuery] string? status,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20,
        CancellationToken ct = default)
    {
        var result = await _mediator.Send(
            new GetMyWorkflowTasksQuery(UserId, status, page, pageSize), ct);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error?.Message);
    }

    /// <summary>Complete a workflow task with APPROVED, REJECTED, or RETURNED decision.</summary>
    [HttpPost("tasks/{taskId}/complete")]
    public async Task<IActionResult> CompleteTask(
        string taskId,
        [FromBody] CompleteTaskRequest req,
        CancellationToken ct)
    {
        var result = await _mediator.Send(
            new CompleteWorkflowTaskCommand(taskId, req.Decision, req.Comments, UserId), ct);

        return result.IsSuccess
            ? Ok(result.Value)
            : BadRequest(new { error = result.Error?.Code, message = result.Error?.Message });
    }
}

// ── Request DTOs ─────────────────────────────────────────────────────────────

public sealed record StartWorkflowRequest(
    string TemplateName,
    string EntityType,
    string EntityId,
    Dictionary<string, string>? Context = null
);

public sealed record CompleteTaskRequest(
    string Decision,
    string Comments = ""
);
