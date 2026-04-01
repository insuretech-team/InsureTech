using System.Collections.Concurrent;
using Microsoft.Extensions.Logging;

namespace PoliSync.RulesEngine.Services;

public class InMemoryWorkflowRepository : IWorkflowRepository
{
    private readonly ConcurrentDictionary<string, Models.Workflow> _workflows = new();
    private readonly ILogger<InMemoryWorkflowRepository> _logger;

    public InMemoryWorkflowRepository(ILogger<InMemoryWorkflowRepository> logger)
    {
        _logger = logger;
        SeedDefaultWorkflows();
    }

    public Task<Models.Workflow?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _workflows.TryGetValue(id, out var workflow);
        return Task.FromResult(workflow);
    }

    public Task<Models.Workflow?> GetByNameAsync(string name, CancellationToken cancellationToken = default)
    {
        var workflow = _workflows.Values.FirstOrDefault(w => 
            w.WorkflowName.Equals(name, StringComparison.OrdinalIgnoreCase));
        return Task.FromResult(workflow);
    }

    public Task<IEnumerable<Models.Workflow>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_workflows.Values.AsEnumerable());
    }

    public Task<Models.Workflow> CreateAsync(Models.Workflow workflow, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(workflow.Id))
        {
            workflow.Id = Guid.NewGuid().ToString();
        }
        workflow.CreatedAt = DateTime.UtcNow;
        workflow.UpdatedAt = DateTime.UtcNow;
        
        _workflows[workflow.Id] = workflow;
        _logger.LogInformation("Created workflow: {WorkflowId} - {WorkflowName}", workflow.Id, workflow.WorkflowName);
        
        return Task.FromResult(workflow);
    }

    public Task<Models.Workflow?> UpdateAsync(Models.Workflow workflow, CancellationToken cancellationToken = default)
    {
        if (!_workflows.ContainsKey(workflow.Id))
        {
            return Task.FromResult<Models.Workflow?>(null);
        }

        workflow.UpdatedAt = DateTime.UtcNow;
        _workflows[workflow.Id] = workflow;
        _logger.LogInformation("Updated workflow: {WorkflowId}", workflow.Id);
        
        return Task.FromResult<Models.Workflow?>(workflow);
    }

    public Task<bool> DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var result = _workflows.TryRemove(id, out _);
        if (result)
        {
            _logger.LogInformation("Deleted workflow: {WorkflowId}", id);
        }
        return Task.FromResult(result);
    }

    private void SeedDefaultWorkflows()
    {
        // Underwriting Eligibility Workflow
        var underwritingWorkflow = new Models.Workflow
        {
            Id = Guid.NewGuid().ToString(),
            WorkflowName = "UnderwritingEligibility",
            Description = "Evaluates if an applicant is eligible for insurance",
            Rules =
            [
                new Models.Rule
                {
                    RuleName = "AgeCheck",
                    Expression = "input1.age >= 18 AND input1.age <= 65",
                    SuccessEvent = "Age eligible",
                    ErrorMessage = "Applicant age must be between 18 and 65"
                },
                new Models.Rule
                {
                    RuleName = "CreditScoreCheck",
                    Expression = "input1.creditScore >= 600",
                    SuccessEvent = "Credit score acceptable",
                    ErrorMessage = "Credit score below minimum threshold"
                },
                new Models.Rule
                {
                    RuleName = "NoFraudHistory",
                    Expression = "input1.hasFraudHistory == false",
                    SuccessEvent = "No fraud history",
                    ErrorMessage = "Applicant has fraud history"
                }
            ]
        };

        // Claims Approval Workflow
        var claimsWorkflow = new Models.Workflow
        {
            Id = Guid.NewGuid().ToString(),
            WorkflowName = "ClaimsApproval",
            Description = "Evaluates if a claim should be approved",
            Rules =
            [
                new Models.Rule
                {
                    RuleName = "PolicyActive",
                    Expression = "input1.isPolicyActive == true",
                    SuccessEvent = "Policy is active",
                    ErrorMessage = "Policy is not active"
                },
                new Models.Rule
                {
                    RuleName = "WithinCoverage",
                    Expression = "input1.claimAmount <= input1.coverageLimit",
                    SuccessEvent = "Claim within coverage",
                    ErrorMessage = "Claim exceeds coverage limit"
                },
                new Models.Rule
                {
                    RuleName = "WaitingPeriodPassed",
                    Expression = "input1.daysSincePolicyStart >= input1.waitingPeriodDays",
                    SuccessEvent = "Waiting period satisfied",
                    ErrorMessage = "Waiting period not yet satisfied"
                }
            ]
        };

        // Discount Calculation Workflow
        var discountWorkflow = new Models.Workflow
        {
            Id = Guid.NewGuid().ToString(),
            WorkflowName = "DiscountCalculation",
            Description = "Calculates applicable discounts",
            Rules =
            [
                new Models.Rule
                {
                    RuleName = "LoyaltyDiscount",
                    Expression = "input1.yearsAsCustomer >= 5",
                    SuccessEvent = "10",
                    ErrorMessage = "Not eligible for loyalty discount"
                },
                new Models.Rule
                {
                    RuleName = "MultiPolicyDiscount",
                    Expression = "input1.policyCount >= 2",
                    SuccessEvent = "15",
                    ErrorMessage = "Not eligible for multi-policy discount"
                },
                new Models.Rule
                {
                    RuleName = "SafeDriverDiscount",
                    Expression = "input1.yearsWithoutClaim >= 3",
                    SuccessEvent = "20",
                    ErrorMessage = "Not eligible for safe driver discount"
                }
            ]
        };

        _workflows[underwritingWorkflow.Id] = underwritingWorkflow;
        _workflows[claimsWorkflow.Id] = claimsWorkflow;
        _workflows[discountWorkflow.Id] = discountWorkflow;

        _logger.LogInformation("Seeded {Count} default workflows", _workflows.Count);
    }
}
