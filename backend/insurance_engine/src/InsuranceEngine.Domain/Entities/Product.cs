using System;
using System.Collections.Generic;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Domain.Entities;

public class Product
{
    public Guid Id { get; set; }
    public string ProductCode { get; set; } = string.Empty;
    public string ProductName { get; set; } = string.Empty;
    public string? ProductNameBn { get; set; }
    public string? Description { get; set; }
    public string? DescriptionBn { get; set; }
    public ProductCategory Category { get; set; }
    public ProductStatus Status { get; set; }
    public decimal MinSumInsured { get; set; }
    public decimal MaxSumInsured { get; set; }
    public int MinAge { get; set; }
    public int MaxAge { get; set; }
    public int MinTenureMonths { get; set; }
    public int MaxTenureMonths { get; set; }
    
    public Guid InsurerId { get; set; }
    public Insurer? Insurer { get; set; }
    
    public List<ProductPlan> Plans { get; set; } = new();
    public List<RiskAssessmentQuestion> Questions { get; set; } = new();
    public List<PricingRule> PricingRules { get; set; } = new();
    
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
