namespace PoliSync.Products.Domain;

/// <summary>
/// Product plan variant (e.g., Basic, Silver, Gold).
/// Maps to 'product_plans' table in insurance_schema.
/// </summary>
public class ProductPlan
{
    public Guid PlanId { get; private set; }
    public Guid ProductId { get; private set; }
    public string PlanName { get; private set; } = string.Empty;
    public string? PlanDescription { get; private set; }
    public long PremiumAmount { get; private set; }     // in paisa
    public long MinSumInsured { get; private set; }     // in paisa
    public long MaxSumInsured { get; private set; }     // in paisa
    public string? Attributes { get; private set; }     // JSONB
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }

    // Navigation
    public Product? Product { get; private set; }

    private ProductPlan() { }

    public static ProductPlan Create(
        Guid productId,
        string planName,
        long premiumAmount,
        long minSumInsured,
        long maxSumInsured,
        string? planDescription = null,
        string? attributes = null)
    {
        return new ProductPlan
        {
            PlanId = Guid.NewGuid(),
            ProductId = productId,
            PlanName = planName,
            PlanDescription = planDescription,
            PremiumAmount = premiumAmount,
            MinSumInsured = minSumInsured,
            MaxSumInsured = maxSumInsured,
            Attributes = attributes,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }
}
