using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain;

/// <summary>
/// PricingRule value object representing a single pricing rule with conditions and actions.
/// </summary>
public record PricingRule(
    Guid RuleId,
    string RuleName,
    string RuleType,
    List<RuleCondition> Conditions,
    RuleAction Action,
    bool ApplyAll,
    int Priority) : ValueObject
{
    protected override IEnumerable<object?> GetEqualityComponents()
    {
        yield return RuleId;
        yield return RuleName;
        yield return RuleType;
        yield return ApplyAll;
        yield return Priority;
        foreach (var condition in Conditions)
            yield return condition;
        yield return Action;
    }
}

/// <summary>
/// RuleCondition value object representing a condition to evaluate in a pricing rule.
/// </summary>
public record RuleCondition(
    string Field,
    string Operator,
    string Value) : ValueObject
{
    protected override IEnumerable<object?> GetEqualityComponents()
    {
        yield return Field;
        yield return Operator;
        yield return Value;
    }
}

/// <summary>
/// RuleAction value object representing the action to apply when a rule matches.
/// Supported ActionTypes: MULTIPLY, ADD, SET, DISCOUNT
/// </summary>
public record RuleAction(
    string ActionType,
    decimal Value) : ValueObject
{
    protected override IEnumerable<object?> GetEqualityComponents()
    {
        yield return ActionType;
        yield return Value;
    }
}
