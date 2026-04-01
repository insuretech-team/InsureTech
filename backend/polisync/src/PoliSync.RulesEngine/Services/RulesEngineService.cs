using Microsoft.Extensions.Logging;
using RulesEngineAlias = RulesEngine.RulesEngine;

namespace PoliSync.RulesEngine.Services;

public class RulesEngineService : IRulesEngineService
{
    private readonly RulesEngineAlias _rulesEngine;
    private readonly IWorkflowRepository _workflowRepository;
    private readonly IRuleEvaluationService _evaluationService;
    private readonly ILogger<RulesEngineService> _logger;

    public RulesEngineService(
        RulesEngineAlias rulesEngine,
        IWorkflowRepository workflowRepository,
        IRuleEvaluationService evaluationService,
        ILogger<RulesEngineService> logger)
    {
        _rulesEngine = rulesEngine;
        _workflowRepository = workflowRepository;
        _evaluationService = evaluationService;
        _logger = logger;
    }

    public async Task<Models.WorkflowEvaluationResult> EvaluateWorkflowAsync(
        string workflowName, 
        Dictionary<string, object> inputs, 
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Evaluating workflow: {WorkflowName}", workflowName);
        return await _evaluationService.EvaluateAsync(workflowName, inputs, cancellationToken);
    }

    public async Task<IEnumerable<Models.Workflow>> GetWorkflowsAsync(CancellationToken cancellationToken = default)
    {
        return await _workflowRepository.GetAllAsync(cancellationToken);
    }

    public async Task<Models.Workflow?> GetWorkflowAsync(string workflowId, CancellationToken cancellationToken = default)
    {
        return await _workflowRepository.GetByIdAsync(workflowId, cancellationToken);
    }

    public async Task<Models.Workflow> CreateWorkflowAsync(Models.Workflow workflow, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Creating workflow: {WorkflowName}", workflow.WorkflowName);
        
        // Add to RulesEngine
        var reWorkflow = ConvertToReWorkflow(workflow);
        _rulesEngine.AddOrUpdateWorkflow(reWorkflow);
        
        return await _workflowRepository.CreateAsync(workflow, cancellationToken);
    }

    public async Task<Models.Workflow?> UpdateWorkflowAsync(string workflowId, Models.Workflow workflow, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Updating workflow: {WorkflowId}", workflowId);
        
        workflow.Id = workflowId;
        workflow.UpdatedAt = DateTime.UtcNow;
        
        // Update in RulesEngine
        var reWorkflow = ConvertToReWorkflow(workflow);
        _rulesEngine.AddOrUpdateWorkflow(reWorkflow);
        
        return await _workflowRepository.UpdateAsync(workflow, cancellationToken);
    }

    public async Task<bool> DeleteWorkflowAsync(string workflowId, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Deleting workflow: {WorkflowId}", workflowId);
        return await _workflowRepository.DeleteAsync(workflowId, cancellationToken);
    }

    private static global::RulesEngine.Models.Workflow ConvertToReWorkflow(Models.Workflow workflow)
    {
        return new global::RulesEngine.Models.Workflow
        {
            WorkflowName = workflow.WorkflowName,
            Rules = workflow.Rules.Select(r => new global::RulesEngine.Models.Rule
            {
                RuleName = r.RuleName,
                Expression = r.Expression,
                RuleExpressionType = (global::RulesEngine.Models.RuleExpressionType)(int)r.ExpressionType,
                SuccessEvent = r.SuccessEvent,
                ErrorMessage = r.ErrorMessage,
                Rules = r.Rules?.Select(ConvertToReRule).ToList()
            }).ToList()
        };
    }

    private static global::RulesEngine.Models.Rule ConvertToReRule(Models.Rule rule)
    {
        return new global::RulesEngine.Models.Rule
        {
            RuleName = rule.RuleName,
            Expression = rule.Expression,
                RuleExpressionType = (global::RulesEngine.Models.RuleExpressionType)(int)rule.ExpressionType,
                SuccessEvent = rule.SuccessEvent,
                ErrorMessage = rule.ErrorMessage
        };
    }
}
