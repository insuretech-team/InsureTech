namespace PoliSync.RulesEngine.Models;

public class Workflow
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    public string WorkflowName { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public List<Rule> Rules { get; set; } = [];
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;
    public string CreatedBy { get; set; } = string.Empty;
    public bool IsActive { get; set; } = true;
}

public class Rule
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    public string RuleName { get; set; } = string.Empty;
    public string Expression { get; set; } = string.Empty;
    public RuleExpressionType ExpressionType { get; set; } = RuleExpressionType.LambdaExpression;
    public string SuccessEvent { get; set; } = string.Empty;
    public string ErrorMessage { get; set; } = string.Empty;
    public ErrorType ErrorType { get; set; } = ErrorType.Error;
    public List<Rule>? Rules { get; set; }
    public Dictionary<string, object>? Properties { get; set; }
}

public enum RuleExpressionType
{
    LambdaExpression = 0,
    CustomExpression = 1
}

public enum ErrorType
{
    Error = 0,
    Warning = 1
}

public class RuleEvaluationRequest
{
    public string WorkflowName { get; set; } = string.Empty;
    public Dictionary<string, object> Inputs { get; set; } = [];
}

public class RuleEvaluationResult
{
    public string RuleName { get; set; } = string.Empty;
    public bool IsSuccess { get; set; }
    public string? SuccessEvent { get; set; }
    public string? ErrorMessage { get; set; }
    public List<string>? ChildResults { get; set; }
}

public class WorkflowEvaluationResult
{
    public string WorkflowName { get; set; } = string.Empty;
    public bool IsSuccess { get; set; }
    public List<RuleEvaluationResult> Results { get; set; } = [];
    public Dictionary<string, object>? Outputs { get; set; }
}
