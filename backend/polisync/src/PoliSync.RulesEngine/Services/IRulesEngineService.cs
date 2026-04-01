namespace PoliSync.RulesEngine.Services;

public interface IRulesEngineService
{
    Task<Models.WorkflowEvaluationResult> EvaluateWorkflowAsync(string workflowName, Dictionary<string, object> inputs, CancellationToken cancellationToken = default);
    Task<IEnumerable<Models.Workflow>> GetWorkflowsAsync(CancellationToken cancellationToken = default);
    Task<Models.Workflow?> GetWorkflowAsync(string workflowId, CancellationToken cancellationToken = default);
    Task<Models.Workflow> CreateWorkflowAsync(Models.Workflow workflow, CancellationToken cancellationToken = default);
    Task<Models.Workflow?> UpdateWorkflowAsync(string workflowId, Models.Workflow workflow, CancellationToken cancellationToken = default);
    Task<bool> DeleteWorkflowAsync(string workflowId, CancellationToken cancellationToken = default);
}

public interface IWorkflowRepository
{
    Task<Models.Workflow?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<Models.Workflow?> GetByNameAsync(string name, CancellationToken cancellationToken = default);
    Task<IEnumerable<Models.Workflow>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<Models.Workflow> CreateAsync(Models.Workflow workflow, CancellationToken cancellationToken = default);
    Task<Models.Workflow?> UpdateAsync(Models.Workflow workflow, CancellationToken cancellationToken = default);
    Task<bool> DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public interface IRuleEvaluationService
{
    Task<Models.WorkflowEvaluationResult> EvaluateAsync(string workflowName, Dictionary<string, object> inputs, CancellationToken cancellationToken = default);
}
