namespace PoliSync.Products.Domain;

/// <summary>
/// Dynamic pricing configuration for a product.
/// Maps to 'pricing_configs' table in insurance_schema.
/// Rules stored as JSONB.
/// </summary>
public class PricingConfig
{
    public Guid PricingConfigId { get; private set; }
    public Guid ProductId { get; private set; }
    public string Rules { get; private set; } = "[]";  // JSONB array of PricingRule
    public DateTime EffectiveFrom { get; private set; }
    public DateTime? EffectiveTo { get; private set; }
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }

    // Navigation
    public Product? Product { get; private set; }

    private PricingConfig() { }

    public static PricingConfig Create(
        Guid productId,
        string rules,
        DateTime effectiveFrom,
        DateTime? effectiveTo = null)
    {
        return new PricingConfig
        {
            PricingConfigId = Guid.NewGuid(),
            ProductId = productId,
            Rules = rules,
            EffectiveFrom = effectiveFrom,
            EffectiveTo = effectiveTo,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }
}
