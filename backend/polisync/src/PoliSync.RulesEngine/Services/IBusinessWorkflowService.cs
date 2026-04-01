using Insuretech.Workflow.Entity.V1;
using RuleExpressionType = Insuretech.Workflow.Entity.V1.RuleExpressionType;

namespace PoliSync.RulesEngine.Services;

public interface IBusinessWorkflowService
{
    Task<BusinessWorkflowExecution> EvaluateWorkflowAsync(
        string workflowName, 
        Dictionary<string, object> inputs,
        string entityType,
        string entityId,
        string? executedBy = null,
        CancellationToken cancellationToken = default);
    
    Task<IEnumerable<BusinessWorkflowDefinition>> GetWorkflowsAsync(CancellationToken cancellationToken = default);
    Task<BusinessWorkflowDefinition?> GetWorkflowAsync(string workflowId, CancellationToken cancellationToken = default);
    Task<BusinessWorkflowDefinition?> GetWorkflowByNameAsync(string workflowName, CancellationToken cancellationToken = default);
    Task<BusinessWorkflowDefinition> CreateWorkflowAsync(BusinessWorkflowDefinition workflow, CancellationToken cancellationToken = default);
    Task<BusinessWorkflowDefinition?> UpdateWorkflowAsync(string workflowId, BusinessWorkflowDefinition workflow, CancellationToken cancellationToken = default);
    Task<bool> DeleteWorkflowAsync(string workflowId, bool permanent = false, CancellationToken cancellationToken = default);
    Task<IEnumerable<BusinessWorkflowExecution>> GetExecutionHistoryAsync(string workflowId, CancellationToken cancellationToken = default);
}

public interface IBusinessWorkflowRepository
{
    Task<BusinessWorkflowDefinition?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<BusinessWorkflowDefinition?> GetByNameAsync(string name, CancellationToken cancellationToken = default);
    Task<IEnumerable<BusinessWorkflowDefinition>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<BusinessWorkflowDefinition> CreateAsync(BusinessWorkflowDefinition workflow, CancellationToken cancellationToken = default);
    Task<BusinessWorkflowDefinition?> UpdateAsync(BusinessWorkflowDefinition workflow, CancellationToken cancellationToken = default);
    Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default);
    Task<BusinessWorkflowExecution> LogExecutionAsync(BusinessWorkflowExecution execution, CancellationToken cancellationToken = default);
    Task<IEnumerable<BusinessWorkflowExecution>> GetExecutionsByWorkflowIdAsync(string workflowId, CancellationToken cancellationToken = default);
}

public interface IBusinessRuleEvaluationService
{
    Task<BusinessWorkflowExecution> EvaluateAsync(
        BusinessWorkflowDefinition workflow,
        Dictionary<string, object> inputs,
        string entityType,
        string entityId,
        string? executedBy = null,
        CancellationToken cancellationToken = default);
    
    Task<(bool IsValid, string? ErrorMessage)> ValidateRuleAsync(string expression, RuleExpressionType expressionType);
}
