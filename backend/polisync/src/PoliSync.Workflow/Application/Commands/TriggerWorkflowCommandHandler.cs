using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Workflow.Domain;
using PoliSync.Workflow.Infrastructure;
using Google.Protobuf.WellKnownTypes;

namespace PoliSync.Workflow.Application.Commands;

/// <summary>
/// Handles TriggerWorkflowCommand:
///   1. Resolves the right template using IWorkflowTemplateProvider (context-aware routing)
///   2. Ensures that template is registered in the Go engine (idempotent)
///   3. Starts a workflow instance for the entity
/// </summary>
public sealed class TriggerWorkflowCommandHandler
    : IRequestHandler<TriggerWorkflowCommand, Result<TriggerWorkflowResult>>
{
    private readonly IWorkflowTemplateProvider _templateProvider;
    private readonly IWorkflowDataGateway _gateway;
    private readonly IMediator _mediator;
    private readonly ILogger<TriggerWorkflowCommandHandler> _logger;

    public TriggerWorkflowCommandHandler(
        IWorkflowTemplateProvider templateProvider,
        IWorkflowDataGateway gateway,
        IMediator mediator,
        ILogger<TriggerWorkflowCommandHandler> logger)
    {
        _templateProvider = templateProvider;
        _gateway = gateway;
        _mediator = mediator;
        _logger = logger;
    }

    public async Task<Result<TriggerWorkflowResult>> Handle(
        TriggerWorkflowCommand request,
        CancellationToken cancellationToken)
    {
        var ctx = request.Context;

        // 1. Resolve the right template for this entity + context
        var template = _templateProvider.Resolve(ctx);
        if (template is null)
        {
            _logger.LogDebug(
                "No workflow template for {EntityType}/{EntityId} — skipping workflow",
                ctx.EntityType, ctx.EntityId);

            return Result.Ok(new TriggerWorkflowResult(
                WorkflowInstanceId: string.Empty,
                TemplateName: string.Empty,
                EntityType: ctx.EntityType,
                EntityId: ctx.EntityId,
                WasTriggered: false));
        }

        // 2. Ensure template is registered in Go engine (idempotent)
        var registerResult = await _mediator.Send(
            new RegisterWorkflowTemplateCommand(template), cancellationToken);

        if (registerResult.IsFailure)
        {
            _logger.LogError(
                "Failed to register template '{Name}': {Error}",
                template.Name, registerResult.Error?.Message);
            return Result.Fail<TriggerWorkflowResult>(
                "TEMPLATE_REGISTRATION_FAILED",
                $"Could not register workflow template '{template.Name}'");
        }

        var definitionId = registerResult.Value!.DefinitionId;

        // 3. Build context struct for Go engine
        var contextStruct = BuildContextStruct(ctx, template.Name);

        // 4. Start the workflow instance
        var instanceId = await _gateway.StartWorkflowAsync(
            definitionId,
            ctx.EntityType,
            ctx.EntityId,
            contextStruct,
            cancellationToken);

        if (instanceId is null)
        {
            return Result.Fail<TriggerWorkflowResult>(
                "WORKFLOW_START_FAILED",
                $"Failed to start workflow '{template.Name}' for {ctx.EntityType}/{ctx.EntityId}");
        }

        _logger.LogInformation(
            "Workflow started: instance={InstanceId} template='{Template}' entity={EntityType}/{EntityId} amount={AmountPaisa}p",
            instanceId, template.Name, ctx.EntityType, ctx.EntityId, ctx.AmountPaisa);

        return Result.Ok(new TriggerWorkflowResult(
            instanceId,
            template.Name,
            ctx.EntityType,
            ctx.EntityId,
            WasTriggered: true));
    }

    private static Struct BuildContextStruct(WorkflowTriggerContext ctx, string templateName)
    {
        var s = new Struct();
        s.Fields["entity_type"]    = Value.ForString(ctx.EntityType);
        s.Fields["entity_id"]      = Value.ForString(ctx.EntityId);
        s.Fields["initiated_by"]   = Value.ForString(ctx.InitiatedBy);
        s.Fields["template_name"]  = Value.ForString(templateName);
        s.Fields["portal"]         = Value.ForString(ctx.Portal);
        s.Fields["amount_paisa"]   = Value.ForNumber(ctx.AmountPaisa);

        if (ctx.SubType is not null)
            s.Fields["sub_type"] = Value.ForString(ctx.SubType);

        foreach (var kv in ctx.Metadata)
            s.Fields[kv.Key] = Value.ForString(kv.Value);

        return s;
    }
}
