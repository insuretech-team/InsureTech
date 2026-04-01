using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Workflow.Infrastructure;
using Google.Protobuf.WellKnownTypes;

namespace PoliSync.Workflow.Application.Commands;

/// <summary>
/// Handles StartWorkflowCommand by:
/// 1. Resolving the workflow definition ID from the named template via the Go engine
/// 2. Starting a workflow instance for the entity
/// 3. Returning the workflow instance ID for tracking
/// </summary>
public sealed class StartWorkflowCommandHandler
    : IRequestHandler<StartWorkflowCommand, Result<StartWorkflowResult>>
{
    private readonly IWorkflowDataGateway _gateway;
    private readonly ILogger<StartWorkflowCommandHandler> _logger;

    public StartWorkflowCommandHandler(
        IWorkflowDataGateway gateway,
        ILogger<StartWorkflowCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<StartWorkflowResult>> Handle(
        StartWorkflowCommand request,
        CancellationToken cancellationToken)
    {
        // Resolve definition ID by name
        var definitionId = await _gateway.ResolveDefinitionIdByNameAsync(
            request.TemplateName, cancellationToken);

        if (definitionId is null)
            return Result.Fail<StartWorkflowResult>(
                "WORKFLOW_DEFINITION_NOT_FOUND",
                $"No active workflow definition found with name '{request.TemplateName}'");

        // Build context struct for Go engine
        var contextFields = new Dictionary<string, Value>();
        contextFields["initiated_by"] = Value.ForString(request.InitiatedBy);
        contextFields["entity_type"] = Value.ForString(request.EntityType);
        if (request.Context is not null)
        {
            foreach (var kv in request.Context)
                contextFields[kv.Key] = Value.ForString(kv.Value);
        }
        var contextStruct = new Struct();
        foreach (var kv in contextFields)
            contextStruct.Fields[kv.Key] = kv.Value;

        var instanceId = await _gateway.StartWorkflowAsync(
            definitionId,
            request.EntityType,
            request.EntityId,
            contextStruct,
            cancellationToken);

        if (instanceId is null)
            return Result.Fail<StartWorkflowResult>(
                "WORKFLOW_START_FAILED",
                $"Failed to start workflow '{request.TemplateName}' for {request.EntityType}/{request.EntityId}");

        _logger.LogInformation(
            "Started workflow instance {InstanceId} for {EntityType}/{EntityId} using template {Template}",
            instanceId, request.EntityType, request.EntityId, request.TemplateName);

        return Result.Ok(new StartWorkflowResult(instanceId, request.EntityType, request.EntityId));
    }
}
