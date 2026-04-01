using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Workflow.Infrastructure;

namespace PoliSync.Workflow.Application.Commands;

public sealed class CompleteWorkflowTaskCommandHandler
    : IRequestHandler<CompleteWorkflowTaskCommand, Result<CompleteWorkflowTaskResult>>
{
    private readonly IWorkflowDataGateway _gateway;
    private readonly ILogger<CompleteWorkflowTaskCommandHandler> _logger;

    public CompleteWorkflowTaskCommandHandler(
        IWorkflowDataGateway gateway,
        ILogger<CompleteWorkflowTaskCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<CompleteWorkflowTaskResult>> Handle(
        CompleteWorkflowTaskCommand request,
        CancellationToken cancellationToken)
    {
        if (request.Decision is not ("APPROVED" or "REJECTED" or "RETURNED"))
            return Result.Fail<CompleteWorkflowTaskResult>(
                "INVALID_DECISION",
                "Decision must be APPROVED, REJECTED, or RETURNED");

        var result = await _gateway.CompleteTaskAsync(
            request.TaskId,
            request.Decision,
            request.Comments,
            request.CompletedBy,
            cancellationToken);

        if (!result.Success)
            return Result.Fail<CompleteWorkflowTaskResult>(
                result.ErrorCode ?? "COMPLETE_TASK_FAILED",
                result.ErrorMessage ?? "Failed to complete workflow task");

        _logger.LogInformation(
            "Task {TaskId} completed by {UserId} with decision {Decision}",
            request.TaskId, request.CompletedBy, request.Decision);

        return Result.Ok(new CompleteWorkflowTaskResult(
            request.TaskId,
            request.Decision,
            result.WorkflowInstanceId ?? string.Empty,
            result.WorkflowCompleted));
    }
}
