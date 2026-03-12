using System;
using System.Collections.Generic;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Domain;

/// <summary>
/// Product pricing configuration. Maps to 'pricing_configs' table.
/// </summary>
public class PricingConfig
{
    public Guid Id { get; set; }
    public Guid ProductId { get; set; }
    public Product? Product { get; set; }

    /// <summary>
    /// Pricing rules stored as JSONB
    /// </summary>
    public List<PricingRule> Rules { get; set; } = new();

    public DateTime EffectiveFrom { get; set; }
    public DateTime? EffectiveTo { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}

/// <summary>
/// Individual pricing rule within a PricingConfig.
/// Stored as JSONB array in pricing_configs.rules column.
/// </summary>
public class PricingRule
{
    public Guid Id { get; set; }
    public string RuleName { get; set; } = string.Empty;
    public RuleType Type { get; set; }
    public List<RuleCondition> Conditions { get; set; } = new();
    public RuleAction Action { get; set; } = new();
}

/// <summary>
/// Condition for evaluating a pricing rule
/// </summary>
public class RuleCondition
{
    public string Field { get; set; } = string.Empty;    // e.g., "age", "district"
    public string Operator { get; set; } = string.Empty; // e.g., ">=", "<=", "in"
    public string Value { get; set; } = string.Empty;    // Comparison value
}

/// <summary>
/// Action to take when rule conditions are met
/// </summary>
public class RuleAction
{
    public ActionType Type { get; set; }
    public double Value { get; set; } // Percentage or fixed amount
}
