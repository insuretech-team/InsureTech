namespace PoliSync.RulesEngine.Dtos;

public record CreateWorkflowRequest(
    string WorkflowName,
    string Description,
    List<RuleDto> Rules
);

public record UpdateWorkflowRequest(
    string WorkflowName,
    string Description,
    List<RuleDto> Rules,
    bool IsActive
);

public record RuleDto(
    string RuleName,
    string Expression,
    RuleExpressionTypeDto ExpressionType,
    string SuccessEvent,
    string ErrorMessage,
    ErrorTypeDto ErrorType,
    List<RuleDto>? ChildRules
);

public enum RuleExpressionTypeDto
{
    LambdaExpression = 0,
    CustomExpression = 1
}

public enum ErrorTypeDto
{
    Error = 0,
    Warning = 1
}

public record WorkflowDto(
    string Id,
    string WorkflowName,
    string Description,
    List<RuleDto> Rules,
    DateTime CreatedAt,
    DateTime UpdatedAt,
    string CreatedBy,
    bool IsActive
);

public record EvaluateRulesRequest(
    string WorkflowName,
    Dictionary<string, object> Inputs
);

public record RuleEvaluationResultDto(
    string RuleName,
    bool IsSuccess,
    string? SuccessEvent,
    string? ErrorMessage,
    List<string>? ChildResults
);

public record WorkflowEvaluationResultDto(
    string WorkflowName,
    bool IsSuccess,
    List<RuleEvaluationResultDto> Results,
    Dictionary<string, object>? Outputs
);
