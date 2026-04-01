using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Domain;

namespace PoliSync.ApiHost.Controllers;

/// <summary>
/// Admin REST API for managing dynamic workflow templates at runtime.
/// Allows adding, removing, and listing templates without restarting the service.
/// Templates registered here are immediately available for new workflow triggers.
///
/// Note: Runtime-registered templates survive until the next restart unless
/// persisted via the Go workflow-engine (RegisterWorkflowTemplateCommand does this).
/// </summary>
[ApiController]
[Route("api/admin/workflow-templates")]
public sealed class WorkflowTemplateAdminController : ControllerBase
{
    private readonly IWorkflowTemplateProvider _templateProvider;
    private readonly IMediator _mediator;

    public WorkflowTemplateAdminController(
        IWorkflowTemplateProvider templateProvider,
        IMediator mediator)
    {
        _templateProvider = templateProvider;
        _mediator = mediator;
    }

    /// <summary>List all currently registered workflow templates.</summary>
    [HttpGet]
    public IActionResult ListTemplates()
    {
        var templates = _templateProvider.GetAllTemplates()
            .Select(t => new
            {
                t.Name,
                t.EntityType,
                t.WorkflowType,
                t.Description,
                StepCount = t.Steps.Count,
                Steps = t.Steps.Select(s => new
                {
                    s.Name,
                    s.Type,
                    s.AssignRole,
                    s.DueHours,
                    s.Order
                })
            });
        return Ok(templates);
    }

    /// <summary>
    /// Register a new template dynamically and persist it to the Go workflow-engine.
    /// If a template with the same name exists, it is overwritten.
    /// </summary>
    [HttpPost]
    public async Task<IActionResult> RegisterTemplate(
        [FromBody] RegisterTemplateRequest req,
        CancellationToken ct)
    {
        if (string.IsNullOrWhiteSpace(req.Name) || string.IsNullOrWhiteSpace(req.EntityType))
            return BadRequest(new { error = "Name and EntityType are required" });

        var template = new WorkflowTemplate
        {
            Name = req.Name,
            EntityType = req.EntityType.ToUpperInvariant(),
            WorkflowType = req.WorkflowType ?? "APPROVAL",
            Description = req.Description ?? string.Empty,
            Steps = req.Steps?.Select(s => new WorkflowStepTemplate
            {
                Name = s.Name,
                Type = s.Type ?? "APPROVAL",
                AssignRole = s.AssignRole ?? string.Empty,
                AssignTo = s.AssignTo ?? string.Empty,
                DueHours = s.DueHours > 0 ? s.DueHours : 72,
                Order = s.Order
            }).ToList() ?? [],
            Conditions = new WorkflowConditions
            {
                FailFastOnRejection = req.Conditions?.FailFastOnRejection ?? true,
                RequireAllApprovals = req.Conditions?.RequireAllApprovals ?? true,
                AutoApproveAfterHours = req.Conditions?.AutoApproveAfterHours,
                EscalateToRole = req.Conditions?.EscalateToRole
            }
        };

        // Register in-memory (immediately available)
        _templateProvider.Register(template);

        // Persist to Go workflow-engine (survives restarts)
        var result = await _mediator.Send(new RegisterWorkflowTemplateCommand(template), ct);

        if (result.IsFailure)
            return BadRequest(new { error = result.Error?.Code, message = result.Error?.Message });

        return Ok(new
        {
            result.Value!.DefinitionId,
            result.Value.Name,
            result.Value.WasCreated,
            Message = result.Value.WasCreated
                ? "Template registered and persisted to workflow-engine"
                : "Template already exists in workflow-engine — in-memory copy updated"
        });
    }

    /// <summary>Remove a template from the in-memory registry (does not delete from Go engine).</summary>
    [HttpDelete("{name}")]
    public IActionResult RemoveTemplate(string name)
    {
        var removed = _templateProvider.Remove(name);
        return removed
            ? Ok(new { message = $"Template '{name}' removed from in-memory registry" })
            : NotFound(new { message = $"Template '{name}' not found in registry" });
    }

    /// <summary>
    /// Test-trigger a workflow for a given entity without a real entity existing.
    /// Useful for validating template routing logic (dev/staging only).
    /// </summary>
    [HttpPost("test-trigger")]
    public async Task<IActionResult> TestTrigger(
        [FromBody] TestTriggerRequest req,
        CancellationToken ct)
    {
        var result = await _mediator.Send(new TriggerWorkflowCommand(
            new WorkflowTriggerContext
            {
                EntityType  = req.EntityType.ToUpperInvariant(),
                EntityId    = req.EntityId ?? Guid.NewGuid().ToString(),
                InitiatedBy = req.InitiatedBy ?? "admin-test",
                AmountPaisa = req.AmountPaisa,
                SubType     = req.SubType,
                Portal      = req.Portal ?? "SYSTEM"
            }), ct);

        return result.IsSuccess
            ? Ok(result.Value)
            : BadRequest(new { error = result.Error?.Code, message = result.Error?.Message });
    }
}

// ── Request DTOs ─────────────────────────────────────────────────────────────

public sealed record RegisterTemplateRequest(
    string Name,
    string EntityType,
    string? WorkflowType = "APPROVAL",
    string? Description = null,
    List<StepRequest>? Steps = null,
    ConditionsRequest? Conditions = null
);

public sealed record StepRequest(
    string Name,
    string? Type = "APPROVAL",
    string? AssignRole = null,
    string? AssignTo = null,
    int DueHours = 72,
    int Order = 1
);

public sealed record ConditionsRequest(
    bool FailFastOnRejection = true,
    bool RequireAllApprovals = true,
    int? AutoApproveAfterHours = null,
    string? EscalateToRole = null
);

public sealed record TestTriggerRequest(
    string EntityType,
    string? EntityId = null,
    string? InitiatedBy = null,
    long AmountPaisa = 0,
    string? SubType = null,
    string? Portal = "SYSTEM"
);
