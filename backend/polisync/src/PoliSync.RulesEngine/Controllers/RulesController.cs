using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using PoliSync.RulesEngine.Dtos;
using PoliSync.RulesEngine.Services;

namespace PoliSync.RulesEngine.Controllers;

[ApiController]
[Route("api/[controller]")]
public class RulesController : ControllerBase
{
    private readonly IRulesEngineService _rulesEngineService;
    private readonly ILogger<RulesController> _logger;

    public RulesController(IRulesEngineService rulesEngineService, ILogger<RulesController> logger)
    {
        _rulesEngineService = rulesEngineService;
        _logger = logger;
    }

    [HttpGet("workflows")]
    public async Task<ActionResult<IEnumerable<WorkflowDto>>> GetWorkflows(CancellationToken cancellationToken)
    {
        var workflows = await _rulesEngineService.GetWorkflowsAsync(cancellationToken);
        var dtos = workflows.Select(MapToDto);
        return Ok(dtos);
    }

    [HttpGet("workflows/{id}")]
    public async Task<ActionResult<WorkflowDto>> GetWorkflow(string id, CancellationToken cancellationToken)
    {
        var workflow = await _rulesEngineService.GetWorkflowAsync(id, cancellationToken);
        if (workflow == null)
        {
            return NotFound();
        }
        return Ok(MapToDto(workflow));
    }

    [HttpPost("workflows")]
    public async Task<ActionResult<WorkflowDto>> CreateWorkflow(
        [FromBody] CreateWorkflowRequest request, 
        CancellationToken cancellationToken)
    {
        var workflow = new Models.Workflow
        {
            WorkflowName = request.WorkflowName,
            Description = request.Description,
            Rules = request.Rules.Select(MapToModel).ToList()
        };

        var created = await _rulesEngineService.CreateWorkflowAsync(workflow, cancellationToken);
        return CreatedAtAction(nameof(GetWorkflow), new { id = created.Id }, MapToDto(created));
    }

    [HttpPut("workflows/{id}")]
    public async Task<ActionResult<WorkflowDto>> UpdateWorkflow(
        string id,
        [FromBody] UpdateWorkflowRequest request,
        CancellationToken cancellationToken)
    {
        var workflow = new Models.Workflow
        {
            WorkflowName = request.WorkflowName,
            Description = request.Description,
            Rules = request.Rules.Select(MapToModel).ToList(),
            IsActive = request.IsActive
        };

        var updated = await _rulesEngineService.UpdateWorkflowAsync(id, workflow, cancellationToken);
        if (updated == null)
        {
            return NotFound();
        }
        return Ok(MapToDto(updated));
    }

    [HttpDelete("workflows/{id}")]
    public async Task<IActionResult> DeleteWorkflow(string id, CancellationToken cancellationToken)
    {
        var result = await _rulesEngineService.DeleteWorkflowAsync(id, cancellationToken);
        if (!result)
        {
            return NotFound();
        }
        return NoContent();
    }

    [HttpPost("evaluate")]
    public async Task<ActionResult<WorkflowEvaluationResultDto>> EvaluateRules(
        [FromBody] EvaluateRulesRequest request,
        CancellationToken cancellationToken)
    {
        try
        {
            var result = await _rulesEngineService.EvaluateWorkflowAsync(
                request.WorkflowName, 
                request.Inputs, 
                cancellationToken);
            
            return Ok(MapToDto(result));
        }
        catch (InvalidOperationException ex)
        {
            _logger.LogWarning(ex, "Workflow not found: {WorkflowName}", request.WorkflowName);
            return NotFound(new { error = ex.Message });
        }
    }

    private static WorkflowDto MapToDto(Models.Workflow workflow)
    {
        return new WorkflowDto(
            workflow.Id,
            workflow.WorkflowName,
            workflow.Description,
            workflow.Rules.Select(MapToDto).ToList(),
            workflow.CreatedAt,
            workflow.UpdatedAt,
            workflow.CreatedBy,
            workflow.IsActive
        );
    }

    private static RuleDto MapToDto(Models.Rule rule)
    {
        return new RuleDto(
            rule.RuleName,
            rule.Expression,
            (RuleExpressionTypeDto)(int)rule.ExpressionType,
            rule.SuccessEvent,
            rule.ErrorMessage,
            (ErrorTypeDto)(int)rule.ErrorType,
            rule.Rules?.Select(MapToDto).ToList()
        );
    }

    private static Models.Rule MapToModel(RuleDto dto)
    {
        return new Models.Rule
        {
            RuleName = dto.RuleName,
            Expression = dto.Expression,
            ExpressionType = (Models.RuleExpressionType)(int)dto.ExpressionType,
            SuccessEvent = dto.SuccessEvent,
            ErrorMessage = dto.ErrorMessage,
            ErrorType = (Models.ErrorType)(int)dto.ErrorType,
            Rules = dto.ChildRules?.Select(MapToModel).ToList()
        };
    }

    private static WorkflowEvaluationResultDto MapToDto(Models.WorkflowEvaluationResult result)
    {
        return new WorkflowEvaluationResultDto(
            result.WorkflowName,
            result.IsSuccess,
            result.Results.Select(r => new RuleEvaluationResultDto(
                r.RuleName,
                r.IsSuccess,
                r.SuccessEvent,
                r.ErrorMessage,
                r.ChildResults
            )).ToList(),
            result.Outputs
        );
    }
}
