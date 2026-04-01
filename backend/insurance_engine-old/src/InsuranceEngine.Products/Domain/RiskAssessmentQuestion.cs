using System;

namespace InsuranceEngine.Products.Domain;

public class RiskAssessmentQuestion
{
    public Guid Id { get; set; }
    public Guid ProductId { get; set; }
    public Product? Product { get; set; }
    
    public string QuestionText { get; set; } = string.Empty;
    public string? QuestionTextBn { get; set; }
    public string OptionsJson { get; set; } = "[]"; 
    public int Weight { get; set; }
}
