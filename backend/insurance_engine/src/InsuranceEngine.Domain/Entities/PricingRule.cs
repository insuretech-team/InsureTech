using System;

namespace InsuranceEngine.Domain.Entities;

public class PricingRule
{
    public Guid Id { get; set; }
    public Guid ProductId { get; set; }
    public Product? Product { get; set; }
    
    public string RuleName { get; set; } = string.Empty;
    public string RuleExpression { get; set; } = string.Empty; // Logical expression
    public decimal AdjustmentAmount { get; set; }
    public bool IsPercentage { get; set; }
}
