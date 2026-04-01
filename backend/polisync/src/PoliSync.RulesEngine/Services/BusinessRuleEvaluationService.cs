using System.Diagnostics;
using System.Text.Json;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Workflow.Entity.V1;
using Microsoft.Extensions.Logging;
using RulesEngineAlias = RulesEngine.RulesEngine;

namespace PoliSync.RulesEngine.Services;

public class BusinessRuleEvaluationService : IBusinessRuleEvaluationService
{
    private readonly RulesEngineAlias _rulesEngine;
    private readonly ILogger<BusinessRuleEvaluationService> _logger;

    public BusinessRuleEvaluationService(
        RulesEngineAlias rulesEngine,
        ILogger<BusinessRuleEvaluationService> logger)
    {
        _rulesEngine = rulesEngine;
        _logger = logger;
    }

    public async Task<BusinessWorkflowExecution> EvaluateAsync(
        BusinessWorkflowDefinition workflow,
        Dictionary<string, object> inputs,
        string entityType,
        string entityId,
        string? executedBy = null,
        CancellationToken cancellationToken = default)
    {
        var stopwatch = Stopwatch.StartNew();
        var executionId = Guid.NewGuid().ToString();

        _logger.LogInformation(
            "Evaluating business workflow: {WorkflowName} for {EntityType}:{EntityId}",
            workflow.WorkflowName, entityType, entityId);

        try
        {
            // Parse rules from JSON config
            var rules = ParseRulesFromConfig(workflow.RulesConfig);

            // Create workflow for RulesEngine
            var reWorkflow = new global::RulesEngine.Models.Workflow
            {
                WorkflowName = workflow.WorkflowName,
                Rules = rules
            };

            // Add or update workflow in RulesEngine
            _rulesEngine.AddOrUpdateWorkflow(reWorkflow);

            // Prepare inputs for RulesEngine
            var ruleParams = new List<global::RulesEngine.Models.RuleParameter>();
            foreach (var input in inputs)
            {
                ruleParams.Add(new global::RulesEngine.Models.RuleParameter(input.Key, input.Value));
            }

            // Execute all rules
            var results = await _rulesEngine.ExecuteAllRulesAsync(
                workflow.WorkflowName, 
                ruleParams.ToArray());

            stopwatch.Stop();

            // Convert results to proto format
            var ruleResults = results.Select(r => new BusinessRuleResult
            {
                RuleName = r.Rule.RuleName,
                IsSuccess = r.IsSuccess,
                SuccessEvent = r.Rule.SuccessEvent ?? string.Empty,
                ErrorMessage = r.ExceptionMessage ?? r.Rule.ErrorMessage ?? string.Empty,
                ExecutionTimeMs = 0 // Individual timing not available from RulesEngine
            }).ToList();

            var isSuccess = results.All(r => r.IsSuccess);

            _logger.LogInformation(
                "Workflow {WorkflowName} evaluation completed. Success: {IsSuccess}, " +
                "Rules passed: {Passed}/{Total}, Time: {ElapsedMs}ms",
                workflow.WorkflowName, isSuccess,
                results.Count(r => r.IsSuccess), results.Count,
                stopwatch.ElapsedMilliseconds);

            // Build outputs from successful rules
            var outputs = new Dictionary<string, object>();
            foreach (var result in results.Where(r => r.IsSuccess && !string.IsNullOrEmpty(r.Rule.SuccessEvent)))
            {
                outputs[result.Rule.RuleName] = result.Rule.SuccessEvent;
            }

            return new BusinessWorkflowExecution
            {
                ExecutionId = executionId,
                BusinessWorkflowId = workflow.BusinessWorkflowId,
                EntityType = entityType,
                EntityId = entityId,
                IsSuccess = isSuccess,
                Results = { ruleResults },
                InputsJson = JsonSerializer.Serialize(inputs),
                OutputsJson = JsonSerializer.Serialize(outputs),
                ExecutedBy = executedBy ?? "system",
                ExecutedAt = Timestamp.FromDateTime(DateTime.UtcNow),
                ExecutionTimeMs = (int)stopwatch.ElapsedMilliseconds
            };
        }
        catch (Exception ex)
        {
            stopwatch.Stop();
            _logger.LogError(ex, 
                "Error evaluating workflow {WorkflowName}: {ErrorMessage}",
                workflow.WorkflowName, ex.Message);

            return new BusinessWorkflowExecution
            {
                ExecutionId = executionId,
                BusinessWorkflowId = workflow.BusinessWorkflowId,
                EntityType = entityType,
                EntityId = entityId,
                IsSuccess = false,
                Results = 
                {
                    new BusinessRuleResult
                    {
                        RuleName = "System",
                        IsSuccess = false,
                        ErrorMessage = $"Evaluation error: {ex.Message}"
                    }
                },
                InputsJson = JsonSerializer.Serialize(inputs),
                OutputsJson = JsonSerializer.Serialize(new Dictionary<string, object>()),
                ExecutedBy = executedBy ?? "system",
                ExecutedAt = Timestamp.FromDateTime(DateTime.UtcNow),
                ExecutionTimeMs = (int)stopwatch.ElapsedMilliseconds
            };
        }
    }

    public Task<(bool IsValid, string? ErrorMessage)> ValidateRuleAsync(
        string expression, 
        Insuretech.Workflow.Entity.V1.RuleExpressionType expressionType)
    {
        try
        {
            // Create a test rule
            var testRule = new global::RulesEngine.Models.Rule
            {
                RuleName = "ValidationTest",
                Expression = expression,
                RuleExpressionType = (global::RulesEngine.Models.RuleExpressionType)(int)expressionType,
                SuccessEvent = "VALID"
            };

            var testWorkflow = new global::RulesEngine.Models.Workflow
            {
                WorkflowName = "ValidationWorkflow",
                Rules = new List<global::RulesEngine.Models.Rule> { testRule }
            };

            // Try to add to RulesEngine - this will validate syntax
            _rulesEngine.AddOrUpdateWorkflow(testWorkflow);

            // Try to execute with test data
            var testInput = new global::RulesEngine.Models.RuleParameter("input1", new Dictionary<string, object>());
            var result = _rulesEngine.ExecuteAllRulesAsync("ValidationWorkflow", testInput).Result;

            return Task.FromResult<(bool, string?)>((true, null));
        }
        catch (Exception ex)
        {
            return Task.FromResult<(bool, string?)>((false, ex.Message));
        }
    }

    private static List<global::RulesEngine.Models.Rule> ParseRulesFromConfig(string rulesConfig)
    {
        if (string.IsNullOrEmpty(rulesConfig))
        {
            return new List<global::RulesEngine.Models.Rule>();
        }

        var ruleConfigs = JsonSerializer.Deserialize<List<RuleConfig>>(rulesConfig);
        if (ruleConfigs == null)
        {
            return new List<global::RulesEngine.Models.Rule>();
        }

        return ruleConfigs.Select(rc => new global::RulesEngine.Models.Rule
        {
            RuleName = rc.RuleName,
            Expression = rc.Expression,
                RuleExpressionType = System.Enum.Parse<global::RulesEngine.Models.RuleExpressionType>(rc.RuleExpressionType),
                SuccessEvent = rc.SuccessEvent,
                ErrorMessage = rc.ErrorMessage
        }).ToList();
    }

    private class RuleConfig
    {
        public string RuleName { get; set; } = string.Empty;
        public string Expression { get; set; } = string.Empty;
        public string RuleExpressionType { get; set; } = "LambdaExpression";
        public string SuccessEvent { get; set; } = string.Empty;
        public string ErrorMessage { get; set; } = string.Empty;
        public string ErrorType { get; set; } = "Error";
    }
}
