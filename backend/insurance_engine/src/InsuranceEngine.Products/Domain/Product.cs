using InsuranceEngine.SharedKernel.Domain;
using Insuretech.Products.Entity.V1;

namespace InsuranceEngine.Products.Domain;

public sealed class Product : AggregateRoot<Guid>
{
    public string ProductCode { get; private set; } = string.Empty;
    public string ProductName { get; private set; } = string.Empty;
    public ProductCategory Category { get; private set; }
    public LocalizedString Description { get; private set; } = default!;
    public Money BasePremium { get; private set; } = default!;
    public Money MinSumInsured { get; private set; } = default!;
    public Money MaxSumInsured { get; private set; } = default!;
    public int MinTenureMonths { get; private set; }
    public int MaxTenureMonths { get; private set; }
    public ProductStatus Status { get; private set; }
    public List<string> AssessmentQuestions { get; private set; } = new();
    public DateTime CreatedAt { get; private set; }

    private Product(Guid id, string code, string name, ProductCategory category, LocalizedString description, 
                  Money basePremium, Money minSum, Money maxSum, int minTenure, int maxTenure)
    {
        Id = id;
        ProductCode = code;
        ProductName = name;
        Category = category;
        Description = description;
        BasePremium = basePremium;
        MinSumInsured = minSum;
        MaxSumInsured = maxSum;
        MinTenureMonths = minTenure;
        MaxTenureMonths = maxTenure;
        Status = ProductStatus.Active;
        CreatedAt = DateTime.UtcNow;
    }

    public static Product Create(string code, string name, ProductCategory category, string enDesc, string bnDesc, 
                               Money basePremium, Money minSum, Money maxSum, int minTenure, int maxTenure)
    {
        // Domain Validation
        if (string.IsNullOrWhiteSpace(code)) throw new ArgumentException("Product code is required");
        if (minSum.Amount > maxSum.Amount) throw new ArgumentException("Min Sum Insured cannot be greater than Max Sum Insured");

        return new Product(Guid.NewGuid(), code, name, category, new LocalizedString(enDesc, bnDesc), basePremium, minSum, maxSum, minTenure, maxTenure);
    }

    public Money CalculatePremium(int applicantAge, int tenureMonths)
    {
        // FR-024: Dynamic Premium Logic (Simple implementation for demonstration)
        decimal multiplier = 1.0m;
        
        if (applicantAge > 50) multiplier += 0.2m; // 20% increase for senior citizens
        if (tenureMonths > 12) multiplier -= 0.05m; // 5% discount for long term
        
        return Money.FromDecimal(BasePremium.ToDecimal() * multiplier);
    }

    public void SetAssessmentQuestions(IEnumerable<string> questions)
    {
        AssessmentQuestions = questions.ToList();
    }

    public void Deactivate()
    {
        Status = ProductStatus.Inactive;
    }
}
