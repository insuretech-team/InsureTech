using Microsoft.Extensions.Logging;
using RulesEngine.Models;
using RulesEngineAlias = RulesEngine.RulesEngine;

namespace PoliSync.RulesEngine.Services;

public class RuleEvaluationService : IRuleEvaluationService
{
    private readonly RulesEngineAlias _rulesEngine;
    private readonly IWorkflowRepository _workflowRepository;
    private readonly ILogger<RuleEvaluationService> _logger;

    public RuleEvaluationService(
        RulesEngineAlias rulesEngine,
        IWorkflowRepository workflowRepository,
        ILogger<RuleEvaluationService> logger)
    {
        _rulesEngine = rulesEngine;
        _workflowRepository = workflowRepository;
        _logger = logger;
    }

    public async Task<Models.WorkflowEvaluationResult> EvaluateAsync(
        string workflowName, 
        Dictionary<string, object> inputs, 
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Evaluating rules for workflow: {WorkflowName}", workflowName);

        // Ensure workflow is loaded
        var workflow = await _workflowRepository.GetByNameAsync(workflowName, cancellationToken);
        if (workflow == null)
        {
            throw new InvalidOperationException($"Workflow '{workflowName}' not found");
        }

        // Convert inputs to named parameters
        var inputsList = new List<RuleParameter>();
        foreach (var input in inputs)
        {
            inputsList.Add(new RuleParameter(input.Key, input.Value));
        }

        // Execute rules
        var results = await _rulesEngine.ExecuteAllRulesAsync(workflowName, inputsList.ToArray());

        // Convert results
        var evaluationResults = results.Select(r => new Models.RuleEvaluationResult
        {
            RuleName = r.Rule.RuleName,
            IsSuccess = r.IsSuccess,
            SuccessEvent = r.Rule.SuccessEvent,
            ErrorMessage = r.ExceptionMessage ?? r.Rule.ErrorMessage,
            ChildResults = r.ChildResults?.Select(c => c.Rule.RuleName).ToList()
        }).ToList();

        var isSuccess = results.All(r => r.IsSuccess);

        _logger.LogInformation(
            "Workflow {WorkflowName} evaluation completed. Success: {IsSuccess}, Rules evaluated: {RuleCount}",
            workflowName, isSuccess, evaluationResults.Count);

        return new Models.WorkflowEvaluationResult
        {
            WorkflowName = workflowName,
            IsSuccess = isSuccess,
            Results = evaluationResults,
            Outputs = isSuccess ? inputs : null
        };
    }
}
