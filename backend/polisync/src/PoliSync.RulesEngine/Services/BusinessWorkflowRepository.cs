using System.Collections.Concurrent;
using System.Text.Json;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Workflow.Entity.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.RulesEngine.Services;

public class BusinessWorkflowRepository : IBusinessWorkflowRepository
{
    private readonly ConcurrentDictionary<string, BusinessWorkflowDefinition> _workflows = new();
    private readonly ConcurrentDictionary<string, BusinessWorkflowExecution> _executions = new();
    private readonly ILogger<BusinessWorkflowRepository> _logger;

    public BusinessWorkflowRepository(ILogger<BusinessWorkflowRepository> logger)
    {
        _logger = logger;
        SeedDefaultWorkflows();
    }

    public Task<BusinessWorkflowDefinition?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _workflows.TryGetValue(id, out var workflow);
        return Task.FromResult(workflow);
    }

    public Task<BusinessWorkflowDefinition?> GetByNameAsync(string name, CancellationToken cancellationToken = default)
    {
        var workflow = _workflows.Values.FirstOrDefault(w => 
            w.WorkflowName.Equals(name, StringComparison.OrdinalIgnoreCase));
        return Task.FromResult(workflow);
    }

    public Task<IEnumerable<BusinessWorkflowDefinition>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_workflows.Values.AsEnumerable());
    }

    public Task<BusinessWorkflowDefinition> CreateAsync(BusinessWorkflowDefinition workflow, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(workflow.BusinessWorkflowId))
        {
            workflow.BusinessWorkflowId = Guid.NewGuid().ToString();
        }
        workflow.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        workflow.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        workflow.Version = 1;
        
        _workflows[workflow.BusinessWorkflowId] = workflow;
        _logger.LogInformation("Created business workflow: {WorkflowId} - {WorkflowName}", 
            workflow.BusinessWorkflowId, workflow.WorkflowName);
        
        return Task.FromResult(workflow);
    }

    public Task<BusinessWorkflowDefinition?> UpdateAsync(BusinessWorkflowDefinition workflow, CancellationToken cancellationToken = default)
    {
        if (!_workflows.ContainsKey(workflow.BusinessWorkflowId))
        {
            return Task.FromResult<BusinessWorkflowDefinition?>(null);
        }

        workflow.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        workflow.Version++;
        _workflows[workflow.BusinessWorkflowId] = workflow;
        
        _logger.LogInformation("Updated business workflow: {WorkflowId} to version {Version}", 
            workflow.BusinessWorkflowId, workflow.Version);
        
        return Task.FromResult<BusinessWorkflowDefinition?>(workflow);
    }

    public Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default)
    {
        if (permanent)
        {
            var result = _workflows.TryRemove(id, out _);
            if (result)
            {
                _logger.LogInformation("Permanently deleted business workflow: {WorkflowId}", id);
            }
            return Task.FromResult(result);
        }
        else
        {
            if (_workflows.TryGetValue(id, out var workflow))
            {
                workflow.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                workflow.Status = BusinessWorkflowStatus.Inactive;
                _logger.LogInformation("Soft deleted business workflow: {WorkflowId}", id);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }
    }

    public Task<BusinessWorkflowExecution> LogExecutionAsync(BusinessWorkflowExecution execution, CancellationToken cancellationToken = default)
    {
        _executions[execution.ExecutionId] = execution;
        return Task.FromResult(execution);
    }

    public Task<IEnumerable<BusinessWorkflowExecution>> GetExecutionsByWorkflowIdAsync(string workflowId, CancellationToken cancellationToken = default)
    {
        var executions = _executions.Values
            .Where(e => e.BusinessWorkflowId == workflowId)
            .OrderByDescending(e => e.ExecutedAt)
            .AsEnumerable();
        return Task.FromResult(executions);
    }

    private void SeedDefaultWorkflows()
    {
        // Underwriting Eligibility Workflow
        var underwritingWorkflow = new BusinessWorkflowDefinition
        {
            BusinessWorkflowId = Guid.NewGuid().ToString(),
            WorkflowName = "UnderwritingEligibility",
            Description = "Evaluates if an applicant is eligible for insurance coverage",
            WorkflowType = BusinessWorkflowType.UnderwritingEligibility,
            Status = BusinessWorkflowStatus.Active,
            RulesConfig = JsonSerializer.Serialize(new[]
            {
                new
                {
                    RuleName = "AgeCheck",
                    Expression = "input1.age >= 18 AND input1.age <= 65",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "ELIGIBLE",
                    ErrorMessage = "Applicant age must be between 18 and 65",
                    ErrorType = "Error"
                },
                new
                {
                    RuleName = "CreditScoreCheck",
                    Expression = "input1.creditScore >= 600",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "CREDIT_OK",
                    ErrorMessage = "Credit score below minimum threshold of 600",
                    ErrorType = "Error"
                },
                new
                {
                    RuleName = "NoFraudHistory",
                    Expression = "input1.hasFraudHistory == false",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "NO_FRAUD",
                    ErrorMessage = "Applicant has fraud history",
                    ErrorType = "Error"
                }
            }),
            CreatedBy = "system"
        };

        // Claims Approval Workflow
        var claimsWorkflow = new BusinessWorkflowDefinition
        {
            BusinessWorkflowId = Guid.NewGuid().ToString(),
            WorkflowName = "ClaimsApproval",
            Description = "Evaluates if a claim should be approved",
            WorkflowType = BusinessWorkflowType.ClaimsApproval,
            Status = BusinessWorkflowStatus.Active,
            RulesConfig = JsonSerializer.Serialize(new[]
            {
                new
                {
                    RuleName = "PolicyActive",
                    Expression = "input1.isPolicyActive == true",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "POLICY_ACTIVE",
                    ErrorMessage = "Policy is not active",
                    ErrorType = "Error"
                },
                new
                {
                    RuleName = "WithinCoverage",
                    Expression = "input1.claimAmount <= input1.coverageLimit",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "WITHIN_COVERAGE",
                    ErrorMessage = "Claim exceeds coverage limit",
                    ErrorType = "Error"
                },
                new
                {
                    RuleName = "WaitingPeriodPassed",
                    Expression = "input1.daysSincePolicyStart >= input1.waitingPeriodDays",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "WAITING_PASSED",
                    ErrorMessage = "Waiting period not yet satisfied",
                    ErrorType = "Error"
                }
            }),
            CreatedBy = "system"
        };

        // Discount Calculation Workflow
        var discountWorkflow = new BusinessWorkflowDefinition
        {
            BusinessWorkflowId = Guid.NewGuid().ToString(),
            WorkflowName = "DiscountCalculation",
            Description = "Calculates applicable discounts based on customer profile",
            WorkflowType = BusinessWorkflowType.DiscountCalculation,
            Status = BusinessWorkflowStatus.Active,
            RulesConfig = JsonSerializer.Serialize(new[]
            {
                new
                {
                    RuleName = "LoyaltyDiscount",
                    Expression = "input1.yearsAsCustomer >= 5",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "10",
                    ErrorMessage = "Not eligible for loyalty discount",
                    ErrorType = "Warning"
                },
                new
                {
                    RuleName = "MultiPolicyDiscount",
                    Expression = "input1.policyCount >= 2",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "15",
                    ErrorMessage = "Not eligible for multi-policy discount",
                    ErrorType = "Warning"
                },
                new
                {
                    RuleName = "SafeDriverDiscount",
                    Expression = "input1.yearsWithoutClaim >= 3",
                    RuleExpressionType = "LambdaExpression",
                    SuccessEvent = "20",
                    ErrorMessage = "Not eligible for safe driver discount",
                    ErrorType = "Warning"
                }
            }),
            CreatedBy = "system"
        };

        _workflows[underwritingWorkflow.BusinessWorkflowId] = underwritingWorkflow;
        _workflows[claimsWorkflow.BusinessWorkflowId] = claimsWorkflow;
        _workflows[discountWorkflow.BusinessWorkflowId] = discountWorkflow;

        _logger.LogInformation("Seeded {Count} default business workflows", _workflows.Count);
    }
}
