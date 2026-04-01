using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Workflow.Entity.V1;
using Insuretech.Workflow.Services.V1;
using Microsoft.Extensions.Logging;
using PoliSync.RulesEngine.Services;

namespace PoliSync.RulesEngine.GrpcServices;

public class BusinessWorkflowGrpcService : Insuretech.Workflow.Services.V1.BusinessWorkflowService.BusinessWorkflowServiceBase
{
    private readonly IBusinessWorkflowService _workflowService;
    private readonly ILogger<BusinessWorkflowGrpcService> _logger;

    public BusinessWorkflowGrpcService(
        IBusinessWorkflowService workflowService,
        ILogger<BusinessWorkflowGrpcService> logger)
    {
        _workflowService = workflowService;
        _logger = logger;
    }

    public override async Task<EvaluateBusinessWorkflowResponse> EvaluateBusinessWorkflow(
        EvaluateBusinessWorkflowRequest request,
        ServerCallContext context)
    {
        try
        {
            // Convert Struct to Dictionary
            var inputs = request.Inputs.Fields.ToDictionary(
                f => f.Key,
                f => ConvertValue(f.Value));

            var execution = await _workflowService.EvaluateWorkflowAsync(
                request.WorkflowName,
                inputs,
                request.EntityType,
                request.EntityId,
                request.ExecutedBy,
                context.CancellationToken);

            return new EvaluateBusinessWorkflowResponse
            {
                ExecutionId = execution.ExecutionId,
                WorkflowName = request.WorkflowName,
                IsSuccess = execution.IsSuccess,
                Results = { execution.Results },
                Outputs = string.IsNullOrEmpty(execution.OutputsJson) 
                    ? new Struct() 
                    : JsonParser.Default.Parse<Struct>(execution.OutputsJson),
                ExecutionTimeMs = execution.ExecutionTimeMs,
                ExecutedAt = execution.ExecutedAt
            };
        }
        catch (InvalidOperationException ex)
        {
            _logger.LogWarning(ex, "Workflow not found: {WorkflowName}", request.WorkflowName);
            return new EvaluateBusinessWorkflowResponse
            {
                Error = new Error
                {
                    Code = "WORKFLOW_NOT_FOUND",
                    Message = ex.Message,
                    HttpStatusCode = 404
                }
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error evaluating workflow: {WorkflowName}", request.WorkflowName);
            return new EvaluateBusinessWorkflowResponse
            {
                Error = new Error
                {
                    Code = "EVALUATION_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<CreateBusinessWorkflowResponse> CreateBusinessWorkflow(
        CreateBusinessWorkflowRequest request,
        ServerCallContext context)
    {
        try
        {
            var workflow = new BusinessWorkflowDefinition
            {
                WorkflowName = request.WorkflowName,
                Description = request.Description,
                WorkflowType = request.WorkflowType,
                Status = BusinessWorkflowStatus.Active,
                RulesConfig = request.Rules.ToString(),
                InputSchema = request.InputSchema?.ToString() ?? string.Empty,
                OutputSchema = request.OutputSchema?.ToString() ?? string.Empty,
                Metadata = { ConvertStructToDictionary(request.Metadata) },
                CreatedBy = request.CreatedBy
            };

            var created = await _workflowService.CreateWorkflowAsync(workflow, context.CancellationToken);

            return new CreateBusinessWorkflowResponse
            {
                Workflow = created
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating workflow: {WorkflowName}", request.WorkflowName);
            return new CreateBusinessWorkflowResponse
            {
                Error = new Error
                {
                    Code = "CREATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<GetBusinessWorkflowResponse> GetBusinessWorkflow(
        GetBusinessWorkflowRequest request,
        ServerCallContext context)
    {
        var workflow = await _workflowService.GetWorkflowAsync(
            request.BusinessWorkflowId, 
            context.CancellationToken);

        if (workflow == null)
        {
            return new GetBusinessWorkflowResponse
            {
                Error = new Error
                {
                    Code = "WORKFLOW_NOT_FOUND",
                    Message = $"Workflow {request.BusinessWorkflowId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetBusinessWorkflowResponse
        {
            Workflow = workflow
        };
    }

    public override async Task<GetBusinessWorkflowResponse> GetBusinessWorkflowByName(
        GetBusinessWorkflowByNameRequest request,
        ServerCallContext context)
    {
        var workflow = await _workflowService.GetWorkflowByNameAsync(
            request.WorkflowName,
            context.CancellationToken);

        if (workflow == null)
        {
            return new GetBusinessWorkflowResponse
            {
                Error = new Error
                {
                    Code = "WORKFLOW_NOT_FOUND",
                    Message = $"Workflow '{request.WorkflowName}' not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetBusinessWorkflowResponse
        {
            Workflow = workflow
        };
    }

    public override async Task<ListBusinessWorkflowsResponse> ListBusinessWorkflows(
        ListBusinessWorkflowsRequest request,
        ServerCallContext context)
    {
        var workflows = await _workflowService.GetWorkflowsAsync(context.CancellationToken);

        // Apply filters
        var filtered = workflows.AsEnumerable();
        if (request.WorkflowType != BusinessWorkflowType.Unspecified)
        {
            filtered = filtered.Where(w => w.WorkflowType == request.WorkflowType);
        }
        if (request.Status != BusinessWorkflowStatus.Unspecified)
        {
            filtered = filtered.Where(w => w.Status == request.Status);
        }

        var workflowList = filtered.ToList();

        return new ListBusinessWorkflowsResponse
        {
            Workflows = { workflowList },
            TotalCount = workflowList.Count
        };
    }

    public override async Task<UpdateBusinessWorkflowResponse> UpdateBusinessWorkflow(
        UpdateBusinessWorkflowRequest request,
        ServerCallContext context)
    {
        try
        {
            var workflow = new BusinessWorkflowDefinition
            {
                BusinessWorkflowId = request.BusinessWorkflowId,
                WorkflowName = request.WorkflowName,
                Description = request.Description,
                WorkflowType = request.WorkflowType,
                Status = request.Status,
                RulesConfig = request.Rules.ToString(),
                InputSchema = request.InputSchema?.ToString() ?? string.Empty,
                OutputSchema = request.OutputSchema?.ToString() ?? string.Empty,
                Metadata = { ConvertStructToDictionary(request.Metadata) }
            };

            var updated = await _workflowService.UpdateWorkflowAsync(
                request.BusinessWorkflowId,
                workflow,
                context.CancellationToken);

            if (updated == null)
            {
                return new UpdateBusinessWorkflowResponse
                {
                    Error = new Error
                    {
                        Code = "WORKFLOW_NOT_FOUND",
                        Message = $"Workflow {request.BusinessWorkflowId} not found",
                        HttpStatusCode = 404
                    }
                };
            }

            return new UpdateBusinessWorkflowResponse
            {
                Workflow = updated
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating workflow: {WorkflowId}", request.BusinessWorkflowId);
            return new UpdateBusinessWorkflowResponse
            {
                Error = new Error
                {
                    Code = "UPDATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<DeleteBusinessWorkflowResponse> DeleteBusinessWorkflow(
        DeleteBusinessWorkflowRequest request,
        ServerCallContext context)
    {
        var result = await _workflowService.DeleteWorkflowAsync(
            request.BusinessWorkflowId,
            request.Permanent,
            context.CancellationToken);

        return new DeleteBusinessWorkflowResponse
        {
            Success = result
        };
    }

    public override async Task<GetExecutionHistoryResponse> GetExecutionHistory(
        GetExecutionHistoryRequest request,
        ServerCallContext context)
    {
        var executions = await _workflowService.GetExecutionHistoryAsync(
            request.BusinessWorkflowId,
            context.CancellationToken);

        var executionList = executions.ToList();

        return new GetExecutionHistoryResponse
        {
            Executions = { executionList },
            TotalCount = executionList.Count
        };
    }

    public override Task<ValidateRuleExpressionResponse> ValidateRuleExpression(
        ValidateRuleExpressionRequest request,
        ServerCallContext context)
    {
        // This would need to be implemented with the evaluation service
        // For now, return success
        return Task.FromResult(new ValidateRuleExpressionResponse
        {
            IsValid = true,
            TestPassed = true
        });
    }

    private static object ConvertValue(Value value)
    {
        return value.KindCase switch
        {
            Value.KindOneofCase.StringValue => value.StringValue,
            Value.KindOneofCase.NumberValue => value.NumberValue,
            Value.KindOneofCase.BoolValue => value.BoolValue,
            Value.KindOneofCase.StructValue => value.StructValue,
            Value.KindOneofCase.ListValue => value.ListValue,
            _ => value.ToString()
        };
    }

    private static Dictionary<string, string> ConvertStructToDictionary(Google.Protobuf.WellKnownTypes.Struct? structValue)
    {
        if (structValue == null)
            return new Dictionary<string, string>();
        
        return structValue.Fields.ToDictionary(
            kvp => kvp.Key, 
            kvp => kvp.Value.ToString() ?? string.Empty);
    }
}
