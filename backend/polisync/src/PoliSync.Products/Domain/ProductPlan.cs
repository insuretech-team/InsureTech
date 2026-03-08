using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain;

/// <summary>
/// ProductPlan entity owned by Product. Represents a specific plan variant of a product.
/// </summary>
public class ProductPlan : Entity
{
    public Guid Id { get; private set; }
    public Guid ProductId { get; private set; }
    public string PlanCode { get; private set; } = string.Empty;
    public string PlanName { get; private set; } = string.Empty;
    public string Description { get; private set; } = string.Empty;
    public long BasePremiumPaisa { get; private set; }
    public long SumInsuredPaisa { get; private set; }
    public string Currency { get; private set; } = "BDT";
    public List<string> Features { get; private set; } = [];
    public bool IsActive { get; private set; }
    public DateTimeOffset CreatedAt { get; private set; }
    public DateTimeOffset UpdatedAt { get; private set; }

    // For EF Core
    protected ProductPlan() { }

    private ProductPlan(
        Guid id,
        Guid productId,
        string planCode,
        string planName,
        string description,
        long basePremiumPaisa,
        long sumInsuredPaisa,
        List<string> features,
        string currency)
    {
        Id = id;
        ProductId = productId;
        PlanCode = planCode;
        PlanName = planName;
        Description = description;
        BasePremiumPaisa = basePremiumPaisa;
        SumInsuredPaisa = sumInsuredPaisa;
        Features = features;
        Currency = currency;
        IsActive = true;
        CreatedAt = DateTimeOffset.UtcNow;
        UpdatedAt = DateTimeOffset.UtcNow;
    }

    /// <summary>
    /// Factory method to create a new ProductPlan with validation.
    /// </summary>
    public static Result<ProductPlan> Create(
        Guid productId,
        string planCode,
        string planName,
        string description,
        long basePremiumPaisa,
        long sumInsuredPaisa,
        List<string> features,
        string currency)
    {
        if (string.IsNullOrWhiteSpace(planCode))
            return Result<ProductPlan>.Fail("INVALID_CODE", "Plan code cannot be empty.");

        if (string.IsNullOrWhiteSpace(planName))
            return Result<ProductPlan>.Fail("INVALID_NAME", "Plan name cannot be empty.");

        if (basePremiumPaisa <= 0)
            return Result<ProductPlan>.Fail("INVALID_PREMIUM", "Plan premium must be greater than zero.");

        if (sumInsuredPaisa <= 0)
            return Result<ProductPlan>.Fail("INVALID_SUM_INSURED", "Plan sum insured must be greater than zero.");

        var plan = new ProductPlan(
            Guid.NewGuid(),
            productId,
            planCode,
            planName,
            description,
            basePremiumPaisa,
            sumInsuredPaisa,
            features,
            currency);

        return Result<ProductPlan>.Ok(plan);
    }
}
