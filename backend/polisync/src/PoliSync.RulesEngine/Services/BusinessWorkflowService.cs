using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Workflow.Entity.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.RulesEngine.Services;

public class BusinessWorkflowService : IBusinessWorkflowService
{
    private readonly IBusinessWorkflowRepository _repository;
    private readonly IBusinessRuleEvaluationService _evaluationService;
    private readonly ILogger<BusinessWorkflowService> _logger;

    public BusinessWorkflowService(
        IBusinessWorkflowRepository repository,
        IBusinessRuleEvaluationService evaluationService,
        ILogger<BusinessWorkflowService> logger)
    {
        _repository = repository;
        _evaluationService = evaluationService;
        _logger = logger;
    }

    public async Task<BusinessWorkflowExecution> EvaluateWorkflowAsync(
        string workflowName,
        Dictionary<string, object> inputs,
        string entityType,
        string entityId,
        string? executedBy = null,
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation(
            "Evaluating business workflow: {WorkflowName} for {EntityType}:{EntityId}",
            workflowName, entityType, entityId);

        // Get workflow by name
        var workflow = await _repository.GetByNameAsync(workflowName, cancellationToken);
        if (workflow == null)
        {
            throw new InvalidOperationException($"Business workflow '{workflowName}' not found");
        }

        if (workflow.Status != BusinessWorkflowStatus.Active)
        {
            throw new InvalidOperationException($"Business workflow '{workflowName}' is not active");
        }

        // Evaluate rules
        var execution = await _evaluationService.EvaluateAsync(
            workflow, inputs, entityType, entityId, executedBy, cancellationToken);

        // Log execution
        await _repository.LogExecutionAsync(execution, cancellationToken);

        return execution;
    }

    public Task<IEnumerable<BusinessWorkflowDefinition>> GetWorkflowsAsync(CancellationToken cancellationToken = default)
    {
        return _repository.GetAllAsync(cancellationToken);
    }

    public Task<BusinessWorkflowDefinition?> GetWorkflowAsync(string workflowId, CancellationToken cancellationToken = default)
    {
        return _repository.GetByIdAsync(workflowId, cancellationToken);
    }

    public Task<BusinessWorkflowDefinition?> GetWorkflowByNameAsync(string workflowName, CancellationToken cancellationToken = default)
    {
        return _repository.GetByNameAsync(workflowName, cancellationToken);
    }

    public Task<BusinessWorkflowDefinition> CreateWorkflowAsync(BusinessWorkflowDefinition workflow, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Creating business workflow: {WorkflowName}", workflow.WorkflowName);
        return _repository.CreateAsync(workflow, cancellationToken);
    }

    public Task<BusinessWorkflowDefinition?> UpdateWorkflowAsync(
        string workflowId, 
        BusinessWorkflowDefinition workflow, 
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Updating business workflow: {WorkflowId}", workflowId);
        workflow.BusinessWorkflowId = workflowId;
        return _repository.UpdateAsync(workflow, cancellationToken);
    }

    public Task<bool> DeleteWorkflowAsync(
        string workflowId, 
        bool permanent = false, 
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Deleting business workflow: {WorkflowId} (permanent: {Permanent})", 
            workflowId, permanent);
        return _repository.DeleteAsync(workflowId, permanent, cancellationToken);
    }

    public async Task<IEnumerable<BusinessWorkflowExecution>> GetExecutionHistoryAsync(
        string workflowId, 
        CancellationToken cancellationToken = default)
    {
        return await _repository.GetExecutionsByWorkflowIdAsync(workflowId, cancellationToken);
    }
}
