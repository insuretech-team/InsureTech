using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Workflow.Infrastructure;

namespace PoliSync.Workflow.Application.Commands;

/// <summary>
/// Registers a dynamic workflow template into the Go workflow engine.
/// If a definition with the same name already exists, returns the existing ID (idempotent).
/// </summary>
public sealed class RegisterWorkflowTemplateCommandHandler
    : IRequestHandler<RegisterWorkflowTemplateCommand, Result<RegisterWorkflowTemplateResult>>
{
    private readonly IWorkflowDataGateway _gateway;
    private readonly ILogger<RegisterWorkflowTemplateCommandHandler> _logger;

    public RegisterWorkflowTemplateCommandHandler(
        IWorkflowDataGateway gateway,
        ILogger<RegisterWorkflowTemplateCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<RegisterWorkflowTemplateResult>> Handle(
        RegisterWorkflowTemplateCommand request,
        CancellationToken cancellationToken)
    {
        var template = request.Template;

        // Check if already exists — idempotent
        var existingId = await _gateway.ResolveDefinitionIdByNameAsync(
            template.Name, cancellationToken);

        if (existingId is not null)
        {
            _logger.LogDebug(
                "Workflow template '{Name}' already registered as {Id}",
                template.Name, existingId);

            return Result.Ok(new RegisterWorkflowTemplateResult(
                existingId, template.Name, WasCreated: false));
        }

        // Create new definition in Go engine
        var definitionId = await _gateway.CreateDefinitionAsync(
            template.Name,
            template.Description,
            template.WorkflowType,
            template.EntityType,
            template.SerializeSteps(),
            template.SerializeConditions(),
            cancellationToken);

        if (definitionId is null)
            return Result.Fail<RegisterWorkflowTemplateResult>(
                "REGISTRATION_FAILED",
                $"Failed to register workflow template '{template.Name}'");

        _logger.LogInformation(
            "Registered workflow template '{Name}' as {Id} (entity={EntityType}, steps={StepCount})",
            template.Name, definitionId, template.EntityType, template.Steps.Count);

        return Result.Ok(new RegisterWorkflowTemplateResult(
            definitionId, template.Name, WasCreated: true));
    }
}
