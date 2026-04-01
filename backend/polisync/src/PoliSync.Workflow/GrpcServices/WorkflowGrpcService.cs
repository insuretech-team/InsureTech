using Grpc.Core;
using Insuretech.Workflow.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Application.Queries;
using Google.Protobuf.WellKnownTypes;

namespace PoliSync.Workflow.GrpcServices;

/// <summary>
/// PoliSync-side gRPC facade for the workflow engine.
/// This service acts as a smart proxy: it forwards requests to the Go
/// workflow-engine via MediatR commands/queries, applying C# business rules,
/// authorization, and domain validation before/after Go calls.
/// </summary>
public sealed class WorkflowGrpcService : Insuretech.Workflow.Services.V1.WorkflowService.WorkflowServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<WorkflowGrpcService> _logger;

    public WorkflowGrpcService(IMediator mediator, ILogger<WorkflowGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<StartWorkflowResponse> StartWorkflow(
        StartWorkflowRequest request,
        ServerCallContext context)
    {
        var result = await _mediator.Send(new StartWorkflowCommand(
            TemplateName: request.WorkflowDefinitionId, // treat as name for C# layer
            EntityType: request.EntityType,
            EntityId: request.EntityId,
            InitiatedBy: context.RequestHeaders.GetValue("x-user-id") ?? "system",
            Context: null));

        if (result.IsFailure)
        {
            _logger.LogWarning("StartWorkflow failed: {Error}", result.Error);
            return new StartWorkflowResponse
            {
                Error = new Insuretech.Common.V1.Error
                {
                    Code = "400",
                    Message = result.Error?.Message ?? "Failed to start workflow"
                }
            };
        }

        return new StartWorkflowResponse
        {
            WorkflowInstanceId = result.Value!.WorkflowInstanceId,
            Message = "Workflow started successfully"
        };
    }

    public override async Task<GetWorkflowInstanceResponse> GetWorkflowInstance(
        GetWorkflowInstanceRequest request,
        ServerCallContext context)
    {
        var result = await _mediator.Send(new GetWorkflowInstanceQuery(request.WorkflowInstanceId));

        if (result.IsFailure)
            throw new RpcException(new Status(StatusCode.NotFound, result.Error?.Message ?? "Not found"));

        return new GetWorkflowInstanceResponse
        {
            WorkflowInstance = result.Value!.Instance,
            Tasks = { result.Value.Tasks }
        };
    }

    public override async Task<GetMyTasksResponse> GetMyTasks(
        GetMyTasksRequest request,
        ServerCallContext context)
    {
        var userId = context.RequestHeaders.GetValue("x-user-id") ?? string.Empty;
        var result = await _mediator.Send(new GetMyWorkflowTasksQuery(
            userId, request.Status, request.Page, request.PageSize));

        if (result.IsFailure)
            return new GetMyTasksResponse();

        return new GetMyTasksResponse
        {
            Tasks = { result.Value!.Tasks },
            TotalCount = result.Value.TotalCount
        };
    }

    public override async Task<CompleteWorkflowTaskResponse> CompleteTask(
        CompleteWorkflowTaskRequest request,
        ServerCallContext context)
    {
        var userId = context.RequestHeaders.GetValue("x-user-id") ?? string.Empty;
        var result = await _mediator.Send(new CompleteWorkflowTaskCommand(
            request.TaskId,
            request.Decision,
            request.Comments,
            userId));

        if (result.IsFailure)
        {
            return new CompleteWorkflowTaskResponse
            {
                Error = new Insuretech.Common.V1.Error
                {
                    Code = "400",
                    Message = result.Error?.Message ?? "Failed to complete task"
                }
            };
        }

        return new CompleteWorkflowTaskResponse
        {
            Message = $"Task completed with decision: {request.Decision}"
        };
    }

    public override async Task<GetWorkflowHistoryResponse> GetWorkflowHistory(
        GetWorkflowHistoryRequest request,
        ServerCallContext context)
    {
        var result = await _mediator.Send(new GetWorkflowHistoryQuery(
            request.EntityType, request.EntityId));

        if (result.IsFailure)
            return new GetWorkflowHistoryResponse();

        return new GetWorkflowHistoryResponse
        {
            WorkflowInstances = { result.Value! }
        };
    }
}
