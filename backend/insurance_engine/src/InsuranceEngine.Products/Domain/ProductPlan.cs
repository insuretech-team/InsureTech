using System;

namespace InsuranceEngine.Products.Domain;

public class ProductPlan
{
    public Guid Id { get; set; }
    public Guid ProductId { get; set; }
    public Product? Product { get; set; }
    
    public string PlanName { get; set; } = string.Empty;
    public string? PlanNameBn { get; set; }
    public string? Description { get; set; }
    public string? DescriptionBn { get; set; }
    
    public decimal PremiumAmount { get; set; }
    public decimal SumInsured { get; set; }
    
    public bool IsUnitWise { get; set; }
    public decimal? UnitPrice { get; set; }
}
