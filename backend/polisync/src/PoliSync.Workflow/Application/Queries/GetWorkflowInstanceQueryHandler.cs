using MediatR;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Workflow.Infrastructure;
using Insuretech.Workflow.Entity.V1;

namespace PoliSync.Workflow.Application.Queries;

public sealed class GetWorkflowInstanceQueryHandler
    : IRequestHandler<GetWorkflowInstanceQuery, Result<GetWorkflowInstanceResult>>,
      IRequestHandler<GetWorkflowHistoryQuery, Result<IReadOnlyList<WorkflowInstance>>>,
      IRequestHandler<GetMyWorkflowTasksQuery, Result<GetMyTasksResult>>
{
    private readonly IWorkflowDataGateway _gateway;

    public GetWorkflowInstanceQueryHandler(IWorkflowDataGateway gateway)
        => _gateway = gateway;

    public async Task<Result<GetWorkflowInstanceResult>> Handle(
        GetWorkflowInstanceQuery request,
        CancellationToken cancellationToken)
    {
        var result = await _gateway.GetWorkflowInstanceAsync(
            request.WorkflowInstanceId, cancellationToken);

        if (result is null)
            return Result.Fail<GetWorkflowInstanceResult>(
                "NOT_FOUND", $"Workflow instance {request.WorkflowInstanceId} not found");

        return Result.Ok(result);
    }

    public async Task<Result<IReadOnlyList<WorkflowInstance>>> Handle(
        GetWorkflowHistoryQuery request,
        CancellationToken cancellationToken)
    {
        var instances = await _gateway.GetWorkflowHistoryAsync(
            request.EntityType, request.EntityId, cancellationToken);

        return Result.Ok(instances);
    }

    public async Task<Result<GetMyTasksResult>> Handle(
        GetMyWorkflowTasksQuery request,
        CancellationToken cancellationToken)
    {
        var result = await _gateway.GetMyTasksAsync(
            request.UserId, request.Status, request.Page, request.PageSize, cancellationToken);

        return Result.Ok(result);
    }
}
